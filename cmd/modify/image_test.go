// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/storj-up/cmd/testdata"
)

func TestSetImage(t *testing.T) {
	dir := t.TempDir()

	st, rt, err := testdata.InitCompose(dir)
	require.NoError(t, err)

	err = setImage(st, rt, []string{"cockroach", "xxx"})
	require.NoError(t, err)

	require.NoError(t, rt.Write())

	result, err := os.ReadFile(filepath.Join(dir, "docker-compose.yaml"))
	require.NoError(t, err)

	require.Contains(t, string(result), "image: xxx")

}
