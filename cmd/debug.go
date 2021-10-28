package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

func DebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "debug [service ...]",
		Short: "turn on local debugging with DLV",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, SetDebug, "GO_DLV=true", args)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func NoDebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "no-debug [service ...]",
		Short: "turn off local debugging with DLV",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, UnsetDebug, "GO_DLV", args)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(DebugCmd())
	rootCmd.AddCommand(NoDebugCmd())
}

func SetDebug(composeService *types.ServiceConfig, arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	composeService.Environment[parts[0]] = &parts[1]
	composeService.Ports = append(composeService.Ports, types.ServicePortConfig{
		Mode:      "ingress",
		Target:    2345,
		Published: 2345,
		Protocol:  "tcp",
	})
	return nil
}

func UnsetDebug(composeService *types.ServiceConfig, arg string) error {
	delete(composeService.Environment, arg)
	for i, port := range composeService.Ports {
		if port.Target == 2345 {
			composeService.Ports = append(composeService.Ports[:i], composeService.Ports[i+1:]...)
		}
	}
	return nil
}
