# Examples

## Basic File Search

```go
package main

import (
    "fmt"
    "io/fs"
    "path/filepath"
    
    "github.com/verse91/fuzzyvn"
)

func main() {
    // Scan directory
    var files []string
    filepath.WalkDir("/home/user", func(path string, d fs.DirEntry, err error) error {
        if err == nil && !d.IsDir() {
            files = append(files, path)
        }
        return nil
    })

    // Create searcher
    searcher := fuzzyvn.NewSearcher(files)

    // Search
    results := searcher.Search("readme")
    for _, r := range results {
        fmt.Println(r)
    }
}
```

## HTTP Server with Caching

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

## Neovim Integration (Lua)

```lua
local fuzzyvn = require("fuzzyvn")

-- Initialize
local searcher = fuzzyvn.new_searcher(vim.fn.glob("**/*", false, true))

-- Search function
local function search(query)
    return searcher:search(query)
end

-- Record selection
local function on_select(query, path)
    searcher:record_selection(query, path)
    vim.cmd("edit " .. path)
end

-- Telescope integration
require("telescope").setup({
    defaults = {
        finder = function(query)
            return search(query)
        end,
        attach_mappings = function(_, map)
            map("i", "<CR>", function(prompt_bufnr)
                local selection = action_state.get_selected_entry()
                on_select(query, selection.value)
            end)
            return true
        end,
    },
})
```

## Custom Scoring

```go
// Get cache and customize
cache := searcher.GetCache()

// Increase boost for cached results
cache.SetBoostScore(10000)  // Default: 5000

// Keep more queries in cache
cache.SetMaxQueries(500)    // Default: 100
```

## Rebuilding Index

```go
// File watcher detected changes
func onFileSystemChange() {
    // Preserve cache
    cache := searcher.GetCache()
    
    // Rescan files
    newFiles := scanDirectory("/data")
    
    // Rebuild with same cache
    mu.Lock()
    searcher = fuzzyvn.NewSearcherWithCache(newFiles, cache)
    mu.Unlock()
}
```

## Getting Recent Activity

```go
cache := searcher.GetCache()

// Get recent search queries
recentQueries := cache.GetRecentQueries(10)
// ["main.go", "config", "readme", ...]

// Get recently selected files
recentFiles := cache.GetAllRecentFiles(5)
// ["/project/main.go", "/project/config.yaml", ...]

// Get cached files for current query
cachedForQuery := cache.GetCachedFiles("main", 5)
// Files previously selected for "main" or similar queries
```

## Vietnamese Text Examples

```go
searcher := fuzzyvn.NewSearcher([]string{
    "/docs/Báo_cáo_tháng_1.pdf",
    "/docs/Hợp_đồng_thuê_nhà.docx",
    "/music/Sơn Tùng - Lạc Trôi.mp3",
})

// All these queries work:
searcher.Search("bao cao")      // matches "Báo_cáo"
searcher.Search("hop dong")     // matches "Hợp_đồng"
searcher.Search("son tung")     // matches "Sơn Tùng"
searcher.Search("lac troi")     // matches "Lạc Trôi"

// Typos also work:
searcher.Search("bao coa")      // typo: "coa" → "cao"
searcher.Search("sontung")      // missing space
searcher.Search("sont ung")     // wrong space
```

