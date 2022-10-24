// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package dockerfiles

import _ "embed"

// StorjDocker is a Dockerfile for core services.
//
//go:embed storj.Dockerfile
var StorjDocker []byte

// EdgeDocker is a Dockerfile for edge services.
//
//go:embed edge.Dockerfile
var EdgeDocker []byte
