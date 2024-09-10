// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package recipe

import (
	_ "embed"
)

//go:embed minimal.yaml
var minimal []byte

//go:embed db.yaml
var db []byte

//go:embed postgres.yaml
var postgres []byte

//go:embed spanner.yaml
var spanner []byte

//go:embed edge.yaml
var edge []byte

//go:embed core.yaml
var satelliteCore []byte

//go:embed admin.yaml
var satelliteAdmin []byte

//go:embed tracing.yaml
var tracing []byte

//go:embed billing.yaml
var billing []byte

//go:embed gc.yaml
var gc []byte

//go:embed rangedloop.yaml
var rangedLoop []byte

//go:embed uplink.yaml
var uplink []byte

//go:embed repair.yaml
var repair []byte

//go:embed audit.yaml
var audit []byte

//go:embed mailserver.yaml
var mailserver []byte

// Defaults is a map for recipes included in the binary.
var Defaults = map[string][]byte{
	"minimal":         minimal,
	"satellite-core":  satelliteCore,
	"satellite-admin": satelliteAdmin,
	"db":              db,
	"postgres":        postgres,
	"spanner":         spanner,
	"edge":            edge,
	"tracing":         tracing,
	"billing":         billing,
	"gc":              gc,
	"ranged-loop":     rangedLoop,
	"uplink":          uplink,
	"repair":          repair,
	"audit":           audit,
	"mailserver":      mailserver,
}
