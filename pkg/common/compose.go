package common

import (
	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
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
	for service := authservice; service <= grafana; service++ {
		if key&(1<<service) != 0 {
			result = append(result, service.String())
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

func CreateBind(volume string) types.ServiceVolumeConfig {
	mountPoints := strings.Split(volume, ":")

	return types.ServiceVolumeConfig{
		Type: "bind",
		Source: mountPoints[0],
		Target: mountPoints[1],
		ReadOnly: false,
		Consistency: "",
		Bind: &types.ServiceVolumeBind{
			Propagation: "",
			CreateHostPath: true,
		},
	}
}