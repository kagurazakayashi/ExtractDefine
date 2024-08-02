package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// readFile 函數從指定的檔案路徑讀取檔案，並返回包含檔案每一行的字串切片。
// 如果在打開檔案或讀取過程中發生錯誤，則返回錯誤資訊。
// filePath：要讀取的檔案路徑。
// 返回：包含檔案每一行的字串切片和可能的錯誤。
func readFile(filePath string) ([]string, error) {
	// 嘗試打開指定路徑的檔案
	file, err := os.Open(filePath)
	if err != nil {
		// 如果檔案無法打開，返回錯誤
		return nil, err
	}
	// 延遲關閉檔案，確保在函數返回後檔案被正確關閉
	defer file.Close()

	var lines []string
	// 建立一個新的掃描器來逐行讀取檔案
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 將每行讀取到的內容追加到字串切片中
		lines = append(lines, scanner.Text())
	}

	// 返回包含檔案每行內容的字串切片及掃描器可能遇到的錯誤
	return lines, scanner.Err()
}

// removeComments 函式會移除指定字串中的註解部分。
// 給定一行字串和註解符號，此函式會返回去除註解後的字串。
//
// 參數：
// - line: 包含原始內容的一行字串。
// - commentChar: 表示註解開始的字元或字串（例如 "//" 或 "#"）。
//
// 返回值：
// 返回一個字串，其中包含去除註解後的內容。
func removeComments(line string, commentChar string) string {
	// `#.*$`
	// 使用 strings.Index 函式尋找註解符號在字串中的位置。
	if index := strings.Index(line, commentChar); index != -1 {
		// 如果找到註解符號，返回該符號之前的子字串，表示去除了註解。
		return line[:index]
	}
	// 如果沒有找到註解符號，則返回原始字串。
	return line
}

// removeCommentsC 函式移除給定的 C 語言風格註解（單行註解和區塊註解）。
// 參數：
//   - lines: 包含原始代碼行的字串切片。
//
// 返回值：
//   - 返回移除了註解後的代碼行切片。
func removeCommentsC(lines []string) []string {
	var result []string     // 儲存處理後的結果行
	inBlockComment := false // 標記是否在區塊註解中

	for _, line := range lines { // 遍歷每一行代碼
		newLine := "" // 儲存處理後的單行結果
		i := 0        // 當前行的字符索引

		for i < len(line) { // 遍歷行內的每個字符
			if inBlockComment { // 如果當前在區塊註解中
				if i+1 < len(line) && line[i] == '*' && line[i+1] == '/' {
					// 發現區塊註解結束符號 "*/"
					inBlockComment = false
					i += 2 // 跳過結束符號
				} else {
					i++ // 繼續跳過區塊註解內的字符
				}
			} else {
				if i+1 < len(line) && line[i] == '/' && line[i+1] == '*' {
					// 發現區塊註解開始符號 "/*"
					inBlockComment = true
					i += 2 // 跳過開始符號
				} else if i+1 < len(line) && line[i] == '/' && line[i+1] == '/' {
					// 發現單行註解 "//" 則忽略該行後續內容
					break
				} else {
					// 非註解內容添加到新行
					newLine += string(line[i])
					i++
				}
			}
		}

		// 如果新行不在區塊註解中且不為空行，則將其添加到結果中
		if !inBlockComment && strings.TrimSpace(newLine) != "" {
			result = append(result, newLine)
		}
	}

	return result // 返回處理後的代碼行切片
}

// listFilesMatchingPattern 列出指定資料夾中符合萬用字元模式的所有檔案的絕對路徑
//
// 這個函式接受一個資料夾路徑以及一組以逗號分隔的萬用字元模式，然後返回與這些模式匹配的檔案的絕對路徑清單。
// 如果在處理過程中出現錯誤，將返回錯誤訊息。
//
// 參數:
// - folderPath: 資料夾的路徑。
// - patterns: 萬用字元模式的字串，模式之間以逗號分隔。
//
// 返回值:
// - []string: 符合萬用字元模式的檔案的絕對路徑清單。
// - error: 如果發生錯誤，將返回錯誤訊息。
func listFilesMatchingPattern(folderPath, patterns string) ([]string, error) {
	var matches []string // 用於儲存符合模式的檔案清單

	// 取得資料夾的絕對路徑
	absFolderPath, err := filepath.Abs(folderPath)
	if err != nil {
		return nil, err // 如果無法取得絕對路徑，返回錯誤
	}

	// 將模式字串以逗號分割為清單
	patternList := strings.Split(patterns, ",")

	// 針對每個模式進行匹配
	for _, pattern := range patternList {
		// 去掉模式字串兩端的空格
		pattern = strings.TrimSpace(pattern)
		// 建立搜尋模式的完整路徑
		searchPattern := filepath.Join(absFolderPath, pattern)
		// 使用 Glob 函式來匹配符合的檔案
		patternMatches, err := filepath.Glob(searchPattern)
		if err != nil {
			return nil, err // 如果匹配過程中出現錯誤，返回錯誤訊息
		}
		// 將匹配到的檔案路徑加入結果清單中
		matches = append(matches, patternMatches...)
	}

	// 返回匹配到的檔案清單
	return matches, nil
}

// normalizePath 將路徑標準化，使其在不同作業系統間具可攜性。
// 此函數將 Windows 的反斜線 "\" 替換為 Unix 標準的斜線 "/"
// 並將斜線路徑轉換為當前作業系統的格式。
//
// 參數:
// - path: 要標準化的路徑字串。
//
// 回傳值:
// - string: 標準化後的路徑字串。
func normalizePath(path string) string {
	// 將所有反斜線 "\" 替換為正斜線 "/"，以符合 Unix 系統的標準。
	path = strings.ReplaceAll(path, "\\", "/")

	// 將正斜線路徑轉換為當前作業系統的路徑格式。
	path = filepath.FromSlash(path)

	// 回傳標準化後的路徑字串。
	return path
}

// completeFilePath 函式用於將相對路徑轉換為完整的文件路徑。
// 如果輸入的 filePath 已經是絕對路徑，則直接返回該絕對路徑。
// 否則，將 folderPath 和 filePath 合併生成完整的文件路徑並返回。
func completeFilePath(folderPath string, filePath string) string {
	// 判斷 filePath 是否為絕對路徑
	if filepath.IsAbs(filePath) {
		// 如果是絕對路徑，直接返回
		return filePath
	}
	// 如果不是絕對路徑，將 folderPath 和 filePath 合併成完整路徑
	return filepath.Join(folderPath, filePath)
}
