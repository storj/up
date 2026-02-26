// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"testing"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/stretchr/testify/require"

	recipe2 "storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func TestAddFlag(t *testing.T) {
	dir := t.TempDir()
	st := recipe2.Stack([]recipe2.Recipe{
		{
			Name: "base",
			Add: []*recipe2.Service{
				{
					Command: []string{
						"one", "--flag=sg",
					},
				},
			},
		},
	})

	rt := compose.NewEmptyCompose(dir)
	err := runtime.ApplyRecipes(st, rt, []string{"base"}, 0)
	require.NoError(t, err)

	err = addFlag(st, rt, []string{"base", "nf=2"})
	require.NoError(t, err)

	s := rt.GetServices()[0]
	rawService := s.(*compose.Service)
	err = rawService.TransformRaw(func(config *types.ServiceConfig) error {
		require.Equal(t, types.ShellCommand{"one", "--flag=sg", "--nf=2"}, config.Command)
		return nil
	})
	require.NoError(t, err)

}

func TestRemoveFlag(t *testing.T) {
	dir := t.TempDir()
	st := recipe2.Stack([]recipe2.Recipe{
		{
			Name: "base",
			Add: []*recipe2.Service{
				{
					Command: []string{
						"one", "--flag=sg",
					},
				},
			},
		},
	})

	rt := compose.NewEmptyCompose(dir)
	err := runtime.ApplyRecipes(st, rt, []string{"base"}, 0)
	require.NoError(t, err)

	err = removeFlag(st, rt, []string{"base", "flag"})
	require.NoError(t, err)

	s := rt.GetServices()[0]
	rawService := s.(*compose.Service)
	err = rawService.TransformRaw(func(config *types.ServiceConfig) error {
		require.Equal(t, types.ShellCommand{"one"}, config.Command)
		return nil
	})
	require.NoError(t, err)

}
