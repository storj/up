package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
	"strconv"
)

var scaleCmd = &cobra.Command{
	Use:   "scale <number> [service ...]",
	Short: "Static scale of service or services",
	Long: "This command creates multiple instances of the service or services. After this scale services couldn't be scaled up with `docker-compose scale any more`. " +
		"But also not required to scale up and down and it's possible to do per instance local bindmount",
	RunE: func(cmd *cobra.Command, args []string) error {
		composeProject, err := common.UpdateEach(ComposeFile, scale, args[0], args[1:])
		if err != nil {
			return err
		}
		return common.WriteComposeFile(composeProject)
	},
}

func init() {
	rootCmd.AddCommand(scaleCmd)
}

func scale(composeService *types.ServiceConfig, scale string) error {
	instances, err := strconv.ParseUint(scale, 10, 64)
	if err != nil {
		return errs.Wrap(err)
	}
	if instances == 1 {
		composeService.Deploy = nil
	} else if composeService.Deploy == nil {
		composeService.Deploy = &types.DeployConfig{Replicas: &instances}
	} else {
		*composeService.Deploy.Replicas = instances
	}
	return nil
}