// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zeebo/errs/v2"

	"storj.io/common/uuid"
	pkg "storj.io/storj-up/pkg"
	"storj.io/storj/satellite/console/consolewasm"
)

const (
	encKeyVersionByte = byte(77) // magic number EncryptionKey encoding
	secKeyVersionByte = byte(78) // magic number SecretKey encoding

	password = "password"
	secret   = "Welcome1"
	filename = ".creds"
)

var (
	base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)

	retry   int
	export  bool
	s3      bool
	persist bool

	satelliteHost   string
	consoleHost     string
	authServiceHost string

	credentials Credentials
)

// EncryptionKey is an encryption key that an access/secret are encrypted with.
type EncryptionKey [16]byte

// SecretKey is the secret key used to sign requests.
type SecretKey [32]byte

// Credentials is the structure of the credentials file.
type Credentials struct {
	StorjUser     string `json:"email"`
	StorjPassword string `json:"password"`

	ProjectID        string `json:"ProjectID"`
	Cookie           string `json:"Cookie"`
	ApiKey           string `json:"ApiKey"`
	EncryptionSecret string `json:"EncryptionSecret"`
	Grant            string `json:"Grant"`

	AccessKey string `json:"AccessKey,omitempty"`
	SecretKey string `json:"SecretKey,omitempty"`
	Endpoint  string `json:"Endpoint,omitempty"`
}

func credentialsCmd() *cobra.Command {
	credentialsCmd := &cobra.Command{
		Use:   "credentials",
		Args:  cobra.NoArgs,
		Short: "generate test user with credentials",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := os.Stat(filename)
			if err != nil || persist {
				err = executeWithRetry(context.Background(), generateCredentials)
				if err != nil {
					return err
				}
			} else {
				err = loadCredentials()
				if err != nil {
					return err
				}
				if s3 && (credentials.AccessKey == "" || credentials.SecretKey == "" || credentials.Endpoint == "") {
					err = attemptUpdateDockerHost()
					if err != nil {
						return err
					}
					err = executeWithRetry(context.Background(), registerS3Credentials)
					if err != nil {
						return err
					}
				}
			}
			err = printCredentials()
			if err != nil {
				return err
			}
			if persist {
				err = persistCredentials()
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	pflags := credentialsCmd.PersistentFlags()
	pflags.IntVarP(&retry, "retry", "r", 300, "Number of retry with 1 second interval. Default 300 = 5 minutes.")
	pflags.StringVarP(&credentials.StorjUser, "email", "m", "test@storj.io", "The email of the test user to use/create")
	pflags.StringVarP(&satelliteHost, "satellite", "s", "localhost:7777", "The host and port of of the satellite api to connect. Defaults to localhost or STORJ_DOCKER_HOST if set.")
	pflags.StringVarP(&consoleHost, "console", "c", "localhost:10000", "The host and port of of the satellite api console to connect. Defaults to localhost or STORJ_DOCKER_HOST if set.")
	pflags.StringVarP(&authServiceHost, "authservice", "a", "http://localhost:8888", "Host of the auth service. Defaults to localhost or STORJ_DOCKER_HOST if set.")
	pflags.BoolVarP(&export, "export", "e", false, "Turn it off to get bash compatible output with export statements.")
	pflags.BoolVarP(&s3, "s3", "", false, "Register S3 credentials with authservice. IMPORTANT: Proper registration requires this command to be executed INSIDE containers.")
	pflags.BoolVarP(&persist, "persist", "p", false, "Persist credentials to disk for reuse. If persisted credentials are found, they are returned instead of regenerating, however repeated calls with persist flag will regenerate and persist new credentials.")
	pflags.VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(flag.Name, flag)
	})

	return credentialsCmd
}

func init() {
	RootCmd.AddCommand(credentialsCmd())
}

func executeWithRetry(ctx context.Context, f func(ctx context.Context) error) error {
	var err error
	for i := 0; i < retry; i++ {
		err = f(ctx)
		if err == nil {
			return nil
		}
		if !export {
			fmt.Println("#Server is not yet available. Retry in 1 sec...", err)
		}
		time.Sleep(1 * time.Second)
	}
	return errs.Errorf("Failed after %v retries. Last error: %w", retry, err)
}

func attemptUpdateDockerHost() error {
	dockerHost := os.Getenv("STORJ_DOCKER_HOST")
	if dockerHost != "" {
		satelliteUrl, err := url.Parse("http://" + satelliteHost)
		if err != nil {
			return err
		}
		consoleUrl, err := url.Parse("http://" + consoleHost)
		if err != nil {
			return err
		}
		authUrl, err := url.Parse(authServiceHost)
		if err != nil {
			return err
		}
		satelliteHost = strings.Replace(satelliteHost, satelliteUrl.Hostname(), dockerHost, 1)
		consoleHost = strings.Replace(consoleHost, consoleUrl.Hostname(), dockerHost, 1)
		authServiceHost = strings.Replace(authServiceHost, authUrl.Hostname(), dockerHost, 1)
	}
	return nil
}

func generateCredentials(ctx context.Context) error {
	err := attemptUpdateDockerHost()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	satelliteNodeURL, err := pkg.GetSatelliteID(ctx, satelliteHost)
	if err != nil {
		return errs.Wrap(err)
	}

	consoleEndpoint := pkg.NewConsoleEndpoints(consoleHost, credentials.StorjUser)
	err = consoleEndpoint.Login(ctx)
	if err != nil {
		return errs.Wrap(err)
	}

	credentials.ProjectID, credentials.Cookie, err = consoleEndpoint.GetOrCreateProject(ctx)
	if err != nil {
		return errs.Wrap(err)
	}

	credentials.ApiKey, err = consoleEndpoint.CreateAPIKey(ctx, credentials.ProjectID)
	if err != nil {
		return errs.Wrap(err)
	}

	projectUUID, err := uuid.FromString(credentials.ProjectID)
	if err != nil {
		return errs.Wrap(err)
	}

	credentials.Grant, err = consolewasm.GenAccessGrant(satelliteNodeURL+"@"+satelliteHost, credentials.ApiKey, secret, base64.StdEncoding.EncodeToString(projectUUID.Bytes()), true)
	if err != nil {
		return errs.Wrap(err)
	}

	if s3 {
		err = registerS3Credentials(ctx)
		if err != nil {
			return errs.Wrap(err)
		}
	} else {
		credentials.AccessKey = newEncryptionKey()
		credentials.SecretKey = newSecretKey()
		credentials.Endpoint = ""
	}

	return err
}

func registerS3Credentials(ctx context.Context) error {
	if _, err := os.Stat("docker-compose.yaml"); err == nil {
		fmt.Println("Looks like you have a docker-compose.yaml. I suspect you execute this command from the host, not from the container. Please note that S3 compatible access Grant should use the container network host (satellite-api). Therefore it should be executed from the container. (docker-compose exec satellite-api storj-up credentials -s3)")
	}
	var err error
	credentials.AccessKey, credentials.SecretKey, credentials.Endpoint, err = pkg.RegisterAccess(ctx, authServiceHost, credentials.Grant)
	if err != nil {
		return errs.Wrap(err)
	}
	return err
}

func persistCredentials() error {
	credentials.StorjPassword = password
	credentials.EncryptionSecret = secret
	file, err := json.MarshalIndent(&credentials, "", "  ")
	if err != nil {
		return errs.Wrap(err)
	}
	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func loadCredentials() error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return errs.Wrap(err)
	}
	err = json.Unmarshal(file, &credentials)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func printCredentials() error {
	if !export {
		fmt.Printf("User: %s\n", credentials.StorjUser)
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("ProjectID: %s\n", credentials.ProjectID)
		fmt.Printf("Cookie: _tokenKey=%s\n", credentials.Cookie)
		fmt.Printf("API key: %s\n", credentials.ApiKey)
		fmt.Println()

		fmt.Printf("Encryption secret: %s \n", secret)
		fmt.Printf("Grant: %s\n", credentials.Grant)
		fmt.Println()

		fmt.Printf("Access key: %s\n", credentials.AccessKey)
		fmt.Printf("Secret key: %s\n", credentials.SecretKey)
		fmt.Printf("Endpoint: %s\n", credentials.Endpoint)
	} else {
		fmt.Printf("export STORJ_USER=%s\n", credentials.StorjUser)
		fmt.Printf("export STORJ_USER_PASSWORD=%s\n", password)
		fmt.Printf("export STORJ_PROJECT_ID=%s\n", credentials.ProjectID)
		fmt.Printf("export STORJ_SESSION_COOKIE=Cookie: _tokenKey=%s\n", credentials.Cookie)
		fmt.Printf("export STORJ_API_KEY=%s\n", credentials.ApiKey)

		fmt.Printf("export STORJ_ENCRYPTION_SECRET=%s\n", secret)
		fmt.Printf("export STORJ_ACCESS=%s\n", credentials.Grant)
		fmt.Printf("export UPLINK_ACCESS=%s\n", credentials.Grant)

		fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", credentials.AccessKey)
		fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", credentials.SecretKey)
		fmt.Printf("export STORJ_GATEWAY=%s\n", credentials.Endpoint)
	}
	return nil
}

func newEncryptionKey() string {
	key := EncryptionKey{encKeyVersionByte}
	if _, err := rand.Read(key[:]); err != nil {
		return ""
	}
	return strings.ToLower(
		base32Encoding.EncodeToString(
			append([]byte{encKeyVersionByte}, key[:]...),
		),
	)
}

func newSecretKey() string {
	secretKey := SecretKey{secKeyVersionByte}
	if _, err := rand.Read(secretKey[:]); err != nil {
		return ""
	}
	return strings.ToLower(
		base32Encoding.EncodeToString(
			append([]byte{secKeyVersionByte}, secretKey[:]...),
		),
	)
}
