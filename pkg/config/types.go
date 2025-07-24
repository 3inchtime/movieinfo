package config

import "time"

// Config 应用配置结构
type Config struct {
	App      AppConfig      `yaml:"app" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log" validate:"required"`
	JWT      JWTConfig      `yaml:"jwt" validate:"required"`
}

// AppConfig 应用基础配置
type AppConfig struct {
	Name        string `yaml:"name" validate:"required"`
	Version     string `yaml:"version" validate:"required"`
	Environment string `yaml:"environment" validate:"required,oneof=development testing production"`
	Debug       bool   `yaml:"debug"`
	Port        int    `yaml:"port" validate:"required,min=1,max=65535"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `yaml:"driver" validate:"required"`
	Host         string `yaml:"host" validate:"required"`
	Port         int    `yaml:"port" validate:"required,min=1,max=65535"`
	Username     string `yaml:"username" validate:"required"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database" validate:"required"`
	Charset      string `yaml:"charset"`
	MaxOpenConns int    `yaml:"max_open_conns" validate:"min=1"`
	MaxIdleConns int    `yaml:"max_idle_conns" validate:"min=1"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required,min=1,max=65535"`
	Password string `yaml:"password"`
	Database int    `yaml:"database" validate:"min=0"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string     `yaml:"level" validate:"required,oneof=debug info warn error"`
	Format string     `yaml:"format" validate:"required,oneof=json text"`
	Output string     `yaml:"output" validate:"required,oneof=stdout stderr file"`
	File   FileConfig `yaml:"file"`
}

// FileConfig 文件输出配置
type FileConfig struct {
	Path       string `yaml:"path" validate:"required"`
	MaxSize    int    `yaml:"max_size" validate:"min=1"`    // MB
	MaxBackups int    `yaml:"max_backups" validate:"min=0"` // 保留文件数
	MaxAge     int    `yaml:"max_age" validate:"min=1"`     // 天数
	Compress   bool   `yaml:"compress"`                     // 是否压缩
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	ExpireTime time.Duration `yaml:"expire_time" validate:"required"`
	Issuer     string        `yaml:"issuer" validate:"required"`
}
