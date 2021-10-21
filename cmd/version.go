package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"strings"
)

var versionCmd = &cobra.Command{
	Use:   "version <version> [service ...]",
	Short: "Set version (docker image tag) for specified services",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, updateVersion, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func updateVersion(composeService *types.ServiceConfig, version string) error {
	newImage := strings.ReplaceAll(composeService.Image, "@sha256", "")
	newImage = strings.Split(newImage, ":")[0] + ":" + version
	composeService.Image = newImage
	return nil
}