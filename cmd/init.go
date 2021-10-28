package cmd

import (
	_ "embed"
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/cmd/files/templates"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [service ...]",
		Short: "creates/overwrites local docker compose file with services",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := initCompose(templates.ComposeTemplate, args)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(composeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(InitCmd())
}

func initCompose(templateBytes []byte, services []string) (*types.Project, error) {
	templateComposeProject, err := common.LoadComposeFromBytes(templateBytes)
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
		if strings.Contains(servicesString, service.Name) {
			composeServices = append(composeServices, service)
		}
	}
	templateComposeProject.Services = composeServices

	return templateComposeProject, nil
}
