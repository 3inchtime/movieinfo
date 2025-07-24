package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnvString 获取字符串环境变量
func GetEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt 获取整数环境变量
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool 获取布尔环境变量
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetEnvDuration 获取时间间隔环境变量
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// GetEnvStringSlice 获取字符串切片环境变量
func GetEnvStringSlice(key string, defaultValue []string, separator string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, separator)
	}
	return defaultValue
}

// ValidateConfig 验证配置
func ValidateConfig(config *Config) error {
	var errors []string

	// 验证必需的环境变量
	if config.JWT.Secret == "" {
		errors = append(errors, "JWT secret is required")
	}

	if config.Database.Password == "" && !strings.Contains(config.App.Environment, "test") {
		errors = append(errors, "Database password is required for non-test environments")
	}

	// 验证端口范围
	ports := []int{
		config.App.Port,
		config.Database.Port,
		config.Redis.Port,
	}

	for _, port := range ports {
		if port < 1 || port > 65535 {
			errors = append(errors, fmt.Sprintf("Invalid port: %d", port))
		}
	}

	// 验证日志级别
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, config.Log.Level) {
		errors = append(errors, fmt.Sprintf("Invalid log level: %s", config.Log.Level))
	}

	if len(errors) > 0 {
		return fmt.Errorf("config validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PrintConfig 打印配置信息（隐藏敏感信息）
func PrintConfig(config *Config) {
	fmt.Println("=== MovieInfo Configuration ===")
	fmt.Printf("App Name: %s\n", config.App.Name)
	fmt.Printf("Version: %s\n", config.App.Version)
	fmt.Printf("Environment: %s\n", config.App.Environment)
	fmt.Printf("Debug Mode: %t\n", config.App.Debug)
	fmt.Printf("Port: %d\n", config.App.Port)
	fmt.Println()

	fmt.Println("=== Database ===")
	fmt.Printf("Driver: %s\n", config.Database.Driver)
	fmt.Printf("Host: %s:%d\n", config.Database.Host, config.Database.Port)
	fmt.Printf("Database: %s\n", config.Database.Database)
	fmt.Printf("Username: %s\n", config.Database.Username)
	fmt.Printf("Password: %s\n", maskPassword(config.Database.Password))
	fmt.Println()

	fmt.Println("=== Redis ===")
	fmt.Printf("Host: %s:%d\n", config.Redis.Host, config.Redis.Port)
	fmt.Printf("Database: %d\n", config.Redis.Database)
	fmt.Printf("Password: %s\n", maskPassword(config.Redis.Password))
	fmt.Println()

	fmt.Println("=== Log ===")
	fmt.Printf("Level: %s\n", config.Log.Level)
	fmt.Printf("Format: %s\n", config.Log.Format)
	fmt.Printf("Output: %s\n", config.Log.Output)
	if config.Log.Output == "file" {
		fmt.Printf("File Path: %s\n", config.Log.File.Path)
	}
	fmt.Println()

	fmt.Println("=== JWT ===")
	fmt.Printf("Secret: %s\n", maskPassword(config.JWT.Secret))
	fmt.Printf("Expire Time: %s\n", config.JWT.ExpireTime)
	fmt.Printf("Issuer: %s\n", config.JWT.Issuer)
	fmt.Println()
}

// maskPassword 隐藏密码
func maskPassword(password string) string {
	if password == "" {
		return "<empty>"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
