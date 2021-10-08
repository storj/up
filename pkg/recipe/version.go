package recipe

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

func versionCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "version [version]",
		Short: "Set version (docker image tag) for specified services",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.UpdateEach(service, func(service *common.ServiceConfig) error {
				return SetVersion(service, args[0])
			})
		},
	})
}

func SetVersion(service *common.ServiceConfig, version string) error {
	service.Image = strings.ReplaceAll(service.Image, "@sha256", "")
	service.Image = strings.Split(service.Image, ":")[0] + ":" + version
	return nil
}
