package common

type Key uint

var ServiceDict = map[string]uint {
	"authservice"      :1,
	"cockroach"        :2,
	"gateway-mt"       :4,
	"gatewaymt"        :4,
	"linksharing"      :8,
	"redis"            :16,
	"satellite-admin"  :32,
	"satelliteadmin"   :32,
	"satellite-api"    :64,
	"satelliteapi"     :64,
	"satellite-core"   :128,
	"satellitecore"    :128,
	"storagenode"      :256,
	"uplink"           :512,
	"versioncontrol"   :1024,
	"prometheus"       :2048,
	"grafana"          :4096,
	"minimal"          :64 + 256,
	"edge"             :1 + 4 + 8,
	"db"               :2 + 16,
	"monitor"          :2048 + 4096,
	"core"             :32 + 64 + 128 + 256 + 1024,
	"storj"            :1 + 4 + 8 + 32 + 64 + 128 + 256 + 512 + 1024,
}

const (
	authservice Key = iota 		// 1
	cockroach				 	// 2
	gatewaymt					// 4
	linksharing					// 8
	redis						// 16
	satelliteadmin				// 32
	satelliteapi				// 64
	satellitecore				// 128
	storagenode					// 256
	uplink						// 512
	versioncontrol				// 1024
	prometheus                  // 2048
	grafana                     // 4096
)

//go:generate stringer -type=Key
