// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package dockerfiles

import (
	_ "embed"
	"strings"
)

// StorjDocker is a Dockerfile for core services.
//
//go:embed storj.Dockerfile
var StorjDocker []byte

// EdgeDocker is a Dockerfile for edge services.
//
//go:embed edge.Dockerfile
var EdgeDocker []byte

//go:embed build.last
var buildTag string

//go:embed base.last
var baseTag string

// BuildTag returns the current build image tag.
func BuildTag() string { return strings.TrimSpace(buildTag) }

// BaseTag returns the current base image tag.
func BaseTag() string { return strings.TrimSpace(baseTag) }
