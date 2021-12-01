// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/gateway-mt/pkg/auth"
	"storj.io/gateway-mt/pkg/linksharing"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj/satellite"
	"storj.io/storj/storagenode"
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
	configTypes := map[string]reflect.Type{
		"storagenode":     reflect.TypeOf(storagenode.Config{}),
		"satellite-api":   reflect.TypeOf(satellite.Config{}),
		"satellite-admin": reflect.TypeOf(satellite.Config{}),
		"satellite-core":  reflect.TypeOf(satellite.Config{}),
		"linksharing":     reflect.TypeOf(linksharing.Config{}),
		"authservice":     reflect.TypeOf(auth.Config{}),
	}
	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return err
	}

	emptySelection := true
	for _, s := range resolvedServices {
		if configType, found := configTypes[s]; found {
			printConfigStruct("STORJ", configType)
			fmt.Println()
			emptySelection = false
		}
	}
	if emptySelection {
		return errs.New("Couldn't find config type with selector %s. "+
			"Command is supported for the following services: %s",
			strings.Join(services, ","),
			strings.Join(keys(configTypes), ", "))
	}
	return nil
}

func keys(types map[string]reflect.Type) []string {
	var res []string
	for k := range types {
		res = append(res, k)
	}
	return res
}

func printConfigStruct(prefix string, configType reflect.Type) {
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			printConfigStruct(prefix+"_"+camelToUpperCase(field.Name), field.Type)
		} else {
			defaultValue := ""
			if field.Tag.Get("default") != "" {
				defaultValue = fmt.Sprintf("(default: %s)", field.Tag.Get("default"))
			}
			fmt.Printf("%-70s %s %s\n", prefix+"_"+camelToUpperCase(field.Name), field.Tag.Get("help"), defaultValue)
		}
	}

}

func camelToUpperCase(name string) string {
	smallCapital := regexp.MustCompile("([a-z])([A-Z])")
	name = smallCapital.ReplaceAllString(name, "${1}_$2")
	return strings.ToUpper(name)
}
