## Các test đã viết:

| Test | Mô tả |
|------|-------|
| `TestNormalize` | Chuẩn hóa tiếng Việt (đ→d, bỏ dấu) |
| `TestNormalize_IY_Equivalence` | i/y phải giống nhau (kỷ=kỉ) |
| `TestLevenshteinRatio` | Khoảng cách Levenshtein |
| `TestNewSearcher` | Khởi tạo Searcher |
| `TestSearcher_Search_Basic` | Tìm kiếm cơ bản |
| `TestSearcher_Search_Vietnamese` | Tìm tiếng Việt không dấu |
| `TestSearcher_Search_IY_Equivalence` | Search "ky"="ki" |
| `TestSearcher_Search_Typo` | Chịu lỗi gõ (mian→main) |
| `TestQueryCache_RecordSelection` | Ghi nhận lựa chọn |
| `TestQueryCache_GetBoostScores_ExactMatch` | Boost khớp chính xác |
| `TestQueryCache_GetBoostScores_SimilarQuery` | Boost query tương tự |
| `TestQueryCache_GetRecentQueries` | Lấy queries gần đây |
| `TestQueryCache_GetCachedFiles` | Lấy files đã cache |
| `TestQueryCache_GetAllRecentFiles` | Lấy tất cả files gần đây |
| `TestQueryCache_LRU_Eviction` | Xóa cache cũ khi đầy |
| `TestQueryCache_Clear` | Xóa toàn bộ cache |
| `TestSearcher_RecordSelection_BoostsResults` | Cache boost kết quả |
| `TestNewSearcherWithCache` | Giữ cache khi rebuild |
| `TestSearcher_ClearCache` | Xóa cache của Searcher |

Chạy test: `go test -v`