# 4.3 CRUD操作实现

## 4.3.1 概述

CRUD操作（Create、Read、Update、Delete）是数据层的核心功能，它提供了对数据的基本操作接口。对于MovieInfo项目，我们需要实现高效、安全、可维护的CRUD操作，支持事务处理、批量操作和性能优化。

## 4.3.2 为什么需要专业的CRUD实现？

### 4.3.2.1 **数据安全**
- **SQL注入防护**：使用参数化查询防止SQL注入攻击
- **权限控制**：在数据层实现访问权限控制
- **数据验证**：确保数据的完整性和一致性
- **审计日志**：记录数据变更的审计信息

### 4.3.2.2 **性能优化**
- **查询优化**：优化SQL查询语句和执行计划
- **批量操作**：支持批量插入、更新和删除
- **连接复用**：高效的数据库连接管理
- **缓存集成**：与缓存层的无缝集成

### 4.3.2.3 **开发效率**
- **代码复用**：通用的CRUD操作模板
- **类型安全**：强类型的数据操作接口
- **错误处理**：统一的错误处理和异常管理
- **测试支持**：便于单元测试和集成测试

### 4.3.2.4 **业务支持**
- **事务处理**：支持复杂的业务事务
- **软删除**：支持数据的软删除和恢复
- **版本控制**：数据版本的管理和追踪
- **关联查询**：高效的关联数据查询

## 4.3.3 CRUD架构设计

### 4.3.3.1 **Repository模式架构**

```
Repository层架构
├── 接口定义层 (Interface Layer)
│   ├── 基础Repository接口 (BaseRepository)
│   ├── 用户Repository接口 (UserRepository)
│   ├── 电影Repository接口 (MovieRepository)
│   └── 评论Repository接口 (CommentRepository)
├── 实现层 (Implementation Layer)
│   ├── MySQL实现 (MySQLRepository)
│   ├── 缓存装饰器 (CacheDecorator)
│   ├── 事务管理器 (TransactionManager)
│   └── 查询构建器 (QueryBuilder)
├── 工具层 (Utility Layer)
│   ├── SQL构建器 (SQLBuilder)
│   ├── 条件构建器 (ConditionBuilder)
│   ├── 分页工具 (Paginator)
│   └── 排序工具 (Sorter)
└── 监控层 (Monitoring Layer)
    ├── 性能监控 (PerformanceMonitor)
    ├── 查询日志 (QueryLogger)
    ├── 错误统计 (ErrorStats)
    └── 慢查询检测 (SlowQueryDetector)
```

### 4.3.3.2 **基础Repository接口**

#### 4.3.3.2.1 通用Repository接口
```go
// internal/repository/interface.go
package repository

import (
    "context"
    "database/sql"
)

// BaseRepository 基础Repository接口
type BaseRepository[T any] interface {
    // 基础CRUD操作
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id int64) (*T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id int64) error
    
    // 批量操作
    CreateBatch(ctx context.Context, entities []*T) error
    GetByIDs(ctx context.Context, ids []int64) ([]*T, error)
    UpdateBatch(ctx context.Context, entities []*T) error
    DeleteBatch(ctx context.Context, ids []int64) error
    
    // 查询操作
    List(ctx context.Context, opts *ListOptions) ([]*T, error)
    Count(ctx context.Context, opts *CountOptions) (int64, error)
    Exists(ctx context.Context, id int64) (bool, error)
    
    // 事务操作
    WithTx(tx *sql.Tx) BaseRepository[T]
}

// ListOptions 列表查询选项
type ListOptions struct {
    Offset    int                    `json:"offset"`
    Limit     int                    `json:"limit"`
    OrderBy   string                 `json:"order_by"`
    Order     string                 `json:"order"`     // ASC, DESC
    Filters   map[string]interface{} `json:"filters"`
    Includes  []string               `json:"includes"`  // 关联查询
}

// CountOptions 计数查询选项
type CountOptions struct {
    Filters map[string]interface{} `json:"filters"`
}

// PaginationResult 分页结果
type PaginationResult[T any] struct {
    Items      []*T  `json:"items"`
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    TotalPages int   `json:"total_pages"`
    HasNext    bool  `json:"has_next"`
    HasPrev    bool  `json:"has_prev"`
}

// NewPaginationResult 创建分页结果
func NewPaginationResult[T any](items []*T, total int64, page, pageSize int) *PaginationResult[T] {
    totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
    
    return &PaginationResult[T]{
        Items:      items,
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
    }
}
```

#### 4.3.3.2.2 用户Repository接口
```go
// internal/repository/user.go
package repository

import (
    "context"
    
    "github.com/yourname/movieinfo/internal/models"
)

// UserRepository 用户Repository接口
type UserRepository interface {
    BaseRepository[models.User]
    
    // 用户特定查询
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByUsername(ctx context.Context, username string) (*models.User, error)
    GetByEmailOrUsername(ctx context.Context, emailOrUsername string) (*models.User, error)
    
    // 用户状态操作
    UpdateStatus(ctx context.Context, userID int64, status models.UserStatus) error
    UpdateLastLogin(ctx context.Context, userID int64, ip string) error
    VerifyEmail(ctx context.Context, userID int64) error
    
    // 用户统计
    GetActiveUsersCount(ctx context.Context) (int64, error)
    GetUsersByStatus(ctx context.Context, status models.UserStatus, opts *ListOptions) ([]*models.User, error)
    
    // 用户关联查询
    GetUserWithRatings(ctx context.Context, userID int64) (*models.User, error)
    GetUserWithComments(ctx context.Context, userID int64) (*models.User, error)
    GetUserWithFavorites(ctx context.Context, userID int64) (*models.User, error)
}
```

#### 4.3.3.2.3 电影Repository接口
```go
// internal/repository/movie.go
package repository

import (
    "context"
    
    "github.com/yourname/movieinfo/internal/models"
)

// MovieRepository 电影Repository接口
type MovieRepository interface {
    BaseRepository[models.Movie]
    
    // 电影特定查询
    GetByIMDBID(ctx context.Context, imdbID string) (*models.Movie, error)
    GetByTMDBID(ctx context.Context, tmdbID int) (*models.Movie, error)
    SearchByTitle(ctx context.Context, title string, opts *ListOptions) ([]*models.Movie, error)
    
    // 电影分类查询
    GetByCategory(ctx context.Context, categoryID int64, opts *ListOptions) ([]*models.Movie, error)
    GetByCategories(ctx context.Context, categoryIDs []int64, opts *ListOptions) ([]*models.Movie, error)
    
    // 电影排序查询
    GetTopRated(ctx context.Context, limit int) ([]*models.Movie, error)
    GetMostPopular(ctx context.Context, limit int) ([]*models.Movie, error)
    GetLatestReleases(ctx context.Context, limit int) ([]*models.Movie, error)
    GetRecentlyAdded(ctx context.Context, limit int) ([]*models.Movie, error)
    
    // 电影统计操作
    UpdateViewCount(ctx context.Context, movieID int64) error
    UpdateRating(ctx context.Context, movieID int64, rating float64, count int) error
    IncrementCommentCount(ctx context.Context, movieID int64) error
    DecrementCommentCount(ctx context.Context, movieID int64) error
    
    // 电影关联查询
    GetMovieWithCategories(ctx context.Context, movieID int64) (*models.Movie, error)
    GetMovieWithActors(ctx context.Context, movieID int64) (*models.Movie, error)
    GetMovieWithRatings(ctx context.Context, movieID int64) (*models.Movie, error)
    GetMovieWithComments(ctx context.Context, movieID int64, opts *ListOptions) (*models.Movie, error)
    
    // 电影推荐
    GetSimilarMovies(ctx context.Context, movieID int64, limit int) ([]*models.Movie, error)
    GetRecommendedMovies(ctx context.Context, userID int64, limit int) ([]*models.Movie, error)
}
```

### 4.3.3.3 **基础Repository实现**

#### 4.3.3.3.1 MySQL基础实现
```go
// internal/repository/mysql/base.go
package mysql

import (
    "context"
    "database/sql"
    "fmt"
    "reflect"
    "strings"
    
    "github.com/yourname/movieinfo/internal/models"
    "github.com/yourname/movieinfo/internal/repository"
    "github.com/yourname/movieinfo/pkg/database"
    "github.com/yourname/movieinfo/pkg/logger"
)

// BaseRepositoryImpl 基础Repository实现
type BaseRepositoryImpl[T models.Model] struct {
    db        *database.Pool
    tx        *sql.Tx
    tableName string
    logger    logger.Logger
}

// NewBaseRepository 创建基础Repository
func NewBaseRepository[T models.Model](db *database.Pool, tableName string) *BaseRepositoryImpl[T] {
    return &BaseRepositoryImpl[T]{
        db:        db,
        tableName: tableName,
        logger:    logger.GetGlobalLogger(),
    }
}

// Create 创建实体
func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
    // 验证实体
    if err := (*entity).Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // 构建插入SQL
    columns, placeholders, values := r.buildInsertSQL(entity)
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", 
        r.tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))
    
    // 执行插入
    var result sql.Result
    var err error
    
    if r.tx != nil {
        result, err = r.tx.ExecContext(ctx, query, values...)
    } else {
        result, err = r.db.Exec(ctx, query, values...)
    }
    
    if err != nil {
        r.logger.Error("Failed to create entity",
            logger.String("table", r.tableName),
            logger.String("query", query),
            logger.Error(err),
        )
        return fmt.Errorf("failed to create entity: %w", err)
    }
    
    // 获取插入的ID
    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("failed to get last insert id: %w", err)
    }
    
    (*entity).SetID(id)
    
    r.logger.Debug("Entity created successfully",
        logger.String("table", r.tableName),
        logger.Int64("id", id),
    )
    
    return nil
}

// GetByID 根据ID获取实体
func (r *BaseRepositoryImpl[T]) GetByID(ctx context.Context, id int64) (*T, error) {
    query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? AND deleted_at IS NULL", r.tableName)
    
    var row *sql.Row
    if r.tx != nil {
        row = r.tx.QueryRowContext(ctx, query, id)
    } else {
        row = r.db.QueryRow(ctx, query, id)
    }
    
    entity := new(T)
    err := r.scanEntity(row, entity)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, repository.ErrNotFound
        }
        r.logger.Error("Failed to get entity by ID",
            logger.String("table", r.tableName),
            logger.Int64("id", id),
            logger.Error(err),
        )
        return nil, fmt.Errorf("failed to get entity: %w", err)
    }
    
    return entity, nil
}

// Update 更新实体
func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity *T) error {
    // 验证实体
    if err := (*entity).Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // 构建更新SQL
    setClauses, values := r.buildUpdateSQL(entity)
    values = append(values, (*entity).GetID())
    
    query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ? AND deleted_at IS NULL", 
        r.tableName, strings.Join(setClauses, ","))
    
    // 执行更新
    var result sql.Result
    var err error
    
    if r.tx != nil {
        result, err = r.tx.ExecContext(ctx, query, values...)
    } else {
        result, err = r.db.Exec(ctx, query, values...)
    }
    
    if err != nil {
        r.logger.Error("Failed to update entity",
            logger.String("table", r.tableName),
            logger.Int64("id", (*entity).GetID()),
            logger.Error(err),
        )
        return fmt.Errorf("failed to update entity: %w", err)
    }
    
    // 检查是否有行被更新
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return repository.ErrNotFound
    }
    
    r.logger.Debug("Entity updated successfully",
        logger.String("table", r.tableName),
        logger.Int64("id", (*entity).GetID()),
    )
    
    return nil
}

// Delete 删除实体（软删除）
func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, id int64) error {
    query := fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", r.tableName)
    
    var result sql.Result
    var err error
    
    if r.tx != nil {
        result, err = r.tx.ExecContext(ctx, query, id)
    } else {
        result, err = r.db.Exec(ctx, query, id)
    }
    
    if err != nil {
        r.logger.Error("Failed to delete entity",
            logger.String("table", r.tableName),
            logger.Int64("id", id),
            logger.Error(err),
        )
        return fmt.Errorf("failed to delete entity: %w", err)
    }
    
    // 检查是否有行被删除
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return repository.ErrNotFound
    }
    
    r.logger.Debug("Entity deleted successfully",
        logger.String("table", r.tableName),
        logger.Int64("id", id),
    )
    
    return nil
}

// List 列表查询
func (r *BaseRepositoryImpl[T]) List(ctx context.Context, opts *repository.ListOptions) ([]*T, error) {
    query, args := r.buildListSQL(opts)
    
    var rows *sql.Rows
    var err error
    
    if r.tx != nil {
        rows, err = r.tx.QueryContext(ctx, query, args...)
    } else {
        rows, err = r.db.Query(ctx, query, args...)
    }
    
    if err != nil {
        r.logger.Error("Failed to list entities",
            logger.String("table", r.tableName),
            logger.String("query", query),
            logger.Error(err),
        )
        return nil, fmt.Errorf("failed to list entities: %w", err)
    }
    defer rows.Close()
    
    var entities []*T
    for rows.Next() {
        entity := new(T)
        if err := r.scanEntity(rows, entity); err != nil {
            return nil, fmt.Errorf("failed to scan entity: %w", err)
        }
        entities = append(entities, entity)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows iteration error: %w", err)
    }
    
    return entities, nil
}

// Count 计数查询
func (r *BaseRepositoryImpl[T]) Count(ctx context.Context, opts *repository.CountOptions) (int64, error) {
    query, args := r.buildCountSQL(opts)
    
    var count int64
    var err error
    
    if r.tx != nil {
        err = r.tx.QueryRowContext(ctx, query, args...).Scan(&count)
    } else {
        err = r.db.QueryRow(ctx, query, args...).Scan(&count)
    }
    
    if err != nil {
        r.logger.Error("Failed to count entities",
            logger.String("table", r.tableName),
            logger.String("query", query),
            logger.Error(err),
        )
        return 0, fmt.Errorf("failed to count entities: %w", err)
    }
    
    return count, nil
}

// WithTx 使用事务
func (r *BaseRepositoryImpl[T]) WithTx(tx *sql.Tx) repository.BaseRepository[T] {
    return &BaseRepositoryImpl[T]{
        db:        r.db,
        tx:        tx,
        tableName: r.tableName,
        logger:    r.logger,
    }
}

// buildInsertSQL 构建插入SQL
func (r *BaseRepositoryImpl[T]) buildInsertSQL(entity *T) ([]string, []string, []interface{}) {
    // 使用反射获取字段信息
    v := reflect.ValueOf(*entity)
    t := reflect.TypeOf(*entity)
    
    var columns []string
    var placeholders []string
    var values []interface{}
    
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        
        // 跳过ID字段（自增）和时间戳字段（自动设置）
        if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" {
            continue
        }
        
        // 获取数据库字段名
        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" {
            continue
        }
        
        columns = append(columns, dbTag)
        placeholders = append(placeholders, "?")
        values = append(values, value.Interface())
    }
    
    return columns, placeholders, values
}

// buildUpdateSQL 构建更新SQL
func (r *BaseRepositoryImpl[T]) buildUpdateSQL(entity *T) ([]string, []interface{}) {
    // 使用反射获取字段信息
    v := reflect.ValueOf(*entity)
    t := reflect.TypeOf(*entity)
    
    var setClauses []string
    var values []interface{}
    
    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        
        // 跳过ID、创建时间和删除时间字段
        if field.Name == "ID" || field.Name == "CreatedAt" || field.Name == "DeletedAt" {
            continue
        }
        
        // 获取数据库字段名
        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" {
            continue
        }
        
        // 更新时间字段特殊处理
        if field.Name == "UpdatedAt" {
            setClauses = append(setClauses, fmt.Sprintf("%s = NOW()", dbTag))
            continue
        }
        
        setClauses = append(setClauses, fmt.Sprintf("%s = ?", dbTag))
        values = append(values, value.Interface())
    }
    
    return setClauses, values
}
```

## 4.3.4 总结

CRUD操作实现为MovieInfo项目提供了完整的数据访问解决方案。通过Repository模式、泛型支持和完善的错误处理，我们建立了一个类型安全、高性能的数据访问层。

**关键设计要点**：
1. **Repository模式**：清晰的接口定义和实现分离
2. **泛型支持**：类型安全的通用CRUD操作
3. **事务支持**：完整的事务处理机制
4. **性能优化**：批量操作和查询优化
5. **错误处理**：统一的错误处理和日志记录

**CRUD优势**：
- **类型安全**：编译时类型检查
- **代码复用**：通用的CRUD模板
- **性能优化**：高效的SQL生成和执行
- **易于测试**：清晰的接口便于单元测试

**下一步**：基于这个CRUD基础，我们将实现数据库迁移系统，支持数据库结构的版本管理和自动迁移。
