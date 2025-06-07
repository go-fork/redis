# Release Notes - v0.1.2

## Overview
Phiên bản v0.1.2 tập trung vào việc cải thiện xử lý lỗi và tài liệu hướng dẫn.

## What's New
### 🚀 Features
- Thêm tham số `enabled` cho cả Redis Client và Universal Client

### 🐛 Bug Fixes
- Sửa lỗi panic trong phương thức Boot của ServiceProvider
- Cải thiện xử lý lỗi khi cấu hình không hợp lệ

### 🔧 Improvements
- Cập nhật unit test cho phần xử lý lỗi
- Tối ưu quy trình kiểm tra trạng thái kết nối Redis

### 📚 Documentation
- Cập nhật hướng dẫn cấu hình trong README.md và docs/
- Thêm chi tiết về tham số cấu hình mới

## Breaking Changes
### ⚠️ Important Notes
- Không có thay đổi gây ảnh hưởng (non-breaking changes)

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
### Updated
- github.com/redis/go-redis/v9: v9.8.0 → v9.9.0
- go.fork.vn/config: v0.1.2 → v0.1.3
- go.fork.vn/di: v0.1.2 → v0.1.3

### Added
- Thêm scripts tự động hóa cho quản lý phát hành

### Removed
- Không có

## Performance
- Không có thay đổi đáng kể về hiệu năng trong phiên bản này

## Security
- Cập nhật dependencies với các bản vá bảo mật mới nhất

## Testing
- Thêm test case kiểm tra xử lý lỗi panic trong Boot method
- Cải thiện độ phủ test cho các tình huống lỗi và trường hợp ngoại lệ

## Contributors
Cảm ơn tất cả những người đóng góp đã giúp phát hành phiên bản này:
- @go-fork

## Download
- Source code: [go.fork.vn/redis@v0.1.2]
- Documentation: [pkg.go.dev/go.fork.vn/redis@v0.1.2]

---
Release Date: 2025-06-04
