package cmd

import (
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build image on-the-fly instead of using pre-baked image",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(BuildCmd)
}
