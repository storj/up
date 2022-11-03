// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package container

import (
	"strconv"

	"github.com/compose-spec/compose-go/types"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/cmd"
	"storj.io/storj-up/pkg/recipe"
	"storj.io/storj-up/pkg/runtime/runtime"
)

func scaleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scale <selector> <number>",
		Short: "static scale of services",
		Args:  cobra.MinimumNArgs(2),
		Long: "This command creates multiple instances of the service or services. After this scale services couldn't be scaled up with `docker compose scale` any more. " +
			"But also not required to scale up and down and it's possible to do per instance local bindmount",
		RunE: cmd.ExecuteStorjUP(func(st recipe.Stack, rt runtime.Runtime, args []string) error {
			return cmd.ChangeCompose(st, rt, args[:len(args)-1], func(composeService *types.ServiceConfig) error {
				return scale(composeService, args[0])
			})
		}),
	}
}

func init() {
	cmd.RootCmd.AddCommand(scaleCmd())
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
