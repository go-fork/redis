package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrClientNotEnabled được trả về khi Redis client không được bật trong cấu hình.
	ErrClientNotEnabled = errors.New("redis client is not enabled")
	// ErrUniversalNotEnabled được trả về khi Redis universal client không được bật trong cấu hình.
	ErrUniversalNotEnabled = errors.New("redis universal client is not enabled")
)

// Manager định nghĩa interface để quản lý các kết nối Redis.
//
// Manager cung cấp các phương thức để tạo và quản lý Redis client connections,
// bao gồm cả standard client và universal client (hỗ trợ cluster/sentinel).
// Manager đảm bảo rằng connections được tái sử dụng và quản lý lifecycle một cách hiệu quả.
type Manager interface {
	// Client trả về Redis client instance để thực hiện các thao tác Redis.
	//
	// Phương thức này tạo một kết nối Redis standard client dựa trên cấu hình
	// được cung cấp khi khởi tạo Manager. Client được cache để tái sử dụng
	// trong các lần gọi tiếp theo.
	//
	// Returns:
	//   - *redis.Client: Redis client instance nếu thành công
	//   - error: Lỗi nếu không thể tạo hoặc kết nối client
	//
	// Exceptions:
	//   - Trả về error ErrClientNotEnabled nếu client không được bật trong cấu hình
	//   - Trả về error nếu không thể ping Redis server
	//   - Trả về error nếu cấu hình không hợp lệ
	Client() (*redis.Client, error)

	// ClientPing kiểm tra kết nối của Redis client.
	//
	// Phương thức này thực hiện lệnh PING để xác minh rằng kết nối
	// đến Redis server vẫn hoạt động bình thường.
	//
	// Parameters:
	//   - ctx: Context để control timeout và cancellation của operation
	//
	// Returns:
	//   - error: nil nếu ping thành công, ngược lại trả về lỗi
	//
	// Exceptions:
	//   - Trả về error nếu không thể lấy client
	//   - Trả về error nếu ping command thất bại
	//   - Trả về error nếu context bị timeout hoặc cancelled
	ClientPing(ctx context.Context) error

	// UniversalClient trả về Redis universal client instance.
	//
	// Universal client hỗ trợ cả single-node, cluster và sentinel deployments.
	// Client được cache để tái sử dụng trong các lần gọi tiếp theo.
	//
	// Returns:
	//   - *redis.UniversalClient: Universal client instance nếu thành công
	//   - error: Lỗi nếu không thể tạo hoặc kết nối client
	//
	// Exceptions:
	//   - Trả về error ErrUniversalNotEnabled nếu universal client không được bật trong cấu hình
	//   - Trả về error nếu không thể ping Redis server
	//   - Trả về error nếu cấu hình không hợp lệ
	UniversalClient() (*redis.UniversalClient, error)

	// UniversalPing kiểm tra kết nối của Redis universal client.
	//
	// Phương thức này thực hiện lệnh PING để xác minh rằng kết nối
	// đến Redis cluster/sentinel vẫn hoạt động bình thường.
	//
	// Parameters:
	//   - ctx: Context để control timeout và cancellation của operation
	//
	// Returns:
	//   - error: nil nếu ping thành công, ngược lại trả về lỗi
	//
	// Exceptions:
	//   - Trả về error nếu không thể lấy universal client
	//   - Trả về error nếu ping command thất bại
	//   - Trả về error nếu context bị timeout hoặc cancelled
	UniversalPing(ctx context.Context) error

	// Close đóng tất cả các kết nối Redis đang hoạt động.
	//
	// Phương thức này đảm bảo rằng tất cả resources được giải phóng
	// một cách an toàn và các connections được đóng đúng cách.
	// Sau khi gọi Close(), Manager không thể sử dụng lại được.
	//
	// Returns:
	//   - error: nil nếu đóng thành công, ngược lại trả về lỗi đầu tiên gặp phải
	//
	// Exceptions:
	//   - Trả về error nếu việc đóng client connection thất bại
	//   - Trả về error nếu việc đóng universal client connection thất bại
	Close() error
}

// manager là implementation của Manager interface.
//
// manager chịu trách nhiệm quản lý các kết nối Redis và cung cấp
// các client instances cho việc thực hiện các thao tác Redis.
// Struct này cache các connections để tái sử dụng và đảm bảo
// rằng các resources được quản lý hiệu quả.
type manager struct {
	config    Config                 // Cấu hình Redis cho manager
	client    *redis.Client          // Cached Redis standard client
	universal *redis.UniversalClient // Cached Redis universal client
}

// NewManager tạo một Manager instance mới với cấu hình được cung cấp.
//
// Phương thức này khởi tạo một manager mới với Config được truyền vào.
// Manager sẽ sử dụng config này để tạo các Redis connections khi cần thiết.
//
// Parameters:
//   - config: Cấu hình Redis chứa thông tin client và universal settings
//
// Returns:
//   - Manager: Instance của Manager interface
//
// Examples:
//
//	config := DefaultConfig()
//	manager := NewManager(config)
//	defer manager.Close()
func NewManager(config Config) Manager {
	return &manager{
		config: config,
	}
}

// Client trả về Redis client instance để thực hiện các thao tác Redis.
//
// Phương thức này tạo một kết nối Redis standard client dựa trên cấu hình
// được cung cấp khi khởi tạo Manager. Client được cache để tái sử dụng
// trong các lần gọi tiếp theo nếu connection thành công.
//
// Returns:
//   - *redis.Client: Redis client instance nếu thành công
//   - error: Lỗi nếu không thể tạo hoặc kết nối client
//
// Exceptions:
//   - Trả về (nil, nil) nếu client không được bật trong cấu hình
//   - Trả về error nếu không thể ping Redis server trong 2 giây
//   - Trả về error nếu cấu hình không hợp lệ
func (m *manager) Client() (*redis.Client, error) {
	// Trả về cached client nếu đã có sẵn
	if m.client != nil {
		return m.client, nil
	}

	// Kiểm tra xem client có được bật trong cấu hình không
	if !m.config.Client.Enabled {
		return nil, ErrClientNotEnabled
	}

	// Tạo Redis client mới với các tùy chọn cấu hình
	client := redis.NewClient(m.config.Options())

	// Kiểm tra kết nối với timeout 2 giây
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Ping để xác minh kết nối
	if err := client.Ping(ctx).Err(); err != nil {
		// Đóng client ngay lập tức và không cache khi ping thất bại
		client.Close()
		return nil, err
	}

	// Cache client chỉ khi connection thành công
	m.client = client
	return m.client, nil
}

// ClientPing kiểm tra kết nối của Redis client.
//
// Phương thức này thực hiện lệnh PING để xác minh rằng kết nối
// đến Redis server vẫn hoạt động bình thường. Nó sử dụng Client()
// để lấy instance và thực hiện ping command.
//
// Parameters:
//   - ctx: Context để control timeout và cancellation của operation
//
// Returns:
//   - error: nil nếu ping thành công, ngược lại trả về lỗi
//
// Exceptions:
//   - Trả về error nếu không thể lấy client từ Client()
//   - Trả về error nếu ping command thất bại
//   - Trả về error nếu context bị timeout hoặc cancelled
func (m *manager) ClientPing(ctx context.Context) error {
	// Lấy client instance
	client, err := m.Client()
	if err != nil {
		return err
	}

	// Thực hiện lệnh ping
	return client.Ping(ctx).Err()
}

// UniversalClient trả về Redis universal client instance.
//
// Universal client hỗ trợ cả single-node, cluster và sentinel deployments.
// Client được cache để tái sử dụng trong các lần gọi tiếp theo.
// Phương thức này kiểm tra cache trước khi tạo connection mới.
//
// Returns:
//   - *redis.UniversalClient: Universal client instance nếu thành công
//   - error: Lỗi nếu không thể tạo hoặc kết nối client
//
// Exceptions:
//   - Trả về (nil, ErrUniversalNotEnabled) nếu universal client không được bật trong cấu hình
//   - Trả về error nếu không thể ping Redis server trong 2 giây
//   - Trả về error nếu cấu hình không hợp lệ
func (m *manager) UniversalClient() (*redis.UniversalClient, error) {
	// Trả về cached universal client nếu đã có sẵn
	if m.universal != nil {
		return m.universal, nil
	}

	// Kiểm tra xem universal client có được bật trong cấu hình không
	if !m.config.Universal.Enabled {
		return nil, ErrUniversalNotEnabled
	}

	// Tạo universal client mới với các tùy chọn cấu hình
	universal := redis.NewUniversalClient(m.config.UniversalOptions())

	// Kiểm tra kết nối với timeout 2 giây
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Ping để xác minh kết nối
	if err := universal.Ping(ctx).Err(); err != nil {
		// Đóng universal client ngay lập tức và không cache khi ping thất bại
		universal.Close()
		return nil, err
	}

	// Cache universal client chỉ khi connection thành công
	m.universal = &universal
	return m.universal, nil
}

// UniversalPing kiểm tra kết nối của Redis universal client.
//
// Phương thức này thực hiện lệnh PING để xác minh rằng kết nối
// đến Redis cluster/sentinel vẫn hoạt động bình thường. Nó sử dụng
// UniversalClient() để lấy instance và thực hiện ping command.
//
// Parameters:
//   - ctx: Context để control timeout và cancellation của operation
//
// Returns:
//   - error: nil nếu ping thành công, ngược lại trả về lỗi
//
// Exceptions:
//   - Trả về error nếu không thể lấy universal client từ UniversalClient()
//   - Trả về error nếu ping command thất bại
//   - Trả về error nếu context bị timeout hoặc cancelled
func (m *manager) UniversalPing(ctx context.Context) error {
	client, err := m.UniversalClient()
	if err != nil {
		return err
	}
	return (*client).Ping(ctx).Err()
}

// Close đóng tất cả các kết nối Redis đang hoạt động.
//
// Phương thức này đảm bảo rằng tất cả resources được giải phóng
// một cách an toàn và các connections được đóng đúng cách.
// Sau khi gọi Close(), các cached connections được set về nil.
//
// Returns:
//   - error: nil nếu đóng thành công, ngược lại trả về lỗi đầu tiên gặp phải
//
// Exceptions:
//   - Trả về error nếu việc đóng client connection thất bại
//   - Trả về error nếu việc đóng universal client connection thất bại
//   - Nếu có nhiều lỗi, chỉ trả về lỗi đầu tiên
//
// Notes:
//   - Manager không thể tái sử dụng các cached connections sau khi Close()
//   - Các lần gọi tiếp theo đến Client() hoặc UniversalClient() sẽ tạo connections mới
func (m *manager) Close() error {
	// Tạo slice để lưu trữ các lỗi
	var errs []error

	// Đóng client connection nếu tồn tại
	if m.client != nil {
		if err := m.client.Close(); err != nil {
			errs = append(errs, err)
		}
		m.client = nil // Set về nil sau khi đóng
	}

	// Đóng universal client connection nếu tồn tại
	if m.universal != nil {
		if err := (*m.universal).Close(); err != nil {
			errs = append(errs, err)
		}
		m.universal = nil // Set về nil sau khi đóng
	}

	// Trả về lỗi đầu tiên nếu có lỗi xảy ra
	if len(errs) > 0 {
		return errs[0] // Trả về lỗi đầu tiên
	}
	return nil
}
