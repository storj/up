// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveServices(t *testing.T) {
	res, err := ResolveServices([]string{"db", "satellite-api"})
	require.NoError(t, err)
	require.Equal(t, []string{"spanner", "redis", "satellite-api"}, res)
}
