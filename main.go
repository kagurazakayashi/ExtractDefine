//go:generate goversioninfo -icon=ico/icon.ico -manifest=main.exe.manifest -arm=true
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	lang "github.com/kagurazakayashi/libNyaruko_Go/nyai18n"
	logs "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
)

var (
	cMakeListsPath   string
	cMakeListsDir    string
	configFile       string
	filter           []string
	filterStr        string
	defaultDefineStr string
	logLevel         logs.LogLevel
)

func main() {
	// 配置語言
	exePath, err := os.Executable()
	if err != nil {
		log.Panicln(err)
		return
	}
	lang.AutoSetLanguage(true)
	exePath = strings.Replace(exePath, ".exe", "", 1) + ".i18n.ini"
	err = lang.LoadLanguageFile(exePath, false)
	if err != nil {
		log.Panicln(err)
	}

	// 定義並解析命令行參數
	// -c 參數：指定配置文件 (.yaml, 會覆蓋命令行參數)。
	flag.StringVar(&configFile, "c", "", lang.GetMultilingualText("c_configFile"))
	// -i 參數：指定 CMakeLists.txt 文件的路徑，將從這裡開始搜索宏定義。
	flag.StringVar(&cMakeListsPath, "i", "", lang.GetMultilingualText("c_cMakeListsPath"))
	// -d 參數：決定是否顯示詳細信息。
	var iLogLevel int
	flag.IntVar(&iLogLevel, "d", 0, lang.GetMultilingualText("c_iLogLevel"))
	// -f 參數：只需要這些宏的資訊。
	flag.StringVar(&filterStr, "f", "", lang.GetMultilingualText("c_filterStr"))
	// -e 參數：指定一些宏。
	flag.StringVar(&defaultDefineStr, "e", "", lang.GetMultilingualText("c_defaultDefine"))

	var cLang string
	flag.StringVar(&cLang, "l", exePath, "Language (defined in "+exePath+" )")

	// 解析命令行參數
	flag.Parse()

	err = lang.LoadLanguageFile(cLang, false)
	if err != nil {
		log.Panicln(err)
	}

	// 配置輸出
	logLevel = logs.LogLevel(iLogLevel)
	if len(filterStr) > 0 {
		filter = strings.Split(filterStr, ",")
	}
	if len(defaultDefineStr) > 0 {
		for _, v := range strings.Split(defaultDefineStr, ",") {
			kv := strings.Split(v, "=")
			if len(kv) == 2 {
				macroDic[kv[0]] = kv[1]
			}
		}
	}

	// 載入配置文件
	config, err := loadConfigFile()

	// 根據載入的配置文件設置相應的參數
	loadConfig(config)

	// 如果配置文件載入過程中發生錯誤，則直接返回
	if err != nil {
		return
	}

	logs.LogC(logLevel, logs.Info, lang.GetMultilingualText("Title"))
	logs.LogC(logLevel, logs.Debug, "https://github.com/kagurazakayashi/ExtractDefine")

	// 檢查指定的 CMakeLists.txt 文件是否存在，如果不存在則輸出錯誤訊息並退出程序
	if !fileExists(cMakeListsPath) {
		logs.LogC(logLevel, logs.Error, "CMakeLists.txt", lang.GetMultilingualText("NoFile"), ":", cMakeListsPath)
		return
	}

	// 取得 CMakeLists.txt 文件所在的資料夾路徑
	cMakeListsDir = filepath.Dir(cMakeListsPath)
	autoDefineDic["CMAKE_C_FLAGS"] = ""
	autoDefineDic["CMAKE_SOURCE_DIR"] = cMakeListsDir
	autoDefineDic["CMAKE_CURRENT_LIST_DIR"] = cMakeListsDir
	autoDefineDic["ROOT_DIR"] = cMakeListsDir
	autoDefineDic["__root_dir"] = cMakeListsDir
	// 日誌輸出工程文件夾路徑
	logs.LogC(logLevel, logs.Debug, lang.GetMultilingualText("Folder"), ":", cMakeListsDir)

	// 生成要解析的文件列表並開始解析
	makeFileList()

	// 過濾器
	var useFilter bool = len(filter) > 0
	var numView string = ""
	var filterDic []string = []string{}
	if useFilter {
		for k, v := range macroDic {
			for _, f := range filter {
				if k == f {
					filterDic = append(filterDic, fmt.Sprintf("%s=%s", k, v))
					break
				}
			}
		}
		numView = fmt.Sprintf("%d / %d", len(filterDic), len(macroDic))
	} else {
		numView = fmt.Sprintf("%d", len(macroDic))
	}

	// 解析完成後，輸出解析到的宏定義的數量
	logs.LogC(logLevel, logs.Info, lang.GetMultilingualText("RunEnd"), "(", numView, "):")
	logs.ResetColorE()

	// 遍歷並輸出所有解析到的宏定義及其對應的值
	if useFilter {
		for _, v := range filterDic {
			fmt.Println(v)
		}
	} else {
		for k, v := range macroDic {
			fmt.Printf("%s=%s\n", k, v)
		}
	}
}

// fileExists 函數用於檢查指定路徑的檔案是否存在。
// 參數：
// - path: 字串類型，表示檔案的路徑。
// 返回值：
// - bool: 如果檔案存在，返回 true；否則返回 false。
func fileExists(path string) bool {
	// 使用 os.Stat 函數來獲取指定路徑的檔案資訊，
	// 並檢查檔案是否存在。
	_, err := os.Stat(path)

	// 如果錯誤為 os.IsNotExist，表示檔案不存在，返回 false。
	if os.IsNotExist(err) {
		return false
	}

	// 如果沒有錯誤（err 為 nil），表示檔案存在，返回 true。
	// 如果出現其他錯誤，返回 false。
	return err == nil
}
