// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	osruntime "runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
	"golang.org/x/exp/slices"

	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "build and start all services",
		RunE: func(cmd *cobra.Command, args []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return errs.Wrap(err)
			}
			runtime, err := FromDir(pwd)
			if err != nil {
				return errs.Wrap(err)
			}
			st, err := recipe.GetStack()
			if err != nil {
				return errs.Wrap(err)
			}

			err = runtime.Reload(st)
			if err != nil {
				return errs.Wrap(err)
			}
			err = buildServices(runtime.GetServices())
			if err != nil {
				return errs.Wrap(err)
			}
			err = startServices()
			if err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}

func init() {
	RootCmd.AddCommand(startCmd())
}

func buildServices(services []runtime.Service) error {
	var buildpaths []string
	for _, service := range services {
		for _, mount := range service.GetVolumes() {
			if mount.MountType == "bind" && filepath.Dir(mount.Target) == filepath.Clean("/var/lib/storj/go/bin") {
				var codeSource string
				var foundCodeSource bool
				serviceCodeSource, foundCodeSource := service.GetENV()["STORJ_UP_LOCAL_BINARY_SOURCE"]
				if foundCodeSource {
					codeSource = *serviceCodeSource
				} else {
					codeSource, foundCodeSource = os.LookupEnv("STORJ_UP_LOCAL_" + service.ID().Name)
				}
				if foundCodeSource && len(codeSource) > 0 && !slices.Contains(buildpaths, codeSource) {
					buildpaths = append(buildpaths, codeSource)
					var cmd *exec.Cmd
					if osruntime.GOOS != "linux" {
						cmd = exec.Command("go", "install")
						cmd.Env = os.Environ()
						cmd.Env = append(cmd.Env, "GOOS=linux")
						cmd.Env = append(cmd.Env, "GOARCH=amd64")
						if strings.Contains(strings.ToLower(service.ID().Name), "storagenode") {
							cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
						}
					} else {
						cmd = exec.Command("go", "install")
					}
					if cmd == nil {
						return errors.New("unable to run go build command")
					}

					err := runCommand(cmd, filepath.Clean(codeSource))
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func startServices() error {
	cmd := exec.Command("docker", "compose", "up", "-d")
	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	return runCommand(cmd, curDir)
}

func runCommand(cmd *exec.Cmd, runPath string) error {

	curDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(runPath)
	if err != nil {
		return err
	}

	stderr, _ := cmd.StderrPipe()
	fmt.Println("*** Storj-Up Running " + strings.Join(cmd.Args, " ") + " from " + runPath + " ***")
	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}

	err = os.Chdir(curDir)
	if err != nil {
		return err
	}

	return nil
}
