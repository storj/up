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
	_, hasGateway := service.Services["gateway-mt"]
	if hasGateway {
		delete(service.Services["gateway-mt"].Environment, "STORJ_WAIT_FOR_SATELLITE")
		if !hasAuthservice {
			authUrl := "https://auth." + region + ".storjshare.io"
			service.Services["gateway-mt"].Environment["STORJ_AUTH_URL"] = &authUrl
		}
	}
	return nil
}
