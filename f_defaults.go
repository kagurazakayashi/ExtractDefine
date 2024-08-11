package main

import (
	"path/filepath"
	"strconv"
	"strings"

	lang "github.com/kagurazakayashi/libNyaruko_Go/nyai18n"
	log "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
	logs "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
)

// loadDefaultsFile 函式從指定的檔案路徑載入預設值檔案。
// path 參數為檔案的路徑字串。
func loadDefaultsFile(path string) {
	// if !fileExists(path) {
	// 	LogC("错误: 文件不存在: ", path)
	// 	return
	// }

	// 讀取指定路徑的檔案內容
	data, err := readFile(path)

	// 如果讀取檔案時發生錯誤，或是檔案內容為空，則記錄錯誤並返回
	if err != nil {
		logs.LogC(logLevel, logs.Error, lang.GetMultilingualText("UnableReadFile")+": ", err)
		return
	}
	if len(data) == 0 {
		path, _ = filepath.Abs(path)
		logs.LogC(logLevel, logs.Warning, lang.GetMultilingualText("EmptyFile")+": ", path)
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
		var logStr string = lang.GetMultilingualText("DefineChange")
		logStr = strings.Replace(logStr, "%ROW%", strconv.Itoa(i), 1)
		logStr = strings.Replace(logStr, "%MODE%", modeStr, 1)
		logStr = strings.Replace(logStr, "%NAME%", kv[0], 1)
		logStr = strings.Replace(logStr, "%VAL%", macroDic[kv[0]], 1)
		logStr = strings.Replace(logStr, "%TOTAL%", strconv.Itoa(len(macroDic)), 1)
		logs.LogC(logLevel, log.Info, logStr)
	}
}
