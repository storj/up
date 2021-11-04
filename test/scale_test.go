package test

import (
	"strconv"
	"testing"

	"github.com/compose-spec/compose-go/types"
	"github.com/stretchr/testify/require"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
)

func TestScale(t *testing.T) {

	k := &types.Project{
		Services: []types.ServiceConfig{
			{Name: "storagenode",
				Image: "foobar",
			},
			{Name: "satellite-api"},
		},
	}

	var scale uint64 = 3
	actual, err := common.UpdateEach(k, cmd.Scale, strconv.Itoa(int(scale)), []string{"storagenode"})
	require.Nil(t, err)

	expected := &types.Project{
		Services: []types.ServiceConfig{
			{Name: "storagenode",
				Image: "foobar",
				Deploy: &types.DeployConfig{
					Mode:     "",
					Replicas: &scale,
				},
			},
			{Name: "satellite-api"},
		},
	}

	require.Equal(t, expected, actual)
}

func TestScaleDown(t *testing.T) {

	var scaleUp uint64 = 3
	var scaleDown uint64 = 1
	k := &types.Project{
		Services: []types.ServiceConfig{
			{Name: "storagenode",
				Image: "foobar",
				Deploy: &types.DeployConfig{
					Mode:     "",
					Replicas: &scaleUp,
				},
			},
			{Name: "satellite-api"},
		},
	}

	actual, err := common.UpdateEach(k, cmd.Scale, strconv.Itoa(int(scaleDown)), []string{"storagenode"})
	require.Nil(t, err)

	expected := &types.Project{
		Services: []types.ServiceConfig{
			{Name: "storagenode",
				Image: "foobar",
			},
			{Name: "satellite-api"},
		},
	}

	require.Equal(t, expected, actual)
}
