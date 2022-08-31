// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

var Config = map[string][]ConfigKey{}

type ConfigKey struct {
	Name        string
	Description string
	Default     string
}
