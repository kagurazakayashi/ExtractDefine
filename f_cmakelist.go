package main

import (
	"strings"

	lang "github.com/kagurazakayashi/libNyaruko_Go/nyai18n"
	logs "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
)

// replacePlaceholders 這個函式會將輸入字串中的佔位符取代成對應的值。
// 它會遍歷 macroDic，並將每個佔位符（以${key}的格式）替換為字典中對應的值。
// 輸入：input 字串，可能包含一或多個佔位符。
// 輸出：一個字串，所有佔位符已被替換為相應的值。
func replacePlaceholders(input string) string {
	var varList []string = []string{}
	// 找出所有 ${变量名}
	var input2 = input
	for strings.Contains(input2, "${") {
		var start int = strings.Index(input2, "${")
		var end int = strings.Index(input2, "}")
		if start == -1 || end == -1 {
			break
		}
		var varName string = input2[start+2 : end]
		varList = append(varList, varName)
		input2 = input2[:start] + input2[end+1:]
	}

	// 替换变量
	for _, varName := range varList {
		var varValue string = macroDic[varName]
		if varValue == "" {
			// 如果字典中沒有對應的值
			var logStr string = lang.GetMultilingualText("NoVar") + ": " + varName
			logs.LogC(logLevel, logs.Warning, logStr)
			input = strings.ReplaceAll(input, "${"+varName+"}", varName)
		} else {
			var logStr string = lang.GetMultilingualText("ReplaceVar")
			logStr = strings.Replace(logStr, "%NAME%", varName, 1)
			logStr = strings.Replace(logStr, "%VAL%", varValue, 1)
			logs.LogC(logLevel, logs.Debug, logStr)
			input = strings.ReplaceAll(input, "${"+varName+"}", varValue)
		}
	}
	return input
}

// parseSingleLineSet 解析單行的設定語句，並提取設定的鍵和值。
// 該函式會檢查輸入行是否符合特定格式，並根據格式提取相關資訊。
// 參數:
// - line: 要解析的原始字串。
// - key: 設定的鍵名稱，用來匹配行的前綴。
// 返回值:
// - string: 提取的鍵（如果匹配成功）。
// - []string: 提取的值列表（如果匹配成功）。
// - bool: 如果成功匹配並提取資訊，返回 true；否則返回 false。
func parseSingleLineSet(line string, key string) (string, []string, bool) {
	// `^\s*set\s*\(\s*(\w+)\s*(.*?)\s*\)\s*$`

	// 移除行首和行尾的空白字符
	line = strings.TrimSpace(line)

	// 檢查行是否以“key(”開頭並以“)結尾
	if strings.HasPrefix(line, key+"(") && strings.HasSuffix(line, ")") {
		// 去除鍵的前綴和括號，提取內容
		content := strings.TrimSuffix(strings.TrimPrefix(line, key+"("), ")")
		// 將內容分割成字段
		parts := strings.Fields(content)

		// 如果內容至少包含一個鍵和一個值
		if len(parts) > 1 {
			// 提取鍵
			key := parts[0]
			// 提取值
			values := parts[1:]
			// 去除每個值的引號，填充佔位符
			for i, value := range values {
				values[i] = strings.Trim(value, `"`)
				values[i] = replacePlaceholders(values[i])
			}
			// 返回提取的鍵、值列表和匹配成功的標誌
			return key, values, true
		}
	}

	// 處理“key”和“(”之間存在空格的情況
	if strings.HasPrefix(line, key+" ") {
		// 找到左括號的位置
		index := strings.Index(line, "(")
		// 如果左括號存在且行以右括號結尾
		if index != -1 && strings.HasSuffix(line, ")") {
			// 去除括號之間的空白字符並提取內容
			content := strings.TrimSuffix(strings.TrimSpace(line[index+1:]), ")")
			// 將內容分割成字段
			parts := strings.Fields(content)

			// 如果內容至少包含一個鍵和一個值
			if len(parts) > 1 {
				// 提取鍵
				key := parts[0]
				// 提取值
				values := parts[1:]
				// 去除每個值的引號，填充佔位符
				for i, value := range values {
					values[i] = strings.Trim(value, `"`)
					values[i] = replacePlaceholders(values[i])
				}
				// 返回提取的鍵、值列表和匹配成功的標誌
				return key, values, true
			}
		}
	}

	// 如果不符合任何條件，返回空值和匹配失敗的標誌
	return "", nil, false
}

// isMultiLineSetStart 檢查一行字串是否為多行設定的起始點。
// 這個函數會判斷字串是否以指定的鍵開頭，並且開頭後是否有一個左括號（“(”），
// 但該行並沒有以右括號（“)”）結尾。如果符合條件，則返回括號之後的內容並標示為多行開始。
// 否則，返回空字串和 false。
//
// 參數：
// - line: 需要檢查的行內容。
// - key: 關鍵字，用於判斷該行是否是多行設定的開始。
//
// 返回：
// - 字串：如果符合條件，返回左括號之後的內容；否則返回空字串。
// - 布林值：如果該行符合多行設定的開始條件，返回 true；否則返回 false。
func isMultiLineSetStart(line string, key string) (string, bool) {
	// `^\s*set\s*\(\s*(\w+)\s*$`

	// 移除行前後的空白字元
	line = strings.TrimSpace(line)

	// 檢查該行是否以“key(”開頭，並且不是以“)結尾
	if strings.HasPrefix(line, key+"(") && !strings.HasSuffix(line, ")") {
		// 取出括號之後的內容
		content := strings.TrimPrefix(line, key+"(")
		return strings.TrimSpace(content), true
	}

	// 處理“key”和“(”之間有空格的情況
	if strings.HasPrefix(line, key+" ") {
		// 找到左括號的位置
		index := strings.Index(line, "(")
		if index != -1 && !strings.HasSuffix(line, ")") {
			// 取出左括號之後的內容
			content := strings.TrimSpace(line[index+1:])
			return content, true
		}
	}

	// 如果不符合多行設定的開始條件，返回空字串和 false
	return "", false
}

// isMultiLineSetEnd 函式檢查是否為多行設置的結尾
// 輸入：line 字串，表示要檢查的行內容
// 輸出：布林值，若行內容為 ")"，則返回 true，否則返回 false
func isMultiLineSetEnd(line string) bool {
	// `^\s*\)\s*$`

	// 移除行首和行尾的空白字符
	line = strings.TrimSpace(line)

	// 檢查行內容是否為 ")"，若是則表示多行設置結束
	return line == ")"
}

// parseCMakeLists 解析給定的CMakeLists.txt文件內容，提取set和idf_component_register語句。
// 它將解析的結果以鍵值對的形式存儲在map中。
// 參數:
//   - lines: 包含CMakeLists.txt文件的每一行的字符串切片。
//
// 返回值:
//   - 返回一個map，鍵為語句的名稱（例如: set、idf_component_register），值為對應的值列表。
func parseCMakeLists(lines []string) map[string][]string {
	result := make(map[string][]string) // 儲存解析結果的map
	var currentKey string               // 當前正在解析的語句的鍵
	var currentValues []string          // 當前正在解析的語句的值列表

	for _, line := range lines {
		// 刪除註釋並修剪空白字符
		line = removeComments(line, "#") // 移除行中的註釋
		line = strings.TrimSpace(line)   // 去除行首尾的空白字符
		if line == "" {
			continue // 忽略空行
		}

		// 處理單行 set 語句
		if key, values, matched := parseSingleLineSet(line, "set"); matched {
			result[key] = values // 將解析出的鍵值對存儲到結果中
			continue
		}
		if key, values, matched := parseSingleLineSet(line, "idf_component_register"); matched {
			result[key] = values // 將解析出的idf_component_register鍵值對存儲到結果中
			continue
		}

		// 處理多行 set 語句
		if key, matched := isMultiLineSetStart(line, "set"); matched {
			currentKey = key           // 開始解析多行 set 語句
			currentValues = []string{} // 初始化值列表
		} else if key, matched := isMultiLineSetStart(line, "idf_component_register"); matched {
			currentKey = key           // 開始解析多行 idf_component_register 語句
			currentValues = []string{} // 初始化值列表
		} else if isMultiLineSetEnd(line) {
			if currentKey != "" {
				result[currentKey] = currentValues // 結束多行語句的解析並存儲結果
				currentKey = ""                    // 重置當前鍵
			}
		} else if currentKey != "" {
			line = strings.TrimSpace(line)
			if line != "" {
				currentValues = append(currentValues, strings.Trim(line, `"`)) // 添加解析的值到當前值列表中
			}
		}
	}

	return result // 返回解析結果
}
