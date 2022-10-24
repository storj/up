// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package container

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func enableDebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <selector>",
		Short: "turn on local debugging with DLV",
		Long:  "Add environment variable which will activate the DLV debug. Container won't start until the agent is connected. " + cmd.SelectorHelp,
		RunE:  cmd.ExecuteStorjUP(enableDebug),
	}
}

func disableDebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable [service ...]",
		Short: "turn off local debugging with DLV",
		RunE:  cmd.ExecuteStorjUP(disableDebug),
	}
}

func init() {
	debugCmd := cobra.Command{
		Use:   "debug",
		Short: "enable/disable local DLV based go debug",
	}

	debugCmd.AddCommand(enableDebugCmd())
	debugCmd.AddCommand(disableDebugCmd())
	cmd.RootCmd.AddCommand(&debugCmd)
}

func enableDebug(st recipe.Stack, rt runtime.Runtime, args []string) error {
	return runtime.ModifyService(st, rt, args, func(s runtime.Service) error {
		composeService, ok := s.(*compose.Service)
		if !ok {
			return errs.Errorf("this subcommand is supported only for compose based environments")
		}
		return composeService.TransformRaw(func(composeService *types.ServiceConfig) error {
			tr := "true"
			composeService.Environment["GO_DLV"] = &tr
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
		})
	})
}

func disableDebug(st recipe.Stack, rt runtime.Runtime, args []string) error {
	return runtime.ModifyService(st, rt, args[:len(args)-1], func(s runtime.Service) error {
		composeService, ok := s.(*compose.Service)
		if !ok {
			return errs.Errorf("this subcommand is supported only for compose based environments")
		}
		return composeService.TransformRaw(func(composeService *types.ServiceConfig) error {
			delete(composeService.Environment, "GO_DLV")
			for i, port := range composeService.Ports {
				if port.Target == 2345 {
					composeService.Ports = append(composeService.Ports[:i], composeService.Ports[i+1:]...)
				}
			}
			return nil
		})
	})
}
