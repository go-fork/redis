package redis

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config là cấu trúc dữ liệu cấu hình chính cho Redis.
type Config struct {
	// Client là cấu hình cho Redis Standard Client.
	Client *ClientConfig `mapstructure:"client"`

	// Universal là cấu hình cho Redis Universal Client (hỗ trợ Cluster, Sentinel, và standalone).
	Universal *UniversalConfig `mapstructure:"universal"`
}

// ClientConfig chứa cấu hình cho Redis Standard Client.
type ClientConfig struct {
	// Enabled xác định liệu Redis Standard Client có được kích hoạt hay không.
	Enabled bool `mapstructure:"enabled"`

	// Network type: "tcp" or "unix"
	Network string `mapstructure:"network"`

	// Addr là địa chỉ máy chủ Redis theo format host:port.
	Addr string `mapstructure:"addr"`

	// Username cho Redis 6+ ACL authentication
	Username string `mapstructure:"username"`

	// Password là mật khẩu xác thực với máy chủ Redis.
	Password string `mapstructure:"password"`

	// DB là số của database Redis sẽ được sử dụng.
	DB int `mapstructure:"db"`

	// ClientName set via CLIENT SETNAME command
	ClientName string `mapstructure:"client_name"`

	// Protocol RESP version (2 or 3)
	Protocol int `mapstructure:"protocol"`

	// DialTimeout cho establishing new connections
	DialTimeout time.Duration `mapstructure:"dial_timeout"`

	// ReadTimeout cho socket reads (-1=no timeout, -2=disable SetReadDeadline)
	ReadTimeout time.Duration `mapstructure:"read_timeout"`

	// WriteTimeout cho socket writes (-1=no timeout, -2=disable SetWriteDeadline)
	WriteTimeout time.Duration `mapstructure:"write_timeout"`

	// ContextTimeoutEnabled controls whether client respects context timeouts
	ContextTimeoutEnabled bool `mapstructure:"context_timeout_enabled"`

	// MaxRetries before giving up (-1 disables retries)
	MaxRetries int `mapstructure:"max_retries"`

	// MinRetryBackoff minimum backoff between retries (-1 disables)
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`

	// MaxRetryBackoff maximum backoff between retries (-1 disables)
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`

	// PoolFIFO type: true=FIFO, false=LIFO
	PoolFIFO bool `mapstructure:"pool_fifo"`

	// PoolSize base number of socket connections
	PoolSize int `mapstructure:"pool_size"`

	// PoolTimeout amount of time client waits for connection when pool is busy
	PoolTimeout time.Duration `mapstructure:"pool_timeout"`

	// MinIdleConns minimum number of idle connections
	MinIdleConns int `mapstructure:"min_idle_conns"`

	// MaxIdleConns maximum number of idle connections
	MaxIdleConns int `mapstructure:"max_idle_conns"`

	// MaxActiveConns maximum number of connections allocated by pool
	MaxActiveConns int `mapstructure:"max_active_conns"`

	// ConnMaxIdleTime maximum amount of time a connection may be idle
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`

	// ConnMaxLifetime maximum amount of time a connection may be reused
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`

	// TLS configuration
	TLS *TLSConfig `mapstructure:"tls"`

	// DisableIdentity disable CLIENT SETINFO command on connect
	DisableIdentity bool `mapstructure:"disable_identity"`

	// IdentitySuffix added to client name for identification
	IdentitySuffix string `mapstructure:"identity_suffix"`

	// UnstableResp3 enables unstable mode for Redis Search module
	UnstableResp3 bool `mapstructure:"unstable_resp3"`
}

// UniversalConfig chứa cấu hình cho Redis Universal Client.
type UniversalConfig struct {
	// Enabled xác định liệu Redis Universal Client có được kích hoạt hay không.
	Enabled bool `mapstructure:"enabled"`

	// Addrs là danh sách các địa chỉ máy chủ Redis (host:port)
	Addrs []string `mapstructure:"addrs"`

	// Username cho Redis 6+ ACL authentication
	Username string `mapstructure:"username"`

	// Password là mật khẩu xác thực với máy chủ Redis.
	Password string `mapstructure:"password"`

	// SentinelUsername cho Sentinel authentication (if different from Redis)
	SentinelUsername string `mapstructure:"sentinel_username"`

	// SentinelPassword cho Sentinel authentication (if different from Redis)
	SentinelPassword string `mapstructure:"sentinel_password"`

	// DB là số của database Redis sẽ được sử dụng.
	DB int `mapstructure:"db"`

	// ClientName set via CLIENT SETNAME command
	ClientName string `mapstructure:"client_name"`

	// Protocol RESP version (2 or 3)
	Protocol int `mapstructure:"protocol"`

	// DialTimeout cho establishing new connections
	DialTimeout time.Duration `mapstructure:"dial_timeout"`

	// ReadTimeout cho socket reads (-1=no timeout, -2=disable SetReadDeadline)
	ReadTimeout time.Duration `mapstructure:"read_timeout"`

	// WriteTimeout cho socket writes (-1=no timeout, -2=disable SetWriteDeadline)
	WriteTimeout time.Duration `mapstructure:"write_timeout"`

	// ContextTimeoutEnabled controls whether client respects context timeouts
	ContextTimeoutEnabled bool `mapstructure:"context_timeout_enabled"`

	// MaxRetries before giving up (-1 disables retries)
	MaxRetries int `mapstructure:"max_retries"`

	// MinRetryBackoff minimum backoff between retries (-1 disables)
	MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`

	// MaxRetryBackoff maximum backoff between retries (-1 disables)
	MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`

	// PoolFIFO type: true=FIFO, false=LIFO
	PoolFIFO bool `mapstructure:"pool_fifo"`

	// PoolSize base number of socket connections
	PoolSize int `mapstructure:"pool_size"`

	// PoolTimeout amount of time client waits for connection when pool is busy
	PoolTimeout time.Duration `mapstructure:"pool_timeout"`

	// MinIdleConns minimum number of idle connections
	MinIdleConns int `mapstructure:"min_idle_conns"`

	// MaxIdleConns maximum number of idle connections
	MaxIdleConns int `mapstructure:"max_idle_conns"`

	// MaxActiveConns maximum number of connections allocated by pool
	MaxActiveConns int `mapstructure:"max_active_conns"`

	// ConnMaxIdleTime maximum amount of time a connection may be idle
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`

	// ConnMaxLifetime maximum amount of time a connection may be reused
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`

	// TLS configuration
	TLS *TLSConfig `mapstructure:"tls"`

	// MaxRedirects maximum number of retries for failed requests
	MaxRedirects int `mapstructure:"max_redirects"`

	// ReadOnly enables read-only commands on slave nodes
	ReadOnly bool `mapstructure:"read_only"`

	// RouteByLatency allows routing read-only commands to the closest node
	RouteByLatency bool `mapstructure:"route_by_latency"`

	// RouteRandomly allows routing read-only commands to random node
	RouteRandomly bool `mapstructure:"route_randomly"`

	// MasterName Sentinel master name
	MasterName string `mapstructure:"master_name"`

	// DisableIdentity disable CLIENT SETINFO command on connect
	DisableIdentity bool `mapstructure:"disable_identity"`

	// IdentitySuffix added to client name for identification
	IdentitySuffix string `mapstructure:"identity_suffix"`

	// UnstableResp3 enables unstable mode for Redis Search module
	UnstableResp3 bool `mapstructure:"unstable_resp3"`

	// IsClusterMode can be used when only one Addr is provided
	IsClusterMode bool `mapstructure:"is_cluster_mode"`
}

// TLSConfig contains TLS configuration for Redis connections.
type TLSConfig struct {
	// CertFile path to client certificate file
	CertFile string `mapstructure:"cert_file"`

	// KeyFile path to client private key file
	KeyFile string `mapstructure:"key_file"`

	// CAFile path to Certificate Authority file
	CAFile string `mapstructure:"ca_file"`

	// ServerName used to verify the hostname on returned certificates
	ServerName string `mapstructure:"server_name"`

	// InsecureSkipVerify controls whether a client verifies server certificate
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
}

// Options trả về redis.Options từ Config.Client.
func (c *Config) Options() *redis.Options {
	if c.Client == nil || !c.Client.Enabled {
		return nil
	}

	client := c.Client
	opts := &redis.Options{
		Network:               client.Network,
		Addr:                  client.Addr,
		Username:              client.Username,
		Password:              client.Password,
		DB:                    client.DB,
		ClientName:            client.ClientName,
		Protocol:              client.Protocol,
		DialTimeout:           client.DialTimeout,
		ReadTimeout:           client.ReadTimeout,
		WriteTimeout:          client.WriteTimeout,
		ContextTimeoutEnabled: client.ContextTimeoutEnabled,
		MaxRetries:            client.MaxRetries,
		MinRetryBackoff:       client.MinRetryBackoff,
		MaxRetryBackoff:       client.MaxRetryBackoff,
		PoolFIFO:              client.PoolFIFO,
		PoolSize:              client.PoolSize,
		PoolTimeout:           client.PoolTimeout,
		MinIdleConns:          client.MinIdleConns,
		MaxIdleConns:          client.MaxIdleConns,
		MaxActiveConns:        client.MaxActiveConns,
		ConnMaxIdleTime:       client.ConnMaxIdleTime,
		ConnMaxLifetime:       client.ConnMaxLifetime,
		DisableIdentity:       client.DisableIdentity,
		IdentitySuffix:        client.IdentitySuffix,
		UnstableResp3:         client.UnstableResp3,
	}

	// Configure TLS if specified
	if client.TLS != nil {
		tlsConfig, err := client.TLS.buildTLSConfig()
		if err == nil {
			opts.TLSConfig = tlsConfig
		}
	}

	return opts
}

// UniversalOptions trả về redis.UniversalOptions từ Config.Universal.
func (c *Config) UniversalOptions() *redis.UniversalOptions {
	if c.Universal == nil || !c.Universal.Enabled {
		return nil
	}

	universal := c.Universal
	opts := &redis.UniversalOptions{
		Addrs:                 universal.Addrs,
		Username:              universal.Username,
		Password:              universal.Password,
		SentinelUsername:      universal.SentinelUsername,
		SentinelPassword:      universal.SentinelPassword,
		DB:                    universal.DB,
		ClientName:            universal.ClientName,
		Protocol:              universal.Protocol,
		DialTimeout:           universal.DialTimeout,
		ReadTimeout:           universal.ReadTimeout,
		WriteTimeout:          universal.WriteTimeout,
		ContextTimeoutEnabled: universal.ContextTimeoutEnabled,
		MaxRetries:            universal.MaxRetries,
		MinRetryBackoff:       universal.MinRetryBackoff,
		MaxRetryBackoff:       universal.MaxRetryBackoff,
		PoolFIFO:              universal.PoolFIFO,
		PoolSize:              universal.PoolSize,
		PoolTimeout:           universal.PoolTimeout,
		MinIdleConns:          universal.MinIdleConns,
		MaxIdleConns:          universal.MaxIdleConns,
		MaxActiveConns:        universal.MaxActiveConns,
		ConnMaxIdleTime:       universal.ConnMaxIdleTime,
		ConnMaxLifetime:       universal.ConnMaxLifetime,
		MaxRedirects:          universal.MaxRedirects,
		ReadOnly:              universal.ReadOnly,
		RouteByLatency:        universal.RouteByLatency,
		RouteRandomly:         universal.RouteRandomly,
		MasterName:            universal.MasterName,
		DisableIdentity:       universal.DisableIdentity,
		IdentitySuffix:        universal.IdentitySuffix,
		UnstableResp3:         universal.UnstableResp3,
		IsClusterMode:         universal.IsClusterMode,
	}

	// Configure TLS if specified
	if universal.TLS != nil {
		tlsConfig, err := universal.TLS.buildTLSConfig()
		if err == nil {
			opts.TLSConfig = tlsConfig
		}
	}

	return opts
}

// buildTLSConfig creates a tls.Config from TLSConfig.
func (t *TLSConfig) buildTLSConfig() (*tls.Config, error) {
	if t == nil {
		return nil, nil
	}
	tlsConfig := &tls.Config{
		ServerName:         t.ServerName,
		InsecureSkipVerify: t.InsecureSkipVerify,
	}
	// Load client certificate and key if specified
	if t.CertFile != "" && t.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate if specified
	if t.CAFile != "" {
		caCert, err := os.ReadFile(t.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

// DefaultConfig trả về cấu hình mặc định cho Redis.
func DefaultConfig() *Config {
	return &Config{
		Client: &ClientConfig{
			Enabled:               false,
			Network:               "tcp",
			Addr:                  "localhost:6379",
			Username:              "",
			Password:              "",
			DB:                    0,
			ClientName:            "",
			Protocol:              3,
			DialTimeout:           5 * time.Second,
			ReadTimeout:           3 * time.Second,
			WriteTimeout:          3 * time.Second,
			ContextTimeoutEnabled: false,
			MaxRetries:            3,
			MinRetryBackoff:       8 * time.Millisecond,
			MaxRetryBackoff:       512 * time.Millisecond,
			PoolFIFO:              false,
			PoolSize:              10,
			PoolTimeout:           4 * time.Second,
			MinIdleConns:          0,
			MaxIdleConns:          0,
			MaxActiveConns:        0,
			ConnMaxIdleTime:       30 * time.Minute,
			ConnMaxLifetime:       0,
			TLS:                   nil,
			DisableIdentity:       false,
			IdentitySuffix:        "",
			UnstableResp3:         false,
		},
		Universal: &UniversalConfig{
			Enabled:               false,
			Addrs:                 []string{"localhost:6379"},
			Username:              "",
			Password:              "",
			SentinelUsername:      "",
			SentinelPassword:      "",
			DB:                    0,
			ClientName:            "",
			Protocol:              3,
			DialTimeout:           5 * time.Second,
			ReadTimeout:           3 * time.Second,
			WriteTimeout:          3 * time.Second,
			ContextTimeoutEnabled: false,
			MaxRetries:            3,
			MinRetryBackoff:       8 * time.Millisecond,
			MaxRetryBackoff:       512 * time.Millisecond,
			PoolFIFO:              false,
			PoolSize:              10,
			PoolTimeout:           4 * time.Second,
			MinIdleConns:          0,
			MaxIdleConns:          0,
			MaxActiveConns:        0,
			ConnMaxIdleTime:       30 * time.Minute,
			ConnMaxLifetime:       0,
			TLS:                   nil,
			MaxRedirects:          3,
			ReadOnly:              false,
			RouteByLatency:        false,
			RouteRandomly:         false,
			MasterName:            "mymaster",
			DisableIdentity:       false,
			IdentitySuffix:        "",
			UnstableResp3:         false,
			IsClusterMode:         false,
		},
	}
}

// Validate validates the ClientConfig configuration.
func (c *ClientConfig) Validate() error {
	if c == nil {
		return errors.New("ClientConfig cannot be nil")
	}

	// Validate network type
	if c.Network != "" && c.Network != "tcp" && c.Network != "unix" {
		return fmt.Errorf("invalid network type: %s, must be 'tcp' or 'unix'", c.Network)
	}

	// Validate address
	if strings.TrimSpace(c.Addr) == "" {
		return errors.New("addr cannot be empty")
	}

	// Validate database number
	if c.DB < 0 {
		return fmt.Errorf("db must be non-negative, got: %d", c.DB)
	}

	// Validate protocol version
	if c.Protocol != 0 && c.Protocol != 2 && c.Protocol != 3 {
		return fmt.Errorf("invalid protocol version: %d, must be 2 or 3", c.Protocol)
	}

	// Validate pool configuration
	if c.PoolSize < 0 {
		return fmt.Errorf("pool_size must be non-negative, got: %d", c.PoolSize)
	}

	if c.MinIdleConns < 0 {
		return fmt.Errorf("min_idle_conns must be non-negative, got: %d", c.MinIdleConns)
	}

	if c.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns must be non-negative, got: %d", c.MaxIdleConns)
	}

	if c.MaxActiveConns < 0 {
		return fmt.Errorf("max_active_conns must be non-negative, got: %d", c.MaxActiveConns)
	}

	if c.MinIdleConns > c.MaxIdleConns && c.MaxIdleConns > 0 {
		return fmt.Errorf("min_idle_conns (%d) cannot be greater than max_idle_conns (%d)", c.MinIdleConns, c.MaxIdleConns)
	}

	// Validate retry configuration
	if c.MaxRetries < -1 {
		return fmt.Errorf("max_retries must be >= -1, got: %d", c.MaxRetries)
	}

	// Validate TLS configuration if present
	if c.TLS != nil {
		if err := c.TLS.validate(); err != nil {
			return fmt.Errorf("TLS configuration error: %w", err)
		}
	}
	return nil
}

// Validate validates the UniversalConfig configuration.
func (c *UniversalConfig) Validate() error {
	if c == nil {
		return errors.New("UniversalConfig cannot be nil")
	}

	// Validate addresses
	if len(c.Addrs) == 0 {
		return errors.New("addrs cannot be empty")
	}

	for i, addr := range c.Addrs {
		if strings.TrimSpace(addr) == "" {
			return fmt.Errorf("addrs[%d] cannot be empty", i)
		}
	}

	// Validate database number
	if c.DB < 0 {
		return fmt.Errorf("db must be non-negative, got: %d", c.DB)
	}

	// Validate protocol version
	if c.Protocol != 0 && c.Protocol != 2 && c.Protocol != 3 {
		return fmt.Errorf("invalid protocol version: %d, must be 2 or 3", c.Protocol)
	}

	// Validate pool configuration
	if c.PoolSize < 0 {
		return fmt.Errorf("pool_size must be non-negative, got: %d", c.PoolSize)
	}

	if c.MinIdleConns < 0 {
		return fmt.Errorf("min_idle_conns must be non-negative, got: %d", c.MinIdleConns)
	}

	if c.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns must be non-negative, got: %d", c.MaxIdleConns)
	}

	if c.MaxActiveConns < 0 {
		return fmt.Errorf("max_active_conns must be non-negative, got: %d", c.MaxActiveConns)
	}

	if c.MinIdleConns > c.MaxIdleConns && c.MaxIdleConns > 0 {
		return fmt.Errorf("min_idle_conns (%d) cannot be greater than max_idle_conns (%d)", c.MinIdleConns, c.MaxIdleConns)
	}

	// Validate retry configuration
	if c.MaxRetries < -1 {
		return fmt.Errorf("max_retries must be >= -1, got: %d", c.MaxRetries)
	}

	// Validate redirects for cluster mode
	if c.MaxRedirects < 0 {
		return fmt.Errorf("max_redirects must be non-negative, got: %d", c.MaxRedirects)
	}

	// Validate master name for Sentinel mode
	if c.MasterName != "" && strings.TrimSpace(c.MasterName) == "" {
		return errors.New("master_name cannot be empty string when specified")
	}

	// Validate TLS configuration if present
	if c.TLS != nil {
		if err := c.TLS.validate(); err != nil {
			return fmt.Errorf("TLS configuration error: %w", err)
		}
	}

	return nil
}

// Validate validates the main Config structure.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("Config cannot be nil")
	}

	// Validate client configuration if present (regardless of enabled status)
	if c.Client != nil {
		if err := c.Client.Validate(); err != nil {
			return fmt.Errorf("client configuration error: %w", err)
		}
	}

	// Validate universal configuration if present (regardless of enabled status)
	if c.Universal != nil {
		if err := c.Universal.Validate(); err != nil {
			return fmt.Errorf("universal configuration error: %w", err)
		}
	}

	return nil
}

// validate validates the TLS configuration.
func (t *TLSConfig) validate() error {
	if t == nil {
		return nil
	}

	// If cert file is specified, key file must also be specified
	if t.CertFile != "" && t.KeyFile == "" {
		return errors.New("key_file must be specified when cert_file is provided")
	}

	if t.KeyFile != "" && t.CertFile == "" {
		return errors.New("cert_file must be specified when key_file is provided")
	}

	// Validate certificate files exist if specified
	if t.CertFile != "" {
		if _, err := os.Stat(t.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("cert_file does not exist: %s", t.CertFile)
		}
	}

	if t.KeyFile != "" {
		if _, err := os.Stat(t.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("key_file does not exist: %s", t.KeyFile)
		}
	}

	if t.CAFile != "" {
		if _, err := os.Stat(t.CAFile); os.IsNotExist(err) {
			return fmt.Errorf("ca_file does not exist: %s", t.CAFile)
		}
	}

	return nil
}

func (c *Config) Enabled() bool {
	// Check if either Client or Universal is enabled
	if c.Client != nil && c.Client.Enabled {
		return true
	}
	if c.Universal != nil && c.Universal.Enabled {
		return true
	}
	return false
}
