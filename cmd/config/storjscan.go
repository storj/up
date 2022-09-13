// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

func storjscanConfig() []Option {
	return []Option{

		{
			Name:        "STORJ_DEBUG_ADDRESS",
			Description: "",
			Default:     "",
		},
		{
			Name:        "STORJ_DEBUG_CONTROL_TITLE",
			Description: "",
			Default:     "",
		},
		{
			Name:        "STORJ_DEBUG_CONTROL",
			Description: "expose control panel",
			Default:     "",
		},
		{
			Name:        "STORJ_TOKENS_ENDPOINT",
			Description: "Ethereum RPC endpoint",
			Default:     "",
		},
		{
			Name:        "STORJ_TOKENS_CONTRACT",
			Description: "Address of the STORJ token to scan for transactions",
			Default:     "0xb64ef51c888972c908cfacf59b47c1afbc0ab8ac",
		},
		{
			Name:        "STORJ_TOKEN_PRICE_INTERVAL",
			Description: "how often to run the chore",
			Default:     "1m",
		},
		{
			Name:        "STORJ_TOKEN_PRICE_PRICE_WINDOW",
			Description: "max allowable duration between the requested and available ticker price timestamps",
			Default:     "1m",
		},
		{
			Name:        "STORJ_TOKEN_PRICE_COINMARKETCAP_CONFIG_BASE_URL",
			Description: "base URL for ticker price API",
			Default:     "https://pro-api.coinmarketcap.com",
		},
		{
			Name:        "STORJ_TOKEN_PRICE_COINMARKETCAP_CONFIG_APIKEY",
			Description: "API Key used to access coinmarketcap",
			Default:     "",
		},
		{
			Name:        "STORJ_TOKEN_PRICE_COINMARKETCAP_CONFIG_TIMEOUT",
			Description: "coinmarketcap API response timeout",
			Default:     "10s",
		},
		{
			Name:        "STORJ_TOKEN_PRICE_USE_TEST_PRICES",
			Description: "use test prices instead of coninmaketcap",
			Default:     "false",
		},
		{
			Name:        "STORJ_API_ADDRESS",
			Description: "public address to listen on",
			Default:     ":10000",
		},
		{
			Name:        "STORJ_API_KEYS",
			Description: "List of user:secret pairs to connect to service endpoints.",
			Default:     "",
		},
	}
}
