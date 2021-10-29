package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

func SetEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setenv [KEY=VALUE] [service ...]",
		Short: "set environment variable / parameter in a container",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, SetEnv, args[0], args[1:])
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func UnsetEnvCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unsetenv [KEY] [service ...]",
		Short: "remove environment variable / parameter in a container",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, UnsetEnv, args[0], args[1:])
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(SetEnvCmd())
	rootCmd.AddCommand(UnsetEnvCmd())
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
