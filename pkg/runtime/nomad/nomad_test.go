// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package nomad

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func TestPersist(t *testing.T) {

	dir := t.TempDir()
	c, err := NewNomad(dir, "test")
	require.NoError(t, err)

	r := recipe.Service{
		Name:  "satellite-api",
		Image: "img.dev.storj.io/storjup/storj",
	}

	s, err := c.AddService(r)
	require.NoError(t, err)

	err = s.Persist("/some/dir")
	require.NoError(t, err)

	err = c.Write()
	require.NoError(t, err)

	file, err := os.ReadFile(filepath.Join(dir, "storj.hcl"))
	require.NoError(t, err)
	require.Contains(t, string(file), "volumes      = [\"/tmp/satellite-api/0/dir:/some/dir\"]")

}

func TestNomad(t *testing.T) {
	tempDir := t.TempDir()
	defer os.Remove(tempDir) //nolint:errcheck

	rt, err := NewNomad(tempDir, "test")
	require.NoError(t, err)

	st, err := recipe.GetEmbeddedStack()
	require.NoError(t, err)

	err = runtime.ApplyRecipes(st, rt, st.AllRecipeNames(), 0)
	require.NoError(t, err)
}
