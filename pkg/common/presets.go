// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import "fmt"

// Key is an identfier of one service or service group.
type Key uint

// ServiceDict maps each service and groups to bitmap.
var ServiceDict = map[string]uint{
	"authservice":     1,
	"cockroach":       2,
	"gateway-mt":      4,
	"gatewaymt":       4,
	"linksharing":     8,
	"redis":           16,
	"satellite-admin": 32,
	"satelliteadmin":  32,
	"satellite-api":   64,
	"satelliteapi":    64,
	"satellite-core":  128,
	"satellitecore":   128,
	"storagenode":     256,
	"uplink":          512,
	"versioncontrol":  1024,
	"storjscan":       2048,
	"geth":            4096,
	"prometheus":      8192,
	"grafana":         16384,
	"app-base-dev":    32768,
	"app-base-ubuntu": 65536,
	"app-edge":        131072,
	"app-storj":       262144,
	"minimal":         64 + 256,
	"edge":            1 + 4 + 8,
	"db":              2 + 16,
	"billing":         2048 + 4096,
	"monitor":         8192 + 16384,
	"core":            32 + 64 + 128 + 256 + 1024,
	"storj":           1 + 4 + 8 + 32 + 64 + 128 + 256 + 512 + 1024,
}

// BinaryDict contains the name of executable binaries for Storj service.
var BinaryDict = map[string]string{
	"authservice":     "authservice",
	"gateway-mt":      "gateway-mt",
	"linksharing":     "linksharing",
	"satellite-admin": "satellite",
	"satellite-api":   "satellite",
	"satellite-core":  "satellite",
	"storagenode":     "storagenode",
	"uplink":          "uplink",
	"versioncontrol":  "versioncontrol",
	"storjscan":       "storjscan",
}

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

var serviceNameHelper = map[string]string{
	"authservice":    "authservice",
	"cockroach":      "cockroach",
	"gatewaymt":      "gateway-mt",
	"linksharing":    "linksharing",
	"redis":          "redis",
	"satelliteadmin": "satellite-admin",
	"satelliteapi":   "satellite-api",
	"satellitecore":  "satellite-core",
	"storagenode":    "storagenode",
	"uplink":         "uplink",
	"versioncontrol": "versioncontrol",
	"storjscan":      "storjscan",
	"geth":           "geth",
	"prometheus":     "prometheus",
	"grafana":        "grafana",
	"appbasedev":     "app-base-dev",
	"appbaseubuntu":  "app-base-ubuntu",
	"appedge":        "app-edge",
	"appstorj":       "app-storj",
}

const (
	authservice    Key = iota // 1
	cockroach                 // 2
	gatewaymt                 // 4
	linksharing               // 8
	redis                     // 16
	satelliteadmin            // 32
	satelliteapi              // 64
	satellitecore             // 128
	storagenode               // 256
	uplink                    // 512
	versioncontrol            // 1024
	storjscan                 // 2048
	geth                      // 4096
	prometheus                // 8192
	grafana                   // 16384
	appbasedev                // 32768
	appbaseubuntu             // 65536
	appedge                   // 131072
	appstorj                  // 262144
)

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
func ResolveServices(services []string) ([]string, error) {
	var result []string
	var key uint
	for _, service := range services {
		value, found := ServiceDict[service]
		if !found {
			return nil, fmt.Errorf("invalid service selector %s, please run `storj-up services` to find supported values", service)
		}
		key |= value
	}
	for service := authservice; service <= appstorj; service++ {
		if key&(1<<service) != 0 {
			result = append(result, serviceNameHelper[service.String()])
		}
	}
	return result, nil
}

// GetSelectors returns with selectors and associated services (in case of group definition).
func GetSelectors() map[string][]string {
	valueToName := map[uint]string{}
	for name, value := range ServiceDict {
		valueToName[value] = name
	}

	selectors := map[string][]string{}
	for name, value := range ServiceDict {
		services := []string{}
		for v, n := range valueToName {
			if value&v == v && n != name {
				services = append(services, n)
			}
		}
		selectors[name] = services
	}
	return selectors
}

//go:generate stringer -type=Key
