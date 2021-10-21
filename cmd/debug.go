package cmd

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug [service ...]",
	Short: "Turn on local debugging with DLV",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, SetEnv, "GO_DLV=true", args)
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}
var noDebugCmd = &cobra.Command{
	Use:   "no-debug [service ...]",
	Short: "Turn off local debugging with DLV",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, UnsetEnv, "GO_DLV", args)
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(noDebugCmd)
}