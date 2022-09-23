// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"embed"
	"io/ioutil"
	"os"
	"path/filepath"

	dockefiles "storj.io/storj-up/cmd/files/docker"
	"storj.io/storj-up/cmd/files/templates"
)

// ConfigFiles is a type with an embedded FS and a folder for directories,
// and a file and filename for single embedded files.
type ConfigFiles struct {
	file     []byte
	filename string
	fs       embed.FS
	folder   string
}

// AddFiles extract the config files associated with the provided service.
func AddFiles(service string) error {
	configFS := ResolveEmbeds(service)
	for _, configDir := range configFS {
		if configDir.fs != (embed.FS{}) {
			err := recurseFileExtract(configDir.fs, configDir.folder)
			if err != nil {
				return err
			}
		}
		if configDir.file != nil {
			err := ExtractFile("", configDir.filename, configDir.file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func recurseFileExtract(fs embed.FS, folder string) error {
	entries, err := fs.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			err = recurseFileExtract(fs, folder+"/"+entry.Name())
			if err != nil {
				return err
			}
		} else {
			fileContent, err := fs.ReadFile(folder + "/" + entry.Name())
			if err != nil {
				return err
			}
			err = ExtractFile(folder, entry.Name(), fileContent)
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
	case "app-edge":
		return []ConfigFiles{{
			file:     dockefiles.EdgeDocker,
			filename: "edge.Dockerfile",
		}}
	case "app-storj":
		return []ConfigFiles{{
			file:     dockefiles.StorjDocker,
			filename: "storj.Dockerfile",
		}}
	case "app-storjscan":
		return []ConfigFiles{{
			file:     dockefiles.StorjscanDocker,
			filename: "storjscan.Dockerfile",
		}}
	default:
		return nil
	}
}
