package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
	"io/ioutil"
	"strconv"
	"strings"
)

var scaleCmd = &cobra.Command{
	Use:   "scale <number> [service or services]",
	Short: "Static scale of service or services",
	Long: "This command creates multiple instances of the service or services. After this scale services couldn't be scaled up with `docker-compose scale any more`. " +
		"But also not required to scale up and down and it's possible to do per instance local bindmount",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Scale(args)
	},
}

func init() {
	rootCmd.AddCommand(scaleCmd)
}

func Scale(args []string) error {
	instances, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return errs.Wrap(err)
	}

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
			if instances == 1 {
				currentComposeServices[i].Deploy = nil
			} else if currentComposeServices[i].Deploy == nil {
				currentComposeServices[i].Deploy = &types.DeployConfig{Replicas: &instances}
			} else {
				*currentComposeServices[i].Deploy.Replicas = instances
			}
			currentComposeProject.Services[i] = currentComposeServices[i]
		}
	}

	resolvedServices, err := yaml.Marshal(&common.ComposeFile{Version: "3.4", Services: currentComposeProject.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}