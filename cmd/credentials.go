// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
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
	password = "123a123"
	secret   = "Welcome1"
	filename = ".creds"
)

var (
	retry   int
	export  bool
	s3      bool
	persist bool

	satellite   string
	console     string
	authservice string

	credentials Credentials
)

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
		Short: "generate test user with credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
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
					err = executeWithRetry(context.Background(), generateS3Credentials)
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
	pflags.StringVarP(&satellite, "satellite", "s", "localhost:7777", "The host and port of of the satellite api to connect. Defaults to localhost or STORJ_DOCKER_HOST if set.")
	pflags.StringVarP(&console, "console", "c", "localhost:10000", "The host and port of of the satellite api console to connect. Defaults to localhost or STORJ_DOCKER_HOST if set.")
	pflags.StringVarP(&authservice, "authservice", "a", "http://localhost:8888", "Host of the auth service. Defaults to localhost or STORJ_DOCKER_HOST if set.")
	pflags.BoolVarP(&export, "export", "e", false, "Turn it off to get bash compatible output with export statements.")
	pflags.BoolVarP(&s3, "s3", "", false, "Generate S3 credentials. IMPORTANT: this command MUST be executed INSIDE containers as gateway will use it.")
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
	for i := 0; i < retry; i++ {
		err := f(ctx)
		if err == nil {
			return nil
		}
		if !export {
			fmt.Println("#Server is not yet available. Retry in 1 sec...", err)
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("Failed after " + strconv.Itoa(retry) + " retries")
}

func attemptUpdateDockerHost() error {
	dockerHost := os.Getenv("STORJ_DOCKER_HOST")
	if dockerHost != "" {
		satelliteUrl, err := url.Parse("http://" + satellite)
		if err != nil {
			return err
		}
		consoleUrl, err := url.Parse("http://" + console)
		if err != nil {
			return err
		}
		authUrl, err := url.Parse(authservice)
		if err != nil {
			return err
		}
		satellite = strings.Replace(satellite, satelliteUrl.Hostname(), dockerHost, 1)
		console = strings.Replace(console, consoleUrl.Hostname(), dockerHost, 1)
		authservice = strings.Replace(authservice, authUrl.Hostname(), dockerHost, 1)
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

	satelliteNodeURL, err := pkg.GetSatelliteID(ctx, satellite)
	if err != nil {
		return errs.Wrap(err)
	}

	consoleEndpoint := pkg.NewConsoleEndpoints(console, credentials.StorjUser)
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

	credentials.Grant, err = consolewasm.GenAccessGrant(satelliteNodeURL+"@"+satellite, credentials.ApiKey, secret, base64.StdEncoding.EncodeToString(projectUUID.Bytes()))
	if err != nil {
		return errs.Wrap(err)
	}

	if s3 {
		err = generateS3Credentials(ctx)
		if err != nil {
			return err
		}
	}

	return err
}

func generateS3Credentials(ctx context.Context) error {
	if _, err := os.Stat("docker-compose.yaml"); err == nil {
		fmt.Println("Looks like you have a docker-compose.yaml. I suspect you execute this command from the host, not from the container. Please note that S3 compatible access Grant should use the container network host (satellite-api). Therefore it should be executed from the container. (docker-compose exec satellite-api storj-up credentials -s3)")
	}
	var err error
	credentials.AccessKey, credentials.SecretKey, credentials.Endpoint, err = pkg.RegisterAccess(ctx, authservice, credentials.Grant)
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

		if s3 {
			fmt.Printf("Access key: %s\n", credentials.AccessKey)
			fmt.Printf("Secret key: %s\n", credentials.SecretKey)
			fmt.Printf("Endpoint: %s\n", credentials.Endpoint)
		}
	} else {
		fmt.Printf("export STORJ_USER=%s\n", credentials.StorjUser)
		fmt.Printf("export STORJ_USER_PASSWORD=%s\n", password)
		fmt.Printf("export STORJ_PROJECT_ID=%s\n", credentials.ProjectID)
		fmt.Printf("export STORJ_SESSION_COOKIE=Cookie: _tokenKey=%s\n", credentials.Cookie)

		fmt.Printf("export STORJ_API_KEY=%s\n", credentials.ApiKey)

		fmt.Printf("export STORJ_ENCRYPTION_SECRET=%s\n", secret)
		fmt.Printf("export STORJ_ACCESS=%s\n", credentials.Grant)
		fmt.Printf("export UPLINK_ACCESS=%s\n", credentials.Grant)

		if s3 {
			fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", credentials.AccessKey)
			fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", credentials.SecretKey)
			fmt.Printf("export STORJ_GATEWAY=%s\n", credentials.Endpoint)
		}
	}
	return nil
}
