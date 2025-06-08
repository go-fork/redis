# Changelog

## [Unreleased]

### Added
- 📚 **NEW**: Tài liệu hoàn chỉnh và toàn diện cho thư viện Redis Client
- 📚 **NEW**: README.md với hướng dẫn sử dụng chi tiết và ví dụ thực tế
- 📚 **NEW**: Documentation system với 4 tài liệu chính:
  - `docs/index.md` - Tài liệu chính thức và tổng quan
  - `docs/overview.md` - Kiến trúc và nguyên lý hoạt động chi tiết
  - `docs/configuration.md` - Hướng dẫn cấu hình đầy đủ với examples
  - `docs/client_universal.md` - So sánh Standard vs Universal clients
- 📚 **NEW**: Mermaid diagrams cho visualization kiến trúc và workflows
- 📚 **NEW**: Configuration examples và use cases thực tế
- 📚 **NEW**: Performance benchmarks và monitoring guidelines
- 📚 **NEW**: Migration strategies và troubleshooting guides

### Documentation
- 📚 **COMPLETE**: Hoàn thành documentation system với 4 tài liệu chính
- 📚 **ADDED**: API examples cho tất cả major use cases
- 📚 **ADDED**: Architecture diagrams với Mermaid
- 📚 **ADDED**: Configuration templates cho development và production
- 📚 **ADDED**: Performance tuning guidelines
- 📚 **ADDED**: Error handling và debugging guides
- 📚 **ADDED**: Framework integration examples (Gin, Echo)
- 📚 **ADDED**: Testing strategies và examples

### Technical Improvements
- 🏗️ **ARCHITECTURE**: Documented complete system architecture
- ⚡ **PERFORMANCE**: Added performance characteristics và optimization guides
- 🔒 **SECURITY**: TLS/mTLS configuration examples và best practices
- 🔧 **MONITORING**: Observability patterns và metrics collection guides
- 🧪 **TESTING**: Comprehensive testing documentation và strategies

## v0.1.2 - 2025-06-04

### Added
- Thêm thư mục `.github` với các workflow tự động hóa (CI, release, update-deps)
- Thêm cấu trúc quản lý phiên bản với thư mục `releases`
- Thêm scripts tự động hóa quản lý phiên bản trong thư mục `scripts`
- Sửa lỗi panic trong phương thức Boot của ServiceProvider

### Changed
- Nâng cấp `github.com/redis/go-redis/v9` từ v9.8.0 lên v9.9.0
- Nâng cấp `go.fork.vn/config` từ v0.1.2 lên v0.1.3
- Nâng cấp `go.fork.vn/di` từ v0.1.2 lên v0.1.3
- Cải thiện tài liệu và hướng dẫn cấu hình

## v0.1.1 - 2025-06-02

### Changed
- Cập nhật `go.fork.vn/config` từ v0.1.0 lên v0.1.1
- Cập nhật `go.fork.vn/di` từ v0.1.0 lên v0.1.1
- Cập nhật ServiceProvider để phù hợp với di v0.1.1, thay đổi kiểu tham số cho Register và Boot từ interface{} thành di.Application

### Added
- Thêm trường `Enabled` cho ClientConfig và UniversalConfig để kiểm soát việc khởi tạo client
- Cải thiện logic trong ServiceProvider chỉ đăng ký các client đã được kích hoạt
- Bổ sung kiểm tra tính hợp lệ cho cấu hình trong provider
- Thêm test cases để kiểm tra các client bị tắt/bật

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

[Unreleased]: github.com/go-fork/redis/compare/v0.1.2...HEAD
[v0.1.2]: github.com/go-fork/redis/compare/v0.1.1...v0.1.2
[v0.1.1]: github.com/go-fork/redis/compare/v0.1.0...v0.1.1
[v0.1.0]: github.com/go-fork/redis/releases/tag/v0.1.0
