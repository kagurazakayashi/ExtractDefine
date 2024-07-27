package main

import (
	"strings"
)

// `^\s*set\s*\(\s*(\w+)\s*(.*?)\s*\)\s*$`
func parseSingleLineSet(line string, key string) (string, []string, bool) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, key+"(") && strings.HasSuffix(line, ")") {
		content := strings.TrimSuffix(strings.TrimPrefix(line, key+"("), ")")
		parts := strings.Fields(content)
		if len(parts) > 1 {
			key := parts[0]
			values := parts[1:]
			for i, value := range values {
				values[i] = strings.Trim(value, `"`)
			}
			return key, values, true
		}
	}

	// 處理“set”和“(”之間存在空格的情況
	if strings.HasPrefix(line, key+" ") {
		index := strings.Index(line, "(")
		if index != -1 && strings.HasSuffix(line, ")") {
			content := strings.TrimSuffix(strings.TrimSpace(line[index+1:]), ")")
			parts := strings.Fields(content)
			if len(parts) > 1 {
				key := parts[0]
				values := parts[1:]
				for i, value := range values {
					values[i] = strings.Trim(value, `"`)
				}
				return key, values, true
			}
		}
	}

	return "", nil, false
}

// `^\s*set\s*\(\s*(\w+)\s*$`
func isMultiLineSetStart(line string, key string) (string, bool) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, key+"(") && !strings.HasSuffix(line, ")") {
		content := strings.TrimPrefix(line, key+"(")
		return strings.TrimSpace(content), true
	}

	// 處理“set”和“(”之間存在空格的情況
	if strings.HasPrefix(line, key+" ") {
		index := strings.Index(line, "(")
		if index != -1 && !strings.HasSuffix(line, ")") {
			content := strings.TrimSpace(line[index+1:])
			return content, true
		}
	}

	return "", false
}

// `^\s*\)\s*$`
func isMultiLineSetEnd(line string) bool {
	line = strings.TrimSpace(line)
	return line == ")"
}

func parseCMakeLists(lines []string) map[string][]string {
	result := make(map[string][]string)
	var currentKey string
	var currentValues []string

	for _, line := range lines {
		// 刪除註釋並修剪空格
		line = removeComments(line, "#")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 處理單行 set 語句
		if key, values, matched := parseSingleLineSet(line, "set"); matched {
			result[key] = values
			continue
		}
		if key, values, matched := parseSingleLineSet(line, "idf_component_register"); matched {
			result[key] = values
			continue
		}

		// 處理多行 set 語句
		if key, matched := isMultiLineSetStart(line, "set"); matched {
			currentKey = key
			currentValues = []string{}
		} else if key, matched := isMultiLineSetStart(line, "idf_component_register"); matched {
			currentKey = key
			currentValues = []string{}
		} else if isMultiLineSetEnd(line) {
			if currentKey != "" {
				result[currentKey] = currentValues
				currentKey = ""
			}
		} else if currentKey != "" {
			line = strings.TrimSpace(line)
			if line != "" {
				currentValues = append(currentValues, strings.Trim(line, `"`))
			}
		}
	}

	return result
}
