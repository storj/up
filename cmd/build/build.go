// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package build

import (
	"strings"

	"github.com/compose-spec/compose-go/v2/types"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/common"
	dockerfiles "storj.io/storj-up/pkg/files/docker"
	"storj.io/storj-up/pkg/files/templates"
)

var skipFrontend bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Args:  cobra.NoArgs,
	Short: "Build image on-the-fly instead of using pre-baked image",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	buildCmd.PersistentFlags().BoolVarP(&skipFrontend, "skipfrontend", "s", false, "Skip building the frontend")
	cmd.RootCmd.AddCommand(buildCmd)
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
		for serviceName, service := range composeProject.Services {
			if strings.EqualFold(service.Name, buildType) {
				err = setArg(&service, "TYPE="+remoteType)
				if err != nil {
					return errs.Wrap(err)
				}
				if skipFrontend {
					err = setArg(&service, "SKIP_FRONTEND_BUILD=true")
					if err != nil {
						return errs.Wrap(err)
					}
				}
				switch remoteType {
				case github:
					err = setArg(&service, "BRANCH="+branch)
					if err != nil {
						return errs.Wrap(err)
					}
					err = setArg(&service, "SOURCE=branch")
					if err != nil {
						return errs.Wrap(err)
					}
					if commit != "" {
						err = setArg(&service, "COMMIT="+commit)
						if err != nil {
							return errs.Wrap(err)
						}
						err = setArg(&service, "SOURCE=commit")
						if err != nil {
							return errs.Wrap(err)
						}
					}
				case gerrit:
					err = setArg(&service, "REF="+ref)
					if err != nil {
						return errs.Wrap(err)
					}
					err = setArg(&service, "SOURCE=none")
					if err != nil {
						return errs.Wrap(err)
					}
				case local:
					if path == "" {
						path = "."
					}
					err = setArg(&service, "SOURCE=none")
					if err != nil {
						return errs.Wrap(err)
					}
					err = setArg(&service, "PATH="+path)
					if err != nil {
						return errs.Wrap(err)
					}
				default:
					return errs.Errorf("Unsupported remote: %s", remoteType)
				}
				composeProject.Services[serviceName] = service
			}
		}
	}

	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return err
	}

	for _, service := range resolvedServices {
		for serviceName, composeService := range composeProject.Services {
			if common.ServiceMatches(composeService.Name, service) {
				composeService.Image = strings.Split(common.BuildDict[service], "-")[1]
				composeProject.Services[serviceName] = composeService
			}
		}
	}
	return common.WriteComposeFile(".", composeProject)
}

func addToCompose(compose *types.Project, template *types.Project, services []string) (*types.Project, error) {
	if compose == nil {
		compose = &types.Project{Services: make(types.Services)}
	}
	if compose.Services == nil {
		compose.Services = make(types.Services)
	}

	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}
	for _, service := range resolvedServices {
		if !containsService(compose.Services, service) {
			newService, err := template.GetService(service)
			if err != nil {
				return nil, err
			}
			compose.Services[service] = newService
		}
	}
	return compose, nil
}

// containsService checks if the service is included in the services map.
func containsService(services types.Services, name string) bool {
	for _, s := range services {
		if s.Name == name {
			return true
		}
	}
	return false
}
