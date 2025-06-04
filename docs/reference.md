# API Reference - Redis Provider v0.1.2

## Package Overview

```go
package redis

import "go.fork.vn/redis"
```

Redis Provider cung cấp tích hợp Redis hoàn chỉnh cho Fork Framework với hỗ trợ multiple connections, clustering, và sentinel.

## Core Interfaces

### ServiceProvider

```go
type ServiceProvider interface {
    di.ServiceProvider
    Register(app contracts.Application) error
    Boot(app contracts.Application) error
    Provides() []string
    When() []string
}
```

ServiceProvider quản lý lifecycle của Redis services trong ứng dụng.

#### Methods

##### `Register(app contracts.Application) error`

Đăng ký Redis services vào DI container.

**Parameters:**
- `app contracts.Application`: Application instance

**Returns:**
- `error`: Error nếu registration thất bại

**Services Registered:**
- `redis.manager`: Manager singleton
- `redis`: Alias cho manager
- `redis.connection.{name}`: Named connections

##### `Boot(app contracts.Application) error`

Bootstrap Redis provider sau khi all services đã được registered.

**Parameters:**
- `app contracts.Application`: Application instance

**Returns:**
- `error`: Error nếu bootstrap thất bại

##### `Provides() []string`

Trả về danh sách services mà provider cung cấp.

**Returns:**
- `[]string`: Service names

##### `When() []string`

Trả về danh sách dependencies mà provider cần.

**Returns:**
- `[]string`: Dependency names

#### Constructor

```go
func NewServiceProvider() *ServiceProvider
```

Tạo Redis ServiceProvider mới.

**Returns:**
- `*ServiceProvider`: ServiceProvider instance

**Example:**
```go
provider := redis.NewServiceProvider()
app.RegisterProviders(provider)
```

### Manager

```go
type Manager interface {
    Connection(name ...string) Connection
    DefaultConnection() Connection
    Extend(name string, resolver func(config map[string]interface{}) Connection)
    Config() map[string]interface{}
    SetDefaultConnection(name string)
    Close() error
    HealthCheck(ctx context.Context) error
}
```

Manager quản lý multiple Redis connections và cung cấp connection factory.

#### Methods

##### `Connection(name ...string) Connection`

Lấy Redis connection theo tên. Nếu không có tên, trả về default connection.

**Parameters:**
- `name ...string`: Connection name (optional)

**Returns:**
- `Connection`: Redis connection instance

**Example:**
```go
// Default connection
conn := manager.Connection()

// Named connection
cluster := manager.Connection("cluster")
```

##### `DefaultConnection() Connection`

Lấy default Redis connection.

**Returns:**
- `Connection`: Default connection instance

##### `Extend(name string, resolver func(config map[string]interface{}) Connection)`

Đăng ký custom connection driver.

**Parameters:**
- `name string`: Driver name
- `resolver func(config map[string]interface{}) Connection`: Factory function

**Example:**
```go
manager.Extend("custom", func(config map[string]interface{}) Connection {
    return NewCustomConnection(config)
})
```

##### `Config() map[string]interface{}`

Lấy Redis configuration.

**Returns:**
- `map[string]interface{}`: Configuration map

##### `SetDefaultConnection(name string)`

Đặt default connection name.

**Parameters:**
- `name string`: Connection name

##### `Close() error`

Đóng tất cả Redis connections.

**Returns:**
- `error`: Error nếu close thất bại

##### `HealthCheck(ctx context.Context) error`

Kiểm tra health của tất cả connections.

**Parameters:**
- `ctx context.Context`: Context for timeout/cancellation

**Returns:**
- `error`: Error nếu health check thất bại

#### Constructors

```go
func NewManager(config map[string]interface{}) *Manager
```

Tạo Manager mới với configuration.

**Parameters:**
- `config map[string]interface{}`: Redis configuration

**Returns:**
- `*Manager`: Manager instance

```go
func NewManagerWithConfig(config *Config) *Manager
```

Tạo Manager mới với Config struct.

**Parameters:**
- `config *Config`: Typed configuration

**Returns:**
- `*Manager`: Manager instance

### Connection

```go
type Connection interface {
    redis.Cmdable
    Client() redis.UniversalClient
    Config() ConnectionConfig
    Ping(ctx context.Context) error
    Close() error
    Pipeline() redis.Pipeliner
    TxPipeline() redis.Pipeliner
}
```

Connection là abstraction layer cho Redis operations.

#### Methods

##### `Client() redis.UniversalClient`

Lấy underlying Redis client.

**Returns:**
- `redis.UniversalClient`: Redis client instance

##### `Config() ConnectionConfig`

Lấy connection configuration.

**Returns:**
- `ConnectionConfig`: Configuration struct

##### `Ping(ctx context.Context) error`

Test connection health.

**Parameters:**
- `ctx context.Context`: Context for timeout/cancellation

**Returns:**
- `error`: Error nếu connection unhealthy

##### `Close() error`

Đóng connection.

**Returns:**
- `error`: Error nếu close thất bại

##### `Pipeline() redis.Pipeliner`

Tạo Redis pipeline cho batch operations.

**Returns:**
- `redis.Pipeliner`: Pipeline instance

**Example:**
```go
pipe := conn.Pipeline()
pipe.Set(ctx, "key1", "value1", 0)
pipe.Set(ctx, "key2", "value2", 0)
results, err := pipe.Exec(ctx)
```

##### `TxPipeline() redis.Pipeliner`

Tạo transactional pipeline.

**Returns:**
- `redis.Pipeliner`: Transactional pipeline instance

#### Redis.Cmdable Methods

Connection implements `redis.Cmdable` interface, cung cấp tất cả Redis commands:

**String Commands:**
- `Set(ctx, key, value, expiration)` - Set string value
- `Get(ctx, key)` - Get string value  
- `MSet(ctx, pairs...)` - Set multiple strings
- `MGet(ctx, keys...)` - Get multiple strings
- `Incr(ctx, key)` - Increment integer value
- `Decr(ctx, key)` - Decrement integer value

**Hash Commands:**
- `HSet(ctx, key, field, value)` - Set hash field
- `HGet(ctx, key, field)` - Get hash field
- `HMSet(ctx, key, fields)` - Set multiple hash fields
- `HGetAll(ctx, key)` - Get all hash fields

**List Commands:**
- `LPush(ctx, key, values...)` - Push to list head
- `RPush(ctx, key, values...)` - Push to list tail
- `LPop(ctx, key)` - Pop from list head
- `RPop(ctx, key)` - Pop from list tail
- `LLen(ctx, key)` - Get list length

**Set Commands:**
- `SAdd(ctx, key, members...)` - Add set members
- `SRem(ctx, key, members...)` - Remove set members
- `SMembers(ctx, key)` - Get all set members
- `SCard(ctx, key)` - Get set size

**Sorted Set Commands:**
- `ZAdd(ctx, key, members...)` - Add sorted set members
- `ZRem(ctx, key, members...)` - Remove sorted set members
- `ZRange(ctx, key, start, stop)` - Get range by rank
- `ZRangeByScore(ctx, key, min, max)` - Get range by score

**Key Commands:**
- `Del(ctx, keys...)` - Delete keys
- `Exists(ctx, keys...)` - Check if keys exist
- `Expire(ctx, key, expiration)` - Set key expiration
- `TTL(ctx, key)` - Get key time to live

## Configuration

### Config Struct

```go
type Config struct {
    Default     string                            `mapstructure:"default"`
    Connections map[string]ConnectionConfig       `mapstructure:"connections"`
}
```

Main configuration structure cho Redis provider.

#### Fields

- **Default** `string`: Default connection name
- **Connections** `map[string]ConnectionConfig`: Connection definitions

### ConnectionConfig Struct

```go
type ConnectionConfig struct {
    Driver          string        `mapstructure:"driver"`
    Host            string        `mapstructure:"host"`
    Port            int           `mapstructure:"port"`
    Hosts           []string      `mapstructure:"hosts"`
    Password        string        `mapstructure:"password"`
    Username        string        `mapstructure:"username"`
    Database        int           `mapstructure:"database"`
    MasterName      string        `mapstructure:"master_name"`
    Timeout         time.Duration `mapstructure:"timeout"`
    DialTimeout     time.Duration `mapstructure:"dial_timeout"`
    ReadTimeout     time.Duration `mapstructure:"read_timeout"`
    WriteTimeout    time.Duration `mapstructure:"write_timeout"`
    PoolSize        int           `mapstructure:"pool_size"`
    MinIdleConns    int           `mapstructure:"min_idle_conns"`
    MaxConnAge      time.Duration `mapstructure:"max_conn_age"`
    PoolTimeout     time.Duration `mapstructure:"pool_timeout"`
    IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
    MaxRetries      int           `mapstructure:"max_retries"`
    MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`
    MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`
    TLS             *TLSConfig    `mapstructure:"tls"`
}
```

Configuration cho individual Redis connection.

#### Fields

**Connection Settings:**
- **Driver** `string`: Connection driver ("redis", "cluster", "sentinel")
- **Host** `string`: Redis host (for single instance)
- **Port** `int`: Redis port (for single instance)
- **Hosts** `[]string`: Multiple hosts (for cluster/sentinel)
- **Password** `string`: Authentication password
- **Username** `string`: Authentication username (Redis 6+)
- **Database** `int`: Database number (0-15)
- **MasterName** `string`: Master name (for sentinel)

**Timeout Settings:**
- **Timeout** `time.Duration`: General timeout (default: 10s)
- **DialTimeout** `time.Duration`: Connection timeout (default: 5s)
- **ReadTimeout** `time.Duration`: Read timeout (default: 3s)
- **WriteTimeout** `time.Duration`: Write timeout (default: 3s)

**Pool Settings:**
- **PoolSize** `int`: Maximum number of connections (default: 10)
- **MinIdleConns** `int`: Minimum idle connections (default: 5)
- **MaxConnAge** `time.Duration`: Maximum connection age (default: 30m)
- **PoolTimeout** `time.Duration`: Pool get timeout (default: 4s)
- **IdleTimeout** `time.Duration`: Idle connection timeout (default: 5m)

**Retry Settings:**
- **MaxRetries** `int`: Maximum retry attempts (default: 3)
- **MinRetryBackoff** `time.Duration`: Minimum retry backoff (default: 100ms)
- **MaxRetryBackoff** `time.Duration`: Maximum retry backoff (default: 3s)

**Security:**
- **TLS** `*TLSConfig`: TLS configuration

### TLSConfig Struct

```go
type TLSConfig struct {
    Enabled            bool   `mapstructure:"enabled"`
    CertFile           string `mapstructure:"cert_file"`
    KeyFile            string `mapstructure:"key_file"`
    CAFile             string `mapstructure:"ca_file"`
    InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify"`
}
```

TLS configuration cho secure connections.

#### Fields

- **Enabled** `bool`: Enable TLS
- **CertFile** `string`: Client certificate file path
- **KeyFile** `string`: Client private key file path
- **CAFile** `string`: CA certificate file path
- **InsecureSkipVerify** `bool`: Skip certificate verification

## Utility Functions

### Configuration Helpers

```go
func DefaultConfig() *Config
```

Tạo default configuration.

**Returns:**
- `*Config`: Default configuration instance

```go
func DefaultConnectionConfig() ConnectionConfig
```

Tạo default connection configuration.

**Returns:**
- `ConnectionConfig`: Default connection configuration

### Connection Builders

```go
func NewRedisConnection(config ConnectionConfig) Connection
```

Tạo standard Redis connection.

**Parameters:**
- `config ConnectionConfig`: Connection configuration

**Returns:**
- `Connection`: Redis connection instance

```go
func NewClusterConnection(config ConnectionConfig) Connection
```

Tạo Redis Cluster connection.

**Parameters:**
- `config ConnectionConfig`: Connection configuration

**Returns:**
- `Connection`: Cluster connection instance

```go
func NewSentinelConnection(config ConnectionConfig) Connection
```

Tạo Redis Sentinel connection.

**Parameters:**
- `config ConnectionConfig`: Connection configuration

**Returns:**
- `Connection`: Sentinel connection instance

## Error Types

### RedisError

```go
type RedisError struct {
    Operation string
    Key       string
    Err       error
}

func (e *RedisError) Error() string
func (e *RedisError) Unwrap() error
```

Custom error type cho Redis operations.

#### Fields

- **Operation** `string`: Redis operation name
- **Key** `string`: Redis key involved
- **Err** `error`: Underlying error

### ConnectionError

```go
type ConnectionError struct {
    ConnectionName string
    Err           error
}

func (e *ConnectionError) Error() string
func (e *ConnectionError) Unwrap() error
```

Error type cho connection issues.

#### Fields

- **ConnectionName** `string`: Connection name
- **Err** `error`: Underlying error

## Constants

### Driver Names

```go
const (
    DriverRedis    = "redis"
    DriverCluster  = "cluster" 
    DriverSentinel = "sentinel"
)
```

### Service Names

```go
const (
    ServiceManager    = "redis.manager"
    ServiceConnection = "redis.connection"
    ServiceRedis      = "redis"
)
```

### Default Values

```go
const (
    DefaultHost           = "127.0.0.1"
    DefaultPort           = 6379
    DefaultDatabase       = 0
    DefaultTimeout        = 10 * time.Second
    DefaultDialTimeout    = 5 * time.Second
    DefaultReadTimeout    = 3 * time.Second
    DefaultWriteTimeout   = 3 * time.Second
    DefaultPoolSize       = 10
    DefaultMinIdleConns   = 5
    DefaultMaxConnAge     = 30 * time.Minute
    DefaultPoolTimeout    = 4 * time.Second
    DefaultIdleTimeout    = 5 * time.Minute
    DefaultMaxRetries     = 3
    DefaultMinRetryBackoff = 100 * time.Millisecond
    DefaultMaxRetryBackoff = 3 * time.Second
)
```

## Examples

### Basic Usage

```go
// Create manager
config := redis.DefaultConfig()
config.Connections["default"] = redis.ConnectionConfig{
    Driver:   "redis",
    Host:     "localhost",
    Port:     6379,
    Database: 0,
}

manager := redis.NewManagerWithConfig(config)

// Get connection
conn := manager.Connection()

// Basic operations
ctx := context.Background()
err := conn.Set(ctx, "key", "value", 0).Err()
if err != nil {
    return err
}

val, err := conn.Get(ctx, "key").Result()
if err != nil {
    return err
}
```

### Cluster Configuration

```go
config := redis.DefaultConfig()
config.Connections["cluster"] = redis.ConnectionConfig{
    Driver: "cluster",
    Hosts: []string{
        "cluster-node-1:7000",
        "cluster-node-2:7000", 
        "cluster-node-3:7000",
    },
    Password: "cluster-password",
}

manager := redis.NewManagerWithConfig(config)
cluster := manager.Connection("cluster")
```

### Sentinel Configuration

```go
config := redis.DefaultConfig()
config.Connections["sentinel"] = redis.ConnectionConfig{
    Driver:     "sentinel",
    MasterName: "mymaster",
    Hosts: []string{
        "sentinel-1:26379",
        "sentinel-2:26379",
        "sentinel-3:26379",
    },
    Password: "sentinel-password",
}

manager := redis.NewManagerWithConfig(config)
sentinel := manager.Connection("sentinel")
```

### Pipeline Operations

```go
conn := manager.Connection()
pipe := conn.Pipeline()

pipe.Set(ctx, "key1", "value1", 0)
pipe.Set(ctx, "key2", "value2", 0)
pipe.Incr(ctx, "counter")

results, err := pipe.Exec(ctx)
if err != nil {
    return err
}

for _, result := range results {
    fmt.Printf("Result: %v\n", result)
}
```

### Health Monitoring

```go
// Single connection health check
conn := manager.Connection()
if err := conn.Ping(ctx); err != nil {
    log.Printf("Connection unhealthy: %v", err)
}

// Manager-level health check (all connections)
if err := manager.HealthCheck(ctx); err != nil {
    log.Printf("Redis health check failed: %v", err)
}
```

### Custom Driver Extension

```go
// Define custom connection
type CustomConnection struct {
    config ConnectionConfig
    client redis.UniversalClient
}

func NewCustomConnection(config map[string]interface{}) Connection {
    // Parse config and create custom connection
    return &CustomConnection{
        config: parseConfig(config),
        client: createCustomClient(config),
    }
}

// Register custom driver
manager.Extend("custom", NewCustomConnection)

// Use custom connection
custom := manager.Connection("custom")
```
