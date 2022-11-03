// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/recipe"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "print all the configured services with user versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return errs.Wrap(err)
			}
			runtime, err := FromDir(pwd)
			if err != nil {
				return errs.Wrap(err)
			}
			st, err := recipe.GetStack()
			if err != nil {
				return errs.Wrap(err)
			}

			err = runtime.Reload(st)
			if err != nil {
				return errs.Wrap(err)
			}
			services := runtime.GetServices()
			for _, s := range services {
				fmt.Printf("%s\n", s.ID())
			}

			return nil
		},
	}
}

func init() {
	RootCmd.AddCommand(listCmd())
}
