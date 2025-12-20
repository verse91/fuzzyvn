# Bây giờ mới import
import sys
import random
import os
try
    from datasets import load_dataset
    import pandas
except ImportError:
    print("❌ Lỗi: Chưa cài thư viện 'datasets' hoặc 'pandas'.")
    print("👉 Hãy chạy bằng lệnh: make gen")
    sys.exit(1)
# Cấu hình số lượng
NUM_CODE_PATHS = 70000
NUM_VN_PATHS = 30000
OUTPUT_FILE = "demo/gen_data/test_paths_100k.txt" # Sửa lại đường dẫn cho đúng vị trí chạy

paths = []

print("🚀 Bắt đầu tải dữ liệu mẫu từ Hugging Face...")

print("1. Đang tải đường dẫn code từ 'bigcode/the-stack-smol'...")
try:
    ds_code = load_dataset("bigcode/the-stack-smol", data_dir="data/python", split="train", streaming=True)
    count = 0
    for sample in ds_code:
        # Lấy đường dẫn thực tế
        repo = sample.get("repository_name", "unknown_repo")
        path = sample.get("file_path", f"file_{count}.py")
        paths.append(f"{repo}/{path}")
        count += 1
        if count >= NUM_CODE_PATHS:
            break
except Exception as e:
    print(f"⚠️ Lỗi khi tải dataset code: {e}")
    print("-> Sẽ dùng dữ liệu giả lập cho phần code.")
    # Fallback nếu lỗi mạng
    for i in range(NUM_CODE_PATHS):
        paths.append(f"github.com/user/repo/src/main_{i}.go")

print("2. Đang tạo đường dẫn tiếng Việt từ 'ura-hcmut/vietnamese-news'...")
try:
    ds_vn = load_dataset("ura-hcmut/vietnamese-news", split="train", streaming=True)
    extensions = [".pdf", ".docx", ".xlsx", ".pptx", ".txt"]
    folders = ["Tài liệu", "Báo cáo", "Hợp đồng", "Nhân sự", "Kế toán", "Dự án"]

    count = 0
    for sample in ds_vn:
        title = sample["title"]
        # Làm sạch tiêu đề
        safe_name = title.replace(" ", "_").replace("/", "-").replace('"', '').replace("'", "")[:60]

        folder = random.choice(folders)
        ext = random.choice(extensions)

        full_path = f"{folder}/{safe_name}{ext}"
        paths.append(full_path)

        count += 1
        if count >= NUM_VN_PATHS:
            break
except Exception as e:
    print(f"⚠️ Lỗi khi tải dataset VN: {e}")
    # Fallback
    for i in range(NUM_VN_PATHS):
        paths.append(f"Tài liệu/Báo_cáo_tài_chính_{i}.docx")

print(f"3. Đang trộn và ghi {len(paths)} dòng ra file {OUTPUT_FILE}...")
random.shuffle(paths)

# Đảm bảo thư mục tồn tại
os.makedirs(os.path.dirname(OUTPUT_FILE), exist_ok=True)

with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
    for p in paths:
        f.write(p + "\n")

print("✅ Hoàn tất! Giờ bạn có thể chạy 'make gen' hoặc 'go run demo/gen_data/gen_data.go'")
