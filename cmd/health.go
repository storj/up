// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"database/sql"
	"fmt"
	"time"

	// imported for using postgres.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
)

var table string
var number int

func healthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "wait until cluster is healthy (10 storagenodes are registered in the db)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkHealth(table, number)
		},
	}
	cmd.PersistentFlags().StringVarP(&table, "table", "t", "nodes", "table to use for health check")
	cmd.PersistentFlags().IntVarP(&number, "number", "n", 10, "number of entries to expect in the table")
	return cmd
}

func init() {
	RootCmd.AddCommand(healthCmd())
}

// checkHealth polls the database until all storagenodes are checked in.
func checkHealth(table string, records int) error {
	prevCount := -1
	for {
		time.Sleep(1 * time.Second)
		db, err := sql.Open("pgx", "host=localhost port=26257 user=root dbname=master sslmode=disable")
		if err != nil {
			fmt.Printf("Couldn't connect to the database: %s\n", err.Error())
			continue
		}

		count, err := dbRecordCount(db, table)
		_ = db.Close()
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

func dbRecordCount(db *sql.DB, table string) (int, error) {
	row := db.QueryRow("select count(*) from " + table)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
