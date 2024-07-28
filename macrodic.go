package main

var macroDic map[string]string = make(map[string]string)

// / 0:已添加 1:已存在 2:已覆盖
func macroDicAdd(key string, value string) int8 {
	var i int8 = 0
	val, exists := macroDic[key]
	macroDic[key] = value
	if exists {
		i++
		if val != value {
			i++
		}
	}
	return i
}

func macroDicAddStr(key string, value string) string {
	var mode int8 = macroDicAdd(key, value)
	if mode == 2 {
		return "覆盖(!)"
	} else if mode == 1 {
		return "已存在"
	} else {
		return "新增"
	}
}

func macroDicKeyExists(key string) bool {
	_, exists := macroDic[key]
	return exists
}
