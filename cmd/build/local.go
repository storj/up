// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package build

import (
	"github.com/spf13/cobra"
)

var path string

const (
	local = "local"
)

func localCmd() *cobra.Command {
	// NOTE cobra doesn't have a way to document positional parameters:
	// https://github.com/spf13/cobra/issues/378
	localCmd := &cobra.Command{
		Use:   "local <selector>",
		Short: "build local src repo for use inside the container",
		Long: `build local src repo for use inside the container for the indicated
services through positional arguments. See the list of supported service running
` + "`storj-up services`.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := updateCompose(args, local)
			if err != nil {
				return err
			}
			return nil
		},
	}
	localCmd.Flags().StringVarP(&path, "path", "p", "", "The path to the local repo")
	return localCmd
}

func init() {
	buildCmd.AddCommand(localCmd())
}
