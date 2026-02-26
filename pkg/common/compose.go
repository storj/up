// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"context"
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/loader"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/goccy/go-yaml"
	"github.com/zeebo/errs/v2"

	"storj.io/storj-up/pkg/common/composedb"
)

const (
	// ComposeFileName filename used for the docker compose file.
	ComposeFileName = "docker-compose.yaml"
)

// Store is the VersionStore used for compose file history.
var Store = composedb.ComposeHistory{DB: composedb.FileDatabase{}}

// ComposeFile is the simplified structure of one compose file.
type ComposeFile struct {
	Version  string // used for compatibility with Compose V1
	Services types.Services
	Networks types.Networks
}

// LoadComposeFromFile parses docker-compose file from the current directory.
func LoadComposeFromFile(dir string, filename string) (*types.Project, error) {
	options := cli.ProjectOptions{
		Name:        "storj-up",
		ConfigPaths: []string{filepath.Join(dir, filename)},
	}

	return cli.ProjectFromOptions(context.Background(), &options)
}

// LoadComposeFromBytes loads docker-compose definition from bytes.
func LoadComposeFromBytes(composeBytes []byte) (*types.Project, error) {
	return loader.LoadWithContext(context.Background(), types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Content: composeBytes,
			},
		},
		WorkingDir: ".",
	})
}

// CreateBind can create a new volume binding object.
func CreateBind(source string, target string) types.ServiceVolumeConfig {
	return types.ServiceVolumeConfig{
		Type:        "bind",
		Source:      source,
		Target:      target,
		ReadOnly:    false,
		Consistency: "",
		Bind: &types.ServiceVolumeBind{
			Propagation: "",
		},
	}
}

// WriteComposeFile persists current docker-compose project to docker-compose.yaml.
func WriteComposeFile(dir string, compose *types.Project) error {
	prevCompose, _ := LoadComposeFromFile(dir, ComposeFileName)
	err := WriteComposeFileNoHistory(dir, compose)
	if err != nil {
		return err
	}
	if os.Getenv("STORJUP_NO_HISTORY") != "" {
		return nil
	}
	if prevCompose != nil {
		prevComposeBytes, err := yaml.Marshal(prevCompose)
		if err != nil {
			return err
		}
		_, err = Store.SaveCurrentVersion(prevComposeBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteComposeFileNoHistory persists current docker-compose project to docker-compose.yaml without saving a record
// of the current compose file.
func WriteComposeFileNoHistory(dir string, compose *types.Project) error {
	resolvedServices, err := yaml.Marshal(&ComposeFile{Version: "3.4", Services: compose.Services, Networks: compose.Networks})
	if err != nil {
		return errs.Wrap(err)
	}
	if err = os.WriteFile(filepath.Join(dir, ComposeFileName), resolvedServices, 0o644); err != nil {
		return err
	}
	return nil
}
