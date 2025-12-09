# Thuật toán tìm kiếm

FuzzyVN sử dụng thuật toán tìm kiếm đa giai đoạn kết hợp ba kỹ thuật.

## Giai đoạn 1: Fuzzy Matching

Sử dụng `github.com/sahilm/fuzzy` cho khớp substring.

**Đầu vào:** Query đã chuẩn hóa so với đường dẫn file đã chuẩn hóa

**Tính điểm:** Dựa trên vị trí ký tự khớp và khoảng cách

**Ví dụ:**
```
Query: "main"
File:  "project/src/main_server.go"
Khớp: Điểm dựa trên "main" xuất hiện trong tên file
```

## Giai đoạn 2: Khoảng cách Levenshtein

Xử lý lỗi gõ và tên file tương tự.

### Cấu hình ngưỡng

```go
threshold = (queryLength / 3) + 1
if threshold < 3 {
    threshold = 3
}
```

| Độ dài Query | Số lỗi cho phép |
|--------------|-----------------|
| 1-8 ký tự | 3 lỗi |
| 9-11 ký tự | 4 lỗi |
| 12-14 ký tự | 5 lỗi |
| v.v. | +1 mỗi 3 ký tự |

### Kỹ thuật kiểm tra kép

So sánh query với hai phần của tên file:

1. **Độ dài chính xác**: `filename[:queryLen]`
2. **Độ dài mở rộng**: `filename[:queryLen+1]`

Điều này xử lý cả:
- Lỗi gõ: `"mian"` → `"main"` (cùng độ dài)
- Thiếu ký tự: `"main"` → `"maain"` (khớp mở rộng)

### Tính điểm

```go
score = 10000 - (distance × 100) - (lengthDiff / 2)
```

- Điểm cơ bản: 10000
- Phạt mỗi lỗi: -100
- Phạt tên file dài hơn: `-lengthDiff/2`

## Giai đoạn 3: Tiêm Cache

File từ cache được thêm vào ngay cả khi không khớp fuzzy/Levenshtein.

```go
for cachedPath, boost := range cacheBoosts {
    if idx, exists := filePathToIdx[cachedPath]; exists {
        if _, alreadyInResults := uniqueResults[idx]; !alreadyInResults {
            uniqueResults[idx] = boost
        }
    }
}
```

## Tổng hợp điểm

Điểm cuối cùng kết hợp tất cả giai đoạn:

```go
finalScore = fuzzyScore + cacheBoost
```

Nếu file chỉ tìm thấy qua Levenshtein:
```go
finalScore = levenshteinScore + cacheBoost
```

Nếu file chỉ được tiêm từ cache:
```go
finalScore = cacheBoost
```

## Sắp xếp

Kết quả được sắp xếp theo:
1. Điểm (giảm dần)
2. Độ dài đường dẫn (tăng dần) khi hòa điểm

```go
sort.SliceStable(results, func(i, j int) bool {
    if results[i].Score == results[j].Score {
        return len(results[i].Str) < len(results[j].Str)
    }
    return results[i].Score > results[j].Score
})
```

## Chuẩn hóa tiếng Việt

Trước khi khớp, văn bản được chuẩn hóa:

1. **Phân tách NFD**: Tách ký tự cơ bản khỏi dấu
2. **Xóa dấu**: Loại bỏ combining marks
3. **Tổng hợp NFC**: Ghép lại
4. **Xử lý Đ/đ**: Thay bằng D/d

```go
"Đường Nguyễn Huệ" → "Duong Nguyen Hue"
```

## Hiệu năng

| Thao tác | Độ phức tạp |
|----------|-------------|
| Fuzzy match | O(n × m) với n=số file, m=độ dài query |
| Levenshtein | O(n × k²) với k=độ dài query |
| Cache lookup | O(q × f) với q=số query cache, f=file mỗi query |
| Sắp xếp | O(r log r) với r=số kết quả |

Với 10,000 file và query trung bình, thời gian phản hồi: 5-20ms.

