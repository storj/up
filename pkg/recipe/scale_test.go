package recipe

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/cmd"
	"github.com/elek/sjr/pkg/common"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestScale(t *testing.T) {

	k := &common.ComposeFile{
		Version: "3.4",
		Services: []types.ServiceConfig{
			{Name: "storagenode",
			Image: "foobar",
			},
			{Name: "satellite-api"},
		},
	}

	var scale uint64 = 3
	err := cmd.Scale([]string{strconv.Itoa(int(scale)), "storagenode"})
	require.Nil(t, err)

	expected := &common.ComposeFile{
		Version: "3.4",
		Services: []types.ServiceConfig{
			{Name: "storagenode",
				Image: "foobar",
				Deploy: &types.DeployConfig{
					Replicas: &scale,
				},
			},
			{Name: "satellite-api"},
		},
	}

	require.Equal(t, expected, k)
}

func TestScaleDown(t *testing.T) {

	var scaleUp uint64 = 3
	var scaleDown uint64 = 1
	k := &common.ComposeFile{
		Version: "3.4",
		Services: []types.ServiceConfig{
			{Name: "storagenode",
				Image: "foobar",
				Deploy: &types.DeployConfig{
					Replicas: &scaleUp,
				},
			},
			{Name: "satellite-api"},
		},
	}

	err := cmd.Scale([]string{strconv.Itoa(int(scaleDown)), "storagenode"})
	require.Nil(t, err)

	expected := &common.ComposeFile{
		Version: "3.4",
		Services: []types.ServiceConfig{
			{Name: "storagenode",
				Image: "foobar",
			},
			{Name: "satellite-api"},
		},
	}

	require.Equal(t, expected, k)
}
