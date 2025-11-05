// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/storj-up/cmd/testdata"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func TestPersistCockroach(t *testing.T) {
	dir := t.TempDir()

	st, rt, err := testdata.InitCompose(dir)
	require.NoError(t, err)

	err = runtime.ApplyRecipes(st, rt, []string{"cockroach"}, 0)
	require.NoError(t, err)

	err = persist(st, rt, []string{"cockroach"})
	require.NoError(t, err)

	require.NoError(t, rt.Write())

	result, err := os.ReadFile(filepath.Join(dir, "docker-compose.yaml"))
	require.NoError(t, err)

	require.Contains(t, string(result), "source: "+filepath.Join(dir, "/cockroach/cockroach-data"))
	require.Contains(t, string(result), "target: /cockroach/cockroach-data")

}
