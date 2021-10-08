package recipe

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

func imageCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "image",
		Short: "Set image to custom variable",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.UpdateEach(service, func(service *common.ServiceConfig) error {
				return SetImage(service, args[0])
			})
		},
	})
}

func SetImage(service *common.ServiceConfig, image string) error {
	service.Image = image
	return nil
}
