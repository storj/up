package cmd

import (
	"github.com/compose-spec/compose-go/types"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
	"strconv"
)

func ScaleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scale <number> [service ...]",
		Short: "static scale of services",
		Long: "This command creates multiple instances of the service or services. After this scale services couldn't be scaled up with `docker-compose scale any more`. " +
			"But also not required to scale up and down and it's possible to do per instance local bindmount",
		RunE: func(cmd *cobra.Command, args []string) error {
			composeProject, err := common.LoadCompose(ComposeFile)
			updatedComposeProject, err := common.UpdateEach(composeProject, Scale, args[0], args[1:])
			if err != nil {
				return err
			}
			return common.WriteComposeFile(updatedComposeProject)
		},
	}
}

func init() {
	rootCmd.AddCommand(ScaleCmd())
}

func Scale(composeService *types.ServiceConfig, scale string) error {
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