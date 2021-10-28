package cmd

import (
	"fmt"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var svcCmd = &cobra.Command{
	Use:   "services <group or service names>",
	Short: "return service names given in args",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := common.ResolveServices(args)
		fmt.Println(services)
		return err
	},
}

func init() {
	rootCmd.AddCommand(svcCmd)
}
