// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package modify

import (
	"strings"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version <selector>... <version>",
		Short: "set version (docker image tag) for specified services",
		Args:  cobra.MinimumNArgs(2),
		RunE: cmd.ExecuteStorjUP(func(st recipe.Stack, rt runtime.Runtime, args []string) error {
			selector, version := common.SplitArgsSelector1(args)
			return cmd.ChangeCompose(st, rt, selector, func(composeService *types.ServiceConfig) error {
				return updateVersion(composeService, version)
			})
		}),
	}
}

func init() {
	cmd.RootCmd.AddCommand(versionCmd())
}

func updateVersion(composeService *types.ServiceConfig, version string) error {
	newImage := strings.ReplaceAll(composeService.Image, "@sha256", "")
	newImage = strings.Split(newImage, ":")[0] + ":" + version
	composeService.Image = newImage
	return nil
}
