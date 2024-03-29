// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package nomad

import (
	"path/filepath"
	"strconv"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zeebo/errs"

	"storj.io/storj-up/pkg/runtime/runtime"
)

type service struct {
	task   *hclwrite.Block
	id     runtime.ServiceInstance
	env    *hclwrite.Block
	cmds   []string
	render func(string) (string, error)
	labels []string
}

var _ runtime.Service = (*service)(nil)

func (s *service) GetENV() map[string]*string {
	// TODO implement me
	return nil
}

func (s *service) GetVolumes() []runtime.VolumeMount {
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
	return errs.New("RemoveFlag for Nomad is not yet implemented")
}

func (s *service) Persist(dir string) error {
	cfg := s.task.Body().FirstMatchingBlock("config", []string{})
	sourceDir := filepath.Join("/tmp", s.id.Name, strconv.Itoa(s.ID().Instance), filepath.Base(dir))
	cfg.Body().SetAttributeValue("volumes", cty.ListVal(
		[]cty.Value{
			cty.StringVal(sourceDir + ":" + dir),
		}))
	return nil
}

func (s *service) ChangeImage(change func(string) string) error {
	cfg := s.task.Body().FirstMatchingBlock("config", []string{})
	cfg.Body().SetAttributeValue("image", cty.StringVal(change("todo")))
	return nil
}

func (s *service) AddPortForward(runtime.PortMap) error {
	return nil
}

// RemovePortForward implements runtime.Service.
func (s *service) RemovePortForward(runtime.PortMap) error {
	return nil
}

func (s *service) ID() runtime.ServiceInstance {
	return s.id
}

func (s *service) AddConfig(key string, value string) error {
	v, err := s.render(value)
	s.env.Body().SetAttributeValue(key, cty.StringVal(v))
	return err
}

func (s *service) AddFlag(flag string) error {
	f, err := s.render(flag)
	s.cmds = append(s.cmds, f)
	s.updateCmds()
	return err
}

func (s *service) AddEnvironment(key string, value string) error {
	v, err := s.render(value)
	s.env.Body().SetAttributeValue(key, cty.StringVal(v))
	return err
}

func (s *service) updateCmds() {
	if len(s.cmds) > 0 {
		s.task.Body().FirstMatchingBlock("config", []string{}).Body().SetAttributeValue("args", toCtyStrList(s.cmds))
	}
}
