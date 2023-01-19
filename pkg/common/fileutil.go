// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"fmt"
	"os"
)

// IsRegularFile checks to see there is a regular file at the given path and
// returns an error otherwise.
func IsRegularFile(path string) error {
	fi, err := os.Lstat(path)
	switch {
	case err == nil:
		if typ := fi.Mode() & os.ModeType; typ != 0 {
			return fmt.Errorf("%s is not a regular file", path)
		}
		return nil
	case os.IsNotExist(err):
		return fmt.Errorf("%s is not found", path)
	default:
		return fmt.Errorf("failed to stat %s: %w", path, err)
	}
}
