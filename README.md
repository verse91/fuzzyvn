<div align="center">
 <img width="20%" width="1920" height="1920" alt="gopher-min" src="https://github.com/user-attachments/assets/a7f7729e-2e34-4ecc-8866-c8c85d93f233" />

  <h1>FuzzyVN</h1>

  [![License: 0BSD](https://img.shields.io/badge/License-0BSD-blue?style=for-the-badge&logo=github&logoColor=white)](./LICENSE.md)
  [![Status](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge&logo=github&logoColor=white)]()
  [![Documentation](https://img.shields.io/badge/docs-available-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./fuzzyvn.go)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./.github/CONTRIBUTING.md)
</div>
<p><b>FuzzyVN là thư viện tìm kiếm file bằng kỹ thuật chính là fuzzy search được tối ưu cho tiếng Việt, và còn nhanh hơn với tiếng Anh. Kết hợp nhiều thuật toán tìm kiếm với hệ thống cache thông minh để cho kết quả nhanh và chính xác</b></p>

> [!IMPORTANT]
> Package này chỉ nên dùng ở local hoặc side project.
> Vui lòng không được sử dụng trong production.
> Mình sẽ không chịu bất kỳ trách nhiệm nào khi bạn sử dụng nó.

<br>

<div align="center">

<img width="70%" width="1414" height="1425" alt="image" src="https://github.com/user-attachments/assets/9266cc9a-1b06-491f-ab17-2f0cbd9dcabb" />

<img width="70%" width="1379" height="1406" alt="image" src="https://github.com/user-attachments/assets/d874a0a8-8642-4d3b-a35c-2c44bb0d9647" />

<img width="70%" width="1320" height="630" alt="image" src="https://github.com/user-attachments/assets/c26711db-c3cd-4d44-b03a-3915b05a03ee" />
</div>

<div align="center"><i>Bạn có thể test qua phần demo</i></div>

## Tính năng

- **Tối ưu cho tiếng Việt**
- **Xử lí lỗi chính tả**
- **Đa thuật toán**
- **Hệ thống cache**
- **Xử lí lỗi gõ**
- **Thread-Safe**
- **Xử lý parallel cho dataset lớn**

## Cài đặt

```bash
go get github.com/verse91/fuzzyvn
```

**Yêu cầu**: Go 1.21+

**Dependencies**: Chỉ cần `golang.org/x/text` để normalize tiếng Việt

## Benchmark
> [!NOTE]
> Benchmark trên laptop thường với AMD Ryzen 7 PRO 7840HS (16 threads)
> ***9ms, 0.5MB RAM, ~200 allocs cho 100k file***

```bash
Search 'son tung' trong 99987 files... tìm thấy 20 kết quả trong 12.778477ms
Search 'ky niem' trong 99987 files... tìm thấy 20 kết quả trong 5.134529ms
Search 'lac troi' trong 99987 files... tìm thấy 20 kết quả trong 7.271664ms

--- Benchmark Info ---
Đã load 99987 files từ ổ cứng
----------------------
goos: linux
goarch: amd64
pkg: github.com/verse91/fuzzyvn
cpu: AMD Ryzen 7 PRO 7840HS w/ Radeon 780M Graphics
BenchmarkSearch_RealWorld/Search/50k_Files-16         	     294	   4568294 ns/op	  270863 B/op	     208 allocs/op
BenchmarkSearch_RealWorld/Search/100k_Files-16        	     129	   9255966 ns/op	  515881 B/op	     212 allocs/op
BenchmarkSearch_RealWorld/Search/100K_Files_Typo-16   	     128	   9205201 ns/op	  413723 B/op	     210 allocs/op
BenchmarkNewSearcher-16                               	     307	   3787931 ns/op	17656070 B/op	   13012 allocs/op
BenchmarkSearch/100_files-16                          	   37952	     31143 ns/op	   45721 B/op	      60 allocs/op
BenchmarkSearch/1000_files-16                         	   10000	    136732 ns/op	   53063 B/op	     198 allocs/op
BenchmarkSearch/10000_files-16                        	    1528	    814648 ns/op	  206543 B/op	     219 allocs/op
BenchmarkSearchVietnamese/tiếng_Việt_có_dấu-16        	    4102	    272148 ns/op	   85530 B/op	     236 allocs/op
BenchmarkSearchVietnamese/tiếng_Việt_không_dấu-16                	    4596	    265909 ns/op	   84361 B/op	     232 allocs/op
BenchmarkSearchWithCache-16                                      	   10000	    113917 ns/op	   53374 B/op	     202 allocs/op
BenchmarkNormalize-16                                            	   74479	     15439 ns/op	   46152 B/op	      39 allocs/op
BenchmarkLevenshteinRatio-16                                     	 4532257	       311.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkRecordSelection-16                                      	  457828	      2698 ns/op	    8789 B/op	       8 allocs/op
BenchmarkGetBoostScores-16                                       	   36382	     32572 ns/op	   18858 B/op	     213 allocs/op
PASS
ok  	github.com/verse91/fuzzyvn	32.195s
```

| Operation                           | Time    | Memory  | Notes            |
| ----------------------------------- | ------- | ------- | ---------------- |
| NewSearcher                         | 3.78ms  | 17.65MB | Load 99987 files |
| Search 100 files                    | 31µs    | 45KB    |                  |
| Search 1K files                     | 136µs   | 53KB    |                  |
| Search 10K files                    | 814µs   | 206KB   |                  |
| Search 50K files                    | 4.56ms  | 270KB   | Bình thường      |
| **Search 100K files**                   | **9.25ms**  | **516KB**   | Bình thường   |
| **Search 100K files, typo**             | **9.20ms**  | **414KB**   | Sai chính tả  |
| Tiếng Việt có dấu                   | 272µs   | 86KB    |                  |
| Tiếng Việt không dấu                | 265µs   | 84KB    |                  |
| Search với Cache                    | 113µs   | 53KB    |                  |
| Normalize                           | 15µs    | 46KB    |                  |
| LevenshteinRatio                    | 311ns   | 0B      |                  |
| RecordSelection                     | 2.6µs   | 8.7KB   |                  |
| GetBoostScores                      | 32.5µs  | 18KB    |                  |

```bash
go test -bench=BenchmarkSearch -benchmem
```
hoặc
```bash
make bench
```

### Luồng tìm kiếm

```
              Query người dùng
                      ↓
        Normalize (bỏ dấu, lowercase)
                      ↓
┌─────────────┬──────────────┬─────────────┐
│ Fuzzy Match │ Levenshtein  │ Cache Boost │
│ (substring) │ (sửa lỗi gõ) │ (lịch sử)   │
└─────────────┴──────────────┴─────────────┘
                      ↓
              Tính điểm tổng hợp
                      ↓
               Sắp xếp theo điểm
                      ↓
                Top 20 kết quả
```

## Kiến trúc

```
┌─────────────────────────────────────────┐
│                Searcher                 │
├─────────────────────────────────────────┤
│ Originals[]     - File paths gốc        │
│ Normalized[]    - Đã normalize          │
│ FilenamesOnly[] - Chỉ tên file          │
│ Cache          - QueryCache             │
└─────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────┐
│                QueryCache               │
├─────────────────────────────────────────┤
│ entries{}      - query → CacheEntry[]   │
│ queryOrder[]   - LRU tracking           │
│ maxQueries     - Limit (100)            │
│ maxPerQuery    - Files per query (5)    │
│ boostScore     - Boost factor (5000)    │
└─────────────────────────────────────────┘
```

## Nếu bạn muốn tự phát triển

### Nhớ tạo data từ demo trước khi chạy (mình không up lên đây vì lí do dung lương)
```bash
make gen
```
> [!WARNING]
> Phải tạo data trước khi chạy nếu không sẽ lỗi
### Chạy demo
```bash
make demo
```
Kết quả:
> Server running at http://localhost:8080
> Scanning files from directory: ./test_data
> Indexed 99987 files. Cache: 0 queries

### Test
```bash
make test
```
### Benchmark
```go
make bench
```
hoặc benchmark cụ thể

```go
go test -bench=BenchmarkLevenshteinRatio -benchmem -count=1
```
## Cách dùng

### Ví dụ cơ bản

```go
package main

import (
    "fmt"
    "github.com/verse91/fuzzyvn"
)

func main() {
    // Tạo searcher với danh sách file
    files := []string{
        "/home/user/Documents/Báo_cáo_tháng_1.pdf",
        "/home/user/Music/Sơn_Tùng_MTP.mp3",
        "/home/user/Code/main.go",
    }
    searcher := fuzzyvn.NewSearcher(files)

    // Tìm kiếm
    results := searcher.Search("bao cao")
    for _, path := range results {
        fmt.Println(path)
    }
    // Output: /home/user/Documents/Báo_cáo_tháng_1.pdf
}
```

### Ví dụ với Cache

```go
// Người dùng tìm kiếm
results := searcher.Search("main")

// Người dùng chọn file
selectedFile := results[0]
searcher.RecordSelection("main", selectedFile)

// Lần tìm kiếm sau, file này được ưu tiên
results = searcher.Search("mai")  // Gõ sai, vẫn lên đầu nhờ cache
```

### Ví dụ với HTTP Server

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

    mu.Lock()
    searcher.RecordSelection(req.Query, req.Path)
    mu.Unlock()

    w.WriteHeader(http.StatusOK)
}

func main() {
    // Scan và index files
    files := scanDirectory("/data")
    searcher = fuzzyvn.NewSearcher(files)

    http.HandleFunc("/api/search", searchHandler)
    http.HandleFunc("/api/select", selectHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Tài liệu

### API chính

#### `NewSearcher(items []string) *Searcher`
Tạo searcher mới từ danh sách file paths

```go
searcher := fuzzyvn.NewSearcher(files)
```

#### `Search(query string) []string`
Tìm kiếm và trả về top 20 kết quả phù hợp nhất (hardcode 20)

```go
results := searcher.Search("readme")
```

#### `RecordSelection(query, filePath string)`
Lưu lại file mà người dùng đã chọn để cải thiện kết quả tương lai

```go
searcher.RecordSelection("main", "/project/main.go")
```

#### `GetCache() *QueryCache`
Lấy cache object để tùy chỉnh hoặc xem thống kê

```go
cache := searcher.GetCache()
cache.SetBoostScore(10000)      // Tăng boost
cache.SetMaxQueries(500)        // Lưu nhiều query hơn
recentQueries := cache.GetRecentQueries(10)
```

### QueryCache Methods

```go
cache := searcher.GetCache()

// Cấu hình
cache.SetBoostScore(score int)        // Mặc định: 5000
cache.SetMaxQueries(n int)            // Mặc định: 100

// Thống kê
cache.GetRecentQueries(limit int) []string
cache.GetAllRecentFiles(limit int) []string
cache.GetCachedFiles(query string, limit int) []string
cache.Size() int
cache.Clear()
```

### Utility Functions

```go
// Normalize string (bỏ dấu tiếng Việt)
normalized := fuzzyvn.Normalize("Tiếng Việt")
// Output: "Tieng Viet"

// Tính khoảng cách Levenshtein
distance := fuzzyvn.LevenshteinRatio("hello", "helo")
// Output: 1

// Fuzzy find trong slice
matches := fuzzyvn.FuzzyFind("pattern", targets)
```

## Cách hoạt động

### Điểm số (Scoring)

Mỗi kết quả nhận điểm từ nhiều nguồn:

1. **Fuzzy Score** (0-1000+)
   - Match ký tự: +16 mỗi ký tự
   - Ký tự đầu tiên: +24
   - Match liên tiếp: +28
   - Sau word boundary: +20
   - Sau slash: +24
   - CamelCase: +16
   - Gap penalty: -2 mỗi ký tự

2. **Word Bonus** (0-9000+)
   - +3000 cho mỗi từ khớp hoàn toàn
   - Cho phép 1 lỗi với từ ≥3 ký tự

3. **Levenshtein Score** (0-10000)
   - Cho phép ~33% lỗi
   - 10000 - (lỗi × 100)

4. **Cache Boost** (0-10000+)
   - Dựa trên số lần chọn
   - Độ tương đồng query
   - Công thức: `(boostScore × similarity × selectCount) / 100`

### Cache System

Cache hoạt động theo cơ chế LRU (Least Recently Used):

```go
// Mỗi query lưu tối đa 5 files
// Mỗi file có selectCount (số lần chọn)
// Query có độ tương đồng cao được tận dụng cache
```

**Ví dụ**:
- User search `"màn hình"` → chọn `"dell-monitor.pdf"`
- User search `"man hinh"` → `"dell-monitor.pdf"` lên top (similarity 95%)
- User search `"màn hình dell"` → vẫn boost (contains)

## Các trường hợp sử dụng

### 1. File Explorer / Launcher

```go
// Quét thư mục home
files := scanDirectory("/home/user")
searcher := fuzzyvn.NewSearcher(files)

// User gõ, realtime search
results := searcher.Search(userInput)
```

### 2. Document Management

```go
// Index tài liệu công ty
docs := scanWithExtensions("/company/docs", []string{".pdf", ".docx"})
searcher := fuzzyvn.NewSearcher(docs)

// Tìm hợp đồng
contracts := searcher.Search("hop dong")
```

### 3. Code Search

```go
// Index source code
code := scanIgnoreDirs("/project", []string{"node_modules", ".git"})
searcher := fuzzyvn.NewSearcher(code)

// Tìm file main
mains := searcher.Search("main")
```

### 4. Media Library

```go
// Index nhạc
music := scanWithExtensions("/music", []string{".mp3", ".flac"})
searcher := fuzzyvn.NewSearcher(music)

// Tìm bài hát
songs := searcher.Search("son tung")
```

## Ví dụ nâng cao

### Rebuild Index khi file thay đổi

```go
func watchAndRebuild(searcher **fuzzyvn.Searcher) {
    watcher := setupFileWatcher()

    for event := range watcher.Events {
        // Giữ lại cache
        cache := (*searcher).GetCache()

        // Quét lại
        newFiles := scanDirectory("/data")

        // Rebuild với cache cũ
        *searcher = fuzzyvn.NewSearcherWithCache(newFiles, cache)
    }
}
```

### Tùy chỉnh cho domain cụ thể

```go
searcher := fuzzyvn.NewSearcher(files)
cache := searcher.GetCache()

// Tăng boost cho người dùng power user
cache.SetBoostScore(15000)

// Lưu nhiều lịch sử hơn
cache.SetMaxQueries(1000)
```

### Integration với CLI tool

```go
func main() {
    files := scanDirectory(os.Getenv("HOME"))
    searcher := fuzzyvn.NewSearcher(files)

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Search> ")
        query, _ := reader.ReadString('\n')
        query = strings.TrimSpace(query)

        results := searcher.Search(query)
        for i, r := range results {
            fmt.Printf("[%d] %s\n", i, r)
        }

        fmt.Print("Select> ")
        input, _ := reader.ReadString('\n')
        idx, _ := strconv.Atoi(strings.TrimSpace(input))

        if idx >= 0 && idx < len(results) {
            searcher.RecordSelection(query, results[idx])
            // Open file...
        }
    }
}
```

## Đóng góp

Vui lòng theo chuẩn [Contributing](.github/CONTRIBUTING.md) khi tạo một contribution qua pull request.

## License

Package này được cấp phép bởi giấy phép [0BSD License](LICENSE). Bạn có thể sửa, xóa, thêm hay làm bất cứ thứ gì bạn muốn với nó.
