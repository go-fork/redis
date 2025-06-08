# Changelog

All notable changes to the Redis Provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete comprehensive documentation system ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- README.md with detailed usage guide and real examples ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- Documentation system with 4 main documents ([32d785a](https://github.com/go-fork/redis/commit/32d785a)):
  - `docs/index.md` - Official documentation and overview
  - `docs/overview.md` - Architecture and detailed operation principles
  - `docs/configuration.md` - Complete configuration guide with examples
  - `docs/client_universal.md` - Standard vs Universal clients comparison
- Mermaid diagrams for architecture visualization ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- Configuration examples and real use cases ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- Performance benchmarks and monitoring guidelines ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- Migration strategies and troubleshooting guides ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- API examples for all major use cases ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- Framework integration examples (Gin, Echo) ([32d785a](https://github.com/go-fork/redis/commit/32d785a))
- Testing strategies and comprehensive examples ([32d785a](https://github.com/go-fork/redis/commit/32d785a))

### Changed
- Complete restructure according to Go Redis v9.9.0 standards ([6334efa](https://github.com/go-fork/redis/commit/6334efa))

## v0.1.2 - 2025-06-04

### Added
- GitHub workflows for automation (CI, release, update-deps) ([b71677c](https://github.com/go-fork/redis/commit/b71677c))
- Version management structure with releases directory ([b71677c](https://github.com/go-fork/redis/commit/b71677c))
- Automated version management scripts ([b71677c](https://github.com/go-fork/redis/commit/b71677c))
- Fix panic in ServiceProvider Boot method ([b71677c](https://github.com/go-fork/redis/commit/b71677c))

### Changed
- Upgrade `github.com/redis/go-redis/v9` from v9.8.0 to v9.9.0 ([991cf88](https://github.com/go-fork/redis/commit/991cf88))
- Upgrade `go.fork.vn/config` from v0.1.2 to v0.1.3 ([991cf88](https://github.com/go-fork/redis/commit/991cf88))
- Upgrade `go.fork.vn/di` from v0.1.2 to v0.1.3 ([991cf88](https://github.com/go-fork/redis/commit/991cf88))
- Improve documentation and configuration guides ([991cf88](https://github.com/go-fork/redis/commit/991cf88))

## v0.1.1 - 2025-06-02

### Added
- `Enabled` field for ClientConfig and UniversalConfig to control client initialization ([72e9de5](https://github.com/go-fork/redis/commit/72e9de5))
- Improved logic in ServiceProvider to only register enabled clients ([72e9de5](https://github.com/go-fork/redis/commit/72e9de5))
- Configuration validation in provider ([72e9de5](https://github.com/go-fork/redis/commit/72e9de5))
- Test cases to check disabled/enabled clients ([72e9de5](https://github.com/go-fork/redis/commit/72e9de5))

### Changed
- Update `go.fork.vn/config` from v0.1.0 to v0.1.1 ([72e9de5](https://github.com/go-fork/redis/commit/72e9de5))
- Update `go.fork.vn/di` from v0.1.0 to v0.1.1 ([72e9de5](https://github.com/go-fork/redis/commit/72e9de5))
- Update ServiceProvider to match di v0.1.1, change parameter type for Register and Boot from interface{} to di.Application ([72e9de5](https://github.com/go-fork/redis/commit/72e9de5))

## v0.1.0 - 2025-05-31

### Added
- **Redis Client Management**: Comprehensive Redis client management system for Go applications ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Multiple Client Types**: Support for standard Redis client and Universal client (Cluster, Sentinel, standalone) ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **DI Integration**: Seamless integration with Dependency Injection container ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Configuration Support**: Integration with configuration provider for easy setup ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Connection Management**: Advanced connection pool management and configuration ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Error Handling**: Comprehensive error handling and connection reliability ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Testing Support**: Mock implementations and testing utilities ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Performance Optimization**: Optimized connection pooling and resource management ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Security**: Support for authentication and TLS connections ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Monitoring**: Built-in metrics and monitoring capabilities ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Cluster Support**: Full Redis Cluster support with automatic failover ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Sentinel Support**: Redis Sentinel integration for high availability ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Pipeline Support**: Efficient command pipelining for batch operations ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Pub/Sub**: Complete publish/subscribe messaging support ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Lua Scripts**: Support for server-side Lua script execution ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- **Stream Support**: Redis Streams for real-time data processing ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- LICENSE file for MIT license ([2ccf76f](https://github.com/go-fork/redis/commit/2ccf76f))

### Fixed
- Regenerate go.sum with updated dependencies ([8a6413b](https://github.com/go-fork/redis/commit/8a6413b))

### Technical Details
- Initial release as standalone module `go.fork.vn/redis` ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- Repository located at `github.com/go-fork/redis` ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- Built with Go 1.23.9 ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- Full test coverage and documentation included ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- Integration with go-redis/v9 for optimal performance ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- Thread-safe client management ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- Easy mock regeneration with `mockery --name Manager` command ([80c998c](https://github.com/go-fork/redis/commit/80c998c))

### Dependencies
- `github.com/redis/go-redis/v9`: High-performance Redis client ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- `go.fork.vn/di`: Dependency injection integration ([80c998c](https://github.com/go-fork/redis/commit/80c998c))
- `go.fork.vn/config`: Configuration management ([80c998c](https://github.com/go-fork/redis/commit/80c998c))

[Unreleased]: github.com/go-fork/redis/compare/v0.1.2...HEAD
[v0.1.2]: github.com/go-fork/redis/compare/v0.1.1...v0.1.2
[v0.1.1]: github.com/go-fork/redis/compare/v0.1.0...v0.1.1
[v0.1.0]: github.com/go-fork/redis/releases/tag/v0.1.0
