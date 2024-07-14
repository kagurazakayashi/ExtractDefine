package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	cMakeListsPath string
	cMakeListsDir  string
)

func main() {
	log.Println("宏定义提取工具")
	flag.StringVar(&cMakeListsPath, "i", "./CMakeLists.txt", "指定一个 CMakeLists.txt 文件，将从这里开始搜索宏定义。")
	flag.Parse()
	var err error
	cMakeListsPath, err = filepath.Abs(cMakeListsPath)
	if err != nil || !fileExists(cMakeListsPath) {
		log.Fatalf("错误：指定的 CMakeLists.txt 文件不存在: %s", cMakeListsPath)
	}
	cMakeListsDir = filepath.Dir(cMakeListsPath)
	log.Println("工程文件夹: ", cMakeListsDir)
	makeFileList()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
