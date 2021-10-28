package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "build local src directories for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, buildLocalSrc, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	buildCmd.AddCommand(localCmd)
}

func buildLocalSrc(_ *types.ServiceConfig, _ string) error {
	// do magic here
	return nil
}