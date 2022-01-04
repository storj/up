// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/common"
)

func init() {
	cmd := cobra.Command{
		Use:     "flags",
		Aliases: []string{"flag"},
		Short:   "set/unset flags on the startup command",
	}
	rootCmd.AddCommand(&cmd)
	cmd.AddCommand(setFlagCmd())
	cmd.AddCommand(removeFlagCmd())
}

func setFlagCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set <selector> KEY=VALUE",
		Aliases: []string{"add"},
		Short:   "Set (or add) command line flags on the startup command on the container(s). " + selectorHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
			if err != nil {
				return err
			}
			selector, args, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, func(config *types.ServiceConfig, s string) error {
				config.Command = setFlag(config.Command, s)
				return nil
			}, args[0], selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func removeFlagCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <selector> KEY",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove command line flags of the startup command on the container(s). " + selectorHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
			if err != nil {
				return err
			}
			selector, args, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, func(config *types.ServiceConfig, s string) error {
				config.Command = removeFlag(config.Command, s)
				return nil
			}, args[0], selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func setFlag(command []string, s string) []string {
	var res []string
	exist := false

	// s is in format foobar=value
	parts := strings.SplitN(s, "=", 2)

	for _, c := range command {
		if strings.HasPrefix(c, "-"+parts[0]+"=") || strings.HasPrefix(c, "--"+parts[0]+"=") {
			f := s
			if f[0] != '-' {
				f = "-" + f
			}

			// existing flag with two --
			if len(c) > 2 && c[1] == '-' {
				f = "-" + f
			}
			res = append(res, f)
			exist = true
		} else {
			res = append(res, c)
		}
	}
	if !exist {
		res = append(res, "--"+s)
	}
	return res
}

func removeFlag(command []string, s string) []string {
	var res []string

	// s is in format foobar=value
	parts := strings.SplitN(s, "=", 2)

	for _, c := range command {
		if !strings.HasPrefix(c, "-"+parts[0]+"=") && !strings.HasPrefix(c, "--"+parts[0]+"=") {
			res = append(res, c)
		}
	}
	return res

}
