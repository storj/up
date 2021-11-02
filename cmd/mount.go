package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"storj.io/storj-up/pkg/common"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var subdir string

func MountCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:   "mount",
		Short: "Use local compiled binares, bind-mounted to the containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadComposeFromFile(ComposeFile)
			if err != nil {
				return err
			}

			selector, _, err := common.ParseArgumentsWithSelector(args, 0)
			if err != nil {
				return err
			}

			updatedComposeProject, err := common.UpdateEach(composeProject, mountBinaries, "", selector)
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
	source := path.Join(path.Join(goBinPath, subdir), common.BinaryDict[composeService.Name])
	target := path.Join("/var/lib/storj/go/bin", common.BinaryDict[composeService.Name])
	for _, volume := range composeService.Volumes {
		if volume.Type == "bind" &&
			volume.Source == source &&
			volume.Target == target {
			return nil
		}
	}
	composeService.Volumes = append(composeService.Volumes, common.CreateBind(source, target))
	return nil
}
