package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image <image> [service ...]",
	Short: "use a prebuilt docker image",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, setImage, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
}

func setImage(composeService *types.ServiceConfig, image string) error {
	composeService.Image = image
	return nil
}