package common

import (
	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"strings"
)

type ComposeFile struct {
	Version   string
	Services  types.Services
}

func CreateComposeProject(filename string) (*types.Project, error) {
	options := cli.ProjectOptions{
		Name:	filename,
		ConfigPaths: []string{"./" + filename},
	}

	return cli.ProjectFromOptions(&options)
}

func ResolveServices(services []string) ([]string, error) {
	var result []string
	var key uint
	for _, service := range services {
		key |= ServiceDict[service]
	}
	for service := authservice; service <= appstorj; service++ {
		if key&(1<<service) != 0 {
			result = append(result, serviceNameHelper[service.String()])
		}
	}
	return result, nil
}

func ContainsService(s []types.ServiceConfig, e string) bool {
	for _, a := range s {
		if a.Name == e {
			return true
		}
	}
	return false
}

func CreateBind(source string, target string) types.ServiceVolumeConfig {
	return types.ServiceVolumeConfig{
		Type: "bind",
		Source: source,
		Target: target,
		ReadOnly: false,
		Consistency: "",
		Bind: &types.ServiceVolumeBind{
			Propagation: "",
			CreateHostPath: true,
		},
	}
}

func WriteComposeFile(compose *types.Project) error {
	resolvedServices, err := yaml.Marshal(&ComposeFile{Version: "3.4", Services: compose.Services})
	if err = ioutil.WriteFile("docker-compose.yaml", resolvedServices, 0644); err != nil {
		return err
	}
	return nil
}

func LoadCompose(composeDir string) (*types.Project, error) {
	currentComposeProject, err := CreateComposeProject(composeDir)
	if err != nil {
		return nil, err
	}

	return currentComposeProject, nil
}

func UpdateEach(compose *types.Project, cmd func(*types.ServiceConfig, string) error, arg string, services []string) (*types.Project, error) {
	resolvedServices, err := ResolveServices(services)
	if err != nil {
		return nil, err
	}

	for _, service := range resolvedServices {
		for i, composeService := range compose.AllServices() {
			if strings.EqualFold(service, composeService.Name) {
				cmd(&compose.Services[i], arg)
			}
		}
	}
	return compose, nil
}