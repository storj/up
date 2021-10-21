package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

var setEnvCmd = &cobra.Command{
	Use:   "setenv [KEY=VALUE] [service ...]",
	Short: "Set environment variable / parameter in a container",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, SetEnv, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

var unsetEnvCmd = &cobra.Command{
	Use:   "unsetenv [KEY] [service ...]",
	Short: "Remove environment variable / parameter in a container",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, UnsetEnv, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(setEnvCmd)
	rootCmd.AddCommand(unsetEnvCmd)
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