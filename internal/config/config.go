package config

import (
	"github.com/3inchtime/movieinfo/pkg/config"
	"github.com/3inchtime/movieinfo/pkg/logger"
)

// AppConfig 应用内部配置
type AppConfig struct {
	*config.Config
	Logger *logger.Config `yaml:"logger"`
}

// NewConfig 创建新的应用配置
func NewConfig(configPath string) (*AppConfig, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Config: cfg,
		Logger: &logger.Config{
			Level:  cfg.Log.Level,
			Format: cfg.Log.Format,
			Output: cfg.Log.Output,
			File:   logger.FileConfig(cfg.Log.File),
		},
	}, nil
}

// GetLoggerConfig 获取日志配置
func (c *AppConfig) GetLoggerConfig() *logger.Config {
	return c.Logger
}
