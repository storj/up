// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package standalone

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zeebo/errs/v2"

	"storj.io/common/identity"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

//go:embed ca.cert
var caCert []byte

//go:embed ca.key
var caKey []byte

//go:embed identity.cert
var identityCert []byte

//go:embed identity.key
var identityKey []byte

// Standalone is a runtime.Runtime implementation for shell based execution.
type Standalone struct {
	dir        string
	services   []*service
	variables  map[string]map[string]string
	clean      bool
	Intellij   bool
	ProjectDir string
}

// Paths contains directories required for storj-up standalone instances.
type Paths struct {
	ScriptDir  string
	StorjDir   string
	GatewayDir string
	CleanDir   bool
}

// Reload implements runtime.Runtime.
func (c *Standalone) Reload(stack recipe.Stack) error {
	scripts, err := find(c.dir, ".sh")
	if err != nil {
		return err
	}
	for _, script := range scripts {
		for _, recipe := range stack {
			for _, service := range recipe.Add {
				if strings.Contains(script, service.Name) {
					_, err = c.AddService(*service)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	for _, script := range scripts {
		for _, service := range c.services {
			scriptPath := strings.TrimSuffix(script, filepath.Ext(script))
			sInst := runtime.ServiceInstanceFromIndexedName(filepath.Base(scriptPath))
			if service.id.Name == sInst.Name && service.id.Instance == sInst.Instance {
				err := c.reloadEnvironment(service, script)
				if err != nil {
					return err
				}
				err = c.reloadConfig(service)
				if err != nil {
					return err
				}
				err = os.Remove(script)
				if err != nil {
					return err
				}
				_ = os.Remove(scriptPath + ".run.xml")
			}
		}
	}
	return nil
}

func (c *Standalone) reloadEnvironment(service *service, script string) error {
	if service.Environment == nil {
		service.Environment = make(map[string]string)
	}
	file, err := os.Open(script)
	if err != nil {
		return err
	}
	defer func() { err = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimPrefix(line, "export ")
			if env, value, ok := strings.Cut(line, "="); ok {
				service.Environment[env] = strings.Trim(value, "\"")
			}
		}
	}
	return nil
}

func (c *Standalone) reloadConfig(service *service) error {
	serviceDir := filepath.Join(c.dir, service.id.Name, strconv.Itoa(service.id.Instance))
	configFile := filepath.Join(serviceDir, "config.yaml")
	if _, err := os.Stat(configFile); err == nil {
		file, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer func() { err = file.Close() }()
		service.config = nil
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			service.config = append(service.config, line)
		}
	}
	return nil
}

// GetServices implements runtime.Runtime.
func (c *Standalone) GetServices() []runtime.Service {
	k := make([]runtime.Service, len(c.services))
	for ix, s := range c.services {
		k[ix] = s
	}
	return k
}

// Get implements runtime.Runtime.
func (c *Standalone) Get(service runtime.ServiceInstance, name string) string {
	if name == "identityDir" {
		return service.String()
	}
	if name == "accessGrant" {
		sat := runtime.ServiceInstanceFromStr("satellite-api/0")
		key, err := common.GetTestAPIKey(fmt.Sprintf("%s@%s:%d", common.Satellite0Identity, c.GetHost(sat, "external"), c.GetPort(sat, "public").External))
		if err != nil {
			return err.Error()
		}
		return key
	}
	return c.variables[service.Name][name]
}

// GetHost implements runtime.Runtime.
func (c *Standalone) GetHost(service runtime.ServiceInstance, hostType string) string {
	return "localhost"
}

// GetPort implements runtime.Runtime.
func (c *Standalone) GetPort(service runtime.ServiceInstance, portType string) runtime.PortMap {
	port, err := runtime.PortConvention(service, portType)
	if err != nil {
		panic(err.Error())
	}
	return runtime.PortMap{Internal: port, External: port}
}

var (
	_ runtime.Runtime = &Standalone{}
)

// AddService implements runtime.Runtime.
func (c *Standalone) AddService(recipe recipe.Service) (runtime.Service, error) {
	i := c.serviceCount(recipe.Name)

	id := runtime.NewServiceInstance(recipe.Name, i)
	s := &service{
		id: id,
		render: func(s string) (string, error) {
			return runtime.Render(c, id, s)
		},
		config:      []string{},
		Command:     []string{},
		labels:      recipe.Label,
		Environment: map[string]string{},
	}
	if s.labels == nil {
		s.labels = []string{}
	}

	serviceDir := filepath.Join(c.dir, s.id.Name, strconv.Itoa(s.id.Instance))
	if c.clean {
		_ = os.RemoveAll(serviceDir)
	}
	_ = os.MkdirAll(serviceDir, 0755)

	var configFile string
	if recipe.HasLabel("storj") {

		if _, err := os.Stat(filepath.Join(serviceDir, "identity.cert")); os.IsNotExist(err) {
			err := c.generateIdentity(s.id.Name, s.id.Instance)
			if err != nil {
				return nil, err
			}
		}

		configFile = filepath.Join(serviceDir, "config.yaml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) && len(recipe.Command) > 0 {
			args := []string{"setup", "--config-dir=" + filepath.Dir(configFile)}
			if id.Name == "storagenode" {
				args = append(args, "--identity-dir", filepath.Dir(configFile))
			}
			cmd := exec.Command(recipe.Command[0], args...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(string(out))
				return nil, errs.Wrap(err)
			}

			cfg, err := os.ReadFile(configFile)
			if err != nil {
				return nil, errs.Wrap(err)
			}
			s.config = strings.Split(string(cfg), "\n")
			for ix, line := range s.config {
				line = strings.TrimSpace(line)
				if len(line) > 0 && line[0] != '#' {
					// to have same config everywhere, we turn off generated defaults
					s.config[ix] = "#" + line
				}
			}
		}
	}

	err := runtime.InitFromRecipe(s, recipe)
	if err != nil {
		return s, errs.Wrap(err)
	}

	if recipe.HasLabel("storj") {
		err := s.AddFlag("--config-dir=" + "\"" + filepath.Dir(configFile) + "\"")
		if err != nil {
			return s, err
		}
	}

	c.services = append(c.services, s)
	return s, nil
}

func (c *Standalone) serviceCount(name string) int {
	i := 0
	for _, o := range c.services {
		if o.id.Name == runtime.ServiceInstanceFromIndexedName(name).Name {
			i++
		}
	}
	return i
}

// NewStandalone returns with a new runtime, starting services without any container isolation (like storj-sim).
func NewStandalone(paths Paths) (*Standalone, error) {
	s := &Standalone{
		clean:    paths.CleanDir,
		dir:      paths.ScriptDir,
		services: []*service{},
		variables: map[string]map[string]string{
			"cockroach": {
				"main":     "cockroach://root@localhost:26257/master?sslmode=disable",
				"metainfo": "cockroach://root@localhost:26257/metainfo?sslmode=disable",
				"dir":      filepath.Join(paths.ScriptDir, "cockroach", "0", "data"),
			},
			"storagenode": {
				"staticDir": filepath.Join(paths.StorjDir, "web/storagenode"),
			},
			"redis": {
				"url": "redis://localhost:6379",
			},
			"satellite-api": {
				"mailTemplateDir": filepath.Join(paths.StorjDir, "web/satellite/static/emails"),
				"staticDir":       filepath.Join(paths.StorjDir, "web/satellite"),
			},
			"satellite-core": {
				"mailTemplateDir": filepath.Join(paths.StorjDir, "web/satellite/static/emails"),
			},
			"satellite-admin": {
				"staticDir": filepath.Join(paths.StorjDir, "web/satellite"),
			},
			"linksharing": {
				"webDir":    filepath.Join(paths.GatewayDir, "pkg/linksharing/web"),
				"staticDir": filepath.Join(paths.GatewayDir, "pkg/linksharing/web/static"),
			},
		},
	}
	s.variables["satellite-api"]["identity"] = common.Satellite0Identity
	s.variables["satellite-core"]["identity"] = common.Satellite0Identity
	s.variables["satellite-admin"]["identity"] = common.Satellite0Identity
	return s, nil
}

func (c *Standalone) generateIdentity(name string, index int) error {

	serviceDir := filepath.Join(c.dir, name, strconv.Itoa(index))

	caCertPath := filepath.Join(serviceDir, "ca.cert")
	caKeyPath := filepath.Join(serviceDir, "ca.key")
	identCertPath := filepath.Join(serviceDir, "identity.cert")
	identKeyPath := filepath.Join(serviceDir, "identity.key")

	// we use hardcoded identity for satellite-api for predictable grants.
	if name == "satellite-api" && index == 0 {
		return errs.Combine(
			os.MkdirAll(serviceDir, 0755),
			os.WriteFile(identCertPath, identityCert, 0644),
			os.WriteFile(identKeyPath, identityKey, 0644),
			os.WriteFile(caCertPath, caCert, 0644),
			os.WriteFile(caKeyPath, caKey, 0644),
		)

	}

	caConfig := identity.CASetupConfig{
		CertPath:      caCertPath,
		KeyPath:       caKeyPath,
		Difficulty:    0,
		Concurrency:   4,
		VersionNumber: 0,
	}

	status, err := caConfig.Status()
	if err != nil {
		return err
	}
	if status != identity.NoCertNoKey {
		return errs.Errorf("CA certificate and/or key already exists, NOT overwriting!")
	}

	identConfig := identity.SetupConfig{
		CertPath: identCertPath,
		KeyPath:  identKeyPath,
	}

	status, err = identConfig.Status()
	if err != nil {
		return err
	}
	if status != identity.NoCertNoKey {
		return errs.Errorf("Identity certificate and/or key already exists, NOT overwriting!")
	}

	ca, caerr := caConfig.Create(context.Background(), os.Stdout)
	if caerr != nil {
		return caerr
	}

	_, iderr := identConfig.Create(ca)
	if iderr != nil {
		return iderr
	}

	return nil
}

func (c *Standalone) uniqueName(s *service) string {
	u := ""
	if c.serviceCount(s.id.Name) > 1 {
		u += strconv.Itoa(s.id.Instance + 1)
	}
	return s.id.Name + u
}

func find(root, ext string) ([]string, error) {
	var a []string
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}
