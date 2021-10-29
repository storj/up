package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/cmd/files/docker"
	"github.com/elek/sjr/cmd/files/templates"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

var repository, branch, ref string
var github = "github"
var gerrit = "gerrit"

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "build from a remote src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func GithubCmd() *cobra.Command {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "build github src repo for use inside the container",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := updateCompose(args, github)
			if err != nil {
				return err
			}
			for buildType := range common.ResolveBuilds(args) {
				err = updateDockerFile(buildType, github)
			}

			if err != nil {
				return err
			}
			return nil
		},
	}
	githubCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "main", "The branch to checkout and build")
	return githubCmd
}

func GerritCmd() *cobra.Command {
	gerritCmd := &cobra.Command{
		Use:   "gerrit",
		Short: "build gerrit src repo for use inside the container",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := updateCompose(args, gerrit)
			if err != nil {
				return err
			}
			for buildType := range common.ResolveBuilds(args) {
				err = updateDockerFile(buildType, gerrit)
			}
			if err != nil {
				return err
			}
			return nil
		},
	}
	gerritCmd.PersistentFlags().StringVarP(&ref, "refspec", "f", "", "The github refspec to checkout and build")
	gerritCmd.MarkPersistentFlagRequired("refspec")
	return gerritCmd
}

func init() {
	BuildCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(GithubCmd())
	remoteCmd.AddCommand(GerritCmd())
}

func updateCompose(args []string, remoteType string) error {
	composeProject, err := common.LoadComposeFromFile(ComposeFile)
	templateProject, err := common.LoadComposeFromBytes(templates.ComposeTemplate)
	buildServicesToAdd := []string{"app-base-dev", "app-base-ubuntu"}
	for _, arg := range args {
		buildServicesToAdd = append(buildServicesToAdd, common.BuildDict[arg])
	}
	updatedComposeProject, err := addBuildServices(composeProject, templateProject, buildServicesToAdd, remoteType, args)
	if err != nil {
		return err
	}
	return common.WriteComposeFile(updatedComposeProject)
}

func updateDockerFile(buildType string, remoteType string) error {
	buildArgs := map[string]*string{
		"TYPE":   &remoteType,
		"BRANCH": &branch,
		"REPO":   &repository,
		"REF":    &ref,
	}
	client, err := common.CreateClient()
	if err != nil {
		return err
	}
	var dockerBytes []byte
	if buildType == "edge" {
		dockerBytes = dockerfiles.StorjDocker
	} else if buildType == "storj" {
		dockerBytes = dockerfiles.EdgeDocker
	}
	common.BuildImage(client, []string{"ubuntubase"}, "base.Dockerfile", dockerfiles.BaseDocker, map[string]*string{})
	common.BuildImage(client, []string{"devbase"}, "build.Dockerfile", dockerfiles.BuildDocker, map[string]*string{})
	common.BuildImage(client, []string{buildType}, buildType+".Dockerfile", dockerBytes, buildArgs)
	return nil
}

func addBuildServices(compose *types.Project, template *types.Project, services []string, remoteType string, args []string) (*types.Project, error) {
	resolvedServices, err := common.ResolveServices(services)
	if err != nil {
		return nil, err
	}

	for _, service := range resolvedServices {
		newService, err := template.GetService(service)
		if err != nil {
			return nil, err
		}
		newService.Build.Args = map[string]*string{"TYPE": &remoteType}
		if !common.ContainsService(compose.Services, service) {
			compose.Services = append(compose.Services, newService)
		} else {
			for i, composeService := range compose.Services {
				if strings.EqualFold(composeService.Name, newService.Name) {
					compose.Services[i] = newService
				}
			}
		}
	}
	for _, arg := range args {
		common.UpdateEach(compose, setImage, strings.Split(common.BuildDict[arg], "-")[1], []string{arg})
	}
	return compose, nil
}
