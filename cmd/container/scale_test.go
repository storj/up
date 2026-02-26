// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package container

import (
	"testing"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/stretchr/testify/require"
)

func TestScale(t *testing.T) {
	k := types.ServiceConfig{
		Name:  "storagenode",
		Image: "foobar",
	}

	err := scale(&k, "10")
	require.NoError(t, err)

	require.Equal(t, 10, *k.Deploy.Replicas)
}
