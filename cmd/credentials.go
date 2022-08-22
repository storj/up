// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	pkg "storj.io/storj-up/pkg"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj/satellite/console/consolewasm"
)

var (
	satelliteHost, email, authService string
	export, write                     bool
	retry                             int
)

func credentialsCmd() *cobra.Command {
	credentialsCmd := &cobra.Command{
		Use:   "credentials",
		Short: "generate test user with credentialsCmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			for i := -1; i < retry; i++ {
				err = addCredentials(context.Background())
				if err == nil {
					return nil
				}
				fmt.Println("#Server is not yet available. Retry in 1 sec...", err)
				time.Sleep(1 * time.Second)
			}
			return err
		},
	}

	credentialsCmd.PersistentFlags().IntVarP(&retry, "retry", "r", 300, "Number of retry with 1 second interval. Default 300 = 5 minutes.")
	credentialsCmd.PersistentFlags().StringVarP(&email, "email", "m", "test@storj.io", "The email of the test user to use/create")
	credentialsCmd.PersistentFlags().StringVarP(&satelliteHost, "satellite", "s", "localhost", "The host of the satellite api to connect")
	credentialsCmd.PersistentFlags().StringVarP(&authService, "authservice", "a", "http://localhost:8888", "Host of the auth service")
	credentialsCmd.PersistentFlags().BoolVarP(&export, "export", "e", false, "Turn it off to get bash compatible output with export statements.")
	credentialsCmd.PersistentFlags().BoolVarP(&write, "write", "w", false, "Write the right entries to rclone config file (storjdev, storj")
	return credentialsCmd
}

func init() {
	rootCmd.AddCommand(credentialsCmd())
}

func addCredentials(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	satelliteNodeURL, err := pkg.GetSatelliteID(ctx, satelliteHost+":7777")
	if err != nil {
		return err
	}
	console := pkg.NewConsoleEndpoints(satelliteHost+":10000", email)

	err = console.Login(ctx)
	if err != nil {
		return err
	}
	projectID, cookie, err := console.GetOrCreateProject(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	if !export {
		fmt.Printf("User: %s\n", email)
		fmt.Printf("Password: %s\n", "123a123")
		fmt.Printf("ProjectID: %s\n", projectID)
		fmt.Printf("Cookie: _tokenKey=%s\n", cookie)
	} else {
		fmt.Printf("export STORJ_USER=%s\n", email)
		fmt.Printf("export STORJ_PROJECT_ID=%s\n", projectID)
		fmt.Printf("export STORJ_SESSION_COOKIE=Cookie: _tokenKey=%s\n", cookie)
	}

	apiKey, err := console.CreateAPIKey(ctx, projectID)
	if err != nil {
		return errs.Wrap(err)
	}

	secret := "Welcome1"

	internalSatelliteURL := strings.ReplaceAll(satelliteNodeURL, satelliteHost, "satellite-api")
	internalGrant, err := consolewasm.GenAccessGrant(internalSatelliteURL, apiKey, secret, projectID)
	if err != nil {
		return errs.Wrap(err)
	}

	if !export {
		fmt.Printf("API key: %s\n", apiKey)
		fmt.Println()

		fmt.Println("[internal access from containers]")
		fmt.Printf("Encryption secret: %s \n", secret)
		fmt.Printf("Grant: %s\n", internalGrant)
		fmt.Println()

	}

	grant, err := consolewasm.GenAccessGrant(satelliteNodeURL, apiKey, secret, projectID)
	if err != nil {
		return errs.Wrap(err)
	}

	if !export {
		fmt.Println("\n[from host]")
		fmt.Printf("Encryption secret: %s \n", secret)
		fmt.Printf("Grant: %s\n", grant)
	} else {
		fmt.Printf("export STORJ_ACCESS=%s\n", grant)
		fmt.Printf("export UPLINK_ACCESS=%s\n", grant)
	}

	composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
	if err != nil {
		return err
	}

	if containsService(composeProject.Services, "linksharing") {
		accessKey, secretKey, endpoint, err := pkg.RegisterAccess(ctx, authService, internalGrant)
		if err != nil {
			return errs.Wrap(err)
		}
		if !export {
			fmt.Printf("Access key: %s\n", accessKey)
			fmt.Printf("Secret key: %s\n", secretKey)
			fmt.Printf("Endpoint: %s\n", endpoint)
		} else {
			fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", accessKey)
			fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", secretKey)
			fmt.Printf("export STORJ_GATEWAY=%s\n", endpoint)
		}
		if write {
			err = updateRclone(accessKey, secretKey, endpoint, grant)
			if err != nil {
				return errs.Wrap(err)
			}
		}
	}
	return err
}

func containsService(services types.Services, s string) bool {
	for _, service := range services {
		if service.Name == s {
			return true
		}
	}
	return false
}

func updateRclone(key string, secret string, endpoint string, grant string) (err error) {
	usr, err := user.Current()
	if err != nil {
		return errs.Wrap(err)
	}

	out := strings.Builder{}
	rcloneConf := path.Join(usr.HomeDir, ".config", "rclone", "rclone.conf")

	var content []byte

	_ = os.MkdirAll(path.Dir(rcloneConf), 0o755)
	if _, err := os.Stat(rcloneConf); err == nil {
		content, err = ioutil.ReadFile(rcloneConf)
		if err != nil {
			return errs.Wrap(err)
		}
	} else if !os.IsNotExist(err) {
		return errs.Wrap(err)
	}

	section := regexp.MustCompile(`\[(.*)]`)
	currentSection := ""
	updatedS3 := false
	updatedNative := false
	for _, line := range strings.Split(string(content), "\n") {

		matches := section.FindStringSubmatch(line)
		if len(matches) > 0 {
			currentSection = matches[0]
		}

		if currentSection == "[storjdev]" {
			updatedNative = true
			if strings.HasPrefix(line, "access_grant") {
				out.WriteString("access_grant = " + secret + "\n")
				continue
			}
		}
		if currentSection == "[storjdevs3]" {
			updatedS3 = true
			if strings.HasPrefix(line, "secret_access_key") {
				out.WriteString("secret_access_key = " + secret + "\n")
				continue
			} else if strings.HasPrefix(line, "access_key_id") {
				out.WriteString("access_key_id = " + key + "\n")
				continue
			}
		}
		out.WriteString(line + "\n")
	}
	if !updatedS3 {
		out.WriteString("\n[storjdevs3]\n")
		out.WriteString("type = s3\n")
		out.WriteString("provider = Other \n")
		out.WriteString("access_key_id = " + key + "\n")
		out.WriteString("secret_access_key = " + secret + "\n")
		out.WriteString("endpoint = " + endpoint + "\n")
	}
	if !updatedNative {
		out.WriteString("\n[storjdev]\n")
		out.WriteString("type = tardigrade\n")
		out.WriteString("access_grant = " + grant + "\n")
	}
	err = ioutil.WriteFile(rcloneConf, []byte(out.String()), 0o644)
	return err
}
