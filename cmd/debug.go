// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/common"
)

func enableDebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <selector> ",
		Short: "turn on local debugging with DLV",
		Long:  "Add environment variable which will activate the DLV debug. Container won't start until the agent is connected. " + selectorHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(composeFile)
			if err != nil {
				return err
			}

			selector, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, setDebug, "GO_DLV=true", selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func disableDebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable [service ...]",
		Short: "turn off local debugging with DLV",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(composeFile)
			if err != nil {
				return err
			}

			selector, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}
			updatedComposeProject, err := common.UpdateEach(composeProject, unsetDebug, "GO_DLV", selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	debugCmd := cobra.Command{
		Use:   "debug",
		Short: "enable/disable local DLV based go debug",
	}

	debugCmd.AddCommand(enableDebugCmd())
	debugCmd.AddCommand(disableDebugCmd())
	rootCmd.AddCommand(&debugCmd)
}

// setDebug enables the port-forward and environment variable for one docker service.
func setDebug(composeService *types.ServiceConfig, arg string) error {
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

func unsetDebug(composeService *types.ServiceConfig, arg string) error {
	delete(composeService.Environment, arg)
	for i, port := range composeService.Ports {
		if port.Target == 2345 {
			composeService.Ports = append(composeService.Ports[:i], composeService.Ports[i+1:]...)
		}
	}
	return nil
}
