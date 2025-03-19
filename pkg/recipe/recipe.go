// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package recipe

import (
	"github.com/zeebo/errs/v2"
	"gopkg.in/yaml.v3"
)

// Recipe represents all configuration/runtime definition for a service.
type Recipe struct {
	Name        string
	Description string
	// Higher priority recipes will be applied first, default is 0
	Priority int
	Add      []*Service
	Modify   []*Modification
}

// Service contains all the parameters to run one service.
type Service struct {
	Name          string
	ContainerName string
	Label         []string
	Instance      int
	Image         string
	Command       []string
	Environment   map[string]string
	Config        map[string]string
	Persistence   []string
	Port          []PortDefinition
	File          []File
	Folder        []Folder

	// port forward outside->inside
	PortForwards map[int]int

	// mount rules outside --> inside
	Mounts map[string]string
}

// Modification represents a transformation applied to one or more services.
type Modification struct {
	Match       Matcher
	Flag        FlagModification
	Config      map[string]string
	Environment map[string]string
}

// File represents any configuration file required by a recipe.
type File struct {
	Name string
	Path string
	Data string
}

// Folder represents any configuration folder required by a recipe.
type Folder struct {
	Name string
	Path string
}

// FlagModification represents modification for service command / flags.
type FlagModification struct {
	Add    []string
	Remove []string
}

// Matcher can identify the services to be modified by a recipe.
type Matcher struct {
	Label []string
	Name  string
}

// PortDefinition gives information about used ports.
type PortDefinition struct {
	Name        string
	Description string
	Target      int
	Protocol    string
}

// HasLabel checks if the service has one specific label.
func (s Service) HasLabel(s2 string) bool {
	for _, l := range s.Label {
		if l == s2 {
			return true
		}
	}
	return false
}

// Read loads a recipe from a yaml file.
func Read(data []byte) (Recipe, error) {
	r := Recipe{}
	err := yaml.Unmarshal(data, &r)
	if err != nil {
		return Recipe{}, errs.Wrap(err)
	}
	return r, nil
}
