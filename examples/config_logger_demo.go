package main

import (
	"fmt"
	"log"

	"github.com/3inchtime/movieinfo/internal/config"
	"github.com/3inchtime/movieinfo/pkg/logger"
)

func main() {
	// 测试配置管理系统
	fmt.Println("=== 测试配置管理系统 ===")
	
	// 加载配置
	appConfig, err := config.NewConfig("../configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	fmt.Printf("应用名称: %s\n", appConfig.App.Name)
	fmt.Printf("版本: %s\n", appConfig.App.Version)
	fmt.Printf("环境: %s\n", appConfig.App.Environment)
	fmt.Printf("端口: %d\n", appConfig.App.Port)
	fmt.Printf("数据库主机: %s\n", appConfig.Database.Host)
	fmt.Printf("Redis主机: %s\n", appConfig.Redis.Host)
	
	// 测试日志系统
	fmt.Println("\n=== 测试日志系统 ===")
	
	// 初始化日志系统
	loggerConfig := appConfig.GetLoggerConfig()
	err = logger.Init(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	
	// 测试各种日志级别
	logger.Debug("这是一条调试信息")
	logger.Info("这是一条信息")
	logger.Warn("这是一条警告")
	logger.Error("这是一条错误信息")
	
	// 测试格式化日志
	logger.Infof("应用 %s 版本 %s 启动成功", appConfig.App.Name, appConfig.App.Version)
	
	// 测试带字段的日志
	logger.WithField("user_id", 12345).Info("用户登录成功")
	logger.WithField("request_id", "req-123").WithField("duration", "150ms").Info("请求处理完成")
	
	fmt.Println("\n配置管理系统和日志系统测试完成！")
}