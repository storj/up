// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package compose

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/types"
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
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			return s.project.Services[ix].Environment
		}
	}
	return nil
}

// GetVolumes implements runtime.Service.
func (s *Service) GetVolumes() (mounts []runtime.VolumeMount) {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			for _, mount := range s.project.Services[ix].Volumes {
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
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			s.project.Services[ix].Volumes = append(s.project.Services[ix].Volumes, types.ServiceVolumeConfig{
				Type:   "bind",
				Source: name,
				Target: strings.ReplaceAll(filepath.Join(path, name), string(filepath.Separator), "/"),
				Bind: &types.ServiceVolumeBind{
					CreateHostPath: true,
				},
			})
		}
	}
}

// Labels implements runtime.Service.
func (s *Service) Labels() []string {
	return s.labels
}

// RemoveFlag implements runtime.Service.
func (s *Service) RemoveFlag(flag string) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			var filtered []string
			for _, s := range s.project.Services[ix].Command {
				if !strings.HasPrefix(s, flag+"=") {
					filtered = append(filtered, s)
				}
			}
			s.project.Services[ix].Command = filtered
		}
	}
	return nil
}

// Persist implements runtime.Service.
func (s *Service) Persist(dir string) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			s.project.Services[ix].Volumes = append(s.project.Services[ix].Volumes, types.ServiceVolumeConfig{
				Type:   "bind",
				Source: filepath.Join(s.composeDir, s.project.Services[ix].Name, filepath.Base(dir)),
				Target: dir,
				Bind: &types.ServiceVolumeBind{
					CreateHostPath: true,
				},
			})
		}
	}
	return nil
}

// ChangeImage implements runtime.Service.
func (s *Service) ChangeImage(ch func(string) string) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			s.project.Services[ix].Image = ch(s.project.Services[ix].Image)
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
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			rendered, err := s.render(value)
			if err != nil {
				return err
			}
			s.project.Services[ix].Environment[key] = &rendered
		}
	}
	return nil
}

func filtered(s *Service, ds types.ServiceConfig) bool {
	return (s.id.Name == ds.Name && s.id.Instance == 0) || ds.Name == s.id.Name+strconv.Itoa(s.id.Instance+1)
}

// AddPortForward implements runtime.Service.
func (s *Service) AddPortForward(ports runtime.PortMap) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			s.project.Services[ix].Ports = append(s.project.Services[ix].Ports, types.ServicePortConfig{
				Mode:      "ingress",
				Target:    uint32(ports.Internal),
				Published: uint32(ports.External),
				Protocol:  ports.Protocol,
			})
		}
	}
	return nil
}

// RemovePortForward implements runtime.Service.
func (s *Service) RemovePortForward(ports runtime.PortMap) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			service := &s.project.Services[ix]
			i := slices.IndexFunc(service.Ports, func(port types.ServicePortConfig) bool {
				return port.Target == uint32(ports.Internal)
			})
			if i >= 0 {
				service.Ports = slices.Delete(service.Ports, i, i+1)
			}
		}
	}
	return nil
}

// AddNetwork implements runtime.ManageableNetwork.
func (s *Service) AddNetwork(networkAlias string) error {
	if _, ok := s.project.Networks[networkAlias]; !ok {
		s.project.Networks[networkAlias] = types.NetworkConfig{
			Name: networkAlias,
			External: types.External{
				External: true,
			},
			Driver: "default",
		}
	}
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			if s.project.Services[ix].Networks[networkAlias] == nil {
				s.project.Services[ix].Networks[networkAlias] = &types.ServiceNetworkConfig{}
			}
		}
	}
	return nil
}

// RemoveNetwork implements runtime.ManageableNetwork.
func (s *Service) RemoveNetwork(networkAlias string) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			delete(s.project.Services[ix].Networks, networkAlias)
		}
	}
	return nil
}

// AddFlag implements runtime.Service.
func (s *Service) AddFlag(flag string) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			rendered, err := s.render(flag)
			if err != nil {
				return err
			}
			eqIndex := strings.Index(rendered, "=")
			if eqIndex >= 0 {
				commandIndex := slices.IndexFunc(s.project.Services[ix].Command, func(command string) bool {
					return strings.HasPrefix(command, rendered[:eqIndex+1])
				})
				if commandIndex >= 0 {
					s.project.Services[ix].Command = slices.Delete(s.project.Services[ix].Command, commandIndex, commandIndex+1)
				}
			}
			s.project.Services[ix].Command = append(s.project.Services[ix].Command, rendered)
		}
	}
	return nil
}

// AddEnvironment registers new environment variable to be used. For normal configs, use AddConfig to be more general.
func (s *Service) AddEnvironment(key string, value string) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			rendered, err := s.render(value)
			if err != nil {
				return err
			}
			s.project.Services[ix].Environment[key] = ptrStr(rendered)
		}
	}
	return nil
}

// TransformRaw enables to apply transformations on original raw docker service.
func (s *Service) TransformRaw(apply func(config *types.ServiceConfig) error) error {
	for ix, ds := range s.project.Services {
		if filtered(s, ds) {
			err := apply(&s.project.Services[ix])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
