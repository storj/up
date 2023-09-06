// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import (
	"github.com/zeebo/errs"
)

var basePorts = map[string]int{
	"satellite-api":        10000,
	"satellite-gc":         10100,
	"satellite-bf":         10200,
	"satellite-core":       10300,
	"satellite-admin":      10400,
	"satellite-rangedloop": 10500,
	"satellite-repair":     10600,
	"satellite-audit":      10700,
	"storagenode":          30000,
	"gateway-mt":           20000,
	"authservice":          21000,
	"linksharing":          22200,
}

// PortConvention defines port numbers for any services.
func PortConvention(instance ServiceInstance, portType string) (int, error) {

	port, found := basePorts[instance.Name]
	if !found {
		return 0, errs.New("No base port defined for " + instance.String())
	}
	port += instance.Instance * 10
	switch portType {
	case "console":
		port += 0
	case "public":
		port++
	case "private":
		port += 2
	case "debug":
		port += 9
	}
	return port, nil
}
