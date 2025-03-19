// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import (
	"strings"

	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/recipe"
)

// InitFromRecipe can fill standard fields of the service (flags, configs, ...) based on the recipe.
func InitFromRecipe(service Service, recipe recipe.Service) error {
	err := service.ChangeImage(func(s string) string {
		return recipe.Image
	})
	if err != nil {
		return err
	}

	for k, v := range recipe.Config {
		err := service.AddConfig(k, v)
		if err != nil {
			return err
		}
	}

	for k, v := range recipe.Environment {
		err := service.AddEnvironment(k, v)
		if err != nil {
			return err
		}
	}

	for _, v := range recipe.Command {
		err := service.AddFlag(v)
		if err != nil {
			return err
		}
	}

	for _, port := range recipe.Port {
		err := service.AddPortForward(PortMap{Internal: port.Target, External: port.Target, Protocol: port.Protocol})
		if err != nil {
			return err
		}
	}
	for _, f := range recipe.File {
		err := service.UseFile(f.Path, f.Name, f.Data)
		if err != nil {
			return err
		}
	}
	for _, f := range recipe.Folder {
		err := service.UseFolder(f.Path, f.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

// ModifyFromRecipe applies the modification defined by a recipe to a service.
func ModifyFromRecipe(service Service, mod recipe.Modification) error {
	for _, f := range mod.Flag.Add {
		err := service.AddFlag(f)
		if err != nil {
			return err
		}
	}
	for k, v := range mod.Config {
		err := service.AddConfig(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Match checks if matcher selects the given service.
func Match(service Service, matcher recipe.Matcher) bool {
	if service == nil {
		panic("asd")
	}
	for _, l := range service.Labels() {
		for _, m := range matcher.Label {
			if l == m {
				return true
			}
		}
	}
	for _, name := range strings.Split(matcher.Name, ",") {
		if strings.TrimSpace(name) == service.ID().Name {
			return true
		}
	}
	return false
}

// ApplyRecipes can apply full recipes and other services (partial recipes) based on the selectors.
func ApplyRecipes(st recipe.Stack, rt Runtime, selector []string, instanceOverride int) error {
	// First, collect all recipes to apply
	type recipeToApply struct {
		name     string
		recipe   recipe.Recipe
		priority int
		isRecipe bool
	}

	// Collect all recipes to apply, including their priorities
	toApply := []recipeToApply{}

	for _, name := range selector {
		rcp, err := st.Get(name)
		if err == nil {
			// It's a recipe, store it
			toApply = append(toApply, recipeToApply{
				name:     name,
				recipe:   rcp,
				priority: rcp.Priority,
				isRecipe: true,
			})
			continue
		}

		// Check if it's an individual service
		found := false
		for _, r := range st {
			for _, s := range r.Add {
				if s.Name == name {
					found = true
					toApply = append(toApply, recipeToApply{
						name:     name,
						priority: 0, // Individual services default to priority 0
						isRecipe: false,
					})
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			return errs.Errorf("Couldn't find recipe or service in any recipe with the name %s. Please execute `storj-up services` to list available recipes/services", name)
		}
	}

	// Sort by priority (higher priority first)
	// Stable sort preserves original order within same priority
	for i := 0; i < len(toApply); i++ {
		for j := i + 1; j < len(toApply); j++ {
			if toApply[i].priority < toApply[j].priority {
				toApply[i], toApply[j] = toApply[j], toApply[i]
			}
		}
	}

	// Apply recipes in priority order
	for _, item := range toApply {
		if item.isRecipe {
			err := ApplyRecipeToRuntime(rt, item.recipe)
			if err != nil {
				return err
			}
		} else {
			// It's an individual service
			added := 0
			for _, r := range st {
				for _, s := range r.Add {
					if s.Name == item.name {
						if instanceOverride != 0 {
							s.Instance = instanceOverride
						}

						err := AddServiceToRuntime(rt, *s)
						if err != nil {
							return errs.Wrap(err)
						}

						added++
					}
				}
			}
			if added == 0 {
				return errs.Errorf("Couldn't find service with name %s", item.name)
			}
		}
	}

	return nil
}

// ApplyRecipeToRuntime can add all services from recipe and modifies existing ones based on rules.
func ApplyRecipeToRuntime(c Runtime, r recipe.Recipe) error {
	for _, s := range r.Add {
		err := AddServiceToRuntime(c, *s)
		if err != nil {
			return err
		}
	}
	for _, mod := range r.Modify {
		for _, service := range c.GetServices() {
			if Match(service, mod.Match) {
				err := ModifyFromRecipe(service, *mod)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// AddServiceToRuntime helps to create new service in runtime based on a recipe.
func AddServiceToRuntime(c Runtime, r recipe.Service) error {
	instance := r.Instance
	if instance == 0 {
		instance = 1
	}
	for i := 0; i < instance; i++ {
		_, err := c.AddService(r)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}
