# FuzzyVN

FuzzyVN là thư viện tìm kiếm file bằng kỹ thuật chính là fuzzy search được tối ưu cho tiếng Việt. Kết hợp nhiều thuật toán tìm kiếm với hệ thống cache thông minh để cho kết quả nhanh và chính xác.
<div align="center">
  
<img width="70%" width="1414" height="1425" alt="image" src="https://github.com/user-attachments/assets/9266cc9a-1b06-491f-ab17-2f0cbd9dcabb" />

<img width="70%" width="1379" height="1406" alt="image" src="https://github.com/user-attachments/assets/d874a0a8-8642-4d3b-a35c-2c44bb0d9647" />

<img width="70%" width="1320" height="630" alt="image" src="https://github.com/user-attachments/assets/c26711db-c3cd-4d44-b03a-3915b05a03ee" />
</div>

<div align="center"><i>Bạn có thể test qua phần demo</i></div>

## Mục lục

- [Cài đặt](./docs/installation.md)
- [Bắt đầu nhanh](./docs/quickstart.md)
- [API](./docs/api.md)
- [Hệ thống Cache](./docs/cache.md)
- [Thuật toán](./docs/algorithm.md)
- [Ví dụ](./docs/examples.md)
- [Test](./docs/test.md)

## Tính năng

- **Hỗ trợ tiếng Việt**: Xử lý dấu tiếng Việt (chuyển "Đường" thành "Duong")
- **Đa thuật toán**: Kết hợp fuzzy matching + Levenshtein distance
- **Cache thông minh**: Học từ lựa chọn của người dùng để đẩy kết quả liên quan lên đầu
- **Chịu lỗi gõ**: Xử lý lỗi gõ phím thường gặp
- **Thread-Safe**: An toàn khi truy cập đồng thời

## Kiến trúc

```
┌─────────────────────────────────────────────────────────┐
│                      Searcher                           │
├─────────────────────────────────────────────────────────┤
│  Originals[]     - Đường dẫn file gốc                   │
│  Normalized[]    - Đã chuẩn hóa cho fuzzy search        │
│  FilenamesOnly[] - Chỉ tên file cho Levenshtein         │
│  Cache           - Cache query để boost kết quả         │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    QueryCache                           │
├─────────────────────────────────────────────────────────┤
│  entries{}       - query → []CacheEntry                 │
│  queryOrder[]    - Thứ tự LRU                           │
│  maxQueries      - Tối đa queries cache (100)           │
│  maxPerQuery     - Tối đa files mỗi query (5)           │
│  boostScore      - Hệ số boost (5000)                   │
└─────────────────────────────────────────────────────────┘
```

## Luồng tìm kiếm

```
Query người dùng
    │
    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Fuzzy Match  │ +  │ Levenshtein  │ +  │ Cache Boost  │
│ (substring)  │    │ (sửa lỗi gõ) │    │ (lịch sử)    │
└──────────────┘    └──────────────┘    └──────────────┘
    │                      │                    │
    └──────────────────────┼────────────────────┘
                           ▼
                    Gộp kết quả
                           │
                           ▼
                    Sắp xếp theo điểm
                           │
                           ▼
                    Top 20 kết quả
```
## Benchmark
- Lưu ý rằng benchmark này chỉ thực hiện trên một laptop bình thường

```bash
> go test -bench=BenchmarkSearch -benchmem

Search 'son tung' trong 2938 files... tìm thấy 20 kết quả trong 1.51568ms
Search 'ky niem' trong 2938 files... tìm thấy 20 kết quả trong 1.042098ms
Search 'lac troi' trong 2938 files... tìm thấy 20 kết quả trong 979.804µs
goos: linux
goarch: amd64
pkg: github.com/verse91/fuzzyvn
cpu: AMD Ryzen 7 PRO 7840HS w/ Radeon 780M Graphics 
BenchmarkNewSearcher-16         	     402	   3012563 ns/op	17585024 B/op	   13005 allocs/op
BenchmarkSearch/100_files-16    	   27104	     44289 ns/op	   56752 B/op	     486 allocs/op
BenchmarkSearch/1000_files-16   	    4700	    247986 ns/op	  360982 B/op	    2368 allocs/op
BenchmarkSearch/10000_files-16  	     516	   2414495 ns/op	 3393652 B/op	   20491 allocs/op
BenchmarkSearchVietnamese/tiếng_Việt_có_dấu-16                	    2473	    449232 ns/op	  502476 B/op	    2309 allocs/op
BenchmarkSearchVietnamese/tiếng_Việt_không_dấu-16             	    2697	    443719 ns/op	  501403 B/op	    2305 allocs/op
BenchmarkSearchWithCache-16                                   	    4078	    294842 ns/op	  361230 B/op	    2372 allocs/op
BenchmarkNormalize-16                                         	   85342	     13355 ns/op	   46152 B/op	      39 allocs/op
BenchmarkLevenshteinRatio-16                                  	 3456462	       337.5 ns/op	     320 B/op	       4 allocs/op
BenchmarkRecordSelection-16                                   	  528936	      2278 ns/op	    8789 B/op	       8 allocs/op
BenchmarkGetBoostScores-16                                    	   41019	     29497 ns/op	   25112 B/op	     311 allocs/op
PASS
ok  	github.com/verse91/fuzzyvn	15.259s

```

## Kết quả benchmark

| Benchmark | Time | Memory |
|-----------|------|--------|
| **NewSearcher** | 3.0ms | 17.5MB |
| **Search 100 files** | 44µs | 56KB |
| **Search 1,000 files** | 248µs | 360KB |
| **Search 10,000 files** | **2.4ms** | 3.4MB |
| **Tiếng Việt có dấu** | 449µs | 502KB |
| **Tiếng Việt không dấu** | 444µs | 501KB |
| **Search với Cache** | 295µs | 361KB |
| **Normalize** | 13µs | 46KB |
| **LevenshteinRatio** | 338ns | 320B |
| **RecordSelection** | 2.3µs | 8.7KB |
| **GetBoostScores** | 29µs | 25KB |


# Cách dùng

## Tìm kiếm file cơ bản

```go
package main

import (
    "fmt"
    "io/fs"
    "path/filepath"
    
    "github.com/verse91/fuzzyvn"
)

func main() {
    // Quét thư mục
    var files []string
    filepath.WalkDir("/home/user", func(path string, d fs.DirEntry, err error) error {
        if err == nil && !d.IsDir() {
            files = append(files, path)
        }
        return nil
    })

    // Tạo searcher
    searcher := fuzzyvn.NewSearcher(files)

    // Tìm kiếm
    results := searcher.Search("readme")
    for _, r := range results {
        fmt.Println(r)
    }
}
```

## HTTP Server với Cache

```go
package main

import (
    "encoding/json"
    "net/http"
    "sync"
    
    "github.com/verse91/fuzzyvn"
)

var (
    searcher *fuzzyvn.Searcher
    mu       sync.RWMutex
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    
    mu.RLock()
    results := searcher.Search(query)
    mu.RUnlock()
    
    json.NewEncoder(w).Encode(results)
}

func selectHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Query string `json:"query"`
        Path  string `json:"path"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    
    mu.RLock()
    searcher.RecordSelection(req.Query, req.Path)
    mu.RUnlock()
    
    w.WriteHeader(http.StatusOK)
}

func main() {
    files := scanDirectory("/data")
    searcher = fuzzyvn.NewSearcher(files)
    
    http.HandleFunc("/search", searchHandler)
    http.HandleFunc("/select", selectHandler)
    http.ListenAndServe(":8080", nil)
}
```


## Tùy chỉnh điểm số

```go
// Lấy cache và tùy chỉnh
cache := searcher.GetCache()

// Tăng boost cho kết quả cache
cache.SetBoostScore(10000)  // Mặc định: 5000

// Giữ nhiều query hơn trong cache
cache.SetMaxQueries(500)    // Mặc định: 100
```

## Rebuild Index

```go
// File watcher phát hiện thay đổi
func onFileSystemChange() {
    // Giữ lại cache
    cache := searcher.GetCache()
    
    // Quét lại file
    newFiles := scanDirectory("/data")
    
    // Rebuild với cache cũ
    mu.Lock()
    searcher = fuzzyvn.NewSearcherWithCache(newFiles, cache)
    mu.Unlock()
}
```

## Lấy hoạt động gần đây

```go
cache := searcher.GetCache()

// Lấy query tìm kiếm gần đây
recentQueries := cache.GetRecentQueries(10)
// ["main.go", "config", "readme", ...]

// Lấy file đã chọn gần đây
recentFiles := cache.GetAllRecentFiles(5)
// ["/project/main.go", "/project/config.yaml", ...]

// Lấy file cache cho query hiện tại
cachedForQuery := cache.GetCachedFiles("main", 5)
// File đã chọn trước đó cho "main" hoặc query tương tự
```

## Ví dụ tiếng Việt

```go
searcher := fuzzyvn.NewSearcher([]string{
    "/docs/Báo_cáo_tháng_1.pdf",
    "/docs/Hợp_đồng_thuê_nhà.docx",
    "/music/Sơn Tùng - Lạc Trôi.mp3",
})

// Tất cả các query này đều hoạt động:
searcher.Search("bao cao")      // khớp "Báo_cáo"
searcher.Search("hop dong")     // khớp "Hợp_đồng"
searcher.Search("son tung")     // khớp "Sơn Tùng"
searcher.Search("lac troi")     // khớp "Lạc Trôi"

// Gõ sai cũng được:
searcher.Search("bao coa")      // lỗi gõ: "coa" → "cao"
searcher.Search("sontung")      // thiếu dấu cách
searcher.Search("sont ung")     // dấu cách sai chỗ
```

## Quét thư mục nâng cao

```go
// Quét với filter extension
func scanWithExtensions(root string, exts []string) []string {
    var files []string
    extMap := make(map[string]bool)
    for _, ext := range exts {
        extMap[ext] = true
    }

    filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil || d.IsDir() {
            return nil
        }
        if len(exts) == 0 || extMap[filepath.Ext(path)] {
            files = append(files, path)
        }
        return nil
    })
    return files
}

// Quét và bỏ qua thư mục
func scanIgnoreDirs(root string, ignore []string) []string {
    var files []string
    ignoreMap := make(map[string]bool)
    for _, dir := range ignore {
        ignoreMap[dir] = true
    }

    filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        if d.IsDir() && ignoreMap[d.Name()] {
            return filepath.SkipDir
        }
        if !d.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files
}

// Sử dụng
files := scanWithExtensions("/project", []string{".go", ".md"})
files = scanIgnoreDirs("/project", []string{"node_modules", ".git", "vendor"})
```

