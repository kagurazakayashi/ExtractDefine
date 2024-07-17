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

func removeCommentsC(lines []string) []string {
	var result []string
	inBlockComment := false
	for _, line := range lines {
		newLine := ""
		i := 0
		for i < len(line) {
			if inBlockComment {
				if i+1 < len(line) && line[i] == '*' && line[i+1] == '/' {
					inBlockComment = false
					i += 2
				} else {
					i++
				}
			} else {
				if i+1 < len(line) && line[i] == '/' && line[i+1] == '*' {
					inBlockComment = true
					i += 2
				} else if i+1 < len(line) && line[i] == '/' && line[i+1] == '/' {
					break
				} else {
					newLine += string(line[i])
					i++
				}
			}
		}
		if !inBlockComment && strings.TrimSpace(newLine) != "" {
			result = append(result, newLine)
		}
	}
	return result
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

func normalizePath(path string) string {
	// var osPathSeparator string = string(os.PathSeparator)
	path = strings.ReplaceAll(path, "\\", "/")
	path = filepath.FromSlash(path)
	return path
}

func completeFilePath(folderPath string, filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	return filepath.Join(folderPath, filePath)
}
