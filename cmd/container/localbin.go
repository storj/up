// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package container

import (
	"os"
	"path"
	"path/filepath"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

var dir, subdir, command string

// BinaryDict contains the executable app for each supported service.
var BinaryDict = map[string]string{
	"authservice":     "authservice",
	"gateway-mt":      "gateway-mt",
	"linksharing":     "linksharing",
	"satellite-admin": "satellite",
	"satellite-api":   "satellite",
	"satellite-core":  "satellite",
	"storagenode":     "storagenode",
	"uplink":          "uplink",
	"versioncontrol":  "versioncontrol",
	"storjscan":       "storjscan",
}

func localBinCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:     "local-bin <selector>",
		Aliases: []string{"local", "localbin"},
		Short:   "Use local compiled binares, bind-mounted to the containers.",
		RunE: cmd.ExecuteStorjUP(func(st recipe.Stack, rt runtime.Runtime, args []string) error {
			return cmd.ChangeCompose(st, rt, args, mountBinaries)
		}),
	}
	mountCmd.PersistentFlags().StringVarP(&dir, "dir", "d", filepath.Join(os.Getenv("GOPATH"), "bin"), "path where binaries are located")
	mountCmd.PersistentFlags().StringVarP(&subdir, "subdir", "s", "", "sub directory of the path where binaries are located")
	mountCmd.PersistentFlags().StringVarP(&command, "command", "", "", "command to mount. If not specified, the name of the service will be used (eg. gateway-mt binary for the gateway-mt service)")
	return mountCmd
}

func localWebSatelliteCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:     "local-websatellite /path/to/web/satellite",
		Aliases: []string{"local-ws"},
		Short:   "Use local web/satellite npm app for the satellite-api container.",
		Args:    cobra.ExactArgs(1),
		RunE: cmd.ExecuteStorjUP(func(st recipe.Stack, rt runtime.Runtime, args []string) error {
			return cmd.ChangeCompose(st, rt, []string{"satellite-api"}, func(composeService *types.ServiceConfig) error {
				return mountWebSatellite(composeService, args[0])
			})
		}),
	}

	return mountCmd
}

func init() {
	cmd.RootCmd.AddCommand(localBinCmd())
	cmd.RootCmd.AddCommand(localWebSatelliteCmd())
}

func mountBinaries(composeService *types.ServiceConfig) error {
	execName := BinaryDict[composeService.Name]
	if command != "" {
		execName = command
	}
	source := filepath.Join(path.Join(dir, subdir), execName)
	target := filepath.Join("/var/lib/storj/go/bin", execName)
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
