package main

import (
	"os"
	"path/filepath"
	"strings"

	lang "github.com/kagurazakayashi/libNyaruko_Go/nyai18n"
	log "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
	logs "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
)

// makeFileVals 根據提供的關鍵字（findKey）從 cMakeListsConfigs 字典中找出對應的值。
// 若找到對應的鍵，則將其值（字符串切片）添加到結果切片中並返回。
// cMakeListsConfigs 是一個字典，鍵為字符串，值為字符串切片。
// findKey 是用來查找的鍵。
// 回傳值為找到的值的字符串切片，若無對應鍵，則回傳空的字符串切片。
func makeFileVals(cMakeListsConfigs map[string][]string, findKey string) []string {
	// 初始化一個空的字符串切片來存儲結果
	var result []string = []string{}

	// 遍歷 cMakeListsConfigs 字典中的每個鍵值對
	for key, values := range cMakeListsConfigs {
		// 如果當前的鍵與要查找的鍵相同
		if key == findKey {
			// 將該鍵對應的值（字符串切片）添加到結果中
			result = append(result, values...)
		}
	}

	// 返回結果切片
	return result
}

// loadCMakeLists 讀取並分析指定路徑下的 CMakeLists.txt 檔案。
// 這個函式會根據 CMakeLists.txt 中的配置項目，處理 SDK 配置檔、元件目錄、
// 包含的其他 CMake 指令碼，以及專案中要編譯的原始檔列表。
func loadCMakeLists(path string) {
	// 如果 detailed 為 true，則記錄當前正在加載的 CMakeList 檔案路徑。
	logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("ProcessFolder")+" CMakeList :", path)
	// 讀取 CMakeLists.txt 檔案的內容。
	data, err := readFile(path)
	if err != nil || len(data) == 0 {
		// 如果讀取失敗，記錄錯誤並返回。
		logs.LogC(logLevel, log.Clash, lang.GetMultilingualText("Error")+":"+lang.GetMultilingualText("UnableFile"), err)
		return
	}

	// 取得檔案所在的目錄。
	var dir string = filepath.Dir(path)
	// 解析 CMakeLists.txt 中的配置項目。
	var cMakeListsConfigs map[string][]string = parseCMakeLists(data)

	// 處理 SDKCONFIG_DEFAULTS 配置項。
	// SDKCONFIG_DEFAULTS 用於指定預設的 SDK 配置檔案。
	var d_SDKCONFIG_DEFAULTS []string = makeFileVals(cMakeListsConfigs, "SDKCONFIG_DEFAULTS")
	if len(d_SDKCONFIG_DEFAULTS) > 0 {
		logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("ProcessItem")+": SDKCONFIG_DEFAULTS :", path)
		// 逐一處理每個 SDK 配置檔案。
		for i, file := range d_SDKCONFIG_DEFAULTS {
			if file == "." || file == ".." {
				continue
			}
			// 補全文件路徑並分析該配置檔案。
			file = completeFilePath(dir, file)
			logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("AnalyzeFile"), i+1, ":", file)
			loadDefaultsFile(file)
		}
	}

	// 處理 EXTRA_COMPONENT_DIRS 配置項。
	// EXTRA_COMPONENT_DIRS 用於指定額外的元件目錄。
	var d_EXTRA_COMPONENT_DIRS []string = makeFileVals(cMakeListsConfigs, "EXTRA_COMPONENT_DIRS")
	if len(d_EXTRA_COMPONENT_DIRS) > 0 {
		logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("ProcessItem")+": EXTRA_COMPONENT_DIRS :", path)
		// 逐一處理每個指定的元件目錄。
		for _, sub := range d_EXTRA_COMPONENT_DIRS {
			if sub == "." || sub == ".." {
				continue
			}
			// 補全資料夾路徑並找到該路徑下的所有 CMakeLists.txt 檔案。
			sub = completeFilePath(dir, sub)
			var dirList []string = findCMakeLists(sub)
			if len(dirList) == 0 {
				continue
			}
			// 對每個找到的 CMakeLists.txt 檔案進行遞迴加載。
			for _, nowDir := range dirList {
				loadCMakeLists(nowDir)
			}
		}
	}

	// 處理 includes 配置項。
	// includes 用於包含其他 CMake 指令碼檔案。
	var d_includes []string = makeFileVals(cMakeListsConfigs, "includes")
	if len(d_includes) > 0 {
		logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("ProcessItem")+": includes:", path)
		// 逐一處理每個包含的 CMake 指令碼檔案。
		for i, sub := range d_includes {
			sub = strings.TrimSpace(sub)
			if sub == "." || sub == ".." {
				continue
			}
			// 補全資料夾路徑並取得該資料夾下的 .h 和 .c 檔案。
			sub = completeFilePath(dir, sub)
			fileList, err := listFilesMatchingPattern(sub, "*.h,*.c")
			if err != nil {
				logs.LogC(logLevel, logs.Error, lang.GetMultilingualText("Error")+":"+lang.GetMultilingualText("UnableFolder"), err)
				return
			}
			logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("Folder"), i+1, ":", sub)
			// 對每個找到的檔案進行分析。
			for i, file := range fileList {
				logs.LogC(logLevel, logs.Debug, "includes:", lang.GetMultilingualText("AnalyzeFile"), i+1, ":", file)
				loadHCFile(file)
			}
		}
	}

	// 處理 srcs 配置項。
	// SRCS 用於指定專案中要編譯的原始檔列表。
	var d_srcs []string = makeFileVals(cMakeListsConfigs, "srcs")
	if len(d_srcs) > 0 {
		logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("ProcessItem")+": srcs:", path)
		// 逐一處理每個 .c 檔案。
		for i, sub := range d_srcs {
			sub = strings.TrimSpace(sub)
			sub = completeFilePath(dir, sub)
			if strings.HasSuffix(sub, ".c") {
				repOK, repVal := replacePlaceholders(sub)
				if !repOK {
					logs.LogC(logLevel, logs.Warning, "srcs:", lang.GetMultilingualText("NoVarInter"), i+1, ":", sub)
					continue
				}
				sub = repVal
				logs.LogC(logLevel, logs.Debug, "srcs:", lang.GetMultilingualText("AnalyzeFile"), i+1, ":", sub)
				loadHCFile(sub)
			}
		}
	}
}

// 建立檔案列表
func makeFileList() {
	loadCMakeLists(cMakeListsPath)
}

// findCMakeLists 函數會在指定的根目錄下遞迴搜尋所有的文件，
// 並返回一個包含所有 CMakeLists.txt 文件絕對路徑的字串切片。
// 如果在遍歷過程中發生錯誤，會記錄錯誤訊息並返回一個空的字串切片。
func findCMakeLists(root string) []string {
	// 初始化一個字串切片，用來儲存找到的 CMakeLists.txt 文件路徑
	var cmakeFiles []string = []string{}

	// 使用 filepath.Walk 函數遞迴遍歷指定的根目錄
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// 如果遍歷某個文件或文件夾時發生錯誤，記錄錯誤訊息並返回錯誤
		if err != nil {
			logs.LogC(logLevel, logs.Error, lang.GetMultilingualText("UnableFolder")+":", err)
			return err
		}
		// 檢查當前文件是否為非目錄且名稱為 "CMakeLists.txt"
		if !info.IsDir() && info.Name() == "CMakeLists.txt" {
			// 移除路徑兩端的空白字符
			path = strings.TrimSpace(path)
			// 將路徑轉換為絕對路徑
			absPath, err := filepath.Abs(path)
			// 如果轉換過程中出現錯誤，記錄錯誤訊息並返回錯誤
			if err != nil {
				logs.LogC(logLevel, logs.Error, lang.GetMultilingualText("UnableFile")+":", err)
				return err
			}
			// 將絕對路徑加入到切片中
			cmakeFiles = append(cmakeFiles, absPath)
		}
		// 返回 nil 表示遍歷過程沒有發生錯誤
		return nil
	})

	// 如果在遍歷過程中出現錯誤，記錄錯誤訊息並返回一個空的字串切片
	if err != nil {
		logs.LogC(logLevel, logs.Error, lang.GetMultilingualText("UnableFolder")+":", err)
		return []string{}
	}

	// 返回包含所有找到的 CMakeLists.txt 文件絕對路徑的字串切片
	return cmakeFiles
}
