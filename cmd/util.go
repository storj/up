// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"storj.io/common/identity"
	up "storj.io/storj-up/pkg"
)

// WaitForTCPConnection waits until a TCP connection can be established to the given address.
func WaitForTCPConnection(address string, interval time.Duration) error {
	fmt.Printf("Trying TCP connection to %s\n", address)
	for {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			time.Sleep(interval)
			continue
		}
		_ = conn.Close()
		return nil
	}
}

// WaitForPostgresOrCockroach waits for Postgres or Cockroach DB to be ready.
func WaitForPostgresOrCockroach(dbURL string) error {
	// For CockroachDB/Postgres format: cockroach://root@cockroach:26257/master?sslmode=disable
	// Extract host and port using regex
	parts := strings.Split(dbURL, "@")
	if len(parts) < 2 {
		return fmt.Errorf("invalid connection string format")
	}

	hostPortPart := strings.Split(parts[1], "/")[0]
	hostPort := hostPortPart
	if strings.Contains(hostPort, "?") {
		hostPort = strings.Split(hostPort, "?")[0]
	}

	fmt.Printf("Waiting for database connection at %s\n", hostPort)

	// Wait for TCP connection
	if err := WaitForTCPConnection(hostPort, 300*time.Millisecond); err != nil {
		return err
	}

	fmt.Printf("Database is ready: %s\n", dbURL)
	return nil
}

func init() {
	utilCmd := cobra.Command{
		Use:     "util",
		Aliases: []string{"utility"},
		Short:   "Small utilities to help container management",
	}

	{
		nodeIDCmd := &cobra.Command{
			Use:   "node-id <identity_file>",
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
			Use:   "wait-for-satellite <satellite_id>",
			Short: "Wait until satellite can be called and return with the full NodeURL",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx := context.Background()
				for {
					satellite, err := up.GetSatelliteID(ctx, args[0])
					if err != nil {
						fmt.Printf("Couldn't connect to satellite. Retrying... %s\n", err.Error())
						time.Sleep(time.Second)
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
			Use:   "wait-for-port <port>",
			Short: "Wait until port is opened",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return WaitForTCPConnection(args[0], 300*time.Millisecond)
			},
		}
		utilCmd.AddCommand(waitForPortCmd)
	}

	{
		waitForDBCmd := &cobra.Command{
			Use:   "wait-for-db <connection_string>",
			Short: "Wait until the database is ready to accept connections",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				dbURL := args[0]

				if strings.HasPrefix(dbURL, "cockroach://") ||
					strings.HasPrefix(dbURL, "postgres://") ||
					strings.HasPrefix(dbURL, "tidb://") {
					return WaitForPostgresOrCockroach(dbURL)
				}
				fmt.Println("Unsupported database type: " + dbURL)
				return fmt.Errorf("unsupported database type")
			},
		}
		utilCmd.AddCommand(waitForDBCmd)
	}

	RootCmd.AddCommand(&utilCmd)
}
