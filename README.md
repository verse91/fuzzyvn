# FuzzyVN Documentation

FuzzyVN is a Vietnamese-optimized fuzzy file finder library for Go. It combines multiple search algorithms with intelligent caching to provide fast, accurate file search results.

## Table of Contents

- [Installation](./docs/installation.md)
- [Quick Start](./docs/quickstart.md)
- [API Reference](./docs/api.md)
- [Cache System](./docs/cache.md)
- [Search Algorithm](./docs/algorithm.md)
- [Examples](./docs/examples.md)

## Features

- **Vietnamese Support**: Handles Vietnamese diacritics (converts "Đường" to "Duong")
- **Multi-Algorithm Search**: Combines fuzzy matching + Levenshtein distance
- **Smart Caching**: Learns from user selections to boost relevant results
- **Typo Tolerance**: Handles common typing errors
- **Thread-Safe**: Safe for concurrent access

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Searcher                           │
├─────────────────────────────────────────────────────────┤
│  Originals[]     - Original file paths                  │
│  Normalized[]    - Normalized for fuzzy search          │
│  FilenamesOnly[] - Filenames only for Levenshtein       │
│  Cache           - Query cache for boosting             │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    QueryCache                           │
├─────────────────────────────────────────────────────────┤
│  entries{}       - query → []CacheEntry                 │
│  queryOrder[]    - LRU ordering                         │
│  maxQueries      - Maximum cached queries (100)         │
│  maxPerQuery     - Max files per query (5)              │
│  boostScore      - Boost multiplier (5000)              │
└─────────────────────────────────────────────────────────┘
```

## Search Flow

```
User Query
    │
    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Fuzzy Match  │ +  │ Levenshtein  │ +  │ Cache Boost  │
│ (substring)  │    │ (typo fix)   │    │ (history)    │
└──────────────┘    └──────────────┘    └──────────────┘
    │                      │                    │
    └──────────────────────┼────────────────────┘
                           ▼
                    Merged Results
                           │
                           ▼
                    Sort by Score
                           │
                           ▼
                    Top 20 Results
```

