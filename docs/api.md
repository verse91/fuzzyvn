# Tham khảo API

## Searcher

### Kiểu dữ liệu

```go
type Searcher struct {
    Originals     []string     // Đường dẫn file gốc
    Normalized    []string     // Chuỗi đã chuẩn hóa cho fuzzy search
    FilenamesOnly []string     // Chỉ tên file cho Levenshtein
    Cache         *QueryCache  // Cache query
}

type MatchResult struct {
    Str   string  // Đường dẫn file
    Score int     // Điểm khớp
}
```

### Hàm khởi tạo

#### NewSearcher

```go
func NewSearcher(items []string) *Searcher
```

Tạo Searcher mới với danh sách đường dẫn file.

**Tham số:**
- `items`: Danh sách đường dẫn file cần index

**Trả về:**
- `*Searcher`: Instance searcher mới với cache rỗng

---

#### NewSearcherWithCache

```go
func NewSearcherWithCache(items []string, cache *QueryCache) *Searcher
```

Tạo Searcher mới với cache có sẵn.

**Tham số:**
- `items`: Danh sách đường dẫn file cần index
- `cache`: QueryCache có sẵn để tái sử dụng

**Trả về:**
- `*Searcher`: Instance searcher mới với cache được cung cấp

---

### Phương thức

#### Search

```go
func (s *Searcher) Search(query string) []string
```

Tìm kiếm file khớp với query.

**Tham số:**
- `query`: Chuỗi tìm kiếm (dấu tiếng Việt được xử lý tự động)

**Trả về:**
- `[]string`: Top 20 đường dẫn file khớp nhất, sắp xếp theo độ liên quan

---

#### RecordSelection

```go
func (s *Searcher) RecordSelection(query, filePath string)
```

Ghi nhận người dùng đã chọn file nào cho query nào. File này sẽ được boost trong các lần tìm kiếm tương tự sau.

**Tham số:**
- `query`: Query tìm kiếm
- `filePath`: File đã được chọn

---

#### ClearCache

```go
func (s *Searcher) ClearCache()
```

Xóa toàn bộ cache.

---

#### GetCache

```go
func (s *Searcher) GetCache() *QueryCache
```

Lấy cache để tái sử dụng khi rebuild searcher.

---

## QueryCache

### Kiểu dữ liệu

```go
type QueryCache struct {
    // Các trường nội bộ (không export)
}

type CacheEntry struct {
    FilePath    string  // Đường dẫn file đã cache
    SelectCount int     // Số lần được chọn
}
```

### Hàm khởi tạo

#### NewQueryCache

```go
func NewQueryCache() *QueryCache
```

Tạo cache rỗng với cấu hình mặc định:
- `maxQueries`: 100
- `maxPerQuery`: 5
- `boostScore`: 5000

---

### Phương thức

#### SetMaxQueries

```go
func (c *QueryCache) SetMaxQueries(n int)
```

Đặt số lượng query tối đa được cache (loại bỏ theo LRU).

---

#### SetBoostScore

```go
func (c *QueryCache) SetBoostScore(score int)
```

Đặt điểm boost cơ bản cho kết quả từ cache.

---

#### RecordSelection

```go
func (c *QueryCache) RecordSelection(query, filePath string)
```

Ghi nhận lựa chọn query-file.

---

#### GetBoostScores

```go
func (c *QueryCache) GetBoostScores(query string) map[string]int
```

Lấy điểm boost cho các file khớp với query tương tự.

**Trả về:**
- `map[string]int`: Đường dẫn file → điểm boost

---

#### GetCachedFiles

```go
func (c *QueryCache) GetCachedFiles(query string, limit int) []string
```

Lấy các file đã cache cho query tương tự.

**Tham số:**
- `query`: Query tìm kiếm hiện tại
- `limit`: Số file tối đa trả về

**Trả về:**
- `[]string`: Đường dẫn file đã cache, sắp xếp theo độ liên quan

---

#### GetRecentQueries

```go
func (c *QueryCache) GetRecentQueries(limit int) []string
```

Lấy các query gần đây nhất (thứ tự MRU - mới nhất trước).

---

#### GetAllRecentFiles

```go
func (c *QueryCache) GetAllRecentFiles(limit int) []string
```

Lấy tất cả file đã cache gần đây, không phụ thuộc query.

---

#### Size

```go
func (c *QueryCache) Size() int
```

Trả về số lượng query đã cache.

---

#### Clear

```go
func (c *QueryCache) Clear()
```

Xóa toàn bộ cache.

---

## Hàm tiện ích

#### Normalize

```go
func Normalize(s string) string
```

Chuẩn hóa văn bản tiếng Việt bằng cách bỏ dấu.

**Ví dụ:**
- `"Đường"` → `"Duong"`
- `"Nguyễn"` → `"Nguyen"`
- `"café"` → `"cafe"`

---

#### LevenshteinRatio

```go
func LevenshteinRatio(s1, s2 string) int
```

Tính khoảng cách Levenshtein (số phép sửa) giữa hai chuỗi.

**Trả về:**
- `int`: Số phép sửa (thêm, xóa, thay thế)
