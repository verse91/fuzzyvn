# Search Algorithm

FuzzyVN uses a multi-stage search algorithm combining three techniques.

## Stage 1: Fuzzy Matching

Uses `github.com/sahilm/fuzzy` for substring matching.

**Input:** Normalized query against normalized file paths

**Scoring:** Based on character match positions and gaps

**Example:**
```
Query: "main"
File:  "project/src/main_server.go"
Match: Score based on "main" appearing in filename
```

## Stage 2: Levenshtein Distance

Handles typos and similar filenames.

### Configuration

```go
threshold = (queryLength / 3) + 1
if threshold < 3 {
    threshold = 3
}
```

| Query Length | Allowed Errors |
|--------------|----------------|
| 1-8 chars | 3 errors |
| 9-11 chars | 4 errors |
| 12-14 chars | 5 errors |
| etc. | +1 per 3 chars |

### Double Check Technique

Compares query against two slices of filename:

1. **Exact length**: `filename[:queryLen]`
2. **Extended length**: `filename[:queryLen+1]`

This handles both:
- Typos: `"mian"` → `"main"` (same length)
- Missing chars: `"main"` → `"maain"` (extended match)

### Scoring

```go
score = 10000 - (distance × 100) - (lengthDiff / 2)
```

- Base score: 10000
- Penalty per edit: -100
- Penalty for longer filename: `-lengthDiff/2`

## Stage 3: Cache Injection

Files from cache are injected even if they don't match fuzzy/Levenshtein.

```go
for cachedPath, boost := range cacheBoosts {
    if idx, exists := filePathToIdx[cachedPath]; exists {
        if _, alreadyInResults := uniqueResults[idx]; !alreadyInResults {
            uniqueResults[idx] = boost
        }
    }
}
```

## Score Aggregation

Final score combines all stages:

```go
finalScore = fuzzyScore + cacheBoost
```

If file was found by Levenshtein only:
```go
finalScore = levenshteinScore + cacheBoost
```

If file was injected from cache only:
```go
finalScore = cacheBoost
```

## Sorting

Results are sorted by:
1. Score (descending)
2. Path length (ascending) for tie-breaking

```go
sort.SliceStable(results, func(i, j int) bool {
    if results[i].Score == results[j].Score {
        return len(results[i].Str) < len(results[j].Str)
    }
    return results[i].Score > results[j].Score
})
```

## Vietnamese Normalization

Before any matching, text is normalized:

1. **NFD decomposition**: Separate base char from diacritics
2. **Remove diacritics**: Strip combining marks
3. **NFC composition**: Recompose
4. **Handle Đ/đ**: Replace with D/d

```go
"Đường Nguyễn Huệ" → "Duong Nguyen Hue"
```

## Performance

| Operation | Complexity |
|-----------|------------|
| Fuzzy match | O(n × m) where n=files, m=query |
| Levenshtein | O(n × k²) where k=query length |
| Cache lookup | O(q × f) where q=cached queries, f=files per query |
| Sort | O(r log r) where r=results |

For 10,000 files with average query, typical response time: 5-20ms.

