package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/common"
)

func SvcCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "services <selector>",
		Short: "return service names given in args. Without argument it prints out all the possble service selectors",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) > 0 {
				resolvedServices, err := common.ResolveServices(args)
				if err != nil {
					return err
				}
				fmt.Println(strings.Join(resolvedServices, "\n"))
			} else {
				fmt.Println("Available services:")
				fmt.Println()
				for k, v := range common.GetSelectors() {
					if len(v) == 0 {
						fmt.Printf("%s\n", k)
					}

				}
				fmt.Println()
				fmt.Println("Available group selectors (and resolutions):")
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
