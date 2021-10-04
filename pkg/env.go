package sjr

import (
	"github.com/spf13/cobra"
	"strings"
)

func envCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "env [KEY=VALUE]",
		Short: "Set environment variable / parameter in a container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return UpdateEach(service, func(service *ServiceConfig) error {
				return SetEnv(service, args[0])
			})
		},
	})
}

func SetEnv(service *ServiceConfig, arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	service.Environment[parts[0]] = &parts[1]
	return nil
}
