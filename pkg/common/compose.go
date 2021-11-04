package common

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/goccy/go-yaml"
)

type ComposeFile struct {
	Version  string
	Services types.Services
}

func LoadComposeFromFile(filename string) (*types.Project, error) {
	options := cli.ProjectOptions{
		Name:        filename,
		ConfigPaths: []string{"./" + filename},
	}

	return cli.ProjectFromOptions(&options)
}

func LoadComposeFromBytes(composeBytes []byte) (*types.Project, error) {
	return loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Content: composeBytes,
			},
		},
		WorkingDir: ".",
	})
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
		Type:        "bind",
		Source:      source,
		Target:      target,
		ReadOnly:    false,
		Consistency: "",
		Bind: &types.ServiceVolumeBind{
			Propagation:    "",
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

func UpdateEach(compose *types.Project, cmd func(*types.ServiceConfig, string) error, arg string, services []string) (*types.Project, error) {
	resolvedServices := ResolveServices(services)

	if len(resolvedServices) == 0 {
		return nil, fmt.Errorf("no service is selected for update. Try to use the right selector instead of \"%s\"", strings.Join(services, ","))
	}

	for _, service := range resolvedServices {
		for i, composeService := range compose.AllServices() {
			if strings.EqualFold(service, composeService.Name) {
				err := cmd(&compose.Services[i], arg)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return compose, nil
}
