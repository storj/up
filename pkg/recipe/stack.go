// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package recipe

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zeebo/errs/v2"
)

// Stack is the list of known recipes.
type Stack []Recipe

// GetEmbeddedStack returns with a stack includes only embedded recipes.
func GetEmbeddedStack() (Stack, error) {
	res := Stack{}
	for name, src := range Defaults {
		r, err := Read(src)
		if err != nil {
			return res, errs.Errorf("Error on reading %s (embedded recipe), %v", name, err)
		}
		res = append(res, r)
	}
	return res, nil
}

// GetStack returns Stack with all the known recipe definitions.
func GetStack() (Stack, error) {
	res, err := GetEmbeddedStack()
	if err != nil {
		return res, err
	}

	recipeDir := ""
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		recipeDir = filepath.Join(xdgConfigHome, "storj-up", "recipes")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			recipeDir = filepath.Join(home, ".config", "storj-up", "recipes")
		}
	}
	if _, err := os.Stat(recipeDir); err != nil {
		// nolint:nilerr
		return res, nil
	}
	configDir, err := os.ReadDir(recipeDir)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	for _, e := range configDir {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
			receipeFile := filepath.Join(recipeDir, e.Name())
			content, err := os.ReadFile(receipeFile)
			if err != nil {
				return nil, errs.Wrap(err)
			}
			r, err := Read(content)
			if err != nil {
				return res, errs.Errorf("Error on reading %s %v", receipeFile, err)
			}
			res = append(res, r)
		}
	}
	return res, nil
}

// Get returns with first recipe based on the the name.
func (s Stack) Get(name string) (Recipe, error) {
	for _, receipe := range s {
		if receipe.Name == name {
			return receipe, nil
		}
	}
	return Recipe{}, errs.Errorf("No such recipe %s", name)
}

// FindRecipeByName returns with the recipe based on a name.
func (s Stack) FindRecipeByName(name string) (*Service, error) {
	for _, r := range s {
		for _, service := range r.Add {
			if service.Name == name {
				return service, nil
			}
		}
	}
	return nil, errs.Errorf("Couldn't find recipe which includes service %s", name)
}

// AllRecipeNames returns with all internal recipe names.
func (s Stack) AllRecipeNames() []string {
	var k []string
	for _, r := range s {
		k = append(k, r.Name)
	}
	return k
}
