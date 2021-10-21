package cmd

import (
	"fmt"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var svcCmd = &cobra.Command{
	Use:   "services",
	Short: "return services given in args",
	RunE: func(cmd *cobra.Command, args []string) error {
		return svcList(args)
	},
}

func init() {
	rootCmd.AddCommand(svcCmd)
}

func svcList(services []string) error {
	names, err := common.ResolveServices(services)
	fmt.Println(names)
	return err
}