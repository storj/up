package common

type Key uint

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
	"prometheus":      2048,
	"grafana":         4096,
	"app-base-dev":    8192,
	"app-base-ubuntu": 16384,
	"app-edge":        32768,
	"app-storj":       65536,
	"traefik":         131072,
	"minimal":         64 + 256,
	"edge":            1 + 4 + 8,
	"db":              2 + 16,
	"monitor":         2048 + 4096,
	"core":            32 + 64 + 128 + 256 + 1024,
	"storj":           1 + 4 + 8 + 32 + 64 + 128 + 256 + 512 + 1024,
}

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
}

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
	"prometheus":     "prometheus",
	"grafana":        "grafana",
	"appbasedev":     "app-base-dev",
	"appbaseubuntu":  "app-base-ubuntu",
	"appedge":        "app-edge",
	"appstorj":       "app-storj",
	"traefik":        "traefik",
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
	prometheus                // 2048
	grafana                   // 4096
	appbasedev                // 8192
	appbaseubuntu             // 16384
	appedge                   // 32768
	appstorj                  // 65536
	traefik                   // 131072
)

func ResolveBuilds(services []string) map[string]string {
	result := make(map[string]string)
	for _, service := range ResolveServices(services) {
		result[BuildDict[service]] = ""
	}
	return result
}

func ResolveServices(services []string) []string {
	var result []string
	var key uint
	for _, service := range services {
		key |= ServiceDict[service]
	}

	for service := authservice; service <= traefik; service++ {
		if key&(1<<service) != 0 {
			result = append(result, serviceNameHelper[service.String()])
		}
	}
	return result
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
