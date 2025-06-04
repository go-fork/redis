# Redis Package v0.1.1 Release Notes

## Summary
Version v0.1.1 của Redis Package cung cấp các cải tiến về tính ổn định, hiệu suất và cập nhật dependencies. Đây là một bản phát hành bảo trì nhỏ không thêm tính năng mới, nhưng cải thiện độ ổn định và khả năng tích hợp với các module khác trong hệ sinh thái Fork.

## Changes

### Dependencies
- Cập nhật `go.fork.vn/config` từ v0.1.0 lên v0.1.1
- Cập nhật `go.fork.vn/di` từ v0.1.0 lên v0.1.1

### Interface Improvements
- Cải tiến ServiceProvider interface để tương thích với di v0.1.1
- Phương thức `Register` và `Boot` trong ServiceProvider giờ nhận tham số kiểu `di.Application` thay vì `interface{}`
- Cải thiện xử lý lỗi và kiểm tra tính hợp lệ trong Provider

### Performance
- Tối ưu hóa connection pooling cho Redis client
- Cải thiện hiệu suất với các batch operations
- Giảm chi phí allocations khi xử lý nhiều thao tác Redis

### Documentation
- Cập nhật tham chiếu phiên bản trong tất cả tài liệu
- Thêm file MIGRATION_v0.1.1.md với hướng dẫn nâng cấp

## Compatibility
Phiên bản này tương thích ngược một phần với v0.1.0. Nếu bạn đã triển khai tùy chỉnh ServiceProvider, bạn cần thực hiện một số thay đổi nhỏ như được mô tả trong MIGRATION_v0.1.1.md. Nếu bạn chỉ sử dụng Redis client, không cần thực hiện bất kỳ thay đổi nào khi nâng cấp.

## Installation
```bash
go get go.fork.vn/redis@v0.1.1
```

hoặc cập nhật trong file go.mod:
```go
require go.fork.vn/redis v0.1.1
```

## Support
Vui lòng báo cáo bất kỳ vấn đề nào tại: github.com/go-fork/redis/issues
