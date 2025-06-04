# Redis Provider - Hướng dẫn Sử dụng

## Giới thiệu

Redis Provider cung cấp một cách đơn giản và hiệu quả để tích hợp Redis vào ứng dụng Go của bạn. Package hỗ trợ đầy đủ các tính năng của Redis bao gồm standalone, cluster và sentinel mode.

## Cài đặt

```bash
go get go.fork.vn/redis@v0.1.2
```

## Import

```go
import "go.fork.vn/redis"
```

## Cấu hình cơ bản

### 1. Redis Standalone

Tạo file cấu hình `config/app.yaml`:

```yaml
redis:
  client:
    enabled: true
    host: "127.0.0.1"
    port: 6379
    password: ""
    db: 0
    prefix: "app:"
    timeout: 5  # seconds
    dial_timeout: 5  # seconds
    read_timeout: 3  # seconds
    write_timeout: 3  # seconds
    pool_size: 10
    min_idle_conns: 5
```

### 2. Redis Cluster

```yaml
redis:
  universal:
    enabled: true
    addresses:
      - "127.0.0.1:7000"
      - "127.0.0.1:7001"
      - "127.0.0.1:7002"
    password: ""
    db: 0
    prefix: "app:"
    timeout: 5  # seconds
    dial_timeout: 5  # seconds
    read_timeout: 3  # seconds
    write_timeout: 3  # seconds
    max_retries: 3
    min_retry_backoff: 8  # milliseconds
    max_retry_backoff: 512  # milliseconds
    pool_size: 10
    min_idle_conns: 5
    
    # Bật mode Cluster
    cluster_mode: true
    max_redirects: 3
```

### 3. Redis Sentinel

```yaml
redis:
  universal:
    enabled: true
    addresses:
      - "127.0.0.1:26379"
      - "127.0.0.1:26380"
      - "127.0.0.1:26381"
    password: ""
    db: 0
    prefix: "app:"
    timeout: 5  # seconds
    dial_timeout: 5  # seconds
    read_timeout: 3  # seconds
    write_timeout: 3  # seconds
    max_retries: 3
    min_retry_backoff: 8  # milliseconds
    max_retry_backoff: 512  # milliseconds
    pool_size: 10
    min_idle_conns: 5
    
    # Bật mode Sentinel
    sentinel_mode: true
    master_name: "mymaster"
```

## Khởi tạo Application

### Cách 1: Sử dụng DI Container (Khuyến nghị)

```go
package main

import (
    "context"
    "log"
    
    "go.fork.vn/di"
    "go.fork.vn/config"
    "go.fork.vn/redis"
)

func main() {
    // Tạo container
    container := di.NewContainer()
    
    // Tạo app với container
    app := &Application{container: container}
    
    // Đăng ký providers
    configProvider := config.NewServiceProvider()
    redisProvider := redis.NewServiceProvider()
    
    configProvider.Register(app)
    redisProvider.Register(app)
    
    configProvider.Boot(app)
    redisProvider.Boot(app)
    
    // Sử dụng Redis
    useRedis(container)
}

type Application struct {
    container *di.Container
}

func (a *Application) Container() *di.Container {
    return a.container
}

func useRedis(container *di.Container) {
    // Lấy manager từ container
    manager := container.MustMake("redis.manager").(redis.Manager)
    
    // Sử dụng Standard Client
    client, err := manager.Client()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx := context.Background()
    
    // Set một key
    err = client.Set(ctx, "greeting", "Hello Redis!", 0).Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // Get key
    value, err := client.Get(ctx, "greeting").Result()
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Value: %s", value)
}
```

### Cách 2: Khởi tạo trực tiếp

```go
package main

import (
    "context"
    "log"
    
    "go.fork.vn/redis"
)

func main() {
    // Tạo config
    config := redis.DefaultConfig()
    config.Client.Host = "localhost"
    config.Client.Port = 6379
    config.Client.Password = "your_password"
    
    // Tạo manager
    manager := redis.NewManagerWithConfig(config)
    defer manager.Close()
    
    // Kiểm tra kết nối
    ctx := context.Background()
    if err := manager.Ping(ctx); err != nil {
        log.Fatal("Redis connection failed:", err)
    }
    
    // Sử dụng client
    client, err := manager.Client()
    if err != nil {
        log.Fatal(err)
    }
    
    // Thực hiện các thao tác Redis
    basicOperations(ctx, client)
}
```

## Các thao tác cơ bản

### 1. String Operations

```go
func stringOperations(ctx context.Context, client *redis.Client) {
    // SET
    err := client.Set(ctx, "user:1:name", "John Doe", time.Hour).Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // GET
    name, err := client.Get(ctx, "user:1:name").Result()
    if err == redis.Nil {
        log.Println("Key không tồn tại")
    } else if err != nil {
        log.Fatal(err)
    } else {
        log.Printf("Name: %s", name)
    }
    
    // INCR
    counter, err := client.Incr(ctx, "page_views").Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Page views: %d", counter)
    
    // MSET (Multiple SET)
    err = client.MSet(ctx, 
        "user:1:email", "john@example.com",
        "user:1:age", "30",
    ).Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // MGET (Multiple GET)
    values, err := client.MGet(ctx, "user:1:name", "user:1:email", "user:1:age").Result()
    if err != nil {
        log.Fatal(err)
    }
    
    for i, value := range values {
        log.Printf("Value %d: %v", i, value)
    }
}
```

### 2. Hash Operations

```go
func hashOperations(ctx context.Context, client *redis.Client) {
    userKey := "user:1"
    
    // HSET
    err := client.HSet(ctx, userKey, map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   30,
    }).Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // HGET
    name, err := client.HGet(ctx, userKey, "name").Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("User name: %s", name)
    
    // HGETALL
    user, err := client.HGetAll(ctx, userKey).Result()
    if err != nil {
        log.Fatal(err)
    }
    
    for field, value := range user {
        log.Printf("%s: %s", field, value)
    }
    
    // HINCRBY
    newAge, err := client.HIncrBy(ctx, userKey, "age", 1).Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("New age: %d", newAge)
}
```

### 3. List Operations

```go
func listOperations(ctx context.Context, client *redis.Client) {
    listKey := "tasks"
    
    // LPUSH (Push to left)
    err := client.LPush(ctx, listKey, "task1", "task2", "task3").Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // RPUSH (Push to right)
    err = client.RPush(ctx, listKey, "task4", "task5").Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // LRANGE (Get range)
    tasks, err := client.LRange(ctx, listKey, 0, -1).Result()
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("All tasks:")
    for i, task := range tasks {
        log.Printf("%d: %s", i, task)
    }
    
    // LPOP (Pop from left)
    task, err := client.LPop(ctx, listKey).Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Popped task: %s", task)
    
    // LLEN (Get length)
    length, err := client.LLen(ctx, listKey).Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("List length: %d", length)
}
```

### 4. Set Operations

```go
func setOperations(ctx context.Context, client *redis.Client) {
    setKey := "user:1:tags"
    
    // SADD
    err := client.SAdd(ctx, setKey, "golang", "redis", "backend", "database").Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // SMEMBERS
    tags, err := client.SMembers(ctx, setKey).Result()
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("User tags:")
    for _, tag := range tags {
        log.Printf("- %s", tag)
    }
    
    // SISMEMBER
    isMember, err := client.SIsMember(ctx, setKey, "golang").Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Is 'golang' member: %t", isMember)
    
    // SCARD (Get cardinality)
    count, err := client.SCard(ctx, setKey).Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Tags count: %d", count)
}
```

### 5. Sorted Set Operations

```go
func sortedSetOperations(ctx context.Context, client *redis.Client) {
    leaderboardKey := "leaderboard"
    
    // ZADD
    members := []redis.Z{
        {Score: 100, Member: "player1"},
        {Score: 85, Member: "player2"},
        {Score: 95, Member: "player3"},
        {Score: 120, Member: "player4"},
    }
    
    err := client.ZAdd(ctx, leaderboardKey, members...).Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // ZREVRANGE (Get range in reverse order - highest score first)
    topPlayers, err := client.ZRevRangeWithScores(ctx, leaderboardKey, 0, 2).Result()
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Top 3 players:")
    for i, player := range topPlayers {
        log.Printf("%d. %s: %.0f", i+1, player.Member, player.Score)
    }
    
    // ZSCORE
    score, err := client.ZScore(ctx, leaderboardKey, "player1").Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Player1 score: %.0f", score)
    
    // ZRANK
    rank, err := client.ZRevRank(ctx, leaderboardKey, "player1").Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Player1 rank: %d", rank+1) // +1 vì rank bắt đầu từ 0
}
```

## Xử lý JSON

```go
import (
    "encoding/json"
    "time"
)

type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
}

func jsonOperations(ctx context.Context, client *redis.Client) {
    user := User{
        ID:      1,
        Name:    "John Doe", 
        Email:   "john@example.com",
        Created: time.Now(),
    }
    
    // Serialize to JSON
    userJSON, err := json.Marshal(user)
    if err != nil {
        log.Fatal(err)
    }
    
    // Store in Redis
    err = client.Set(ctx, "user:1", userJSON, time.Hour).Err()
    if err != nil {
        log.Fatal(err)
    }
    
    // Retrieve from Redis
    userJSONStr, err := client.Get(ctx, "user:1").Result()
    if err != nil {
        log.Fatal(err)
    }
    
    // Deserialize from JSON
    var retrievedUser User
    err = json.Unmarshal([]byte(userJSONStr), &retrievedUser)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Retrieved user: %+v", retrievedUser)
}
```

## Pub/Sub (Publisher/Subscriber)

```go
func pubSubExample(ctx context.Context, client *redis.Client) {
    // Publisher
    go func() {
        for i := 0; i < 10; i++ {
            message := fmt.Sprintf("Message %d", i)
            err := client.Publish(ctx, "notifications", message).Err()
            if err != nil {
                log.Printf("Publish error: %v", err)
                return
            }
            log.Printf("Published: %s", message)
            time.Sleep(time.Second)
        }
    }()
    
    // Subscriber
    pubsub := client.Subscribe(ctx, "notifications")
    defer pubsub.Close()
    
    // Nhận tin nhắn
    ch := pubsub.Channel()
    
    for msg := range ch {
        log.Printf("Received: %s from channel %s", msg.Payload, msg.Channel)
        
        // Thoát sau 5 tin nhắn
        if msg.Payload == "Message 4" {
            break
        }
    }
}
```

## Pipeline và Transaction

### Pipeline

```go
func pipelineExample(ctx context.Context, client *redis.Client) {
    pipe := client.Pipeline()
    
    // Thêm các lệnh vào pipeline
    incr := pipe.Incr(ctx, "pipeline_counter")
    pipe.Expire(ctx, "pipeline_counter", time.Hour)
    set := pipe.Set(ctx, "pipeline_key", "pipeline_value", time.Hour)
    get := pipe.Get(ctx, "pipeline_key")
    
    // Thực thi tất cả lệnh cùng một lúc
    _, err := pipe.Exec(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Lấy kết quả
    log.Printf("Incr result: %d", incr.Val())
    log.Printf("Set result: %s", set.Val())
    log.Printf("Get result: %s", get.Val())
}
```

### Transaction (MULTI/EXEC)

```go
func transactionExample(ctx context.Context, client *redis.Client) {
    // Watch key để đảm bảo không bị thay đổi trong transaction
    err := client.Watch(ctx, func(tx *redis.Tx) error {
        // Kiểm tra giá trị hiện tại
        val, err := tx.Get(ctx, "counter").Int()
        if err != nil && err != redis.Nil {
            return err
        }
        
        // Tạo transaction
        _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
            // Tăng counter
            pipe.Set(ctx, "counter", val+1, 0)
            // Ghi log
            pipe.LPush(ctx, "counter_log", fmt.Sprintf("Incremented to %d at %s", val+1, time.Now().Format(time.RFC3339)))
            return nil
        })
        
        return err
    }, "counter")
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Transaction completed successfully")
}
```

## Lua Scripts

```go
func luaScriptExample(ctx context.Context, client *redis.Client) {
    // Script Lua để increment một key và trả về giá trị mới
    script := redis.NewScript(`
        local key = KEYS[1]
        local increment = tonumber(ARGV[1])
        local current = redis.call('GET', key)
        
        if current == false then
            current = 0
        else
            current = tonumber(current)
        end
        
        local new_value = current + increment
        redis.call('SET', key, new_value)
        
        return new_value
    `)
    
    // Thực thi script
    result, err := script.Run(ctx, client, []string{"lua_counter"}, 5).Result()
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Lua script result: %v", result)
}
```

## Error Handling

```go
func errorHandlingExample(ctx context.Context, client *redis.Client) {
    // Thử get một key không tồn tại
    val, err := client.Get(ctx, "nonexistent_key").Result()
    if err == redis.Nil {
        log.Println("Key không tồn tại")
    } else if err != nil {
        log.Printf("Lỗi Redis: %v", err)
    } else {
        log.Printf("Giá trị: %s", val)
    }
    
    // Timeout handling
    timeoutCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    _, err = client.Get(timeoutCtx, "some_key").Result()
    if err != nil {
        if err == context.DeadlineExceeded {
            log.Println("Timeout khi thực hiện lệnh Redis")
        } else {
            log.Printf("Lỗi khác: %v", err)
        }
    }
}
```

## Best Practices

### 1. Connection Management

```go
func connectionManagement() {
    manager := redis.NewManager()
    
    // Luôn đóng kết nối khi không cần thiết
    defer func() {
        if err := manager.Close(); err != nil {
            log.Printf("Error closing Redis: %v", err)
        }
    }()
    
    // Kiểm tra kết nối trước khi sử dụng
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := manager.Ping(ctx); err != nil {
        log.Fatalf("Redis connection failed: %v", err)
    }
}
```

### 2. Key Naming Convention

```go
const (
    UserPrefix     = "user:"
    SessionPrefix  = "session:"
    CachePrefix    = "cache:"
)

func keyNaming(userID int, sessionID string) {
    client, _ := redis.NewManager().Client()
    ctx := context.Background()
    
    // Sử dụng prefix rõ ràng
    userKey := fmt.Sprintf("%s%d", UserPrefix, userID)
    sessionKey := fmt.Sprintf("%s%s", SessionPrefix, sessionID)
    cacheKey := fmt.Sprintf("%sproduct:%d", CachePrefix, 123)
    
    // TTL cho các loại key khác nhau
    client.Set(ctx, userKey, "user_data", 24*time.Hour)      // User data: 24h
    client.Set(ctx, sessionKey, "session_data", 2*time.Hour) // Session: 2h
    client.Set(ctx, cacheKey, "cached_data", 15*time.Minute) // Cache: 15m
}
```

### 3. Batch Operations

```go
func batchOperations(ctx context.Context, client *redis.Client) {
    // Sử dụng pipeline cho multiple operations
    pipe := client.Pipeline()
    
    // Thêm nhiều operations
    for i := 0; i < 100; i++ {
        key := fmt.Sprintf("batch_key_%d", i)
        value := fmt.Sprintf("value_%d", i)
        pipe.Set(ctx, key, value, time.Hour)
    }
    
    // Thực thi tất cả cùng lúc
    _, err := pipe.Exec(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Batch operations completed")
}
```

### 4. Monitoring và Health Check

```go
func healthCheck(manager redis.Manager) bool {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := manager.Ping(ctx); err != nil {
        log.Printf("Redis health check failed: %v", err)
        return false
    }
    
    return true
}

func monitoringExample(client *redis.Client) {
    ctx := context.Background()
    
    // Đếm số lượng keys
    keyCount, err := client.DBSize(ctx).Result()
    if err == nil {
        log.Printf("Total keys in database: %d", keyCount)
    }
    
    // Kiểm tra memory usage
    memInfo, err := client.Info(ctx, "memory").Result()
    if err == nil {
        log.Printf("Memory info: %s", memInfo)
    }
}
```

## Production Tips

### 1. Configuration cho Production

```yaml
redis:
  client:
    host: "redis.production.com"
    port: 6379
    password: "${REDIS_PASSWORD}"  # Sử dụng environment variable
    db: 0
    timeout: 10
    dial_timeout: 5
    read_timeout: 3
    write_timeout: 3
    pool_size: 20                  # Tăng pool size cho production
    min_idle_conns: 10             # Đảm bảo có sẵn kết nối
```

### 2. Error Recovery

```go
func resilientRedisOperation(client *redis.Client, key, value string) error {
    ctx := context.Background()
    maxRetries := 3
    baseDelay := 100 * time.Millisecond
    
    for i := 0; i < maxRetries; i++ {
        err := client.Set(ctx, key, value, time.Hour).Err()
        if err == nil {
            return nil
        }
        
        // Exponential backoff
        delay := baseDelay * time.Duration(1<<uint(i))
        log.Printf("Redis operation failed (attempt %d/%d): %v. Retrying in %v...", 
                   i+1, maxRetries, err, delay)
        time.Sleep(delay)
    }
    
    return fmt.Errorf("redis operation failed after %d attempts", maxRetries)
}
```

### 3. Graceful Shutdown

```go
func gracefulShutdown(manager redis.Manager) {
    // Đăng ký signal handler
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        log.Println("Shutting down Redis connections...")
        
        if err := manager.Close(); err != nil {
            log.Printf("Error closing Redis: %v", err)
        } else {
            log.Println("Redis connections closed successfully")
        }
        
        os.Exit(0)
    }()
}
```

## Troubleshooting

### Common Issues

1. **Connection timeout**: Kiểm tra network connectivity và Redis server status
2. **Authentication failed**: Xác minh password trong cấu hình
3. **Memory issues**: Monitor Redis memory usage và set appropriate TTL
4. **High latency**: Kiểm tra network, increase timeout values nếu cần

### Debug Commands

```go
func debugCommands(client *redis.Client) {
    ctx := context.Background()
    
    // Kiểm tra connection
    pong, err := client.Ping(ctx).Result()
    log.Printf("PING: %s, Error: %v", pong, err)
    
    // Thông tin server
    info, err := client.Info(ctx).Result()
    log.Printf("INFO: %s, Error: %v", info, err)
    
    // Kiểm tra config
    config, err := client.ConfigGet(ctx, "*").Result()
    log.Printf("CONFIG: %v, Error: %v", config, err)
}
```
