package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
)

var repository, branch, ref string

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "build from a remote src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "build github src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, buildRemoteGithubSrc, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

var gerritCmd = &cobra.Command{
	Use:   "gerrit",
	Short: "build gerrit src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, buildRemoteGerritSrc, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	buildCmd.AddCommand(remoteCmd)
	remoteCmd.AddCommand(githubCmd)
	remoteCmd.AddCommand(gerritCmd)

	githubCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "https://github.com/storj/{gateway-mt/storj}.git", "Git repository to clone before build.")
	githubCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "main", "The branch to checkout and build")
	gerritCmd.PersistentFlags().StringVarP(&ref, "refspec", "f", "", "The github refspec to checkout and build")
	gerritCmd.MarkPersistentFlagRequired("refspec")
}

func buildRemoteGithubSrc(_ *types.ServiceConfig, _ string) error {
	// do magic here
	return nil
}

	func buildRemoteGerritSrc(_ *types.ServiceConfig, _ string) error {
	// do magic here
	return nil
}