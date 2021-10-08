package common

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

