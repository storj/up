// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

// Config contains all possible configuration for each known services.
var Config = map[string][]Option{}

// Option represents one possible configuration options for a service.
type Option struct {
	Name        string
	Description string
	Default     string
}
