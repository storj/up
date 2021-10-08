package recipe

import (
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

func debugCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "debug",
		Short: "Turn on local debugging with DLV",
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.UpdateEach(service, DebugEnable)
		},
	})
	command.AddCommand(&cobra.Command{
		Use:   "no-debug",
		Short: "Turn off local debugging with DLV",
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.UpdateEach(service, DebugDisable)
		},
	})
}

func DebugEnable(service *common.ServiceConfig) error {
	tr := "true"
	service.Environment["GO_DLV"] = &tr
	return nil
}

func DebugDisable(service *common.ServiceConfig) error {
	delete(service.Environment, "GO_DLV")
	return nil
}
