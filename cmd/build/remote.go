// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package build

import (
	"strings"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/common"
	dockerfiles "storj.io/storj-up/pkg/files/docker"
	"storj.io/storj-up/pkg/files/templates"
)

var branch, commit, ref string

const (
	github = "github"
	gerrit = "gerrit"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "build from a remote src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func githubCmd() *cobra.Command {
	// NOTE cobra doesn't have a way to document positional parameters:
	// https://github.com/spf13/cobra/issues/378
	githubCmd := &cobra.Command{
		Use:   "github <selector>",
		Short: "build github src repo for use inside the container",
		Long: `build github src repo for use inside the container for the indicated
services through positional arguments. See the list of supported service running
` + "`storj-up services`.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := updateCompose(args, github)
			if err != nil {
				return err
			}
			return nil
		},
	}
	githubCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "main", "The branch to checkout and build")
	githubCmd.PersistentFlags().StringVarP(&commit, "commit", "c", "", "The commit to checkout and build")
	return githubCmd
}

func gerritCmd() *cobra.Command {
	// NOTE cobra doesn't have a way to document positional parameters:
	// https://github.com/spf13/cobra/issues/378
	gerritCmd := &cobra.Command{
		Use:   "gerrit <selector>",
		Short: "build gerrit src repo for use inside the container",
		Long: `build gerrit src repo for use inside the container for the indicated
services through positional arguments. See the list of supported service running
` + "`storj-up services`.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := updateCompose(args, gerrit)
			if err != nil {
				return err
			}
			return nil
		},
	}
	gerritCmd.PersistentFlags().StringVarP(&ref, "refspec", "f", "", "The gerrit refspec to checkout and build")
	_ = gerritCmd.MarkPersistentFlagRequired("refspec")
	return gerritCmd
}

func init() {
	buildCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(githubCmd())
	remoteCmd.AddCommand(gerritCmd())
}

func updateCompose(services []string, remoteType string) error {
	err := common.ExtractFile("", "storj.Dockerfile", dockerfiles.StorjDocker)
	if err != nil {
		return err
	}

	err = common.ExtractFile("", "edge.Dockerfile", dockerfiles.EdgeDocker)
	if err != nil {
		return err
	}

	composeProject, err := common.LoadComposeFromFile("./", common.ComposeFileName)
	if err != nil {
		return err
	}
	templateProject, err := common.LoadComposeFromBytes(templates.ComposeTemplate)
	if err != nil {
		return err
	}

	resolvedBuilds, err := common.ResolveBuilds(services)
	if err != nil {
		return err
	}

	for buildType := range resolvedBuilds {
		_, err = addToCompose(composeProject, templateProject, []string{buildType})
		if err != nil {
			return err
		}
		for i, service := range composeProject.AllServices() {
			if strings.EqualFold(service.Name, buildType) {
				err = setArg(&composeProject.Services[i], "TYPE="+remoteType)
				if err != nil {
					return errs.Wrap(err)
				}
				switch remoteType {
				case github:
					err = setArg(&composeProject.Services[i], "BRANCH="+branch)
					if err != nil {
						return errs.Wrap(err)
					}
					err = setArg(&composeProject.Services[i], "SOURCE=branch")
					if err != nil {
						return errs.Wrap(err)
					}
					if commit != "" {
						err = setArg(&composeProject.Services[i], "COMMIT="+commit)
						if err != nil {
							return errs.Wrap(err)
						}
						err = setArg(&composeProject.Services[i], "SOURCE=commit")
						if err != nil {
							return errs.Wrap(err)
						}
					}
				case gerrit:
					err = setArg(&composeProject.Services[i], "REF="+ref)
					if err != nil {
						return errs.Wrap(err)
					}
					err = setArg(&composeProject.Services[i], "SOURCE=none")
					if err != nil {
						return errs.Wrap(err)
					}
				default:
					return errs.Errorf("Unsupported remote: %s", remoteType)
				}
			}
		}
	}

	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return err
	}

	for _, service := range resolvedServices {
		for i, composeService := range composeProject.AllServices() {
			if common.ServiceMatches(composeService.Name, service) {
				composeProject.Services[i].Image = strings.Split(common.BuildDict[service], "-")[1]
			}
		}
	}
	return common.WriteComposeFile(".", composeProject)
}

func addToCompose(compose *types.Project, template *types.Project, services []string) (*types.Project, error) {
	if compose == nil {
		compose = &types.Project{Services: []types.ServiceConfig{}}
	}

	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}
	for _, service := range resolvedServices {
		if !common.ContainsService(compose.Services, service) {
			newService, err := template.GetService(service)
			if err != nil {
				return nil, err
			}
			compose.Services = append(compose.Services, newService)
		}
	}
	return compose, nil
}
