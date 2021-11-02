package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"storj.io/storj-up/pkg/common"
	"github.com/spf13/cobra"
)

func LocalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "local",
		Short: "build local src directories for use inside the container",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, buildLocalSrc, args[0], args[1:])
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	BuildCmd.AddCommand(LocalCmd())
}

func buildLocalSrc(_ *types.ServiceConfig, _ string) error {
	// do magic here
	return nil
}
