// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolveService(t *testing.T) {
	services, err := ResolveServices([]string{"minimal", "db"})
	require.NoError(t, err)
	require.Equal(t,
		[]string{"cockroach", "redis", "satellite-api", "storagenode"},
		services)
}
