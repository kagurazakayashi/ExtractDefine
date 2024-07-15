package main

import (
	"log"
	"strings"
)

func loadDefaultsFile(path string) {
	data, err := readFile(path)
	if err != nil || len(data) == 0 {
		log.Println("错误：无法读取文件 ", path, err)
		return
	}
	parseDefaults(data)
}

func parseDefaults(lines []string) {
	// macroDic
	for i, line := range lines {
		line = removeComments(line, "#")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var kv []string = strings.Split(line, "=")
		var kvLen int = len(kv)
		if kvLen <= 0 || kvLen > 2 {
			continue
		}
		var modeStr string = ""
		if kvLen == 1 {
			modeStr = macroDicAddStr(kv[0], "")
		} else if kvLen == 2 {
			modeStr = macroDicAddStr(kv[0], kv[1])
		}
		log.Printf("在行 %d %s %s , 值为 %s (已存储 %d)\n", i, modeStr, kv[0], macroDic[kv[0]], len(macroDic))
	}
}
