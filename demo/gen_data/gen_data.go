package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	SourceListPath = "test_paths_100k.txt"
	TargetDir      = "../test_data"
	MaxWorkers     = 50
)

func main() {
	start := time.Now()

	file, err := os.Open(SourceListPath)
	if err != nil {
		log.Fatalf("Không thể mở file nguồn: %v", err)
	}
	defer file.Close()

	if err := os.MkdirAll(TargetDir, 0o755); err != nil {
		log.Fatalf("Không thể tạo thư mục gốc: %v", err)
	}

	sem := make(chan struct{}, MaxWorkers)
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	count := 0

	fmt.Printf("Đang tạo file vào thư mục: %s\n", TargetDir)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		cleanPath := strings.TrimPrefix(line, "test_data/")
		cleanPath = strings.TrimPrefix(cleanPath, "test_data\\")

		count++
		wg.Add(1)

		sem <- struct{}{}

		go func(pathStr string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := createDummyFile(pathStr); err != nil {
				log.Printf("Lỗi tạo file [%s]: %v\n", pathStr, err)
			}
		}(cleanPath)

		if count%5000 == 0 {
			fmt.Printf("...đã xử lý %d đường dẫn\n", count)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("\nHoàn tất! Đã tạo %d file trong %v\n", count, duration)
}

func createDummyFile(relativePath string) error {
	finalPath := filepath.Join(TargetDir, relativePath)

	dir := filepath.Dir(finalPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
