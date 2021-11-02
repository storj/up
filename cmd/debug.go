package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

func DebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "debug <selector> ",
		Short: "turn on local debugging with DLV",
		Long:  "Add environment variable which will activate the DLV debug. Container won't start until the agent is connected. " + selectorHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			if err != nil {
				return err
			}

			selector, _, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, SetDebug, "GO_DLV=true", selector)
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
			if err != nil {
				return err
			}

			selector, _, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}
			updatedComposeProject, err := common.UpdateEach(composeProject, UnsetDebug, "GO_DLV", selector)
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
	for _, portConfig := range composeService.Ports {
		if portConfig.Mode == "ingress" &&
			portConfig.Target == 2345 &&
			portConfig.Published == 2345 &&
			portConfig.Protocol == "tcp" {
			return nil
		}
	}
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
