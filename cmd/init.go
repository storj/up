package cmd

import (
	_ "embed"
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

var initCmd = &cobra.Command{
	Use:   "init [service ...]",
	Short: "Creates/overwrites local docker-compose.yaml with service. You can use predefined groups as arguments.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := initCompose(TemplateFile, args)
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initCompose(templateDir string, services []string) (*types.Project, error) {
	templateComposeProject, err := common.CreateComposeProject(templateDir)
	if err != nil {
		return nil, err
	}

	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}
	servicesString := strings.Join(resolvedServices[:], ",")

	composeServices := templateComposeProject.AllServices()[:0]
	for _, service := range templateComposeProject.AllServices() {
		if strings.Contains(servicesString, strings.ReplaceAll(service.Name, "-", "")) {
			composeServices = append(composeServices, service)
		}
	}
	templateComposeProject.Services = composeServices

	return templateComposeProject, nil
}
