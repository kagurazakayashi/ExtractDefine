package main

import (
	"log"
	"strings"
)

func loadHCFile(path string) {
	data, err := readFile(path)
	if err != nil || len(data) == 0 {
		log.Println("错误：无法读取文件 ", err)
		return
	}
	parseHC(data)
}

func parseHC(lines []string) {
	// #if defined(CONFIG_EXAMPLE_CONNECT_IPV6)
	// #define CFG_WEB_NET_IPV6 // 支持IPv6
	// #endif
	var ifName string = ""
	var noSave string = ""
	const endif string = "#endif"
	lines = removeCommentsC(lines)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#define") {
			var spaceIndex int = strings.Index(line, " ")
			if spaceIndex < 0 {
				continue
			}
			line = line[spaceIndex+1:]
			var kv []string = strings.Split(line, " ")
			var kvLen int = len(kv)
			if kvLen <= 0 || kvLen > 2 {
				continue
			}
			var modeStr string = ""
			if len(noSave) > 0 {
				if detailed && noSave != endif {
					log.Printf("忽略行 %d 定义 %s , 值为 %s (因为不满足 %s )\n", i, kv[0], macroDic[kv[0]], noSave)
				}
				noSave = ""
				continue
			}
			if kvLen == 1 {
				modeStr = macroDicAddStr(kv[0], "")
			} else if kvLen == 2 {
				modeStr = macroDicAddStr(kv[0], kv[1])
			}
			if detailed {
				log.Printf("从行 %d %s定义 %s , 值为 %s (已存储 %d)\n", i, modeStr, kv[0], macroDic[kv[0]], len(macroDic))
			}
		} else if strings.HasPrefix(line, "#if") {
			var spaceIndex int = strings.Index(line, " ")
			if spaceIndex < 0 {
				continue
			}
			line = line[spaceIndex+1:]
			if strings.HasPrefix(line, "defined") {
				var ka []string = strings.Split(line, "(")
				if len(ka) != 2 {
					continue
				}
				ifName = ka[1][:len(ka[1])-1]
				if macroDicKeyExists(ifName) {
					noSave = ""
				} else {
					noSave = ifName
					ifName = ""
				}
			}
		} else if strings.HasPrefix(line, endif) {
			noSave = endif
			ifName = ""
		}
	}
}
