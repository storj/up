package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"storj.io/storj-up/pkg/common"
	"strings"
)

func SvcCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "services <selector>",
		Short: "return service names given in args. Without argument it prints out all the possble service selectors",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) > 0 {
				fmt.Println(strings.Join(common.ResolveServices(args), "\n"))
			} else {
				for k, v := range common.GetSelectors() {
					if len(v) == 0 {
						fmt.Printf("%s\n", k)
					}

				}
				fmt.Println()

				for k, v := range common.GetSelectors() {
					if len(v) > 0 {
						fmt.Printf("%s => %s\n", k, v)
					}

				}
			}

			return nil
		},
	}
}

func init() {
	rootCmd.AddCommand(SvcCmd())
}
