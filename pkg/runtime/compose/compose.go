// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package compose

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

// Compose represents the runtime.Runtime implementation for docker-compose based environments.
type Compose struct {
	project   *types.Project
	dir       string
	variables map[string]map[string]string
}

// Reload implements runtime.Runtime.
func (c *Compose) Reload(stack recipe.Stack) error {
	options := cli.ProjectOptions{
		Name:        common.ComposeFileName,
		ConfigPaths: []string{filepath.Join(c.dir, common.ComposeFileName)},
	}
	composeProject, err := cli.ProjectFromOptions(&options)
	if err != nil {
		return err
	}
	c.project = composeProject
	for i, service := range c.project.Services {
		for _, recipe := range stack {
			for _, baseService := range recipe.Add {
				if strings.Contains(service.Name, baseService.Name) && len(baseService.Label) > 0 {
					c.project.Services[i].Extensions = map[string]interface{}{"labels": baseService.Label}
				}
			}
		}
	}
	return nil
}

// GetServices  implements runtime.Runtime.
func (c *Compose) GetServices() []runtime.Service {
	k := make([]runtime.Service, len(c.project.Services))
	for ix, s := range c.project.Services {
		id := runtime.ServiceInstanceFromIndexedName(s.Name)
		var labels []string
		if s.Extensions["labels"] != nil {
			labels = s.Extensions["labels"].([]string)
		}
		k[ix] = &Service{
			id:         id,
			project:    c.project,
			composeDir: c.dir,
			render: func(s string) (string, error) {
				return runtime.Render(c, id, s)
			},
			labels: labels,
		}
	}
	return k
}

// Get implements runtime.Runtime.
func (c *Compose) Get(service runtime.ServiceInstance, name string) string {
	if name == "accessGrant" {
		sat := runtime.ServiceInstanceFromStr("satellite-api/0")
		key, err := common.GetTestAPIKey(fmt.Sprintf("%s@%s:%d", common.Satellite0Identity, c.GetHost(sat, "internal"), c.GetPort(sat, "public").Internal))
		if err != nil {
			return err.Error()
		}
		return key
	}
	return c.variables[service.Name][name]
}

// GetHost implements runtime.Runtime.
func (c *Compose) GetHost(service runtime.ServiceInstance, hostType string) string {
	switch hostType {
	case "listen":
		return "0.0.0.0"
	case "internal":
		if service.Name == "storagenode" {
			return service.Name + strconv.Itoa(service.Instance+1)
		}
		return service.Name
	case "external":
		return "localhost"
	}
	return "???"
}

// GetPort implements runtime.Runtime.
func (c *Compose) GetPort(service runtime.ServiceInstance, portType string) runtime.PortMap {
	if portType == "debug" {
		return runtime.PortMap{Internal: 11111, External: 11111}
	}
	switch service.Name {
	case "satellite-api":
		switch portType {
		case "public":
			return runtime.PortMap{Internal: 7777, External: 7777}
		case "console":
			return runtime.PortMap{Internal: 10000, External: 10000}
		}
	case "storagenode":
		p, _ := runtime.PortConvention(service, portType)
		return runtime.PortMap{Internal: p, External: p}
	case "gateway-mt":
		if portType == "public" {
			return runtime.PortMap{Internal: 9999, External: 9999}
		}
	case "authservice":
		if portType == "public" {
			return runtime.PortMap{Internal: 8888, External: 8888}
		}
	case "linksharing":
		if portType == "public" {
			return runtime.PortMap{Internal: 9090, External: 9090}
		}
	case "satellite-admin":
		if portType == "console" {
			return runtime.PortMap{Internal: 8080, External: 9080}
		}
	}

	return runtime.PortMap{Internal: -1, External: -1}
}

var _ runtime.Runtime = &Compose{}

// NewCompose creates a new compose runtime.
func NewCompose(dir string) (*Compose, error) {
	return &Compose{
		dir:     dir,
		project: &types.Project{},
		variables: map[string]map[string]string{
			"cockroach": {
				"main":     "cockroach://root@cockroach:26257/master?sslmode=disable",
				"metainfo": "cockroach://root@cockroach:26257/metainfo?sslmode=disable",
				"dir":      "/tmp/cockroach",
			},
			"storagenode": {
				"identityDir": "/var/lib/storj/.local/share/storj/identity/storagenode/",
				"staticDir":   "/var/lib/storj/web/storagenode",
			},
			"redis": {
				"url": "redis://redis:6379",
			},
			"satellite-api": {
				"mailTemplateDir": "/var/lib/storj/storj/web/satellite/static/emails/",
				"staticDir":       "/var/lib/storj/storj/web/satellite/",
				"identityDir":     "/var/lib/storj/.local/share/storj/identity/satellite-api/",
				"identity":        common.Satellite0Identity,
			},
			"satellite-core": {
				"mailTemplateDir": "/var/lib/storj/storj/web/satellite/static/emails/",
				"identityDir":     "/var/lib/storj/.local/share/storj/identity/satellite-api/",
			},
			"satellite-admin": {
				"staticDir":   "/var/lib/storj/storj/satellite/admin/ui/build",
				"identityDir": "/var/lib/storj/.local/share/storj/identity/satellite-api/",
			},
			"satellite-gc": {
				"identityDir": "/var/lib/storj/.local/share/storj/identity/satellite-api/",
			},
			"satellite-bf": {
				"identityDir": "/var/lib/storj/.local/share/storj/identity/satellite-api/",
			},
			"satellite-rangedloop": {
				"identityDir": "/var/lib/storj/.local/share/storj/identity/satellite-api/",
			},
			"linksharing": {
				"webDir":    "/var/lib/storj/pkg/linksharing/web/",
				"staticDir": "/var/lib/storj/pkg/linksharing/web/static",
			},
		},
	}, nil
}

// NewEmptyCompose creates a Compose service without any initialization.
func NewEmptyCompose(dir string) *Compose {
	return &Compose{
		dir:     dir,
		project: &types.Project{},
	}
}

// AddService implements runtime.Runtime.
func (c *Compose) AddService(recipe recipe.Service) (runtime.Service, error) {
	cmd := recipe.Command
	if recipe.Name == "cockroach" {
		recipe.Command = cmd[1:]
	}

	index := c.serviceCount(recipe.Name)
	name := recipe.Name
	if index > 0 {
		name += strconv.Itoa(index + 1)
	}
	if index == 1 {
		for ix, ds := range c.project.Services {
			if ds.Name == recipe.Name {
				c.project.Services[ix].Name = ds.Name + "1"
			}
		}
	}

	one := uint64(1)
	s := types.ServiceConfig{
		Name:        name,
		Command:     []string{},
		Environment: map[string]*string{},
		Deploy: &types.DeployConfig{
			Replicas: &one,
		},
		Ports:      []types.ServicePortConfig{},
		Image:      recipe.Image,
		Extensions: map[string]interface{}{"labels": recipe.Label},
	}

	if recipe.Name == "storagenode" || recipe.Name == "satellite-core" || recipe.Name == "satellite-admin" {
		s.Environment["STORJ_ROLE"] = ptrStr(recipe.Name)
		s.Environment["STORJ_WAIT_FOR_SATELLITE"] = ptrStr("true")
	} else if recipe.Name == "satellite-api" {
		s.Environment["STORJ_ROLE"] = ptrStr(recipe.Name)
		s.Environment["STORJ_WAIT_FOR_DB"] = ptrStr("true")
	} else if recipe.Name == "authservice" {
		s.Environment["STORJ_ROLE"] = ptrStr(recipe.Name)
	}

	c.project.Services = append(c.project.Services, s)

	id := runtime.NewServiceInstance(recipe.Name, index)
	r := &Service{
		id:         id,
		composeDir: c.dir,
		project:    c.project,
		render: func(s string) (string, error) {
			return runtime.Render(c, id, s)
		},
	}

	err := runtime.InitFromRecipe(r, recipe)
	if err != nil {
		return r, err
	}

	for _, public := range []string{"satellite-api", "gateway-mt", "linksharing", "authservice"} {
		if public == recipe.Name {
			err := r.AddPortForward(c.GetPort(id, "public"))
			if err != nil {
				return r, err
			}
		}
	}
	if recipe.Name == "satellite-api" || recipe.Name == "satellite-admin" {
		err := r.AddPortForward(c.GetPort(id, "console"))
		if err != nil {
			return r, err
		}
	}

	if recipe.Name == "satellite-api" {
		err := errs.Combine(
			r.AddEnvironment("STORJ_ROLE", "satellite-api"),
			r.AddEnvironment("STORJ_IDENTITY_DIR", "{{ Environment .This \"identityDir\"}}"))
		if err != nil {
			return nil, err
		}
	}
	if strings.HasPrefix(recipe.Name, "storagenode") {
		err := errs.Combine(
			r.AddPortForward(c.GetPort(id, "console")),
			r.AddEnvironment("STORJ_ROLE", "storagenode"),
			r.AddEnvironment("STORJ_IDENTITY_DIR", "{{ Environment .This \"identityDir\"}}"))
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (c *Compose) serviceCount(name string) int {
	i := 0

	for _, ds := range c.project.Services {
		if m, err := regexp.Match(name+"\\d*", []byte(ds.Name)); m && err == nil {
			i++
		}
	}
	return i
}

func ptrStr(name string) *string {
	return &name
}

// Write implements runtime.Runtime.
func (c *Compose) Write() error {
	return common.WriteComposeFile(c.dir, c.project)

}
