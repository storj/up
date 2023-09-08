// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package standalone

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"

	"storj.io/storj-up/pkg/runtime/runtime"
)

type service struct {
	id          runtime.ServiceInstance
	Command     []string
	ConfigFile  string
	render      func(string) (string, error)
	config      []string
	Environment map[string]string
	labels      []string
}

var _ runtime.Service = (*service)(nil)

func (s *service) GetVolumes() []runtime.VolumeMount {
	// TODO implement me
	return nil
}

func (s *service) GetENV() map[string]*string {
	// TODO implement me
	return nil
}

func (s *service) UseFolder(path string, name string) error {
	// TODO: folders are not yet supported
	return nil
}

func (s *service) UseFile(path string, name string, data string) error {
	// TODO: files are not yet extracted, but we accept recipes with files
	return nil
}

func (s *service) Labels() []string {
	return s.labels
}

func (s *service) RemoveFlag(flag string) error {
	fmt.Println("RemoveFlag for Standalone is not yet implemented!")
	return nil
}

func (s *service) Persist(dir string) error {
	// NOOP: Standalone runner doesn't use containers. No need to persist directories
	return nil
}

func (s *service) ChangeImage(func(string) string) error {
	// NOOP: we accept any container
	return nil
}

func (s *service) AddPortForward(runtime.PortMap) error {
	// all ports are available, by default...
	return nil
}

func (s *service) RemovePortForward(runtime.PortMap) error {
	// all ports are available, by default...
	return nil
}

func (s *service) ID() runtime.ServiceInstance {
	return s.id
}

func (s *service) AddConfig(key string, value string) error {
	value, err := s.render(value)
	if err != nil {
		return err
	}

	for ix, line := range s.config {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(strings.TrimPrefix(parts[0], "#"))
			if camelToUpperCase(k) == key {
				s.config[ix] = fmt.Sprintf("%s: %s", k, strings.TrimSpace(value))
				return nil
			}
		}
	}
	s.config = append(s.config, fmt.Sprintf("%s: %s", key, value))
	return nil
}

func (s *service) AddFlag(flag string) error {
	f, err := s.render(flag)
	if err != nil {
		return err
	}
	eqIndex := strings.Index(f, "=")
	if eqIndex >= 0 {
		commandIndex := slices.IndexFunc(s.Command, func(command string) bool {
			return strings.HasPrefix(command, f[:eqIndex+1])
		})
		if commandIndex >= 0 {
			s.Command = slices.Delete(s.Command, commandIndex, commandIndex+1)
		}
	}
	s.Command = append(s.Command, f)
	return nil
}

func (s *service) AddEnvironment(key string, value string) error {
	v, err := s.render(value)
	if err != nil {
		return err
	}
	s.Environment[key] = v
	return nil
}

func camelToUpperCase(name string) string {
	smallCapital := regexp.MustCompile("([a-z])([A-Z])")
	name = smallCapital.ReplaceAllString(name, "${1}_$2")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return strings.ToUpper("STORJ_" + name)
}
