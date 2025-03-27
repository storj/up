// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	spannerdb "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"

	"storj.io/common/identity"
	up "storj.io/storj-up/pkg"
)

// SpannerDBInfo holds parsed information about a Spanner spannerdb URL.
type SpannerDBInfo struct {
	URL          string
	ProjectID    string
	InstanceName string
	DatabaseName string
	ProjectPath  string
	InstancePath string
	DatabasePath string
}

// ParseSpannerURL extracts components from a Spanner URL.
func ParseSpannerURL(dbURL string) (*SpannerDBInfo, error) {
	// Format: spanner://projects/test-project/instances/test-instance/databases/master
	parts := strings.Split(strings.TrimPrefix(dbURL, "spanner://"), "/")
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid spanner URL format: expected format spanner://projects/<project>/instances/<instance>/databases/<spannerdb>")
	}

	projectID := parts[1]
	instanceName := parts[3]
	databaseName := parts[5]

	projectPath := fmt.Sprintf("projects/%s", projectID)
	instancePath := fmt.Sprintf("%s/instances/%s", projectPath, instanceName)
	databasePath := fmt.Sprintf("%s/databases/%s", instancePath, databaseName)

	return &SpannerDBInfo{
		URL:          dbURL,
		ProjectID:    projectID,
		InstanceName: instanceName,
		DatabaseName: databaseName,
		ProjectPath:  projectPath,
		InstancePath: instancePath,
		DatabasePath: databasePath,
	}, nil
}

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

// WaitForSpannerInstance waits for a Spanner instance to be ready.
func WaitForSpannerInstance(ctx context.Context, dbInfo *SpannerDBInfo, maxAttempts int) error {
	fmt.Printf("Waiting for Spanner instance: %s\n", dbInfo.InstancePath)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Checking for Spanner instance (attempt %d/%d)...\n", attempt, maxAttempts)

		instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
		if err != nil {
			fmt.Printf("Error creating instance admin client (will retry): %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		defer func(instanceAdmin *instance.InstanceAdminClient) {
			err := instanceAdmin.Close()
			if err != nil {
				fmt.Printf("Error closing instance admin client: %v\n", err)
			}
		}(instanceAdmin)

		// Build the request to list instances
		req := &instancepb.ListInstancesRequest{
			Parent: dbInfo.ProjectPath,
		}

		// List instances
		it := instanceAdmin.ListInstances(ctx, req)
		instanceFound := false

		for {
			resp, err := it.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				fmt.Printf("Error listing instances (will retry): %v\n", err)
				break
			}

			fmt.Printf("Found instance: %s\n", resp.GetName())
			if resp.GetName() == dbInfo.InstancePath {
				instanceFound = true
				break
			}
		}

		if instanceFound {
			fmt.Printf("Spanner instance '%s' is ready!\n", dbInfo.InstancePath)
			return nil
		}

		fmt.Printf("Instance '%s' not found or not ready, waiting...\n", dbInfo.InstancePath)
		time.Sleep(time.Second)
	}

	return fmt.Errorf("timed out waiting for Spanner instance '%s' after %d attempts", dbInfo.InstancePath, maxAttempts)
}

// WaitForSpannerDatabase waits for a Spanner spannerdb to be ready.
func WaitForSpannerDatabase(ctx context.Context, dbInfo *SpannerDBInfo, maxAttempts int) error {
	fmt.Printf("Waiting for Spanner spannerdb: %s\n", dbInfo.DatabasePath)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Checking for spannerdb (attempt %d/%d)...\n", attempt, maxAttempts)

		dbAdmin, err := spannerdb.NewDatabaseAdminClient(ctx)
		if err != nil {
			fmt.Printf("Error creating spannerdb admin client (will retry): %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		defer func(dbAdmin *spannerdb.DatabaseAdminClient) {
			err := dbAdmin.Close()
			if err != nil {
				fmt.Printf("Error closing spannerdb admin client: %v\n", err)
			}
		}(dbAdmin)

		// List databases in the instance
		dbReq := &databasepb.ListDatabasesRequest{
			Parent: dbInfo.InstancePath,
		}

		dbIt := dbAdmin.ListDatabases(ctx, dbReq)
		dbFound := false

		for {
			db, err := dbIt.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				fmt.Printf("Error listing databases (will retry): %v\n", err)
				break
			}

			fmt.Printf("Found spannerdb: %s\n", db.GetName())
			if db.GetName() == dbInfo.DatabasePath {
				dbFound = true
				break
			}
		}

		if dbFound {
			fmt.Printf("Spanner spannerdb '%s' is ready!\n", dbInfo.DatabasePath)
			return nil
		}

		fmt.Printf("Database '%s' not found or not ready, waiting...\n", dbInfo.DatabasePath)
		time.Sleep(time.Second)
	}

	return fmt.Errorf("timed out waiting for Spanner spannerdb '%s' after %d attempts", dbInfo.DatabasePath, maxAttempts)
}

// WaitForSpannerDB waits for both the Spanner instance and spannerdb to be ready.
func WaitForSpannerDB(dbURL string) error {
	ctx := context.Background()

	// First ensure the spanner service is up on port 9010
	fmt.Println("Checking if Spanner service is up...")
	if err := WaitForTCPConnection("spanner:9010", time.Second); err != nil {
		return fmt.Errorf("failed to connect to spanner service: %w", err)
	}

	// Parse the Spanner URL
	dbInfo, err := ParseSpannerURL(dbURL)
	if err != nil {
		return err
	}

	// Wait for instance to be ready
	if err := WaitForSpannerInstance(ctx, dbInfo, 60); err != nil {
		return err
	}

	// Wait for spannerdb to be ready
	return WaitForSpannerDatabase(ctx, dbInfo, 30)
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
			Short: "Wait until spannerdb is ready to accept connections",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				dbURL := args[0]

				// Handle different spannerdb connection string types
				if strings.HasPrefix(dbURL, "spanner://") {
					fmt.Printf("Detected Spanner spannerdb URL: %s\n", dbURL)
					return WaitForSpannerDB(dbURL)
				} else if strings.HasPrefix(dbURL, "cockroach://") ||
					strings.HasPrefix(dbURL, "postgres://") {
					return WaitForPostgresOrCockroach(dbURL)
				} else {
					// For other spannerdb types
					fmt.Println("Unsupported spannerdb type: " + dbURL)
					return fmt.Errorf("unsupported spannerdb type")
				}
			},
		}
		utilCmd.AddCommand(waitForDBCmd)
	}

	RootCmd.AddCommand(&utilCmd)
}
