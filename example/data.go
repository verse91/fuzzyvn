//go:build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 1. KHO TỪ VỰNG (Dictionary)
var (
	// Mảng âm nhạc
	artists   = []string{"Sơn Tùng M-TP", "Đen Vâu", "Mỹ Tâm", "Hà Anh Tuấn", "Phan Mạnh Quỳnh", "Hoàng Thùy Linh", "Binz", "JustaTee", "Jack", "Erik", "Đức Phúc", "Vũ", "Chillies", "Ngọt"}
	songVbs   = []string{"Yêu", "Thương", "Nhớ", "Quên", "Chờ", "Đợi", "Mơ", "Khóc", "Cười", "Đi", "Về", "Lạc", "Trôi", "Bay", "Chạy", "Ngủ", "Thức"}
	songAdjs  = []string{"Vội Vàng", "Mong Manh", "Xa Vời", "Ngọt Ngào", "Đắng Cay", "Bình Yên", "Lặng Lẽ", "Ồn Ào", "Cuối Cùng", "Đầu Tiên", "Mãi Mãi", "Vô Tận"}
	songNouns = []string{"Mùa Thu", "Cơn Mưa", "Nỗi Buồn", "Hạnh Phúc", "Kỷ Niệm", "Giấc Mơ", "Con Đường", "Thành Phố", "Đại Dương", "Bầu Trời", "Ánh Sáng", "Bóng Tối"}

	// Mảng Công việc (IT & Office)
	depts     = []string{"Dev_Team", "HR_Nhan_Su", "Accounting_Ke_Toan", "Marketing", "Board_Of_Directors", "Sales_Team", "Design_Team"}
	projects  = []string{"Aegirus_WAF", "E_Swap_System", "Twitter_Extension", "Internal_Portal", "Customer_App", "Payment_Gateway", "Microservice_Core"}
	docTypes  = []string{"Báo Cáo", "Quy Trình", "Hợp Đồng", "Biên Bản", "Đề Xuất", "Kế Hoạch", "Hướng Dẫn", "Tài Liệu", "Sơ Đồ", "Kiến Trúc"}
	docStatus = []string{"Draft", "Final", "Signed", "Review", "v1.0", "v2.1", "Approved", "Rejected"}

	// Mảng Sách truyện
	bookCats  = []string{"Tiên Hiệp", "Kiếm Hiệp", "Ngôn Tình", "Trinh Thám", "Kinh Tế", "Kỹ Năng Sống", "Lịch Sử", "Văn Học"}
	bookAdjs  = []string{"Huyền Bí", "Vĩ Đại", "Bí Ẩn", "Tàn Khốc", "Rực Rỡ", "Đen Tối", "Thần Thánh", "Tuyệt Thế", "Vô Song"}
	bookNouns = []string{"Chiến Thần", "Bá Chủ", "Thiên Hạ", "Đế Vương", "Sát Thủ", "Bảo Vật", "Bí Kíp", "Trường Sinh", "Định Mệnh", "Thanh Xuân"}
)

var (
	totalFiles int
	mu         sync.Mutex
	rootDir    = "./massive_data"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Dang xoa du lieu cu...")
	os.RemoveAll(rootDir)
	fmt.Printf("Dang tao 10,000+ file tai: %s ...\n", rootDir)

	var wg sync.WaitGroup

	// Tạo 3 luồng lớn để sinh dữ liệu song song
	wg.Add(3)
	go generateMusic(&wg, 3000) // 3000 bài hát
	go generateDocs(&wg, 4000)  // 4000 tài liệu công việc
	go generateBooks(&wg, 3000) // 3000 quyển truyện

	wg.Wait()
	fmt.Printf("\n✅ XONG! Tong so file: %d\n", totalFiles)
	fmt.Println("Hay chay: go run main.go")
}

// --- GENERATOR: MUSIC ---
func generateMusic(wg *sync.WaitGroup, count int) {
	defer wg.Done()
	basePath := filepath.Join(rootDir, "Music_Lossless")

	for i := 0; i < count; i++ {
		artist := artists[rand.Intn(len(artists))]
		// Tên bài hát: Động từ + Danh từ + Tính từ (Vd: Yêu Mùa Thu Vội Vàng)
		songName := fmt.Sprintf("%s %s %s",
			songVbs[rand.Intn(len(songVbs))],
			songNouns[rand.Intn(len(songNouns))],
			songAdjs[rand.Intn(len(songAdjs))],
		)

		folder := filepath.Join(basePath, artist)
		os.MkdirAll(folder, 0o755)

		ext := ".mp3"
		if rand.Intn(2) == 0 {
			ext = ".flac"
		}

		filename := fmt.Sprintf("%s - %s%s", songName, artist, ext)
		createFile(filepath.Join(folder, filename))
	}
}

// --- GENERATOR: OFFICE DOCS ---
func generateDocs(wg *sync.WaitGroup, count int) {
	defer wg.Done()
	basePath := filepath.Join(rootDir, "Work_Documents")

	for i := 0; i < count; i++ {
		dept := depts[rand.Intn(len(depts))]
		proj := projects[rand.Intn(len(projects))]

		// Folder theo Phòng ban -> Dự án
		folder := filepath.Join(basePath, dept, proj)
		os.MkdirAll(folder, 0o755)

		// Tên file: Loại + Năm + Trạng thái (Vd: Báo Cáo_2025_Final.pdf)
		dType := docTypes[rand.Intn(len(docTypes))]
		status := docStatus[rand.Intn(len(docStatus))]
		year := 2020 + rand.Intn(6) // 2020-2025

		exts := []string{".docx", ".xlsx", ".pdf", ".pptx", ".txt"}
		ext := exts[rand.Intn(len(exts))]

		// Đôi khi thêm tên dự án vào tên file cho dễ search
		filename := fmt.Sprintf("%s_%s_%d_%s%s", dType, proj, year, status, ext)

		// Nếu là folder Dev_Team thì sinh code
		if dept == "Dev_Team" {
			ext = ".go"
			if rand.Intn(2) == 0 {
				ext = ".ts"
			}
			filename = fmt.Sprintf("%s_module_%d%s", proj, rand.Intn(1000), ext)
		}

		createFile(filepath.Join(folder, filename))
	}
}

// --- GENERATOR: BOOKS ---
func generateBooks(wg *sync.WaitGroup, count int) {
	defer wg.Done()
	basePath := filepath.Join(rootDir, "Ebook_Library")

	for i := 0; i < count; i++ {
		cat := bookCats[rand.Intn(len(bookCats))]

		// Tên truyện: Tính từ + Danh từ (Vd: Tuyệt Thế Chiến Thần)
		title := fmt.Sprintf("%s %s",
			bookAdjs[rand.Intn(len(bookAdjs))],
			bookNouns[rand.Intn(len(bookNouns))],
		)

		folder := filepath.Join(basePath, cat)
		os.MkdirAll(folder, 0o755)

		chap := rand.Intn(1000) + 1
		exts := []string{".pdf", ".epub", ".mobi"}
		ext := exts[rand.Intn(len(exts))]

		filename := fmt.Sprintf("%s - Chương %d%s", title, chap, ext)
		createFile(filepath.Join(folder, filename))
	}
}

func createFile(path string) {
	f, err := os.Create(path)
	if err == nil {
		f.Close()
		mu.Lock()
		totalFiles++
		mu.Unlock()
		if totalFiles%1000 == 0 {
			fmt.Printf("..da tao %d files..\n", totalFiles)
		}
	}
}
