// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package container

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

var dir, subdir, command, mountSource, mountTarget string
var mountService []string

type mount struct {
	service    []string
	targetPath string
}

var frontendDict = map[string]mount{
	"web.satellite": {
		service:    []string{"satellite-api"},
		targetPath: "/var/lib/storj/storj/web/satellite",
	},
	"web.multinode": {
		service:    []string{"storagenode"},
		targetPath: "/var/lib/storj/web/multinode",
	},
	"web.storagenode": {
		service:    []string{"storagenode"},
		targetPath: "/var/lib/storj/web/storagenode",
	},
	"admin.ui": {
		service:    []string{"satellite-admin"},
		targetPath: "/var/lib/storj/storj/satellite/admin/ui",
	},
}

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
	mountCmd.PersistentFlags().StringVarP(&command, "command", "c", "", "command to mount. If not specified, the name of the service will be used (eg. gateway-mt binary for the gateway-mt service)")
	return mountCmd
}

func localWebCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:     "local-websource -s /path/to/web/directory <optional selector>",
		Aliases: []string{"local-ws"},
		Short:   "Use local web/* npm app for the service container.",
		RunE: cmd.ExecuteStorjUP(func(st recipe.Stack, rt runtime.Runtime, args []string) error {
			if mountTarget == "" {
				err := resolveTarget()
				if err != nil {
					return err
				}
			}
			if len(args) != 0 {
				mountService = args
			}
			if len(mountService) == 0 {
				return errors.New("unable to determine service for mount. please provide the service as an argument")
			}
			return cmd.ChangeCompose(st, rt, mountService, mountWebDir)
		}),
	}
	mountCmd.PersistentFlags().StringVarP(&mountSource, "source", "s", "", "local path to the web directory")
	mountCmd.PersistentFlags().StringVarP(&mountTarget, "target", "t", "", "path where web directory will be mounted")
	_ = mountCmd.MarkPersistentFlagRequired("source")
	return mountCmd
}

func init() {
	cmd.RootCmd.AddCommand(localBinCmd())
	cmd.RootCmd.AddCommand(localWebCmd())
}

func mountBinaries(composeService *types.ServiceConfig) error {
	execName := BinaryDict[stripNumeric(composeService.Name)]
	if command != "" {
		execName = command
	}
	source := filepath.Join(path.Join(dir, subdir), execName)
	target := filepath.ToSlash(filepath.Join("/var/lib/storj/go/bin", execName))

	if err := common.IsRegularFile(source); err != nil {
		return err
	}

	for i, volume := range composeService.Volumes {
		if volume.Type == "bind" &&
			volume.Target == target {
			composeService.Volumes = append(composeService.Volumes[:i], composeService.Volumes[i+1:]...)
		}
	}
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(source, target))
	return nil
}

func resolveTarget() (err error) {
	for k, v := range frontendDict {
		match, err := regexp.MatchString(k, mountSource)
		if err != nil {
			return err
		}
		if match {
			mountTarget = v.targetPath
			mountService = v.service
			return nil
		}
	}
	return errors.New("unable to determine target mount directory. use -t to specify")
}

func mountWebDir(composeService *types.ServiceConfig) error {
	for _, volume := range composeService.Volumes {
		if volume.Type == "bind" &&
			volume.Source == mountSource &&
			volume.Target == mountTarget {
			return nil
		}
	}
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(mountSource, mountTarget))
	return nil
}

func stripNumeric(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		b := s[i]
		if b < '0' || '9' < b {
			return s[:i+1]
		}
	}
	return ""
}
