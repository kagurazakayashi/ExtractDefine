package main

import (
	"log"
	"strings"
)

// loadDefaultsFile 函式從指定的檔案路徑載入預設值檔案。
// path 參數為檔案的路徑字串。
func loadDefaultsFile(path string) {
	// 讀取指定路徑的檔案內容
	data, err := readFile(path)

	// 如果讀取檔案時發生錯誤，或是檔案內容為空，則記錄錯誤並返回
	if err != nil || len(data) == 0 {
		log.Println("错误: 无法读取文件: ", err)
		return
	}

	// 解析檔案內容，載入預設值
	parseDefaults(data)
}

// parseDefaults 函式接收一個字串切片（lines），解析每一行並將結果儲存到 macroDic 字典中。
// lines: 包含需要解析的設定檔內容的字串切片
func parseDefaults(lines []string) {
	// 迴圈遍歷每一行資料，進行處理
	for i, line := range lines {
		// 去除行首尾的空白字符
		line = strings.TrimSpace(line)
		// 移除註解內容（以 "#" 開頭的部分）
		line = removeComments(line, "#")
		// 若該行為空，則跳過
		if line == "" {
			continue
		}
		// 將行以 "=" 分割為鍵值對
		var kv []string = strings.Split(line, "=")
		var kvLen int = len(kv)
		// 若分割結果無效，則跳過該行
		if kvLen <= 0 || kvLen > 2 {
			continue
		}
		// modeStr 用來存放處理後的字串
		var modeStr string = ""
		// 如果只有鍵，沒有值，則設定為空字串
		if kvLen == 1 {
			modeStr = macroDicAddStr(kv[0], "")
		} else if kvLen == 2 {
			// 如果有鍵和值，則儲存鍵值對
			modeStr = macroDicAddStr(kv[0], kv[1])
		}
		// 如果 detailed 模式啟用，則輸出詳細資訊
		if detailed {
			log.Printf("从行 %d %s定义 %s , 值为 %s (已存储 %d)\n", i, modeStr, kv[0], macroDic[kv[0]], len(macroDic))
		}
	}
}
