// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"os"
	"path/filepath"
)

// ExtractFile extract embedded file, if doesn't exist.
func ExtractFile(path, fileName string, content []byte) error {
	newpath := filepath.Join(".", path)
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		return err
	}
	if _, err := os.Stat(newpath + "/" + fileName); os.IsNotExist(err) {
		return os.WriteFile(newpath+"/"+fileName, content, 0o644)
	}
	return nil
}
