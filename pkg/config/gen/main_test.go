// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpper(t *testing.T) {
	require.Equal(t, "KEY_PATH", camelToUpperCase("keyPath"))
}
