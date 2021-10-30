package cmd

import (
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

func updateCompose(services []string, remoteType string) error {
	composeProject, err := common.LoadComposeFromFile(ComposeFile)
	if err != nil {
		return err
	}
	templateProject, err := common.LoadComposeFromBytes(templates.ComposeTemplate)
	if err != nil {
		return err
	}
	for buildType := range common.ResolveBuilds(services) {
		AddToCompose(composeProject, templateProject, []string{buildType})
		if err != nil {
			return err
		}
		for i, service := range composeProject.AllServices() {
			if strings.EqualFold(service.Name, buildType) {
				SetArg(&composeProject.Services[i], "TYPE="+remoteType)
				if remoteType == github {
					SetArg(&composeProject.Services[i], "BRANCH="+branch)
				} else if remoteType == gerrit {
					SetArg(&composeProject.Services[i], "REF="+ref)
				}
			}
		}
	}
	for _, service := range common.ResolveServices(services) {
		for i, composeService := range composeProject.AllServices() {
			if strings.EqualFold(composeService.Name, service) {
				SetImage(&composeProject.Services[i], strings.Split(common.BuildDict[service], "-")[1])
			}
		}
	}
	return common.WriteComposeFile(composeProject)
}