package recipe

import (
	"fmt"
	"github.com/elek/sjr/pkg/common"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"
	"strconv"
)

func scaleCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "scale <number>",
		Short: "Static scale of service",
		Long: "This command creates multiple instances from the save service. After this scale services couldn't be scaled up with `docker-compose scale any more`. " +
			"But also not required to scale up and down and it's possible to do per instance local bindmount",
		RunE: func(cmd *cobra.Command, args []string) error {
			instances, err := strconv.Atoi(args[0])
			if err != nil {
				return errs.Wrap(err)
			}
			return common.Update(service, func(compose *common.SimplifiedCompose) error {
				return Scale(service, compose, instances)
			})
		},
	})
}

func Scale(service string, compose *common.SimplifiedCompose, instances int) error {
	existing := compose.FilterPrefix(service)
	if len(existing) == 0 {
		return fmt.Errorf("Couldn't find any active instance from service %s", service)
	}
	if instances == 1 && len(existing) <= 1 {
		return fmt.Errorf("Couldn't scale down to 1 instances as there is only %d services", len(existing))
	}

	var template *common.ServiceConfig
	for _, v := range existing {
		template = v
		break
	}
	//delete existing
	for k, _ := range existing {
		delete(compose.Services, k)
	}
	//add new
	for i := 1; i <= instances; i++ {
		name := service
		if instances > 1 {
			name = fmt.Sprintf("%s%d", service, i)
		}
		//TODO: do a deep copy to have per/service configuration
		compose.Services[name] = template
	}
	return nil
}
