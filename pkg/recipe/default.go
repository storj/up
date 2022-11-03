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

//go:embed edge.yaml
var edge []byte

//go:embed tracing.yaml
var tracing []byte

//go:embed billing.yaml
var billing []byte

//go:embed gc.yaml
var gc []byte

// Defaults is a map for recipes included in the binary.
var Defaults = map[string][]byte{
	"minimal": minimal,
	"db":      db,
	"edge":    edge,
	"tracing": tracing,
	"billing": billing,
	"gc":      gc,
}
