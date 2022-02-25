// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package composedb

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	composeHistoryRelativePath = "./.history"
	composeHistoryFileEXT      = ".yaml"
)

// FileDatabase implements is an abstraction of ioutil.WriterFile.
type FileDatabase struct{}

// Write implements the Writer interface for a flat filesystem database.
func (db FileDatabase) Write(filename string, data []byte) error {
	err := createDBDirIfNotExist()
	if err != nil {
		return err
	}
	path, err := filepath.Abs(composeHistoryRelativePath)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(path, filename)+composeHistoryFileEXT, data, 0o644)
}

// Read implements the Reader interface for a flat filesystem database.
func (db FileDatabase) Read(filename string) ([]byte, error) {
	err := createDBDirIfNotExist()
	if err != nil {
		return nil, err
	}
	path, err := filepath.Abs(composeHistoryRelativePath)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(filepath.Join(path, filename) + composeHistoryFileEXT)
}

// Delete implements the Delete interface for a flat filesystem database.
func (db FileDatabase) Delete(filename string) error {
	err := createDBDirIfNotExist()
	if err != nil {
		return err
	}
	path, err := filepath.Abs(composeHistoryRelativePath)
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(path, filename) + composeHistoryFileEXT)
}

// GetObjectVersions returns the name and modified time of all objects stored in the DB.
func (db FileDatabase) GetObjectVersions() ([]Version, error) {
	err := createDBDirIfNotExist()
	if err != nil {
		return nil, err
	}
	path, err := filepath.Abs(composeHistoryRelativePath)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	versions := make([]Version, 0, len(files))
	for _, file := range files {
		filename := file.Name()
		lastUpdated := file.ModTime()
		ext := filepath.Ext(filename)
		filename = filename[:len(filename)-len(ext)]
		versions = append(versions, Version{
			ID:          filename,
			LastUpdated: lastUpdated,
		})
	}
	return versions, err
}

// createDBDirIfNotExist creates the path to the stored files if not created.
func createDBDirIfNotExist() error {
	path, err := filepath.Abs(composeHistoryRelativePath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
