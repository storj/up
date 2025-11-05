// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
)

var table, host, user, dbname, dbtype string
var number, timeout, port int
var spannerURL string

func healthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Args:  cobra.NoArgs,
		Short: "wait until cluster is healthy (10 storagenodes are registered in the db)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return checkHealthWithTimeout(dbtype, table, number, timeout)
		},
	}
	cmd.PersistentFlags().StringVarP(&table, "table", "t", "nodes", "table to use for health check")
	cmd.PersistentFlags().IntVarP(&number, "number", "n", 10, "number of entries to expect in the table")
	cmd.PersistentFlags().IntVarP(&timeout, "duration", "d", 0, "time to wait (in seconds) for health check")
	cmd.PersistentFlags().StringVarP(&host, "host", "", "localhost", "host/ip address to use for health check. Defaults to localhost or STORJ_DOCKER_HOST if set.")
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 9010, "port to use for health check")
	cmd.Flags().StringVarP(&user, "user", "u", "root", "user to connect to the DB as")
	cmd.Flags().StringVarP(&dbname, "dbname", "", "master", "DB name to connect to")
	cmd.Flags().StringVarP(&dbtype, "dbtype", "", "spanner", "database type (spanner, postgres, or cockroach)")
	cmd.Flags().StringVarP(&spannerURL, "spanner-url", "", "", "URL for Spanner connection in format spanner://projects/PROJECT/instances/INSTANCE/databases/DATABASE")
	return cmd
}

func init() {
	RootCmd.AddCommand(healthCmd())
}

// checkHealthWithTimeout polls the database until all storagenodes are checked in, or the timeout is exceeded.
// A timeout of 0 (default) means no timeout.
func checkHealthWithTimeout(dbtype, table string, records, timeout int) error {
	if timeout == 0 {
		return checkHealth(dbtype, table, records)
	}

	c1 := make(chan error, 1)
	go func() {
		c1 <- checkHealth(dbtype, table, records)
	}()

	select {
	case err := <-c1:
		return err
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("health check failed: duration limit reached")
	}
}

// checkHealth polls the database until all storagenodes are checked in.
func checkHealth(dbtype, table string, records int) error {
	switch strings.ToLower(dbtype) {
	case "cockroach", "postgres":
		return checkHealthCockroachPostgres(table, records)
	case "spanner":
		return checkHealthSpanner(table, records)
	default:
		return fmt.Errorf("unsupported database type: %s", dbtype)
	}
}

// checkHealthCockroachPostgres polls a Cockroach or Postgres database until the required records are present.
func checkHealthCockroachPostgres(table string, records int) error {
	prevCount := -1
	defaultHost := os.Getenv("STORJ_DOCKER_HOST")
	if host == "127.0.0.1" && defaultHost != "" {
		host = defaultHost
	}

	for {
		time.Sleep(1 * time.Second)
		db, err := pgx.Connect(context.TODO(), "host="+host+" port="+strconv.Itoa(port)+" user="+user+" dbname="+dbname+" sslmode=disable")
		if err != nil {
			fmt.Printf("Couldn't connect to the database: %s\n", err.Error())
			continue
		}

		row := db.QueryRow(context.TODO(), "select count(*) from "+table)
		var count int
		err = row.Scan(&count)
		_ = db.Close(context.TODO())

		if err != nil {
			fmt.Printf("Couldn't query database for records: %s\n", err.Error())
			continue
		}

		if count == records {
			fmt.Println()
			fmt.Println(table, "has", records, "records")
			return nil
		}

		if count != prevCount {
			fmt.Printf("Found only %d records in the database ", count)
		} else {
			fmt.Print(".")
		}
		prevCount = count
	}
}

// checkHealthSpanner polls a Spanner database until the required records are present.
func checkHealthSpanner(table string, records int) error {
	prevCount := -1

	// If no spanner URL was provided, try to construct one
	if spannerURL == "" {
		// Default project and instance if not provided
		project := os.Getenv("SPANNER_PROJECT_ID")
		if project == "" {
			project = "test-project"
		}

		instance := os.Getenv("SPANNER_INSTANCE")
		if instance == "" {
			instance = "test-instance"
		}

		spannerURL = fmt.Sprintf("spanner://projects/%s/instances/%s/databases/%s", project, instance, dbname)
	}

	// Parse the Spanner URL
	dbInfo, err := ParseSpannerURL(spannerURL)
	if err != nil {
		return fmt.Errorf("invalid Spanner URL: %w", err)
	}

	// Configure Spanner emulator connection
	emulatorHost := host
	if host == "127.0.0.1" {
		emulatorHost = "localhost"
	}
	emulatorAddr := fmt.Sprintf("%s:9010", emulatorHost)
	err = os.Setenv("SPANNER_EMULATOR_HOST", emulatorAddr)
	if err != nil {
		return fmt.Errorf("failed to set SPANNER_EMULATOR_HOST: %w", err)
	}
	fmt.Printf("Connecting to Spanner emulator at %s\n", emulatorAddr)

	ctx := context.Background()

	for {
		time.Sleep(3 * time.Second)

		clientCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		client, err := spanner.NewClient(clientCtx, dbInfo.DatabasePath)

		if err != nil {
			cancel()
			fmt.Printf("Connection to Spanner failed: %s\n", err.Error())
			continue
		}

		// Check if the table exists
		iter := client.Single().Query(clientCtx, spanner.Statement{
			SQL: `SELECT table_name FROM information_schema.tables WHERE table_name = @table_name`,
			Params: map[string]interface{}{
				"table_name": table,
			},
		})

		tableExists := false
		_, err = iter.Next()
		if err == nil {
			tableExists = true
		} else if !errors.Is(err, iterator.Done) {
			fmt.Printf("Error checking table existence: %s\n", err.Error())
		}
		iter.Stop()

		if !tableExists {
			fmt.Printf("Table %s doesn't exist yet, waiting...\n", table)
			client.Close()
			cancel()
			continue
		}

		// Get record count
		iter = client.Single().Query(clientCtx, spanner.Statement{
			SQL: fmt.Sprintf("SELECT COUNT(*) FROM %s", table),
		})

		row, err := iter.Next()
		iter.Stop()

		if err != nil {
			client.Close()
			cancel()
			fmt.Printf("Error querying record count: %s\n", err.Error())
			continue
		}

		var count int64
		if err := row.Column(0, &count); err != nil {
			client.Close()
			cancel()
			fmt.Printf("Error reading count value: %s\n", err.Error())
			continue
		}

		client.Close()
		cancel()

		// Check if we have enough records
		if int(count) == records {
			fmt.Printf("\n%s has %d records\n", table, records)
			return nil
		}

		// Show progress
		if int(count) != prevCount {
			fmt.Printf("Found only %d records in the Spanner database ", count)
		} else {
			fmt.Print(".")
		}
		prevCount = int(count)
	}
}
