// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func imageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "image <selector> <image>",
		Short: "Change container image for one more more services",
		Long:  "Change image of specified services." + cmd.SelectorHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE:  cmd.ExecuteStorjUP(setImage),
	}
}

func init() {
	cmd.RootCmd.AddCommand(imageCmd())
}

func setImage(st recipe.Stack, rt runtime.Runtime, args []string) error {
	return runtime.ModifyService(st, rt, args[:len(args)-1], func(s runtime.Service) error {
		return s.ChangeImage(func(s string) string {
			return args[len(args)-1]
		})
	})
}
