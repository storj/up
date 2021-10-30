package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

func SetArgCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setarg [KEY=VALUE] [service ...]",
		Short: "set container arg",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, SetArg, args[0], args[1:])
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func UnsetArgCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unsetarg [KEY] [service ...]",
		Short: "remove container arg",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, UnsetArg, args[0], args[1:])
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(SetArgCmd())
	rootCmd.AddCommand(UnsetArgCmd())
}

func SetArg(composeService *types.ServiceConfig, arg string) error {
	if composeService.Build == nil {
		composeService.Build = &types.BuildConfig{
			Args: map[string]*string{},
		}
	} else if composeService.Build.Args == nil {
		composeService.Build.Args = map[string]*string{}
	}
	parts := strings.SplitN(arg, "=", 2)
	composeService.Build.Args[parts[0]] = &parts[1]
	return nil
}

func UnsetArg(composeService *types.ServiceConfig, arg string) error {
	delete(composeService.Build.Args, arg)
	return nil
}
