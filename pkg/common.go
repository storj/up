package sjr

import (
	_ "embed"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var (
	RootCmd = &cobra.Command{}
	presets = createPresets()
)

//go:embed docker-compose.template.yaml
var composeTemplate []byte

func createPresets() map[string][]string {
	presets := map[string][]string{}
	presets["minimal"] = []string{"satellite-api", "storagenode"}
	presets["edge"] = []string{"gateway-mt", "linksharing", "authservice"}
	presets["db"] = []string{"cockroach", "redis"}
	presets["monitor"] = []string{"prometheus", "grafana"}
	presets["core"] = append(presets["minimal"], "satellite-core", "satellite-admin", "versioncontrol")
	presets["storj"] = append(presets["core"], presets["edge"]...)
	presets["storj"] = append(presets["storj"], "uplink")
	return presets
}

func init() {
	for _, service := range presets["storj"] {
		serviceCmd := &cobra.Command{
			Use:   service,
			Short: fmt.Sprintf("Customize the %s service", service),
		}

		addMutatorCommands(service, serviceCmd)
		RootCmd.AddCommand(serviceCmd)
	}
	for group, values := range presets {

		short := fmt.Sprintf("Customize all %s services (%s)", group, strings.Join(values, ", "))
		if group == "all" {
			short = "Customize all activated service"
		}

		serviceCmd := &cobra.Command{
			Use:   group,
			Short: short,
		}
		addMutatorCommands(group, serviceCmd)
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
	current, err := ReadCompose("docker-compose.yaml")
	if err != nil {
		return err
	}
	for name, s := range current.Services {
		fmt.Printf("%s (%s)\n", name, s.Image)
	}
	return nil
}

func add(group string) error {

	template, err := ParseCompose(composeTemplate)
	if err != nil {
		return err
	}

	current, err := ReadCompose("docker-compose.yaml")
	if err != nil {
		return err
	}

	for k, v := range template.Services {
		if _, found := current.Services[k]; !found && selected(group, k) {
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
	content, err := ParseCompose(composeTemplate)
	if err != nil {
		return err
	}

	filtered := make(map[string]*ServiceConfig, 0)
	for k, v := range content.Services {
		if selected(group, k) {
			filtered[k] = v
		}
	}

	content.Services = filtered
	out, err := yaml.Marshal(content)
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
		return extractFile(prometheusConfig, "prometheus.yml")
	}
	return nil
}

func addMutatorCommands(group string, serviceCmd *cobra.Command) {
	debugCmd(group, serviceCmd)
	imageCmd(group, serviceCmd)
	localCmd(group, serviceCmd)
	versionCmd(group, serviceCmd)
	localEntrypointCmd(group, serviceCmd)
	envCmd(group, serviceCmd)
	buildCmd(group, serviceCmd)
}

func selected(selector string, service string) bool {
	for _, part := range strings.Split(selector, ",") {
		selector := strings.TrimSpace(part)
		if selector == "all" {
			return true
		}
		if selector == service {
			return true
		}
		if group, found := presets[selector]; found {
			for _, s := range group {
				if s == service {
					return true
				}
			}
		}
	}
	return false
}

func extractFile(content []byte, fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return ioutil.WriteFile(fileName, content, 0644)
	}
	fmt.Printf("File %s exists/couldn't be checked. Skipping to write\n", fileName)
	return nil
}
