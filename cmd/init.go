package cmd

import (
	_ "embed"
	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"storj.io/storj-up/cmd/files/templates"
	"storj.io/storj-up/pkg/common"
	"strings"
)

func InitCmd() *cobra.Command {
	return &cobra.Command{
		Use: "init [selector]",
		Short: "Generate docker-compose file with selected services. " + selectorHelp + ". Without argument it generates " +
			"full Storj cluster with databases (storj,db)",
		RunE: func(cmd *cobra.Command, args []string) error {

			selector, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}

			composeProject, err := initCompose(templates.ComposeTemplate, selector)
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

	if len(services) == 0 {
		services = []string{"storj", "db"}
	}
	servicesString := strings.Join(common.ResolveServices(services)[:], ",")

	composeServices := templateComposeProject.AllServices()[:0]
	for _, service := range templateComposeProject.AllServices() {
		if strings.Contains(servicesString, service.Name) {
			composeServices = append(composeServices, service)
		}
	}
	templateComposeProject.Services = composeServices

	return templateComposeProject, nil
}
