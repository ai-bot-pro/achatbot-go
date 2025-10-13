package main

import (
	"fmt"
	"os"

	"achatbot/pkg/consts"
)

func main() {
	fmt.Println("Project Path Constants Example")
	fmt.Println("==============================")

	// 显示所有路径常量
	fmt.Printf("SRC_PATH: %s\n", consts.SRC_PATH)
	fmt.Printf("DIR_PATH: %s\n", consts.DIR_PATH)
	fmt.Printf("LOG_DIR: %s\n", consts.LOG_DIR)
	fmt.Printf("CONFIG_DIR: %s\n", consts.CONFIG_DIR)
	fmt.Printf("MODELS_DIR: %s\n", consts.MODELS_DIR)
	fmt.Printf("RECORDS_DIR: %s\n", consts.RECORDS_DIR)
	fmt.Printf("VIDEOS_DIR: %s\n", consts.VIDEOS_DIR)
	fmt.Printf("ASSETS_DIR: %s\n", consts.ASSETS_DIR)
	fmt.Printf("RESOURCES_DIR: %s\n", consts.RESOURCES_DIR)

	// 演示如何使用这些路径
	fmt.Println("\nUsing path constants:")

	// 检查日志目录是否存在，如果不存在则创建
	for _, dir := range []string{
		consts.LOG_DIR, consts.RECORDS_DIR, consts.VIDEOS_DIR, consts.CONFIG_DIR, consts.MODELS_DIR,
		consts.ASSETS_DIR, consts.RESOURCES_DIR} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("Creating log directory: %s\n", dir)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("Error creating directory: %v\n", err)
			} else {
				fmt.Printf("Directory created successfully: %s\n", dir)
			}
		} else {
			fmt.Printf("Log directory already exists: %s\n", dir)
		}

	}
}
