package cmd

import (
	"errors"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"heckel.io/ntfy/server"
	"heckel.io/ntfy/util"
	"log"
	"time"
)

var flagsServe = []cli.Flag{
	&cli.StringFlag{Name: "config", Aliases: []string{"c"}, EnvVars: []string{"NTFY_CONFIG_FILE"}, Value: "/etc/ntfy/server.yml", DefaultText: "/etc/ntfy/server.yml", Usage: "config file"},
	altsrc.NewStringFlag(&cli.StringFlag{Name: "base-url", Aliases: []string{"B"}, EnvVars: []string{"NTFY_BASE_URL"}, Usage: "externally visible base URL for this host (e.g. https://ntfy.sh)"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "listen-http", Aliases: []string{"l"}, EnvVars: []string{"NTFY_LISTEN_HTTP"}, Value: server.DefaultListenHTTP, Usage: "ip:port used to as HTTP listen address"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "listen-https", Aliases: []string{"L"}, EnvVars: []string{"NTFY_LISTEN_HTTPS"}, Usage: "ip:port used to as HTTPS listen address"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "key-file", Aliases: []string{"K"}, EnvVars: []string{"NTFY_KEY_FILE"}, Usage: "private key file, if listen-https is set"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "cert-file", Aliases: []string{"E"}, EnvVars: []string{"NTFY_CERT_FILE"}, Usage: "certificate file, if listen-https is set"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "firebase-key-file", Aliases: []string{"F"}, EnvVars: []string{"NTFY_FIREBASE_KEY_FILE"}, Usage: "Firebase credentials file; if set additionally publish to FCM topic"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "cache-file", Aliases: []string{"C"}, EnvVars: []string{"NTFY_CACHE_FILE"}, Usage: "cache file used for message caching"}),
	altsrc.NewDurationFlag(&cli.DurationFlag{Name: "cache-duration", Aliases: []string{"b"}, EnvVars: []string{"NTFY_CACHE_DURATION"}, Value: server.DefaultCacheDuration, Usage: "buffer messages for this time to allow `since` requests"}),
	altsrc.NewDurationFlag(&cli.DurationFlag{Name: "keepalive-interval", Aliases: []string{"k"}, EnvVars: []string{"NTFY_KEEPALIVE_INTERVAL"}, Value: server.DefaultKeepaliveInterval, Usage: "interval of keepalive messages"}),
	altsrc.NewDurationFlag(&cli.DurationFlag{Name: "manager-interval", Aliases: []string{"m"}, EnvVars: []string{"NTFY_MANAGER_INTERVAL"}, Value: server.DefaultManagerInterval, Usage: "interval of for message pruning and stats printing"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "smtp-sender-addr", EnvVars: []string{"NTFY_SMTP_SENDER_ADDR"}, Usage: "SMTP server address (host:port) for outgoing emails"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "smtp-sender-user", EnvVars: []string{"NTFY_SMTP_SENDER_USER"}, Usage: "SMTP user (if e-mail sending is enabled)"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "smtp-sender-pass", EnvVars: []string{"NTFY_SMTP_SENDER_PASS"}, Usage: "SMTP password (if e-mail sending is enabled)"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "smtp-sender-from", EnvVars: []string{"NTFY_SMTP_SENDER_FROM"}, Usage: "SMTP sender address (if e-mail sending is enabled)"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "smtp-server-listen", EnvVars: []string{"NTFY_SMTP_SERVER_LISTEN"}, Usage: "SMTP server address (ip:port) for incoming emails, e.g. :25"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "smtp-server-domain", EnvVars: []string{"NTFY_SMTP_SERVER_DOMAIN"}, Usage: "SMTP domain for incoming e-mail, e.g. ntfy.sh"}),
	altsrc.NewStringFlag(&cli.StringFlag{Name: "smtp-server-addr-prefix", EnvVars: []string{"NTFY_SMTP_SERVER_ADDR_PREFIX"}, Usage: "SMTP email address prefix for topics to prevent spam (e.g. 'ntfy-')"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "global-topic-limit", Aliases: []string{"T"}, EnvVars: []string{"NTFY_GLOBAL_TOPIC_LIMIT"}, Value: server.DefaultGlobalTopicLimit, Usage: "total number of topics allowed"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "visitor-subscription-limit", EnvVars: []string{"NTFY_VISITOR_SUBSCRIPTION_LIMIT"}, Value: server.DefaultVisitorSubscriptionLimit, Usage: "number of subscriptions per visitor"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "visitor-request-limit-burst", EnvVars: []string{"NTFY_VISITOR_REQUEST_LIMIT_BURST"}, Value: server.DefaultVisitorRequestLimitBurst, Usage: "initial limit of requests per visitor"}),
	altsrc.NewDurationFlag(&cli.DurationFlag{Name: "visitor-request-limit-replenish", EnvVars: []string{"NTFY_VISITOR_REQUEST_LIMIT_REPLENISH"}, Value: server.DefaultVisitorRequestLimitReplenish, Usage: "interval at which burst limit is replenished (one per x)"}),
	altsrc.NewIntFlag(&cli.IntFlag{Name: "visitor-email-limit-burst", EnvVars: []string{"NTFY_VISITOR_EMAIL_LIMIT_BURST"}, Value: server.DefaultVisitorEmailLimitBurst, Usage: "initial limit of e-mails per visitor"}),
	altsrc.NewDurationFlag(&cli.DurationFlag{Name: "visitor-email-limit-replenish", EnvVars: []string{"NTFY_VISITOR_EMAIL_LIMIT_REPLENISH"}, Value: server.DefaultVisitorEmailLimitReplenish, Usage: "interval at which burst limit is replenished (one per x)"}),
	altsrc.NewBoolFlag(&cli.BoolFlag{Name: "behind-proxy", Aliases: []string{"P"}, EnvVars: []string{"NTFY_BEHIND_PROXY"}, Value: false, Usage: "if set, use X-Forwarded-For header to determine visitor IP address (for rate limiting)"}),
}

var cmdServe = &cli.Command{
	Name:      "serve",
	Usage:     "Run the ntfy server",
	UsageText: "ntfy serve [OPTIONS..]",
	Action:    execServe,
	Flags:     flagsServe,
	Before:    initConfigFileInputSource("config", flagsServe),
	Description: `Run the ntfy server and listen for incoming requests

The command will load the configuration from /etc/ntfy/server.yml. Config options can 
be overridden using the command line options.

Examples:
  ntfy serve                      # Starts server in the foreground (on port 80)
  ntfy serve --listen-http :8080  # Starts server with alternate port`,
}

func execServe(c *cli.Context) error {
	if c.NArg() > 0 {
		return errors.New("no arguments expected, see 'ntfy serve --help' for help")
	}

	// Read all the options
	baseURL := c.String("base-url")
	listenHTTP := c.String("listen-http")
	listenHTTPS := c.String("listen-https")
	keyFile := c.String("key-file")
	certFile := c.String("cert-file")
	firebaseKeyFile := c.String("firebase-key-file")
	cacheFile := c.String("cache-file")
	cacheDuration := c.Duration("cache-duration")
	keepaliveInterval := c.Duration("keepalive-interval")
	managerInterval := c.Duration("manager-interval")
	smtpSenderAddr := c.String("smtp-sender-addr")
	smtpSenderUser := c.String("smtp-sender-user")
	smtpSenderPass := c.String("smtp-sender-pass")
	smtpSenderFrom := c.String("smtp-sender-from")
	smtpServerListen := c.String("smtp-server-listen")
	smtpServerDomain := c.String("smtp-server-domain")
	smtpServerAddrPrefix := c.String("smtp-server-addr-prefix")
	globalTopicLimit := c.Int("global-topic-limit")
	visitorSubscriptionLimit := c.Int("visitor-subscription-limit")
	visitorRequestLimitBurst := c.Int("visitor-request-limit-burst")
	visitorRequestLimitReplenish := c.Duration("visitor-request-limit-replenish")
	visitorEmailLimitBurst := c.Int("visitor-email-limit-burst")
	visitorEmailLimitReplenish := c.Duration("visitor-email-limit-replenish")
	behindProxy := c.Bool("behind-proxy")

	// Check values
	if firebaseKeyFile != "" && !util.FileExists(firebaseKeyFile) {
		return errors.New("if set, FCM key file must exist")
	} else if keepaliveInterval < 5*time.Second {
		return errors.New("keepalive interval cannot be lower than five seconds")
	} else if managerInterval < 5*time.Second {
		return errors.New("manager interval cannot be lower than five seconds")
	} else if cacheDuration > 0 && cacheDuration < managerInterval {
		return errors.New("cache duration cannot be lower than manager interval")
	} else if keyFile != "" && !util.FileExists(keyFile) {
		return errors.New("if set, key file must exist")
	} else if certFile != "" && !util.FileExists(certFile) {
		return errors.New("if set, certificate file must exist")
	} else if listenHTTPS != "" && (keyFile == "" || certFile == "") {
		return errors.New("if listen-https is set, both key-file and cert-file must be set")
	} else if smtpSenderAddr != "" && (baseURL == "" || smtpSenderUser == "" || smtpSenderPass == "" || smtpSenderFrom == "") {
		return errors.New("if smtp-sender-addr is set, base-url, smtp-sender-user, smtp-sender-pass and smtp-sender-from must also be set")
	} else if smtpServerListen != "" && smtpServerDomain == "" {
		return errors.New("if smtp-server-listen is set, smtp-server-domain must also be set")
	}

	// Run server
	conf := server.NewConfig()
	conf.BaseURL = baseURL
	conf.ListenHTTP = listenHTTP
	conf.ListenHTTPS = listenHTTPS
	conf.KeyFile = keyFile
	conf.CertFile = certFile
	conf.FirebaseKeyFile = firebaseKeyFile
	conf.CacheFile = cacheFile
	conf.CacheDuration = cacheDuration
	conf.KeepaliveInterval = keepaliveInterval
	conf.ManagerInterval = managerInterval
	conf.SMTPSenderAddr = smtpSenderAddr
	conf.SMTPSenderUser = smtpSenderUser
	conf.SMTPSenderPass = smtpSenderPass
	conf.SMTPSenderFrom = smtpSenderFrom
	conf.SMTPServerListen = smtpServerListen
	conf.SMTPServerDomain = smtpServerDomain
	conf.SMTPServerAddrPrefix = smtpServerAddrPrefix
	conf.GlobalTopicLimit = globalTopicLimit
	conf.VisitorSubscriptionLimit = visitorSubscriptionLimit
	conf.VisitorRequestLimitBurst = visitorRequestLimitBurst
	conf.VisitorRequestLimitReplenish = visitorRequestLimitReplenish
	conf.VisitorEmailLimitBurst = visitorEmailLimitBurst
	conf.VisitorEmailLimitReplenish = visitorEmailLimitReplenish
	conf.BehindProxy = behindProxy
	s, err := server.New(conf)
	if err != nil {
		log.Fatalln(err)
	}
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Exiting.")
	return nil
}
