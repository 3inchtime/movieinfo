package logger

// Config 日志配置
type Config struct {
	// 基础配置
	Level  string `yaml:"level" validate:"required,oneof=debug info warn error"`
	Format string `yaml:"format" validate:"required,oneof=json text"`
	Output string `yaml:"output" validate:"required,oneof=stdout stderr file"`

	// 文件输出配置
	File FileConfig `yaml:"file"`
}

// FileConfig 文件输出配置
type FileConfig struct {
	Path       string `yaml:"path" validate:"required"`
	MaxSize    int    `yaml:"max_size" validate:"min=1"`    // MB
	MaxBackups int    `yaml:"max_backups" validate:"min=0"` // 保留文件数
	MaxAge     int    `yaml:"max_age" validate:"min=1"`     // 天数
	Compress   bool   `yaml:"compress"`                     // 是否压缩
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
		File: FileConfig{
			Path:       "logs/app.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
	}
}
