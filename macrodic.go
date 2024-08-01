package main

// 宏定義字典，用來儲存關鍵字與其對應的值
var macroDic map[string]string = make(map[string]string)

// macroDicAdd 函數將一個新的鍵值對添加到宏定義字典中。
// 如果鍵已存在，根據情況返回不同的狀態碼。
// 返回值：
// 0: 成功新增鍵值對。
// 1: 鍵已存在，未進行覆蓋。
// 2: 鍵已存在，且值被覆蓋。
func macroDicAdd(key string, value string) int8 {
	var i int8 = 0
	val, exists := macroDic[key] // 檢查鍵是否已存在
	macroDic[key] = value        // 將鍵值對添加到字典中
	if exists {
		i++ // 鍵已存在，狀態碼加 1
		if val != value {
			i++ // 鍵已存在且值不同，狀態碼再加 1
		}
	}
	return i // 返回狀態碼
}

// macroDicAddStr 函數將一個新的鍵值對添加到宏定義字典中，並根據操作結果返回對應的字串說明。
// 返回值：
// "新增": 成功新增鍵值對。
// "已存在": 鍵已存在，未進行覆蓋。
// "覆蓋": 鍵已存在，且值被覆蓋。
func macroDicAddStr(key string, value string) string {
	var mode int8 = macroDicAdd(key, value)
	if mode == 2 {
		return "覆蓋"
	} else if mode == 1 {
		return "已存在"
	} else {
		return "新增"
	}
}

// macroDicRemove 函數從宏定義字典中移除指定的鍵值對。
// 如果鍵存在，返回被移除的值；如果鍵不存在，返回空字串。
func macroDicRemove(key string) string {
	val, exists := macroDic[key] // 檢查鍵是否存在
	if exists {
		delete(macroDic, key) // 刪除鍵值對
		return val            // 返回被刪除的值
	}
	return "" // 如果鍵不存在，返回空字串
}

// macroDicKeyExists 函數檢查指定的鍵是否存在於宏定義字典中。
// 返回值：
// true: 鍵存在。
// false: 鍵不存在。
func macroDicKeyExists(key string) bool {
	_, exists := macroDic[key] // 檢查鍵是否存在
	return exists
}
