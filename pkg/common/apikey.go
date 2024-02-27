// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/zeebo/errs/v2"

	"storj.io/common/macaroon"
	"storj.io/storj/satellite/console/consolewasm"
)

// Satellite0Identity is a standard test identity included in our compose files.
var Satellite0Identity = "12whfK1EDvHJtajBiAUeajQLYcWqxcQmdYQU5zX5cCf6bAxfgu4"

// GetTestAPIKey can calculate an access grant for the predefined test users/project.
func GetTestAPIKey(satelliteID string) (string, error) {
	key, err := macaroon.FromParts([]byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3,
	}, []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4,
	})
	if err != nil {
		return "", errs.Wrap(err)
	}

	idHash := sha256.Sum256([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	base64Salt := base64.StdEncoding.EncodeToString(idHash[:])

	accessGrant, err := consolewasm.GenAccessGrant(satelliteID, key.Serialize(), "password", base64Salt)
	if err != nil {
		return "", errs.Wrap(err)
	}

	return accessGrant, nil
}
