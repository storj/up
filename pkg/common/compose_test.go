package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ResolveService(t *testing.T) {
	require.Equal(t,
		[]string{"cockroach", "redis", "satellite-api", "storagenode"},
		ResolveServices([]string{"minimal", "db"}))
}
