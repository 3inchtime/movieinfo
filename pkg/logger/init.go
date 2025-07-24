package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

// Init 初始化日志系统
func Init(config *Config) error {
	// 确保日志目录存在
	if config.Output == "file" {
		logDir := filepath.Dir(config.File.Path)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// 初始化全局日志器
	return InitGlobalLogger(config)
}
