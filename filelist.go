package main

import (
	"bufio"
	"log"
	"os"
)

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func makeFileVals(cMakeListsConfigs map[string][]string, findKey string) []string {
	var result []string = []string{}
	for key, values := range cMakeListsConfigs {
		if key == findKey {
			result = append(result, values...)
		}
	}
	if len(result) > 0 {
		log.Printf("    找到 %s : %s", findKey, result)
	}
	return result
}

func loadCMakeLists(path string) {
	log.Println("加载 CMakeList :", path)
	data, err := readFile(path)
	if err != nil || len(data) == 0 {
		log.Println("错误：无法读取文件 ", path, err)
	}
	var cMakeListsConfigs map[string][]string = ParseCMakeLists(data)
	// log.Println("CMakeList set 列表: ", cMakeListsConfigs)

	/* 在 CMakeLists.txt 檔案中，SDKCONFIG_DEFAULTS 變數用於指定一個或多個預設的 SDK 配置檔案。這些檔案包含了一些預設的配置選項，用於初始化 SDK 的配置。這些配置選項通常是在 sdkconfig 檔案中定義的，用於控制編譯時的一些行為和設定。 */
	var d_SDKCONFIG_DEFAULTS []string = makeFileVals(cMakeListsConfigs, "SDKCONFIG_DEFAULTS")
	if len(d_SDKCONFIG_DEFAULTS) > 0 {
		/* 內容是一個個 .defaults 檔案。
		讀取裡面指定的每個 .defaults 檔案。 */
	}
	/* EXTRA_COMPONENT_DIRS 是在 CMakeLists.txt 檔案中用於指定額外元件目錄的變數。ESP-IDF 專案通常由多個元件組成，這些元件包含在 components 目錄中。預設情況下，CMake 會自動在專案根目錄的 components 目錄中搜索元件。但如果你的元件不在預設的 components 目錄中，或者你有額外的元件目錄，你可以使用 EXTRA_COMPONENT_DIRS 來指定這些目錄。 */
	var d_EXTRA_COMPONENT_DIRS []string = makeFileVals(cMakeListsConfigs, "EXTRA_COMPONENT_DIRS")
	if len(d_EXTRA_COMPONENT_DIRS) > 0 {
		/* 內容是一個個資料夾路徑。
		遍歷資料夾，包括子資料夾。直到找出哪層還有 CMakeLists.txt ，處理這層的 CMakeLists.txt 。 */
	}
	/* 在 CMakeLists.txt 檔案中，include 命令用於包含其他 CMake 指令碼檔案。這些檔案通常包含額外的設定、宏定義、函式定義等，可以在你的 CMake 構建過程中使用。透過包含其他檔案，你可以將 CMake 配置分成多個部分，使其更易於管理和維護。 */
	var d_includes []string = makeFileVals(cMakeListsConfigs, "includes")
	if len(d_includes) > 0 {
		/* 內容是一個個資料夾路徑，這些資料夾裡面是 .h 檔案。
		从这些文件中提取宏定义。 */
	}
	/* 在 CMakeLists.txt 檔案中，SRCS 變數用於指定專案中要編譯的原始檔列表。透過定義 SRCS 變數，CMake 可以知道哪些原始檔需要包含在編譯過程中。 */
	var d_srcs []string = makeFileVals(cMakeListsConfigs, "srcs")
	if len(d_srcs) > 0 {
		/* 內容是一個個 .c 檔案。
		從這些檔案中提取宏定義。 */
	}
}

func makeFileList() {
	loadCMakeLists(cMakeListsPath)
}
