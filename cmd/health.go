// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"database/sql"
	"fmt"
	"time"

	// imported for using postgres.
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
)

func healthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "wait until cluster is healthy (10 storagenodes are registered in the db)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkHealth(10)
		},
	}
}

func init() {
	rootCmd.AddCommand(healthCmd())
}

// checkHealth polls the database until all storagenodes are checked in.
func checkHealth(requiredStorageNodes int) error {
	for {
		time.Sleep(1 * time.Second)
		db, err := sql.Open("postgres", "host=localhost port=26257 user=root dbname=master sslmode=disable")
		if err != nil {
			fmt.Printf("Couldn't connect to the database: %s\n", err.Error())
			continue
		}

		count, err := registeredNodeCount(db)
		_ = db.Close()
		if err != nil {
			fmt.Printf("Couldn't query database for nodes: %s\n", err.Error())
			continue
		}
		if count == requiredStorageNodes {
			fmt.Println("Storj cluster is healthy")
			return nil
		}
		fmt.Printf("Found only %d storagenodes in the database\n", count)
	}
}

func registeredNodeCount(db *sql.DB) (int, error) {
	row := db.QueryRow("select count(*) from nodes")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
