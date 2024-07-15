package main

var macroDic map[string]string = make(map[string]string)

func macroDicAdd(key string, value string) bool {
	val, exists := macroDic[key]
	macroDic[key] = value
	return exists && len(val) > 0
}

func macroDicAddStr(key string, value string) string {
	if macroDicAdd(key, value) {
		return "覆盖(!)定义"
	} else {
		return "新增定义"
	}
}
