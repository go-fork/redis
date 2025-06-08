# Tổng quan kiến trúc Redis Client Library

## Giới thiệu

Go Fork Redis Client Library được thiết kế theo nguyên tắc **Dependency Injection** và **Configuration Management**, cung cấp một giải pháp toàn diện để quản lý kết nối Redis trong các ứng dụng Go hiện đại.

## Kiến trúc tổng thể

```mermaid
graph TD
    subgraph "Application Layer"
        A[Application Code]
        B[Business Logic]
        C[Service Layer]
    end
    
    subgraph "Redis Client Library"
        D[ServiceProvider]
        E[Manager Interface]
        F[Configuration System]
        
        subgraph "Client Management"
            G[Standard Client]
            H[Universal Client]
            I[Connection Pool]
        end
        
        subgraph "Configuration"
            J[ClientConfig]
            K[UniversalConfig]
            L[TLSConfig]
        end
    end
    
    subgraph "Infrastructure"
        M[Redis Standalone]
        N[Redis Cluster]
        O[Redis Sentinel]
    end
    
    A --> D
    B --> E
    C --> E
    D --> F
    E --> G
    E --> H
    G --> I
    H --> I
    F --> J
    F --> K
    F --> L
    I --> M
    I --> N
    I --> O
```

## Thành phần chính

### 1. ServiceProvider

**Mục đích**: Quản lý vòng đời và khởi tạo các thành phần Redis client

```mermaid
sequenceDiagram
    participant App as Application
    participant DI as DI Container
    participant SP as ServiceProvider
    participant Config as Configuration
    participant Manager as Manager
    
    App->>DI: Register ServiceProvider
    DI->>SP: Boot()
    SP->>Config: Load Configuration
    Config-->>SP: Return Config
    SP->>Manager: Create Manager
    SP->>DI: Register Manager
    DI-->>App: Ready to use
```

**Chức năng chính**:
- Đăng ký dependencies vào DI container
- Khởi tạo Redis clients dựa trên cấu hình
- Quản lý lifecycle của connections
- Validation cấu hình

### 2. Manager Interface

**Mục đích**: Cung cấp interface thống nhất để thao tác với Redis

```go
type Manager interface {
    // Basic operations
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Del(ctx context.Context, keys ...string) (int64, error)
    
    // Advanced operations
    Exists(ctx context.Context, keys ...string) (int64, error)
    Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
    TTL(ctx context.Context, key string) (time.Duration, error)
    
    // Hash operations
    HGet(ctx context.Context, key, field string) (string, error)
    HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
    HGetAll(ctx context.Context, key string) (map[string]string, error)
    
    // List operations
    LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
    RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
    LPop(ctx context.Context, key string) (string, error)
    RPop(ctx context.Context, key string) (string, error)
    
    // Set operations
    SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
    SMembers(ctx context.Context, key string) ([]string, error)
    
    // Pub/Sub
    Publish(ctx context.Context, channel string, message interface{}) (int64, error)
    Subscribe(ctx context.Context, channels ...string) *redis.PubSub
    
    // Pipeline
    Pipeline() redis.Pipeliner
    TxPipeline() redis.Pipeliner
    
    // Client access
    Client() redis.UniversalClient
    Close() error
}
```

### 3. Configuration System

**Mục đích**: Quản lý cấu hình linh hoạt cho các client types

```mermaid
graph LR
    subgraph "Configuration Hierarchy"
        A[Config Root]
        B[ClientConfig]
        C[UniversalConfig]
        D[TLSConfig]
        
        A --> B
        A --> C
        B --> D
        C --> D
    end
    
    subgraph "Client Types"
        E[Standard Client]
        F[Universal Client]
        
        B --> E
        C --> F
    end
    
    subgraph "Redis Infrastructure"
        G[Standalone]
        H[Cluster]
        I[Sentinel]
        
        E --> G
        F --> G
        F --> H
        F --> I
    end
```

## Luồng hoạt động

### 1. Khởi tạo ứng dụng

```mermaid
sequenceDiagram
    participant Main as main()
    participant DI as DI Container
    participant Config as Config Provider
    participant Redis as Redis Provider
    participant Manager as Manager
    
    Main->>DI: Create DI Container
    Main->>Config: Create Config Provider
    Main->>DI: Register Config Provider
    Main->>Redis: Create Redis Provider
    Main->>DI: Register Redis Provider
    Main->>DI: Boot Application
    
    DI->>Config: Boot Config Provider
    Config->>Config: Load Configuration
    DI->>Redis: Boot Redis Provider
    Redis->>Config: Get Redis Configuration
    Config-->>Redis: Return Config
    Redis->>Manager: Create Manager
    Redis->>DI: Register Manager
    
    Main->>DI: Resolve Manager
    DI-->>Main: Return Manager Instance
```

### 2. Xử lý request

```mermaid
sequenceDiagram
    participant App as Application
    participant Manager as Manager
    participant Client as Redis Client
    participant Pool as Connection Pool
    participant Redis as Redis Server
    
    App->>Manager: Get/Set/Del operation
    Manager->>Client: Execute command
    Client->>Pool: Get connection
    Pool-->>Client: Return connection
    Client->>Redis: Send command
    Redis-->>Client: Return result
    Client->>Pool: Return connection
    Client-->>Manager: Return result
    Manager-->>App: Return response
```

### 3. Connection pooling

```mermaid
graph TD
    subgraph "Connection Pool"
        A[Pool Manager]
        B[Active Connections]
        C[Idle Connections]
        D[Connection Factory]
    end
    
    subgraph "Configuration"
        E[PoolSize: 10]
        F[MinIdleConns: 2]
        G[MaxIdleConns: 5]
        H[MaxActiveConns: 20]
        I[ConnMaxLifetime: 1h]
        J[ConnMaxIdleTime: 30m]
    end
    
    subgraph "Redis Servers"
        K[Redis 1]
        L[Redis 2]
        M[Redis 3]
    end
    
    A --> B
    A --> C
    A --> D
    E --> A
    F --> A
    G --> A
    H --> A
    I --> A
    J --> A
    
    D --> K
    D --> L
    D --> M
```

## Patterns và Best Practices

### 1. Dependency Injection Pattern

```go
// Đăng ký services
func setupDI() *di.Container {
    app := di.New()
    
    // Config provider
    app.Register("config", config.NewProvider())
    
    // Redis provider
    app.Register("redis", redis.NewServiceProvider())
    
    return app
}

// Sử dụng trong handler
func (h *Handler) GetUser(ctx context.Context, userID string) (*User, error) {
    // Manager được inject thông qua DI
    userKey := fmt.Sprintf("user:%s", userID)
    
    data, err := h.redis.Get(ctx, userKey)
    if err != nil {
        return nil, err
    }
    
    var user User
    if err := json.Unmarshal([]byte(data), &user); err != nil {
        return nil, err
    }
    
    return &user, nil
}
```

### 2. Configuration Pattern

```yaml
# config/app.yaml
redis:
  # Standard client for simple use cases
  client:
    enabled: false
    addr: "localhost:6379"
    password: ""
    db: 0
    
  # Universal client for production
  universal:
    enabled: true
    addrs:
      - "redis-cluster-1:6379"
      - "redis-cluster-2:6379"  
      - "redis-cluster-3:6379"
    password: "${REDIS_PASSWORD}"
    pool_size: 20
    min_idle_conns: 5
    max_idle_conns: 10
    tls:
      cert_file: "/etc/ssl/redis-client.crt"
      key_file: "/etc/ssl/redis-client.key"
      ca_file: "/etc/ssl/ca.crt"
```

### 3. Error Handling Pattern

```go
func (m *manager) Get(ctx context.Context, key string) (string, error) {
    result, err := m.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return "", ErrKeyNotFound
        }
        return "", fmt.Errorf("redis get error: %w", err)
    }
    return result, nil
}
```

## Performance Characteristics

### 1. Connection Pooling

```mermaid
graph TD
    A[Request] --> B{Pool Available?}
    B -->|Yes| C[Get Connection]
    B -->|No| D{Under Max?}
    D -->|Yes| E[Create Connection]
    D -->|No| F[Wait for Connection]
    
    C --> G[Execute Command]
    E --> G
    F --> G
    
    G --> H[Return Connection]
    H --> I[Response]
```

### 2. Memory Management

- **Connection Reuse**: Tái sử dụng connections để giảm overhead
- **Idle Timeout**: Tự động đóng connections không sử dụng
- **Max Lifetime**: Giới hạn tuổi thọ connection để tránh memory leak
- **FIFO/LIFO Pool**: Tối ưu pattern sử dụng connection

### 3. Failure Handling

```mermaid
sequenceDiagram
    participant App as Application
    participant Manager as Manager
    participant Client as Redis Client
    participant Redis as Redis Server
    
    App->>Manager: Execute Command
    Manager->>Client: Send Command
    Client->>Redis: Network Call
    Redis-->>Client: Connection Error
    
    Client->>Client: Retry Logic
    Client->>Redis: Retry Command
    Redis-->>Client: Success Response
    Client-->>Manager: Return Result
    Manager-->>App: Return Response
    
    Note over Client: MaxRetries: 3
    Note over Client: MinRetryBackoff: 8ms
    Note over Client: MaxRetryBackoff: 512ms
```

## Tích hợp với các framework

### 1. Gin Framework

```go
func main() {
    // Setup DI
    app := setupDI()
    
    // Setup Gin
    r := gin.Default()
    
    // Middleware để inject dependencies
    r.Use(func(c *gin.Context) {
        var manager redis.Manager
        app.Resolve("redis.manager", &manager)
        c.Set("redis", manager)
        c.Next()
    })
    
    r.GET("/user/:id", getUserHandler)
    r.Run()
}
```

### 2. Echo Framework

```go
func main() {
    // Setup DI
    app := setupDI()
    
    // Setup Echo
    e := echo.New()
    
    // Middleware
    e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            var manager redis.Manager
            app.Resolve("redis.manager", &manager)
            c.Set("redis", manager)
            return next(c)
        }
    })
    
    e.GET("/user/:id", getUserHandler)
    e.Start(":8080")
}
```

## Monitoring và Observability

### 1. Metrics Collection

```go
// Custom metrics wrapper
type MetricsManager struct {
    redis.Manager
    metrics prometheus.Registerer
}

func (m *MetricsManager) Get(ctx context.Context, key string) (string, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        // Record metrics
        m.recordMetric("get", duration)
    }()
    
    return m.Manager.Get(ctx, key)
}
```

### 2. Logging

```go
// Structured logging
func (m *manager) Get(ctx context.Context, key string) (string, error) {
    logger := log.With(
        "operation", "get",
        "key", key,
        "trace_id", getTraceID(ctx),
    )
    
    logger.Debug("executing redis get command")
    
    result, err := m.client.Get(ctx, key).Result()
    if err != nil {
        logger.Error("redis get failed", "error", err)
        return "", err
    }
    
    logger.Debug("redis get completed", "result_length", len(result))
    return result, nil
}
```

---

**Tiếp theo**: [Cấu hình chi tiết](configuration.md)
