package redis

import (
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// skipIntegrationTest kiểm tra nếu thử nghiệm tích hợp nên được bỏ qua
func skipIntegrationTest(t *testing.T) {
	if os.Getenv("REDIS_INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test. Set REDIS_INTEGRATION_TEST=1 to run.")
	}
}

func TestNewManager(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Client.Enabled = true // Ensure client is enabled for this test
	manager := NewManager(cfg)
	if manager == nil {
		t.Error("Expected manager to be initialized, got nil")
	}

	config := manager.GetConfig()
	if config == nil {
		t.Error("Expected config to be initialized, got nil")
		return // Return early to avoid nil dereference
	}

	if config.Client == nil {
		t.Error("Expected client config to be initialized, got nil")
	}

	if config.Universal == nil {
		t.Error("Expected universal config to be initialized, got nil")
	}
}

func TestNewManager_Custom(t *testing.T) {
	customConfig := DefaultConfig()
	customConfig.Client.Host = "custom-redis-host"
	customConfig.Client.Port = 6380
	customConfig.Client.Password = "custom-password"

	manager := NewManager(customConfig)
	config := manager.GetConfig()

	if config.Client.Host != "custom-redis-host" {
		t.Errorf("Expected host to be 'custom-redis-host', got '%s'", config.Client.Host)
	}
	if config.Client.Port != 6380 {
		t.Errorf("Expected port to be 6380, got %d", config.Client.Port)
	}
	if config.Client.Password != "custom-password" {
		t.Errorf("Expected password to be 'custom-password', got '%s'", config.Client.Password)
	}
}

func TestManager_Client(t *testing.T) {
	config := DefaultConfig()
	config.Client.Enabled = true // Ensure client is enabled for this test
	manager := NewManager(config)

	client, err := manager.Client()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err.Error())
	}

	if client == nil {
		t.Error("Expected client to be initialized, got nil")
	}

}

func TestManager_UniversalClient(t *testing.T) {
	// Create a custom config with Universal.Enabled = true
	config := DefaultConfig()
	config.Universal.Enabled = true
	manager := NewManager(config)

	client, err := manager.UniversalClient()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Error("Expected universal client to be initialized, got nil")
	}

	// Second call should return cached client
	cachedClient, err := manager.UniversalClient()
	if err != nil {
		t.Fatalf("Expected no error on cached universal client, got %v", err)
	}

	if client != cachedClient {
		t.Error("Expected cached universal client to be the same instance")
	}
}

func TestManager_Close(t *testing.T) {
	// Create a manager with Universal.Enabled = true
	config := DefaultConfig()
	config.Universal.Enabled = true
	manager := NewManager(config)

	// Close all connections
	err := manager.Close()
	if err != nil {
		t.Fatalf("Failed to close connections: %v", err)
	}
}

func TestManager_Ping_Integration(t *testing.T) {
	skipIntegrationTest(t)
	cfg := DefaultConfig()
	cfg.Client.Enabled = true // Ensure client is enabled for this test
	manager := NewManager(cfg)
	err := manager.Ping(context.Background())

	if err != nil {
		t.Fatalf("Failed to ping Redis server: %v", err)
	}
}

func TestManager_ClusterPing_Integration(t *testing.T) {
	skipIntegrationTest(t)

	// For this test, you need a Redis cluster or a single Redis instance
	// that works with the universal client
	config := DefaultConfig()
	manager := NewManager(config)

	err := manager.ClusterPing(context.Background())
	if err != nil {
		t.Fatalf("Failed to ping Redis cluster: %v", err)
	}
}

// MockRedisClient tạo một client giả để test
func MockRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func TestManager_Client_ErrorCases(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Client = nil
	manager := NewManager(cfg)
	client, err := manager.Client()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration is missing")

	cfg = DefaultConfig()
	cfg.Client.Enabled = false
	manager = NewManager(cfg)
	client, err = manager.Client()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is disabled")
}

func TestManager_UniversalClient_ErrorCases(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Universal = nil
	manager := NewManager(cfg)
	client, err := manager.UniversalClient()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration is missing")

	cfg = DefaultConfig()
	cfg.Universal.Enabled = false
	manager = NewManager(cfg)
	client, err = manager.UniversalClient()
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is disabled")
}

func TestManager_Close_ErrorCases(t *testing.T) {
	// Đóng khi chưa có client nào không lỗi
	cfg := DefaultConfig()
	manager := NewManager(cfg)
	assert.NoError(t, manager.Close())
}

func TestManager_Close_ClientAndUniversalError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Client.Enabled = true
	cfg.Universal.Enabled = true
	m := NewManager(cfg).(*manager)
	// Tạo client và universal client với địa chỉ không hợp lệ để Close trả về lỗi
	badClient := redis.NewClient(&redis.Options{Addr: "localhost:0"})
	badUniversal := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{"localhost:0"}})
	m.client = badClient
	m.universalClient = badUniversal
	_ = badClient.Close() // Đảm bảo client đã đóng để Close trả về lỗi
	_ = badUniversal.Close()
	err := m.Close()
	assert.Error(t, err)
}

func TestManager_Close_OnlyClientError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Client.Enabled = true
	cfg.Universal.Enabled = true
	m := NewManager(cfg).(*manager)
	badClient := redis.NewClient(&redis.Options{Addr: "localhost:0"})
	goodUniversal := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{"localhost:6379"}})
	m.client = badClient
	m.universalClient = goodUniversal
	_ = badClient.Close()
	err := m.Close()
	assert.Error(t, err)
}

func TestManager_Close_OnlyUniversalError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Client.Enabled = true
	cfg.Universal.Enabled = true
	m := NewManager(cfg).(*manager)
	goodClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	badUniversal := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{"localhost:0"}})
	m.client = goodClient
	m.universalClient = badUniversal
	_ = badUniversal.Close()
	err := m.Close()
	assert.Error(t, err)
}

func TestManager_Ping_ClientPingError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Client.Enabled = true
	m := NewManager(cfg).(*manager)
	badClient := redis.NewClient(&redis.Options{Addr: "localhost:0"})
	m.client = badClient
	err := m.Ping(context.Background())
	assert.Error(t, err)
}

func TestManager_ClusterPing_UniversalPingError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Universal.Enabled = true
	m := NewManager(cfg).(*manager)
	badUniversal := redis.NewUniversalClient(&redis.UniversalOptions{Addrs: []string{"localhost:0"}})
	m.universalClient = badUniversal
	err := m.ClusterPing(context.Background())
	assert.Error(t, err)
}

func TestManager_Ping_ClientNil(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Client.Enabled = true
	m := NewManager(cfg).(*manager)
	// Đặt config.Client = nil để Client() trả về lỗi
	m.config.Client = nil
	err := m.Ping(context.Background())
	assert.Error(t, err)
}

func TestManager_ClusterPing_UniversalNil(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Universal.Enabled = true
	m := NewManager(cfg).(*manager)
	m.universalClient = nil
	err := m.ClusterPing(context.Background())
	assert.Error(t, err)
}

func TestManager_ClusterPing_UniversalConfigNil(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Universal = nil
	m := NewManager(cfg).(*manager)
	err := m.ClusterPing(context.Background())
	assert.Error(t, err)
}
