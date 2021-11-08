package common

import (
	"github.com/stretchr/testify/require"
	"testing"
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

	services, err = ResolveServices([]string{"enterprise"})
	require.Error(t, err)
}
