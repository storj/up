// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package composedb

import "time"

// Version is made up of the ID and the last time the record was updated.
type Version struct {
	// todo use VersionID instead of string
	ID          string
	LastUpdated time.Time
}

// VersionStore is the interface for operations performed on the docker-compose yaml file history.
type VersionStore interface {
	// todo use VersionID instead of string
	SaveCurrentVersion(bytes []byte) (string, error)
	RestoreLatestVersion() ([]byte, error)
	// todo use VersionID instead of string
	RestoreVersion(string) ([]byte, error)
	ListVersions() ([]Version, error)
}

// Database is the interface for a backing database used to store docker-compose yaml file history.
type Database interface {
	Read(versionID string) ([]byte, error)
	Write(versionID string, data []byte) error
	Delete(versionID string) error
	GetObjectVersions() ([]Version, error)
}
