// pkg/config/loader.go
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 设置环境变量
	v.SetEnvPrefix("MOVIEINFO")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 绑定环境变量
	bindEnvVars(v)

	// 解析配置
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 应用默认值
	applyDefaults(&config)

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// bindEnvVars 绑定环境变量
func bindEnvVars(v *viper.Viper) {
	v.BindEnv("database.password", "MOVIEINFO_DATABASE_PASSWORD")
	v.BindEnv("redis.password", "MOVIEINFO_REDIS_PASSWORD")
	v.BindEnv("jwt.secret", "MOVIEINFO_JWT_SECRET")
}

// applyDefaults 应用默认值
func applyDefaults(config *Config) {
	if config.App.Port == 0 {
		config.App.Port = 8080
	}

	// 数据库默认值
	if config.Database.Charset == "" {
		config.Database.Charset = "utf8mb4"
	}
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 100
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 10
	}

	// 日志默认值
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Log.Format == "" {
		config.Log.Format = "json"
	}
	if config.Log.Output == "" {
		config.Log.Output = "stdout"
	}
	if config.Log.File.MaxSize == 0 {
		config.Log.File.MaxSize = 100
	}
	if config.Log.File.MaxBackups == 0 {
		config.Log.File.MaxBackups = 10
	}
	if config.Log.File.MaxAge == 0 {
		config.Log.File.MaxAge = 30
	}

	// JWT默认值
	if config.JWT.ExpireTime == 0 {
		config.JWT.ExpireTime = 24 * time.Hour
	}
	if config.JWT.Issuer == "" {
		config.JWT.Issuer = "movieinfo"
	}
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 额外的业务逻辑验证
	if config.Database.Password == "" && !strings.Contains(config.App.Environment, "test") {
		return fmt.Errorf("database password is required for non-test environments")
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}
