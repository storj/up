// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"os"
	"path/filepath"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
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
}

// LoadComposeFromFile parses docker-compose file from the current directory.
func LoadComposeFromFile(dir string, filename string) (*types.Project, error) {
	options := cli.ProjectOptions{
		Name:        filename,
		ConfigPaths: []string{filepath.Join(dir, filename)},
	}

	return cli.ProjectFromOptions(&options)
}

// LoadComposeFromBytes loads docker-compose definition from bytes.
func LoadComposeFromBytes(composeBytes []byte) (*types.Project, error) {
	return loader.Load(types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Content: composeBytes,
			},
		},
		WorkingDir: ".",
	})
}

// ContainsService check if the service is included in the list.
func ContainsService(s []types.ServiceConfig, e string) bool {
	for _, a := range s {
		if a.Name == e {
			return true
		}
	}
	return false
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
	resolvedServices, err := yaml.Marshal(&ComposeFile{Version: "3.4", Services: compose.Services})
	if err != nil {
		return errs.Wrap(err)
	}
	if err = os.WriteFile(filepath.Join(dir, ComposeFileName), resolvedServices, 0o644); err != nil {
		return err
	}
	return nil
}
