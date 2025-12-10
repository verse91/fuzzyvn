# Test

## Chạy tests

```bash
# Test với Makefile
make test

# Chạy tất cả tests
go test -v

# Chạy test cụ thể
go test -v -run TestNormalize

# Benchmark với test Makefile (benchmark chi tiết)
make bench

# Chạy benchmarks
go test -bench=. -benchmem

# Chạy benchmark cụ thể
go test -bench=BenchmarkSearch -benchmem

# Chạy với coverage
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Danh sách Tests

### 1. Tests cho Normalize

| Test | Mô tả |
|------|-------|
| `TestNormalize` | Kiểm tra chuẩn hóa tiếng Việt (bỏ dấu, đ→d) |
| `TestNormalize_IY_Equivalence` | Kiểm tra i/y được coi là giống nhau (kỷ=kỉ) |

```go
// Ví dụ test cases
{"Đường", "Duong"}
{"kỷ niệm", "ki niem"}
{"kỉ niệm", "ki niem"}  // phải giống kỷ niệm
```

### 2. Tests cho LevenshteinRatio

| Test | Mô tả |
|------|-------|
| `TestLevenshteinRatio` | Kiểm tra tính khoảng cách Levenshtein |

```go
// Ví dụ test cases
{"abc", "abc", 0}      // giống nhau
{"main", "mian", 2}    // đảo 2 chữ
{"kitten", "sitting", 3}
```

### 3. Tests cho Searcher

| Test | Mô tả |
|------|-------|
| `TestNewSearcher` | Kiểm tra khởi tạo Searcher |
| `TestSearcher_Search_Basic` | Tìm kiếm cơ bản |
| `TestSearcher_Search_Vietnamese` | Tìm kiếm tiếng Việt không dấu |
| `TestSearcher_Search_IY_Equivalence` | Search "ky niem" = "ki niem" |
| `TestSearcher_Search_Typo` | Chịu lỗi gõ (mian→main) |
| `TestSearchWithVietnameseData` | Test với nhiều file tiếng Việt |
| `TestSearchWithTypos` | Test nhiều lỗi gõ khác nhau |
| `TestSearchCacheBoost` | Kiểm tra cache có boost kết quả |
| `TestSearchWithRealworldData` | Test với dữ liệu thật (~3000 files) |

```go
// Test tiếng Việt
{"bao cao", "Báo_cáo"}      // không dấu → có dấu
{"son tung", "Sơn Tùng"}    // tên ca sĩ
{"ky niem", "Kỷ Niệm"}      // i/y equivalence

// Test lỗi gõ
{"mian", "main.go"}         // đảo chữ
{"conifg", "config.yaml"}   // đảo chữ
{"redame", "README.md"}     // thiếu chữ
```

### 4. Tests cho QueryCache

| Test | Mô tả |
|------|-------|
| `TestQueryCache_RecordSelection` | Ghi nhận lựa chọn |
| `TestQueryCache_GetBoostScores_ExactMatch` | Boost với query khớp chính xác |
| `TestQueryCache_GetBoostScores_SimilarQuery` | Boost với query tương tự |
| `TestQueryCache_GetRecentQueries` | Lấy queries gần đây |
| `TestQueryCache_GetCachedFiles` | Lấy files đã cache |
| `TestQueryCache_GetAllRecentFiles` | Lấy tất cả files gần đây |
| `TestQueryCache_LRU_Eviction` | Kiểm tra xóa cache cũ (LRU) |
| `TestQueryCache_Clear` | Xóa toàn bộ cache |

```go
// Test LRU eviction
cache.SetMaxQueries(3)
cache.RecordSelection("q1", "/a.go")
cache.RecordSelection("q2", "/b.go")
cache.RecordSelection("q3", "/c.go")
cache.RecordSelection("q4", "/d.go")  // q1 bị xóa
```

### 5. Tests tích hợp

| Test | Mô tả |
|------|-------|
| `TestSearcher_RecordSelection_BoostsResults` | Cache boost đưa file lên đầu |
| `TestNewSearcherWithCache` | Giữ cache khi rebuild Searcher |
| `TestSearcher_ClearCache` | Xóa cache của Searcher |

## Benchmarks

### Kết quả benchmark (AMD Ryzen 7 PRO 7840HS)

| Benchmark | Thời gian | Bộ nhớ |
|-----------|-----------|--------|
| `BenchmarkNewSearcher` (1000 files) | ~3ms | 17MB |
| `BenchmarkSearch/100_files` | ~74µs | 38KB |
| `BenchmarkSearch/1000_files` | ~1ms | 255KB |
| `BenchmarkSearch/10000_files` | ~14ms | 2.5MB |
| `BenchmarkSearchVietnamese` | ~1.7ms | 375KB |
| `BenchmarkSearchWithCache` | ~1.1ms | 255KB |
| `BenchmarkNormalize` | ~16µs | 46KB |
| `BenchmarkLevenshteinRatio` | ~456ns | 320B |
| `BenchmarkRecordSelection` | ~2.8µs | 9KB |
| `BenchmarkGetBoostScores` | ~41µs | 25KB |

### Chi tiết benchmarks

```go
// Benchmark tìm kiếm với số lượng files khác nhau
BenchmarkSearch/100_files
BenchmarkSearch/1000_files
BenchmarkSearch/10000_files

// Benchmark tiếng Việt
BenchmarkSearchVietnamese/tiếng_Việt_có_dấu
BenchmarkSearchVietnamese/tiếng_Việt_không_dấu

// Benchmark cache
BenchmarkSearchWithCache
BenchmarkRecordSelection
BenchmarkGetBoostScores
```

## Test với dữ liệu thật

Test `TestSearchWithRealworldData` sử dụng dữ liệu trong `demo/test_data/`:

```
demo/test_data/
├── Music_Lossless/     (~3000 files nhạc Việt)
├── Ebook_Library/      (~3000 files sách)
└── Work_Documents/     (~4000 files tài liệu)
```

Kết quả thực tế:
```
Search 'son tung' trong 2938 files... tìm thấy 20 kết quả trong 5.87ms
Search 'ky niem' trong 2938 files... tìm thấy 20 kết quả trong 5.62ms
Search 'lac troi' trong 2938 files... tìm thấy 20 kết quả trong 6.19ms
```

## Viết test mới

```go
func TestMyFeature(t *testing.T) {
    // Chuẩn bị dữ liệu
    files := []string{
        "/project/main.go",
        "/project/config.yaml",
    }
    
    // Tạo searcher
    searcher := NewSearcher(files)
    
    // Thực hiện search
    results := searcher.Search("main")
    
    // Kiểm tra kết quả
    if !slices.Contains(results, "/project/main.go") {
        t.Error("Không tìm thấy main.go")
    }
}

func BenchmarkMyFeature(b *testing.B) {
    files := generateTestFiles(1000)
    searcher := NewSearcher(files)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        searcher.Search("main")
    }
}
```

## Helper functions

```go
// Tạo files test
generateTestFiles(n int) []string

// Tạo files tiếng Việt test
generateVietnameseTestFiles(n int) []string

// Quét thư mục test_data
scanTestData(root string) []string
```
