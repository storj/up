package cmd

import (
	"github.com/spf13/cobra"
)

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "mount local binaries to the default docker image",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mountBinaries()
	},
}

func init() {
	rootCmd.AddCommand(mountCmd)
}

func mountBinaries() error {
	//magic here
	return nil
}