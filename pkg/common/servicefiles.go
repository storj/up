// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"

	"storj.io/storj-up/cmd/files/templates"
)

// ConfigFiles is a type with an embedded FS and a folder.
type ConfigFiles struct {
	fs     embed.FS
	folder string
}

// AddFiles extract the config files associated with the provided service.
func AddFiles(service string) error {
	configFS := ResolveEmbeds(service)
	for _, configDir := range configFS {
		err := recurseFileExtract(configDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func recurseFileExtract(configFiles ConfigFiles) error {
	entries, err := configFiles.fs.ReadDir(configFiles.folder)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			err = recurseFileExtract(ConfigFiles{configFiles.fs, configFiles.folder + "/" + entry.Name()})
			if err != nil {
				return err
			}
		} else {
			fileContent, err := configFiles.fs.ReadFile(configFiles.folder + "/" + entry.Name())
			if err != nil {
				return err
			}
			err = ExtractFile(configFiles.folder, entry.Name(), fileContent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ExtractFile extract embedded file, if doesn't exist.
func ExtractFile(path, fileName string, content []byte) error {
	newpath := filepath.Join(".", path)
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		return err
	}
	if _, err := os.Stat(newpath + "/" + fileName); os.IsNotExist(err) {
		return ioutil.WriteFile(newpath+"/"+fileName, content, 0o644)
	}
	return nil
}

// ResolveEmbeds returns the embedded file system and folder name of the provided service.
func ResolveEmbeds(service string) []ConfigFiles {
	switch service {
	case "geth":
		return []ConfigFiles{
			{
				fs:     templates.GethData,
				folder: "geth/geth-config",
			},
			{
				fs:     templates.StorjscanData,
				folder: "storjscan/test-contract",
			},
		}
	case "prometheus":
		return []ConfigFiles{{
			fs:     templates.PrometheusYaml,
			folder: "prometheus",
		}}
	default:
		return nil
	}
}
