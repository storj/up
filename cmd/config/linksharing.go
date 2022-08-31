// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

func init() {
	Config["linksharing"] = []ConfigKey{
		{
			Name:        "STORJ_SERVER_NAME",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_ADDRESS",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_ADDRESS_TLS",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_TRAFFIC_LOGGING",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_TLSCONFIG",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_SHUTDOWN_TIMEOUT",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_URLBASES",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_TEMPLATES",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_STATIC_SOURCES_PATH",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_TXT_RECORD_TTL",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_BASE_URL",
			Description: "base url to use for resolving access key ids",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_TOKEN",
			Description: "auth token for giving access to the auth service",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_TIMEOUT",
			Description: "how long to wait for a single auth service connection",
			Default:     "10s",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_BACK_OFF_DELAY",
			Description: "The active time between retries, typically not set",
			Default:     "0ms",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_BACK_OFF_MAX",
			Description: "The maximum total time to allow retries",
			Default:     "5m",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_BACK_OFF_MIN",
			Description: "The minimum time between retries",
			Default:     "100ms",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_CACHE_EXPIRATION",
			Description: "how long to keep cached access grants in cache",
			Default:     "24h",
		}, {
			Name:        "STORJ_HANDLER_AUTH_SERVICE_CONFIG_CACHE_CAPACITY",
			Description: "how many cached access grants to keep in cache",
			Default:     "1000",
		}, {
			Name:        "STORJ_HANDLER_DNSSERVER",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_REDIRECT_HTTPS",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_LANDING_REDIRECT_TARGET",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_UPLINK",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_CONNECTION_POOL_CAPACITY",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_CONNECTION_POOL_KEY_CAPACITY",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_CONNECTION_POOL_IDLE_EXPIRATION",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_USE_QOS_AND_CC",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_CLIENT_TRUSTED_IPS_LIST",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_USE_CLIENT_IPHEADERS",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_STANDARD_RENDERS_CONTENT",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_HANDLER_STANDARD_VIEWS_HTML",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_GEO_LOCATION_DB",
			Description: "",
			Default:     "",
		},
	}
}
