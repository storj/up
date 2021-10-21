package cmd

import (
	"fmt"
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

var entryPointCmd = &cobra.Command{
	Use:   "local-entrypoint [service ...]",
	Short: "Bind mount entrypoint.sh to use local modifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, updateEntryPoint, fmt.Sprintf("%s**%s", "./entrypoint.sh", "/var/lib/storj/entrypoint.sh"), args)
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(entryPointCmd)
}

func updateEntryPoint(composeService *types.ServiceConfig, arg string) error {
	parts := strings.SplitN(arg, "**", 2)
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(parts[0], parts[1]))
	return nil
}