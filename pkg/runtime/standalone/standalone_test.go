// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package standalone

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/common/identity"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func TestIdentity(t *testing.T) {
	nodeID, err := identity.NodeIDFromCertPath("identity.cert")
	require.NoError(t, err)
	require.Equal(t, common.Satellite0Identity, nodeID.String())
}

func TestStandalone(t *testing.T) {
	tempDir := t.TempDir()
	defer os.Remove(tempDir) //nolint:errcheck
	rt, err := NewStandalone(Paths{
		ScriptDir:  tempDir,
		StorjDir:   tempDir,
		GatewayDir: tempDir,
		CleanDir:   false,
	})
	require.NoError(t, err)

	st, err := recipe.GetEmbeddedStack()
	require.NoError(t, err)

	err = runtime.ApplyRecipes(st, rt, st.AllRecipeNames())
	require.NoError(t, err)
}
