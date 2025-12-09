# API Reference

## Searcher

### Types

```go
type Searcher struct {
    Originals     []string     // Original file paths
    Normalized    []string     // Normalized strings for fuzzy search
    FilenamesOnly []string     // Filenames only for Levenshtein
    Cache         *QueryCache  // Query cache
}

type MatchResult struct {
    Str   string  // File path
    Score int     // Match score
}
```

### Functions

#### NewSearcher

```go
func NewSearcher(items []string) *Searcher
```

Creates a new Searcher with the given file paths.

**Parameters:**
- `items`: List of file paths to index

**Returns:**
- `*Searcher`: New searcher instance with empty cache

---

#### NewSearcherWithCache

```go
func NewSearcherWithCache(items []string, cache *QueryCache) *Searcher
```

Creates a new Searcher with an existing cache.

**Parameters:**
- `items`: List of file paths to index
- `cache`: Existing QueryCache to reuse

**Returns:**
- `*Searcher`: New searcher instance with provided cache

---

### Methods

#### Search

```go
func (s *Searcher) Search(query string) []string
```

Searches for files matching the query.

**Parameters:**
- `query`: Search query (Vietnamese diacritics are handled automatically)

**Returns:**
- `[]string`: Top 20 matching file paths, sorted by relevance

---

#### RecordSelection

```go
func (s *Searcher) RecordSelection(query, filePath string)
```

Records that a user selected a file for a query. This boosts the file in future similar searches.

**Parameters:**
- `query`: The search query
- `filePath`: The file that was selected

---

#### ClearCache

```go
func (s *Searcher) ClearCache()
```

Clears all cached query-file associations.

---

#### GetCache

```go
func (s *Searcher) GetCache() *QueryCache
```

Returns the cache for reuse when rebuilding the searcher.

---

## QueryCache

### Types

```go
type QueryCache struct {
    // Internal fields (not exported)
}

type CacheEntry struct {
    FilePath    string  // Cached file path
    SelectCount int     // Number of times selected
}
```

### Functions

#### NewQueryCache

```go
func NewQueryCache() *QueryCache
```

Creates a new empty cache with default settings:
- `maxQueries`: 100
- `maxPerQuery`: 5
- `boostScore`: 5000

---

### Methods

#### SetMaxQueries

```go
func (c *QueryCache) SetMaxQueries(n int)
```

Sets the maximum number of queries to cache (LRU eviction).

---

#### SetBoostScore

```go
func (c *QueryCache) SetBoostScore(score int)
```

Sets the base boost score for cached results.

---

#### RecordSelection

```go
func (c *QueryCache) RecordSelection(query, filePath string)
```

Records a query-file selection.

---

#### GetBoostScores

```go
func (c *QueryCache) GetBoostScores(query string) map[string]int
```

Returns boost scores for files matching similar queries.

**Returns:**
- `map[string]int`: File path → boost score

---

#### GetCachedFiles

```go
func (c *QueryCache) GetCachedFiles(query string, limit int) []string
```

Returns cached files for similar queries.

**Parameters:**
- `query`: Current search query
- `limit`: Maximum files to return

**Returns:**
- `[]string`: Cached file paths sorted by relevance

---

#### GetRecentQueries

```go
func (c *QueryCache) GetRecentQueries(limit int) []string
```

Returns the most recent queries (MRU order).

---

#### GetAllRecentFiles

```go
func (c *QueryCache) GetAllRecentFiles(limit int) []string
```

Returns all recently cached files regardless of query.

---

#### Size

```go
func (c *QueryCache) Size() int
```

Returns the number of cached queries.

---

#### Clear

```go
func (c *QueryCache) Clear()
```

Clears all cache entries.

---

## Utility Functions

#### Normalize

```go
func Normalize(s string) string
```

Normalizes Vietnamese text by removing diacritics.

**Examples:**
- `"Đường"` → `"Duong"`
- `"Nguyễn"` → `"Nguyen"`
- `"café"` → `"cafe"`

---

#### LevenshteinRatio

```go
func LevenshteinRatio(s1, s2 string) int
```

Calculates the Levenshtein edit distance between two strings.

**Returns:**
- `int`: Number of edits (insertions, deletions, substitutions)

