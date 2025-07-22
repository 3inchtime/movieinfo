# 4.4 数据库迁移

## 4.4.1 概述

数据库迁移是管理数据库结构变更的重要机制，它确保数据库结构在不同环境和版本间的一致性。对于MovieInfo项目，我们需要实现一个完整的迁移系统，支持版本控制、自动执行和回滚操作。

## 4.4.2 为什么需要数据库迁移？

### 4.4.2.1 **版本控制**
- **结构版本化**：将数据库结构变更纳入版本控制
- **变更追踪**：记录每次数据库结构的变更历史
- **团队协作**：确保团队成员使用相同的数据库结构
- **环境一致性**：保证开发、测试、生产环境的一致性

### 4.4.2.2 **部署自动化**
- **自动执行**：部署时自动执行数据库变更
- **依赖管理**：管理迁移脚本间的依赖关系
- **状态检查**：检查数据库当前状态和目标状态
- **错误处理**：迁移失败时的错误处理和恢复

### 4.4.2.3 **数据安全**
- **备份机制**：迁移前自动备份数据
- **回滚支持**：支持迁移的回滚操作
- **事务保护**：在事务中执行迁移操作
- **验证检查**：迁移后的数据完整性验证

### 4.4.2.4 **运维支持**
- **批量执行**：支持批量执行多个迁移
- **状态监控**：监控迁移执行状态和进度
- **日志记录**：详细的迁移执行日志
- **性能优化**：大数据量迁移的性能优化

## 4.4.3 迁移系统架构设计

### 4.4.3.1 **迁移系统组件**

```
数据库迁移系统
├── 迁移管理器 (Migration Manager)
│   ├── 迁移执行器 (Migration Executor)
│   ├── 版本管理器 (Version Manager)
│   ├── 状态跟踪器 (State Tracker)
│   └── 回滚管理器 (Rollback Manager)
├── 迁移定义 (Migration Definitions)
│   ├── 结构迁移 (Schema Migrations)
│   ├── 数据迁移 (Data Migrations)
│   ├── 索引迁移 (Index Migrations)
│   └── 约束迁移 (Constraint Migrations)
├── 迁移工具 (Migration Tools)
│   ├── SQL生成器 (SQL Generator)
│   ├── 备份工具 (Backup Tool)
│   ├── 验证工具 (Validation Tool)
│   └── 性能分析器 (Performance Analyzer)
└── 监控系统 (Monitoring System)
    ├── 执行监控 (Execution Monitor)
    ├── 进度跟踪 (Progress Tracker)
    ├── 错误报告 (Error Reporter)
    └── 性能统计 (Performance Stats)
```

### 4.4.3.2 **迁移文件结构**

#### 4.4.3.2.1 迁移文件命名规范
```
迁移文件命名格式: YYYYMMDDHHMMSS_description.sql
例如:
├── 20240101120000_create_users_table.sql
├── 20240101120100_create_movies_table.sql
├── 20240101120200_create_comments_table.sql
├── 20240101120300_create_ratings_table.sql
├── 20240101120400_add_indexes.sql
└── 20240101120500_insert_initial_data.sql
```

#### 4.4.3.2.2 迁移文件内容结构
```sql
-- migrations/20240101120000_create_users_table.sql

-- +migrate Up
-- 创建用户表
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    bio TEXT,
    gender TINYINT,
    birthday DATE,
    location VARCHAR(100),
    website VARCHAR(255),
    status TINYINT DEFAULT 1,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    last_login_ip VARCHAR(45),
    movies_watched INT DEFAULT 0,
    reviews_count INT DEFAULT 0,
    followers_count INT DEFAULT 0,
    following_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
-- 删除用户表
DROP TABLE IF EXISTS users;
```

### 4.4.3.3 **迁移管理器实现**

#### 4.4.3.3.1 迁移接口定义
```go
// internal/migration/interface.go
package migration

import (
    "context"
    "time"
)

// Migration 迁移接口
type Migration interface {
    // 获取迁移信息
    GetID() string
    GetDescription() string
    GetTimestamp() time.Time
    
    // 执行迁移
    Up(ctx context.Context, executor Executor) error
    Down(ctx context.Context, executor Executor) error
    
    // 验证迁移
    Validate(ctx context.Context, executor Executor) error
}

// Executor 迁移执行器接口
type Executor interface {
    // SQL执行
    Exec(ctx context.Context, query string, args ...interface{}) error
    Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error)
    
    // 事务管理
    Begin(ctx context.Context) (Transaction, error)
    
    // 状态查询
    GetCurrentVersion(ctx context.Context) (string, error)
    SetCurrentVersion(ctx context.Context, version string) error
    
    // 迁移历史
    GetMigrationHistory(ctx context.Context) ([]MigrationRecord, error)
    AddMigrationRecord(ctx context.Context, record *MigrationRecord) error
}

// Transaction 事务接口
type Transaction interface {
    Exec(ctx context.Context, query string, args ...interface{}) error
    Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error)
    Commit() error
    Rollback() error
}

// MigrationRecord 迁移记录
type MigrationRecord struct {
    ID          string    `json:"id" db:"id"`
    Description string    `json:"description" db:"description"`
    Timestamp   time.Time `json:"timestamp" db:"timestamp"`
    ExecutedAt  time.Time `json:"executed_at" db:"executed_at"`
    ExecutionTime time.Duration `json:"execution_time" db:"execution_time"`
    Status      string    `json:"status" db:"status"` // success, failed, rollback
    Error       string    `json:"error,omitempty" db:"error"`
}

// Manager 迁移管理器接口
type Manager interface {
    // 迁移操作
    Migrate(ctx context.Context) error
    MigrateTo(ctx context.Context, targetVersion string) error
    Rollback(ctx context.Context, steps int) error
    RollbackTo(ctx context.Context, targetVersion string) error
    
    // 状态查询
    GetCurrentVersion(ctx context.Context) (string, error)
    GetPendingMigrations(ctx context.Context) ([]Migration, error)
    GetMigrationHistory(ctx context.Context) ([]MigrationRecord, error)
    
    // 验证操作
    ValidateAll(ctx context.Context) error
    ValidateMigration(ctx context.Context, migrationID string) error
    
    // 工具操作
    GenerateMigration(description string) (string, error)
    CreateMigrationFile(id, description string) (string, error)
}
```

#### 4.4.3.3.2 迁移管理器实现
```go
// internal/migration/manager.go
package migration

import (
    "context"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "sort"
    "strings"
    "time"
    
    "github.com/yourname/movieinfo/pkg/logger"
)

// ManagerImpl 迁移管理器实现
type ManagerImpl struct {
    executor      Executor
    migrationsDir string
    migrations    map[string]Migration
    logger        logger.Logger
}

// NewManager 创建迁移管理器
func NewManager(executor Executor, migrationsDir string) *ManagerImpl {
    return &ManagerImpl{
        executor:      executor,
        migrationsDir: migrationsDir,
        migrations:    make(map[string]Migration),
        logger:        logger.GetGlobalLogger(),
    }
}

// LoadMigrations 加载迁移文件
func (m *ManagerImpl) LoadMigrations() error {
    files, err := ioutil.ReadDir(m.migrationsDir)
    if err != nil {
        return fmt.Errorf("failed to read migrations directory: %w", err)
    }
    
    for _, file := range files {
        if !strings.HasSuffix(file.Name(), ".sql") {
            continue
        }
        
        migration, err := m.parseMigrationFile(filepath.Join(m.migrationsDir, file.Name()))
        if err != nil {
            return fmt.Errorf("failed to parse migration file %s: %w", file.Name(), err)
        }
        
        m.migrations[migration.GetID()] = migration
    }
    
    m.logger.Info("Loaded migrations",
        logger.Int("count", len(m.migrations)),
    )
    
    return nil
}

// Migrate 执行所有待执行的迁移
func (m *ManagerImpl) Migrate(ctx context.Context) error {
    pendingMigrations, err := m.GetPendingMigrations(ctx)
    if err != nil {
        return fmt.Errorf("failed to get pending migrations: %w", err)
    }
    
    if len(pendingMigrations) == 0 {
        m.logger.Info("No pending migrations")
        return nil
    }
    
    m.logger.Info("Starting migration",
        logger.Int("pending_count", len(pendingMigrations)),
    )
    
    for _, migration := range pendingMigrations {
        if err := m.executeMigration(ctx, migration, true); err != nil {
            return fmt.Errorf("failed to execute migration %s: %w", migration.GetID(), err)
        }
    }
    
    m.logger.Info("Migration completed successfully")
    return nil
}

// MigrateTo 迁移到指定版本
func (m *ManagerImpl) MigrateTo(ctx context.Context, targetVersion string) error {
    currentVersion, err := m.GetCurrentVersion(ctx)
    if err != nil {
        return fmt.Errorf("failed to get current version: %w", err)
    }
    
    if currentVersion == targetVersion {
        m.logger.Info("Already at target version", logger.String("version", targetVersion))
        return nil
    }
    
    // 获取需要执行的迁移
    migrations := m.getMigrationsToExecute(currentVersion, targetVersion)
    
    if len(migrations) == 0 {
        m.logger.Info("No migrations to execute")
        return nil
    }
    
    // 判断是向上迁移还是向下迁移
    isUpgrade := targetVersion > currentVersion
    
    m.logger.Info("Starting targeted migration",
        logger.String("from", currentVersion),
        logger.String("to", targetVersion),
        logger.Bool("upgrade", isUpgrade),
        logger.Int("migration_count", len(migrations)),
    )
    
    for _, migration := range migrations {
        if err := m.executeMigration(ctx, migration, isUpgrade); err != nil {
            return fmt.Errorf("failed to execute migration %s: %w", migration.GetID(), err)
        }
    }
    
    m.logger.Info("Targeted migration completed successfully")
    return nil
}

// Rollback 回滚指定步数的迁移
func (m *ManagerImpl) Rollback(ctx context.Context, steps int) error {
    if steps <= 0 {
        return fmt.Errorf("rollback steps must be positive")
    }
    
    history, err := m.GetMigrationHistory(ctx)
    if err != nil {
        return fmt.Errorf("failed to get migration history: %w", err)
    }
    
    // 按执行时间倒序排列
    sort.Slice(history, func(i, j int) bool {
        return history[i].ExecutedAt.After(history[j].ExecutedAt)
    })
    
    if len(history) < steps {
        return fmt.Errorf("cannot rollback %d steps, only %d migrations in history", steps, len(history))
    }
    
    m.logger.Info("Starting rollback",
        logger.Int("steps", steps),
    )
    
    for i := 0; i < steps; i++ {
        record := history[i]
        migration, exists := m.migrations[record.ID]
        if !exists {
            return fmt.Errorf("migration %s not found", record.ID)
        }
        
        if err := m.executeMigration(ctx, migration, false); err != nil {
            return fmt.Errorf("failed to rollback migration %s: %w", migration.GetID(), err)
        }
    }
    
    m.logger.Info("Rollback completed successfully")
    return nil
}

// executeMigration 执行单个迁移
func (m *ManagerImpl) executeMigration(ctx context.Context, migration Migration, isUp bool) error {
    startTime := time.Now()
    
    // 开始事务
    tx, err := m.executor.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // 创建事务执行器
    txExecutor := &TransactionExecutor{tx: tx}
    
    // 执行迁移
    if isUp {
        err = migration.Up(ctx, txExecutor)
    } else {
        err = migration.Down(ctx, txExecutor)
    }
    
    if err != nil {
        m.logger.Error("Migration execution failed",
            logger.String("migration_id", migration.GetID()),
            logger.Bool("is_up", isUp),
            logger.Error(err),
        )
        return err
    }
    
    // 更新迁移记录
    record := &MigrationRecord{
        ID:            migration.GetID(),
        Description:   migration.GetDescription(),
        Timestamp:     migration.GetTimestamp(),
        ExecutedAt:    time.Now(),
        ExecutionTime: time.Since(startTime),
        Status:        "success",
    }
    
    if isUp {
        if err := m.executor.AddMigrationRecord(ctx, record); err != nil {
            return fmt.Errorf("failed to add migration record: %w", err)
        }
        
        if err := m.executor.SetCurrentVersion(ctx, migration.GetID()); err != nil {
            return fmt.Errorf("failed to set current version: %w", err)
        }
    } else {
        // 回滚时删除迁移记录
        if err := m.removeMigrationRecord(ctx, migration.GetID()); err != nil {
            return fmt.Errorf("failed to remove migration record: %w", err)
        }
        
        // 设置为前一个版本
        previousVersion, err := m.getPreviousVersion(ctx, migration.GetID())
        if err != nil {
            return fmt.Errorf("failed to get previous version: %w", err)
        }
        
        if err := m.executor.SetCurrentVersion(ctx, previousVersion); err != nil {
            return fmt.Errorf("failed to set current version: %w", err)
        }
    }
    
    // 提交事务
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    action := "up"
    if !isUp {
        action = "down"
    }
    
    m.logger.Info("Migration executed successfully",
        logger.String("migration_id", migration.GetID()),
        logger.String("action", action),
        logger.Duration("execution_time", record.ExecutionTime),
    )
    
    return nil
}

// GetPendingMigrations 获取待执行的迁移
func (m *ManagerImpl) GetPendingMigrations(ctx context.Context) ([]Migration, error) {
    currentVersion, err := m.GetCurrentVersion(ctx)
    if err != nil {
        return nil, err
    }
    
    var pending []Migration
    for _, migration := range m.migrations {
        if migration.GetID() > currentVersion {
            pending = append(pending, migration)
        }
    }
    
    // 按ID排序
    sort.Slice(pending, func(i, j int) bool {
        return pending[i].GetID() < pending[j].GetID()
    })
    
    return pending, nil
}

// GetCurrentVersion 获取当前版本
func (m *ManagerImpl) GetCurrentVersion(ctx context.Context) (string, error) {
    return m.executor.GetCurrentVersion(ctx)
}

// GetMigrationHistory 获取迁移历史
func (m *ManagerImpl) GetMigrationHistory(ctx context.Context) ([]MigrationRecord, error) {
    return m.executor.GetMigrationHistory(ctx)
}

// ValidateAll 验证所有迁移
func (m *ManagerImpl) ValidateAll(ctx context.Context) error {
    for _, migration := range m.migrations {
        if err := migration.Validate(ctx, m.executor); err != nil {
            return fmt.Errorf("validation failed for migration %s: %w", migration.GetID(), err)
        }
    }
    
    m.logger.Info("All migrations validated successfully")
    return nil
}

// GenerateMigration 生成迁移文件
func (m *ManagerImpl) GenerateMigration(description string) (string, error) {
    timestamp := time.Now().Format("20060102150405")
    id := fmt.Sprintf("%s_%s", timestamp, strings.ReplaceAll(description, " ", "_"))
    
    filename, err := m.CreateMigrationFile(id, description)
    if err != nil {
        return "", err
    }
    
    m.logger.Info("Migration file generated",
        logger.String("id", id),
        logger.String("filename", filename),
    )
    
    return filename, nil
}

// CreateMigrationFile 创建迁移文件
func (m *ManagerImpl) CreateMigrationFile(id, description string) (string, error) {
    filename := fmt.Sprintf("%s.sql", id)
    filepath := filepath.Join(m.migrationsDir, filename)
    
    content := fmt.Sprintf(`-- %s

-- +migrate Up
-- TODO: Add your migration SQL here


-- +migrate Down
-- TODO: Add your rollback SQL here

`, description)
    
    if err := ioutil.WriteFile(filepath, []byte(content), 0644); err != nil {
        return "", fmt.Errorf("failed to create migration file: %w", err)
    }
    
    return filepath, nil
}
```

## 4.4.4 总结

数据库迁移系统为MovieInfo项目提供了完整的数据库版本管理解决方案。通过自动化的迁移执行、完善的回滚机制和详细的历史记录，我们建立了一个可靠、高效的数据库变更管理系统。

**关键设计要点**：
1. **版本控制**：完整的数据库结构版本管理
2. **自动执行**：自动化的迁移执行和状态跟踪
3. **事务保护**：事务中执行迁移确保数据一致性
4. **回滚支持**：完善的回滚机制和错误恢复
5. **监控日志**：详细的执行日志和性能监控

**迁移优势**：
- **版本一致性**：确保各环境数据库结构一致
- **自动化部署**：支持自动化的数据库部署
- **安全可靠**：事务保护和回滚机制
- **团队协作**：便于团队协作和变更管理

**下一步**：基于完整的数据层基础，我们将开始业务逻辑层的开发，首先实现用户服务的业务逻辑。
