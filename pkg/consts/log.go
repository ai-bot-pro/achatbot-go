package consts

import (
	"os"
	"path/filepath"
	"runtime"
)

// 定义项目路径常量
var (
	// SRC_PATH 源代码路径
	SRC_PATH string

	// DIR_PATH 项目根路径
	DIR_PATH string

	// 日志目录
	LOG_DIR string

	// 配置目录
	CONFIG_DIR string

	// 模型目录
	MODELS_DIR string

	// 录音目录
	RECORDS_DIR string

	// 视频目录
	VIDEOS_DIR string

	// 资源目录
	ASSETS_DIR string

	// 资源文件目录
	RESOURCES_DIR string
)

func init() {
	// 获取当前文件的目录作为源代码路径
	_, currentFile, _, _ := runtime.Caller(0)
	SRC_PATH = filepath.Dir(currentFile)

	// 项目根路径是源代码路径的上级目录
	DIR_PATH = filepath.Join(SRC_PATH, "..", "..")

	// 检查环境变量
	if os.Getenv("ACHATBOT_PKG") != "" {
		homeDir, _ := os.UserHomeDir()
		DIR_PATH = filepath.Join(homeDir, ".achatbot-go")
	}

	// 初始化各目录路径
	LOG_DIR = filepath.Join(DIR_PATH, "log")
	CONFIG_DIR = filepath.Join(DIR_PATH, "config")
	MODELS_DIR = filepath.Join(DIR_PATH, "models")
	RECORDS_DIR = filepath.Join(DIR_PATH, "records")
	VIDEOS_DIR = filepath.Join(DIR_PATH, "videos")
	ASSETS_DIR = filepath.Join(DIR_PATH, "assets")
	RESOURCES_DIR = filepath.Join(DIR_PATH, "resources")
}
