// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"testing"

	"github.com/compose-spec/compose-go/types"
	"github.com/stretchr/testify/require"

	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func TestAddPort(t *testing.T) {
	dir := t.TempDir()
	st := recipe.Stack([]recipe.Recipe{
		{
			Name: "base",
			Add: []*recipe.Service{
				{
					Port: []recipe.PortDefinition{
						{
							Name:     "test port",
							Target:   80,
							Protocol: "tcp",
						},
					},
				},
			},
		},
	})

	rt := compose.NewEmptyCompose(dir)
	err := runtime.ApplyRecipes(st, rt, []string{"base"})
	require.NoError(t, err)

	err = addPort(st, rt, []string{"base", "90"})
	require.NoError(t, err)

	s := rt.GetServices()[0]
	rawService := s.(*compose.Service)
	err = rawService.TransformRaw(func(config *types.ServiceConfig) error {
		require.Equal(t, []types.ServicePortConfig{{
			Mode:       "ingress",
			HostIP:     "",
			Target:     80,
			Published:  80,
			Protocol:   "tcp",
			Extensions: nil,
		}, {
			Mode:       "ingress",
			HostIP:     "",
			Target:     90,
			Published:  90,
			Protocol:   "tcp",
			Extensions: nil,
		}}, config.Ports)
		return nil
	})
	require.NoError(t, err)

}

func TestRemovePort(t *testing.T) {
	dir := t.TempDir()
	st := recipe.Stack([]recipe.Recipe{
		{
			Name: "base",
			Add: []*recipe.Service{
				{
					Port: []recipe.PortDefinition{
						{
							Name:     "test port 1",
							Target:   80,
							Protocol: "tcp",
						},
						{
							Name:     "test port 2",
							Target:   90,
							Protocol: "tcp",
						},
					},
				},
			},
		},
	})

	rt := compose.NewEmptyCompose(dir)
	err := runtime.ApplyRecipes(st, rt, []string{"base"})
	require.NoError(t, err)

	err = removePort(st, rt, []string{"base", "80"})
	require.NoError(t, err)

	s := rt.GetServices()[0]
	rawService := s.(*compose.Service)
	err = rawService.TransformRaw(func(config *types.ServiceConfig) error {
		require.Equal(t, []types.ServicePortConfig{{
			Mode:       "ingress",
			HostIP:     "",
			Target:     90,
			Published:  90,
			Protocol:   "tcp",
			Extensions: nil,
		}}, config.Ports)
		return nil
	})
	require.NoError(t, err)
}
