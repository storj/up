// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/common"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version <selector> <version>",
		Short: "set version (docker image tag) for specified services",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(composeFile)
			if err != nil {
				return err
			}

			selector, parsedArgs, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, updateVersion, parsedArgs[0], selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(versionCmd())
}

func updateVersion(composeService *types.ServiceConfig, version string) error {
	newImage := strings.ReplaceAll(composeService.Image, "@sha256", "")
	newImage = strings.Split(newImage, ":")[0] + ":" + version
	composeService.Image = newImage
	return nil
}
