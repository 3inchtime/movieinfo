package config

import (
	"fmt"
	"sync"
)

// Manager 配置管理器
type Manager struct {
	config     *Config
	configPath string
	mu         sync.RWMutex
}

// NewManager 创建配置管理器
func NewManager(configPath string) (*Manager, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Manager{
		config:     config,
		configPath: configPath,
	}, nil
}

// Get 获取配置
func (m *Manager) Get() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Reload 重新加载配置
func (m *Manager) Reload() error {
	newConfig, err := LoadConfig(m.configPath)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	m.mu.Lock()
	m.config = newConfig
	m.mu.Unlock()

	fmt.Println("Config reloaded successfully")
	return nil
}

// GetDSN 获取数据库连接字符串
func (m *Manager) GetDSN() string {
	config := m.Get()
	db := config.Database
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
		db.Charset,
	)
}

// GetRedisAddr 获取Redis地址
func (m *Manager) GetRedisAddr() string {
	config := m.Get()
	return fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
}

// IsProduction 判断是否为生产环境
func (m *Manager) IsProduction() bool {
	return m.Get().App.Environment == "production"
}

// IsDevelopment 判断是否为开发环境
func (m *Manager) IsDevelopment() bool {
	return m.Get().App.Environment == "development"
}

// IsTesting 判断是否为测试环境
func (m *Manager) IsTesting() bool {
	return m.Get().App.Environment == "testing"
}
