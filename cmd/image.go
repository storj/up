package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"storj.io/storj-up/pkg/common"
)

const selectorHelp = "<selector> is a service name or group (use `storj-up service` to list available services)"

func ImageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "image <selector> <image>",
		Short: "Change container image for one more more services",
		Long:  "Change image of specified services." + selectorHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			if err != nil {
				return err
			}

			selector, arguments, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, SetImage, arguments[0], selector)
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

func SetImage(composeService *types.ServiceConfig, image string) error {
	composeService.Image = image
	return nil
}
