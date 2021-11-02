package cmd

import (
	"fmt"
	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"storj.io/storj-up/pkg/common"
	"strings"
)

func EntryPointCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "local-entrypoint [service ...]",
		Short: "bind mount entrypoint.sh to use local modifications",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, updateEntryPoint, fmt.Sprintf("%s**%s", "./entrypoint.sh", "/var/lib/storj/entrypoint.sh"), args)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(EntryPointCmd())
}

func updateEntryPoint(composeService *types.ServiceConfig, arg string) error {
	parts := strings.SplitN(arg, "**", 2)
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(parts[0], parts[1]))
	return nil
}
