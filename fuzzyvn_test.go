package fuzzyvn

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"
)

const TestData100kPath = "./demo/test_data/test_paths_100k.txt"

func loadTestFiles(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	files := strings.Split(string(content), "\n")
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}
	return result, nil
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Đường", "Duong"},
		{"đường", "duong"},
		{"Nguyễn", "Nguyen"},
		{"nguyễn", "nguyen"},
		{"Huệ", "Hue"},
		{"café", "cafe"},
		{"kỷ niệm", "ky niem"},
		{"kỉ niệm", "ki niem"},
		{"lý do", "ly do"},
		{"lí do", "li do"},
		{"quy định", "quy dinh"},
		{"qui định", "qui dinh"},
		{"Sơn Tùng", "Son Tung"},
		{"Báo cáo tháng 1", "Bao cao thang 1"},
		{"Hello World", "Hello World"},
		{"vật lý", "vat ly"},
		{"Python", "Python"},
		{"", ""},
	}

	for _, tt := range tests {
		result := Normalize(tt.input)
		if result != tt.expected {
			t.Errorf("Normalize(%q) = %q, muốn %q", tt.input, result, tt.expected)
		}
	}
}

func TestNormalize_YI_NotEquivalent(t *testing.T) {
	pairs := []struct {
		a, b    string
		expectA string
		expectB string
	}{
		{"kỷ niệm", "kỉ niệm", "ky niem", "ki niem"},
		{"lý do", "lí do", "ly do", "li do"},
		{"vật lý", "vật lí", "vat ly", "vat li"},
	}

	for _, pair := range pairs {
		normA := Normalize(pair.a)
		normB := Normalize(pair.b)
		if normA != pair.expectA {
			t.Errorf("Normalize(%q) = %q, muốn %q", pair.a, normA, pair.expectA)
		}
		if normB != pair.expectB {
			t.Errorf("Normalize(%q) = %q, muốn %q", pair.b, normB, pair.expectB)
		}
		if normA == normB {
			t.Errorf("Normalize(%q) = %q KHÔNG nên bằng Normalize(%q) = %q", pair.a, normA, pair.b, normB)
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

func TestSearchWithVietnameseData(t *testing.T) {
	cases := []struct {
		pattern     string
		data        []string
		wantMatches int
		wantFirst   string
	}{
		{
			"bao cao",
			[]string{"Báo_cáo_tháng_1.pdf", "config.yaml", "README.md"},
			1,
			"Báo_cáo_tháng_1.pdf",
		},
		{
			"son tung",
			[]string{"Sơn Tùng - Lạc Trôi.mp3", "Mỹ Tâm - Hãy Về Đây.mp3"},
			1,
			"Sơn Tùng - Lạc Trôi.mp3",
		},
		{
			"ky niem",
			[]string{"Kỷ Niệm Vô Tận.flac", "Kỉ Niệm Xưa.mp3", "config.yaml"},
			2,
			"",
		},
		{
			"ki niem",
			[]string{"Kỷ Niệm Vô Tận.flac", "Kỉ Niệm Xưa.mp3", "config.yaml"},
			2,
			"",
		},
		{
			"duong",
			[]string{"Đường Về Nhà.mp3", "duong_dan.txt", "config.yaml"},
			2,
			"",
		},
		{
			"nguyen",
			[]string{"Nguyễn Văn A.docx", "nguyen_config.yaml"},
			2,
			"",
		},
	}

	for _, c := range cases {
		searcher := NewSearcher(c.data)
		results := searcher.Search(c.pattern)

		if len(results) < c.wantMatches {
			t.Errorf("Search(%q): got %d matches, want at least %d", c.pattern, len(results), c.wantMatches)
		}

		if c.wantFirst != "" && len(results) > 0 && results[0] != c.wantFirst {
			t.Errorf("Search(%q): first match = %q, want %q", c.pattern, results[0], c.wantFirst)
		}
	}
}

func TestSearchWithTypos(t *testing.T) {
	cases := []struct {
		pattern string
		data    []string
		want    string
	}{
		{"mian", []string{"main.go", "config.yaml"}, "main.go"},
		{"conifg", []string{"main.go", "config.yaml"}, "config.yaml"},
		{"redame", []string{"README.md", "main.go"}, "README.md"},
		{"maiin", []string{"main.go", "test.go"}, "main.go"},
	}

	for _, c := range cases {
		searcher := NewSearcher(c.data)
		results := searcher.Search(c.pattern)

		if !slices.Contains(results, c.want) {
			t.Errorf("Search(%q) với typo: không tìm thấy %q", c.pattern, c.want)
		}
	}
}

func TestSearchCacheBoost(t *testing.T) {
	files := []string{
		"/project/main.go",
		"/project/main_test.go",
		"/project/main_server.go",
		"/project/main_client.go",
		"/project/config.yaml",
	}

	searcher := NewSearcher(files)

	results1 := searcher.Search("main")
	firstBefore := results1[0]

	searcher.RecordSelection("main", "/project/main_client.go")
	searcher.RecordSelection("main", "/project/main_client.go")
	searcher.RecordSelection("main", "/project/main_client.go")

	results2 := searcher.Search("main")
	firstAfter := results2[0]

	if firstAfter != "/project/main_client.go" {
		t.Errorf("Sau khi cache, file được chọn nhiều lần phải lên đầu. Got %q", firstAfter)
	}

	if firstBefore == firstAfter {
		t.Log("Lưu ý: kết quả có thể giống nhau nếu file đã ở đầu")
	}
}

func TestSearchWithRealworldData(t *testing.T) {
	t.Run("với test_data music", func(t *testing.T) {
		testDataPath := "demo/test_data/Music_Lossless"
		if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
			t.Skip("Bỏ qua: không có test_data")
		}

		files := scanTestData(testDataPath)
		if len(files) == 0 {
			t.Skip("Bỏ qua: không có files trong test_data")
		}

		searcher := NewSearcher(files)

		cases := []struct {
			pattern     string
			wantMatches int
		}{
			{"son tung", 1},
			{"ky niem", 1},
			{"lac troi", 1},
		}

		for _, c := range cases {
			now := time.Now()
			results := searcher.Search(c.pattern)
			elapsed := time.Since(now)

			fmt.Printf("Search '%s' trong %d files... tìm thấy %d kết quả trong %v\n",
				c.pattern, len(files), len(results), elapsed)

			if len(results) < c.wantMatches {
				t.Errorf("Search(%q): got %d matches, want at least %d", c.pattern, len(results), c.wantMatches)
			}
		}
	})
}

func BenchmarkSearch_RealWorld(b *testing.B) {
	allFiles, err := loadTestFiles(TestData100kPath)
	if err != nil {
		b.Skipf("Lỗi khi tải file %s. Vui lòng đảm bảo file tồn tại: %v", TestData100kPath, err)
	}

	files50k := allFiles[:50000]

	b.Run("Search/50k_Files", func(b *testing.B) {
		searcher := NewSearcher(files50k)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("config")
		}
	})

	b.Run("Search/100k_Files", func(b *testing.B) {
		searcher := NewSearcher(allFiles)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("config")
		}
	})

	b.Run("Search/100k_Files_Typo", func(b *testing.B) {
		searcher := NewSearcher(allFiles)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("conifg")
		}
	})
}

func BenchmarkNewSearcher(b *testing.B) {
	files := generateTestFiles(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewSearcher(files)
	}
}

func BenchmarkSearch(b *testing.B) {
	b.Run("100 files", func(b *testing.B) {
		files := generateTestFiles(100)
		searcher := NewSearcher(files)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("main")
		}
	})

	b.Run("1000 files", func(b *testing.B) {
		files := generateTestFiles(1000)
		searcher := NewSearcher(files)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("main")
		}
	})

	b.Run("10000 files", func(b *testing.B) {
		files := generateTestFiles(10000)
		searcher := NewSearcher(files)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("config")
		}
	})
}

func BenchmarkSearchVietnamese(b *testing.B) {
	files := generateVietnameseTestFiles(1000)
	searcher := NewSearcher(files)

	b.Run("tiếng Việt có dấu", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("báo cáo")
		}
	})

	b.Run("tiếng Việt không dấu", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			searcher.Search("bao cao")
		}
	})
}

func BenchmarkSearchWithCache(b *testing.B) {
	files := generateTestFiles(1000)
	searcher := NewSearcher(files)

	searcher.RecordSelection("main", files[0])
	searcher.RecordSelection("main", files[1])
	searcher.RecordSelection("config", files[2])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searcher.Search("main")
	}
}

func BenchmarkNormalize(b *testing.B) {
	testStrings := []string{
		"Đường Nguyễn Huệ",
		"Báo cáo tháng 1",
		"Sơn Tùng M-TP",
		"Kỷ niệm vô tận",
		"Hello World",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range testStrings {
			Normalize(s)
		}
	}
}

func BenchmarkLevenshteinRatio(b *testing.B) {
	pairs := []struct{ a, b string }{
		{"main", "mian"},
		{"config", "conifg"},
		{"moduleNameResolver", "mnr"},
		{"hello", "hallo"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range pairs {
			LevenshteinRatio(p.a, p.b)
		}
	}
}

func BenchmarkRecordSelection(b *testing.B) {
	cache := NewQueryCache()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.RecordSelection(fmt.Sprintf("query%d", i%100), fmt.Sprintf("/file%d.go", i%1000))
	}
}

func BenchmarkGetBoostScores(b *testing.B) {
	cache := NewQueryCache()
	for i := 0; i < 100; i++ {
		cache.RecordSelection(fmt.Sprintf("query%d", i), fmt.Sprintf("/file%d.go", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetBoostScores("query50")
	}
}

func generateTestFiles(n int) []string {
	files := make([]string, n)
	names := []string{"main", "config", "test", "utils", "helper", "server", "client", "api", "model", "view"}
	exts := []string{".go", ".yaml", ".json", ".md", ".txt"}

	for i := range n {
		name := names[i%len(names)]
		ext := exts[i%len(exts)]
		files[i] = fmt.Sprintf("/project/src/%s_%d%s", name, i, ext)
	}
	return files
}

func generateVietnameseTestFiles(n int) []string {
	files := make([]string, n)
	names := []string{
		"Báo_cáo_tháng",
		"Hợp_đồng_thuê",
		"Đơn_xin_nghỉ",
		"Kế_hoạch_năm",
		"Biên_bản_họp",
		"Quyết_định",
		"Thông_báo",
		"Công_văn",
		"Tờ_trình",
		"Đề_xuất",
	}
	exts := []string{".pdf", ".docx", ".xlsx", ".pptx", ".txt"}

	for i := 0; i < n; i++ {
		name := names[i%len(names)]
		ext := exts[i%len(exts)]
		files[i] = fmt.Sprintf("/documents/%s_%d%s", name, i, ext)
	}
	return files
}

func scanTestData(root string) []string {
	var files []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return files
	}

	for _, entry := range entries {
		path := root + "/" + entry.Name()
		if entry.IsDir() {
			files = append(files, scanTestData(path)...)
		} else {
			files = append(files, path)
		}
	}
	return files
}
