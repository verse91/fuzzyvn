# Tài liệu FuzzyVN

FuzzyVN là thư viện tìm kiếm file mờ (fuzzy finder) được tối ưu cho tiếng Việt. Kết hợp nhiều thuật toán tìm kiếm với hệ thống cache thông minh để cho kết quả nhanh và chính xác.

## Mục lục

- [Cài đặt](./docs/installation.md)
- [Bắt đầu nhanh](./docs/quickstart.md)
- [API](./docs/api.md)
- [Hệ thống Cache](./docs/cache.md)
- [Thuật toán](./docs/algorithm.md)
- [Ví dụ](./docs/examples.md)

## Tính năng

- **Hỗ trợ tiếng Việt**: Xử lý dấu tiếng Việt (chuyển "Đường" thành "Duong")
- **Đa thuật toán**: Kết hợp fuzzy matching + Levenshtein distance
- **Cache thông minh**: Học từ lựa chọn của người dùng để đẩy kết quả liên quan lên đầu
- **Chịu lỗi gõ**: Xử lý lỗi gõ phím thường gặp
- **Thread-Safe**: An toàn khi truy cập đồng thời

## Kiến trúc

```
┌─────────────────────────────────────────────────────────┐
│                      Searcher                           │
├─────────────────────────────────────────────────────────┤
│  Originals[]     - Đường dẫn file gốc                   │
│  Normalized[]    - Đã chuẩn hóa cho fuzzy search        │
│  FilenamesOnly[] - Chỉ tên file cho Levenshtein         │
│  Cache           - Cache query để boost kết quả         │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    QueryCache                           │
├─────────────────────────────────────────────────────────┤
│  entries{}       - query → []CacheEntry                 │
│  queryOrder[]    - Thứ tự LRU                           │
│  maxQueries      - Tối đa queries cache (100)           │
│  maxPerQuery     - Tối đa files mỗi query (5)           │
│  boostScore      - Hệ số boost (5000)                   │
└─────────────────────────────────────────────────────────┘
```

## Luồng tìm kiếm

```
Query người dùng
    │
    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│ Fuzzy Match  │ +  │ Levenshtein  │ +  │ Cache Boost  │
│ (substring)  │    │ (sửa lỗi gõ) │    │ (lịch sử)    │
└──────────────┘    └──────────────┘    └──────────────┘
    │                      │                    │
    └──────────────────────┼────────────────────┘
                           ▼
                    Gộp kết quả
                           │
                           ▼
                    Sắp xếp theo điểm
                           │
                           ▼
                    Top 20 kết quả
```
