# Cache System

The cache system learns from user behavior to improve search results over time.

## How It Works

### Recording Selections

When a user selects a file from search results:

```go
searcher.RecordSelection("main server", "/project/src/main_server.go")
```

The cache stores:
- Query: `"main server"` (normalized)
- File: `/project/src/main_server.go`
- Count: 1 (increments on repeated selections)

### Query Similarity

The cache doesn't require exact query matches. It uses similarity scoring:

| Match Type | Score | Example |
|------------|-------|---------|
| Exact match | 100 | `"main"` = `"main"` |
| Substring (query in cached) | 80 | `"main"` ⊂ `"main server"` |
| Prefix match | 70-100 | `"mai"` → `"main"` |
| Reverse prefix | 50-80 | `"main server"` → `"main"` |
| Word overlap | 50-95 | `"server main"` ↔ `"main server"` |
| Fuzzy (≤30% errors) | 30-60 | `"mian"` ≈ `"main"` |

### Boost Calculation

```
boost = boostScore × similarity × selectCount / 100
```

Default `boostScore` is 5000, so:
- Exact match, selected 1 time: `5000 × 100 × 1 / 100 = 5000`
- Exact match, selected 3 times: `5000 × 100 × 3 / 100 = 15000`
- Prefix match (80%), selected 2 times: `5000 × 80 × 2 / 100 = 8000`

## Configuration

### Max Queries

```go
cache := searcher.GetCache()
cache.SetMaxQueries(200)  // Default: 100
```

When limit is exceeded, oldest queries are evicted (LRU).

### Max Files Per Query

Each query stores up to 5 files by default. When exceeded, the file with lowest `selectCount` is removed.

### Boost Score

```go
cache.SetBoostScore(10000)  // Default: 5000
```

Higher values make cached results more prominent.

## LRU Eviction

Queries are ordered by recency:
- Each selection moves the query to "most recent"
- When `maxQueries` is exceeded, oldest queries are removed

## Persistence

The cache is in-memory only. To persist across restarts:

```go
// Before shutdown
cache := searcher.GetCache()
// Serialize cache.entries and cache.queryOrder to JSON/gob

// On startup
// Deserialize and create new cache
searcher = fuzzyvn.NewSearcherWithCache(files, loadedCache)
```

## Thread Safety

All cache operations are protected by `sync.RWMutex`:
- Read operations: `RLock`
- Write operations: `Lock`

Safe for concurrent access from multiple goroutines.

