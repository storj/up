// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

func init() {
	Config["authservice"] = []Option{
		{
			Name:        "STORJ_ENDPOINT",
			Description: "Gateway endpoint URL to return to clients",
			Default:     "",
		}, {
			Name:        "STORJ_AUTH_TOKEN",
			Description: "auth security token to validate requests",
			Default:     "",
		}, {
			Name:        "STORJ_ALLOWED_SATELLITES",
			Description: "list of satellite NodeURLs allowed for incoming access grants",
			Default:     "https://www.storj.io/dcs-satellites",
		}, {
			Name:        "STORJ_CACHE_EXPIRATION",
			Description: "length of time satellite addresses are cached for",
			Default:     "10m",
		}, {
			Name:        "STORJ_GET_ACCESS_RATE_LIMITERS_MAX_REQS_SECOND",
			Description: "maximum number of allowed operations per second starting when first failure operation happens",
			Default:     "2",
		}, {
			Name:        "STORJ_GET_ACCESS_RATE_LIMITERS_BURST",
			Description: "maximum number of allowed operations to overpass the maximum operations per second",
			Default:     "3",
		}, {
			Name:        "STORJ_GET_ACCESS_RATE_LIMITERS_NUM_LIMITS",
			Description: "maximum number of keys/rate-limit pairs stored in the LRU cache",
			Default:     "1000",
		}, {
			Name:        "STORJ_GET_ACCESS_RATE_LIMITERS_ENABLED",
			Description: "indicates if rate-limiting for GetAccess endpoints is enabled",
			Default:     "false",
		}, {
			Name:        "STORJ_KVBACKEND",
			Description: "key/value store backend url",
			Default:     "",
		}, {
			Name:        "STORJ_MIGRATION",
			Description: "create or update the database schema, and then continue service startup",
			Default:     "false",
		}, {
			Name:        "STORJ_LISTEN_ADDR",
			Description: "public HTTP address to listen on",
			Default:     ":20000",
		}, {
			Name:        "STORJ_LISTEN_ADDR_TLS",
			Description: "public HTTPS address to listen on",
			Default:     ":20001",
		}, {
			Name:        "STORJ_DRPCLISTEN_ADDR",
			Description: "public DRPC address to listen on",
			Default:     ":20002",
		}, {
			Name:        "STORJ_DRPCLISTEN_ADDR_TLS",
			Description: "public DRPC+TLS address to listen on",
			Default:     ":20003",
		}, {
			Name:        "STORJ_LETS_ENCRYPT",
			Description: "use lets-encrypt to handle TLS certificates",
			Default:     "false",
		}, {
			Name:        "STORJ_CERT_FILE",
			Description: "server certificate file",
			Default:     "",
		}, {
			Name:        "STORJ_KEY_FILE",
			Description: "server key file",
			Default:     "",
		}, {
			Name:        "STORJ_PUBLIC_URL",
			Description: "public url for the server, for the TLS certificate",
			Default:     "",
		}, {
			Name:        "STORJ_DELETE_UNUSED_RUN",
			Description: "whether to run unused records deletion chore",
			Default:     "false",
		}, {
			Name:        "STORJ_DELETE_UNUSED_INTERVAL",
			Description: "interval unused records deletion chore waits to start next iteration",
			Default:     "24h",
		}, {
			Name:        "STORJ_DELETE_UNUSED_AS_OF_SYSTEM_INTERVAL",
			Description: "the interval specified in AS OF SYSTEM in unused records deletion chore query as negative interval",
			Default:     "5s",
		}, {
			Name:        "STORJ_DELETE_UNUSED_SELECT_SIZE",
			Description: "batch size of records selected for deletion at a time",
			Default:     "10000",
		}, {
			Name:        "STORJ_DELETE_UNUSED_DELETE_SIZE",
			Description: "batch size of records to delete from selected records at a time",
			Default:     "1000",
		},
	}
}
