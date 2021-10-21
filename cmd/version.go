package cmd

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var versionCmd = &cobra.Command{
	Use:   "version <version> [service or services]",
	Short: "Set version (docker image tag) for specified services",
	RunE: func(cmd *cobra.Command, args []string) error {
		return UpdateVersion(args)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func UpdateVersion(args []string) error {
	services, err := common.ResolveServices(args[1:])
	if err != nil {
		return err
	}
	servicesString := strings.Join(services[:], ",")

	currentComposeProject, err := common.CreateComposeProject("docker-compose.yaml")
	if err != nil {
		return err
	}

	currentComposeServices := currentComposeProject.AllServices()

	for i, currentComposeService := range currentComposeServices {
		if strings.Contains(servicesString, strings.ReplaceAll(currentComposeService.Name, "-", "")) {
			currentComposeService.Image = strings.ReplaceAll(currentComposeService.Image, "@sha256", "")
			currentComposeService.Image = strings.Split(currentComposeService.Image, ":")[0] + ":" + args[0]
			currentComposeProject.Services[i].Image = currentComposeService.Image
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}