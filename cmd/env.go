package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"storj.io/storj-up/pkg/common"
	"strings"
)

func SetEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setenv <selector> KEY=VALUE",
		Short: "Set environment variable / parameter in a container",
		Long:  selectorHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			if err != nil {
				return err
			}

			selector, args, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, SetEnv, args[0], selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func UnsetEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unsetenv <selector> KEY",
		Short: "remove environment variable / parameter in a container",
		Long:  "Remove environment variable from selected containers. " + selectorHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			if err != nil {
				return err
			}

			selector, args, err := common.ParseArgumentsWithSelector(args, 1)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, UnsetEnv, args[0], selector)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	envCmd := cobra.Command{
		Use:   "env",
		Short: "add/remove environment variables (configuration parameter) to specified services",
	}
	envCmd.AddCommand(SetEnvCmd())
	envCmd.AddCommand(UnsetEnvCmd())
	rootCmd.AddCommand(&envCmd)
}

func SetEnv(composeService *types.ServiceConfig, arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	composeService.Environment[parts[0]] = &parts[1]
	return nil
}

func UnsetEnv(composeService *types.ServiceConfig, arg string) error {
	delete(composeService.Environment, arg)
	return nil
}
