package cmd

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var entryPointCmd = &cobra.Command{
	Use:   "local-entrypoint [service or services]",
	Short: "Bind mount entrypoint.sh to use local modifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		return localEntrypoint(args)
	},
}

func init() {
	rootCmd.AddCommand(entryPointCmd)
}

func localEntrypoint(args []string) error {
	services, err := common.ResolveServices(args)
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
			currentComposeProject.Services[i].Volumes = append(currentComposeProject.Services[i].Volumes, common.CreateBind("./entrypoint.sh:/var/lib/storj/entrypoint.sh"))
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}