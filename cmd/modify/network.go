// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func setNetworkCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set <selector>... <NETWORK>",
		Aliases: []string{"setnetwork"},
		Short:   "Set network for a service or services to use",
		Long:    cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(setNetwork),
	}
}

func unsetNetworkCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "unset <selector>... <NETWORK>",
		Aliases: []string{"unsetnetwork", "rm", "delete"},
		Short:   "remove network for a specific service or services",
		Long:    cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(removeNetwork),
	}
}

func init() {
	networkCmd := cobra.Command{
		Use:   "network",
		Short: "set/unset network for specified services",
	}
	networkCmd.AddCommand(setNetworkCmd())
	networkCmd.AddCommand(unsetNetworkCmd())
	cmd.RootCmd.AddCommand(&networkCmd)
}

func setNetwork(st recipe.Stack, rt runtime.Runtime, args []string) error {
	selector, network := common.SplitArgsSelector1(args)
	return runtime.ModifyService(st, rt, selector, func(s runtime.Service) error {
		if t, ok := s.(runtime.ManageableNetwork); ok {
			return t.AddNetwork(network)
		}
		return errs.New("runtime does not support network management")
	})
}

func removeNetwork(st recipe.Stack, rt runtime.Runtime, args []string) error {
	selector, network := common.SplitArgsSelector1(args)
	return runtime.ModifyService(st, rt, selector, func(s runtime.Service) error {
		if t, ok := s.(runtime.ManageableNetwork); ok {
			return t.RemoveNetwork(network)
		}
		return errs.New("runtime does not support network management")
	})
}
