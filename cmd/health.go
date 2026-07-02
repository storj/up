// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

var table, host, user, dbname, dbtype string
var number, timeout, port int

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
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 5432, "port to use for health check")
	cmd.Flags().StringVarP(&user, "user", "u", "", "user to connect to the DB as (default: postgres for postgres, root for cockroach)")
	cmd.Flags().StringVarP(&dbname, "dbname", "", "master", "DB name to connect to")
	cmd.Flags().StringVarP(&dbtype, "dbtype", "", "postgres", "database type (postgres or cockroach)")
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
		if user == "" {
			user = "postgres"
			if strings.ToLower(dbtype) == "cockroach" {
				user = "root"
			}
		}
		return checkHealthCockroachPostgres(table, records)
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
