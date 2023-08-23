// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import (
	"testing"

	"github.com/stretchr/testify/require"

	"storj.io/storj-up/pkg/recipe"
)

func TestInitFromRecipe(t *testing.T) {
	s := NewMockService("satellite-api")
	r := recipe.Service{
		Image:   "img.dev.storj.io/storjup/storj:latest",
		Command: []string{"ls", "-lah"},
		Config: map[string]string{
			"conf1": "val1",
		},
	}
	err := InitFromRecipe(s, r)
	require.NoError(t, err)

	require.Equal(t, "img.dev.storj.io/storjup/storj:latest", s.Image)
	require.Equal(t, []string{"ls", "-lah"}, s.Flag)
	require.Equal(t, "val1", s.Config["conf1"])
}

func TestMatch(t *testing.T) {
	s := NewMockService("satellite-api")
	s.Label = []string{"asd", "storj"}
	matcher := recipe.Matcher{
		Label: []string{
			"storj",
		},
	}
	require.True(t, Match(s, matcher))

}

func TestModifyFromRecipe(t *testing.T) {
	s := NewMockService("satellite-api")
	s.Config["conf1"] = "xxx"

	err := ModifyFromRecipe(s, recipe.Modification{
		Flag: recipe.FlagModification{
			Add: []string{"--flag2=xxx"},
		},
		Config: map[string]string{
			"conf1": "yyy",
		},
	})
	require.NoError(t, err)

	require.Equal(t, "yyy", s.Config["conf1"])
	require.Equal(t, []string{"--flag2=xxx"}, s.Flag)
}

func TestApplyRecipeCreateByRecipe(t *testing.T) {
	rt := NewMockRuntime()
	st := recipe.Stack{
		{
			Name: "minimal",
			Add: []*recipe.Service{
				{
					Name: "satellite-api",
					Config: map[string]string{
						"conf1": "val1",
					},
				},
			},
		},
	}
	err := ApplyRecipes(st, rt, []string{"minimal"}, 0)
	require.NoError(t, err)

	require.Len(t, rt.Services, 1)
	require.Equal(t, "val1", rt.Services[0].(*MockService).Config["conf1"])

}

func TestApplyRecipeModifyExisting(t *testing.T) {
	rt := NewMockRuntime()
	db := NewMockService("db")
	db.Label = []string{"db"}
	rt.Services = append(rt.Services, db)
	st := recipe.Stack{
		{
			Name: "minimal",
			Add: []*recipe.Service{
				{
					Name: "satellite-api",
					Config: map[string]string{
						"conf1": "val1",
					},
				},
			},
			Modify: []*recipe.Modification{
				{
					Match: recipe.Matcher{
						Label: []string{"db"},
					},
					Flag: recipe.FlagModification{
						Add: []string{"--new=sg"},
					},
				},
			},
		},
	}
	err := ApplyRecipes(st, rt, []string{"minimal"}, 0)
	require.NoError(t, err)

	require.Len(t, rt.Services, 2)
	require.Equal(t, []string{"--new=sg"}, rt.Services[0].(*MockService).Flag)

}

func TestApplyRecipeModify(t *testing.T) {
	s := NewMockService("satellite-api")
	s.Config["conf1"] = "xxx"

	err := ModifyFromRecipe(s, recipe.Modification{
		Flag: recipe.FlagModification{
			Add: []string{"--flag2=xxx"},
		},
		Config: map[string]string{
			"conf1": "yyy",
		},
	})
	require.NoError(t, err)

	require.Equal(t, "yyy", s.Config["conf1"])
	require.Equal(t, []string{"--flag2=xxx"}, s.Flag)
}
