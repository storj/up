// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"os"
	"path"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/pkg/common"
)

var dir, subdir string

func localBinCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:     "local-bin <selector>",
		Aliases: []string{"local", "localbin"},
		Short:   "Use local compiled binares, bind-mounted to the containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
			if err != nil {
				return err
			}

			selector, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, mountBinaries, "", selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
	mountCmd.PersistentFlags().StringVarP(&dir, "dir", "d", path.Join(os.Getenv("GOPATH"), "bin"), "path where binaries are located")
	mountCmd.PersistentFlags().StringVarP(&subdir, "subdir", "s", "", "sub directory of the path where binaries are located")
	return mountCmd
}

func localWebSatelliteCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:     "local-websatellite /path/to/web/satellite",
		Aliases: []string{"local-ws"},
		Short:   "Use local web/satellite npm app for the satellite-api container.",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(common.ComposeFileName)
			if err != nil {
				return err
			}

			webSatPath, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, mountWebSatellite, webSatPath[0], []string{"satellite-api"})
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}

	return mountCmd
}

func init() {
	rootCmd.AddCommand(localBinCmd())
	rootCmd.AddCommand(localWebSatelliteCmd())
}

func mountBinaries(composeService *types.ServiceConfig, _ string) error {
	source := path.Join(path.Join(dir, subdir), common.BinaryDict[composeService.Name])
	target := path.Join("/var/lib/storj/go/bin", common.BinaryDict[composeService.Name])
	for i, volume := range composeService.Volumes {
		if volume.Type == "bind" &&
			volume.Target == target {
			composeService.Volumes = append(composeService.Volumes[:i], composeService.Volumes[i+1:]...)
		}
	}
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(source, target))
	return nil
}

func mountWebSatellite(composeService *types.ServiceConfig, webSatPath string) error {
	source := webSatPath
	target := "/var/lib/storj/storj/web/satellite/"
	for _, volume := range composeService.Volumes {
		if volume.Type == "bind" &&
			volume.Source == source &&
			volume.Target == target {
			return nil
		}
	}
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(source, target))
	return nil
}
