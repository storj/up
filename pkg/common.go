package sjr

import (
	_ "embed"
	"fmt"
	"github.com/elek/sjr/pkg/common"
	"github.com/elek/sjr/pkg/recipe"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var (
	RootCmd = &cobra.Command{}
)

//go:embed docker-compose.template.yaml
var composeTemplate []byte


func init() {
	current, err := common.ReadCompose("docker-compose.yaml")
	if err != nil {
		panic("docker-compose.yaml couldn't be read from the local dir " + err.Error())
	}

	for service, _ := range current.Services {
		serviceCmd := &cobra.Command{
			Use:   service,
			Short: fmt.Sprintf("Customize the %s service", service),
		}

		recipe.AddMutatorCommands(service, serviceCmd)
		RootCmd.AddCommand(serviceCmd)
	}

	for group, values := range common.Presets {

		short := fmt.Sprintf("Customize all %s services (%s)", group, strings.Join(values, ", "))
		if group == "all" {
			short = "Customize all activated service"
		}

		serviceCmd := &cobra.Command{
			Use:   group,
			Short: short,
		}
		recipe.AddMutatorCommands(group, serviceCmd)
		RootCmd.AddCommand(serviceCmd)
	}

	RootCmd.AddCommand(&cobra.Command{
		Use:   "init [groups/service]",
		Short: "Creates/overwrites local docker-compose.yaml with service. You can use predefined groups as arguments.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return initEnv("db,storj")
			} else {
				return initEnv(args[0])
			}
		},
	})
	addCmd := &cobra.Command{
		Use:   "add [service name or group]",
		Short: "Add more services to the docker-compose.yaml. You can use predefined groups as arguments.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return add(args[0])
		},
	}
	RootCmd.AddCommand(addCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Print all the configured services from the dockerfile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list()
		},
	}
	RootCmd.AddCommand(listCmd)
}

func list() error {
	current, err := common.ReadCompose("docker-compose.yaml")
	if err != nil {
		return err
	}
	for name, s := range current.Services {
		fmt.Printf("%s (%s)\n", name, s.Image)
	}
	return nil
}

func add(group string) error {

	template, err := common.ParseCompose(composeTemplate)
	if err != nil {
		return err
	}

	current, err := common.ReadCompose("docker-compose.yaml")
	if err != nil {
		return err
	}

	for k, v := range template.FilterPrefixAndGroup(group, common.Presets) {
		if _, found := current.Services[k]; !found {
			current.Services[k] = v
		}
	}

	out, err := yaml.Marshal(current)
	if err != nil {
		return err
	}
	if err = writeDependencies(group); err != nil {
		return err
	}
	if err = ioutil.WriteFile("docker-compose.yaml", out, 0644); err != nil {
		return err
	}
	return nil
}

func initEnv(group string) error {
	template, err := common.ParseCompose(composeTemplate)
	if err != nil {
		return err
	}

	filtered := make(map[string]*common.ServiceConfig, 0)
	for k, v := range template.FilterPrefixAndGroup(group, common.Presets) {
		filtered[k] = v
	}

	template.Services = filtered
	out, err := yaml.Marshal(template)
	if err != nil {
		return err
	}
	if err = writeDependencies(group); err != nil {
		return err
	}
	if err = ioutil.WriteFile("docker-compose.yaml", out, 0644); err != nil {
		return err
	}
	return nil
}

//go:embed prometheus.yml
var prometheusConfig []byte

func writeDependencies(group string) error {
	if strings.Contains(group, "prometheus") {
		return common.ExtractFile(prometheusConfig, "prometheus.yml")
	}
	return nil
}


