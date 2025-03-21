// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/recipe"
)

func serviceCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "recipes",
		Aliases: []string{"recipe", "services"},
		Short:   "Return available recipes and included service names",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println("The following predefined recipes are available:")
			fmt.Println()
			stack, err := recipe.GetStack()
			if err != nil {
				return err
			}
			for _, r := range stack {
				var services []string
				for _, s := range r.Add {
					services = append(services, s.Name)
				}
				fmt.Printf("%-10s %s (%s)\n", r.Name, r.Description, strings.Join(services, ","))
			}
			fmt.Println("")
			fmt.Println("You can use both the service names (like satellite-api) or recipe names (like minimal) when you init/add clusters")
			return nil
		},
	}
}

func init() {
	RootCmd.AddCommand(serviceCmd())
}
