package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

func ImageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "image <image> [service ...]",
		Short: "Use a prebuilt docker image",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, setImage, args[0], args[1:])
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(ImageCmd())
}

func setImage(composeService *types.ServiceConfig, image string) error {
	composeService.Image = image
	return nil
}
