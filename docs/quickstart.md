# Quick Start

## Basic Usage

```go
package main

import (
    "fmt"
    "github.com/verse91/fuzzyvn"
)

func main() {
    // 1. Prepare file list
    files := []string{
        "/home/user/documents/báo_cáo_tháng_1.pdf",
        "/home/user/documents/hợp_đồng_2024.docx",
        "/home/user/music/Sơn Tùng - Chạy Ngay Đi.mp3",
    }

    // 2. Create searcher
    searcher := fuzzyvn.NewSearcher(files)

    // 3. Search
    results := searcher.Search("bao cao")
    
    for _, r := range results {
        fmt.Println(r)
    }
}
```

## With Caching

```go
// Search and get results
results := searcher.Search("son tung")

// User selects a file -> record it
selectedFile := results[0]
searcher.RecordSelection("son tung", selectedFile)

// Next time user searches similar query, 
// the selected file will be boosted to top
results = searcher.Search("sontung")  // typo still works
results = searcher.Search("son")      // partial still works
```

## Preserving Cache Across Rebuilds

```go
// Get cache before rebuild
cache := searcher.GetCache()

// Rebuild with new file list (e.g., after file system changes)
newFiles := scanDirectory("/home/user")
searcher = fuzzyvn.NewSearcherWithCache(newFiles, cache)

// Cache is preserved!
```

