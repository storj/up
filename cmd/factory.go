// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/runtime/compose"
	"storj.io/storj-up/pkg/runtime/nomad"
	"storj.io/storj-up/pkg/runtime/runtime"
	"storj.io/storj-up/pkg/runtime/standalone"
)

// FromDir creates the right runtime based on available file names in the directory.
func FromDir(dir string) (runtime.Runtime, error) {
	_, err := os.Stat(filepath.Join(dir, "docker-compose.yaml"))
	if err == nil {
		return compose.NewCompose(dir)
	}

	_, err = os.Stat(filepath.Join(dir, "storj.hcl"))
	if err == nil {
		return nomad.NewNomad(dir, "storj")
	}

	_, err = os.Stat(filepath.Join(dir, "supervisord.conf"))
	if err == nil {
		storjProjectDir := os.Getenv("STORJ_PROJECT_DIR")
		if storjProjectDir == "" {
			return nil, errs.Errorf("Please set \"STORJ_PROJECT_DIR\" environment variable with the location of your checked out storj/storj project. (Required to use web resources")
		}
		gatewayProjectDir := os.Getenv("GATEWAY_PROJECT_DIR")
		if gatewayProjectDir == "" {
			return nil, errs.Errorf("Please set \"GATEWAY_PROJECT_DIR\" environment variable with the location of your checked out storj/gateway project. (Required to use web resources")
		}
		return standalone.NewStandalone(dir, storjProjectDir, gatewayProjectDir)
	}

	return nil, errors.New("directory doesn't contain supported deployment descriptor")
}
