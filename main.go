package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	cMakeListsPath string
	cMakeListsDir  string
	detailed       bool
)

func main() {
	log.Println("宏定义提取工具")
	flag.StringVar(&cMakeListsPath, "i", "./CMakeLists.txt", "指定一个 CMakeLists.txt 文件，将从这里开始搜索宏定义。")
	flag.BoolVar(&detailed, "d", false, "是否显示详细信息。")
	flag.Parse()
	var err error
	cMakeListsPath, err = filepath.Abs(cMakeListsPath)
	if err != nil || !fileExists(cMakeListsPath) {
		log.Fatalf("错误：指定的 CMakeLists.txt 文件不存在: %s", cMakeListsPath)
	}
	cMakeListsDir = filepath.Dir(cMakeListsPath)
	log.Println("工程文件夹: ", cMakeListsDir)
	makeFileList()
	log.Printf("解析完毕，定义列表 (%d) :", len(macroDic))
	// var maxNumLen int = len(strconv.Itoa(len(macroDic)))
	// var i int = 0
	for k, v := range macroDic {
		// i++
		// var numFormat string = fmt.Sprintf("%%%dd", maxNumLen)
		// numFormat = fmt.Sprintf(numFormat, i)
		// fmt.Printf("%s  %s=%s\n", numFormat, k, v)
		fmt.Printf("%s=%s\n", k, v)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
