// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"

	"storj.io/common/identity"
	up "storj.io/storj-up/pkg"
)

func init() {
	utilCmd := cobra.Command{
		Use:     "util",
		Aliases: []string{"utility"},
		Short:   "Small utilities to help container management",
	}

	{
		nodeIDCmd := &cobra.Command{
			Use:   "node-id",
			Short: "Generated node id string from identity file",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				nodeID, err := identity.NodeIDFromCertPath(args[0])
				if err != nil {
					return err
				}
				fmt.Println(nodeID)
				return nil
			},
		}
		utilCmd.AddCommand(nodeIDCmd)
	}
	{
		waitForSatelliteCmd := &cobra.Command{
			Use:   "wait-for-satellite",
			Short: "Wait until satellite can be called and return with the full NodeURL",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := context.Background()
				for {
					satellite, err := up.GetSatelliteID(ctx, args[0])
					if err != nil {
						println("Couldn't connect to satellite. Retrying... " + err.Error())
						time.Sleep(1 * time.Second)
						continue
					}

					fmt.Println(satellite)
					return nil
				}
			},
		}
		utilCmd.AddCommand(waitForSatelliteCmd)
	}
	{
		waitForPortCmd := &cobra.Command{
			Use:   "wait-for-port",
			Short: "Wait until ports is opened",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("Trying tcp connection to " + args[0])
				for {
					timeout := time.Second
					conn, err := net.DialTimeout("tcp", args[0], timeout)
					if err != nil {
						time.Sleep(300 * time.Millisecond)
						continue
					}
					_ = conn.Close()
					return nil
				}
			},
		}
		utilCmd.AddCommand(waitForPortCmd)
	}
	RootCmd.AddCommand(&utilCmd)
}
