package fuzzyvn

import (
	"slices"
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Đường", "Duong"},
		{"đường", "duong"},
		{"Nguyễn", "Nguien"},
		{"nguyễn", "nguien"},
		{"Huệ", "Hue"},
		{"café", "cafe"},
		{"kỷ niệm", "ki niem"},
		{"kỉ niệm", "ki niem"},
		{"lý do", "li do"},
		{"lí do", "li do"},
		{"quy định", "qui dinh"},
		{"qui định", "qui dinh"},
		{"Sơn Tùng", "Son Tung"},
		{"Báo cáo tháng 1", "Bao cao thang 1"},
		{"Hello World", "Hello World"},
		{"", ""},
	}

	for _, tt := range tests {
		result := Normalize(tt.input)
		if result != tt.expected {
			t.Errorf("Normalize(%q) = %q, muốn %q", tt.input, result, tt.expected)
		}
	}
}

func TestNormalize_IY_Equivalence(t *testing.T) {
	pairs := []struct {
		a, b string
	}{
		{"kỷ niệm", "kỉ niệm"},
		{"lý do", "lí do"},
		{"quy tắc", "qui tắc"},
		{"ký tự", "kí tự"},
		{"mỹ thuật", "mĩ thuật"},
	}

	for _, pair := range pairs {
		normA := Normalize(pair.a)
		normB := Normalize(pair.b)
		if normA != normB {
			t.Errorf("Normalize(%q) = %q, Normalize(%q) = %q, phải giống nhau", pair.a, normA, pair.b, normB)
		}
	}
}

func TestLevenshteinRatio(t *testing.T) {
	tests := []struct {
		s1, s2   string
		expected int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "ab", 1},
		{"abc", "abcd", 1},
		{"main", "mian", 2},
		{"kitten", "sitting", 3},
		{"hello", "hallo", 1},
	}

	for _, tt := range tests {
		result := LevenshteinRatio(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("LevenshteinRatio(%q, %q) = %d, muốn %d", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestNewSearcher(t *testing.T) {
	files := []string{
		"/home/user/main.go",
		"/home/user/config.yaml",
		"/home/user/README.md",
	}

	searcher := NewSearcher(files)

	if len(searcher.Originals) != 3 {
		t.Errorf("Originals có %d phần tử, muốn 3", len(searcher.Originals))
	}

	if len(searcher.Normalized) != 3 {
		t.Errorf("Normalized có %d phần tử, muốn 3", len(searcher.Normalized))
	}

	if len(searcher.FilenamesOnly) != 3 {
		t.Errorf("FilenamesOnly có %d phần tử, muốn 3", len(searcher.FilenamesOnly))
	}

	if searcher.Cache == nil {
		t.Error("Cache không được khởi tạo")
	}
}

func TestSearcher_Search_Basic(t *testing.T) {
	files := []string{
		"/project/main.go",
		"/project/main_test.go",
		"/project/config.yaml",
		"/project/README.md",
	}

	searcher := NewSearcher(files)

	results := searcher.Search("main")
	if len(results) < 2 {
		t.Errorf("Search('main') trả về %d kết quả, muốn ít nhất 2", len(results))
	}

	if !slices.Contains(results, "/project/main.go") {
		t.Error("Search('main') không tìm thấy /project/main.go")
	}
}

func TestSearcher_Search_Vietnamese(t *testing.T) {
	files := []string{
		"/docs/Báo_cáo_tháng_1.pdf",
		"/docs/Hợp_đồng_thuê_nhà.docx",
		"/music/Sơn Tùng - Lạc Trôi.mp3",
		"/music/Mỹ Tâm - Đừng Hỏi Em.mp3",
	}

	searcher := NewSearcher(files)

	tests := []struct {
		query    string
		contains string
	}{
		{"bao cao", "Báo_cáo"},
		{"hop dong", "Hợp_đồng"},
		{"son tung", "Sơn Tùng"},
		{"lac troi", "Lạc Trôi"},
		{"my tam", "Mỹ Tâm"},
	}

	for _, tt := range tests {
		results := searcher.Search(tt.query)
		if len(results) == 0 {
			t.Errorf("Search(%q) không trả về kết quả", tt.query)
			continue
		}

		found := slices.ContainsFunc(results, func(r string) bool {
			return strings.Contains(r, tt.contains)
		})
		if !found {
			t.Errorf("Search(%q) không tìm thấy file chứa %q", tt.query, tt.contains)
		}
	}
}

func TestSearcher_Search_IY_Equivalence(t *testing.T) {
	files := []string{
		"/music/Kỷ Niệm Vô Tận - Vũ.flac",
		"/docs/Lý do nghỉ việc.docx",
	}

	searcher := NewSearcher(files)

	results1 := searcher.Search("ky niem")
	results2 := searcher.Search("ki niem")

	if len(results1) == 0 || len(results2) == 0 {
		t.Error("Search với i/y phải trả về kết quả")
	}

	if len(results1) != len(results2) {
		t.Errorf("Search('ky niem') và Search('ki niem') phải cho cùng số kết quả")
	}
}

func TestSearcher_Search_Typo(t *testing.T) {
	files := []string{
		"/project/main.go",
		"/project/config.yaml",
	}

	searcher := NewSearcher(files)

	results := searcher.Search("mian")
	if !slices.Contains(results, "/project/main.go") {
		t.Error("Search('mian') phải tìm thấy main.go (typo tolerance)")
	}
}

func TestQueryCache_RecordSelection(t *testing.T) {
	cache := NewQueryCache()

	cache.RecordSelection("main", "/project/main.go")

	if cache.Size() != 1 {
		t.Errorf("Cache size = %d, muốn 1", cache.Size())
	}

	cache.RecordSelection("main", "/project/main.go")

	if cache.Size() != 1 {
		t.Errorf("Cache size sau khi chọn lại = %d, muốn 1", cache.Size())
	}

	cache.RecordSelection("config", "/project/config.yaml")

	if cache.Size() != 2 {
		t.Errorf("Cache size = %d, muốn 2", cache.Size())
	}
}

func TestQueryCache_GetBoostScores_ExactMatch(t *testing.T) {
	cache := NewQueryCache()
	cache.RecordSelection("main", "/project/main.go")

	scores := cache.GetBoostScores("main")

	if score, exists := scores["/project/main.go"]; !exists || score == 0 {
		t.Error("GetBoostScores phải trả về score cho file đã cache")
	}
}

func TestQueryCache_GetBoostScores_SimilarQuery(t *testing.T) {
	cache := NewQueryCache()
	cache.RecordSelection("main server", "/project/main_server.go")

	tests := []string{"main", "main ser", "server"}

	for _, query := range tests {
		scores := cache.GetBoostScores(query)
		if score, exists := scores["/project/main_server.go"]; !exists || score == 0 {
			t.Errorf("GetBoostScores(%q) phải trả về score cho query tương tự", query)
		}
	}
}

func TestQueryCache_GetRecentQueries(t *testing.T) {
	cache := NewQueryCache()
	cache.RecordSelection("first", "/a.go")
	cache.RecordSelection("second", "/b.go")
	cache.RecordSelection("third", "/c.go")

	recent := cache.GetRecentQueries(2)

	if len(recent) != 2 {
		t.Errorf("GetRecentQueries(2) trả về %d, muốn 2", len(recent))
	}

	if recent[0] != "third" {
		t.Errorf("Query gần nhất = %q, muốn 'third'", recent[0])
	}
}

func TestQueryCache_GetCachedFiles(t *testing.T) {
	cache := NewQueryCache()
	cache.RecordSelection("main", "/project/main.go")
	cache.RecordSelection("main", "/project/main_test.go")

	files := cache.GetCachedFiles("main", 5)

	if len(files) != 2 {
		t.Errorf("GetCachedFiles trả về %d files, muốn 2", len(files))
	}
}

func TestQueryCache_GetAllRecentFiles(t *testing.T) {
	cache := NewQueryCache()
	cache.RecordSelection("query1", "/a.go")
	cache.RecordSelection("query2", "/b.go")
	cache.RecordSelection("query3", "/c.go")

	files := cache.GetAllRecentFiles(5)

	if len(files) != 3 {
		t.Errorf("GetAllRecentFiles trả về %d files, muốn 3", len(files))
	}

	if files[0] != "/c.go" {
		t.Errorf("File gần nhất = %q, muốn '/c.go'", files[0])
	}
}

func TestQueryCache_LRU_Eviction(t *testing.T) {
	cache := NewQueryCache()
	cache.SetMaxQueries(3)

	cache.RecordSelection("q1", "/a.go")
	cache.RecordSelection("q2", "/b.go")
	cache.RecordSelection("q3", "/c.go")
	cache.RecordSelection("q4", "/d.go")

	if cache.Size() != 3 {
		t.Errorf("Cache size sau eviction = %d, muốn 3", cache.Size())
	}

	scores := cache.GetBoostScores("q1")
	if len(scores) > 0 {
		t.Error("Query cũ nhất (q1) phải bị xóa")
	}
}

func TestQueryCache_Clear(t *testing.T) {
	cache := NewQueryCache()
	cache.RecordSelection("main", "/main.go")
	cache.RecordSelection("config", "/config.yaml")

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Cache size sau Clear = %d, muốn 0", cache.Size())
	}
}

func TestSearcher_RecordSelection_BoostsResults(t *testing.T) {
	files := []string{
		"/project/main.go",
		"/project/main_server.go",
		"/project/main_test.go",
		"/project/config.yaml",
	}

	searcher := NewSearcher(files)

	searcher.RecordSelection("main", "/project/main_test.go")
	searcher.RecordSelection("main", "/project/main_test.go")
	searcher.RecordSelection("main", "/project/main_test.go")

	results := searcher.Search("main")

	if len(results) == 0 {
		t.Fatal("Search không trả về kết quả")
	}

	if results[0] != "/project/main_test.go" {
		t.Errorf("File được cache nhiều lần phải ở đầu, got %q", results[0])
	}
}

func TestNewSearcherWithCache(t *testing.T) {
	files1 := []string{"/a.go", "/b.go"}
	searcher1 := NewSearcher(files1)
	searcher1.RecordSelection("test", "/a.go")

	cache := searcher1.GetCache()

	files2 := []string{"/a.go", "/b.go", "/c.go"}
	searcher2 := NewSearcherWithCache(files2, cache)

	if searcher2.Cache.Size() != 1 {
		t.Error("Cache phải được giữ lại khi dùng NewSearcherWithCache")
	}
}

func TestSearcher_ClearCache(t *testing.T) {
	files := []string{"/main.go"}
	searcher := NewSearcher(files)
	searcher.RecordSelection("main", "/main.go")

	searcher.ClearCache()

	if searcher.Cache.Size() != 0 {
		t.Error("ClearCache phải xóa hết cache")
	}
}
