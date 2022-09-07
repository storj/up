// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	pkg "storj.io/storj-up/pkg"
	"storj.io/storj/satellite/console/consolewasm"
)

var (
	satelliteHost, email, authService string
	export, write                     bool
	retry                             int
	s3                                bool
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
				if !export {
					fmt.Println("#Server is not yet available. Retry in 1 sec...", err)
				}
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
	credentialsCmd.PersistentFlags().BoolVarP(&write, "write", "w", false, "DEPRECATED. Write the right entries to rclone config file (storjdev, storj)")
	credentialsCmd.PersistentFlags().BoolVarP(&s3, "s3", "", false, "Generate S3 credentials. IMPORTANT: this command MUST be executed INSIDE containers as gateway will use it.")
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

	password := "123a123"

	if !export {
		fmt.Printf("User: %s\n", email)
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("ProjectID: %s\n", projectID)
		fmt.Printf("Cookie: _tokenKey=%s\n", cookie)
	} else {
		fmt.Printf("export STORJ_USER=%s\n", email)
		fmt.Printf("export STORJ_USER_PASSWORD=%s\n", password)
		fmt.Printf("export STORJ_PROJECT_ID=%s\n", projectID)
		fmt.Printf("export STORJ_SESSION_COOKIE=Cookie: _tokenKey=%s\n", cookie)
	}

	apiKey, err := console.CreateAPIKey(ctx, projectID)
	if err != nil {
		return errs.Wrap(err)
	}

	secret := "Welcome1"

	if !export {
		fmt.Printf("API key: %s\n", apiKey)
		fmt.Println()
	} else {
		fmt.Printf("export STORJ_API_KEY=%s\n", apiKey)
	}

	grant, err := consolewasm.GenAccessGrant(satelliteNodeURL, apiKey, secret, projectID)
	if err != nil {
		return errs.Wrap(err)
	}

	if !export {
		fmt.Printf("Encryption secret: %s \n", secret)
		fmt.Printf("Grant: %s\n", grant)
	} else {
		fmt.Printf("export STORJ_ENCRYPTION_SECRET=%s\n", secret)
		fmt.Printf("export STORJ_ACCESS=%s\n", grant)
		fmt.Printf("export UPLINK_ACCESS=%s\n", grant)
	}

	if s3 {
		if _, err := os.Stat("docker-compose.yaml"); err == nil {
			fmt.Println("Looks like you have a docker-compose.yaml. I suspect you execute this command from the host, not from the container. Please note that S3 compatible access grant should use the container network host (satellite-api). Therefore it should be executed from the container. (docker-compose exec satellite-api storj-up credentials -s3)")
		}
		accessKey, secretKey, endpoint, err := pkg.RegisterAccess(ctx, authService, grant)
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
			fmt.Println("Write flag is removed. Rclone config examples are printed out by default.")
		}
	}
	return err
}
