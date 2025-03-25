// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package build

import (
	"github.com/spf13/cobra"
)

var branch, commit, ref string

const (
	github = "github"
	gerrit = "gerrit"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Args:  cobra.NoArgs,
	Short: "build from a remote src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func githubCmd() *cobra.Command {
	// NOTE cobra doesn't have a way to document positional parameters:
	// https://github.com/spf13/cobra/issues/378
	githubCmd := &cobra.Command{
		Use:   "github <selector>...",
		Short: "build github src repo for use inside the container",
		Long: `build github src repo for use inside the container for the indicated
services through positional arguments. See the list of supported service running
` + "`storj-up services`.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, selector []string) error {
			err := updateCompose(selector, github)
			if err != nil {
				return err
			}
			return nil
		},
	}
	githubCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "main", "The branch to checkout and build")
	githubCmd.PersistentFlags().StringVarP(&commit, "commit", "c", "", "The commit to checkout and build")
	return githubCmd
}

func gerritCmd() *cobra.Command {
	// NOTE cobra doesn't have a way to document positional parameters:
	// https://github.com/spf13/cobra/issues/378
	gerritCmd := &cobra.Command{
		Use:   "gerrit <selector>...",
		Short: "build gerrit src repo for use inside the container",
		Long: `build gerrit src repo for use inside the container for the indicated
services through positional arguments. See the list of supported service running
` + "`storj-up services`.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, selector []string) error {
			err := updateCompose(selector, gerrit)
			if err != nil {
				return err
			}
			return nil
		},
	}
	gerritCmd.PersistentFlags().StringVarP(&ref, "refspec", "f", "", "The gerrit refspec to checkout and build")
	_ = gerritCmd.MarkPersistentFlagRequired("refspec")
	return gerritCmd
}

func init() {
	buildCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(githubCmd())
	remoteCmd.AddCommand(gerritCmd())
}
