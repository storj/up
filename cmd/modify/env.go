// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func setEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set <selector>... <KEY>=<VALUE>",
		Aliases: []string{"setenv"},
		Short:   "Set environment variable / parameter in a container",
		Long:    cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(setEnv),
	}
}

func unsetEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "unset <selector>... <KEY>",
		Aliases: []string{"unsetenv", "rm", "delete"},
		Short:   "remove environment variable / parameter in a container",
		Long:    "Remove environment variable from selected containers. " + cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(removeEnv),
	}
}

func init() {
	envCmd := cobra.Command{
		Use:   "env",
		Short: "add/remove environment variables to specified services",
	}
	envCmd.AddCommand(setEnvCmd())
	envCmd.AddCommand(unsetEnvCmd())
	cmd.RootCmd.AddCommand(&envCmd)
}

func setEnv(st recipe.Stack, rt runtime.Runtime, args []string) error {
	selector, keyvalue := common.SplitArgsSelector1(args)
	return runtime.ModifyService(st, rt, selector, func(s runtime.Service) error {
		key, value, ok := strings.Cut(keyvalue, "=")
		if !ok {
			return errors.New("expected key=value")
		}
		return s.AddEnvironment(key, value)
	})
}

func removeEnv(st recipe.Stack, rt runtime.Runtime, args []string) error {
	selector, key := common.SplitArgsSelector1(args)
	return runtime.ModifyService(st, rt, selector, func(s runtime.Service) error {
		return s.AddEnvironment(key, "")
	})
}
