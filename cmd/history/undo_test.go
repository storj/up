// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/compose-spec/compose-go/types"
	"github.com/stretchr/testify/require"

	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func Test_Undo(t *testing.T) {

	dir := t.TempDir()

	rt, err := compose.NewCompose(dir)
	require.NoError(t, err)

	st, err := recipe.GetStack()
	require.NoError(t, err)

	err = runtime.ApplyRecipes(st, rt, []string{"satellite-api"}, 0)
	require.NoError(t, err)

	// first modification
	err = rt.GetServices()[0].AddConfig("foo", "bar")
	require.NoError(t, err)
	require.NoError(t, rt.Write())

	// second modification
	err = rt.GetServices()[0].AddConfig("foo", "vok")
	require.NoError(t, err)
	require.NoError(t, rt.Write())

	// reload
	rt, err = compose.NewCompose(dir)
	require.NoError(t, err)
	require.NoError(t, rt.Reload(st))

	// check the current value
	err = rt.GetServices()[0].(*compose.Service).TransformRaw(func(config *types.ServiceConfig) error {
		require.Equal(t, "vok", *config.Environment["foo"])
		return nil
	})
	require.NoError(t, err)

	// revert
	reverted, err := common.Store.RestoreLatestVersion()
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir, common.ComposeFileName), reverted, 0644)
	require.NoError(t, err)

	// reload
	rt, err = compose.NewCompose(dir)
	require.NoError(t, err)
	require.NoError(t, rt.Reload(st))

	// check the current value
	err = rt.GetServices()[0].(*compose.Service).TransformRaw(func(config *types.ServiceConfig) error {
		require.Equal(t, "bar", *config.Environment["foo"])
		return nil
	})
	require.NoError(t, err)

}
