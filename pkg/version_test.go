package sjr

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpper(t *testing.T) {
	require.Equal(t, "KEY_PATH", upper("keyPath"))
}
