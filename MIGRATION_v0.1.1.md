# Redis v0.1.1: Hướng dẫn Nâng cấp

## Tổng quan

Phiên bản Redis v0.1.1 tập trung vào việc cải tiến tính ổn định và cập nhật các dependencies. Nó không bao gồm bất kỳ thay đổi API nào đáng kể, nhưng có một vài điều chỉnh để phù hợp với các phiên bản mới của dependencies.

## Thay đổi chính

1. **Cập nhật dependencies**:
   - `go.fork.vn/config` từ v0.1.0 lên v0.1.1
   - `go.fork.vn/di` từ v0.1.0 lên v0.1.1

2. **Cập nhật interfaces**:
   - Các phương thức `Register` và `Boot` trong `ServiceProvider` cần một tham số kiểu `di.Application` thay vì `interface{}`

## Hướng dẫn nâng cấp

### Tự động nâng cấp

Để nâng cấp tự động module của bạn:

```bash
go get go.fork.vn/redis@v0.1.1
```

hoặc cập nhật file go.mod:

```
require go.fork.vn/redis v0.1.1
```

### Thay đổi với ServiceProvider implementations

Nếu bạn đã tạo một triển khai tùy chỉnh của `ServiceProvider`, bạn cần cập nhật các phương thức `Register` và `Boot` để nhận tham số kiểu `di.Application`:

Từ:
```go
func (p *CustomServiceProvider) Register(app interface{}) {
    // Cũ: sử dụng type assertions
    if appWithContainer, ok := app.(interface{ Container() *di.Container }); ok {
        container := appWithContainer.Container()
        // ...
    }
}

func (p *CustomServiceProvider) Boot(app interface{}) {
    // Cũ: sử dụng type assertions
    // ...
}
```

Thành:
```go
func (p *CustomServiceProvider) Register(app di.Application) {
    // Mới: trực tiếp sử dụng app
    container := app.Container()
    // ...
}

func (p *CustomServiceProvider) Boot(app di.Application) {
    // Mới: trực tiếp sử dụng app
    // ...
}
```

## Kiểm tra tương thích

Sau khi nâng cấp lên v0.1.1, hãy chạy các bài kiểm tra và đảm bảo rằng mọi thứ vẫn hoạt động đúng:

```bash
go test ./...
```

## Vấn đề đã biết

Không có vấn đề đã biết nào khi nâng cấp lên phiên bản này.

## Hỗ trợ

Nếu bạn gặp bất kỳ vấn đề nào khi nâng cấp, vui lòng báo cáo tại:
https://github.com/go-fork/redis/issues
