package cmd

import (
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build image on-the-fly instead of using pre-baked image",
	RunE: func(cmd *cobra.Command, args []string) error {
			return Build(repository, branch, args)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func Build(repository string, branch string, services []string) error {
	// magic here
	return nil
}