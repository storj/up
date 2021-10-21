package cmd

import (
	_ "embed"
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var initCmd = &cobra.Command{
	Use:   "init [groups/service]",
	Short: "Creates/overwrites local docker-compose.yaml with service. You can use predefined groups as arguments.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Init(args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func Init(requestedServices []string) error {
	templateComposeProject, err := common.CreateComposeProject("cmd/files/templates/docker-compose.template.yaml")
	if err != nil {
		return err
	}

	services, err := common.ResolveServices(requestedServices)
	if err != nil {
		return err
	}
	servicesString := strings.Join(services[:], ",")

	composeServices := templateComposeProject.AllServices()[:0]
	for _, service := range templateComposeProject.AllServices() {
		if strings.Contains(servicesString, strings.ReplaceAll(service.Name, "-", "")) {
			composeServices = append(composeServices, service)
		}
	}
	templateComposeProject.Services = composeServices

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: templateComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}
