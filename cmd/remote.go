package cmd

import (
	"github.com/spf13/cobra"
)

var repository, branch, ref string

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "build from a remote src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "build github src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildRemoteGithubSrc()
	},
}

var gerritCmd = &cobra.Command{
	Use:   "gerrit",
	Short: "build gerrit src repo for use inside the container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return buildRemoteGerritSrc()
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

func buildRemoteGithubSrc() error {
	// do magic here
	return nil
}

func buildRemoteGerritSrc() error {
	// do magic here
	return nil
}