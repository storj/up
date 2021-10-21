package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [service ...]",
	Short: "Add more services to the docker-compose.yaml. You can use predefined groups as arguments.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := AddToCompose(ComposeFile, TemplateFile, args)
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func AddToCompose(composeDir string, templateDir string, services []string) (*types.Project, error) {
	templateComposeProject, err := common.CreateComposeProject(templateDir)
	if err != nil {
		return nil, err
	}

	currentComposeProject, err := common.CreateComposeProject(composeDir)
	if err != nil {
		return nil, err
	}

	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}

	for _, service := range resolvedServices {
		if !common.ContainsService(currentComposeProject.Services, service) {
			newService, err := templateComposeProject.GetService(service)
			if err != nil {
				return nil, err
			}
			currentComposeProject.Services = append(currentComposeProject.Services, newService)
		}
	}
	return currentComposeProject, nil
}