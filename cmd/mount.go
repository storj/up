package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "mount local binaries to the default docker image",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, mountBinaries, "", args)
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(mountCmd)
}

func mountBinaries(composeService *types.ServiceConfig, _ string) error {
	goBinPath := path.Join(os.Getenv("GOPATH"), "bin")
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(path.Join(goBinPath, composeService.Name), path.Join("/var/lib/storj/go/bin", composeService.Name)))
	return nil
}