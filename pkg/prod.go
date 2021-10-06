package sjr

import "github.com/spf13/cobra"

func prodCmd(service string, command *cobra.Command) {
	command.AddCommand(&cobra.Command{
		Use:   "userprod",
		Short: "Use production satellite with local edge services",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Update(service, func(compose *SimplifiedCompose) error {
				return UserProd(compose, args[0])
			})
		},
	})
}

func UserProd(service *SimplifiedCompose, region string) error {
	_, hasAuthservice := service.Services["authservice"]

	delete(service.Services["gateway-mt"].Environment, "STORJ_WAIT_FOR_SATELLITE")
	delete(service.Services["authservice"].Environment, "STORJ_WAIT_FOR_SATELLITE")
	delete(service.Services["linksharing"].Environment, "STORJ_WAIT_FOR_SATELLITE")
	if !hasAuthservice {
		authUrl := "https://auth." + region + ".storjshare.io"
		service.Services["gateway-mt"].Environment["STORJ_AUTH_URL"] = &authUrl
	}

	gatewayUrl := "https://eu1.storj.io"
	satellite := "12L9ZFwhzVpuEKMUNUqkaTLGzwY9G24tbiigLiXpmZWKwmcNDDs@eu1.storj.io:7777"
	service.Services["authservice"].Environment["STORJ_ALLOWED_SATELLITES"] = &satellite
	service.Services["authservice"].Environment["STORJ_ENDPOINT"] = &gatewayUrl
	return nil
}
