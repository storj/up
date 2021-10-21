package cmd

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var setEnvCmd = &cobra.Command{
	Use:   "setEnv [KEY=VALUE] [service or services]",
	Short: "Set environment variable / parameter in a container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return SetEnv(args)
	},
}

var unsetEnvCmd = &cobra.Command{
	Use:   "unsetEnv [KEY] [service or services]",
	Short: "Remove environment variable / parameter in a container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return UnsetEnv(args)
	},
}

func init() {
	rootCmd.AddCommand(setEnvCmd)
	rootCmd.AddCommand(unsetEnvCmd)
}

func SetEnv(args []string) error {
	parts := strings.SplitN(args[0], "=", 2)

	services, err := common.ResolveServices(args[1:])
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
			currentComposeProject.Services[i].Environment[parts[0]] = &parts[1]
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}

func UnsetEnv(args []string) error {
	services, err := common.ResolveServices(args[1:])
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
			delete(currentComposeProject.Services[i].Environment, args[0])
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}