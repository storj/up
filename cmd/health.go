package cmd

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func HealthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "wait until cluster is healthy (storagenodes are registered in the db)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkHealth(10)
		},
	}
}

func init() {
	rootCmd.AddCommand(HealthCmd())
}

func checkHealth(requiredStorageNodes int) error {
	for {
		time.Sleep(1 * time.Second)
		db, err := sql.Open("postgres", "host=localhost port=26257 user=root dbname=master sslmode=disable")
		if err != nil {
			fmt.Printf("Couldn't connect to the database: %s\n", err.Error())
			continue
		}

		defer db.Close()
		res, err := db.Query("select count(*) from nodes")
		if err != nil {
			fmt.Printf("Couldn't query the database: %s\n", err.Error())
			continue
		}
		defer res.Close()
		res.Next()
		var count int
		err = res.Scan(&count)
		if err != nil {
			fmt.Printf("Couldn't read results from the database: %s\n", err.Error())

		}
		if count == requiredStorageNodes {
			fmt.Println("Storj cluster is healthy")
			return nil
		}
		fmt.Printf("Found only %d storagenodes in the database\n", count)
	}
}