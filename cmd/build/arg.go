// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package build

import (
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func init() {
	argCmd := cobra.Command{
		Use:   "args",
		Args:  cobra.NoArgs,
		Short: "set/unset build arguments on specified services",
	}
	cmd.RootCmd.AddCommand(&argCmd)
	argCmd.AddCommand(setArgCmd())
	argCmd.AddCommand(unsetArgCmd())
}

func setArgCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <selector>... <KEY>=<VALUE>",
		Short: "Set build arguments on service. Build arguments should be supported by referenced Dockerfile " + cmd.SelectorHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE: cmd.ExecuteStorjUP(func(st recipe.Stack, rt runtime.Runtime, args []string) error {
			selector, keyvalue := common.SplitArgsSelector1(args)
			return cmd.ChangeCompose(st, rt, selector, func(composeService *types.ServiceConfig) error {
				return setArg(composeService, keyvalue)
			})
		}),
	}
}

func unsetArgCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <selector>... <KEY>",
		Short: "remove container arg",
		Args:  cobra.MinimumNArgs(2),
		RunE: cmd.ExecuteStorjUP(func(st recipe.Stack, rt runtime.Runtime, args []string) error {
			selector, key := common.SplitArgsSelector1(args)
			return cmd.ChangeCompose(st, rt, selector, func(composeService *types.ServiceConfig) error {
				return unsetArg(composeService, key)
			})
		}),
	}
}

func setArg(composeService *types.ServiceConfig, arg string) error {
	if composeService.Build == nil {
		composeService.Build = &types.BuildConfig{
			Args: map[string]*string{},
		}
	} else if composeService.Build.Args == nil {
		composeService.Build.Args = map[string]*string{}
	}
	parts := strings.SplitN(arg, "=", 2)
	composeService.Build.Args[parts[0]] = &parts[1]
	return nil
}

func unsetArg(composeService *types.ServiceConfig, arg string) error {
	delete(composeService.Build.Args, arg)
	return nil
}
