package main

import (
	"log"
	"strings"
)

// loadHCFile 從指定路徑載入 .h .c 檔案
// path 參數是文件的路徑
func loadHCFile(path string) {
	// 嘗試從指定的路徑讀取檔案
	data, err := readFile(path)
	// 如果發生錯誤或讀取到的資料為空，則記錄錯誤訊息並返回
	if err != nil || len(data) == 0 {
		log.Println("错误: 无法读取文件 ", err)
		return
	}
	// 解析讀取到的 HC 檔案資料
	parseHC(data)
}

// parseHC 函數用來解析給定的字串陣列，並依據 C 語言的條件編譯指令進行處理。
// 這個函數會處理 #define、#undef、#if 和 #endif 等指令，並根據條件決定是否保存或移除對應的宏定義。
func parseHC(lines []string) {
	var ifName string = ""        // 用來儲存目前處理中的條件編譯宏名稱
	var noSave string = ""        // 如果設定此變數，則表示目前的條件編譯指令應忽略後續的宏定義
	const endif string = "#endif" // 常數，用來表示 #endif 指令

	// 移除程式碼中的註解，返回去除註解後的字串陣列
	lines = removeCommentsC(lines)

	// 遍歷每一行程式碼
	for i, line := range lines {
		line = strings.TrimSpace(line) // 去除行首尾的空白字符

		// 處理 #define 指令
		if strings.HasPrefix(line, "#define") {
			var spaceIndex int = strings.Index(line, " ") // 找到第一個空白字符的位置
			if spaceIndex < 0 {
				continue // 如果沒有找到空白字符，則跳過該行
			}
			line = line[spaceIndex+1:] // 取得 #define 後面的內容

			// 將 #define 的內容按空白字符分割成鍵值對
			var kv []string = strings.Split(line, " ")
			var kvLen int = len(kv)
			if kvLen <= 0 || kvLen > 2 {
				continue // 如果鍵值對數量不符合要求，則跳過該行
			}

			var modeStr string = ""
			if len(noSave) > 0 { // 如果 noSave 不為空，表示要忽略該行定義
				if detailed && noSave != endif {
					log.Printf("忽略行 %d 定义 %s , 值为 %s (因为不满足 %s )\n", i, kv[0], macroDic[kv[0]], noSave)
				}
				noSave = ""
				continue
			}

			// 根據 kv 的長度決定如何儲存宏定義
			if kvLen == 1 {
				modeStr = macroDicAddStr(kv[0], "")
			} else if kvLen == 2 {
				modeStr = macroDicAddStr(kv[0], kv[1])
			}

			if detailed { // 如果詳細模式啟用，則記錄存儲過程
				kv[0] = strings.TrimSpace(kv[0])
				log.Printf("从行 %d %s定义 %s , 值为 %s (已存储 %d)\n", i, modeStr, kv[0], macroDic[kv[0]], len(macroDic))
			}
		} else if strings.HasPrefix(line, "#undef") { // 處理 #undef 指令
			var spaceIndex int = strings.Index(line, " ") // 找到第一個空白字符的位置
			if spaceIndex < 0 {
				continue // 如果沒有找到空白字符，則跳過該行
			}
			line = line[spaceIndex+1:] // 取得 #undef 後面的內容

			// 將 #undef 的內容按空白字符分割成鍵值對
			var kv []string = strings.Split(line, " ")
			var kvLen int = len(kv)
			if kvLen <= 0 || kvLen > 1 {
				continue // 如果鍵值對數量不符合要求，則跳過該行
			}

			if len(noSave) > 0 { // 如果 noSave 不為空，表示要忽略該行取消定義
				if detailed && noSave != endif {
					log.Printf("忽略行 %d 取消定义 %s , 值为 %s (因为不满足 %s )\n", i, kv[0], macroDic[kv[0]], noSave)
				}
				noSave = ""
				continue
			}

			rmVal := macroDicRemove(kv[0]) // 從宏定義字典中移除對應的鍵值
			if detailed {                  // 如果詳細模式啟用，則記錄刪除過程
				kv[0] = strings.TrimSpace(kv[0])
				log.Printf("从行 %d 删除定义 %s , 值为 %s (已存储 %d)\n", i, kv[0], rmVal, len(macroDic))
			}
		} else if strings.HasPrefix(line, "#if") { // 處理 #if 指令
			var spaceIndex int = strings.Index(line, " ")
			if spaceIndex < 0 {
				continue // 如果沒有找到空白字符，則跳過該行
			}
			line = line[spaceIndex+1:]

			if strings.HasPrefix(line, "defined") { // 處理條件為 defined 的情況
				var isExists int8 = -1

				if strings.HasPrefix(line, "||") { // 處理 OR 條件
					var or []string = strings.Split(line, "||")
					for _, orDef := range or {
						orDef = strings.TrimSpace(orDef)
						var ka []string = strings.Split(orDef, "(")
						if len(ka) != 2 {
							break
						}
						ifName = ka[1][:len(ka[1])-1]
						if macroDicKeyExists(ifName) {
							isExists = 1
							break
						} else {
							isExists = 0
						}
					}
				} else if strings.HasPrefix(line, "&&") { // 處理 AND 條件
					var and []string = strings.Split(line, "&&")
					for _, andDef := range and {
						andDef = strings.TrimSpace(andDef)
						var ka []string = strings.Split(andDef, "(")
						if len(ka) != 2 {
							break
						}
						ifName = ka[1][:len(ka[1])-1]
						if macroDicKeyExists(ifName) {
							isExists = 1
						} else {
							isExists = 0
							break
						}
					}
				} else { // 處理單一 defined 條件
					var ka []string = strings.Split(line, "(")
					if len(ka) != 2 {
						continue
					}
					ifName = ka[1][:len(ka[1])-1]
					if macroDicKeyExists(ifName) {
						isExists = 1
					} else {
						isExists = 0
					}
				}

				// 根據 isExists 的值決定是否忽略後續的定義
				if isExists == 1 {
					noSave = ""
				} else if isExists == 0 {
					noSave = ifName
					ifName = ""
				} else {
					continue
				}
			}
		} else if strings.HasPrefix(line, endif) { // 處理 #endif 指令
			noSave = endif
			ifName = ""
		}
	}
}
