// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package nomad

import (
	// to use "go:embed".
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zeebo/errs"

	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

// Nomad is the runtime.Runtime implementation for Hashicorp Nomad based environments.
type Nomad struct {
	dir       string
	External  string
	job       *hclwrite.File
	group     *hclwrite.Block
	jobs      []*service
	variables map[string]map[string]string
}

// NewNomad creates a new runtime.Runtime for Nomad environments.
func NewNomad(dir string, name string) (*Nomad, error) {
	job := hclwrite.NewEmptyFile()
	j := job.Body().AppendNewBlock("job", []string{name})
	j.Body().SetAttributeValue("type", cty.StringVal("service"))
	j.Body().SetAttributeValue("datacenters", cty.ListVal([]cty.Value{cty.StringVal("dc1")}))
	group := j.Body().AppendNewBlock("group", []string{name})
	group.Body().SetAttributeValue("count", cty.NumberIntVal(1))
	group.Body().AppendNewBlock("restart", []string{}).Body().SetAttributeValue("attempts", cty.NumberIntVal(20))
	return &Nomad{
		External: "localhost",
		dir:      dir,
		job:      job,
		group:    group,
		variables: map[string]map[string]string{
			"cockroach": {
				"main":     "cockroach://root@localhost:26257/master?sslmode=disable",
				"metainfo": "cockroach://root@localhost:26257/metainfo?sslmode=disable",
				"dir":      "/tmp/cockroach",
			},
			"storagenode": {
				"identityDir": "/var/lib/storj/.local/share/storj/identity/storagenode/",
				"staticDir":   "/var/lib/storj/web/storagenode",
			},
			"redis": {
				"url": "redis://localhost:6379",
			},
			"satellite-api": {
				"mailTemplateDir": "/var/lib/storj/storj/web/satellite/static/emails/",
				"staticDir":       "/var/lib/storj/storj/web/satellite/",
				"identityDir":     "/var/lib/storj/identities/1",
				"identity":        common.Satellite0Identity,
			},
			"satellite-core": {
				"mailTemplateDir": "/var/lib/storj/storj/web/satellite/static/emails/",
				"identityDir":     "/var/lib/storj/identities/1",
			},
			"satellite-admin": {
				"staticDir":   "/var/lib/storj/storj/web/satellite/",
				"identityDir": "/var/lib/storj/identities/1",
			},
			"satellite-gc": {
				"identityDir": "/var/lib/storj/identities/1",
			},
			"satellite-bf": {
				"identityDir": "/var/lib/storj/identities/1",
			},
			"satellite-rangedloop": {
				"identityDir": "/var/lib/storj/identities/1",
			},
			"linksharing": {
				"webDir":    "/var/lib/storj/pkg/linksharing/web/",
				"staticDir": "/var/lib/storj/pkg/linksharing/web/static",
			},
		},
	}, nil
}

func firstBlockByType(body *hclwrite.Body, name string) *hclwrite.Block {
	for _, b := range body.Blocks() {
		if b.Type() == name {
			return b
		}
	}
	return nil
}

// Reload implements runtime.Runtime.
func (c *Nomad) Reload(stack recipe.Stack) error {
	content, err := os.ReadFile("storj.hcl")
	if err != nil {
		return errs.Wrap(err)
	}

	var diags hcl.Diagnostics
	c.job, diags = hclwrite.ParseConfig(content, "storj.hcl", hcl.Pos{})
	if len(diags) != 0 {
		return errs.Wrap(err)
	}

	c.jobs = make([]*service, 0)

	j := firstBlockByType(c.job.Body(), "job")
	c.group = firstBlockByType(j.Body(), "group")

	for _, b := range c.group.Body().Blocks() {
		if b.Type() != "task" {
			continue
		}
		id := runtime.ServiceInstanceFromIndexedName(b.Labels()[0])
		s := &service{
			id:   id,
			task: b,
			env:  b.Body().FirstMatchingBlock("env", []string{}),
			render: func(s string) (string, error) {
				return runtime.Render(c, id, s)
			},
		}
		recipe, err := stack.FindRecipeByName(id.Name)
		if err != nil {
			s.labels = recipe.Label
		}
		c.jobs = append(c.jobs, s)
	}
	return nil
}

// GetServices implements runtime.Runtime.
func (c *Nomad) GetServices() []runtime.Service {
	k := make([]runtime.Service, len(c.jobs))
	for ix, s := range c.jobs {
		k[ix] = s
	}
	return k
}

// Get implements runtime.Runtime.
func (c *Nomad) Get(s runtime.ServiceInstance, name string) string {
	if name == "accessGrant" {
		sat := runtime.ServiceInstanceFromStr("satellite-api/0")
		key, err := common.GetTestAPIKey(fmt.Sprintf("%s@%s:%d", common.Satellite0Identity, c.GetHost(sat, "external"), c.GetPort(sat, "public").External))
		if err != nil {
			return err.Error()
		}
		return key
	}
	return c.variables[s.Name][name]
}

// GetHost implements runtime.Runtime.
func (c *Nomad) GetHost(service runtime.ServiceInstance, hostType string) string {
	switch hostType {
	case "external":
		return c.External
	case "internal":
		return c.External
	case "listen":
		return "0.0.0.0"
	}
	return "????"
}

// GetPort implements runtime.Runtime.
func (c *Nomad) GetPort(service runtime.ServiceInstance, portType string) runtime.PortMap {
	port, err := runtime.PortConvention(service, portType)
	if err != nil {
		panic(err.Error())
	}
	return runtime.PortMap{Internal: port, External: port}
}

// AddService implements runtime.Runtime.
func (c *Nomad) AddService(rcp recipe.Service) (runtime.Service, error) {
	i := 0
	for _, b := range c.group.Body().Blocks() {
		for _, l := range b.Labels() {
			if strings.HasPrefix(l, rcp.Name) {
				i++
			}
		}
	}

	if i == 1 {
		// second instance, let's rename the first
		c.group.Body().FirstMatchingBlock("task", []string{rcp.Name}).SetLabels([]string{rcp.Name + strconv.Itoa(i)})
	}
	name := rcp.Name
	if i > 0 {
		name += strconv.Itoa(i + 1)
	}

	task := c.group.Body().AppendNewBlock("task", []string{name})
	task.Body().SetAttributeValue("driver", cty.StringVal("docker"))

	memory := int64(500)
	cpu := int64(500)
	if name == "cockroach" {
		cpu = 1000
		memory = 4000
	}
	resources := task.Body().AppendNewBlock("resources", []string{})
	resources.Body().SetAttributeValue("memory", cty.NumberIntVal(memory))
	resources.Body().SetAttributeValue("cpu", cty.NumberIntVal(cpu))

	cfg := task.Body().AppendNewBlock("config", []string{})
	cfg.Body().SetAttributeValue("image", cty.StringVal(rcp.Image))
	cfg.Body().SetAttributeValue("network_mode", cty.StringVal("host"))

	id := runtime.NewServiceInstance(rcp.Name, i)
	s := &service{
		id:   id,
		task: task,
		env:  task.Body().AppendNewBlock("env", []string{}),
		render: func(s string) (string, error) {
			return runtime.Render(c, id, s)
		},
	}

	cmds := rcp.Command
	if cmds == nil {
		cmds = []string{}
	}
	if rcp.Name == "cockroach" {
		cmds = cmds[1:]
	}
	rcp.Command = cmds

	err := runtime.InitFromRecipe(s, rcp)
	if err != nil {
		return nil, err
	}

	c.jobs = append(c.jobs, s)
	if rcp.Name == "satellite-api" {
		err := errs.Combine(
			s.AddEnvironment("STORJ_ROLE", "satellite-api"),
			s.AddEnvironment("STORJ_IDENTITY_DIR", "/var/lib/storj/identities/1"),
			s.AddFlag("--identity-dir=/var/lib/storj/identities/1"))
		if err != nil {
			return nil, err
		}
	}
	if strings.HasPrefix(rcp.Name, "storagenode") {
		err := errs.Combine(
			s.AddEnvironment("STORJ_ROLE", "storagenode"),
			s.AddEnvironment("STORJ_IDENTITY_DIR", "/var/lib/storj/.local/share/storj/identity/storagenode/"),
			s.AddFlag("--identity-dir=/var/lib/storj/.local/share/storj/identity/storagenode/"))

		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func toCtyStrList(values []string) cty.Value {
	flags := make([]cty.Value, 0)
	for _, f := range values {
		flags = append(flags, cty.StringVal(f))
	}
	return cty.ListVal(flags)
}

// Write implements runtime.Runtime.
func (c *Nomad) Write() error {
	return os.WriteFile(filepath.Join(c.dir, "storj.hcl"), c.job.Bytes(), 0644)
}

var _ runtime.Runtime = &Nomad{}
