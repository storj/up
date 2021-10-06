package sjr

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
	"io/ioutil"
	"os/user"
	"path"
	"regexp"
	"storj.io/storj/satellite/console/consolewasm"
	"strings"
	"time"
)

func init() {
	var satelliteHost string
	var authService string
	var email string
	var export, write bool
	credentialsCmd := &cobra.Command{
		Use:   "credentials",
		Short: "Generate test user with credentialsCmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			sateliteNodeUrl, err := GetSatelliteId(ctx, satelliteHost+":7777")
			if err != nil {
				return err
			}
			console := NewConsoleEndpoints(satelliteHost+":10000", email)

			err = console.Login(ctx)
			if err != nil {
				return err
			}
			projectID, err := console.GetOrCreateProject(ctx)
			if err != nil {
				return errs.Wrap(err)
			}
			fmt.Printf("User: %s\n", email)
			fmt.Printf("ProjectID: %s\n", projectID)

			apiKey, err := console.CreateAPIKey(ctx, projectID)
			if err != nil {
				return errs.Wrap(err)
			}

			fmt.Printf("API key: %s (use from host)\n", apiKey)
			grant, err := consolewasm.GenAccessGrant(sateliteNodeUrl, apiKey, "Welcome1", projectID)
			if err != nil {
				return errs.Wrap(err)
			}
			fmt.Printf("Grant: %s\n", grant)

			fmt.Printf("API key: %s (use from container)\n", apiKey)
			internalSatelliteUrl := strings.ReplaceAll(sateliteNodeUrl, satelliteHost, "satellite-api")
			internalGrant, err := consolewasm.GenAccessGrant(internalSatelliteUrl, apiKey, "Welcome1", projectID)
			if err != nil {
				return errs.Wrap(err)
			}
			fmt.Printf("Grant: %s\n", internalGrant)

			accessKey, secretKey, endpoint, err := RegisterAccess(ctx, authService, internalGrant)
			if err != nil {
				return errs.Wrap(err)
			}
			fmt.Printf("Access key: %s\n", accessKey)
			fmt.Printf("Secret key: %s\n", secretKey)
			fmt.Printf("Endpoint: %s\n", endpoint)
			if write {
				err = updateRclone(accessKey, secretKey, endpoint, grant)
				if err != nil {
					return errs.Wrap(err)
				}
			}
			return err
		},
	}
	credentialsCmd.Flags().StringVarP(&email, "email", "m", "test@storj.io", "The email of the test user to use/create")
	credentialsCmd.Flags().StringVarP(&satelliteHost, "satellite", "s", "localhost", "The host of the satellite api to connect")
	credentialsCmd.Flags().StringVarP(&authService, "authservice", "a", "http://localhost:8888", "Host of the auth service")
	credentialsCmd.Flags().BoolVarP(&export, "export", "e", false, "Turn it off to get bash compatible output with export statements.")
	credentialsCmd.Flags().BoolVarP(&export, "write", "w", false, "Write the right entries to rclone config file (storjdev, storj")
	RootCmd.AddCommand(credentialsCmd)
}

func updateRclone(key string, secret string, endpoint string, grant string) error {
	usr, err := user.Current()
	if err != nil {
		return errs.Wrap(err)
	}

	out := strings.Builder{}
	rcloneConf := path.Join(usr.HomeDir, ".config", "rclone", "rclone.conf")
	content, err := ioutil.ReadFile(rcloneConf)
	if err != nil {
		return errs.Wrap(err)
	}

	section := regexp.MustCompile("\\[(.*)\\]")
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
	err = ioutil.WriteFile(rcloneConf, []byte(out.String()), 0644)
	return err
}
