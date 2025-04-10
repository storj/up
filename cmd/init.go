// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/runtime"
	"storj.io/storj-up/pkg/runtime/standalone"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "init [<selector>...] OR init <compose|shell> [<selector>...]",
		Args: cobra.MinimumNArgs(1),
		Short: "Initialize new storj-up stack with the chosen container orchestrator. " + SelectorHelp + ". Without argument it generates " +
			"full Storj cluster with databases (db,minimal,edge)",
	}

	{
		composeCmd := &cobra.Command{
			Use:  "compose [<selector>...]",
			Args: cobra.MinimumNArgs(0),
		}
		composeCmd.RunE = func(cmd *cobra.Command, selector []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			n, err := compose.NewCompose(pwd)
			if err != nil {
				return err
			}
			st, err := recipe.GetStack()
			if err != nil {
				return err
			}
			err = runtime.ApplyRecipes(st, n, normalizedArgs(selector), 0)
			if err != nil {
				return err
			}

			return n.Write()
		}
		cmd.AddCommand(composeCmd)
		cmd.RunE = composeCmd.RunE
	}

	{
		shellCmd := &cobra.Command{
			Use:     "shell [<selector>...]",
			Args:    cobra.MinimumNArgs(0),
			Aliases: []string{"standalone"},
		}
		storjProjDir := shellCmd.Flags().StringP("storjdir", "s", "", "Directory of the storj code.")
		gatewayProjDir := shellCmd.Flags().StringP("gatewaydir", "g", "", "Directory of the gateway code.")
		shellCmd.RunE = func(cmd *cobra.Command, selector []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			storjProjectDir := os.Getenv("STORJ_PROJECT_DIR")
			if *storjProjDir != "" {
				storjProjectDir = *storjProjDir
			}
			if storjProjectDir == "" {
				return errs.Errorf("Please set \"STORJ_PROJECT_DIR\" environment variable or add -s flag with the location of your checked out storj/storj project. (Required to use web resources")
			}
			gatewayProjectDir := os.Getenv("GATEWAY_PROJECT_DIR")
			if *gatewayProjDir != "" {
				gatewayProjectDir = *gatewayProjDir
			}
			if gatewayProjectDir == "" {
				fmt.Println("WARNING: \"GATEWAY_PROJECT_DIR\" environment variable not set! Please set or add -g flag with the location of your checked out storj/gateway-mt project to use web resources.")
				gatewayProjectDir = "/tmp"
			}
			n, err := standalone.NewStandalone(standalone.Paths{
				ScriptDir:  pwd,
				StorjDir:   storjProjectDir,
				GatewayDir: gatewayProjectDir,
				CleanDir:   true,
			})
			if err != nil {
				return err
			}
			st, err := recipe.GetStack()
			if err != nil {
				return err
			}
			err = runtime.ApplyRecipes(st, n, normalizedArgs(selector), 0)
			if err != nil {
				return err
			}

			return n.Write()
		}
		cmd.AddCommand(shellCmd)
	}

	return cmd
}

func normalizedArgs(args []string) []string {
	var res []string
	for _, a := range args {
		for _, p := range strings.Split(a, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				res = append(res, p)
			}
		}
	}
	if len(res) == 0 {
		return []string{"db", "minimal", "edge"}
	}
	return res
}

func init() {
	RootCmd.AddCommand(initCmd())
}
