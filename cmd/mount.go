package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var subdir string

func MountCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:   "mount",
		Short: "mount local binaries to the default docker image",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadCompose(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, mountBinaries, "", args)
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
	mountCmd.PersistentFlags().StringVarP(&subdir, "subdir", "s", "", "sub directory of the go/bin folder where binaries are located")
	return mountCmd
}

func init() {
	rootCmd.AddCommand(MountCmd())
}

func mountBinaries(composeService *types.ServiceConfig, _ string) error {
	goBinPath := path.Join(os.Getenv("GOPATH"), "bin")
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(
		path.Join(
			path.Join(goBinPath, subdir),
			common.BinaryDict[composeService.Name]),
		path.Join("/var/lib/storj/go/bin", common.BinaryDict[composeService.Name])))
	return nil
}