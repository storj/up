// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/config"
)

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "configs <service>",
		Aliases: []string{"config"},
		Short:   "Print out available configuration for specific service",
		RunE: func(cmd *cobra.Command, args []string) error {
			selector, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}
			return printConfigs(selector)
		},
	}
}

func init() {
	cmd.RootCmd.AddCommand(configCmd())
}

func printConfigs(services []string) error {
	emptySelection := true
	allOptions := config.All()

	for _, s := range services {
		if configs, found := allOptions[s]; found {
			printConfigStruct(configs)
			fmt.Println()
			emptySelection = false
		}
	}
	if emptySelection {
		return errs.New("Couldn't find config type with service name %s. "+
			"Command is supported for the following services: %s",
			strings.Join(services, ","),
			strings.Join(keys(allOptions), ", "))
	}
	return nil
}

func printConfigStruct(configs []config.Option) {
	for _, c := range configs {
		def := ""
		if c.Default != "" {
			def = fmt.Sprintf("(default: %s)", c.Default)
		}
		fmt.Printf("%-70s %s %s\n", c.Name, c.Description, def)
	}
}

func keys(types map[string][]config.Option) []string {
	var res []string
	for k := range types {
		res = append(res, k)
	}
	return res
}
