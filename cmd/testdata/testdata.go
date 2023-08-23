// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package testdata

import (
	_ "embed"
	"os"
	"path/filepath"

	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/runtime"
)

//go:embed docker-compose.yaml
var composeFile []byte

func InitCompose(dir string) (st recipe.Stack, rt runtime.Runtime, err error) {
	err = os.WriteFile(filepath.Join(dir, "docker-compose.yaml"), composeFile, 0644)
	if err != nil {
		return
	}
	rt, err = compose.NewCompose(dir)
	if err != nil {
		return
	}

	st, err = recipe.GetStack()
	if err != nil {
		return
	}
	err = runtime.ApplyRecipes(st, rt, []string{"db", "minimal"}, 0)
	return
}
