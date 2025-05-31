# Release Notes - Redis v0.1.0

**Release Date:** May 30, 2025  
**Module:** `go.fork.vn/redis`

## 🚀 Major Changes

### Module Migration
- **BREAKING CHANGE**: Migrated module path from `github.com/go-fork/providers/redis` to `go.fork.vn/redis`
- All import statements must be updated to use the new module path
- This marks the transition to the official go.fork.vn domain

## 📚 Documentation

### New Documentation Structure
- **Added** comprehensive documentation in `docs/` folder
- **Added** `docs/redis.md`: Complete technical documentation covering architecture, interfaces, and implementation details
- **Added** `docs/usage.md`: Practical usage guide with examples and best practices
- **Updated** `CHANGELOG.md` with v0.1.0 release information

### Documentation Highlights
- Detailed API reference with all methods of the Manager interface
- Complete configuration guide for standalone, cluster, and sentinel modes
- Comprehensive examples for all Redis data types (String, Hash, List, Set, Sorted Set)
- Advanced features: Pub/Sub, Pipeline, Transactions, Lua Scripts
- Error handling patterns and best practices
- Production deployment guidelines
- Performance optimization tips
- Monitoring and troubleshooting guide

## 🔧 Technical Details

### API Compatibility
- **Maintained** full API compatibility with v0.0.1
- **No breaking changes** to the Manager interface
- **No breaking changes** to ServiceProvider implementation
- All existing code will work without modifications (except import paths)

### Core Features (Unchanged)
- Support for Redis Standard Client and Universal Client
- Automatic connection pooling and management
- Integration with go.fork.vn/di dependency injection
- Configuration support through go.fork.vn/config
- Support for Redis standalone, cluster, and sentinel modes
- Thread-safe operations
- Comprehensive error handling

### Redis Client Support
- **Standard Client**: Optimized for single Redis instance
- **Universal Client**: Supports Cluster, Sentinel, and standalone modes
- **Automatic Failover**: Built-in support for high availability setups
- **Connection Pooling**: Efficient connection management

## 📦 Installation

### New Installation Command
```bash
go get go.fork.vn/redis@v0.1.0
```

### Import Statement
```go
import "go.fork.vn/redis"
```

## 🔄 Migration Guide

### Update Import Statements
**Old:**
```go
import "github.com/go-fork/providers/redis"
```

**New:**
```go
import "go.fork.vn/redis"
```

### Update go.mod
```go
module your-app

go 1.23

require (
    go.fork.vn/redis v0.1.0
    // other dependencies...
)
```

### Find and Replace
Use your IDE's find and replace functionality:
- Find: `github.com/go-fork/providers/redis`
- Replace: `go.fork.vn/redis`

## 📋 Configuration

### Basic Configuration Example
```yaml
redis:
  client:
    host: "127.0.0.1"
    port: 6379
    password: ""
    db: 0
    timeout: 10
    pool_size: 10
```

### Cluster Configuration Example
```yaml
redis:
  universal:
    addresses:
      - "127.0.0.1:7000"
      - "127.0.0.1:7001" 
      - "127.0.0.1:7002"
    cluster_mode: true
    max_redirects: 3
    pool_size: 10
```

## 🎯 Usage Examples

### Basic Usage with DI Container
```go
// Register providers
app.RegisterProviders(
    config.NewServiceProvider(),
    redis.NewServiceProvider(),
)

// Use Redis
manager := container.MustMake("redis.manager").(redis.Manager)
client, err := manager.Client()
if err != nil {
    log.Fatal(err)
}

// Perform Redis operations
err = client.Set(ctx, "key", "value", time.Hour).Err()
```

### Direct Usage
```go
config := redis.DefaultConfig()
config.Client.Host = "localhost"
config.Client.Port = 6379

manager := redis.NewManagerWithConfig(config)
defer manager.Close()

client, err := manager.Client()
// Use client for Redis operations...
```

## 🔗 Resources

- **Technical Documentation**: See `docs/redis.md`
- **Usage Guide**: See `docs/usage.md` 
- **API Reference**: Detailed in `docs/redis.md`
- **Examples**: Available in `docs/usage.md`
- **Changelog**: See `CHANGELOG.md` for complete history

## 🚀 What's Next

### Upcoming Features (Future Releases)
- Redis Streams support
- Advanced monitoring and metrics
- Enhanced connection pooling options
- Performance optimizations
- Additional configuration options

## 🐛 Known Issues

None reported for this release.

## 📊 Dependencies

- **go.fork.vn/di**: ^0.1.0 - Dependency injection container
- **go.fork.vn/config**: ^0.1.0 - Configuration management
- **github.com/redis/go-redis/v9**: ^9.8.0 - Redis client library
- **github.com/stretchr/testify**: ^1.10.0 - Testing framework

## 🔒 Security

### Best Practices
- Use strong passwords for Redis authentication
- Enable TLS for production deployments
- Implement proper network security (VPN, firewall rules)
- Monitor Redis access logs

### Configuration Security
```yaml
redis:
  client:
    password: "${REDIS_PASSWORD}"  # Use environment variables
    # Never commit passwords to version control
```

## 📈 Performance

### Optimizations
- Connection pooling enabled by default
- Configurable timeout settings
- Support for pipeline operations
- Efficient batch operations

### Recommended Settings for Production
```yaml
redis:
  client:
    pool_size: 20
    min_idle_conns: 10
    dial_timeout: 5
    read_timeout: 3
    write_timeout: 3
```

## 🙏 Credits

This release maintains all existing functionality while establishing the foundation for future development under the go.fork.vn domain.

---

**Full Changelog**: https://github.com/go-fork/providers/compare/redis/v0.0.1...redis/v0.1.0
