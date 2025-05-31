# Changelog

## [Unreleased]

## v0.1.0 - 2025-05-31

### Added
- **Redis Client Management**: Comprehensive Redis client management system for Go applications
- **Multiple Client Types**: Support for standard Redis client and Universal client (Cluster, Sentinel, standalone)
- **DI Integration**: Seamless integration with Dependency Injection container
- **Configuration Support**: Integration with configuration provider for easy setup
- **Connection Management**: Advanced connection pool management and configuration
- **Error Handling**: Comprehensive error handling and connection reliability
- **Testing Support**: Mock implementations and testing utilities
- **Performance Optimization**: Optimized connection pooling and resource management
- **Security**: Support for authentication and TLS connections
- **Monitoring**: Built-in metrics and monitoring capabilities
- **Cluster Support**: Full Redis Cluster support with automatic failover
- **Sentinel Support**: Redis Sentinel integration for high availability
- **Pipeline Support**: Efficient command pipelining for batch operations
- **Pub/Sub**: Complete publish/subscribe messaging support
- **Lua Scripts**: Support for server-side Lua script execution
- **Stream Support**: Redis Streams for real-time data processing

### Technical Details
- Initial release as standalone module `go.fork.vn/redis`
- Repository located at `github.com/go-fork/redis`
- Built with Go 1.23.9
- Full test coverage and documentation included
- Integration with go-redis/v9 for optimal performance
- Thread-safe client management
- Easy mock regeneration with `mockery --name Manager` command

### Dependencies
- `github.com/redis/go-redis/v9`: High-performance Redis client
- `go.fork.vn/di`: Dependency injection integration
- `go.fork.vn/config`: Configuration management

[Unreleased]: https://github.com/go-fork/redis/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/go-fork/redis/releases/tag/v0.1.0
