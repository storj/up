package recipe

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScale(t *testing.T) {

	k := &common.SimplifiedCompose{
		Services: map[string]*common.ServiceConfig{
			"storagenode": &common.ServiceConfig{
				Image: "foobar",
			},
			"satellite-api": &common.ServiceConfig{},
		},
	}

	err := Scale("storagenode", k, 3)
	require.Nil(t, err)

	expected := &common.SimplifiedCompose{
		Services: map[string]*common.ServiceConfig{
			"storagenode1": &common.ServiceConfig{
				Image: "foobar",
			},
			"storagenode2": &common.ServiceConfig{
				Image: "foobar",
			},
			"storagenode3": &common.ServiceConfig{
				Image: "foobar",
			},
			"satellite-api": &common.ServiceConfig{},
		},
	}

	require.Equal(t, expected, k)
}

func TestScaleDown(t *testing.T) {

	k := &common.SimplifiedCompose{
		Services: map[string]*common.ServiceConfig{
			"storagenode1": &common.ServiceConfig{
				Image: "foobar",
			},
			"storagenode2": &common.ServiceConfig{
				Image: "foobar",
			},
			"storagenode3": &common.ServiceConfig{
				Image: "foobar",
			},
			"satellite-api": &common.ServiceConfig{},
		},
	}

	err := Scale("storagenode", k, 1)
	require.Nil(t, err)

	expected := &common.SimplifiedCompose{
		Services: map[string]*common.ServiceConfig{
			"storagenode": &common.ServiceConfig{
				Image: "foobar",
			},
			"satellite-api": &common.ServiceConfig{},
		},
	}

	require.Equal(t, expected, k)
}
