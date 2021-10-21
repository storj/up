package cmd

import (
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "build local src directories for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildLocalSrc()
	},
}

func init() {
	buildCmd.AddCommand(localCmd)
}

func buildLocalSrc() error {
	// do magic here
	return nil
}