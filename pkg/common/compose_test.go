package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ResolveService(t *testing.T) {

	res, err := ResolveServices([]string{"minimal", "db"})
	require.NoError(t, err)
	require.Equal(t, []string{"cockroach", "redis", "satellite-api", "storagenode"}, res)

}
