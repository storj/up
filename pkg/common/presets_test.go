// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolveServices(t *testing.T) {
	services, err := ResolveServices([]string{"gateway-mt"})
	require.NoError(t, err)
	require.Equal(t, []string{"gateway-mt"}, services)

	services, err = ResolveServices([]string{"gatewaymt"})
	require.NoError(t, err)
	require.Equal(t, []string{"gateway-mt"}, services)

	services, err = ResolveServices([]string{"minimal"})
	require.NoError(t, err)
	require.Equal(t, []string{"satellite-api", "storagenode"}, services)

	_, err = ResolveServices([]string{"enterprise"})
	require.Error(t, err)
}
