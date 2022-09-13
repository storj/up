// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

// All returns with all known options for registered services.
func All() map[string][]Option {
	return map[string][]Option{
		"authservice":     authserviceConfig(),
		"linksharing":     linksharingConfig(),
		"satellite-admin": satelliteadminConfig(),
		"satellite-api":   satelliteapiConfig(),
		"satellite-core":  satellitecoreConfig(),
		"storagenode":     storagenodeConfig(),
		"storjscan":       storjscanConfig(),
	}
}
