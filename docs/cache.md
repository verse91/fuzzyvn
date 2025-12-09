# Hệ thống Cache

Hệ thống cache học từ hành vi người dùng để cải thiện kết quả tìm kiếm theo thời gian.

## Cách hoạt động

### Ghi nhận lựa chọn

Khi người dùng chọn file từ kết quả tìm kiếm:

```go
searcher.RecordSelection("main server", "/project/src/main_server.go")
```

Cache lưu:
- Query: `"main server"` (đã chuẩn hóa)
- File: `/project/src/main_server.go`
- Số lần: 1 (tăng khi chọn lại)

### Độ tương đồng Query

Cache không yêu cầu query khớp chính xác. Nó dùng điểm tương đồng:

| Loại khớp | Điểm | Ví dụ |
|-----------|------|-------|
| Khớp chính xác | 100 | `"main"` = `"main"` |
| Substring (query nằm trong cache) | 80 | `"main"` ⊂ `"main server"` |
| Prefix (tiền tố) | 70-100 | `"mai"` → `"main"` |
| Reverse prefix | 50-80 | `"main server"` → `"main"` |
| Từ chung | 50-95 | `"server main"` ↔ `"main server"` |
| Fuzzy (≤30% lỗi) | 30-60 | `"mian"` ≈ `"main"` |

### Tính điểm Boost

```
boost = boostScore × similarity × selectCount / 100
```

Mặc định `boostScore` là 5000, nên:
- Khớp chính xác, chọn 1 lần: `5000 × 100 × 1 / 100 = 5000`
- Khớp chính xác, chọn 3 lần: `5000 × 100 × 3 / 100 = 15000`
- Prefix match (80%), chọn 2 lần: `5000 × 80 × 2 / 100 = 8000`

## Cấu hình

### Số query tối đa

```go
cache := searcher.GetCache()
cache.SetMaxQueries(200)  // Mặc định: 100
```

Khi vượt giới hạn, query cũ nhất bị loại bỏ (LRU).

### Số file tối đa mỗi query

Mỗi query lưu tối đa 5 file. Khi vượt, file có `selectCount` thấp nhất bị xóa.

### Điểm Boost

```go
cache.SetBoostScore(10000)  // Mặc định: 5000
```

Giá trị cao hơn làm kết quả từ cache nổi bật hơn.

## Loại bỏ LRU

Query được sắp xếp theo thời gian gần đây:
- Mỗi lần chọn đưa query lên "gần đây nhất"
- Khi vượt `maxQueries`, query cũ nhất bị xóa

## Lưu trữ

Cache chỉ lưu trong bộ nhớ. Để lưu qua các lần khởi động lại:

```go
// Trước khi tắt
cache := searcher.GetCache()
// Serialize cache.entries và cache.queryOrder ra JSON/gob

// Khi khởi động
// Deserialize và tạo cache mới
searcher = fuzzyvn.NewSearcherWithCache(files, loadedCache)
```

## Thread Safety

Mọi thao tác cache được bảo vệ bởi `sync.RWMutex`:
- Thao tác đọc: `RLock`
- Thao tác ghi: `Lock`

An toàn khi truy cập đồng thời từ nhiều goroutine.
