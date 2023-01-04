// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package standalone

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zeebo/errs"

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
	return errs.New("RemoveFlag for Standalone is not yet implemented")
}

func (s *service) Persist(dir string) error {
	return errs.New("Standalone runner doesn't use containers. No need to persist directories.")
}

func (s *service) ChangeImage(func(string) string) error {
	// NOOP: we accept any container
	return nil
}

func (s *service) AddPortForward(runtime.PortMap) error {
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
			k := parts[0]
			v := strings.TrimSpace(parts[1])
			k = strings.TrimSpace(strings.TrimPrefix(k, "#"))

			if camelToUpperCase(k) == key {
				if len(v) > 0 && rune(v[0]) == '"' {
					value = "\"" + value + "\""
				}
				s.config[ix] = fmt.Sprintf("%s: %s", k, value)
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
