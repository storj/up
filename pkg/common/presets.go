package common

import "strings"

var Presets = createPresets()

func createPresets() map[string][]string {
	presets := map[string][]string{}
	presets["minimal"] = []string{"satellite-api", "storagenode"}
	presets["edge"] = []string{"gateway-mt", "linksharing", "authservice"}
	presets["db"] = []string{"cockroach", "redis"}
	presets["monitor"] = []string{"prometheus", "grafana"}
	presets["core"] = append(presets["minimal"], "satellite-core", "satellite-admin", "versioncontrol")
	presets["storj"] = append(presets["core"], presets["edge"]...)
	presets["storj"] = append(presets["storj"], "uplink")
	return presets
}

func Selected(selector string, service string) bool {
	for _, part := range strings.Split(selector, ",") {
		selector := strings.TrimSpace(part)
		if selector == "all" {
			return true
		}
		if selector == service {
			return true
		}
		if group, found := Presets[selector]; found {
			for _, s := range group {
				if s == service {
					return true
				}
			}
		}
	}
	return false
}
