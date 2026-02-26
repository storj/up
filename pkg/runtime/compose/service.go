// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package compose

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/v2/types"
	"golang.org/x/exp/slices"

	"storj.io/storj-up/pkg/runtime/runtime"
)

// Service is the implementation of runtime.Service for Docker compose files.
type Service struct {
	composeDir string
	id         runtime.ServiceInstance
	project    *types.Project
	render     func(string) (string, error)
	labels     []string
}

var _ runtime.Service = (*Service)(nil)
var _ runtime.ManageableNetwork = (*Service)(nil)

// GetENV implements runtime.Service.
func (s *Service) GetENV() map[string]*string {
	for _, ds := range s.project.Services {
		if filtered(s, ds) {
			return ds.Environment
		}
	}
	return nil
}

// GetVolumes implements runtime.Service.
func (s *Service) GetVolumes() (mounts []runtime.VolumeMount) {
	for _, ds := range s.project.Services {
		if filtered(s, ds) {
			for _, mount := range ds.Volumes {
				mounts = append(mounts, runtime.VolumeMount{
					MountType: mount.Type,
					Source:    mount.Source,
					Target:    mount.Target,
				})
			}
			return mounts
		}
	}
	return nil
}

// UseFolder implements runtime.Service.
func (s *Service) UseFolder(path string, name string) error {
	s.useVolume(path, name)
	return nil
}

// UseFile implements runtime.Service.
func (s *Service) UseFile(path string, name string, data string) error {
	err := os.WriteFile(filepath.Join(s.composeDir, name), []byte(data), 0644)
	if err != nil {
		return err
	}
	s.useVolume(path, name)
	return nil
}

// useVolume adds a bind mount to the service.
func (s *Service) useVolume(path string, name string) {
	if path == "" {
		path = "/tmp"
	}
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			ds.Volumes = append(ds.Volumes, types.ServiceVolumeConfig{
				Type:   "bind",
				Source: name,
				Target: strings.ReplaceAll(filepath.Join(path, name), string(filepath.Separator), "/"),
				Bind: &types.ServiceVolumeBind{
					CreateHostPath: true,
				},
			})
			s.project.Services[serviceName] = ds
		}
	}
}

// Labels implements runtime.Service.
func (s *Service) Labels() []string {
	return s.labels
}

// RemoveFlag implements runtime.Service.
func (s *Service) RemoveFlag(flag string) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			var filteredCmd []string
			for _, cmd := range ds.Command {
				if !strings.HasPrefix(cmd, flag+"=") {
					filteredCmd = append(filteredCmd, cmd)
				}
			}
			ds.Command = filteredCmd
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// Persist implements runtime.Service.
func (s *Service) Persist(dir string) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			ds.Volumes = append(ds.Volumes, types.ServiceVolumeConfig{
				Type:   "bind",
				Source: filepath.Join(s.composeDir, ds.Name, filepath.Base(dir)),
				Target: dir,
				Bind: &types.ServiceVolumeBind{
					CreateHostPath: true,
				},
			})
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// ChangeImage implements runtime.Service.
func (s *Service) ChangeImage(ch func(string) string) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			ds.Image = ch(ds.Image)
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// ID  implements runtime.Service.
func (s *Service) ID() runtime.ServiceInstance {
	return s.id
}

// AddConfig implements runtime.Runtime.
func (s *Service) AddConfig(key string, value string) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			rendered, err := s.render(value)
			if err != nil {
				return err
			}
			ds.Environment[key] = &rendered
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

func filtered(s *Service, ds types.ServiceConfig) bool {
	return (s.id.Name == ds.Name && s.id.Instance == 0) || ds.Name == s.id.Name+strconv.Itoa(s.id.Instance+1)
}

// AddPortForward implements runtime.Service.
func (s *Service) AddPortForward(ports runtime.PortMap) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			ds.Ports = append(ds.Ports, types.ServicePortConfig{
				Mode:      "ingress",
				Target:    uint32(ports.Internal),
				Published: strconv.Itoa(ports.External),
				Protocol:  ports.Protocol,
			})
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// RemovePortForward implements runtime.Service.
func (s *Service) RemovePortForward(ports runtime.PortMap) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			i := slices.IndexFunc(ds.Ports, func(port types.ServicePortConfig) bool {
				return port.Target == uint32(ports.Internal)
			})
			if i >= 0 {
				ds.Ports = slices.Delete(ds.Ports, i, i+1)
			}
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// AddNetwork implements runtime.ManageableNetwork.
func (s *Service) AddNetwork(networkAlias string) error {
	if s.project.Networks == nil {
		s.project.Networks = make(types.Networks)
	}
	if _, ok := s.project.Networks[networkAlias]; !ok {
		s.project.Networks[networkAlias] = types.NetworkConfig{
			Name:     networkAlias,
			External: true,
			Driver:   "default",
		}
	}
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			if ds.Networks == nil {
				ds.Networks = make(map[string]*types.ServiceNetworkConfig)
			}
			if ds.Networks[networkAlias] == nil {
				ds.Networks[networkAlias] = &types.ServiceNetworkConfig{}
			}
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// RemoveNetwork implements runtime.ManageableNetwork.
func (s *Service) RemoveNetwork(networkAlias string) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			delete(ds.Networks, networkAlias)
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// AddFlag implements runtime.Service.
func (s *Service) AddFlag(flag string) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			rendered, err := s.render(flag)
			if err != nil {
				return err
			}
			eqIndex := strings.Index(rendered, "=")
			if eqIndex >= 0 {
				commandIndex := slices.IndexFunc(ds.Command, func(command string) bool {
					return strings.HasPrefix(command, rendered[:eqIndex+1])
				})
				if commandIndex >= 0 {
					ds.Command = slices.Delete(ds.Command, commandIndex, commandIndex+1)
				}
			}
			ds.Command = append(ds.Command, rendered)
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// AddEnvironment registers new environment variable to be used. For normal configs, use AddConfig to be more general.
func (s *Service) AddEnvironment(key string, value string) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			rendered, err := s.render(value)
			if err != nil {
				return err
			}
			ds.Environment[key] = ptrStr(rendered)
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}

// TransformRaw enables to apply transformations on original raw docker service.
func (s *Service) TransformRaw(apply func(config *types.ServiceConfig) error) error {
	for serviceName, ds := range s.project.Services {
		if filtered(s, ds) {
			err := apply(&ds)
			if err != nil {
				return err
			}
			s.project.Services[serviceName] = ds
		}
	}
	return nil
}
