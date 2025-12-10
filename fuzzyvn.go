/*
----------------
Author: verse91
Date: 2025-12-10
Version: v0.1.6
License: 0BSD
----------------

Structures:
fuzzyvn.go
├── Package + Imports (line 43-56)
├── Types (line 58-95)
│   ├── CacheEntry
│   ├── QueryCache
│   ├── Searcher
│   └── MatchResult
├── Utility Functions (line 97-139)
│   ├── Normalize()
│   └── LevenshteinRatio()
├── QueryCache Internal Methods (line 141-216)
│   ├── querySimilarity()  (private)
│   ├── moveToFront()      (private)
│   └── evictIfNeeded()    (private)
├── QueryCache Public Methods (line 218-456)
│   ├── NewQueryCache()
│   ├── SetMaxQueries()
│   ├── SetBoostScore()
│   ├── RecordSelection()
│   ├── GetBoostScores()
│   ├── GetRecentQueries()
│   ├── GetCachedFiles()
│   ├── GetAllRecentFiles()
│   ├── Size()
│   └── Clear()
└── Searcher (line 458-609)

	├── NewSearcher()
	├── NewSearcherWithCache()
	├── Search()
	├── RecordSelection()
	├── GetCache()
	└── ClearCache()
*/
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

// =============================================================================
// Struct
// =============================================================================

/*
  - Boost score tính theo SelectCount: File chọn nhiều lần → điểm boost cao hơn:
    boost = boostScore * similarity * SelectCount / 100
  - Khi vượt giới hạn (maxPerQuery = 5): Entry có SelectCount thấp nhất bị xóa
  - Lần search sau: Files có SelectCount cao được ưu tiên lên đầu
*/
type CacheEntry struct {
	FilePath    string // Đường dẫn file đã chọn
	SelectCount int    // Số lần user chọn file này
}

type QueryCache struct {
	mu          sync.RWMutex            // Dùng RWMutex thay vì Mutex vì search chủ yếu là đọc (99%), tránh race codition
	entries     map[string][]CacheEntry // Key là từ khóa (đã chuẩn hóa), Value là danh sách các CacheEntry
	queryOrder  []string                // Danh sách lưu thứ tự các từ khóa đã tìm, cái nào lâu không dùng thì xóa trước
	maxQueries  int                     // Giới hạn tổng số từ khóa được lưu
	maxPerQuery int                     // Giới hạn số file được lưu cho mỗi từ khóa
	boostScore  int                     // Điểm cho các file hay search
}

type Searcher struct {
	Originals     []string    // Data gốc (có dấu, viết hoa thường lộn xộn bla bla). Dùng để trả về kết quả hiển thị
	Normalized    []string    // Data đã chuẩn hóa cho fuzzy search
	FilenamesOnly []string    // Chỉ chứa tên file đã chuẩn hóa (bỏ đường dẫn). Dùng cho thuật toán Levenshtein (sửa lỗi chính tả)
	Cache         *QueryCache // Để lấy dữ liệu lịch sử
}

/*
- Struct tạm thời dùng để gom kết quả và điểm số lại để sắp xếp trước khi trả về cho người dùng
*/
type MatchResult struct {
	Str   string
	Score int
}

// =============================================================================
// Utility Functions
// =============================================================================

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func countWordMatches(queryWords []string, target string) int {
	targetWords := strings.Fields(target)
	count := 0
	for _, qWord := range queryWords {
		if len(qWord) < 2 {
			continue
		}
		for _, tWord := range targetWords {
			if len(tWord) < 2 {
				continue
			}
			// Exact match
			if qWord == tWord {
				count++
				break
			}
			// Fuzzy match: cho phép 1 lỗi nếu từ >= 3 ký tự
			if len(qWord) >= 3 && len(tWord) >= 3 {
				dist := LevenshteinRatio(qWord, tWord)
				if dist <= 1 {
					count++
					break
				}
			}
		}
	}
	return count
}

func Normalize(s string) string {
	// Tách dấu -> xóa dấu -> ghép lại
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, _ := transform.String(t, s)
	// Chỉ xử lý đ->d, KHÔNG đổi y->i vì sẽ làm sai lệch kết quả
	output = strings.ReplaceAll(output, "đ", "d")
	output = strings.ReplaceAll(output, "Đ", "D")
	return output
}

/*
- Levenshtein Distance: https://viblo.asia/p/khoang-cach-levenshtein-va-fuzzy-query-trong-elasticsearch-jvElaOXAKkw
- Bạn hiểu nôm na là để tính độ sai lệch khi gõ sai, tìm kết quả gần khớp với ý muốn của bạn nhất
- Mính sẽ chỉ triển khai cái nào cần cho tiếng Việt thôi, Trung, Hàn, Nhật,... bỏ qua
- Mục tiêu là biến chuỗi s1 thành s2
- Tại mỗi bước so sánh ký tự, ta có 3 quyền lựa chọn, ta sẽ chọn cái nào tốn ít chi phí nhất (minVal):
+ Xóa bỏ ký tự ở s1 (Chi phí +1)
+ Thêm ký tự vào s1 để giống s2 (Chi phí +1)
+ Thay thế:
> Nếu 2 ký tự giống nhau: Không mất phí (+0)
> Nếu khác nhau: Thay ký tự này bằng ký tự kia (+1)
NOTE: 1 điều lưu ý là ta không cần quan tâm chữ hoa, chữ thường vì đã chuẩn hóa rồi
*/
func LevenshteinRatio(s1, s2 string) int {
	/*
		4 dòng dưới
		Đây là trường hợp biến chuỗi s1 thành "chuỗi rỗng"
		Ví dụ s1 = "ABC", s2 = ""
		Biến "" thành "" mất 0 bước (column[0] = 0)
		Biến "A" thành "" mất 1 bước xóa (column[1] = 1)
		Biến "AB" thành "" mất 2 bước xóa (column[2] = 2)
		Lúc này mảng column trông như thế này: [0, 1, 2, 3, ... len(s1)]
		Tương ứng tăng dần từ 0 đến len(s1) là chi phí biến thành chuỗi rỗng
	*/
	s1Len := len(s1)
	s2Len := len(s2)
	column := make([]int, len(s1)+1)
	for y := 1; y <= s1Len; y++ {
		column[y] = y
	}

	/*
			Ở đây mình sẽ giải thích sơ sơ
			Thay vì dùng ma trận, mình dùng column như 1 stack từ trên xuống vậy, và ta sẽ ghi đè lên cái nào đã dùng
			Chủ yếu để tiết kiệm 1 chút bộ nhớ thôi
			Giờ nhìn ma trận trước
		        /*
		          "" |  A |  B |  C
		        ┌────┬────┬────┬────┐
		      ""│  0 │  1 │  2 │  3 │   ← khởi tạo, từ rỗng thành rỗng cần 0 bước, thành A cần qua chữ A, thành B cần qua A,B, thành C cần qua A,B,C
		        ├────┼────┼────┼────┤
		      A │  1 │  0 │  1 │  2 │   ← A=A (0), còn lại +1 theo cách biến đổi như cách khởi tạo
		        ├────┼────┼────┼────┤
		      X │  2 │  1 │  ? │  2 │   ← chuỗi AX, đổi thành rỗng cần 2 bước,... nhưng X≠B đọc tiếp xuống dưới
		        ├────┼────┼────┼────┤
		      C │  3 │  2 │  2 │  1 │   ← Tương tự, tại ô của B, Biến AX thành AB (tốn 1 bước sửa X->B), sau đó dư chữ C nên phải Xóa C (1 bước nữa)
		        └────┴────┴────┴────┘
			Kết quả tại ô "?" = 1
			Vì ô ? = min(trên, trái, chéo trái) + 1 (+1 khi ta thấy được ký tự khác nhau)
			Còn bạn nhìn vào ô (4,4) (C,C) ta thấy nó bằng 1 vì min(trên, trái, chéo trái) không + 1 vì C-C giống nhau
			Bây giờ, hãy xem chuyện gì xảy ra khi ta ép cái bảng trên vào 1 mảng duy nhất (column)

	*/
	for i := 1; i <= s2Len; i++ {
		column[0] = i    // Ví dụ: "" -> "A" (1 thêm), "" -> "AX" (2 thêm)
		lastKey := i - 1 // Lưu giá trị cũ của ô chéo trên trái ta đã đề cập
		for j := 1; j <= s1Len; j++ {
			/*
					IMPORTANT: Lưu lại giá trị cũ của column[j] trước khi bị ghi đè
					column[j] lúc này đang chứa giá trị của hàng bên trên (i-1)
					Sau khi vòng lặp này kết thúc, giá trị này sẽ trở thành
				    ô chéo trên trái cho vòng lặp tiếp theo (j+1)
			*/
			oldKey := column[j]
			/*
							Tính toán chi phí biến đổi:

									(lastKey)    (column[j] cũ)
				   					  CHÉO      |     TRÊN
				    				   ↘        |      ↓
				           					┌───────┐
				  				   TRÁI ──→ │  ???  │  (Đang tính)
				              (column[j-1]) └───────┘

							NOTE: lastKey = column[j-1]
			*/
			var incr int
			if s1[j-1] != s2[i-1] {
				incr = 1 // Khác nhau: +1 bước thay thế, còn không thì thôi không cần cộng
			}

			// Và đây chính xác là cái min chúng ta đã làm ở trên: min(trên, trái, chéo trái)
			// Xóa. Ví dụ: Name -> Nam
			minVal := column[j] + 1
			// Thêm. Ví dụ: Nam -> Name
			if column[j-1]+1 < minVal {
				minVal = column[j-1] + 1
			}
			// Sửa. Ví dụ: Năm -> Nấm
			if lastKey+incr < minVal {
				minVal = lastKey + incr
			}
			column[j] = minVal
			// Giá trị Trên của ô hiện tại (oldKey) sẽ trở thành
			// giá trị Chéo của ô bên phải
			lastKey = oldKey
		}
	}
	// Trả về chi phí cuối dựa trên độ dài s1 (phần tử cuối). Đọc tới đây mà không hiểu thì hãy xem lại ma trận
	return column[s1Len]
}

// =============================================================================
// QueryCache Internal Methods (Helpers)
// =============================================================================

/*
- querySimilarity: chấm điểm độ tương đồng giữa hai câu truy vấn (q1 và q2) trên thang điểm từ 0 đến 100
- Mục đích là để xem q1 (câu người dùng mới nhập) có đủ giống với q2 (câu đã lưu trong cache) hay không để tận dụng kết quả cũ
- Hàm này hoạt động theo cơ chế "Sàng lọc theo tầng" (Layered Filtering), ưu tiên độ chính xác từ cao xuống thấp
*/
func (c *QueryCache) querySimilarity(q1, q2 string) int {
	if q1 == q2 {
		return 100
	}
	// Ví dụ: "samsung" (q1), cache có "samsung s23" (q2)
	if strings.HasPrefix(q2, q1) {
		return 70 + (30 * len(q1) / len(q2))
	}

	if strings.HasPrefix(q1, q2) && len(q2) >= 2 {
		return 50 + (30 * len(q2) / len(q1))
	}
	// Ví dụ: q1="ip 15", q2="mua ip 15 giá rẻ"
	if len(q1) >= 2 && strings.Contains(q2, q1) {
		return 80
	}
	// Ví dụ: q1="mua ip 15 giá rẻ", q2="ip 15"
	if len(q2) >= 2 && strings.Contains(q1, q2) {
		return 60
	}
	// Bạn hoàn toàn có thể thay đổi cơ chế chấm điểm bên trên

	// Nếu không khớp chuỗi liền mạch, hàm sẽ cắt chuỗi thành từng từ (bằng strings.Fields) để so sánh
	// Logic này xử lý việc đảo từ
	words1 := strings.Fields(q1)
	words2 := strings.Fields(q2)
	// Tách từ ra, tìm xem có bao nhiêu từ giống nhau
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
		/*
			Ví dụ:
			q1: "sơn tùng mtp"
			q2: "mtp sơn tùng"
			Hai chuỗi này Contains sẽ sai, nhưng tách từ thì khớp 3 từ.
			Điểm: 50 + (15 điểm cho mỗi từ trùng). Nếu trùng 3 từ = 95 điểm
			Nếu điểm cao như này thì hoàn toàn khẳng định được đây là từ khóa cần tìm
		*/
		if commonWords > 0 {
			return 50 + (commonWords * 15)
		}
	}

	/*
		Sai chính tả
	*/
	if len(q1) >= 3 && len(q2) >= 3 {
		// Tính khoảng cách Levenshtein
		dist := LevenshteinRatio(q1, q2)
		maxLen := len(q1)
		if len(q2) > maxLen {
			maxLen = len(q2)
		}
		// Tính ngưỡng sai số cho phép (threshold): Khoảng 30% độ dài chuỗi dài nhất
		// Ví dụ chuỗi dài 10 ký tự thì cho phép sai tối đa 3 lỗi
		threshold := maxLen * 30 / 100
		if threshold < 2 {
			threshold = 2
		}
		/*
			Nếu số lỗi nằm trong ngưỡng cho phép: Trả về 60 trừ đi điểm phạt (mỗi lỗi trừ 10 điểm)
			Ví dụ:
			q1: "iphone"
			q2: "ipbone" (Sai 1 ký tự h -> b, dist = 1)
			Điểm: 60 - (1 * 10) = 50 điểm
		*/
		if dist <= threshold {
			return 60 - (dist * 10)
		}
	}

	return 0
}

//
/*
- moveToFront: Đẩy query lên đầu danh sách queryOrder
- Để query được tìm kiếm nhiều nhất sẽ được ưu tiên hơn
- Nhưng mà trong code ta đẩy xuống cuối mảng
*/
func (c *QueryCache) moveToFront(query string) {
	for i, q := range c.queryOrder {
		if q == query {
			c.queryOrder = append(c.queryOrder[:i], c.queryOrder[i+1:]...)
			break
		}
	}
	c.queryOrder = append(c.queryOrder, query)
}

/*
- evictIfNeeded: Xóa query cũ nhất nếu vượt giới hạn maxQueries
*/
func (c *QueryCache) evictIfNeeded() {
	for len(c.queryOrder) > c.maxQueries {
		oldestQuery := c.queryOrder[0]
		c.queryOrder = c.queryOrder[1:]
		delete(c.entries, oldestQuery)
	}
}

// =============================================================================
// QueryCache Public Methods
// =============================================================================

// NewQueryCache: Define một object QueryCache mặc định, có thể tùy chỉnh dựa tùy vào project của bạn
func NewQueryCache() *QueryCache {
	return &QueryCache{
		entries:     make(map[string][]CacheEntry),
		queryOrder:  make([]string, 0),
		maxQueries:  100,
		maxPerQuery: 5,
		boostScore:  5000,
	}
}

// SetMaxQueries: Đặt giới hạn tổng số từ khóa được lưu
// Sau khi giảm maxQueries, nó gọi evictIfNeeded để xóa bớt dữ liệu thừa ngay lập tức
func (c *QueryCache) SetMaxQueries(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxQueries = n
	c.evictIfNeeded()
}

// SetBoostScore: Đặt điểm boost cơ bản cho kết quả từ cache
func (c *QueryCache) SetBoostScore(score int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.boostScore = score
}

/*
- User search "main" -> Thấy danh sách files -> Click chọn "/project/main.go"
-> RecordSelection("main", "/project/main.go")
- NOTE: Hãy nhớ rằng bạn có thể định nghĩa thế nào là "chọn", hiện tại mình chỉ định nghĩa theo phần demo ở thư mục demo
*/
func (c *QueryCache) RecordSelection(query, filePath string) {
	if query == "" || filePath == "" {
		return
	}
	// Dùng lock vì hàm này cần ghi
	c.mu.Lock()
	defer c.mu.Unlock()
	// Chuẩn hóa query, xóa hết dấu, xóa hết các ký tự viết hoa
	// Ví dụ: Cộng đồng Golang Việt Nam -> cong dong golang viet nam
	queryNorm := strings.ToLower(Normalize(query))

	// Phải ưu tiên kiểm tra trong cache trước rồi mới tới các bước tiếp theo
	entries, exists := c.entries[queryNorm]
	// Case 1: Nếu có -> tăng count lên 1 và đẩy nó lên
	if exists {
		for i, entry := range entries {
			if entry.FilePath == filePath {
				c.entries[queryNorm][i].SelectCount++
				c.moveToFront(queryNorm)
				return
			}
		}
	}

	// Case 2: Nếu không có -> tạo mới một CacheEntry với count = 1
	newEntry := CacheEntry{FilePath: filePath, SelectCount: 1}
	if !exists {
		c.entries[queryNorm] = []CacheEntry{newEntry}
		c.queryOrder = append(c.queryOrder, queryNorm)
	} else { // Case 3: Đã có query nhưng mà file đó ta chưa thêm vào
		// Nếu đã đạt giới hạn số file được lưu cho mỗi từ khóa, xóa file có count thấp nhất
		if len(entries) >= c.maxPerQuery {
			minIdx := 0
			minCount := entries[0].SelectCount
			for i, e := range entries {
				if e.SelectCount < minCount {
					minCount = e.SelectCount
					minIdx = i
				}
			}
			// Xóa file có count thấp nhất
			c.entries[queryNorm] = append(entries[:minIdx], entries[minIdx+1:]...)
		}
		// Thêm file mới vào
		c.entries[queryNorm] = append(c.entries[queryNorm], newEntry)
	}

	c.moveToFront(queryNorm) // Đẩy query lên đầu danh sách
	c.evictIfNeeded()
}

/*
- GetBoostScores: Lấy điểm boost cho từng file dựa trên query người dùng
- Ví dụ như ta có query "màn hình"
- List ra sản phẩm có Màn hình Dell Ultrasharp hoặc chỉ đơn giản Dell Ultrasharp thôi, ...
- Nó sẽ học hành vi người dùng nhấn vào ví dụ Dell Ultrasharp mặc dù chả có cái chữ "màn hình" nào ở đây cả
- Nhưng nó vẫn sẽ lưu lại, càng nhiều lần càng cộng điểm
*/
func (c *QueryCache) GetBoostScores(query string) map[string]int {
	// Đoạn này dùng RLock vì chỉ đọc là chính, cho phép nhiều luồng (goroutine) cùng đọc một lúc
	// Điều này giúp hiệu năng cao hơn nhiều so với Lock thường (chỉ cho 1 người vào, dù chỉ để đọc)
	// Nói chung bạn hiểu nôm na là để handle nhiều query cùng một lúc
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]int)
	if query == "" {
		return result
	}

	queryNorm := strings.ToLower(Normalize(query))
	/*
		Kết quả nào càng giống ý định tìm kiếm VÀ càng được chọn nhiều trước đây, thì điểm cộng càng cao
		Dựa vào config boost cơ bản của bạn
		Độ giống nhau dựa vào querySimilarity
		Độ phổ biến dựa vào entry.SelectCount, kiểu như 1 người bấm vào chọn nhiều lần hoặc nhiều người bấm vào chọn
	*/
	for cachedQuery, entries := range c.entries {
		similarity := c.querySimilarity(queryNorm, cachedQuery)
		if similarity > 0 {
			for _, entry := range entries {
				boost := (c.boostScore * similarity * entry.SelectCount) / 100
				/*
									Một file (entry.FilePath) có thể xuất hiện trong nhiều cached query khác nhau
					    			Ví dụ: File "iPhone 15.html" xuất hiện khi tìm "iphone" và cả khi tìm "apple" (đại loại vậy)
									Đoạn code này đảm bảo: Nếu một file được tìm thấy nhiều lần,
									ta chỉ giữ lại điểm Boost cao nhất mà nó đạt được
				*/
				if currentBoost, exists := result[entry.FilePath]; !exists || boost > currentBoost {
					result[entry.FilePath] = boost
				}
			}
		}
	}

	return result
}

/*
- GetRecentQueries: Lấy danh sách query gần đây nhất
- Ví dụ như ta có query "màn hình"
- List ra các query gần đây nhất, ví dụ: "màn hình", "màn hình dell", "màn hình dell ultasharp", ...
- Hãy xem demo để biết
*/
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

/*
- GetCachedFiles: Lấy danh sách file đã lưu trong cache
- Ví dụ như ta có query "màn hình"
- List ra các file đã lưu trong cache như /data/products/dell/dell-ultrasharp.html
- Hãy xem demo để biết
*/
func (c *QueryCache) GetCachedFiles(query string, limit int) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if query == "" || limit <= 0 {
		return []string{}
	}

	queryNorm := Normalize(query)

	type fileScore struct {
		path  string
		score int
	}
	var matches []fileScore
	seen := make(map[string]bool)

	// Ưu tiên cao nhất cho những query đã từng được gõ y hệt
	if entries, ok := c.entries[queryNorm]; ok {
		for _, entry := range entries {
			// Điểm cơ bản cực cao (100) * số lần click
			score := 100 * entry.SelectCount
			matches = append(matches, fileScore{path: entry.FilePath, score: score})
			seen[entry.FilePath] = true
		}
	}

	// Tìm các query liên quan khác. Ví dụ: gõ "màn hình", tìm thấy cả trong lịch sử "màn hình dell"
	for cachedQuery, entries := range c.entries {
		// Bỏ qua nếu là chính nó (đã xử lý ở trên)
		if cachedQuery == queryNorm {
			continue
		}

		// Nếu độ dài chuỗi lệch nhau quá 5 ký tự, khả năng cao là không liên quan -> Bỏ qua để đỡ tốn tài nguyên
		if abs(len(cachedQuery)-len(queryNorm)) > 5 {
			continue
		}

		similarity := c.querySimilarity(queryNorm, cachedQuery)
		if similarity > 0 {
			for _, entry := range entries {
				// Nếu đã có trong phần tìm khớp rồi thì không add lại
				if seen[entry.FilePath] {
					continue
				}
				// Điểm = Độ giống * Độ phổ biến
				score := similarity * entry.SelectCount
				matches = append(matches, fileScore{path: entry.FilePath, score: score})
			}
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	result := make([]string, 0, limit)

	// Reset map seen để dùng cho việc filter kết quả trả về
	seenResult := make(map[string]bool)

	for _, m := range matches {
		if !seenResult[m.path] {
			seenResult[m.path] = true
			result = append(result, m.path)
			if len(result) >= limit {
				break
			}
		}
	}

	return result
}

/*
- GetAllRecentFiles: Lấy lịch sử danh sách file đã lưu trong cache
- List ra các file đã lưu trong cache như /data/products/dell/dell-ultrasharp.html,...
- Màn hình chính, input rỗng -> hiển thị "Recent Files"
- Hãy xem demo để biết
*/
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

/*
- Size: Lấy số lượng query đã lưu trong cache
*/
func (c *QueryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

/*
- Clear: Xóa tất cả query đã lưu trong cache
*/
func (c *QueryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string][]CacheEntry)
	c.queryOrder = make([]string, 0)
}

// =============================================================================
// Searcher
// =============================================================================

/*
- NewSearcher: Tạo Searcher mới
- items: Danh sách đường dẫn file cần index
*/
func NewSearcher(items []string) *Searcher {
	normPaths := make([]string, len(items))
	normNames := make([]string, len(items))

	for i, item := range items {
		filename := filepath.Base(item)
		// Ưu tiên tên file, theo path thì điểm thấp hơn
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

/*
- NewSearcherWithCache: Tạo Searcher mới với cache có sẵn
- items: Danh sáng đường dẫn file cần index
- cache: QueryCache có sẵn để tái sử dụng
*/
func NewSearcherWithCache(items []string, cache *QueryCache) *Searcher {
	s := NewSearcher(items)
	if cache != nil {
		s.Cache = cache
	}
	return s
}

/*
- Hàm quan trọng nhất, kết hợp Fuzzy Search + Levenshtein + Cache Boost
- Có lẽ mình quên nói ở trên là ta phải dùng Rune
- Ví dụ như:
s := "Việt Nam"
fmt.Println(len(s))  // 12 bytes -> SAI (8 mới đúng)
-> Đúng ra ta phải dùng Rune
s := "Việt Nam"
runes := []rune(s)
fmt.Println(len(runes))  // 8 (đúng 8 ký tự)
- Ta cần đếm số ký tự, chứ không tính theo byte được
*/
func (s *Searcher) Search(query string) []string {
	queryNorm := strings.ToLower(Normalize(query))
	queryRunes := []rune(queryNorm)
	queryLen := len(queryRunes) // đếm số ký tự, không phải byte
	queryWords := strings.Fields(queryNorm)

	uniqueResults := make(map[int]int)
	filePathToIdx := make(map[string]int)

	for i, fp := range s.Originals {
		filePathToIdx[fp] = i
	}
	// Ví dụ: User từng search "main" và chọn main.go nhiều lần:
	// cacheBoosts = {"/a/main.go": 5000}
	var cacheBoosts map[string]int
	if s.Cache != nil {
		cacheBoosts = s.Cache.GetBoostScores(query)
	}

	// Search bằng thư viện fuzzy ở trên
	matches := fuzzy.Find(queryNorm, s.Normalized)
	for _, m := range matches {
		// Word bonus tính trên tên file (không phải full path)
		// Vì user search "ve mua thu" thì nên ưu tiên file tên là "ve mua thu"
		// không phải file nằm trong folder có chứa các ký tự đó
		wordMatches := countWordMatches(queryWords, s.FilenamesOnly[m.Index])
		wordBonus := wordMatches * 3000
		uniqueResults[m.Index] = m.Score + wordBonus
	}

	// Ta tính điểm sai chính tả dựa trên Levenshtein
	// Tức là nếu user gõ "maain" hay "mian" thì ta vẫn tính điểm cho "main"
	// Threshold = (queryLen / 3) + 1: cho phép khoảng 1 lỗi mỗi 3 ký tự + 1 lỗi bonus
	// Minimum threshold = 3: query ngắn (2-5 ký tự) vẫn cần đủ độ linh hoạt để match
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

			// So sánh với phần đầu của filename
			targetStr1 := string(runesName[:queryLen])
			d1 := LevenshteinRatio(queryNorm, targetStr1)
			dist = d1

			// So sánh thêm 1 ký tự (phòng trường hợp typo thêm ký tự)
			if len(runesName) > queryLen {
				targetStr2 := string(runesName[:queryLen+1])
				d2 := LevenshteinRatio(queryNorm, targetStr2)
				if d2 < dist {
					dist = d2
				}
			}

			/*
				Ở phần trên ví dụ như "mian", target 1 là "main" target 2 là "maina"
				Ta tính điểm ở target 1, dist = d1 = 2, nhưng ở target 2, dist = d2 = 3
				if d2 < dist {
						dist = d2
					}
				Tức là nếu nhỏ hơn cái d1 thì lấy, còn không thì giữ nguyên
				Kiểu như min(d1, d2)
			*/

			// Nếu điểm sai chính tả nhỏ hơn ngưỡng cho phép thì tính điểm
			// Robust solution khi sai chính tả đi quá xa (hoặc nếu không thì mong bạn có thể mở PR hỗ trợ mình)
			if dist < baseThreshold {
				score := 10000 - (dist * 100)

				lenDiff := len(runesName) - queryLen
				if lenDiff > 0 {
					score -= (lenDiff / 2)
				}

				// Thêm word bonus cho Levenshtein matches
				// Dùng tên file để tính word matches (không phải full path)
				wordMatches := countWordMatches(queryWords, s.FilenamesOnly[i])
				score += wordMatches * 3000

				if oldScore, exists := uniqueResults[i]; !exists || score > oldScore {
					uniqueResults[i] = score
				}
			}
		}
	}
	/*
		Đảm bảo file đã cache luôn xuất hiện trong kết quả, kể cả khi fuzzy/Levenshtein không match
		Thì ví dụ như:
		User search "tiền lương", xong họ chả chọn cái gì liên quan tới tiền lương
		nhưng chọn "bao_cao_tai_chinh_2024.xlsx"
		Hệ thống lưu lại: Query: "tiền lương" -> File: "bao_cao..."
		Xong giờ search lại "tien luong" một lần nữa
		Lúc này cả fuzzy và levenshtein đều không match
		Đoạn code này sẽ giải quyết vấn đề trên
		Nó vẫn in ra "bao_cao_tai_chinh_2024.xlsx", vì trước đây từng có hành vi này
		Và có thể nó sẽ là 1 trong những file user cần
		Đây chỉ là một cơ chế phòng bị cho trường hợp user quên tên file
		vì nó cũng không có độ chính xác quá cao
	*/
	for cachedPath, boost := range cacheBoosts {
		if idx, exists := filePathToIdx[cachedPath]; exists {
			if _, alreadyInResults := uniqueResults[idx]; !alreadyInResults {
				uniqueResults[idx] = boost // Không làm hỏng kết quả khác vì chỉ thêm khi chưa có
			}
		}
	}

	/*
		File: "/a/main.go"
		Fuzzy score: 85
		Cache boost: 5000
		Final score: 85 + 5000 = 5085 -> Lên top
	*/
	var rankedResults []MatchResult
	for idx, score := range uniqueResults {
		filePath := s.Originals[idx]
		finalScore := score

		if boost, exists := cacheBoosts[filePath]; exists {
			if score != boost { // Tránh duplicate
				finalScore += boost
			}
		}

		rankedResults = append(rankedResults, MatchResult{
			Str:   filePath,
			Score: finalScore,
		})
	}
	// Logic:
	// Điểm cao lên trước
	// Cùng điểm, ưu tiên file path ngắn
	sort.SliceStable(rankedResults, func(i, j int) bool {
		if rankedResults[i].Score == rankedResults[j].Score {
			return len(rankedResults[i].Str) < len(rankedResults[j].Str)
		}
		return rankedResults[i].Score > rankedResults[j].Score
	})
	// Trả về top 20, nếu kết quả ít hơn 20 thì show bấy nhiêu thôi
	// Hãy xem demo
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

/*
- RecordSelection: Chỉ để gọi nhanh hơn, ngắn hơn
*/
func (s *Searcher) RecordSelection(query, filePath string) {
	if s.Cache != nil {
		s.Cache.RecordSelection(query, filePath)
	}
}

/*
- GetCache: Lấy object cache
- Ví dụ:
cache := searcher.GetCache()
cache.GetRecentQueries(5)
cache.GetAllRecentFiles(10)
cache.SetMaxQueries(50)
*/
func (s *Searcher) GetCache() *QueryCache {
	return s.Cache
}

func (s *Searcher) ClearCache() {
	if s.Cache != nil {
		s.Cache.Clear()
	}
}
