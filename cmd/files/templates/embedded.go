// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package templates

import _ "embed"

// ComposeTemplate represents the canonical docker-compose with all the possible services.
//go:embed docker-compose.template.yaml
var ComposeTemplate []byte

// PrometheusYaml represents an example prometheus config.
//go:embed prometheus.yml
var PrometheusYaml []byte
