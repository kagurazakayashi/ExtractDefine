package main

import (
	"strings"

	lang "github.com/kagurazakayashi/libNyaruko_Go/nyai18n"
	logs "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
)

/*
switch varName {
	case "CMAKE_C_FLAGS":
		// CMake 中用於指定編譯 C 程式時傳遞給編譯器的選項的變數。它主要用於定義 C 語言編譯器的標誌。
		break
	case "CMAKE_SOURCE_DIR":
		// 表示 CMakeLists.txt 檔案所在的頂級目錄的完整路徑，也就是專案的根目錄。
		break
	case "CMAKE_CURRENT_LIST_DIR":
		// 表示當前處理的 CMake 指令碼（通常是 CMakeLists.txt 檔案或透過 include 引入的其他 CMake 指令碼）所在目錄的完整路徑。
		break
	case "CMAKE_CURRENT_SOURCE_DIR": //
		// 它表示當前處理的 CMakeLists.txt 檔案所在的源目錄的完整路徑。這個變數會隨之變化。
		break
	case "COMPONENT_DIR": //
		// 通常指的是當前元件的目錄，它的值通常由 ESP-IDF 構建系統自動設定，指向包含該元件的 CMakeLists.txt 檔案的目錄。
		break
	case "ROOT_DIR":
		// 表示專案的根目錄，即 CMakeLists.txt 檔案所在的目錄。
		break
	default:
		break
}
*/

var autoDefineDic map[string]string = map[string]string{
	"EXTRACTDEFINE": "1",
}

// var autoDefineArrDic map[string][]string = map[string][]string{}

// replacePlaceholders 這個函式會將輸入字串中的佔位符取代成對應的值。
// 它會遍歷 macroDic，並將每個佔位符（以${key}的格式）替換為字典中對應的值。
// 輸入：input 字串，可能包含一或多個佔位符。
// 輸出：bool:是否被替換, string:所有佔位符已被替換為相應的值。
func replacePlaceholders(input string) (bool, string) {
	var varList []string = []string{}
	// 找出所有 ${變數名}
	var input2 string = input
	var input3 string = input
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

	// 替換變數
	for _, varName := range varList {
		varValue, exists := macroDic[varName]
		if !exists {
			varValue, exists = autoDefineDic[varName]
			if !exists {
				// 如果字典中沒有對應的值
				var logStr string = lang.GetMultilingualText("NoVar") + ": " + varName
				logs.LogC(logLevel, logs.Warning, logStr)
				// input = strings.ReplaceAll(input, "${"+varName+"}", varName)
				return false, input
			} else {
				var logStr string = lang.GetMultilingualText("AutoVar")
				logStr = strings.Replace(logStr, "%NAME%", varName, 1)
				logStr = strings.Replace(logStr, "%VAL%", varValue, 1)
				logs.LogC(logLevel, logs.Info, logStr)
				input = strings.ReplaceAll(input, "${"+varName+"}", varValue)
			}
		} else {
			var logStr string = lang.GetMultilingualText("ReplaceVar")
			logStr = strings.Replace(logStr, "%NAME%", varName, 1)
			logStr = strings.Replace(logStr, "%VAL%", varValue, 1)
			logs.LogC(logLevel, logs.Debug, logStr)
			input = strings.ReplaceAll(input, "${"+varName+"}", varValue)
		}
	}
	if input != input3 {
		var logStr string = lang.GetMultilingualText("ReplaceVars")
		logStr = strings.Replace(logStr, "%STR0%", input3, 1)
		logStr = strings.Replace(logStr, "%STR1%", input, 1)
		logs.LogC(logLevel, logs.Debug, logStr)
	}
	return true, input
}
