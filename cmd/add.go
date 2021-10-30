package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/cmd/files/templates"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

func AddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [service ...]",
		Short: "add more services to the docker compose file. You can use predefined groups as arguments.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			templateProject, err := common.LoadComposeFromBytes(templates.ComposeTemplate)
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
	for _, service := range common.ResolveServices(services) {
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
