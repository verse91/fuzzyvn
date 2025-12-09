# Bắt đầu nhanh

## Cách dùng cơ bản

```go
package main

import (
    "fmt"
    "io/fs"
    "path/filepath"
    
    "github.com/verse91/fuzzyvn"
)

func main() {
    // 1. Quét thư mục lấy danh sách file
    var files []string
    filepath.WalkDir("/home/user/documents", func(path string, d fs.DirEntry, err error) error {
        if err == nil && !d.IsDir() {
            files = append(files, path)
        }
        return nil
    })

    // 2. Tạo searcher
    searcher := fuzzyvn.NewSearcher(files)

    // 3. Tìm kiếm (hỗ trợ tiếng Việt không dấu)
    results := searcher.Search("bao cao")
    
    for _, r := range results {
        fmt.Println(r)
    }
}
```

## Sử dụng Cache

```go
// Tìm kiếm và lấy kết quả
results := searcher.Search("son tung")

// Người dùng chọn file -> ghi nhận vào cache
selectedFile := results[0]
searcher.RecordSelection("son tung", selectedFile)

// Lần sau tìm query tương tự, file đã chọn sẽ được đẩy lên đầu
results = searcher.Search("sontung")  // gõ sai vẫn được
results = searcher.Search("son")      // gõ một phần vẫn được
```

## Giữ Cache khi rebuild

```go
// Lấy cache trước khi rebuild
cache := searcher.GetCache()

// Rebuild với danh sách file mới (vd: sau khi file system thay đổi)
newFiles := scanDirectory("/home/user")
searcher = fuzzyvn.NewSearcherWithCache(newFiles, cache)

// Cache vẫn được giữ nguyên!
```

## Quét thư mục với filter

```go
// Chỉ lấy một số đuôi file
func scanWithFilter(root string, extensions []string) []string {
    var files []string
    extMap := make(map[string]bool)
    for _, ext := range extensions {
        extMap[ext] = true
    }

    filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil || d.IsDir() {
            return nil
        }
        
        ext := filepath.Ext(path)
        if len(extensions) == 0 || extMap[ext] {
            files = append(files, path)
        }
        return nil
    })

    return files
}

// Sử dụng
files := scanWithFilter("/home/user", []string{".go", ".md", ".txt"})
searcher := fuzzyvn.NewSearcher(files)
```

## Bỏ qua thư mục

```go
func scanIgnore(root string, ignoreDirs []string) []string {
    var files []string
    ignore := make(map[string]bool)
    for _, dir := range ignoreDirs {
        ignore[dir] = true
    }

    filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        
        // Bỏ qua thư mục trong danh sách ignore
        if d.IsDir() && ignore[d.Name()] {
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
files := scanIgnore("/project", []string{"node_modules", ".git", "vendor"})
searcher := fuzzyvn.NewSearcher(files)
```

