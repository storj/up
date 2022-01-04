// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package composedb

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"

	"github.com/zeebo/errs"
)

// ComposeHistory used to interact with the previous compose files.
type ComposeHistory struct {
	DB Database
}

const (
	maxSize                  = 5
	composeHistoryObjectName = "docker-compose"
)

type byLatest []Version

func (a byLatest) Len() int           { return len(a) }
func (a byLatest) Less(i, j int) bool { return sortableName(a[i].ID) > sortableName(a[j].ID) }
func (a byLatest) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// RestoreLatestVersion returns the latest version of the saved compose history.
func (s ComposeHistory) RestoreLatestVersion() ([]byte, error) {
	if empty(s.DB) {
		return nil, fmt.Errorf("DB is empty, no history to restore")
	}
	objectNames, err := getObjectNamesSortedByLatest(s.DB)
	if err != nil {
		return nil, err
	}
	bytes, err := s.DB.Read(objectNames[0].ID)
	if err != nil {
		return nil, err
	}
	err = s.DB.Delete(objectNames[0].ID)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// RestoreVersion  returns the requested version of the saved compose history.
func (s ComposeHistory) RestoreVersion(version string) ([]byte, error) {
	// todo add functionality for selective file history restore
	return nil, fmt.Errorf("Unsupported")
}

// ListVersions  lists the latest stored versions of the compose history files.
func (s ComposeHistory) ListVersions() ([]Version, error) {
	if empty(s.DB) {
		return nil, fmt.Errorf("DB is empty, no history to list")
	}
	objectNames, err := getObjectNamesSortedByLatest(s.DB)
	if err != nil {
		return nil, err
	}
	return objectNames, nil
}

// SaveCurrentVersion stores the provided compose history into the compose history records.
func (s ComposeHistory) SaveCurrentVersion(bytes []byte) (string, error) {
	objectNames, err := getObjectNamesSortedByLatest(s.DB)
	if err != nil {
		return "", err
	}
	objectName := newObjectName(objectNames)
	if full(s.DB) {
		defer func() { err = errs.Combine(err, s.DB.Delete(objectNames[len(objectNames)-1].ID)) }()
	}
	err = s.DB.Write(objectName, bytes)
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func empty(db Database) bool {
	objectNames, _ := db.GetObjectVersions()
	return len(objectNames) == 0
}

func full(db Database) bool {
	objectNames, _ := db.GetObjectVersions()
	return len(objectNames) >= maxSize
}

// newObjectName creates a new object name for history to be pushed onto the stack. names are
// sequentially increased as needed. Initial object contains the suffix "1".
func newObjectName(objectNames []Version) string {
	if len(objectNames) == 0 {
		return composeHistoryObjectName + "1"
	}
	numericSuffixAsInt, err := strconv.Atoi(getNumericSuffix(objectNames[0].ID))
	if err != nil {
		return ""
	}
	numericSuffixAsInt++
	return composeHistoryObjectName + strconv.Itoa(numericSuffixAsInt)
}

func getNumericSuffix(objectName string) string {
	// split numeric suffix
	i := len(objectName) - 1
	for ; i >= 0; i-- {
		if '0' > objectName[i] || objectName[i] > '9' {
			break
		}
	}
	i++
	return objectName[i:]
}

// getObjectNamesSortedByLatest returns the current stored objects sorted by objectname (latest first).
func getObjectNamesSortedByLatest(db Database) ([]Version, error) {
	objectVersions, err := db.GetObjectVersions()
	sort.Sort(byLatest(objectVersions))
	return objectVersions, err
}

// sortableName returns a objectname sort key with non-negative integer suffixes.
func sortableName(objectName string) string {
	// split numeric suffix
	i := len(objectName) - 1
	for ; i >= 0; i-- {
		if '0' > objectName[i] || objectName[i] > '9' {
			break
		}
	}
	i++
	// string numeric suffix to uint64 bytes
	// empty string is zero, so integers are plus one
	b64 := make([]byte, 64/8)
	s64 := objectName[i:]
	if len(s64) > 0 {
		u64, err := strconv.ParseUint(s64, 10, 64)
		if err == nil {
			binary.BigEndian.PutUint64(b64, u64+1)
		}
	}
	// prefix + numeric-suffix + ext
	return objectName[:i] + string(b64)
}
