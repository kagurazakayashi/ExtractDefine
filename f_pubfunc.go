package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// `#.*$`
func removeComments(line string, commentChar string) string {
	if index := strings.Index(line, commentChar); index != -1 {
		return line[:index]
	}
	return line
}

// 列出指定資料夾中符合萬用字元模式的所有檔案的絕對路徑
func listFilesMatchingPattern(folderPath, pattern string) ([]string, error) {
	var matches []string
	// 獲取資料夾的絕對路徑
	absFolderPath, err := filepath.Abs(folderPath)
	if err != nil {
		return nil, err
	}
	// 構建搜尋模式
	searchPattern := filepath.Join(absFolderPath, pattern)
	// 使用 Glob 函式匹配檔案
	matches, err = filepath.Glob(searchPattern)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
