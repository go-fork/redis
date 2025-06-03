# Overview - Redis Provider v0.1.1

## Kiến trúc tổng quan

Redis Provider được thiết kế theo nguyên tắc **Clean Architecture** và **Dependency Injection** của Go-Fork Framework.

### Sơ đồ kiến trúc

```
┌─────────────────────────────────────┐
│           Application               │
├─────────────────────────────────────┤
│         Service Providers           │
│  ┌─────────────────────────────────┐ │
│  │     Redis Service Provider     │ │
│  │  ┌─────────────────────────────┐│ │
│  │  │         Manager             ││ │
│  │  │  ┌─────────────────────────┐││ │
│  │  │  │     Connections         │││ │
│  │  │  │ ┌─────┐ ┌───────┐ ┌────┐│││ │
│  │  │  │ │Redis│ │Cluster│ │Sent││││ │
│  │  │  │ └─────┘ └───────┘ └────┘│││ │
│  │  │  └─────────────────────────┘││ │
│  │  └─────────────────────────────┘│ │
│  └─────────────────────────────────┘ │
├─────────────────────────────────────┤
│           DI Container              │
└─────────────────────────────────────┘
```

## Thành phần chính

### 1. ServiceProvider

Quản lý việc đăng ký và khởi tạo Redis services trong DI container.

```go
type ServiceProvider interface {
    di.ServiceProvider
    Register(app contracts.Application) error
    Boot(app contracts.Application) error
    Provides() []string
    When() []string
}
```

**Responsibilities:**
- Đăng ký Redis Manager vào DI container
- Cấu hình connections từ config
- Khởi tạo health checks
- Quản lý lifecycle của connections

### 2. Manager

Interface chính để quản lý multiple Redis connections.

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

**Responsibilities:**
- Quản lý multiple connections
- Connection pooling
- Configuration management
- Health monitoring

### 3. Connection

Abstraction layer cho Redis client operations.

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

**Responsibilities:**
- Redis command execution
- Pipeline operations
- Transaction support
- Connection health checks

## Patterns & Principles

### Service Provider Pattern

Redis Provider implements Laravel-style Service Provider pattern:

```go
// Registration phase
func (p *ServiceProvider) Register(app contracts.Application) error {
    return app.Container().Singleton("redis.manager", func() (interface{}, error) {
        return NewManager(app.Config().Sub("redis"))
    })
}

// Bootstrap phase
func (p *ServiceProvider) Boot(app contracts.Application) error {
    // Perform any necessary bootstrapping
    return nil
}
```

### Manager Pattern

Quản lý multiple connections với configuration-driven approach:

```go
manager := redis.NewManager(config)

// Default connection
conn := manager.Connection()

// Named connection
cluster := manager.Connection("cluster")
sentinel := manager.Connection("sentinel")
```

### Factory Pattern

Dynamic connection creation với pluggable drivers:

```go
// Built-in drivers
drivers := map[string]DriverFactory{
    "redis":     NewRedisConnection,
    "cluster":   NewClusterConnection,
    "sentinel":  NewSentinelConnection,
}

// Custom driver extension
manager.Extend("custom", func(config map[string]interface{}) Connection {
    return NewCustomConnection(config)
})
```

## Configuration Management

### Hierarchical Configuration

```yaml
redis:
  default: "redis"              # Default connection name
  connections:                  # Connection definitions
    redis:                      # Standard Redis
      driver: "redis"
      host: "127.0.0.1"
      port: 6379
      
    cluster:                    # Redis Cluster
      driver: "cluster"
      hosts:
        - "node1:7000"
        - "node2:7000"
        
    sentinel:                   # Redis Sentinel
      driver: "sentinel" 
      master_name: "mymaster"
      hosts:
        - "sentinel1:26379"
        - "sentinel2:26379"
```

### Environment Override

Configuration có thể được override bằng environment variables:

```bash
REDIS_HOST=production-redis.com
REDIS_PORT=6380
REDIS_PASSWORD=secret
```

## Connection Lifecycle

### Connection States

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Created   │───▶│ Connecting  │───▶│  Connected  │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       │                   ▼                   ▼
       │            ┌─────────────┐    ┌─────────────┐
       └───────────▶│    Error    │    │   Closed    │
                    └─────────────┘    └─────────────┘
```

### Automatic Reconnection

- Exponential backoff retry logic
- Circuit breaker pattern
- Health check monitoring
- Graceful degradation

## Performance Considerations

### Connection Pooling

```go
type ConnectionConfig struct {
    PoolSize        int           `mapstructure:"pool_size"`        // Default: 10
    MinIdleConns    int           `mapstructure:"min_idle_conns"`   // Default: 5
    MaxConnAge      time.Duration `mapstructure:"max_conn_age"`     // Default: 30m
    PoolTimeout     time.Duration `mapstructure:"pool_timeout"`     // Default: 4s
    IdleTimeout     time.Duration `mapstructure:"idle_timeout"`     // Default: 5m
}
```

### Pipelining Support

```go
// Pipeline for batch operations
pipe := connection.Pipeline()
pipe.Set(ctx, "key1", "value1", 0)
pipe.Set(ctx, "key2", "value2", 0)
pipe.Incr(ctx, "counter")

results, err := pipe.Exec(ctx)
```

### Memory Optimization

- Lazy connection initialization
- Resource pooling và reuse
- Automatic cleanup
- Memory leak prevention

## Security Features

### TLS/SSL Support

```yaml
redis:
  connections:
    secure:
      driver: "redis"
      host: "secure-redis.com"
      port: 6380
      tls:
        enabled: true
        cert_file: "/path/to/client.crt"
        key_file: "/path/to/client.key"
        ca_file: "/path/to/ca.crt"
        insecure_skip_verify: false
```

### Authentication

```yaml
redis:
  connections:
    authenticated:
      driver: "redis"
      host: "redis.com"
      port: 6379
      password: "${REDIS_PASSWORD}"
      username: "${REDIS_USERNAME}"  # Redis 6+ ACL support
```

## Testing Strategy

### Interface-Based Testing

Sử dụng interfaces để dễ dàng mock trong unit tests:

```go
func TestMyService(t *testing.T) {
    mockManager := mocks.NewMockManager(t)
    mockConnection := mocks.NewMockConnection(t)
    
    mockManager.On("Connection").Return(mockConnection)
    mockConnection.On("Set", mock.Anything, "key", "value", 0).
        Return(redis.NewStatusResult("OK", nil))
    
    service := NewMyService(mockManager)
    err := service.StoreValue("key", "value")
    
    assert.NoError(t, err)
    mockManager.AssertExpectations(t)
}
```

### Integration Testing

```go
func TestRedisIntegration(t *testing.T) {
    // Start test Redis container
    container := testcontainers.SetupRedis(t)
    defer container.Close()
    
    // Configure manager với test connection
    config := redis.DefaultConfig()
    config.Host = container.Host()
    config.Port = container.Port()
    
    manager := redis.NewManagerWithConfig(config)
    connection := manager.Connection()
    
    // Test operations
    ctx := context.Background()
    err := connection.Set(ctx, "test", "value", 0).Err()
    assert.NoError(t, err)
}
```

## Error Handling

### Centralized Error Management

```go
type RedisError struct {
    Operation string
    Key       string
    Err       error
}

func (e *RedisError) Error() string {
    return fmt.Sprintf("redis %s operation failed for key '%s': %v", 
        e.Operation, e.Key, e.Err)
}
```

### Retry Logic

```go
type RetryConfig struct {
    MaxRetries      int           `mapstructure:"max_retries"`
    MinBackoff      time.Duration `mapstructure:"min_backoff"`
    MaxBackoff      time.Duration `mapstructure:"max_backoff"`
    BackoffFunction string        `mapstructure:"backoff_function"` // linear, exponential
}
```

## Monitoring & Observability

### Metrics Collection

- Connection pool utilization
- Command execution latency
- Error rates và types
- Memory usage
- Network I/O statistics

### Health Checks

```go
func (m *manager) HealthCheck(ctx context.Context) error {
    var errs []error
    
    for name, conn := range m.connections {
        if err := conn.Ping(ctx); err != nil {
            errs = append(errs, fmt.Errorf("connection %s: %w", name, err))
        }
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("health check failed: %v", errs)
    }
    
    return nil
}
```

### Logging Integration

Tích hợp với logging system của Go-Fork:

```go
logger := app.Logger()
logger.With("component", "redis").
    With("connection", connectionName).
    Info("Redis connection established")
```

## Extensibility

### Custom Drivers

```go
// Implement custom Redis driver
func NewCustomRedisConnection(config map[string]interface{}) Connection {
    return &customConnection{
        config: config,
        client: createCustomClient(config),
    }
}

// Register custom driver
manager.Extend("custom", NewCustomRedisConnection)
```

### Middleware Support

```go
type Middleware interface {
    Handle(ctx context.Context, cmd redis.Cmder) error
}

// Add middleware to connection
connection.AddMiddleware(NewLoggingMiddleware())
connection.AddMiddleware(NewMetricsMiddleware())
```

## Migration Guide

### From v0.1.0 to v0.1.1

**Breaking Changes:**
1. Manager interface redesigned
2. Configuration structure simplified
3. Connection interface introduced
4. DI integration updated

**Migration Steps:**
1. Update import paths
2. Update configuration files
3. Refactor Manager usage
4. Update tests với new interfaces
5. Regenerate mocks

## Next Steps

- [Reference](reference.md) - Complete API documentation
- [Usage](usage.md) - Detailed usage examples
- [Examples](https://github.com/go-fork/recipes/tree/main/examples/redis)
