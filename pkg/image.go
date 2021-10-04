package sjr

import "github.com/spf13/cobra"

func imageCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "image",
		Short: "Set image to custom variable",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return UpdateEach(service, func(service *ServiceConfig) error {
				return SetImage(service, args[0])
			})
		},
	})
}

func SetImage(service *ServiceConfig, image string) error {
	service.Image = image
	return nil
}
