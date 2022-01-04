// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/common"
)

func init() {
	argCmd := cobra.Command{
		Use:   "args",
		Short: "set/unset build arguments on specified services",
	}
	rootCmd.AddCommand(&argCmd)
	argCmd.AddCommand(setArgCmd())
	argCmd.AddCommand(unsetArgCmd())
}

func setArgCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <selector> KEY=VALUE",
		Short: "Set build arguments on service. Build arguments should be supported by referenced Dockerfile " + selectorHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
			if err != nil {
				return err
			}
			selector, args, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, setArg, args[0], selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func unsetArgCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <selector> KEY",
		Short: "remove container arg",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
			if err != nil {
				return err
			}
			selector, args, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, unsetArg, args[0], selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
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
