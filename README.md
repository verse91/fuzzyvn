<div align="center">
 <img width="20%" width="1920" height="1920" alt="gopher-min" src="https://github.com/user-attachments/assets/a7f7729e-2e34-4ecc-8866-c8c85d93f233" />

  <h1>FuzzyVN</h1>

  [![License: 0BSD](https://img.shields.io/badge/License-0BSD-blue?style=for-the-badge&logo=github&logoColor=white)](./LICENSE.md)
  [![Status](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge&logo=github&logoColor=white)]()
  [![Documentation](https://img.shields.io/badge/docs-available-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./fuzzyvn.go)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=for-the-badge&logo=github&logoColor=white)](./.github/CONTRIBUTING.md)
</div>
<p><b>FuzzyVN là thư viện tìm kiếm file bằng kỹ thuật chính là fuzzy matching được tối ưu cho tiếng Việt, và còn nhanh hơn với tiếng Anh. Kết hợp nhiều thuật toán tìm kiếm với hệ thống cache thông minh để cho kết quả nhanh và chính xác</b></p>

> [!IMPORTANT]
> **FuzzyVN tập trung vào việc tăng khả năng chính xác khi sai chính tả để đánh đổi một phần tốc độ nhưng vẫn đảm bảo tốc độ cần thiết cho dự án.**  
> Fuzzyvn vẫn hỗ trợ rất tốt cho tiếng Anh.  
> FuzzyVN hỗ trợ tốt nhất cho tìm kiếm file theo file path thay vì chỉ mỗi tên file (Single String), có thể sẽ dùng thừa tài nguyên cũng như điểm số có thể sai lệch một chút (sẽ không ảnh hưởng nhiều).  
> Package này chỉ nên dùng ở local hoặc side project.  
> Vui lòng không được sử dụng trong production.  
> Mình sẽ không chịu bất kỳ trách nhiệm nào khi bạn sử dụng nó.

<br>

<div align="center">
 <img width="1320" height="630" alt="image" src="https://github.com/user-attachments/assets/c26711db-c3cd-4d44-b03a-3915b05a03ee" />
</div>

<div align="center"><i>Bạn có thể test qua phần <a href="https://github.com/verse91/fuzzyvn/tree/main/demo">demo</a></i></div>

## Tính năng

- **Tối ưu cho tiếng Việt**
- **Xử lí lỗi chính tả**
- **Đa thuật toán**
- **Hệ thống cache**
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

```bash
Search 'son tung' trong 99989 files... tìm thấy 20 kết quả trong 40.652347ms
Search 'ky niem' trong 99989 files... tìm thấy 20 kết quả trong 44.309346ms
Search 'lac troi' trong 99989 files... tìm thấy 20 kết quả trong 38.466203ms

--- Benchmark Info ---
Đã load 99989 files từ ổ cứng
----------------------
goos: linux
goarch: amd64
pkg: github.com/verse91/fuzzyvn
cpu: AMD Ryzen 7 PRO 7840HS w/ Radeon 780M Graphics
BenchmarkSearch_RealWorld/Search/50k_Files-16         	      97	  16962948 ns/op	  325842 B/op	     133 allocs/op
BenchmarkSearch_RealWorld/Search/100k_Files-16        	      44	  33841020 ns/op	  654466 B/op	     135 allocs/op
BenchmarkSearch_RealWorld/Search/100K_Files_Typo-16   	      42	  32882259 ns/op	  553335 B/op	     132 allocs/op
BenchmarkNewSearcher-16                               	    5857	    171898 ns/op	  152016 B/op	    1012 allocs/op
BenchmarkSearch/100_files-16                          	   38449	     30274 ns/op	   21765 B/op	      45 allocs/op
BenchmarkSearch/1000_files-16                         	    4848	    271137 ns/op	   33196 B/op	     162 allocs/op
BenchmarkSearch/10000_files-16                        	     349	   3492175 ns/op	  276717 B/op	    1144 allocs/op
BenchmarkSearchVietnamese/tiếng_Việt_có_dấu-16        	    2362	    542172 ns/op	   55913 B/op	     171 allocs/op
BenchmarkSearchVietnamese/tiếng_Việt_không_dấu-16                	    1831	    602954 ns/op	   55866 B/op	     167 allocs/op
BenchmarkSearchWithCache-16                                      	    4622	    298519 ns/op	   33439 B/op	     165 allocs/op
BenchmarkNormalize-16                                            	  677215	      1584 ns/op	     136 B/op	       9 allocs/op
BenchmarkLevenshteinRatio-16                                     	 3766005	       333.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkRecordSelection-16                                      	 4060006	       340.6 ns/op	      29 B/op	       2 allocs/op
BenchmarkGetBoostScores-16                                       	   31498	     35362 ns/op	   10097 B/op	     207 allocs/op
PASS
ok  	github.com/verse91/fuzzyvn	25.206s
```

| Operation                  | Time    | Memory | Notes             |
| -------------------------- | ------- | ------ | ----------------- |
| NewSearcher                | 0.17ms  | 148KB  | Load 99,989 files |
| Search 100 files           | 30µs    | 21KB   |                   |
| Search 1K files            | 271µs   | 32KB   |                   |
| Search 10K files           | 3.49ms  | 270KB  |                   |
| Search 50K files           | 16.96ms | 318KB  | Bình thường       |
| Search 100K files          | 33.84ms | 639KB  | Bình thường       |
| Search 100K files, typo    | 32.88ms | 540KB  | Sai chính tả      |
| Vietnamese with accents    | 542µs   | 55KB   |                   |
| Vietnamese without accents | 603µs   | 55KB   |                   |
| Search with cache          | 299µs   | 33KB   |                   |
| Normalize                  | 1.58µs  | 136B   | Gần zero alloc    |
| LevenshteinRatio           | 334ns   | 0B     | Zero allocation   |
| RecordSelection            | 341ns   | 29B    | Gần zero alloc    |
| GetBoostScores             | 35.3µs  | 9.9KB  |                   |


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

<details open>
  <summary><b>Ví dụ cơ bản</b></summary>
<br>

```go
//go:build ignore

package main

import (
	"fmt"

	"github.com/verse91/fuzzyvn"
)

func main() {
	// 1. TẠO SEARCHER
	// Bạn hoàn toàn có thể đọc từ 1 folder, đây chỉ là ví dụ đơn giản
	files := []string{
		"/home/user/Documents/Báo_cáo_tháng_1.pdf",
		"/home/user/Documents/Hợp_đồng_thuê_nhà.docx",
		"/home/user/Music/Sơn_Tùng_MTP.mp3",
		"/home/user/Code/main.go",
		"/home/user/Code/utils.go",
	}

	searcher := fuzzyvn.NewSearcher(files)

	// 2. TÌM KIẾM CƠ BẢN
	fmt.Println("--- Tìm 'bao cao' ---")
	results := searcher.Search("bao cao")
	for _, path := range results {
		fmt.Println("  →", path)
	}
	// Output: /home/user/Documents/Báo_cáo_tháng_1.pdf

	// 3. TÌM KIẾM KHÔNG DẤU
	fmt.Println("\n--- Tìm 'son tung' (không dấu) ---")
	results = searcher.Search("son tung")
	for _, path := range results {
		fmt.Println("  →", path)
	}
	// Output: /home/user/Music/Sơn_Tùng_MTP.mp3

	// 4. SỬA LỖI CHÍNH TẢ (Levenshtein)
	fmt.Println("\n--- Tìm 'maiin' (gõ sai) ---")
	results = searcher.Search("maiin")
	for _, path := range results {
		fmt.Println("  →", path)
	}
	// Output: /home/user/Code/main.go

	// 5. CACHE SYSTEM - Học hành vi người dùng
	fmt.Println("\n--- Cache Demo ---")

	// User tìm "main" và chọn main.go
	searcher.RecordSelection("main", "/home/user/Code/main.go")

	// Chọn thêm 2 lần nữa
	searcher.RecordSelection("main", "/home/user/Code/main.go")
	searcher.RecordSelection("main", "/home/user/Code/main.go")

	// Giờ tìm với từ tương tự → main.go được boost lên top
	fmt.Println("Tìm 'mai' (sau khi đã cache):")
	results = searcher.Search("mai")
	for _, path := range results {
		fmt.Println("  →", path)
	}
	// main.go sẽ lên đầu vì đã được chọn 3 lần

	// 6. XEM THỐNG KÊ CACHE
	cache := searcher.GetCache()

	fmt.Println("\n--- Thống kê ---")
	fmt.Println("Recent queries:", cache.GetRecentQueries(3))
	fmt.Println("Recent files:", cache.GetAllRecentFiles(3))
	fmt.Printf("Tổng queries: %d\n", cache.Size())

	// 7. TÙY CHỈNH CACHE
	cache.SetBoostScore(10000) // Tăng độ ưu tiên cho cache
	cache.SetMaxQueries(200)   // Lưu nhiều queries hơn

	fmt.Println("\n✓ Đã cấu hình cache!")
}
```
</details>

<details>
  <summary><b>Ví dụ với Cache</b></summary>
<br>


```go
// Người dùng tìm kiếm
results := searcher.Search("main")

// Người dùng chọn file
selectedFile := results[0]
searcher.RecordSelection("main", selectedFile)

// Lần tìm kiếm sau, file này được ưu tiên
results = searcher.Search("mai")  // Gõ sai, vẫn lên đầu nhờ cache
```
</details>

<details>
  <summary><b>Ví dụ với HTTP Server</b></summary>
<br>

Xem ví dụ ở [demo](https://github.com/verse91/fuzzyvn/tree/main/demo)
</details>

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
<div align="center">
 <img width="70%" width="1414" height="1425" alt="image" src="https://github.com/user-attachments/assets/9266cc9a-1b06-491f-ab17-2f0cbd9dcabb" />
</div>

Mỗi kết quả nhận điểm từ nhiều nguồn:

1. **Fuzzy Score** (0-1000+)
   - Đầu từ (word start): +80
   - Match liên tiếp (consecutive): +40
   - Match thường: +10
   - Phạt độ dài: -(lenT - lenP)

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
<div align="center">
 <img width="70%" width="1379" height="1406" alt="image" src="https://github.com/user-attachments/assets/d874a0a8-8642-4d3b-a35c-2c44bb0d9647" />
</div>
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

<details>
  <summary><b>1. File Explorer / Launcher</b></summary>
<br>

```go
// Quét thư mục home
files := scanDirectory("/home/user")
searcher := fuzzyvn.NewSearcher(files)

// User gõ, realtime searchautomatically
results := searcher.Search(userInput)
```

</details>

<details>
  <summary><b>2. Document Management</b></summary>
<br>
 
```go
// Index tài liệu công ty
docs := scanWithExtensions("/company/docs", []string{".pdf", ".docx"})
searcher := fuzzyvn.NewSearcher(docs)

// Tìm hợp đồng
contracts := searcher.Search("hop dong")
```

</details>

<details>
  <summary><b>3. Code Search</b></summary>
<br>
 
```go
// Index source code
code := scanIgnoreDirs("/project", []string{"node_modules", ".git"})
searcher := fuzzyvn.NewSearcher(code)

// Tìm file main
mains := searcher.Search("main")
```

</details>


<details>
  <summary><b>4. Media Library</b></summary>
<br>
 
```go
// Index nhạc
music := scanWithExtensions("/music", []string{".mp3", ".flac"})
searcher := fuzzyvn.NewSearcher(music)

// Tìm bài hát
songs := searcher.Search("son tung")
```

</details>

## Ví dụ nâng cao

<details>
  <summary><b>Rebuild Index khi file thay đổi</b></summary>
<br>
 
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

</details>

<details>
  <summary><b>Tùy chỉnh cho domain cụ thể</b></summary>
<br>
 
```go
searcher := fuzzyvn.NewSearcher(files)
cache := searcher.GetCache()

// Tăng boost cho người dùng power user
cache.SetBoostScore(15000)

// Lưu nhiều lịch sử hơn
cache.SetMaxQueries(1000)
```

</details>

<details>
  <summary><b>Integration với CLI tool</b></summary>
<br>
 
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

</details>

## Đóng góp

Vui lòng theo chuẩn [Contributing](.github/CONTRIBUTING.md) khi tạo một contribution qua pull request.

## Giấy phép

Package này được cấp phép bởi giấy phép [0BSD License](LICENSE). Bạn có thể sửa, xóa, thêm hay làm bất cứ thứ gì bạn muốn với nó.
