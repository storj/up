package cmd

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var imageCmd = &cobra.Command{
	Use:   "image <image> [service or services]",
	Short: "use a prebuilt docker image",
	RunE: func(cmd *cobra.Command, args []string) error {
		return UpdateImage(args)
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
}

func UpdateImage(args []string) error {
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
			currentComposeService.Image = args[0]
			currentComposeProject.Services[i].Image = currentComposeService.Image
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}