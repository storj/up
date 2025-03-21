// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func init() {
	cmd.RootCmd.AddCommand(&cobra.Command{
		Use:   "persist <selector>...",
		Short: "Make internal state (database files, storagenode files) persisted between restarts. ",
		Long:  "This is done usually with mounting the directory to the houst. ." + cmd.SelectorHelp,
		Args:  cobra.MinimumNArgs(1),
		RunE:  cmd.ExecuteStorjUP(persist),
	})
}

func persist(st recipe.Stack, rt runtime.Runtime, selectors []string) error {
	return runtime.ModifyService(st, rt, selectors, func(s runtime.Service) error {
		rService, err := st.FindRecipeByName(s.ID().Name)
		if err != nil {
			return err
		}

		if rService.Persistence != nil {
			for _, p := range rService.Persistence {
				err := s.Persist(p)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
