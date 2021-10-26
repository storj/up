package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

func AddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [service ...]",
		Short: "add more services to the docker compose file",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadCompose(ComposeFile)
			templateProject, err := common.LoadCompose(TemplateFile)
			updatedComposeProject, err := AddToCompose(composeProject, templateProject, args)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(AddCmd())
}

func AddToCompose(compose *types.Project, template *types.Project, services []string) (*types.Project, error) {
	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}

	for _, service := range resolvedServices {
		if !common.ContainsService(compose.Services, service) {
			newService, err := template.GetService(service)
			if err != nil {
				return nil, err
			}
			compose.Services = append(compose.Services, newService)
		}
	}
	return compose, nil
}