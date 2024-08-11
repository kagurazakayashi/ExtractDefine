package main

import (
	"os"
	"path/filepath"

	log "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
	logs "github.com/kagurazakayashi/libNyaruko_Go/nyalog"
	"gopkg.in/yaml.v2"
)

// yaml 配置檔案中的內容
type Config struct {
	ExtractDefine int      `yaml:"ExtractDefine"` // 配置檔案屬於哪個軟體及版本號
	CMakeList     string   `yaml:"CMakeList"`     // 入口 CMakeLists.txt 檔案路徑
	LogLevel      int      `yaml:"LogLevel"`      // 執行過程中是否輸出詳細過程
	Filter        []string `yaml:"Filter"`        // 只需要這些宏的資訊
}

// loadConfigFile 讀取並加載配置文件
// 若指定的配置文件存在且格式正確，則返回解析後的 Config 結構體
// 若配置文件路徑錯誤或無法解析，則返回相應的錯誤
func loadConfigFile() (Config, error) {
	var err error
	var config Config

	// 檢查是否有指定配置文件路徑
	if len(configFile) > 0 {
		// 取得配置文件的絕對路徑
		configFile, err = filepath.Abs(configFile)
		if err != nil {
			// 如果無法取得絕對路徑，記錄錯誤並返回
			logs.LogC(logLevel, log.Clash, "无法读取配置文件:", err)
			return config, err
		}

		// 檢查配置文件是否存在
		if fileExists(configFile) {
			// LogC(logLevel, LogLevelDebug, "读取配置文件:", configFile)
			// 讀取配置文件內容
			cfg, err := os.ReadFile(configFile)
			if err != nil {
				// 如果讀取文件失敗，記錄錯誤並返回
				logs.LogC(logLevel, log.Clash, "无法读取配置文件:", err)
				return config, err
			}

			// 解析 YAML 格式的配置文件
			err = yaml.Unmarshal(cfg, &config)
			if err != nil {
				// 如果解析 YAML 失敗，記錄錯誤並返回
				logs.LogC(logLevel, log.Clash, "解析 YAML 文件失败:", err)
				return config, err
			}

			// 成功解析並返回配置
			return config, nil
		} else {
			// 如果配置文件不存在，記錄錯誤並返回
			logs.LogC(logLevel, log.Clash, "配置文件路径不正确:", configFile)
			return config, err
		}
	}

	// 如果未指定配置文件，返回空配置和 nil 錯誤
	return config, nil
}

// loadConfig 函式從給定的 Config 結構體中載入設定。
// 如果設定中的 ExtractDefine 不是 1，則記錄錯誤訊息並返回。
// @param config Config 結構體，包含所需的設定值
func loadConfig(config Config) {
	// 檢查 ExtractDefine 是否等於 1，若否則記錄錯誤訊息並返回
	if config.ExtractDefine != 1 {
		logs.LogC(logLevel, log.Clash, "配置文件版本", config.ExtractDefine, "不正确。")
		return
	}

	// 將設定值指派給全域變數
	cMakeListsPath = config.CMakeList         // 設定 CMakeList 的路徑
	logLevel = logs.LogLevel(config.LogLevel) // 設定詳細模式
	filter = config.Filter                    // 設定過濾器
}
