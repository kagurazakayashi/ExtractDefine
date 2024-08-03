package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	cMakeListsPath string
	cMakeListsDir  string
	detailed       bool
	configFile     string
	filter         []string
	filterStr      string
)

func main() {
	// 定義並解析命令行參數
	// -c 參數：指定配置文件 (.yaml, 會覆蓋命令行參數)。
	flag.StringVar(&configFile, "c", "", "指定一个配置文件 (.yaml 格式, 会覆盖命令行参数)。")
	// -i 參數：指定 CMakeLists.txt 文件的路徑，將從這裡開始搜索宏定義。
	flag.StringVar(&cMakeListsPath, "i", "", "指定一个 CMakeLists.txt 文件，将从这里开始搜索宏定义。")
	// -d 參數：決定是否顯示詳細信息。
	flag.BoolVar(&detailed, "d", false, "是否显示详细信息。")
	// -f 參數：只需要這些宏的資訊。
	flag.StringVar(&filterStr, "f", "", "只需要这些宏的信息(用,分隔)")

	// 解析命令行參數
	flag.Parse()
	if len(filterStr) > 0 {
		filter = strings.Split(filterStr, ",")
	}

	// 載入配置文件
	config, err := loadConfigFile()

	// 根據載入的配置文件設置相應的參數
	loadConfig(config)

	// 如果配置文件載入過程中發生錯誤，則直接返回
	if err != nil {
		return
	}

	if detailed {
		log.Println("C项目生效宏定义提取工具")
	}

	// 檢查指定的 CMakeLists.txt 文件是否存在，如果不存在則輸出錯誤訊息並退出程序
	if !fileExists(cMakeListsPath) {
		log.Println("错误: 指定的 CMakeLists.txt 文件不存在:", cMakeListsPath)
		return
	}

	// 取得 CMakeLists.txt 文件所在的資料夾路徑
	cMakeListsDir = filepath.Dir(cMakeListsPath)
	// 日誌輸出工程文件夾路徑
	if detailed {
		log.Println("工程文件夹: ", cMakeListsDir)
	}

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
	if detailed {
		log.Printf("解析完毕，定义列表 ( %s ):\n", numView)
	}

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
