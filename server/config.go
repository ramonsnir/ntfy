package server

import (
	"time"
)

// Defines default config settings
const (
	DefaultListenHTTP                = ":80"
	DefaultCacheDuration             = 12 * time.Hour
	DefaultKeepaliveInterval         = 30 * time.Second
	DefaultManagerInterval           = time.Minute
	DefaultAtSenderInterval          = 10 * time.Second
	DefaultMinDelay                  = 10 * time.Second
	DefaultMaxDelay                  = 3 * 24 * time.Hour
	DefaultMessageLimit              = 512
	DefaultFirebaseKeepaliveInterval = time.Hour
)

// Defines all the limits
// - global topic limit: max number of topics overall
// - per visitor request limit: max number of PUT/GET/.. requests (here: 60 requests bucket, replenished at a rate of one per 10 seconds)
// - per visitor email limit: max number of emails (here: 16 email bucket, replenished at a rate of one per hour)
// - per visitor subscription limit: max number of subscriptions (active HTTP connections) per per-visitor/IP
const (
	DefaultGlobalTopicLimit             = 5000
	DefaultVisitorRequestLimitBurst     = 60
	DefaultVisitorRequestLimitReplenish = 10 * time.Second
	DefaultVisitorEmailLimitBurst       = 16
	DefaultVisitorEmailLimitReplenish   = time.Hour
	DefaultVisitorSubscriptionLimit     = 30
)

// Config is the main config struct for the application. Use New to instantiate a default config struct.
type Config struct {
	BaseURL                      string
	ListenHTTP                   string
	ListenHTTPS                  string
	KeyFile                      string
	CertFile                     string
	FirebaseKeyFile              string
	CacheFile                    string
	CacheDuration                time.Duration
	KeepaliveInterval            time.Duration
	ManagerInterval              time.Duration
	AtSenderInterval             time.Duration
	FirebaseKeepaliveInterval    time.Duration
	SMTPSenderAddr               string
	SMTPSenderUser               string
	SMTPSenderPass               string
	SMTPSenderFrom               string
	SMTPServerListen             string
	SMTPServerDomain             string
	SMTPServerAddrPrefix         string
	MessageLimit                 int
	MinDelay                     time.Duration
	MaxDelay                     time.Duration
	GlobalTopicLimit             int
	VisitorRequestLimitBurst     int
	VisitorRequestLimitReplenish time.Duration
	VisitorEmailLimitBurst       int
	VisitorEmailLimitReplenish   time.Duration
	VisitorSubscriptionLimit     int
	BehindProxy                  bool
}

// NewConfig instantiates a default new server config
func NewConfig() *Config {
	return &Config{
		BaseURL:                      "",
		ListenHTTP:                   DefaultListenHTTP,
		ListenHTTPS:                  "",
		KeyFile:                      "",
		CertFile:                     "",
		FirebaseKeyFile:              "",
		CacheFile:                    "",
		CacheDuration:                DefaultCacheDuration,
		KeepaliveInterval:            DefaultKeepaliveInterval,
		ManagerInterval:              DefaultManagerInterval,
		MessageLimit:                 DefaultMessageLimit,
		MinDelay:                     DefaultMinDelay,
		MaxDelay:                     DefaultMaxDelay,
		AtSenderInterval:             DefaultAtSenderInterval,
		FirebaseKeepaliveInterval:    DefaultFirebaseKeepaliveInterval,
		GlobalTopicLimit:             DefaultGlobalTopicLimit,
		VisitorRequestLimitBurst:     DefaultVisitorRequestLimitBurst,
		VisitorRequestLimitReplenish: DefaultVisitorRequestLimitReplenish,
		VisitorEmailLimitBurst:       DefaultVisitorEmailLimitBurst,
		VisitorEmailLimitReplenish:   DefaultVisitorEmailLimitReplenish,
		VisitorSubscriptionLimit:     DefaultVisitorSubscriptionLimit,
		BehindProxy:                  false,
	}
}
