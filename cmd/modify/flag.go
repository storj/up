// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func init() {
	c := cobra.Command{
		Use:     "flags",
		Aliases: []string{"flag"},
		Short:   "set/unset flags on the startup command",
	}
	cmd.RootCmd.AddCommand(&c)
	c.AddCommand(setFlagCmd())
	c.AddCommand(removeFlagCmd())
}

func setFlagCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set <selector>.. <KEY>=<VALUE>",
		Aliases: []string{"add"},
		Short:   "Set (or add) command line flags on the startup command on the container(s). " + cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(addFlag),
	}
}

func removeFlagCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <selector>.. <KEY>",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove command line flags of the startup command on the container(s). " + cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(removeFlag),
	}
}

func addFlag(st recipe.Stack, rt runtime.Runtime, args []string) error {
	selector, keyvalue := common.SplitArgsSelector1(args)
	return runtime.ModifyService(st, rt, selector, func(s runtime.Service) error {
		return s.AddFlag("--" + keyvalue)
	})
}

func removeFlag(st recipe.Stack, rt runtime.Runtime, args []string) error {
	selector, key := common.SplitArgsSelector1(args)
	return runtime.ModifyService(st, rt, selector, func(s runtime.Service) error {
		return s.RemoveFlag("--" + key)
	})
}
