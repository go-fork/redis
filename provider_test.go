package redis

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/config/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
)

// setupMockApplication thiết lập một mock Application với Container đã cấu hình
func setupMockApplication(t *testing.T) (*diMocks.MockApplication, di.Container) {
	container := di.New()

	mockApp := diMocks.NewMockApplication(t)
	mockApp.On("Container").Return(container).Maybe()

	return mockApp, container
}

func TestNewServiceProvider(t *testing.T) {
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "NewServiceProvider() không được trả về nil")
}

func TestServiceProviderRegister(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
		var config *Config
		// Handle both *Config and **Config
		if c, ok := args.Get(1).(**Config); ok {
			config = *c
		} else if c, ok := args.Get(1).(*Config); ok {
			config = c
		}
		if config != nil {
			config.Client = &ClientConfig{
				Enabled:      true,
				Host:         "localhost",
				Port:         6379,
				Password:     "",
				DB:           0,
				Prefix:       "test:",
				Timeout:      5,
				DialTimeout:  5,
				ReadTimeout:  3,
				WriteTimeout: 3,
				PoolSize:     10,
				MinIdleConns: 2,
			}
		}
	}).Return(nil).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider với application
	provider.Register(mockApp)

	// Kiểm tra binding "redis"
	managerInstance, err := container.Make("redis")
	assert.NoError(t, err, "ServiceProvider phải đăng ký binding 'redis'")

	manager, ok := managerInstance.(Manager)
	assert.True(t, ok, "Binding 'redis' phải là kiểu Manager, nhưng nhận được %T", managerInstance)

	// Dọn dẹp
	if err := manager.Close(); err != nil {
		t.Logf("Không thể đóng manager: %v", err)
	}
}

func TestServiceProviderBoot(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func() di.Application
		expectPanic bool
	}{
		{
			name: "valid application with redis binding",
			setupMocks: func() di.Application {
				mockApp, container := setupMockApplication(t)
				container.Instance("redis", NewManager(&Config{
					Client: &ClientConfig{
						Enabled: true,
						Host:    "localhost",
						Port:    6379,
					},
				}))
				return mockApp
			},
			expectPanic: false,
		},
		{
			name: "nil application",
			setupMocks: func() di.Application {
				return nil
			},
			expectPanic: false,
		},
		{
			name: "application with nil container",
			setupMocks: func() di.Application {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.On("Container").Return(nil).Maybe()
				return mockApp
			},
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewServiceProvider()
			app := tt.setupMocks()

			if tt.expectPanic {
				assert.Panics(t, func() {
					provider.Boot(app)
				})
			} else {
				assert.NotPanics(t, func() {
					provider.Boot(app)
				})
			}
		})
	}
}

func TestServiceProviderWithConfigError(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager với lỗi
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Return(errors.New("config error")).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Register nên panic khi config manager trả về lỗi
	assert.Panics(t, func() {
		provider.Register(mockApp)
	}, "ServiceProvider.Register nên panic khi config manager trả về lỗi")
}

func TestServiceProviderWithInvalidConfig(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager với cấu hình không hợp lệ
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
		var config *Config
		if c, ok := args.Get(1).(**Config); ok {
			config = *c
		} else if c, ok := args.Get(1).(*Config); ok {
			config = c
		}
		if config != nil {
			config.Client = &ClientConfig{
				Enabled: true,
				Host:    "localhost",
				Port:    -1, // Port không hợp lệ
			}
		}
	}).Return(nil).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Không mong đợi panic vì provider không validate config
	assert.NotPanics(t, func() {
		provider.Register(mockApp)
	}, "ServiceProvider.Register không nên panic khi cấu hình không hợp lệ nếu không có validation")
}

func TestContainerBindingResolution(t *testing.T) {
	// Tạo mock application và container
	mockApp, container := setupMockApplication(t)

	// Tạo mock config manager
	mockConfigManager := mocks.NewMockManager(t)
	mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
		var config *Config
		if c, ok := args.Get(1).(**Config); ok {
			config = *c
		} else if c, ok := args.Get(1).(*Config); ok {
			config = c
		}
		if config != nil {
			config.Client = &ClientConfig{
				Enabled:      true,
				Host:         "localhost",
				Port:         6379,
				Password:     "",
				DB:           0,
				Prefix:       "test:",
				Timeout:      5,
				DialTimeout:  5,
				ReadTimeout:  3,
				WriteTimeout: 3,
				PoolSize:     10,
				MinIdleConns: 2,
			}
		}
	}).Return(nil).Once()

	// Đăng ký config manager vào container
	container.Instance("config", mockConfigManager)

	// Tạo service provider
	provider := NewServiceProvider()

	// Đăng ký provider
	provider.Register(mockApp)

	// Thêm một binding phụ thuộc vào log manager
	container.Bind("custom.logger", func(c di.Container) interface{} {
		// Lấy log manager từ container
		manager, err := c.Make("redis")
		if err != nil {
			t.Fatal("Không thể resolve dependency 'redis':", err)
		}

		// Trả về một struct sử dụng log manager
		return struct {
			LogManager Manager
			Name       string
		}{
			LogManager: manager.(Manager),
			Name:       "CustomLogger",
		}
	})

	// Giải quyết binding
	customLogger, err := container.Make("custom.logger")
	assert.NoError(t, err, "Phải resolve binding 'custom.logger' thành công")

	// Kiểm tra cấu trúc được trả về
	loggerStruct, ok := customLogger.(struct {
		LogManager Manager
		Name       string
	})

	assert.True(t, ok, "Binding 'custom.logger' phải trả về kiểu đúng, nhưng nhận được: %T", customLogger)
	assert.Equal(t, "CustomLogger", loggerStruct.Name, "Tên phải đúng")
	assert.NotNil(t, loggerStruct.LogManager, "LogManager không được là nil")
}

// TestServiceProviderRequires kiểm tra method Requires() trả về giá trị đúng
func TestServiceProviderRequires(t *testing.T) {
	provider := NewServiceProvider()
	requires := provider.Requires()
	// Redis provider phụ thuộc vào config provider
	assert.Equal(t, []string{"config"}, requires, "Redis provider phải phụ thuộc vào provider 'config'")
}

func TestServiceProviderProviders(t *testing.T) {
	provider := NewServiceProvider()
	providers := provider.Providers()
	// Nếu provider không đăng ký service nào, mong đợi rỗng
	assert.Empty(t, providers, "Provider không nên đăng ký service nào nếu Providers trả về rỗng")
}

// TestRegisterWithInvalidInputs kiểm tra các trường hợp đầu vào không hợp lệ cho Register
func TestRegisterWithInvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func() (di.Application, di.Container)
		expectPanic bool
		description string
	}{
		{
			name: "nil application",
			setupMocks: func() (di.Application, di.Container) {
				return nil, nil
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi app là nil",
		},
		{
			name: "application with nil container",
			setupMocks: func() (di.Application, di.Container) {
				mockApp := diMocks.NewMockApplication(t)
				mockApp.On("Container").Return(nil).Once()
				return mockApp, nil
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi container là nil",
		},
		{
			name: "container without config manager",
			setupMocks: func() (di.Application, di.Container) {
				mockApp, container := setupMockApplication(t)
				return mockApp, container
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi config manager không tồn tại",
		},
		{
			name: "container with invalid config manager type",
			setupMocks: func() (di.Application, di.Container) {
				mockApp, container := setupMockApplication(t)
				container.Instance("config", "not-a-config-manager")
				return mockApp, container
			},
			expectPanic: true,
			description: "ServiceProvider.Register nên panic khi config manager có kiểu không đúng",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewServiceProvider()
			app, _ := tt.setupMocks()

			if tt.expectPanic {
				assert.Panics(t, func() {
					provider.Register(app)
				}, tt.description)
			} else {
				assert.NotPanics(t, func() {
					provider.Register(app)
				}, tt.description)
			}
		})
	}
}

// BenchmarkNewServiceProvider đo hiệu suất tạo ServiceProvider mới
func BenchmarkNewServiceProvider(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider := NewServiceProvider()
		_ = provider
	}
}

// BenchmarkServiceProviderRegister đo hiệu suất đăng ký ServiceProvider
func BenchmarkServiceProviderRegister(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Setup cho mỗi iteration để tránh mock conflicts
		container := di.New()

		// Tạo mock application với clean state
		mockApp := diMocks.NewMockApplication(b)
		mockApp.On("Container").Return(container).Once()

		// Tạo mock config manager với clean state
		mockConfigManager := mocks.NewMockManager(b)
		mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
			var config *Config
			if c, ok := args.Get(1).(**Config); ok {
				config = *c
			} else if c, ok := args.Get(1).(*Config); ok {
				config = c
			}
			if config != nil {
				config.Client = &ClientConfig{
					Enabled:      true,
					Host:         "localhost",
					Port:         6379,
					Password:     "",
					DB:           0,
					Prefix:       "test:",
					Timeout:      5,
					DialTimeout:  5,
					ReadTimeout:  3,
					WriteTimeout: 3,
					PoolSize:     10,
					MinIdleConns: 2,
				}
			}
		}).Return(nil).Once()

		container.Instance("config", mockConfigManager)
		provider := NewServiceProvider()
		b.StartTimer()

		provider.Register(mockApp)

		// Cleanup để tránh memory leak
		b.StopTimer()
		if manager, err := container.Make("redis"); err == nil {
			if logManager, ok := manager.(Manager); ok {
				logManager.Close()
			}
		}
	}
}

// BenchmarkServiceProviderBoot đo hiệu suất Boot method
func BenchmarkServiceProviderBoot(b *testing.B) {
	// Setup một lần cho Boot benchmark vì Boot không có side effects
	container := di.New()
	mockApp := &diMocks.MockApplication{}
	mockApp.On("Container").Return(container).Maybe()
	container.Instance("redis", NewManager(&Config{
		Client: &ClientConfig{
			Enabled: true,
			Host:    "localhost",
			Port:    6379,
		},
	}))

	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.Boot(mockApp)
	}
}

// BenchmarkServiceProviderRequires đo hiệu suất Requires method
func BenchmarkServiceProviderRequires(b *testing.B) {
	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.Requires()
	}
}

// BenchmarkServiceProviderProviders đo hiệu suất Providers method
func BenchmarkServiceProviderProviders(b *testing.B) {
	provider := NewServiceProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.Providers()
	}
}

// BenchmarkContainerMakeLog đo hiệu suất resolve log service từ container
func BenchmarkContainerMakeLog(b *testing.B) {
	// Setup một lần
	container := di.New()
	mockApp := &diMocks.MockApplication{}
	mockApp.On("Container").Return(container).Once()

	mockConfigManager := &mocks.MockManager{}
	mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
		var config *Config
		if c, ok := args.Get(1).(**Config); ok {
			config = *c
		} else if c, ok := args.Get(1).(*Config); ok {
			config = c
		}
		if config != nil {
			config.Client = &ClientConfig{
				Enabled:      true,
				Host:         "localhost",
				Port:         6379,
				Password:     "",
				DB:           0,
				Prefix:       "test:",
				Timeout:      5,
				DialTimeout:  5,
				ReadTimeout:  3,
				WriteTimeout: 3,
				PoolSize:     10,
				MinIdleConns: 2,
			}
		}
	}).Return(nil).Once()

	container.Instance("config", mockConfigManager)

	provider := NewServiceProvider()
	provider.Register(mockApp)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager, err := container.Make("redis")
		if err != nil {
			b.Fatal("Failed to make log service:", err)
		}
		_ = manager
	}

	// Cleanup after benchmark
	b.StopTimer()
	if manager, err := container.Make("redis"); err == nil {
		if logManager, ok := manager.(Manager); ok {
			logManager.Close()
		}
	}
}

// BenchmarkCompleteServiceProviderWorkflow đo hiệu suất toàn bộ workflow
func BenchmarkCompleteServiceProviderWorkflow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Setup cho mỗi iteration
		container := di.New()
		mockApp := diMocks.NewMockApplication(b)
		mockApp.On("Container").Return(container).Times(2) // Register + Boot

		mockConfigManager := mocks.NewMockManager(b)
		mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
			var config *Config
			if c, ok := args.Get(1).(**Config); ok {
				config = *c
			} else if c, ok := args.Get(1).(*Config); ok {
				config = c
			}
			if config != nil {
				config.Client = &ClientConfig{
					Enabled:      true,
					Host:         "localhost",
					Port:         6379,
					Password:     "",
					DB:           0,
					Prefix:       "test:",
					Timeout:      5,
					DialTimeout:  5,
					ReadTimeout:  3,
					WriteTimeout: 3,
					PoolSize:     10,
					MinIdleConns: 2,
				}
			}
		}).Return(nil).Once()

		container.Instance("config", mockConfigManager)
		provider := NewServiceProvider()
		b.StartTimer()

		// Toàn bộ workflow: Register -> Boot -> Make
		provider.Register(mockApp)
		provider.Boot(mockApp)
		manager, err := container.Make("redis")
		if err != nil {
			b.Fatal("Failed to make log service:", err)
		}
		_ = manager

		// Cleanup
		b.StopTimer()
		if logManager, ok := manager.(Manager); ok {
			logManager.Close()
		}
	}
}

// BenchmarkParallelServiceProviderRegister đo hiệu suất với concurrent access
func BenchmarkParallelServiceProviderRegister(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Setup cho mỗi goroutine với clean mocks
			container := di.New()
			mockApp := diMocks.NewMockApplication(b)
			mockApp.On("Container").Return(container).Once()

			mockConfigManager := mocks.NewMockManager(b)
			mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
				var config *Config
				if c, ok := args.Get(1).(**Config); ok {
					config = *c
				} else if c, ok := args.Get(1).(*Config); ok {
					config = c
				}
				if config != nil {
					config.Client = &ClientConfig{
						Enabled:      true,
						Host:         "localhost",
						Port:         6379,
						Password:     "",
						DB:           0,
						Prefix:       "test:",
						Timeout:      5,
						DialTimeout:  5,
						ReadTimeout:  3,
						WriteTimeout: 3,
						PoolSize:     10,
						MinIdleConns: 2,
					}
				}
			}).Return(nil).Once()

			container.Instance("config", mockConfigManager)

			provider := NewServiceProvider()
			provider.Register(mockApp)

			// Cleanup
			if manager, err := container.Make("redis"); err == nil {
				if logManager, ok := manager.(Manager); ok {
					logManager.Close()
				}
			}
		}
	})
}

// BenchmarkServiceProviderWithDifferentLogLevels đo hiệu suất với các log level khác nhau
func BenchmarkServiceProviderWithDifferentLogLevels(b *testing.B) {
	logLevels := []string{"debug", "info", "warning", "error", "fatal"}

	for _, level := range logLevels {
		b.Run("Level_"+level, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				container := di.New()
				mockApp := diMocks.NewMockApplication(b)
				mockApp.On("Container").Return(container).Once()

				mockConfigManager := mocks.NewMockManager(b)
				mockConfigManager.On("UnmarshalKey", "redis", mock.Anything).Run(func(args mock.Arguments) {
					var config *Config
					if c, ok := args.Get(1).(**Config); ok {
						config = *c
					} else if c, ok := args.Get(1).(*Config); ok {
						config = c
					}
					if config != nil {
						config.Client = &ClientConfig{
							Enabled:      true,
							Host:         "localhost",
							Port:         6379,
							Password:     "",
							DB:           0,
							Prefix:       "test:",
							Timeout:      5,
							DialTimeout:  5,
							ReadTimeout:  3,
							WriteTimeout: 3,
							PoolSize:     10,
							MinIdleConns: 2,
						}
					}
				}).Return(nil).Once()

				container.Instance("config", mockConfigManager)
				provider := NewServiceProvider()
				b.StartTimer()

				provider.Register(mockApp)

				// Cleanup
				b.StopTimer()
				if manager, err := container.Make("redis"); err == nil {
					if logManager, ok := manager.(Manager); ok {
						logManager.Close()
					}
				}
			}
		})
	}
}
