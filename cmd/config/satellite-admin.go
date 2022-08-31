// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package config

func init() {
	Config["satellite-admin"] = []Option{
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
			Name:        "STORJ_ADMIN_ADDRESS",
			Description: "admin peer http listening address",
			Default:     "",
		}, {
			Name:        "STORJ_ADMIN_STATIC_DIR",
			Description: "an alternate directory path which contains the static assets to serve. When empty, it uses the embedded assets",
			Default:     "",
		}, {
			Name:        "STORJ_ADMIN_AUTHORIZATION_TOKEN",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_PASSWORD_COST",
			Description: "password hashing cost (0=automatic)",
			Default:     "0",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_OPEN_REGISTRATION_ENABLED",
			Description: "enable open registration",
			Default:     "false",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_DEFAULT_PROJECT_LIMIT",
			Description: "default project limits for users",
			Default:     "1",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_TOKEN_EXPIRATION_TIME",
			Description: "expiration time for auth tokens, account recovery tokens, and activation tokens",
			Default:     "24h",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_AS_OF_SYSTEM_TIME_DURATION",
			Description: "default duration for AS OF SYSTEM TIME",
			Default:     "",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_USAGE_LIMITS_STORAGE_FREE",
			Description: "the default free-tier storage usage limit",
			Default:     "150.00GB",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_USAGE_LIMITS_STORAGE_PAID",
			Description: "the default paid-tier storage usage limit",
			Default:     "25.00TB",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_USAGE_LIMITS_BANDWIDTH_FREE",
			Description: "the default free-tier bandwidth usage limit",
			Default:     "150.00GB",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_USAGE_LIMITS_BANDWIDTH_PAID",
			Description: "the default paid-tier bandwidth usage limit",
			Default:     "100.00TB",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_USAGE_LIMITS_SEGMENT_FREE",
			Description: "the default free-tier segment usage limit",
			Default:     "150000",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_USAGE_LIMITS_SEGMENT_PAID",
			Description: "the default paid-tier segment usage limit",
			Default:     "1000000",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_RECAPTCHA_ENABLED",
			Description: "whether or not reCAPTCHA is enabled for user registration",
			Default:     "false",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_RECAPTCHA_SITE_KEY",
			Description: "reCAPTCHA site key",
			Default:     "",
		}, {
			Name:        "STORJ_ADMIN_CONSOLE_CONFIG_RECAPTCHA_SECRET_KEY",
			Description: "reCAPTCHA secret key",
			Default:     "",
		}, {
			Name:        "STORJ_CONTACT_EXTERNAL_ADDRESS",
			Description: "the public address of the node, useful for nodes behind NAT",
			Default:     "",
		}, {
			Name:        "STORJ_CONTACT_TIMEOUT",
			Description: "timeout for pinging storage nodes",
			Default:     "10m0s",
		}, {
			Name:        "STORJ_CONTACT_RATE_LIMIT_INTERVAL",
			Description: "the amount of time that should happen between contact attempts usually",
			Default:     "",
		}, {
			Name:        "STORJ_CONTACT_RATE_LIMIT_BURST",
			Description: "the maximum burst size for the contact rate limit token bucket",
			Default:     "",
		}, {
			Name:        "STORJ_CONTACT_RATE_LIMIT_CACHE_SIZE",
			Description: "the number of nodes or addresses to keep token buckets for",
			Default:     "1000",
		}, {
			Name:        "STORJ_OVERLAY_NODE_NEW_NODE_FRACTION",
			Description: "the fraction of new nodes allowed per request",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_NODE_MINIMUM_VERSION",
			Description: "the minimum node software version for node selection queries",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_NODE_ONLINE_WINDOW",
			Description: "the amount of time without seeing a node before its considered offline",
			Default:     "4h",
		}, {
			Name:        "STORJ_OVERLAY_NODE_DISTINCT_IP",
			Description: "require distinct IPs when choosing nodes for upload",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_NODE_MINIMUM_DISK_SPACE",
			Description: "how much disk space a node at minimum must have to be selected for upload",
			Default:     "500.00MB",
		}, {
			Name:        "STORJ_OVERLAY_NODE_AS_OF_SYSTEM_TIME_ENABLED",
			Description: "enables the use of the AS OF SYSTEM TIME feature in CRDB",
			Default:     "true",
		}, {
			Name:        "STORJ_OVERLAY_NODE_AS_OF_SYSTEM_TIME_DEFAULT_INTERVAL",
			Description: "default duration for AS OF SYSTEM TIME",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_NODE_UPLOAD_EXCLUDED_COUNTRY_CODES",
			Description: "list of country codes to exclude from node selection for uploads",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_NODE_SELECTION_CACHE_DISABLED",
			Description: "disable node cache",
			Default:     "false",
		}, {
			Name:        "STORJ_OVERLAY_NODE_SELECTION_CACHE_STALENESS",
			Description: "how stale the node selection cache can be",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_GEO_IP_DB",
			Description: "the location of the maxmind database containing geoip country information",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_GEO_IP_MOCK_COUNTRIES",
			Description: "a mock list of countries the satellite will attribute to nodes (useful for testing)",
			Default:     "",
		}, {
			Name:        "STORJ_OVERLAY_UPDATE_STATS_BATCH_SIZE",
			Description: "number of update requests to process per transaction",
			Default:     "100",
		}, {
			Name:        "STORJ_OVERLAY_NODE_CHECK_IN_WAIT_PERIOD",
			Description: "the amount of time to wait before accepting a redundant check-in from a node (unmodified info since last check-in)",
			Default:     "2h",
		}, {
			Name:        "STORJ_OVERLAY_REPAIR_EXCLUDED_COUNTRY_CODES",
			Description: "list of country codes to exclude nodes from target repair selection",
			Default:     "",
		}, {
			Name:        "STORJ_STRAY_NODES_ENABLE_DQ",
			Description: "whether nodes will be disqualified if they have not been contacted in some time",
			Default:     "",
		}, {
			Name:        "STORJ_STRAY_NODES_INTERVAL",
			Description: "how often to check for and DQ stray nodes",
			Default:     "",
		}, {
			Name:        "STORJ_STRAY_NODES_MAX_DURATION_WITHOUT_CONTACT",
			Description: "length of time a node can go without contacting satellite before being disqualified",
			Default:     "",
		}, {
			Name:        "STORJ_STRAY_NODES_LIMIT",
			Description: "Max number of nodes to return in a single query. Chore will iterate until rows returned is less than limit",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_DATABASE_URL",
			Description: "the database connection string to use",
			Default:     "postgres://",
		}, {
			Name:        "STORJ_METAINFO_MIN_REMOTE_SEGMENT_SIZE",
			Description: "minimum remote segment size",
			Default:     "1240",
		}, {
			Name:        "STORJ_METAINFO_MAX_INLINE_SEGMENT_SIZE",
			Description: "maximum inline segment size",
			Default:     "4KiB",
		}, {
			Name:        "STORJ_METAINFO_MAX_ENCRYPTED_OBJECT_KEY_LENGTH",
			Description: "maximum encrypted object key length",
			Default:     "1280",
		}, {
			Name:        "STORJ_METAINFO_MAX_SEGMENT_SIZE",
			Description: "maximum segment size",
			Default:     "64MiB",
		}, {
			Name:        "STORJ_METAINFO_MAX_METADATA_SIZE",
			Description: "maximum segment metadata size",
			Default:     "2KiB",
		}, {
			Name:        "STORJ_METAINFO_MAX_COMMIT_INTERVAL",
			Description: "maximum time allowed to pass between creating and committing a segment",
			Default:     "48h",
		}, {
			Name:        "STORJ_METAINFO_MIN_PART_SIZE",
			Description: "minimum allowed part size (last part has no minimum size limit)",
			Default:     "5MiB",
		}, {
			Name:        "STORJ_METAINFO_MAX_NUMBER_OF_PARTS",
			Description: "maximum number of parts object can contain",
			Default:     "10000",
		}, {
			Name:        "STORJ_METAINFO_OVERLAY",
			Description: "toggle flag if overlay is enabled",
			Default:     "true",
		}, {
			Name:        "STORJ_METAINFO_RS_ERASURE_SHARE_SIZE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_RS_MIN",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_RS_REPAIR",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_RS_SUCCESS",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_RS_TOTAL",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_SEGMENT_LOOP_COALESCE_DURATION",
			Description: "how long to wait for new observers before starting iteration",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_SEGMENT_LOOP_RATE_LIMIT",
			Description: "rate limit (default is 0 which is unlimited segments per second)",
			Default:     "0",
		}, {
			Name:        "STORJ_METAINFO_SEGMENT_LOOP_LIST_LIMIT",
			Description: "how many items to query in a batch",
			Default:     "2500",
		}, {
			Name:        "STORJ_METAINFO_SEGMENT_LOOP_AS_OF_SYSTEM_INTERVAL",
			Description: "as of system interval",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_SEGMENT_LOOP_SUSPICIOUS_PROCESSED_RATIO",
			Description: "ratio where to consider processed count as supicious",
			Default:     "0.03",
		}, {
			Name:        "STORJ_METAINFO_RATE_LIMITER_ENABLED",
			Description: "whether rate limiting is enabled.",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_RATE_LIMITER_RATE",
			Description: "request rate per project per second.",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_RATE_LIMITER_CACHE_CAPACITY",
			Description: "number of projects to cache.",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_RATE_LIMITER_CACHE_EXPIRATION",
			Description: "how long to cache the projects limiter.",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_PROJECT_LIMITS_MAX_BUCKETS",
			Description: "max bucket count for a project.",
			Default:     "100",
		}, {
			Name:        "STORJ_METAINFO_PROJECT_LIMITS_VALIDATE_SEGMENT_LIMIT",
			Description: "whether segment limit validation is enabled.",
			Default:     "false",
		}, {
			Name:        "STORJ_METAINFO_PIECE_DELETION_MAX_CONCURRENCY",
			Description: "maximum number of concurrent requests to storage nodes",
			Default:     "100",
		}, {
			Name:        "STORJ_METAINFO_PIECE_DELETION_MAX_CONCURRENT_PIECES",
			Description: "maximum number of concurrent pieces can be processed",
			Default:     "1000000",
		}, {
			Name:        "STORJ_METAINFO_PIECE_DELETION_MAX_PIECES_PER_BATCH",
			Description: "maximum number of pieces per batch",
			Default:     "5000",
		}, {
			Name:        "STORJ_METAINFO_PIECE_DELETION_MAX_PIECES_PER_REQUEST",
			Description: "maximum number pieces per single request",
			Default:     "1000",
		}, {
			Name:        "STORJ_METAINFO_PIECE_DELETION_DIAL_TIMEOUT",
			Description: "timeout for dialing nodes (0 means satellite default)",
			Default:     "3s",
		}, {
			Name:        "STORJ_METAINFO_PIECE_DELETION_FAIL_THRESHOLD",
			Description: "threshold for retrying a failed node",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_PIECE_DELETION_REQUEST_TIMEOUT",
			Description: "timeout for a single delete request",
			Default:     "",
		}, {
			Name:        "STORJ_METAINFO_SERVER_SIDE_COPY",
			Description: "enable code for server-side copy",
			Default:     "",
		}, {
			Name:        "STORJ_ORDERS_ENCRYPTION_KEYS_DEFAULT_ID",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_ORDERS_ENCRYPTION_KEYS_DEFAULT_KEY",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_ORDERS_ENCRYPTION_KEYS_LIST",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_ORDERS_ENCRYPTION_KEYS_KEY_BY_ID",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_ORDERS_EXPIRATION",
			Description: "how long until an order expires",
			Default:     "48h",
		}, {
			Name:        "STORJ_ORDERS_FLUSH_BATCH_SIZE",
			Description: "how many items in the rollups write cache before they are flushed to the database",
			Default:     "",
		}, {
			Name:        "STORJ_ORDERS_FLUSH_INTERVAL",
			Description: "how often to flush the rollups write cache to the database",
			Default:     "",
		}, {
			Name:        "STORJ_ORDERS_NODE_STATUS_LOGGING",
			Description: "deprecated, log the offline/disqualification status of nodes",
			Default:     "false",
		}, {
			Name:        "STORJ_ORDERS_ORDERS_SEMAPHORE_SIZE",
			Description: "how many concurrent orders to process at once. zero is unlimited",
			Default:     "2",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_REPAIR_WEIGHT",
			Description: "weight to apply to audit reputation for total repair reputation calculation",
			Default:     "1.0",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_UPLINK_WEIGHT",
			Description: "weight to apply to audit reputation for total uplink reputation calculation",
			Default:     "1.0",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_LAMBDA",
			Description: "the forgetting factor used to calculate the audit SNs reputation",
			Default:     "0.95",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_WEIGHT",
			Description: "the normalization weight used to calculate the audit SNs reputation",
			Default:     "1.0",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_DQ",
			Description: "the reputation cut-off for disqualifying SNs based on audit history",
			Default:     "0.6",
		}, {
			Name:        "STORJ_REPUTATION_SUSPENSION_GRACE_PERIOD",
			Description: "the time period that must pass before suspended nodes will be disqualified",
			Default:     "",
		}, {
			Name:        "STORJ_REPUTATION_SUSPENSION_DQENABLED",
			Description: "whether nodes will be disqualified if they have been suspended for longer than the suspended grace period",
			Default:     "",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_COUNT",
			Description: "the number of times a node has been audited to not be considered a New Node",
			Default:     "",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_HISTORY_WINDOW_SIZE",
			Description: "The length of time spanning a single audit window",
			Default:     "",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_HISTORY_TRACKING_PERIOD",
			Description: "The length of time to track audit windows for node suspension and disqualification",
			Default:     "",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_HISTORY_GRACE_PERIOD",
			Description: "The length of time to give suspended SNOs to diagnose and fix issues causing downtime. Afterwards, they will have one tracking period to reach the minimum online score before disqualification",
			Default:     "",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_HISTORY_OFFLINE_THRESHOLD",
			Description: "The point below which a node is punished for offline audits. Determined by calculating the ratio of online/total audits within each window and finding the average across windows within the tracking period.",
			Default:     "0.6",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_HISTORY_OFFLINE_DQENABLED",
			Description: "whether nodes will be disqualified if they have low online score after a review period",
			Default:     "",
		}, {
			Name:        "STORJ_REPUTATION_AUDIT_HISTORY_OFFLINE_SUSPENSION_ENABLED",
			Description: "whether nodes will be suspended if they have low online score",
			Default:     "",
		}, {
			Name:        "STORJ_CHECKER_INTERVAL",
			Description: "how frequently checker should check for bad segments",
			Default:     "",
		}, {
			Name:        "STORJ_CHECKER_RELIABILITY_CACHE_STALENESS",
			Description: "how stale reliable node cache can be",
			Default:     "",
		}, {
			Name:        "STORJ_CHECKER_REPAIR_OVERRIDES_LIST",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_CHECKER_NODE_FAILURE_RATE",
			Description: "the probability of a single node going down within the next checker iteration",
			Default:     "0.00005435",
		}, {
			Name:        "STORJ_REPAIRER_MAX_REPAIR",
			Description: "maximum segments that can be repaired concurrently",
			Default:     "",
		}, {
			Name:        "STORJ_REPAIRER_INTERVAL",
			Description: "how frequently repairer should try and repair more data",
			Default:     "",
		}, {
			Name:        "STORJ_REPAIRER_TIMEOUT",
			Description: "time limit for uploading repaired pieces to new storage nodes",
			Default:     "5m0s",
		}, {
			Name:        "STORJ_REPAIRER_DOWNLOAD_TIMEOUT",
			Description: "time limit for downloading pieces from a node for repair",
			Default:     "5m0s",
		}, {
			Name:        "STORJ_REPAIRER_TOTAL_TIMEOUT",
			Description: "time limit for an entire repair job, from queue pop to upload completion",
			Default:     "45m",
		}, {
			Name:        "STORJ_REPAIRER_MAX_BUFFER_MEM",
			Description: "maximum buffer memory (in bytes) to be allocated for read buffers",
			Default:     "4.0 MiB",
		}, {
			Name:        "STORJ_REPAIRER_MAX_EXCESS_RATE_OPTIMAL_THRESHOLD",
			Description: "ratio applied to the optimal threshold to calculate the excess of the maximum number of repaired pieces to upload",
			Default:     "0.05",
		}, {
			Name:        "STORJ_REPAIRER_IN_MEMORY_REPAIR",
			Description: "whether to download pieces for repair in memory (true) or download to disk (false)",
			Default:     "false",
		}, {
			Name:        "STORJ_AUDIT_MAX_RETRIES_STAT_DB",
			Description: "max number of times to attempt updating a statdb batch",
			Default:     "3",
		}, {
			Name:        "STORJ_AUDIT_MIN_BYTES_PER_SECOND",
			Description: "the minimum acceptable bytes that storage nodes can transfer per second to the satellite",
			Default:     "128B",
		}, {
			Name:        "STORJ_AUDIT_MIN_DOWNLOAD_TIMEOUT",
			Description: "the minimum duration for downloading a share from storage nodes before timing out",
			Default:     "5m0s",
		}, {
			Name:        "STORJ_AUDIT_MAX_REVERIFY_COUNT",
			Description: "limit above which we consider an audit is failed",
			Default:     "3",
		}, {
			Name:        "STORJ_AUDIT_CHORE_INTERVAL",
			Description: "how often to run the reservoir chore",
			Default:     "",
		}, {
			Name:        "STORJ_AUDIT_QUEUE_INTERVAL",
			Description: "how often to recheck an empty audit queue",
			Default:     "",
		}, {
			Name:        "STORJ_AUDIT_SLOTS",
			Description: "number of reservoir slots allotted for nodes, currently capped at 3",
			Default:     "3",
		}, {
			Name:        "STORJ_AUDIT_WORKER_CONCURRENCY",
			Description: "number of workers to run audits on segments",
			Default:     "2",
		}, {
			Name:        "STORJ_GARBAGE_COLLECTION_INTERVAL",
			Description: "the time between each send of garbage collection filters to storage nodes",
			Default:     "",
		}, {
			Name:        "STORJ_GARBAGE_COLLECTION_ENABLED",
			Description: "set if garbage collection is enabled or not",
			Default:     "",
		}, {
			Name:        "STORJ_GARBAGE_COLLECTION_INITIAL_PIECES",
			Description: "the initial number of pieces expected for a storage node to have, used for creating a filter",
			Default:     "",
		}, {
			Name:        "STORJ_GARBAGE_COLLECTION_FALSE_POSITIVE_RATE",
			Description: "the false positive rate used for creating a garbage collection bloom filter",
			Default:     "",
		}, {
			Name:        "STORJ_GARBAGE_COLLECTION_CONCURRENT_SENDS",
			Description: "the number of nodes to concurrently send garbage collection bloom filters to",
			Default:     "",
		}, {
			Name:        "STORJ_GARBAGE_COLLECTION_RETAIN_SEND_TIMEOUT",
			Description: "the amount of time to allow a node to handle a retain request",
			Default:     "1m",
		}, {
			Name:        "STORJ_EXPIRED_DELETION_INTERVAL",
			Description: "the time between each attempt to go through the db and clean up expired segments",
			Default:     "",
		}, {
			Name:        "STORJ_EXPIRED_DELETION_ENABLED",
			Description: "set if expired segment cleanup is enabled or not",
			Default:     "",
		}, {
			Name:        "STORJ_EXPIRED_DELETION_LIST_LIMIT",
			Description: "how many expired objects to query in a batch",
			Default:     "100",
		}, {
			Name:        "STORJ_ZOMBIE_DELETION_INTERVAL",
			Description: "the time between each attempt to go through the db and clean up zombie objects",
			Default:     "",
		}, {
			Name:        "STORJ_ZOMBIE_DELETION_ENABLED",
			Description: "set if zombie object cleanup is enabled or not",
			Default:     "true",
		}, {
			Name:        "STORJ_ZOMBIE_DELETION_LIST_LIMIT",
			Description: "how many objects to query in a batch",
			Default:     "100",
		}, {
			Name:        "STORJ_ZOMBIE_DELETION_INACTIVE_FOR",
			Description: "after what time object will be deleted if there where no new upload activity",
			Default:     "24h",
		}, {
			Name:        "STORJ_TALLY_INTERVAL",
			Description: "how frequently the tally service should run",
			Default:     "",
		}, {
			Name:        "STORJ_TALLY_SAVE_ROLLUP_BATCH_SIZE",
			Description: "how large of batches SaveRollup should process at a time",
			Default:     "1000",
		}, {
			Name:        "STORJ_TALLY_READ_ROLLUP_BATCH_SIZE",
			Description: "how large of batches GetBandwidthSince should process at a time",
			Default:     "10000",
		}, {
			Name:        "STORJ_TALLY_LIST_LIMIT",
			Description: "how many objects to query in a batch",
			Default:     "2500",
		}, {
			Name:        "STORJ_TALLY_AS_OF_SYSTEM_INTERVAL",
			Description: "as of system interval",
			Default:     "",
		}, {
			Name:        "STORJ_ROLLUP_INTERVAL",
			Description: "how frequently rollup should run",
			Default:     "",
		}, {
			Name:        "STORJ_ROLLUP_DELETE_TALLIES",
			Description: "option for deleting tallies after they are rolled up",
			Default:     "true",
		}, {
			Name:        "STORJ_ROLLUP_ARCHIVE_INTERVAL",
			Description: "how frequently rollup archiver should run",
			Default:     "",
		}, {
			Name:        "STORJ_ROLLUP_ARCHIVE_ARCHIVE_AGE",
			Description: "age at which a rollup is archived",
			Default:     "2160h",
		}, {
			Name:        "STORJ_ROLLUP_ARCHIVE_BATCH_SIZE",
			Description: "number of records to delete per delete execution. Used only for crdb which is slow without limit.",
			Default:     "500",
		}, {
			Name:        "STORJ_ROLLUP_ARCHIVE_ENABLED",
			Description: "whether or not the rollup archive is enabled.",
			Default:     "true",
		}, {
			Name:        "STORJ_LIVE_ACCOUNTING_STORAGE_BACKEND",
			Description: "what to use for storing real-time accounting data",
			Default:     "",
		}, {
			Name:        "STORJ_LIVE_ACCOUNTING_BANDWIDTH_CACHE_TTL",
			Description: "bandwidth cache key time to live",
			Default:     "5m",
		}, {
			Name:        "STORJ_LIVE_ACCOUNTING_AS_OF_SYSTEM_INTERVAL",
			Description: "as of system interval",
			Default:     "-10s",
		}, {
			Name:        "STORJ_PROJECT_BWCLEANUP_INTERVAL",
			Description: "how often to remove unused project bandwidth rollups",
			Default:     "168h",
		}, {
			Name:        "STORJ_PROJECT_BWCLEANUP_RETAIN_MONTHS",
			Description: "number of months of project bandwidth rollups to retain, not including the current month",
			Default:     "2",
		}, {
			Name:        "STORJ_MAIL_SMTPSERVER_ADDRESS",
			Description: "smtp server address",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_TEMPLATE_PATH",
			Description: "path to email templates source",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_FROM",
			Description: "sender email address",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_AUTH_TYPE",
			Description: "smtp authentication type",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_LOGIN",
			Description: "plain/login auth user login",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_PASSWORD",
			Description: "plain/login auth user password",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_REFRESH_TOKEN",
			Description: "refresh token used to retrieve new access token",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_CLIENT_ID",
			Description: "oauth2 app's client id",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_CLIENT_SECRET",
			Description: "oauth2 app's client secret",
			Default:     "",
		}, {
			Name:        "STORJ_MAIL_TOKEN_URI",
			Description: "uri which is used when retrieving new access token",
			Default:     "",
		}, {
			Name:        "STORJ_PAYMENTS_PROVIDER",
			Description: "payments provider to use",
			Default:     "",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_STRIPE_SECRET_KEY",
			Description: "stripe API secret key",
			Default:     "",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_STRIPE_PUBLIC_KEY",
			Description: "stripe API public key",
			Default:     "",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_STRIPE_FREE_TIER_COUPON_ID",
			Description: "stripe free tier coupon ID",
			Default:     "",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_COINPAYMENTS_PUBLIC_KEY",
			Description: "coinpayments API public key",
			Default:     "",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_COINPAYMENTS_PRIVATE_KEY",
			Description: "coinpayments API private key key",
			Default:     "",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_TRANSACTION_UPDATE_INTERVAL",
			Description: "amount of time we wait before running next transaction update loop",
			Default:     "2m",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_ACCOUNT_BALANCE_UPDATE_INTERVAL",
			Description: "amount of time we wait before running next account balance update loop",
			Default:     "2m",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_CONVERSION_RATES_CYCLE_INTERVAL",
			Description: "amount of time we wait before running next conversion rates update loop",
			Default:     "10m",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_AUTO_ADVANCE",
			Description: "toogle autoadvance feature for invoice creation",
			Default:     "false",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_LISTING_LIMIT",
			Description: "sets the maximum amount of items before we start paging on requests",
			Default:     "100",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_GOB_FLOAT_MIGRATION_BATCH_INTERVAL",
			Description: "amount of time to wait between gob-encoded big.Float database migration batches",
			Default:     "1m",
		}, {
			Name:        "STORJ_PAYMENTS_STRIPE_COIN_PAYMENTS_GOB_FLOAT_MIGRATION_BATCH_SIZE",
			Description: "number of rows with gob-encoded big.Float values to migrate at once",
			Default:     "100",
		}, {
			Name:        "STORJ_PAYMENTS_STORAGE_TBPRICE",
			Description: "price user should pay for storing TB per month",
			Default:     "4",
		}, {
			Name:        "STORJ_PAYMENTS_EGRESS_TBPRICE",
			Description: "price user should pay for each TB of egress",
			Default:     "7",
		}, {
			Name:        "STORJ_PAYMENTS_SEGMENT_PRICE",
			Description: "price user should pay for each segment stored in network per month",
			Default:     "0",
		}, {
			Name:        "STORJ_PAYMENTS_BONUS_RATE",
			Description: "amount of percents that user will earn as bonus credits by depositing in STORJ tokens",
			Default:     "10",
		}, {
			Name:        "STORJ_PAYMENTS_NODE_EGRESS_BANDWIDTH_PRICE",
			Description: "price node receive for storing TB of egress in cents",
			Default:     "2000",
		}, {
			Name:        "STORJ_PAYMENTS_NODE_REPAIR_BANDWIDTH_PRICE",
			Description: "price node receive for storing TB of repair in cents",
			Default:     "1000",
		}, {
			Name:        "STORJ_PAYMENTS_NODE_AUDIT_BANDWIDTH_PRICE",
			Description: "price node receive for storing TB of audit in cents",
			Default:     "1000",
		}, {
			Name:        "STORJ_PAYMENTS_NODE_DISK_SPACE_PRICE",
			Description: "price node receive for storing disk space in cents/TB",
			Default:     "150",
		}, {
			Name:        "STORJ_CONSOLE_ADDRESS",
			Description: "server address of the graphql api gateway and frontend app",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_STATIC_DIR",
			Description: "path to static resources",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_WATCH",
			Description: "whether to load templates on each request",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_EXTERNAL_ADDRESS",
			Description: "external endpoint of the satellite if hosted",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_AUTH_TOKEN",
			Description: "auth token needed for access to registration token creation endpoint",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_AUTH_TOKEN_SECRET",
			Description: "secret used to sign auth tokens",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_CONTACT_INFO_URL",
			Description: "url link to contacts page",
			Default:     "https://forum.storj.io",
		}, {
			Name:        "STORJ_CONSOLE_FRAME_ANCESTORS",
			Description: "allow domains to embed the satellite in a frame, space separated",
			Default:     "tardigrade.io storj.io",
		}, {
			Name:        "STORJ_CONSOLE_LET_US_KNOW_URL",
			Description: "url link to let us know page",
			Default:     "https://storjlabs.atlassian.net/servicedesk/customer/portals",
		}, {
			Name:        "STORJ_CONSOLE_SEO",
			Description: "used to communicate with web crawlers and other web robots",
			Default:     "User-agent: *Disallow: Disallow: /cgi-bin/",
		}, {
			Name:        "STORJ_CONSOLE_SATELLITE_NAME",
			Description: "used to display at web satellite console",
			Default:     "Storj",
		}, {
			Name:        "STORJ_CONSOLE_SATELLITE_OPERATOR",
			Description: "name of organization which set up satellite",
			Default:     "Storj Labs",
		}, {
			Name:        "STORJ_CONSOLE_TERMS_AND_CONDITIONS_URL",
			Description: "url link to terms and conditions page",
			Default:     "https://storj.io/storage-sla/",
		}, {
			Name:        "STORJ_CONSOLE_ACCOUNT_ACTIVATION_REDIRECT_URL",
			Description: "url link for account activation redirect",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_PARTNERED_SATELLITES",
			Description: "names and addresses of partnered satellites in JSON list format",
			Default:     "[[\"US1\",\"https://us1.storj.io\"],[\"EU1\",\"https://eu1.storj.io\"],[\"AP1\",\"https://ap1.storj.io\"]]",
		}, {
			Name:        "STORJ_CONSOLE_GENERAL_REQUEST_URL",
			Description: "url link to general request page",
			Default:     "https://supportdcs.storj.io/hc/en-us/requests/new?ticket_form_id=360000379291",
		}, {
			Name:        "STORJ_CONSOLE_PROJECT_LIMITS_INCREASE_REQUEST_URL",
			Description: "url link to project limit increase request page",
			Default:     "https://supportdcs.storj.io/hc/en-us/requests/new?ticket_form_id=360000683212",
		}, {
			Name:        "STORJ_CONSOLE_GATEWAY_CREDENTIALS_REQUEST_URL",
			Description: "url link for gateway credentials requests",
			Default:     "https://auth.us1.storjshare.io",
		}, {
			Name:        "STORJ_CONSOLE_IS_BETA_SATELLITE",
			Description: "indicates if satellite is in beta",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_BETA_SATELLITE_FEEDBACK_URL",
			Description: "url link for for beta satellite feedback",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_BETA_SATELLITE_SUPPORT_URL",
			Description: "url link for for beta satellite support",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_DOCUMENTATION_URL",
			Description: "url link to documentation",
			Default:     "https://docs.storj.io/",
		}, {
			Name:        "STORJ_CONSOLE_COUPON_CODE_BILLING_UIENABLED",
			Description: "indicates if user is allowed to add coupon codes to account from billing",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_COUPON_CODE_SIGNUP_UIENABLED",
			Description: "indicates if user is allowed to add coupon codes to account from signup",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_FILE_BROWSER_FLOW_DISABLED",
			Description: "indicates if file browser flow is disabled",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_CSPENABLED",
			Description: "indicates if Content Security Policy is enabled",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_LINKSHARING_URL",
			Description: "url link for linksharing requests",
			Default:     "https://link.us1.storjshare.io",
		}, {
			Name:        "STORJ_CONSOLE_PATHWAY_OVERVIEW_ENABLED",
			Description: "indicates if the overview onboarding step should render with pathways",
			Default:     "true",
		}, {
			Name:        "STORJ_CONSOLE_NEW_PROJECT_DASHBOARD",
			Description: "indicates if new project dashboard should be used",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_NEW_NAVIGATION",
			Description: "indicates if new navigation structure should be rendered",
			Default:     "true",
		}, {
			Name:        "STORJ_CONSOLE_NEW_OBJECTS_FLOW",
			Description: "indicates if new objects flow should be used",
			Default:     "true",
		}, {
			Name:        "STORJ_CONSOLE_GENERATED_APIENABLED",
			Description: "indicates if generated console api should be used",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_INACTIVITY_TIMER_ENABLED",
			Description: "indicates if session can be timed out due inactivity",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_INACTIVITY_TIMER_DELAY",
			Description: "inactivity timer delay in seconds",
			Default:     "600",
		}, {
			Name:        "STORJ_CONSOLE_RATE_LIMIT_DURATION",
			Description: "the rate at which request are allowed",
			Default:     "5m",
		}, {
			Name:        "STORJ_CONSOLE_RATE_LIMIT_BURST",
			Description: "number of events before the limit kicks in",
			Default:     "5",
		}, {
			Name:        "STORJ_CONSOLE_RATE_LIMIT_NUM_LIMITS",
			Description: "number of clients whose rate limits we store",
			Default:     "1000",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_PASSWORD_COST",
			Description: "password hashing cost (0=automatic)",
			Default:     "0",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_OPEN_REGISTRATION_ENABLED",
			Description: "enable open registration",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_DEFAULT_PROJECT_LIMIT",
			Description: "default project limits for users",
			Default:     "1",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_TOKEN_EXPIRATION_TIME",
			Description: "expiration time for auth tokens, account recovery tokens, and activation tokens",
			Default:     "24h",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_AS_OF_SYSTEM_TIME_DURATION",
			Description: "default duration for AS OF SYSTEM TIME",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_USAGE_LIMITS_STORAGE_FREE",
			Description: "the default free-tier storage usage limit",
			Default:     "150.00GB",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_USAGE_LIMITS_STORAGE_PAID",
			Description: "the default paid-tier storage usage limit",
			Default:     "25.00TB",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_USAGE_LIMITS_BANDWIDTH_FREE",
			Description: "the default free-tier bandwidth usage limit",
			Default:     "150.00GB",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_USAGE_LIMITS_BANDWIDTH_PAID",
			Description: "the default paid-tier bandwidth usage limit",
			Default:     "100.00TB",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_USAGE_LIMITS_SEGMENT_FREE",
			Description: "the default free-tier segment usage limit",
			Default:     "150000",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_USAGE_LIMITS_SEGMENT_PAID",
			Description: "the default paid-tier segment usage limit",
			Default:     "1000000",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_RECAPTCHA_ENABLED",
			Description: "whether or not reCAPTCHA is enabled for user registration",
			Default:     "false",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_RECAPTCHA_SITE_KEY",
			Description: "reCAPTCHA site key",
			Default:     "",
		}, {
			Name:        "STORJ_CONSOLE_CONFIG_RECAPTCHA_SECRET_KEY",
			Description: "reCAPTCHA secret key",
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
			Name:        "STORJ_GRACEFUL_EXIT_ENABLED",
			Description: "whether or not graceful exit is enabled on the satellite side.",
			Default:     "true",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_CHORE_BATCH_SIZE",
			Description: "size of the buffer used to batch inserts into the transfer queue.",
			Default:     "500",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_CHORE_INTERVAL",
			Description: "how often to run the transfer queue chore.",
			Default:     "",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_ENDPOINT_BATCH_SIZE",
			Description: "size of the buffer used to batch transfer queue reads and sends to the storage node.",
			Default:     "300",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_MAX_FAILURES_PER_PIECE",
			Description: "maximum number of transfer failures per piece.",
			Default:     "5",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_OVERALL_MAX_FAILURES_PERCENTAGE",
			Description: "maximum percentage of transfer failures per node.",
			Default:     "10",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_MAX_INACTIVE_TIME_FRAME",
			Description: "maximum inactive time frame of transfer activities per node.",
			Default:     "168h",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_RECV_TIMEOUT",
			Description: "the minimum duration for receiving a stream from a storage node before timing out",
			Default:     "2h",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_MAX_ORDER_LIMIT_SEND_COUNT",
			Description: "maximum number of order limits a satellite sends to a node before marking piece transfer failed",
			Default:     "10",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_NODE_MIN_AGE_IN_MONTHS",
			Description: "minimum age for a node on the network in order to initiate graceful exit",
			Default:     "6",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_AS_OF_SYSTEM_TIME_INTERVAL",
			Description: "interval for AS OF SYSTEM TIME clause (crdb specific) to read from db at a specific time in the past",
			Default:     "-10s",
		}, {
			Name:        "STORJ_GRACEFUL_EXIT_TRANSFER_QUEUE_BATCH_SIZE",
			Description: "batch size (crdb specific) for deleting and adding items to the transfer queue",
			Default:     "1000",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_AT_REST_GBHOURS_VALUE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_AT_REST_GBHOURS_EXP",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_GET_TB_VALUE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_GET_TB_EXP",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_PUT_TB_VALUE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_PUT_TB_EXP",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_GET_REPAIR_TB_VALUE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_GET_REPAIR_TB_EXP",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_PUT_REPAIR_TB_VALUE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_PUT_REPAIR_TB_EXP",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_GET_AUDIT_TB_VALUE",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_RATES_GET_AUDIT_TB_EXP",
			Description: "",
			Default:     "",
		}, {
			Name:        "STORJ_COMPENSATION_WITHHELD_PERCENTS",
			Description: "comma separated monthly withheld percentage rates",
			Default:     "75,75,75,50,50,50,25,25,25,0,0,0,0,0,0",
		}, {
			Name:        "STORJ_COMPENSATION_DISPOSE_PERCENT",
			Description: "percent of held amount disposed to node after leaving withheld",
			Default:     "50",
		}, {
			Name:        "STORJ_PROJECT_LIMIT_CACHE_CAPACITY",
			Description: "number of projects to cache.",
			Default:     "",
		}, {
			Name:        "STORJ_PROJECT_LIMIT_CACHE_EXPIRATION",
			Description: "how long to cache the project limits.",
			Default:     "",
		}, {
			Name:        "STORJ_ANALYTICS_SEGMENT_WRITE_KEY",
			Description: "segment write key",
			Default:     "",
		}, {
			Name:        "STORJ_ANALYTICS_ENABLED",
			Description: "enable analytics reporting",
			Default:     "false",
		}, {
			Name:        "STORJ_ANALYTICS_HUB_SPOT_APIKEY",
			Description: "hubspot api key",
			Default:     "",
		}, {
			Name:        "STORJ_ANALYTICS_HUB_SPOT_CHANNEL_SIZE",
			Description: "the number of events that can be in the queue before dropping",
			Default:     "1000",
		}, {
			Name:        "STORJ_ANALYTICS_HUB_SPOT_CONCURRENT_SENDS",
			Description: "the number of concurrent api requests that can be made",
			Default:     "4",
		}, {
			Name:        "STORJ_ANALYTICS_HUB_SPOT_DEFAULT_TIMEOUT",
			Description: "the default timeout for the hubspot http client",
			Default:     "10s",
		},
	}
}
