package cmd

import (
	"fmt"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var svcCmd = &cobra.Command{
	Use:   "services",
	Short: "return service names given in args",
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := common.ResolveServices(args)
		fmt.Println(services)
		return err
	},
}

func init() {
	rootCmd.AddCommand(svcCmd)
}