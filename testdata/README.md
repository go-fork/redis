# Redis Test Data

Folder này chứa các file dữ liệu test cho Redis provider theo quy ước Go testing.

## Cấu trúc Files

### Valid Configurations
- `valid_client_config.json`: Cấu hình hợp lệ cho Redis Client
- `valid_universal_config.json`: Cấu hình hợp lệ cho Redis Universal Client
- `complex_config.json`: Cấu hình phức tạp với cả Client và Universal
- `disabled_redis_config.json`: Cấu hình với Redis bị vô hiệu hóa

### Invalid Configurations
- `invalid_client_config.json`: Cấu hình không hợp lệ cho Client (thiếu addr)
- `invalid_universal_config.json`: Cấu hình không hợp lệ cho Universal (thiếu addrs)
- `malformed_config.json`: File JSON bị lỗi format

## Sử dụng trong Tests

Các file này được sử dụng trong unit tests để:
1. Test parsing và validation của config
2. Test các scenarios khác nhau của provider
3. Test error handling với config không hợp lệ
4. Benchmark performance với các config khác nhau

## Quy ước

- Tất cả file config sử dụng format JSON
- File names mô tả rõ ràng mục đích sử dụng
- Bao gồm cả valid và invalid cases để test comprehensive
