// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"storj.io/storj-up/pkg/recipe"
)

// HostResolver helps to find the right hostname (or ip) for a specific service.
type HostResolver interface {
	GetHost(serviceInstance ServiceInstance, hostType string) string
}

// PortResolver helps to find the right port number for a specific service.
type PortResolver interface {
	GetPort(serviceInstance ServiceInstance, portType string) PortMap
}

// VariableGetter returns with any custom template variable.
type VariableGetter interface {
	Get(serviceInstance ServiceInstance, name string) string
}

// Runtime provides methods to read/write/modify any existing runtime definition (like compose/nomad/...)
type Runtime interface {
	HostResolver
	PortResolver
	VariableGetter

	// AddService creates and adds new service instance based on the recipe.
	AddService(recipe.Service) (Service, error)
	Write() error
	GetServices() []Service
	Reload(stack recipe.Stack) error
}

// PortMap defines the internal and external ports to use when port forwarding.
type PortMap struct {
	Internal int
	External int
}

// VolumeMount defines the type source and target fields when mounting a volume to the container.
type VolumeMount struct {
	MountType string
	Source    string
	Target    string
}

// Service is the interface to modify any service.
type Service interface {

	// ID returns with the unique identifier (including index in case of multiple instances are added from a service)
	ID() ServiceInstance
	GetENV() map[string]*string
	GetVolumes() []VolumeMount

	// ChangeImage applies transformation to the container image
	ChangeImage(func(string) string) error

	// AddConfig adds / changes existing configuration. Use it instead of AddEnvironment to be more generic.
	AddConfig(key string, value string) error
	AddFlag(flag string) error
	RemoveFlag(flag string) error

	// AddEnvironment registers new environment variable to be used. For normal configs, use AddConfig to be more general.
	AddEnvironment(key string, value string) error

	AddPortForward(PortMap) error
	Persist(dir string) error
	Labels() []string

	UseFile(path string, name string, data string) error
	UseFolder(path string, name string) error
}

// ServiceInstance is a unique identifier of a service instance.
type ServiceInstance struct {
	Name     string
	Instance int
}

func (i ServiceInstance) String() string {
	return fmt.Sprintf("%s/%d", i.Name, i.Instance)
}

// ServiceInstanceFromStr parses ServiceInstance from 'storagenode/2' like string expression.
func ServiceInstanceFromStr(service string) ServiceInstance {
	parts := strings.Split(service, "/")
	ix := 0
	if len(parts) > 1 {
		ix, _ = strconv.Atoi(parts[1])
	}
	return ServiceInstance{
		Name:     parts[0],
		Instance: ix,
	}
}

var indexedName = regexp.MustCompile(`(\D+)(\d*)`)

// ServiceInstanceFromIndexedName creates ID from string like 'storagenode0'.
func ServiceInstanceFromIndexedName(service string) ServiceInstance {
	submatch := indexedName.FindStringSubmatch(service)
	s := ServiceInstance{}
	if len(submatch) > 2 {
		s.Name = submatch[1]
		s.Instance, _ = strconv.Atoi(submatch[2])
		if s.Instance > 0 {
			s.Instance--
		}
	}
	return s
}

// NewServiceInstance creates the ID from name and instance identifier (number).
func NewServiceInstance(name string, i int) ServiceInstance {
	return ServiceInstance{
		Name:     name,
		Instance: i,
	}
}

// ModifyService applies a function to all services based on selector.
func ModifyService(stack recipe.Stack, rt Runtime, selectors []string, f func(service Service) error) error {
	for _, oneOrMoreSelector := range selectors {

		for _, selector := range strings.Split(oneOrMoreSelector, ",") {
			found := false

			for _, s := range rt.GetServices() {
				if s.ID().Name == selector { // selector can be the generic name of service (eg. storagenode without index)
					err := f(s)
					if err != nil {
						return err
					}
					found = true
				} else if fmt.Sprintf("%s%d", s.ID().Name, s.ID().Instance+1) == selector { // selector can be the exact service name
					err := f(s)
					if err != nil {
						return err
					}
					found = true
				}
			}

			if found {
				continue
			}

			// ok, it was not a service name. Might be a recipe name?
			for _, r := range stack {
				if r.Name == selector {
					for _, rs := range r.Add { // we will modify each of the services defined by the recipe
						for _, s := range rt.GetServices() {
							if s.ID().Name == rs.Name {
								err := f(s)
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}

		}
	}
	return nil
}

// SetFlag adds OR modify existing flag set.
func SetFlag(command []string, s string) []string {
	var res []string
	exist := false

	// s is in format foobar=value
	parts := strings.SplitN(s, "=", 2)

	for _, c := range command {
		if strings.HasPrefix(c, "-"+parts[0]+"=") || strings.HasPrefix(c, "--"+parts[0]+"=") {
			f := s
			if f[0] != '-' {
				f = "-" + f
			}

			// existing flag with two --
			if len(c) > 2 && c[1] == '-' {
				f = "-" + f
			}
			res = append(res, f)
			exist = true
		} else {
			res = append(res, c)
		}
	}
	if !exist {
		res = append(res, "--"+s)
	}
	return res
}

// RemoveFlag remove flags based on prefix.
func RemoveFlag(command []string, s string) []string {
	var res []string

	// s is in format foobar=value
	parts := strings.SplitN(s, "=", 2)

	for _, c := range command {
		if !strings.HasPrefix(c, "-"+parts[0]+"=") && !strings.HasPrefix(c, "--"+parts[0]+"=") {
			res = append(res, c)
		}
	}
	return res
}
