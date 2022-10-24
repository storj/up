// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"strings"

	"storj.io/storj-up/pkg/recipe"
)

// BuildDict stores the name of the container to build for Storj services.
var BuildDict = map[string]string{
	"authservice":     "app-edge",
	"gateway-mt":      "app-edge",
	"linksharing":     "app-edge",
	"satellite-admin": "app-storj",
	"satellite-api":   "app-storj",
	"satellite-core":  "app-storj",
	"storagenode":     "app-storj",
	"uplink":          "app-storj",
	"versioncontrol":  "app-storj",
}

// ResolveBuilds returns with the required docker images to build (as keys in the maps).
func ResolveBuilds(services []string) (map[string]string, error) {
	result := make(map[string]string)
	resolvedServices, err := ResolveServices(services)
	if err != nil {
		return result, err
	}
	for _, service := range resolvedServices {
		result[BuildDict[service]] = ""
	}
	return result, nil
}

// ResolveServices replaces group definition with exact services in the list.
func ResolveServices(selectors []string) ([]string, error) {
	var res []string
	st, err := recipe.GetStack()
	if err != nil {
		return res, err
	}
	for _, oneOrMoreSelector := range selectors {
		for _, selector := range strings.Split(oneOrMoreSelector, ",") {
			resolved := false
			for _, r := range st {
				if r.Name == selector {
					for _, rs := range r.Add {
						res = append(res, rs.Name)
						resolved = true
					}
				}
			}
			if !resolved {
				res = append(res, selector)
			}
		}
	}
	return res, nil
}
