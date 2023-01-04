// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import (
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
		err := service.AddPortForward(PortMap{Internal: port.Target, External: port.Target})
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
	return matcher.Name == service.ID().Name
}

// ApplyRecipes can apply full recipes and other services (partial recipes) based on the selectors.
func ApplyRecipes(st recipe.Stack, rt Runtime, selector []string) error {
	for _, name := range selector {
		rcp, err := st.Get(name)
		if err == nil {
			// it's a recipe, apply it fully
			err = ApplyRecipeToRuntime(rt, rcp)
			if err != nil {
				return err
			}
			continue
		}

		added := 0
		for _, r := range st {
			for _, s := range r.Add {
				if s.Name == name {
					instance := s.Instance
					if instance == 0 {
						instance = 1
					}

					for i := 1; i <= instance; i++ {
						err = AddServiceToRuntime(rt, *s)
						if err != nil {
							return errs.Wrap(err)
						}
					}
					added++
				}

			}

		}
		if added > 0 {
			continue
		}

		return errs.Errorf("Couldn't find recipe or service in any recipe with the name %s. Please execute `storj-up services` to list available recipes/services", name)

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
