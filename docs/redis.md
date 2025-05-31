# Redis Provider - Tài liệu Kỹ thuật

## Tổng quan

Redis Provider là một package cung cấp khả năng tích hợp Redis vào ứng dụng Go một cách dễ dàng và hiệu quả. Package này hỗ trợ cả Redis standalone, Redis Cluster và Redis Sentinel thông qua Universal Client.

## Kiến trúc

### Thành phần chính

1. **ServiceProvider**: Quản lý việc đăng ký Redis services vào DI container
2. **Manager**: Interface chính để quản lý các kết nối Redis
3. **Config**: Cấu trúc cấu hình linh hoạt cho nhiều loại Redis deployment
4. **Client Types**: Hỗ trợ cả Standard Client và Universal Client

### Sơ đồ kiến trúc

```
Application
    ↓
ServiceProvider (redis.NewServiceProvider)
    ↓
Manager (redis.Manager)
    ├── Standard Client (redis.Client)
    └── Universal Client (redis.UniversalClient)
```

## API Reference

### ServiceProvider Interface

```go
type ServiceProvider interface {
    di.ServiceProvider
    Register(app interface{})
    Boot(app interface{})
    Providers() []string
    Requires() []string
}
```

**Khởi tạo:**
```go
provider := redis.NewServiceProvider()
```

**Phương thức:**

- `Register(app interface{})`: Đăng ký Redis services vào DI container
- `Boot(app interface{})`: Khởi động provider (không cần thực hiện gì đặc biệt)
- `Providers() []string`: Trả về danh sách services được cung cấp
- `Requires() []string`: Trả về danh sách dependencies (config)

### Manager Interface

```go
type Manager interface {
    Client() (*redis.Client, error)
    UniversalClient() (redis.UniversalClient, error)
    GetConfig() *Config
    SetConfig(config *Config)
    Close() error
    Ping(ctx context.Context) error
    ClusterPing(ctx context.Context) error
}
```

**Khởi tạo:**
```go
// Với cấu hình mặc định
manager := redis.NewManager()

// Với cấu hình tùy chỉnh
config := redis.DefaultConfig()
// ... cấu hình config
manager := redis.NewManagerWithConfig(config)
```

**Phương thức chính:**

#### `Client() (*redis.Client, error)`
Trả về Redis Standard Client. Thích hợp cho Redis standalone.

```go
client, err := manager.Client()
if err != nil {
    log.Fatal(err)
}

// Sử dụng client
result := client.Set(ctx, "key", "value", 0)
```

#### `UniversalClient() (redis.UniversalClient, error)`
Trả về Redis Universal Client. Hỗ trợ Cluster, Sentinel và standalone.

```go
universalClient, err := manager.UniversalClient()
if err != nil {
    log.Fatal(err)
}

// Sử dụng universal client (API giống Standard Client)
result := universalClient.Set(ctx, "key", "value", 0)
```

#### `Ping(ctx context.Context) error`
Kiểm tra kết nối tới Redis server.

```go
ctx := context.Background()
if err := manager.Ping(ctx); err != nil {
    log.Printf("Redis connection failed: %v", err)
}
```

#### `Close() error`
Đóng tất cả các kết nối Redis.

```go
defer func() {
    if err := manager.Close(); err != nil {
        log.Printf("Error closing Redis connections: %v", err)
    }
}()
```

## Cấu hình

### Cấu trúc Config

```go
type Config struct {
    Client    *ClientConfig    `mapstructure:"client"`
    Universal *UniversalConfig `mapstructure:"universal"`
}
```

### ClientConfig (Standard Client)

```go
type ClientConfig struct {
    Host            string `mapstructure:"host"`
    Port            int    `mapstructure:"port"`
    Password        string `mapstructure:"password"`
    DB              int    `mapstructure:"db"`
    Prefix          string `mapstructure:"prefix"`
    Timeout         int    `mapstructure:"timeout"`
    DialTimeout     int    `mapstructure:"dial_timeout"`
    ReadTimeout     int    `mapstructure:"read_timeout"`
    WriteTimeout    int    `mapstructure:"write_timeout"`
    PoolSize        int    `mapstructure:"pool_size"`
    MinIdleConns    int    `mapstructure:"min_idle_conns"`
}
```

**Cấu hình mặc định:**
- Host: `127.0.0.1`
- Port: `6379`
- DB: `0`
- Timeout: `10` seconds
- DialTimeout: `5` seconds
- ReadTimeout: `3` seconds
- WriteTimeout: `3` seconds
- PoolSize: `10`
- MinIdleConns: `5`

### UniversalConfig (Universal Client)

```go
type UniversalConfig struct {
    Addresses       []string `mapstructure:"addresses"`
    Password        string   `mapstructure:"password"`
    DB              int      `mapstructure:"db"`
    Prefix          string   `mapstructure:"prefix"`
    Timeout         int      `mapstructure:"timeout"`
    DialTimeout     int      `mapstructure:"dial_timeout"`
    ReadTimeout     int      `mapstructure:"read_timeout"`
    WriteTimeout    int      `mapstructure:"write_timeout"`
    MaxRetries      int      `mapstructure:"max_retries"`
    MinRetryBackoff int      `mapstructure:"min_retry_backoff"`
    MaxRetryBackoff int      `mapstructure:"max_retry_backoff"`
    PoolSize        int      `mapstructure:"pool_size"`
    MinIdleConns    int      `mapstructure:"min_idle_conns"`
    ClusterMode     bool     `mapstructure:"cluster_mode"`
    MaxRedirects    int      `mapstructure:"max_redirects"`
    SentinelMode    bool     `mapstructure:"sentinel_mode"`
    MasterName      string   `mapstructure:"master_name"`
}
```

**Cấu hình mặc định:**
- Addresses: `["127.0.0.1:6379"]`
- DB: `0`
- Timeout: `10` seconds
- MaxRetries: `3`
- MinRetryBackoff: `100` milliseconds
- MaxRetryBackoff: `3000` milliseconds
- PoolSize: `10`
- MinIdleConns: `5`
- MaxRedirects: `3`

## DI Container Integration

Khi đăng ký với DI container, Redis Provider tạo ra các services sau:

1. **`redis.manager`**: Instance của Manager interface
2. **`redis.client`**: Instance của *redis.Client (nếu cấu hình client có sẵn)
3. **`redis.universal`**: Instance của redis.UniversalClient (nếu cấu hình universal có sẵn)

### Sử dụng với DI Container

```go
// Trong application boot
app.RegisterProviders(
    config.NewServiceProvider(),
    redis.NewServiceProvider(),
)

// Lấy services từ container
manager := container.MustMake("redis.manager").(redis.Manager)
client := container.MustMake("redis.client").(*redis.Client)
universalClient := container.MustMake("redis.universal").(redis.UniversalClient)
```

## Error Handling

### Các lỗi phổ biến

1. **Missing configuration**: Khi config không được cung cấp
2. **Connection errors**: Khi không thể kết nối tới Redis
3. **Authentication errors**: Khi password không đúng
4. **Timeout errors**: Khi thao tác vượt quá thời gian chờ

### Best practices

```go
// Kiểm tra kết nối trước khi sử dụng
ctx := context.Background()
if err := manager.Ping(ctx); err != nil {
    return fmt.Errorf("redis connection failed: %w", err)
}

// Sử dụng context với timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result := client.Get(ctx, "key")
if err := result.Err(); err != nil {
    if err == redis.Nil {
        // Key không tồn tại
        return nil, nil
    }
    return nil, fmt.Errorf("redis get error: %w", err)
}
```

## Performance

### Connection Pooling

Redis Provider tự động quản lý connection pool để tối ưu hiệu suất:

- **PoolSize**: Số kết nối tối đa trong pool
- **MinIdleConns**: Số kết nối rảnh tối thiểu
- **DialTimeout**: Timeout khi tạo kết nối mới

### Best practices

1. **Sử dụng connection pooling**: Không tạo client mới cho mỗi request
2. **Đặt timeout phù hợp**: Tránh blocking quá lâu
3. **Monitoring**: Theo dõi số lượng kết nối và latency
4. **Graceful shutdown**: Gọi `Close()` khi ứng dụng tắt

## Security

### Authentication
```go
config := redis.DefaultConfig()
config.Client.Password = "your_secure_password"
```

### Network Security
- Sử dụng TLS khi kết nối qua mạng công cộng
- Cấu hình firewall để hạn chế truy cập
- Sử dụng VPN hoặc private network

## Monitoring và Debugging

### Health Check
```go
func healthCheck(manager redis.Manager) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return manager.Ping(ctx)
}
```

### Metrics
Theo dõi các metrics quan trọng:
- Connection pool usage
- Command latency
- Error rate
- Memory usage

## Dependencies

Redis Provider phụ thuộc vào:

1. **go.fork.vn/di**: Dependency injection container
2. **go.fork.vn/config**: Configuration management
3. **github.com/redis/go-redis/v9**: Redis client library

## Thread Safety

Tất cả các components của Redis Provider đều thread-safe:

- Manager instance có thể được sử dụng đồng thời
- Redis clients (Standard và Universal) đều thread-safe
- Connection pooling được quản lý tự động

## Migration

Khi nâng cấp từ các phiên bản cũ:

1. Cập nhật import path từ `github.com/go-fork/providers/redis` sang `go.fork.vn/redis`
2. Kiểm tra breaking changes trong CHANGELOG.md
3. Test kỹ lưỡng trong môi trường staging
