# Redis Provider v0.1.2

Redis Provider là một package cung cấp tích hợp Redis hoàn chỉnh cho Fork Framework, hỗ trợ Redis standalone, Cluster và Sentinel.

## Cài đặt nhanh

```bash
go get go.fork.vn/redis@v0.1.2
```

## Sử dụng cơ bản

```go
// main.go
package main

import (
    "context"
    "log"
    
    "go.fork.vn/core"
    "go.fork.vn/config"
    "go.fork.vn/redis"
)

func main() {
    app := core.NewApplication()
    
    // Đăng ký providers
    app.RegisterProviders(
        config.NewServiceProvider(),
        redis.NewServiceProvider(),
    )
    
    app.Boot()
    
    // Sử dụng Redis
    var manager redis.Manager
    app.Container().Make("redis.manager", &manager)
    
    conn := manager.Connection()
    
    ctx := context.Background()
    conn.Set(ctx, "welcome", "Hello Fork Redis!", 0)
    
    val, err := conn.Get(ctx, "welcome").Result()
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Value: %s", val)
}
```

## Cấu hình nhanh

```yaml
# config/redis.yaml
redis:
  default: "redis"
  connections:
    redis:
      driver: "redis"
      host: "127.0.0.1"
      port: 6379
      database: 0
      timeout: "10s"
```

## Tính năng chính

- ✅ **Đa kết nối**: Hỗ trợ multiple Redis connections
- ✅ **Cluster & Sentinel**: Hỗ trợ Redis Cluster và Sentinel
- ✅ **Connection Pooling**: Quản lý pool kết nối tự động
- ✅ **DI Integration**: Tích hợp hoàn chỉnh với DI container
- ✅ **Health Checks**: Kiểm tra sức khỏe kết nối
- ✅ **TLS Support**: Hỗ trợ kết nối bảo mật
- ✅ **Testing**: Mock interfaces cho unit testing

## Tài liệu chi tiết

- [Overview](overview.md) - Kiến trúc và thiết kế
- [Reference](reference.md) - API Reference đầy đủ  
- [Usage](usage.md) - Hướng dẫn sử dụng chi tiết

## Phiên bản & Compatibility

- **Version**: v0.1.1
- **Go**: 1.23.9+
- **Dependencies**: 
  - go.fork.vn/di v0.1.2
  - go.fork.vn/config v0.1.2
  - github.com/redis/go-redis/v9

## Hỗ trợ

- [GitHub Issues](github.com/go-fork/redis/issues)
- [Examples](github.com/go-fork/recipes/tree/main/examples/redis)
- [Community](github.com/go-fork/redis/discussions)
