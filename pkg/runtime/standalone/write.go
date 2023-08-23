// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package standalone

import (
	_ "embed"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/zeebo/errs/v2"
)

//go:embed intellij.xml
var intelliJTemplate []byte

//go:embed start.sh
var startTemplate []byte

//go:embed supervisord.template
var supervisorTemplate []byte

//go:embed .env
var dotEnvrc []byte

func (c *Standalone) Write() error {
	_ = os.MkdirAll(c.dir, 0755)
	for _, service := range c.services {
		err := c.writeService(service)
		if err != nil {
			return err
		}
		err = c.writeIntelliJRunner(service)
		if err != nil {
			return err
		}
		if len(service.config) > 0 {
			err := os.WriteFile(filepath.Join(c.dir, service.id.Name, strconv.Itoa(service.id.Instance), "config.yaml"), []byte(strings.Join(service.config, "\n")), 0644)
			if err != nil {
				return errs.Wrap(err)
			}
		}
	}
	err := c.writeSupervisor()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(c.dir, ".envrc"), dotEnvrc, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Standalone) writeSupervisor() error {
	f, err := os.OpenFile(filepath.Join(c.dir, "supervisord.conf"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	t, err := template.New("supervisord.template").
		Funcs(map[string]interface{}{
			"UniqueName": c.uniqueName,
			"Add": func(a int, b int) int {
				return a + b
			},
		}).
		Parse(string(supervisorTemplate))
	if err != nil {
		return errs.Wrap(err)
	}

	err = t.Execute(f, struct {
		Services []*service
	}{
		Services: c.services,
	})
	return err
}

func (c *Standalone) writeService(s *service) error {
	f, err := os.OpenFile(filepath.Join(c.dir, c.uniqueName(s)+".sh"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	t, err := template.New("start.sh").
		Funcs(map[string]interface{}{
			"HasPrefix": strings.HasPrefix,
			"Safe": func(p string) string {
				out := ""
				for i := 0; i < len(p); i++ {
					if i > 0 && p[i] == '"' && p[i-1] != '\\' {
						out += "\\"
					}
					out += string(p[i])
				}
				return out
			},
		}).
		Parse(string(startTemplate))
	if err != nil {
		return errs.Wrap(err)
	}

	err = t.Execute(f, struct {
		Service *service
	}{
		Service: s,
	})
	return err
}

var runnerSupported = map[string]string{
	"satellite-api": "storj",
	"satellite-gc":  "storj",
	"satellite":     "storj",
	"storagenode":   "storj",
	"gateway-mt":    "gateway-mt",
	"authservice":   "gateway-mt",
	"linksharing":   "gateway-mt",
}

func (c *Standalone) writeIntelliJRunner(s *service) error {
	pkg, found := runnerSupported[s.ID().Name]
	if !found {
		return nil
	}
	f, err := os.OpenFile(filepath.Join(c.dir, c.uniqueName(s)+".run.xml"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	t, err := template.New("intellij.xml").
		Funcs(map[string]interface{}{
			"HasPrefix":  strings.HasPrefix,
			"UniqueName": c.uniqueName,
			"Tail": func(a []string) []string {
				if len(a) > 0 {
					return a[1:]
				}
				return a
			},
			"Join": strings.Join,
			"Safe": func(p string) string {
				return strings.ReplaceAll(p, "\"", "&quot;")
			},
		}).
		Parse(string(intelliJTemplate))
	if err != nil {
		return errs.Wrap(err)
	}

	executable := s.Command[0]
	if strings.HasPrefix(executable, "satellite") {
		executable = "satellite"
	}
	err = t.Execute(f, struct {
		Service    *service
		Package    string
		Executable string
	}{
		Service:    s,
		Package:    pkg,
		Executable: executable,
	})
	return err
}
