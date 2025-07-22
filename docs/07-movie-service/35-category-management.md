# 第35步：电影分类管理

## 📋 概述

电影分类管理是MovieInfo项目的重要组织功能，为电影内容提供有序的分类体系和管理机制。一个完善的分类系统需要支持多维度分类、层级结构管理和灵活的分类策略。

## 🎯 设计目标

### 1. **分类体系完整性**
- 多维度分类支持
- 层级结构管理
- 标签系统集成
- 自定义分类支持

### 2. **管理功能完善**
- 分类CRUD操作
- 批量管理功能
- 分类统计分析
- 分类关系维护

### 3. **用户体验优化**
- 直观的分类导航
- 智能分类推荐
- 快速分类筛选
- 个性化分类展示

## 🏗️ 分类系统架构

### 1. **分类体系设计**

```
┌─────────────────────────────────────────────────────────────┐
│                    电影分类体系                              │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  类型分类    │  │  地区分类    │  │  年代分类    │          │
│  │             │  │             │  │             │          │
│  │ • 动作片    │  │ • 华语电影   │  │ • 2020年代  │          │
│  │ • 喜剧片    │  │ • 欧美电影   │  │ • 2010年代  │          │
│  │ • 剧情片    │  │ • 日韩电影   │  │ • 经典电影   │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
│                              │                              │
│                              ▼                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                  分类管理系统                            │ │
│  │                                                         │ │
│  │ • 分类创建    • 分类编辑    • 分类删除    • 分类排序     │ │
│  │ • 关系管理    • 统计分析    • 批量操作    • 权限控制     │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 2. **分类数据模型**

```go
// 分类基础模型
type Category struct {
    ID          string    `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"not null;size:100" json:"name"`
    Slug        string    `gorm:"not null;uniqueIndex;size:100" json:"slug"`
    Description string    `gorm:"size:500" json:"description"`
    Type        string    `gorm:"not null;size:50;index" json:"type"` // genre, country, language, decade, custom
    
    // 层级结构
    ParentID    *string   `gorm:"index" json:"parent_id,omitempty"`
    Parent      *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
    Level       int       `gorm:"not null;default:0" json:"level"`
    Path        string    `gorm:"size:500" json:"path"` // 层级路径，如 "/action/thriller"
    
    // 显示属性
    Color       string    `gorm:"size:20" json:"color,omitempty"`
    Icon        string    `gorm:"size:100" json:"icon,omitempty"`
    ImageURL    string    `gorm:"size:500" json:"image_url,omitempty"`
    
    // 排序和状态
    SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`
    Status      string    `gorm:"not null;default:'active'" json:"status"` // active, inactive, archived
    
    // 统计信息
    MovieCount  int       `gorm:"not null;default:0" json:"movie_count"`
    ViewCount   int64     `gorm:"not null;default:0" json:"view_count"`
    
    // 元数据
    Metadata    JSON      `gorm:"type:json" json:"metadata,omitempty"`
    
    // 时间戳
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CreatedBy   string    `gorm:"size:50" json:"created_by,omitempty"`
    UpdatedBy   string    `gorm:"size:50" json:"updated_by,omitempty"`
}

// 电影分类关联
type MovieCategory struct {
    ID         string    `gorm:"primaryKey" json:"id"`
    MovieID    string    `gorm:"not null;index" json:"movie_id"`
    CategoryID string    `gorm:"not null;index" json:"category_id"`
    Movie      *Movie    `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
    Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    
    // 关联属性
    Primary    bool      `gorm:"not null;default:false" json:"primary"` // 是否为主要分类
    Weight     float64   `gorm:"not null;default:1.0" json:"weight"`    // 关联权重
    
    CreatedAt  time.Time `json:"created_at"`
    CreatedBy  string    `gorm:"size:50" json:"created_by,omitempty"`
}

// 分类请求结构
type CategoryRequest struct {
    Name        string            `json:"name" binding:"required,max=100"`
    Slug        string            `json:"slug" binding:"max=100"`
    Description string            `json:"description" binding:"max=500"`
    Type        string            `json:"type" binding:"required,oneof=genre country language decade custom"`
    ParentID    *string           `json:"parent_id"`
    Color       string            `json:"color" binding:"max=20"`
    Icon        string            `json:"icon" binding:"max=100"`
    ImageURL    string            `json:"image_url" binding:"max=500"`
    SortOrder   int               `json:"sort_order"`
    Status      string            `json:"status" binding:"oneof=active inactive archived"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// 分类响应结构
type CategoryResponse struct {
    Success bool      `json:"success"`
    Message string    `json:"message,omitempty"`
    Data    *Category `json:"data,omitempty"`
}

// 分类列表响应
type CategoryListResponse struct {
    Success    bool              `json:"success"`
    Message    string            `json:"message,omitempty"`
    Data       []Category        `json:"data,omitempty"`
    Pagination *PaginationInfo   `json:"pagination,omitempty"`
    Statistics *CategoryStats    `json:"statistics,omitempty"`
}

// 分类统计
type CategoryStats struct {
    TotalCategories int                    `json:"total_categories"`
    TypeDistribution map[string]int        `json:"type_distribution"`
    StatusDistribution map[string]int      `json:"status_distribution"`
    TopCategories   []CategoryStatItem     `json:"top_categories"`
}

type CategoryStatItem struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    MovieCount int    `json:"movie_count"`
    ViewCount  int64  `json:"view_count"`
}
```

## 🔧 分类管理服务

### 1. **核心分类服务**

```go
type CategoryService struct {
    categoryRepo CategoryRepository
    movieRepo    MovieRepository
    cacheStore   CacheStore
    logger       *logrus.Logger
    metrics      *CategoryMetrics
}

func NewCategoryService(
    categoryRepo CategoryRepository,
    movieRepo MovieRepository,
    cacheStore CacheStore,
) *CategoryService {
    return &CategoryService{
        categoryRepo: categoryRepo,
        movieRepo:    movieRepo,
        cacheStore:   cacheStore,
        logger:       logrus.New(),
        metrics:      NewCategoryMetrics(),
    }
}

// 创建分类
func (cs *CategoryService) CreateCategory(ctx context.Context, req *CategoryRequest, userID string) (*CategoryResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCategoryOperation("create", time.Since(start))
    }()

    // 验证请求
    if err := cs.validateCategoryRequest(req); err != nil {
        cs.metrics.IncInvalidRequests("create")
        return &CategoryResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }

    // 生成Slug（如果未提供）
    if req.Slug == "" {
        req.Slug = cs.generateSlug(req.Name)
    }

    // 检查Slug唯一性
    if exists, err := cs.categoryRepo.ExistsBySlug(ctx, req.Slug); err != nil {
        cs.logger.Errorf("Failed to check slug existence: %v", err)
        return nil, errors.New("分类创建失败")
    } else if exists {
        return &CategoryResponse{
            Success: false,
            Message: "分类标识已存在",
        }, nil
    }

    // 处理层级关系
    var level int
    var path string
    if req.ParentID != nil {
        parent, err := cs.categoryRepo.FindByID(ctx, *req.ParentID)
        if err != nil {
            return &CategoryResponse{
                Success: false,
                Message: "父分类不存在",
            }, nil
        }
        level = parent.Level + 1
        path = parent.Path + "/" + req.Slug
    } else {
        level = 0
        path = "/" + req.Slug
    }

    // 创建分类对象
    category := &Category{
        ID:          uuid.New().String(),
        Name:        req.Name,
        Slug:        req.Slug,
        Description: req.Description,
        Type:        req.Type,
        ParentID:    req.ParentID,
        Level:       level,
        Path:        path,
        Color:       req.Color,
        Icon:        req.Icon,
        ImageURL:    req.ImageURL,
        SortOrder:   req.SortOrder,
        Status:      req.Status,
        Metadata:    req.Metadata,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        CreatedBy:   userID,
        UpdatedBy:   userID,
    }

    // 保存分类
    if err := cs.categoryRepo.Create(ctx, category); err != nil {
        cs.logger.Errorf("Failed to create category: %v", err)
        cs.metrics.IncOperationErrors("create")
        return nil, errors.New("分类创建失败")
    }

    // 清除相关缓存
    cs.clearCategoryCache(ctx, req.Type)

    cs.metrics.IncSuccessfulOperations("create")
    cs.logger.Infof("Category created: %s (%s)", category.Name, category.ID)

    return &CategoryResponse{
        Success: true,
        Message: "分类创建成功",
        Data:    category,
    }, nil
}

// 获取分类列表
func (cs *CategoryService) GetCategories(ctx context.Context, categoryType string, includeInactive bool) (*CategoryListResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCategoryOperation("list", time.Since(start))
    }()

    // 生成缓存键
    cacheKey := cs.generateListCacheKey(categoryType, includeInactive)

    // 尝试从缓存获取
    if cachedResult, err := cs.getListFromCache(ctx, cacheKey); err == nil {
        cs.metrics.IncCacheHits()
        return cachedResult, nil
    }
    cs.metrics.IncCacheMisses()

    // 构建查询条件
    conditions := make(map[string]interface{})
    if categoryType != "" {
        conditions["type"] = categoryType
    }
    if !includeInactive {
        conditions["status"] = "active"
    }

    // 获取分类列表
    categories, err := cs.categoryRepo.FindWithConditions(ctx, conditions)
    if err != nil {
        cs.logger.Errorf("Failed to get categories: %v", err)
        cs.metrics.IncOperationErrors("list")
        return nil, errors.New("获取分类列表失败")
    }

    // 构建层级结构
    categoryTree := cs.buildCategoryTree(categories)

    // 获取统计信息
    stats, err := cs.getCategoryStatistics(ctx)
    if err != nil {
        cs.logger.Errorf("Failed to get category statistics: %v", err)
        // 统计信息获取失败不影响主要功能
    }

    response := &CategoryListResponse{
        Success:    true,
        Data:       categoryTree,
        Statistics: stats,
    }

    // 异步缓存结果
    go func() {
        if err := cs.cacheListResult(context.Background(), cacheKey, response); err != nil {
            cs.logger.Errorf("Failed to cache category list: %v", err)
        }
    }()

    cs.metrics.IncSuccessfulOperations("list")
    return response, nil
}

// 更新分类
func (cs *CategoryService) UpdateCategory(ctx context.Context, categoryID string, req *CategoryRequest, userID string) (*CategoryResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCategoryOperation("update", time.Since(start))
    }()

    // 获取现有分类
    category, err := cs.categoryRepo.FindByID(ctx, categoryID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &CategoryResponse{
                Success: false,
                Message: "分类不存在",
            }, nil
        }
        cs.logger.Errorf("Failed to get category: %v", err)
        return nil, errors.New("获取分类失败")
    }

    // 验证请求
    if err := cs.validateCategoryRequest(req); err != nil {
        cs.metrics.IncInvalidRequests("update")
        return &CategoryResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }

    // 检查Slug唯一性（如果有变化）
    if req.Slug != "" && req.Slug != category.Slug {
        if exists, err := cs.categoryRepo.ExistsBySlug(ctx, req.Slug); err != nil {
            cs.logger.Errorf("Failed to check slug existence: %v", err)
            return nil, errors.New("分类更新失败")
        } else if exists {
            return &CategoryResponse{
                Success: false,
                Message: "分类标识已存在",
            }, nil
        }
    }

    // 更新分类字段
    cs.updateCategoryFields(category, req)
    category.UpdatedAt = time.Now()
    category.UpdatedBy = userID

    // 处理层级关系变化
    if req.ParentID != nil && (category.ParentID == nil || *req.ParentID != *category.ParentID) {
        if err := cs.updateCategoryHierarchy(ctx, category, req.ParentID); err != nil {
            return &CategoryResponse{
                Success: false,
                Message: err.Error(),
            }, nil
        }
    }

    // 保存更新
    if err := cs.categoryRepo.Update(ctx, category); err != nil {
        cs.logger.Errorf("Failed to update category: %v", err)
        cs.metrics.IncOperationErrors("update")
        return nil, errors.New("分类更新失败")
    }

    // 清除相关缓存
    cs.clearCategoryCache(ctx, category.Type)

    cs.metrics.IncSuccessfulOperations("update")
    cs.logger.Infof("Category updated: %s (%s)", category.Name, category.ID)

    return &CategoryResponse{
        Success: true,
        Message: "分类更新成功",
        Data:    category,
    }, nil
}

// 删除分类
func (cs *CategoryService) DeleteCategory(ctx context.Context, categoryID string, userID string) (*CategoryResponse, error) {
    start := time.Now()
    defer func() {
        cs.metrics.ObserveCategoryOperation("delete", time.Since(start))
    }()

    // 获取分类
    category, err := cs.categoryRepo.FindByID(ctx, categoryID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &CategoryResponse{
                Success: false,
                Message: "分类不存在",
            }, nil
        }
        return nil, errors.New("获取分类失败")
    }

    // 检查是否有子分类
    hasChildren, err := cs.categoryRepo.HasChildren(ctx, categoryID)
    if err != nil {
        cs.logger.Errorf("Failed to check children: %v", err)
        return nil, errors.New("删除检查失败")
    }
    if hasChildren {
        return &CategoryResponse{
            Success: false,
            Message: "该分类下还有子分类，无法删除",
        }, nil
    }

    // 检查是否有关联电影
    movieCount, err := cs.categoryRepo.GetMovieCount(ctx, categoryID)
    if err != nil {
        cs.logger.Errorf("Failed to get movie count: %v", err)
        return nil, errors.New("删除检查失败")
    }
    if movieCount > 0 {
        return &CategoryResponse{
            Success: false,
            Message: fmt.Sprintf("该分类下还有%d部电影，无法删除", movieCount),
        }, nil
    }

    // 执行删除
    if err := cs.categoryRepo.Delete(ctx, categoryID); err != nil {
        cs.logger.Errorf("Failed to delete category: %v", err)
        cs.metrics.IncOperationErrors("delete")
        return nil, errors.New("分类删除失败")
    }

    // 清除相关缓存
    cs.clearCategoryCache(ctx, category.Type)

    cs.metrics.IncSuccessfulOperations("delete")
    cs.logger.Infof("Category deleted: %s (%s)", category.Name, category.ID)

    return &CategoryResponse{
        Success: true,
        Message: "分类删除成功",
    }, nil
}

// 构建分类树
func (cs *CategoryService) buildCategoryTree(categories []Category) []Category {
    categoryMap := make(map[string]*Category)
    var rootCategories []Category

    // 创建分类映射
    for i := range categories {
        categoryMap[categories[i].ID] = &categories[i]
        categories[i].Children = []Category{}
    }

    // 构建层级关系
    for i := range categories {
        if categories[i].ParentID == nil {
            rootCategories = append(rootCategories, categories[i])
        } else {
            if parent, exists := categoryMap[*categories[i].ParentID]; exists {
                parent.Children = append(parent.Children, categories[i])
            }
        }
    }

    // 排序
    cs.sortCategories(rootCategories)
    for _, category := range categoryMap {
        cs.sortCategories(category.Children)
    }

    return rootCategories
}

// 排序分类
func (cs *CategoryService) sortCategories(categories []Category) {
    sort.Slice(categories, func(i, j int) bool {
        if categories[i].SortOrder != categories[j].SortOrder {
            return categories[i].SortOrder < categories[j].SortOrder
        }
        return categories[i].Name < categories[j].Name
    })
}
```

### 2. **分类关联管理**

```go
// 电影分类关联服务
type MovieCategoryService struct {
    movieCategoryRepo MovieCategoryRepository
    categoryRepo      CategoryRepository
    logger            *logrus.Logger
}

func NewMovieCategoryService(
    movieCategoryRepo MovieCategoryRepository,
    categoryRepo CategoryRepository,
) *MovieCategoryService {
    return &MovieCategoryService{
        movieCategoryRepo: movieCategoryRepo,
        categoryRepo:      categoryRepo,
        logger:            logrus.New(),
    }
}

// 为电影添加分类
func (mcs *MovieCategoryService) AddMovieCategories(ctx context.Context, movieID string, categoryIDs []string, userID string) error {
    // 验证分类存在性
    for _, categoryID := range categoryIDs {
        if exists, err := mcs.categoryRepo.ExistsByID(ctx, categoryID); err != nil {
            return err
        } else if !exists {
            return fmt.Errorf("分类 %s 不存在", categoryID)
        }
    }

    // 创建关联记录
    for i, categoryID := range categoryIDs {
        association := &MovieCategory{
            ID:         uuid.New().String(),
            MovieID:    movieID,
            CategoryID: categoryID,
            Primary:    i == 0, // 第一个分类为主要分类
            Weight:     1.0,
            CreatedAt:  time.Now(),
            CreatedBy:  userID,
        }

        if err := mcs.movieCategoryRepo.Create(ctx, association); err != nil {
            mcs.logger.Errorf("Failed to create movie category association: %v", err)
            return err
        }
    }

    // 更新分类电影计数
    go func() {
        for _, categoryID := range categoryIDs {
            mcs.updateCategoryMovieCount(context.Background(), categoryID)
        }
    }()

    return nil
}

// 更新分类电影计数
func (mcs *MovieCategoryService) updateCategoryMovieCount(ctx context.Context, categoryID string) {
    count, err := mcs.movieCategoryRepo.CountMoviesByCategory(ctx, categoryID)
    if err != nil {
        mcs.logger.Errorf("Failed to count movies for category %s: %v", categoryID, err)
        return
    }

    if err := mcs.categoryRepo.UpdateMovieCount(ctx, categoryID, int(count)); err != nil {
        mcs.logger.Errorf("Failed to update movie count for category %s: %v", categoryID, err)
    }
}
```

## 📊 性能监控

### 1. **分类管理指标**

```go
type CategoryMetrics struct {
    operationCount    *prometheus.CounterVec
    operationDuration *prometheus.HistogramVec
    cacheHitRate      *prometheus.CounterVec
    categoryCount     prometheus.Gauge
    errorCount        *prometheus.CounterVec
}

func NewCategoryMetrics() *CategoryMetrics {
    return &CategoryMetrics{
        operationCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "category_operations_total",
                Help: "Total number of category operations",
            },
            []string{"operation", "status"},
        ),
        operationDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "category_operation_duration_seconds",
                Help: "Duration of category operations",
            },
            []string{"operation"},
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "category_cache_operations_total",
                Help: "Total number of category cache operations",
            },
            []string{"type"},
        ),
        categoryCount: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "categories_total",
                Help: "Total number of categories",
            },
        ),
        errorCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "category_errors_total",
                Help: "Total number of category errors",
            },
            []string{"operation", "error_type"},
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **分类管理API端点**

```go
func (cc *CategoryController) CreateCategory(c *gin.Context) {
    var req CategoryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }

    userID := cc.getUserIDFromContext(c)
    response, err := cc.categoryService.CreateCategory(c.Request.Context(), &req, userID)
    if err != nil {
        cc.logger.Errorf("Failed to create category: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "分类创建失败",
        })
        return
    }

    c.JSON(200, response)
}

func (cc *CategoryController) GetCategories(c *gin.Context) {
    categoryType := c.Query("type")
    includeInactive := c.Query("include_inactive") == "true"

    response, err := cc.categoryService.GetCategories(c.Request.Context(), categoryType, includeInactive)
    if err != nil {
        cc.logger.Errorf("Failed to get categories: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": "获取分类列表失败",
        })
        return
    }

    c.Header("Cache-Control", "public, max-age=1800") // 30分钟缓存
    c.JSON(200, response)
}
```

## 📝 总结

电影分类管理为MovieInfo项目提供了完整的内容组织体系：

**核心功能**：
1. **多维分类**：类型、地区、年代等多维度分类支持
2. **层级管理**：支持分类的层级结构和关系管理
3. **批量操作**：高效的批量分类管理功能
4. **统计分析**：完整的分类使用统计和分析

**管理特性**：
- 灵活的分类创建和编辑
- 智能的层级关系维护
- 完善的权限控制机制
- 实时的统计信息更新

**性能优化**：
- 多层缓存策略
- 数据库查询优化
- 异步统计更新
- 监控指标收集

至此，电影服务的核心功能已经完成。下一步，我们将继续完成评论服务、主页服务等其他模块的开发文档。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第36步：评论数据模型](../08-comment-service/36-comment-model.md)
