// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package build

import (
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
)

var skipFrontend bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build image on-the-fly instead of using pre-baked image",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	buildCmd.PersistentFlags().BoolVarP(&skipFrontend, "skipfrontend", "s", false, "Skip building the frontend")
	cmd.RootCmd.AddCommand(buildCmd)
}
