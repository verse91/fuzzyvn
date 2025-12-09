package fuzzyvn

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/sahilm/fuzzy"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type CacheEntry struct {
	FilePath    string
	SelectCount int
}

type QueryCache struct {
	mu          sync.RWMutex
	entries     map[string][]CacheEntry
	queryOrder  []string
	maxQueries  int
	maxPerQuery int
	boostScore  int
}

func NewQueryCache() *QueryCache {
	return &QueryCache{
		entries:     make(map[string][]CacheEntry),
		queryOrder:  make([]string, 0),
		maxQueries:  100,
		maxPerQuery: 5,
		boostScore:  5000,
	}
}

func (c *QueryCache) SetMaxQueries(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxQueries = n
	c.evictIfNeeded()
}

func (c *QueryCache) SetBoostScore(score int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.boostScore = score
}

func (c *QueryCache) RecordSelection(query, filePath string) {
	if query == "" || filePath == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	queryNorm := strings.ToLower(Normalize(query))

	entries, exists := c.entries[queryNorm]
	if exists {
		for i, entry := range entries {
			if entry.FilePath == filePath {
				c.entries[queryNorm][i].SelectCount++
				c.moveToFront(queryNorm)
				return
			}
		}
	}

	newEntry := CacheEntry{FilePath: filePath, SelectCount: 1}
	if !exists {
		c.entries[queryNorm] = []CacheEntry{newEntry}
		c.queryOrder = append(c.queryOrder, queryNorm)
	} else {
		if len(entries) >= c.maxPerQuery {
			minIdx := 0
			minCount := entries[0].SelectCount
			for i, e := range entries {
				if e.SelectCount < minCount {
					minCount = e.SelectCount
					minIdx = i
				}
			}
			c.entries[queryNorm] = append(entries[:minIdx], entries[minIdx+1:]...)
		}
		c.entries[queryNorm] = append(c.entries[queryNorm], newEntry)
	}

	c.moveToFront(queryNorm)
	c.evictIfNeeded()
}

func (c *QueryCache) GetBoostScores(query string) map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]int)
	if query == "" {
		return result
	}

	queryNorm := strings.ToLower(Normalize(query))

	for cachedQuery, entries := range c.entries {
		similarity := c.querySimilarity(queryNorm, cachedQuery)
		if similarity > 0 {
			for _, entry := range entries {
				boost := (c.boostScore * similarity * entry.SelectCount) / 100
				if currentBoost, exists := result[entry.FilePath]; !exists || boost > currentBoost {
					result[entry.FilePath] = boost
				}
			}
		}
	}

	return result
}

func (c *QueryCache) querySimilarity(q1, q2 string) int {
	if q1 == q2 {
		return 100
	}

	if strings.HasPrefix(q2, q1) {
		return 70 + (30 * len(q1) / len(q2))
	}

	if strings.HasPrefix(q1, q2) && len(q2) >= 2 {
		return 50 + (30 * len(q2) / len(q1))
	}

	if len(q1) >= 2 && strings.Contains(q2, q1) {
		return 80
	}
	if len(q2) >= 2 && strings.Contains(q1, q2) {
		return 60
	}

	words1 := strings.Fields(q1)
	words2 := strings.Fields(q2)
	if len(words1) > 0 && len(words2) > 0 {
		commonWords := 0
		for _, w1 := range words1 {
			for _, w2 := range words2 {
				if w1 == w2 && len(w1) >= 2 {
					commonWords++
					break
				}
			}
		}
		if commonWords > 0 {
			return 50 + (commonWords * 15)
		}
	}

	if len(q1) >= 3 && len(q2) >= 3 {
		dist := LevenshteinRatio(q1, q2)
		maxLen := len(q1)
		if len(q2) > maxLen {
			maxLen = len(q2)
		}
		threshold := maxLen * 30 / 100
		if threshold < 2 {
			threshold = 2
		}
		if dist <= threshold {
			return 60 - (dist * 10)
		}
	}

	return 0
}

func (c *QueryCache) moveToFront(query string) {
	for i, q := range c.queryOrder {
		if q == query {
			c.queryOrder = append(c.queryOrder[:i], c.queryOrder[i+1:]...)
			break
		}
	}
	c.queryOrder = append(c.queryOrder, query)
}

func (c *QueryCache) evictIfNeeded() {
	for len(c.queryOrder) > c.maxQueries {
		oldestQuery := c.queryOrder[0]
		c.queryOrder = c.queryOrder[1:]
		delete(c.entries, oldestQuery)
	}
}

func (c *QueryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string][]CacheEntry)
	c.queryOrder = make([]string, 0)
}

func (c *QueryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

func (c *QueryCache) GetRecentQueries(limit int) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if limit <= 0 || len(c.queryOrder) == 0 {
		return []string{}
	}

	result := make([]string, 0, limit)
	for i := len(c.queryOrder) - 1; i >= 0 && len(result) < limit; i-- {
		result = append(result, c.queryOrder[i])
	}
	return result
}

func (c *QueryCache) GetCachedFiles(query string, limit int) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if query == "" || limit <= 0 {
		return []string{}
	}

	queryNorm := strings.ToLower(Normalize(query))

	type fileScore struct {
		path  string
		score int
	}
	var matches []fileScore

	for cachedQuery, entries := range c.entries {
		similarity := c.querySimilarity(queryNorm, cachedQuery)
		if similarity > 0 {
			for _, entry := range entries {
				score := similarity * entry.SelectCount
				matches = append(matches, fileScore{path: entry.FilePath, score: score})
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	seen := make(map[string]bool)
	result := make([]string, 0, limit)
	for _, m := range matches {
		if !seen[m.path] {
			seen[m.path] = true
			result = append(result, m.path)
			if len(result) >= limit {
				break
			}
		}
	}

	return result
}

func (c *QueryCache) GetAllRecentFiles(limit int) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if limit <= 0 {
		return []string{}
	}

	type fileInfo struct {
		path       string
		queryIndex int
		count      int
	}
	fileMap := make(map[string]*fileInfo)

	for i, query := range c.queryOrder {
		entries := c.entries[query]
		for _, entry := range entries {
			if existing, ok := fileMap[entry.FilePath]; ok {
				if i > existing.queryIndex {
					existing.queryIndex = i
				}
				existing.count += entry.SelectCount
			} else {
				fileMap[entry.FilePath] = &fileInfo{
					path:       entry.FilePath,
					queryIndex: i,
					count:      entry.SelectCount,
				}
			}
		}
	}

	files := make([]*fileInfo, 0, len(fileMap))
	for _, f := range fileMap {
		files = append(files, f)
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].queryIndex != files[j].queryIndex {
			return files[i].queryIndex > files[j].queryIndex
		}
		return files[i].count > files[j].count
	})

	result := make([]string, 0, limit)
	for i := 0; i < len(files) && i < limit; i++ {
		result = append(result, files[i].path)
	}

	return result
}

type Searcher struct {
	Originals     []string
	Normalized    []string
	FilenamesOnly []string
	Cache         *QueryCache
}

type MatchResult struct {
	Str   string
	Score int
}

func NewSearcher(items []string) *Searcher {
	normPaths := make([]string, len(items))
	normNames := make([]string, len(items))

	for i, item := range items {
		filename := filepath.Base(item)
		priorityString := filename + " " + item
		normPaths[i] = strings.ToLower(Normalize(priorityString))
		normNames[i] = strings.ToLower(Normalize(filename))
	}

	return &Searcher{
		Originals:     items,
		Normalized:    normPaths,
		FilenamesOnly: normNames,
		Cache:         NewQueryCache(),
	}
}

func NewSearcherWithCache(items []string, cache *QueryCache) *Searcher {
	s := NewSearcher(items)
	if cache != nil {
		s.Cache = cache
	}
	return s
}

func (s *Searcher) RecordSelection(query, filePath string) {
	if s.Cache != nil {
		s.Cache.RecordSelection(query, filePath)
	}
}

func (s *Searcher) ClearCache() {
	if s.Cache != nil {
		s.Cache.Clear()
	}
}

func (s *Searcher) GetCache() *QueryCache {
	return s.Cache
}

func (s *Searcher) Search(query string) []string {
	queryNorm := strings.ToLower(Normalize(query))
	queryRunes := []rune(queryNorm)
	queryLen := len(queryRunes)

	uniqueResults := make(map[int]int)
	filePathToIdx := make(map[string]int)

	for i, fp := range s.Originals {
		filePathToIdx[fp] = i
	}

	var cacheBoosts map[string]int
	if s.Cache != nil {
		cacheBoosts = s.Cache.GetBoostScores(query)
	}

	matches := fuzzy.Find(queryNorm, s.Normalized)
	for _, m := range matches {
		uniqueResults[m.Index] = m.Score
	}

	if queryLen > 1 {
		baseThreshold := (queryLen / 3) + 1
		if baseThreshold < 3 {
			baseThreshold = 3
		}

		for i, nameNorm := range s.FilenamesOnly {
			runesName := []rune(nameNorm)

			if len(runesName) < queryLen {
				continue
			}

			dist := 100

			targetStr1 := string(runesName[:queryLen])
			d1 := LevenshteinRatio(queryNorm, targetStr1)
			dist = d1

			if len(runesName) > queryLen {
				targetStr2 := string(runesName[:queryLen+1])
				d2 := LevenshteinRatio(queryNorm, targetStr2)
				if d2 < dist {
					dist = d2
				}
			}

			if dist < baseThreshold {
				score := 10000 - (dist * 100)

				lenDiff := len(runesName) - queryLen
				if lenDiff > 0 {
					score -= (lenDiff / 2)
				}

				if oldScore, exists := uniqueResults[i]; !exists || score > oldScore {
					uniqueResults[i] = score
				}
			}
		}
	}

	for cachedPath, boost := range cacheBoosts {
		if idx, exists := filePathToIdx[cachedPath]; exists {
			if _, alreadyInResults := uniqueResults[idx]; !alreadyInResults {
				uniqueResults[idx] = boost
			}
		}
	}

	var rankedResults []MatchResult
	for idx, score := range uniqueResults {
		filePath := s.Originals[idx]
		finalScore := score

		if boost, exists := cacheBoosts[filePath]; exists {
			if score != boost {
				finalScore += boost
			}
		}

		rankedResults = append(rankedResults, MatchResult{
			Str:   filePath,
			Score: finalScore,
		})
	}

	sort.SliceStable(rankedResults, func(i, j int) bool {
		if rankedResults[i].Score == rankedResults[j].Score {
			return len(rankedResults[i].Str) < len(rankedResults[j].Str)
		}
		return rankedResults[i].Score > rankedResults[j].Score
	})

	var results []string
	limit := 20
	if len(rankedResults) < limit {
		limit = len(rankedResults)
	}
	for _, res := range rankedResults[:limit] {
		results = append(results, res.Str)
	}
	return results
}

func Normalize(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, _ := transform.String(t, s)
	output = strings.ReplaceAll(output, "đ", "d")
	output = strings.ReplaceAll(output, "Đ", "D")
	return output
}

func LevenshteinRatio(s1, s2 string) int {
	s1Len := len(s1)
	s2Len := len(s2)
	column := make([]int, len(s1)+1)
	for y := 1; y <= s1Len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2Len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1Len; y++ {
			oldkey := column[y]
			var incr int
			if s1[y-1] != s2[x-1] {
				incr = 1
			}
			minVal := column[y] + 1
			if column[y-1]+1 < minVal {
				minVal = column[y-1] + 1
			}
			if lastkey+incr < minVal {
				minVal = lastkey + incr
			}
			column[y] = minVal
			lastkey = oldkey
		}
	}
	return column[s1Len]
}
