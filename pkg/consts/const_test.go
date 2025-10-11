package consts

import (
	"testing"
)

func TestPathConstants(t *testing.T) {
	// 测试路径常量是否已正确初始化
	if SRC_PATH == "" {
		t.Error("SRC_PATH should not be empty")
	}
	
	if DIR_PATH == "" {
		t.Error("DIR_PATH should not be empty")
	}
	
	if LOG_DIR == "" {
		t.Error("LOG_DIR should not be empty")
	}
	
	if CONFIG_DIR == "" {
		t.Error("CONFIG_DIR should not be empty")
	}
	
	if MODELS_DIR == "" {
		t.Error("MODELS_DIR should not be empty")
	}
	
	if RECORDS_DIR == "" {
		t.Error("RECORDS_DIR should not be empty")
	}
	
	if VIDEOS_DIR == "" {
		t.Error("VIDEOS_DIR should not be empty")
	}
	
	if ASSETS_DIR == "" {
		t.Error("ASSETS_DIR should not be empty")
	}
	
	if RESOURCES_DIR == "" {
		t.Error("RESOURCES_DIR should not be empty")
	}
	
	// 打印路径值用于调试
	t.Logf("SRC_PATH: %s", SRC_PATH)
	t.Logf("DIR_PATH: %s", DIR_PATH)
	t.Logf("LOG_DIR: %s", LOG_DIR)
	t.Logf("CONFIG_DIR: %s", CONFIG_DIR)
	t.Logf("MODELS_DIR: %s", MODELS_DIR)
	t.Logf("RECORDS_DIR: %s", RECORDS_DIR)
	t.Logf("VIDEOS_DIR: %s", VIDEOS_DIR)
	t.Logf("ASSETS_DIR: %s", ASSETS_DIR)
	t.Logf("RESOURCES_DIR: %s", RESOURCES_DIR)
}