// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

func init() {
	Config["storagenode"] = []Option{
		{
			Name:        "STORJ_IDENTITY_CERT_PATH",
			Description: "path to the certificate chain for this identity",
			Default:     "$IDENTITYDIR/identity.cert",
		}, {
			Name:        "STORJ_IDENTITY_KEY_PATH",
			Description: "path to the private key for this identity",
			Default:     "$IDENTITYDIR/identity.key",
		}, {
			Name:        "STORJ_SERVER_CONFIG_REVOCATION_DBURL",
			Description: "url for revocation database (e.g. bolt://some.db OR redis://127.0.0.1:6378?db=2&password=abc123)",
			Default:     "bolt://$CONFDIR/revocations.db",
		}, {
			Name:        "STORJ_SERVER_CONFIG_PEER_CAWHITELIST_PATH",
			Description: "path to the CA cert whitelist (peer identities must be signed by one these to be verified). this will override the default peer whitelist",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_CONFIG_USE_PEER_CAWHITELIST",
			Description: "if true, uses peer ca whitelist checking",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_CONFIG_PEER_IDVERSIONS",
			Description: "identity version(s) the server will be allowed to talk to",
			Default:     "latest",
		}, {
			Name:        "STORJ_SERVER_CONFIG_EXTENSIONS_REVOCATION",
			Description: "if true, client leaves may contain the most recent certificate revocation for the current certificate",
			Default:     "true",
		}, {
			Name:        "STORJ_SERVER_CONFIG_EXTENSIONS_WHITELIST_SIGNED_LEAF",
			Description: "if true, client leaves must contain a valid \"signed certificate extension\" (NB: verified against certs in the peer ca whitelist; i.e. if true, a whitelist must be provided)",
			Default:     "false",
		}, {
			Name:        "STORJ_SERVER_ADDRESS",
			Description: "public address to listen on",
			Default:     ":7777",
		}, {
			Name:        "STORJ_SERVER_PRIVATE_ADDRESS",
			Description: "private address to listen on",
			Default:     "127.0.0.1:7778",
		}, {
			Name:        "STORJ_SERVER_DISABLE_QUIC",
			Description: "disable QUIC listener on a server",
			Default:     "false",
		}, {
			Name:        "STORJ_SERVER_DISABLE_TCPTLS",
			Description: "disable TCP/TLS listener on a server",
			Default:     "",
		}, {
			Name:        "STORJ_SERVER_DEBUG_LOG_TRAFFIC",
			Description: "",
			Default:     "false",
		}, {
			Name:        "STORJ_DEBUG_ADDRESS",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_DEBUG_CONTROL_TITLE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_DEBUG_CONTROL",
			Description: "expose control panel",
			Default:     "",
		}, {
			Name:        "STORJ_PREFLIGHT_LOCAL_TIME_CHECK",
			Description: "whether or not preflight check for local system clock is enabled on the satellite side. When disabling this feature, your storagenode may not setup correctly.",
			Default:     "true",
		}, {
			Name:        "STORJ_PREFLIGHT_DATABASE_CHECK",
			Description: "whether or not preflight check for database is enabled.",
			Default:     "true",
		}, {
			Name:        "STORJ_CONTACT_EXTERNAL_ADDRESS",
			Description: "the public address of the node, useful for nodes behind NAT",
			Default:     "",
		}, {
			Name:        "STORJ_CONTACT_INTERVAL",
			Description: "how frequently the node contact chore should run",
			Default:     "",
		}, {
			Name:        "STORJ_OPERATOR_EMAIL",
			Description: "operator email address",
			Default:     "",
		}, {
			Name:        "STORJ_OPERATOR_WALLET",
			Description: "operator wallet address",
			Default:     "",
		}, {
			Name:        "STORJ_OPERATOR_WALLET_FEATURES",
			Description: "operator wallet features",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE_PATH",
			Description: "path to store data in",
			Default:     "$CONFDIR/storage",
		}, {
			Name:        "STORJ_STORAGE_WHITELISTED_SATELLITES",
			Description: "a comma-separated list of approved satellite node urls (unused)",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE_ALLOCATED_DISK_SPACE",
			Description: "total allocated disk space in bytes",
			Default:     "1TB",
		}, {
			Name:        "STORJ_STORAGE_ALLOCATED_BANDWIDTH",
			Description: "total allocated bandwidth in bytes (deprecated)",
			Default:     "0B",
		}, {
			Name:        "STORJ_STORAGE_KBUCKET_REFRESH_INTERVAL",
			Description: "how frequently Kademlia bucket should be refreshed with node stats",
			Default:     "1h0m0s",
		}, {
			Name:        "STORJ_STORAGE2_DATABASE_DIR",
			Description: "directory to store databases. if empty, uses data path",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_EXPIRATION_GRACE_PERIOD",
			Description: "how soon before expiration date should things be considered expired",
			Default:     "48h0m0s",
		}, {
			Name:        "STORJ_STORAGE2_MAX_CONCURRENT_REQUESTS",
			Description: "how many concurrent requests are allowed, before uploads are rejected. 0 represents unlimited.",
			Default:     "0",
		}, {
			Name:        "STORJ_STORAGE2_DELETE_WORKERS",
			Description: "how many piece delete workers",
			Default:     "1",
		}, {
			Name:        "STORJ_STORAGE2_DELETE_QUEUE_SIZE",
			Description: "size of the piece delete queue",
			Default:     "10000",
		}, {
			Name:        "STORJ_STORAGE2_ORDER_LIMIT_GRACE_PERIOD",
			Description: "how long after OrderLimit creation date are OrderLimits no longer accepted",
			Default:     "1h0m0s",
		}, {
			Name:        "STORJ_STORAGE2_CACHE_SYNC_INTERVAL",
			Description: "how often the space used cache is synced to persistent storage",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_STREAM_OPERATION_TIMEOUT",
			Description: "how long to spend waiting for a stream operation before canceling",
			Default:     "30m",
		}, {
			Name:        "STORJ_STORAGE2_RETAIN_TIME_BUFFER",
			Description: "allows for small differences in the satellite and storagenode clocks",
			Default:     "48h0m0s",
		}, {
			Name:        "STORJ_STORAGE2_REPORT_CAPACITY_THRESHOLD",
			Description: "threshold below which to immediately notify satellite of capacity",
			Default:     "500MB",
		}, {
			Name:        "STORJ_STORAGE2_MAX_USED_SERIALS_SIZE",
			Description: "amount of memory allowed for used serials store - once surpassed, serials will be dropped at random",
			Default:     "1MB",
		}, {
			Name:        "STORJ_STORAGE2_MIN_UPLOAD_SPEED",
			Description: "a client upload speed should not be lower than MinUploadSpeed in bytes-per-second (E.g: 1Mb), otherwise, it will be flagged as slow-connection and potentially be closed",
			Default:     "0Mb",
		}, {
			Name:        "STORJ_STORAGE2_MIN_UPLOAD_SPEED_GRACE_DURATION",
			Description: "if MinUploadSpeed is configured, after a period of time after the client initiated the upload, the server will flag unusually slow upload client",
			Default:     "0h0m10s",
		}, {
			Name:        "STORJ_STORAGE2_MIN_UPLOAD_SPEED_CONGESTION_THRESHOLD",
			Description: "if the portion defined by the total number of alive connection per MaxConcurrentRequest reaches this threshold, a slow upload client will no longer be monitored and flagged",
			Default:     "0.8",
		}, {
			Name:        "STORJ_STORAGE2_TRUST_SOURCES",
			Description: "list of trust sources",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_TRUST_EXCLUSIONS_RULES",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_TRUST_REFRESH_INTERVAL",
			Description: "how often the trust pool should be refreshed",
			Default:     "6h",
		}, {
			Name:        "STORJ_STORAGE2_TRUST_CACHE_PATH",
			Description: "file path where trust lists should be cached",
			Default:     "${CONFDIR}/trust-cache.json",
		}, {
			Name:        "STORJ_STORAGE2_MONITOR_INTERVAL",
			Description: "how frequently Kademlia bucket should be refreshed with node stats",
			Default:     "1h0m0s",
		}, {
			Name:        "STORJ_STORAGE2_MONITOR_VERIFY_DIR_READABLE_INTERVAL",
			Description: "how frequently to verify the location and readability of the storage directory",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_MONITOR_VERIFY_DIR_WRITABLE_INTERVAL",
			Description: "how frequently to verify writability of storage directory",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_MONITOR_MINIMUM_DISK_SPACE",
			Description: "how much disk space a node at minimum has to advertise",
			Default:     "500GB",
		}, {
			Name:        "STORJ_STORAGE2_MONITOR_MINIMUM_BANDWIDTH",
			Description: "how much bandwidth a node at minimum has to advertise (deprecated)",
			Default:     "0TB",
		}, {
			Name:        "STORJ_STORAGE2_MONITOR_NOTIFY_LOW_DISK_COOLDOWN",
			Description: "minimum length of time between capacity reports",
			Default:     "10m",
		}, {
			Name:        "STORJ_STORAGE2_ORDERS_MAX_SLEEP",
			Description: "maximum duration to wait before trying to send orders",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_ORDERS_SENDER_INTERVAL",
			Description: "duration between sending",
			Default:     "",
		}, {
			Name:        "STORJ_STORAGE2_ORDERS_SENDER_TIMEOUT",
			Description: "timeout for sending",
			Default:     "1h0m0s",
		}, {
			Name:        "STORJ_STORAGE2_ORDERS_SENDER_DIAL_TIMEOUT",
			Description: "timeout for dialing satellite during sending orders",
			Default:     "1m0s",
		}, {
			Name:        "STORJ_STORAGE2_ORDERS_CLEANUP_INTERVAL",
			Description: "duration between archive cleanups",
			Default:     "5m0s",
		}, {
			Name:        "STORJ_STORAGE2_ORDERS_ARCHIVE_TTL",
			Description: "length of time to archive orders before deletion",
			Default:     "168h0m0s",
		}, {
			Name:        "STORJ_STORAGE2_ORDERS_PATH",
			Description: "path to store order limit files in",
			Default:     "$CONFDIR/orders",
		}, {
			Name:        "STORJ_COLLECTOR_INTERVAL",
			Description: "how frequently expired pieces are collected",
			Default:     "1h0m0s",
		}, {
			Name:        "STORJ_FILESTORE_WRITE_BUFFER_SIZE",
			Description: "in-memory buffer for uploads",
			Default:     "128KiB",
		}, {
			Name:        "STORJ_PIECES_WRITE_PREALLOC_SIZE",
			Description: "file preallocated for uploading",
			Default:     "4MiB",
		}, {
			Name:        "STORJ_PIECES_DELETE_TO_TRASH",
			Description: "move pieces to trash upon deletion. Warning: if set to false, you risk disqualification for failed audits if a satellite database is restored from backup.",
			Default:     "true",
		}, {
			Name:        "STORJ_RETAIN_MAX_TIME_SKEW",
			Description: "allows for small differences in the satellite and storagenode clocks",
			Default:     "72h0m0s",
		}, {
			Name:        "STORJ_RETAIN_STATUS",
			Description: "allows configuration to enable, disable, or test retain requests from the satellite. Options: (disabled/enabled/debug)",
			Default:     "enabled",
		}, {
			Name:        "STORJ_RETAIN_CONCURRENCY",
			Description: "how many concurrent retain requests can be processed at the same time.",
			Default:     "5",
		}, {
			Name:        "STORJ_NODESTATS_MAX_SLEEP",
			Description: "maximum duration to wait before requesting data",
			Default:     "",
		}, {
			Name:        "STORJ_NODESTATS_REPUTATION_SYNC",
			Description: "how often to sync reputation",
			Default:     "",
		}, {
			Name:        "STORJ_NODESTATS_STORAGE_SYNC",
			Description: "how often to sync storage",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_ADDRESS",
			Description: "server address of the api gateway and frontend app",
			Default:     "127.0.0.1:14002",
		}, {
			Name:        "STORJ_CONSOLE_STATIC_DIR",
			Description: "path to static resources",
			Default:     "",
		}, {
			Name:        "STORJ_VERSION_CLIENT_CONFIG_SERVER_ADDRESS",
			Description: "server address to check its version against",
			Default:     "https://version.storj.io",
		}, {
			Name:        "STORJ_VERSION_CLIENT_CONFIG_REQUEST_TIMEOUT",
			Description: "Request timeout for version checks",
			Default:     "0h1m0s",
		}, {
			Name:        "STORJ_VERSION_CHECK_INTERVAL",
			Description: "Interval to check the version",
			Default:     "0h15m0s",
		}, {
			Name:        "STORJ_BANDWIDTH_INTERVAL",
			Description: "how frequently bandwidth usage rollups are calculated",
			Default:     "1h0m0s",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_CHORE_INTERVAL",
			Description: "how often to run the chore to check for satellites for the node to exit.",
			Default:     "",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_NUM_WORKERS",
			Description: "number of workers to handle satellite exits",
			Default:     "4",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_NUM_CONCURRENT_TRANSFERS",
			Description: "number of concurrent transfers per graceful exit worker",
			Default:     "5",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_MIN_BYTES_PER_SECOND",
			Description: "the minimum acceptable bytes that an exiting node can transfer per second to the new node",
			Default:     "5KB",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_MIN_DOWNLOAD_TIMEOUT",
			Description: "the minimum duration for downloading a piece from storage nodes before timing out",
			Default:     "2m",
		},
	}
}
