package sjr

import (
	"github.com/spf13/cobra"
)

func localEntrypointCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "local-entrypoint",
		Short: "Bind mount entrypoint.sh to use local modifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			return UpdateEach(service, localEntrypoint)
		},
	})
}

func localEntrypoint(service *ServiceConfig) error {
	service.Volumes = append(service.Volumes,
		"./entrypoint.sh:/var/lib/storj/entrypoint.sh")
	return nil
}
