// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/storj-up/cmd/config"
	"storj.io/storj-up/pkg/common"
)

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "configs <selector>",
		Aliases: []string{"config"},
		Short:   "Print out available configuration for specific services",
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
	rootCmd.AddCommand(configCmd())
}

func printConfigs(services []string) error {
	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return err
	}

	emptySelection := true
	for _, s := range resolvedServices {
		if configs, found := config.Config[s]; found {
			printConfigStruct(configs)
			fmt.Println()
			emptySelection = false
		}
	}
	if emptySelection {
		return errs.New("Couldn't find config type with selector %s. "+
			"Command is supported for the following services: %s",
			strings.Join(services, ","),
			strings.Join(keys(config.Config), ", "))
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
