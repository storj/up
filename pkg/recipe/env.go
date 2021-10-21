package recipe

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)


func envCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "env [KEY=VALUE]",
		Short: "Set environment variable / parameter in a container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.UpdateEach(service, func(service *common.ServiceConfig) error {
				return SetEnv(service, args[0])
			})
		},
	})
}

func SetEnv(service *common.ServiceConfig, arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	service.Environment[parts[0]] = &parts[1]
	return nil
}
