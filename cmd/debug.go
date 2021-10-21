package cmd

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var debugCmd = &cobra.Command{
	Use:   "debug [service or services]",
	Short: "Turn on local debugging with DLV",
	RunE: func(cmd *cobra.Command, args []string) error {
		return DebugEnable(args)
	},
}
var noDebugCmd = &cobra.Command{
	Use:   "no-debug [service or services]",
	Short: "Turn off local debugging with DLV",
	RunE: func(cmd *cobra.Command, args []string) error {
		return DebugDisable(args)
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(noDebugCmd)
}

func DebugEnable(args []string) error {
	value := "true"
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
			currentComposeProject.Services[i].Environment["GO_DLV"] = &value
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}

func DebugDisable(args []string) error {
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
			delete(currentComposeProject.Services[i].Environment, "GO_DLV")
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}