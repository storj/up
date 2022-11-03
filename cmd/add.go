// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/runtime/runtime"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "add <selector>",
		Short: "add more services to existing stack. " + SelectorHelp,
		Args:  cobra.MinimumNArgs(1),
		RunE:  ExecuteStorjUP(runtime.ApplyRecipes),
	})
}
