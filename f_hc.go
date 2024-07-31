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
				kv[0] = strings.TrimSpace(kv[0])
				log.Printf("从行 %d %s定义 %s , 值为 %s (已存储 %d)\n", i, modeStr, kv[0], macroDic[kv[0]], len(macroDic))
			}
		} else if strings.HasPrefix(line, "#undef") {
			var spaceIndex int = strings.Index(line, " ")
			if spaceIndex < 0 {
				continue
			}
			line = line[spaceIndex+1:]
			var kv []string = strings.Split(line, " ")
			var kvLen int = len(kv)
			if kvLen <= 0 || kvLen > 1 {
				continue
			}
			if len(noSave) > 0 {
				if detailed && noSave != endif {
					log.Printf("忽略行 %d 取消定义 %s , 值为 %s (因为不满足 %s )\n", i, kv[0], macroDic[kv[0]], noSave)
				}
				noSave = ""
				continue
			}
			rmVal := macroDicRemove(kv[0])
			if detailed {
				kv[0] = strings.TrimSpace(kv[0])
				log.Printf("从行 %d 删除定义 %s , 值为 %s (已存储 %d)\n", i, kv[0], rmVal, len(macroDic))
			}
		} else if strings.HasPrefix(line, "#if") {
			var spaceIndex int = strings.Index(line, " ")
			if spaceIndex < 0 {
				continue
			}
			line = line[spaceIndex+1:]
			if strings.HasPrefix(line, "defined") {
				var isExists int8 = -1
				if strings.HasPrefix(line, "||") {
					// ||
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
				} else if strings.HasPrefix(line, "&&") {
					// &&
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
				} else {
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
				if isExists == 1 {
					noSave = ""
				} else if isExists == 0 {
					noSave = ifName
					ifName = ""
				} else {
					continue
				}
			}
		} else if strings.HasPrefix(line, endif) {
			noSave = endif
			ifName = ""
		}
	}
}
