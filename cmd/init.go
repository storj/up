// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd/files/templates"
	"storj.io/storj-up/pkg/common"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use: "init [selector]",
		Short: "Generate docker-compose file with selected services. " + selectorHelp + ". Without argument it generates " +
			"full Storj cluster with databases (storj,db)",
		RunE: func(cmd *cobra.Command, args []string) error {
			selector, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}

			composeProject, err := initCompose(templates.ComposeTemplate, selector)
			if err != nil {
				return err
			}

			return common.WriteComposeFile(composeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(initCmd())
}

func initCompose(templateBytes []byte, services []string) (*types.Project, error) {
	templateComposeProject, err := common.LoadComposeFromBytes(templateBytes)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		services = []string{"storj", "db"}
	}
	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}

	compose, err := addToCompose(nil, templateComposeProject, resolvedServices)
	if err != nil {
		return nil, err
	}

	if len(compose.Services) == 0 {
		return nil, fmt.Errorf("no service is selected by selector \"%s\", please use `storj-up services` to check available service and group selectors to be used", strings.Join(services, ","))
	}

	templateComposeProject.Services = compose.Services

	return templateComposeProject, nil
}
