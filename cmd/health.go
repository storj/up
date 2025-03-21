// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	// imported for using postgres.
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
)

var table, host, user, dbname string
var number, timeout, port int

func healthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Args:  cobra.NoArgs,
		Short: "wait until cluster is healthy (10 storagenodes are registered in the db)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return checkHealthWithTimeout(table, number, timeout)
		},
	}
	cmd.PersistentFlags().StringVarP(&table, "table", "t", "nodes", "table to use for health check")
	cmd.PersistentFlags().IntVarP(&number, "number", "n", 10, "number of entries to expect in the table")
	cmd.PersistentFlags().IntVarP(&timeout, "duration", "d", 0, "time to wait (in seconds) for health check")
	cmd.PersistentFlags().StringVarP(&host, "host", "", "127.0.0.1", "host/ip address to use for health check. Defaults to 127.0.0.1 or STORJ_DOCKER_HOST if set.")
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 26257, "port to use for health check")
	cmd.Flags().StringVarP(&user, "user", "u", "root", "user to connect to the DB as")
	cmd.Flags().StringVarP(&dbname, "dbname", "", "master", "DB name to connect to")
	return cmd
}

func init() {
	RootCmd.AddCommand(healthCmd())
}

// checkHealthWithTimeout polls the database until all storagenodes are checked in, or the timeout is exceeded. a timeout of 0 (default)
// means no timeout.
func checkHealthWithTimeout(table string, records, timeout int) error {
	if timeout == 0 {
		return checkHealth(table, records)
	}
	c1 := make(chan error, 1)
	go func() {
		err := checkHealth(table, records)
		if err != nil {
			c1 <- err
		}
		c1 <- nil
	}()

	select {
	case err := <-c1:
		if err != nil {
			return err
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		return fmt.Errorf("health check failed. duration limit reached")
	}
	return nil
}

// checkHealth polls the database until all storagenodes are checked in.
func checkHealth(table string, records int) error {
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

		count, err := dbRecordCount(db, table)
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

func dbRecordCount(db *pgx.Conn, table string) (int, error) {
	row := db.QueryRow(context.TODO(), "select count(*) from "+table)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
