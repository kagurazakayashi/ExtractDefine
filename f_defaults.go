package main

import (
	"log"
	"strings"
)

func loadDefaultsFile(path string) {
	data, err := readFile(path)
	if err != nil || len(data) == 0 {
		log.Println("错误：无法读取文件 ", err)
		return
	}
	parseDefaults(data)
}

func parseDefaults(lines []string) {
	// macroDic
	for i, line := range lines {
		line = strings.TrimSpace(line)
		line = removeComments(line, "#")
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
		if detailed {
			log.Printf("从行 %d %s定义 %s , 值为 %s (已存储 %d)\n", i, modeStr, kv[0], macroDic[kv[0]], len(macroDic))
		}
	}
}
