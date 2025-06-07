package redis

import (
	"fmt"

	"go.fork.vn/config"
	"go.fork.vn/di"
)

// ServiceProvider định nghĩa interface cho Redis service provider.
//
// ServiceProvider kế thừa từ di.ServiceProvider và định nghĩa
// các phương thức cần thiết cho một Redis service provider.
type ServiceProvider interface {
	di.ServiceProvider
}

// serviceProvider là implementation của ServiceProvider.
//
// serviceProvider chịu trách nhiệm đăng ký các dịch vụ Redis vào DI container
// và cung cấp Redis manager cùng các client instances cho các module khác trong ứng dụng.
type serviceProvider struct {
	config    *Config  // Cấu hình Redis đã được validate
	manager   *Manager // Redis manager để quản lý các client connections
	providers []string // Danh sách các service đã được đăng ký
}

// NewServiceProvider tạo một instance mới của Redis service provider.
//
// ServiceProvider này sẽ đăng ký các dịch vụ Redis vào DI container
// và quản lý vòng đời của Redis connections.
//
// Returns:
//   - ServiceProvider: Instance mới của Redis service provider
//
// Example:
//
//	provider := NewServiceProvider()
//	app.RegisterProvider(provider)
func NewServiceProvider() ServiceProvider {
	return &serviceProvider{}
}

// Register đăng ký các dependency cần thiết cho Redis provider.
//
// Phương thức này được gọi trong giai đoạn đăng ký services của application lifecycle.
// Nó sẽ:
// 1. Lấy config manager từ DI container
// 2. Đọc và validate Redis configuration
// 3. Khởi tạo Redis manager nếu Redis được enabled
//
// Parameters:
//   - app: Application instance chứa DI container
//
// Panics:
//   - Khi app parameter là nil
//   - Khi config service chưa được đăng ký
//   - Khi không thể parse Redis configuration
//   - Khi Redis configuration không hợp lệ
func (p *serviceProvider) Register(app di.Application) {
	// Kiểm tra tính hợp lệ của app parameter
	if app == nil {
		panic("Application cannot be nil in Redis provider Register")
	}

	// Lấy container từ app
	container := app.Container()

	// Kiểm tra và lấy config manager từ container
	configManager, ok := container.MustMake("config").(config.Manager)
	if !ok {
		panic("Config service unregistered")
	}

	// Tạo cấu hình mặc định và đọc cấu hình từ config
	redisConfig := DefaultConfig()
	err := configManager.UnmarshalKey("redis", redisConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal redis config: %v", err))
	}

	// Validate cấu hình Redis trước khi sử dụng
	if err := redisConfig.Validate(); err != nil {
		panic(fmt.Sprintf("Invalid redis config: %v", err))
	}

	// Lưu cấu hình đã validate
	p.config = redisConfig

	// Nếu Redis được bật, khởi tạo manager
	if p.config.Enabled() {
		manager := NewManager(*p.config)
		p.manager = &manager
		// Thêm service redis vào danh sách providers
		p.providers = append(p.providers, "redis")
	}
}

// Boot khởi động và đăng ký các Redis clients vào DI container.
//
// Phương thức này được gọi trong giai đoạn boot của application lifecycle,
// sau khi tất cả providers đã được đăng ký. Nó sẽ:
// 1. Khởi tạo Redis client nếu được enabled
// 2. Khởi tạo Universal client nếu được enabled
// 3. Đăng ký các clients vào DI container
//
// Application sẽ panic nếu Redis được config enabled nhưng không thể kết nối.
// Điều này đảm bảo fail-fast behavior và tránh runtime errors sau này.
//
// Parameters:
//   - app: Application instance chứa DI container
//
// Panics:
//   - Khi app parameter là nil
//   - Khi không thể tạo Redis client (nếu client enabled)
//   - Khi không thể tạo Universal client (nếu universal enabled)
func (p *serviceProvider) Boot(app di.Application) {
	// Kiểm tra tính hợp lệ của app parameter
	if app == nil {
		panic("Application cannot be nil in Redis provider Boot")
	}

	// Lấy container từ application
	c := app.Container()

	// Chỉ tiến hành nếu config và manager đã được khởi tạo
	if p.config != nil && p.manager != nil {
		// Khởi tạo Redis client nếu được bật trong config
		if p.config.Client.Enabled {
			client, err := (*p.manager).Client()
			if err != nil {
				panic(fmt.Sprintf("Redis Client connection failed: %v", err))
			}
			// Đăng ký client instance vào container
			c.Instance("redis.client", *client)
			p.providers = append(p.providers, "redis.client")
		}

		// Khởi tạo Universal client nếu được bật trong config
		if p.config.Universal.Enabled {
			universal, err := (*p.manager).UniversalClient()
			if err != nil {
				panic(fmt.Sprintf("Redis Universal client connection failed: %v", err))
			}
			// Đăng ký universal client instance vào container
			c.Instance("redis.universal", *universal)
			p.providers = append(p.providers, "redis.universal")
		}
	}
}

// Providers trả về danh sách các service mà provider này đăng ký.
//
// Phương thức này được DI container sử dụng để biết provider này
// đăng ký những service nào. Danh sách có thể bao gồm:
// - "redis": Redis manager service
// - "redis.client": Redis client instance
// - "redis.universal": Redis universal client instance
//
// Returns:
//   - []string: Mảng chứa tên của các service được đăng ký
func (p *serviceProvider) Providers() []string {
	return p.providers
}

// Requires trả về danh sách các dependency mà Redis provider phụ thuộc.
//
// Redis provider phụ thuộc vào config provider để đọc cấu hình Redis
// từ file configuration. Config provider phải được đăng ký trước
// khi Redis provider được khởi tạo.
//
// Returns:
//   - []string: Danh sách các service provider khác mà provider này yêu cầu
func (p *serviceProvider) Requires() []string {
	return []string{
		"config",
	}
}
