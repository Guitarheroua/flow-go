package netconf

import (
	"time"

	"github.com/onflow/flow-go/network/p2p/p2pconf"
)

// Config encapsulation of configuration structs for all components related to the Flow network.
type Config struct {
	UnicastConfig                 `mapstructure:",squash"`
	p2pconf.ResourceManagerConfig `mapstructure:",squash"`
	ConnectionManagerConfig       `mapstructure:",squash"`
	// GossipSubConfig core gossipsub configuration.
	p2pconf.GossipSubConfig `mapstructure:",squash"`
	AlspConfig              `mapstructure:",squash"`

	// NetworkConnectionPruning determines whether connections to nodes
	// that are not part of protocol state should be trimmed
	// TODO: solely a fallback mechanism, can be removed upon reliable behavior in production.
	NetworkConnectionPruning bool `mapstructure:"networking-connection-pruning"`
	// PreferredUnicastProtocols list of unicast protocols in preferred order
	PreferredUnicastProtocols       []string      `mapstructure:"preferred-unicast-protocols"`
	NetworkReceivedMessageCacheSize uint32        `validate:"gt=0" mapstructure:"received-message-cache-size"`
	PeerUpdateInterval              time.Duration `validate:"gt=0s" mapstructure:"peerupdate-interval"`
	UnicastMessageTimeout           time.Duration `validate:"gt=0s" mapstructure:"unicast-message-timeout"`

	DNSCacheTTL time.Duration `validate:"gt=0s" mapstructure:"dns-cache-ttl"`
	// DisallowListNotificationCacheSize size of the queue for notifications about new peers in the disallow list.
	DisallowListNotificationCacheSize uint32 `validate:"gt=0" mapstructure:"disallow-list-notification-cache-size"`
}

type UnicastConfig struct {
	// UnicastRateLimitersConfig configuration for all unicast rate limiters.
	UnicastRateLimitersConfig `mapstructure:",squash"`
	// UnicastCreateStreamRetryDelay initial delay used in the exponential backoff for create stream retries
	UnicastCreateStreamRetryDelay time.Duration `validate:"gt=0s" mapstructure:"unicast-create-stream-retry-delay"`

	// UnicastStreamZeroRetryResetThreshold is the threshold that determines when to reset the stream creation retry budget to the default value.
	//
	// For example the default value of 100 means that if the stream creation retry budget is decreased to 0, then it will be reset to default value
	// when the number of consecutive successful streams reaches 100.
	//
	// This is to prevent the retry budget from being reset too frequently, as the retry budget is used to gauge the reliability of the stream creation.
	// When the stream creation retry budget is reset to the default value, it means that the stream creation is reliable enough to be trusted again.
	// This parameter mandates when the stream creation is reliable enough to be trusted again; i.e., when the number of consecutive successful streams reaches this threshold.
	// Note that the counter is reset to 0 when the stream creation fails, so the value of for example 100 means that the stream creation is reliable enough that the recent
	// 100 stream creations are all successful.
	UnicastStreamZeroRetryResetThreshold uint64 `validate:"gt=0" mapstructure:"unicast-stream-zero-retry-reset-threshold"`

	// UnicastDialZeroRetryResetThreshold is the threshold that determines when to reset the dial retry budget to the default value.
	// For example the threshold of 1 hour means that if the dial retry budget is decreased to 0, then it will be reset to default value
	// when it has been 1 hour since the last successful dial.
	//
	// This is to prevent the retry budget from being reset too frequently, as the retry budget is used to gauge the reliability of the dialing a remote peer.
	// When the dial retry budget is reset to the default value, it means that the dialing is reliable enough to be trusted again.
	// This parameter mandates when the dialing is reliable enough to be trusted again; i.e., when it has been 1 hour since the last successful dial.
	// Note that the last dial attempt timestamp is reset to zero when the dial fails, so the value of for example 1 hour means that the dialing to the remote peer is reliable enough that the last
	// successful dial attempt was 1 hour ago.
	UnicastDialZeroRetryResetThreshold time.Duration `validate:"gt=0s" mapstructure:"unicast-dial-zero-retry-reset-threshold"`

	// UnicastMaxDialRetryAttemptTimes is the maximum number of attempts to be made to connect to a remote node to establish a unicast (1:1) connection before we give up.
	UnicastMaxDialRetryAttemptTimes uint64 `validate:"gt=0" mapstructure:"unicast-max-dial-retry-attempt-times"`

	// UnicastMaxStreamCreationRetryAttemptTimes is the maximum number of attempts to be made to create a stream to a remote node over a direct unicast (1:1) connection before we give up.
	UnicastMaxStreamCreationRetryAttemptTimes uint64 `validate:"gt=0" mapstructure:"unicast-max-stream-creation-retry-attempt-times"`
}

// UnicastRateLimitersConfig unicast rate limiter configuration for the message and bandwidth rate limiters.
type UnicastRateLimitersConfig struct {
	// DryRun setting this to true will disable connection disconnects and gating when unicast rate limiters are configured
	DryRun bool `mapstructure:"unicast-dry-run"`
	// LockoutDuration the number of seconds a peer will be forced to wait before being allowed to successfully reconnect to the node
	// after being rate limited.
	LockoutDuration time.Duration `validate:"gte=0" mapstructure:"unicast-lockout-duration"`
	// MessageRateLimit amount of unicast messages that can be sent by a peer per second.
	MessageRateLimit int `validate:"gte=0" mapstructure:"unicast-message-rate-limit"`
	// BandwidthRateLimit bandwidth size in bytes a peer is allowed to send via unicast streams per second.
	BandwidthRateLimit int `validate:"gte=0" mapstructure:"unicast-bandwidth-rate-limit"`
	// BandwidthBurstLimit bandwidth size in bytes a peer is allowed to send via unicast streams at once.
	BandwidthBurstLimit int `validate:"gte=0" mapstructure:"unicast-bandwidth-burst-limit"`
}

// AlspConfig is the config for the Application Layer Spam Prevention (ALSP) protocol.
type AlspConfig struct {
	// Size of the cache for spam records. There is at most one spam record per authorized (i.e., staked) node.
	// Recommended size is 10 * number of authorized nodes to allow for churn.
	SpamRecordCacheSize uint32 `mapstructure:"alsp-spam-record-cache-size"`

	// SpamReportQueueSize is the size of the queue for spam records. The queue is used to store spam records
	// temporarily till they are picked by the workers. When the queue is full, new spam records are dropped.
	// Recommended size is 100 * number of authorized nodes to allow for churn.
	SpamReportQueueSize uint32 `mapstructure:"alsp-spam-report-queue-size"`

	// DisablePenalty indicates whether applying the penalty to the misbehaving node is disabled.
	// When disabled, the ALSP module logs the misbehavior reports and updates the metrics, but does not apply the penalty.
	// This is useful for managing production incidents.
	// Note: under normal circumstances, the ALSP module should not be disabled.
	DisablePenalty bool `mapstructure:"alsp-disable-penalty"`

	// HeartBeatInterval is the interval between heartbeats sent by the ALSP module. The heartbeats are recurring
	// events that are used to perform critical ALSP tasks, such as updating the spam records cache.
	HearBeatInterval time.Duration `mapstructure:"alsp-heart-beat-interval"`
}
