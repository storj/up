package cmd

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var addCmd = &cobra.Command{
	Use:   "add [service names or groups]",
	Short: "Add more services to the docker-compose.yaml. You can use predefined groups as arguments.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Add(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func Add(requestedServices []string) error {
	templateComposeProject, err := common.CreateComposeProject("cmd/files/templates/docker-compose.template.yaml")
	if err != nil {
		return err
	}

	currentComposeProject, err := common.CreateComposeProject("docker-compose.yaml")
	if err != nil {
		return err
	}

	services, err := common.ResolveServices(requestedServices)
	if err != nil {
		return err
	}

	for _, service := range services {
		if !common.ContainsService(currentComposeProject.Services, service) {
			newService, err := templateComposeProject.GetService(service)
			if err != nil {
				return err
			}
			currentComposeProject.Services = append(currentComposeProject.Services, newService)
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}