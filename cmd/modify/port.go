// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"strconv"

	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

var overrideExternalPort int

func addPortCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <selector> port",
		Aliases: []string{"addport"},
		Short:   "Add port forward to a container",
		Long:    cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(addPort),
	}
	cmd.Flags().IntVarP(&overrideExternalPort, "external", "e", 0, "specify a different external port to forward than the internal port (default: internal port)")
	return cmd
}

func removePortCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <selector> port",
		Aliases: []string{"removeport", "rm", "delete"},
		Short:   "remove port forward from a container",
		Long:    "Remove port forward with the given target port from selected containers. " + cmd.SelectorHelp,
		Args:    cobra.MinimumNArgs(2),
		RunE:    cmd.ExecuteStorjUP(removePort),
	}
}

func init() {
	portCmd := cobra.Command{
		Use:   "port",
		Short: "add/remove port forward in specified services",
	}
	portCmd.AddCommand(addPortCmd())
	portCmd.AddCommand(removePortCmd())
	cmd.RootCmd.AddCommand(&portCmd)
}

func addPort(st recipe.Stack, rt runtime.Runtime, args []string) error {
	return runtime.ModifyService(st, rt, args[:len(args)-1], func(s runtime.Service) error {
		internal, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			return err
		}
		if overrideExternalPort == 0 {
			overrideExternalPort = internal
		}
		return s.AddPortForward(runtime.PortMap{
			Internal: internal,
			External: overrideExternalPort,
			Protocol: "tcp",
		})
	})
}

func removePort(st recipe.Stack, rt runtime.Runtime, args []string) error {
	return runtime.ModifyService(st, rt, args[:len(args)-1], func(s runtime.Service) error {
		port, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			return err
		}
		return s.RemovePortForward(runtime.PortMap{
			Internal: port,
		})
	})
}
