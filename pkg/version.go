package sjr

import (
	"github.com/spf13/cobra"
	"strings"
)

func versionCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "version [version]",
		Short: "Set version (docker image tag) for specified services",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return UpdateEach(service, func(service *ServiceConfig) error {
				return SetVersion(service, args[0])
			})
		},
	})
}

func SetVersion(service *ServiceConfig, version string) error {
	service.Image = strings.ReplaceAll(service.Image, "@sha256", "")
	service.Image = strings.Split(service.Image, ":")[0] + ":" + version
	return nil
}
