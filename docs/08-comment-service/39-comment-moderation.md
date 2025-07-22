# 第39步：评论审核机制

## 📋 概述

评论审核机制是保障平台内容质量和用户体验的关键系统。通过自动化审核和人工审核相结合的方式，确保平台内容的健康性、合规性和高质量。

## 🎯 设计目标

### 1. **内容安全**
- 敏感词检测过滤
- 违规内容识别
- 恶意行为防范
- 法律合规保障

### 2. **审核效率**
- 自动化审核优先
- 智能分级处理
- 快速响应机制
- 批量处理能力

### 3. **用户体验**
- 透明的审核流程
- 及时的反馈通知
- 申诉处理机制
- 公平的审核标准

## 🔧 审核系统架构

### 1. **审核数据模型**

```go
// 审核规则配置
type ModerationRule struct {
    ID          string    `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"not null;size:100" json:"name"`
    Type        string    `gorm:"not null;size:50" json:"type"` // keyword, pattern, ml_model, manual
    Category    string    `gorm:"not null;size:50" json:"category"` // spam, abuse, adult, political
    
    // 规则内容
    Keywords    []string  `gorm:"type:json" json:"keywords,omitempty"`
    Patterns    []string  `gorm:"type:json" json:"patterns,omitempty"`
    ModelConfig JSON      `gorm:"type:json" json:"model_config,omitempty"`
    
    // 处理动作
    Action      string    `gorm:"not null;size:50" json:"action"` // block, flag, review, warn
    Severity    int       `gorm:"not null;default:1" json:"severity"` // 1-5
    
    // 状态信息
    Enabled     bool      `gorm:"not null;default:true" json:"enabled"`
    Priority    int       `gorm:"not null;default:0" json:"priority"`
    
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CreatedBy   string    `gorm:"size:50" json:"created_by"`
}

// 审核记录
type ModerationRecord struct {
    ID          string    `gorm:"primaryKey" json:"id"`
    ContentID   string    `gorm:"not null;index" json:"content_id"`
    ContentType string    `gorm:"not null;size:50" json:"content_type"` // comment, rating, review
    UserID      string    `gorm:"not null;index" json:"user_id"`
    
    // 审核结果
    Status      string    `gorm:"not null;size:50" json:"status"` // pending, approved, rejected, flagged
    Action      string    `gorm:"not null;size:50" json:"action"` // auto_approve, auto_reject, manual_review
    Confidence  float64   `gorm:"not null;default:0" json:"confidence"`
    
    // 审核详情
    RuleMatches []RuleMatch `gorm:"type:json" json:"rule_matches"`
    Reason      string    `gorm:"size:500" json:"reason"`
    ReviewerID  *string   `gorm:"size:50" json:"reviewer_id,omitempty"`
    ReviewNote  string    `gorm:"size:1000" json:"review_note,omitempty"`
    
    // 时间信息
    CreatedAt   time.Time `json:"created_at"`
    ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
    
    // 关联数据
    User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Reviewer    *User     `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// 规则匹配结果
type RuleMatch struct {
    RuleID      string  `json:"rule_id"`
    RuleName    string  `json:"rule_name"`
    Category    string  `json:"category"`
    Severity    int     `json:"severity"`
    Confidence  float64 `json:"confidence"`
    MatchedText string  `json:"matched_text,omitempty"`
    Position    int     `json:"position,omitempty"`
}

// 审核请求
type ModerationRequest struct {
    ContentID   string `json:"content_id" binding:"required"`
    ContentType string `json:"content_type" binding:"required"`
    Content     string `json:"content" binding:"required"`
    UserID      string `json:"user_id" binding:"required"`
    Context     map[string]interface{} `json:"context,omitempty"`
}

// 审核响应
type ModerationResponse struct {
    Success    bool           `json:"success"`
    Status     string         `json:"status"`
    Action     string         `json:"action"`
    Confidence float64        `json:"confidence"`
    Reason     string         `json:"reason,omitempty"`
    Matches    []RuleMatch    `json:"matches,omitempty"`
    RecordID   string         `json:"record_id,omitempty"`
}
```

### 2. **审核服务实现**

```go
type ModerationService struct {
    ruleRepo       ModerationRuleRepository
    recordRepo     ModerationRecordRepository
    keywordFilter  *KeywordFilter
    patternMatcher *PatternMatcher
    mlService      MLModerationService
    notificationSvc NotificationService
    logger         *logrus.Logger
    metrics        *ModerationMetrics
}

func NewModerationService(
    ruleRepo ModerationRuleRepository,
    recordRepo ModerationRecordRepository,
    keywordFilter *KeywordFilter,
    patternMatcher *PatternMatcher,
    mlService MLModerationService,
    notificationSvc NotificationService,
) *ModerationService {
    return &ModerationService{
        ruleRepo:        ruleRepo,
        recordRepo:      recordRepo,
        keywordFilter:   keywordFilter,
        patternMatcher:  patternMatcher,
        mlService:       mlService,
        notificationSvc: notificationSvc,
        logger:          logrus.New(),
        metrics:         NewModerationMetrics(),
    }
}

// 审核内容
func (ms *ModerationService) ModerateContent(ctx context.Context, req *ModerationRequest) (*ModerationResponse, error) {
    start := time.Now()
    defer func() {
        ms.metrics.ObserveModerationDuration(time.Since(start))
    }()

    // 获取适用的审核规则
    rules, err := ms.ruleRepo.FindEnabledRules(ctx)
    if err != nil {
        ms.logger.Errorf("Failed to get moderation rules: %v", err)
        return nil, errors.New("获取审核规则失败")
    }

    // 执行多层审核
    var allMatches []RuleMatch
    var maxSeverity int
    var totalConfidence float64

    // 1. 关键词过滤
    keywordMatches := ms.checkKeywords(req.Content, rules)
    allMatches = append(allMatches, keywordMatches...)

    // 2. 正则表达式匹配
    patternMatches := ms.checkPatterns(req.Content, rules)
    allMatches = append(allMatches, patternMatches...)

    // 3. 机器学习模型检测
    mlMatches, err := ms.checkMLModels(ctx, req.Content, rules)
    if err != nil {
        ms.logger.Errorf("ML moderation failed: %v", err)
        // ML检测失败不影响其他检测
    } else {
        allMatches = append(allMatches, mlMatches...)
    }

    // 4. 用户行为分析
    behaviorMatches := ms.checkUserBehavior(ctx, req.UserID, req.Context)
    allMatches = append(allMatches, behaviorMatches...)

    // 计算综合评分
    for _, match := range allMatches {
        if match.Severity > maxSeverity {
            maxSeverity = match.Severity
        }
        totalConfidence += match.Confidence
    }

    // 确定审核结果
    status, action := ms.determineResult(allMatches, maxSeverity, totalConfidence)

    // 创建审核记录
    record := &ModerationRecord{
        ID:          uuid.New().String(),
        ContentID:   req.ContentID,
        ContentType: req.ContentType,
        UserID:      req.UserID,
        Status:      status,
        Action:      action,
        Confidence:  totalConfidence / float64(len(allMatches)+1),
        RuleMatches: allMatches,
        Reason:      ms.generateReason(allMatches),
        CreatedAt:   time.Now(),
    }

    if err := ms.recordRepo.Create(ctx, record); err != nil {
        ms.logger.Errorf("Failed to create moderation record: %v", err)
        // 记录创建失败不影响审核结果
    }

    // 发送通知（如果需要人工审核）
    if action == "manual_review" {
        go ms.notifyReviewers(context.Background(), record)
    }

    // 更新指标
    ms.metrics.IncModerationCount(status, action)
    if len(allMatches) > 0 {
        ms.metrics.IncViolationCount(allMatches[0].Category)
    }

    response := &ModerationResponse{
        Success:    true,
        Status:     status,
        Action:     action,
        Confidence: record.Confidence,
        Reason:     record.Reason,
        Matches:    allMatches,
        RecordID:   record.ID,
    }

    ms.logger.Infof("Content moderated: %s, status: %s, action: %s, matches: %d", 
        req.ContentID, status, action, len(allMatches))

    return response, nil
}

// 关键词检测
func (ms *ModerationService) checkKeywords(content string, rules []ModerationRule) []RuleMatch {
    var matches []RuleMatch

    for _, rule := range rules {
        if rule.Type != "keyword" || !rule.Enabled {
            continue
        }

        for _, keyword := range rule.Keywords {
            if ms.keywordFilter.Contains(content, keyword) {
                matches = append(matches, RuleMatch{
                    RuleID:      rule.ID,
                    RuleName:    rule.Name,
                    Category:    rule.Category,
                    Severity:    rule.Severity,
                    Confidence:  0.9, // 关键词匹配置信度较高
                    MatchedText: keyword,
                    Position:    ms.keywordFilter.FindPosition(content, keyword),
                })
            }
        }
    }

    return matches
}

// 正则表达式检测
func (ms *ModerationService) checkPatterns(content string, rules []ModerationRule) []RuleMatch {
    var matches []RuleMatch

    for _, rule := range rules {
        if rule.Type != "pattern" || !rule.Enabled {
            continue
        }

        for _, pattern := range rule.Patterns {
            if matched, matchedText, position := ms.patternMatcher.Match(content, pattern); matched {
                matches = append(matches, RuleMatch{
                    RuleID:      rule.ID,
                    RuleName:    rule.Name,
                    Category:    rule.Category,
                    Severity:    rule.Severity,
                    Confidence:  0.8, // 正则匹配置信度中等
                    MatchedText: matchedText,
                    Position:    position,
                })
            }
        }
    }

    return matches
}

// 机器学习模型检测
func (ms *ModerationService) checkMLModels(ctx context.Context, content string, rules []ModerationRule) ([]RuleMatch, error) {
    var matches []RuleMatch

    for _, rule := range rules {
        if rule.Type != "ml_model" || !rule.Enabled {
            continue
        }

        result, err := ms.mlService.Predict(ctx, content, rule.ModelConfig)
        if err != nil {
            ms.logger.Errorf("ML model prediction failed for rule %s: %v", rule.ID, err)
            continue
        }

        if result.IsViolation {
            matches = append(matches, RuleMatch{
                RuleID:     rule.ID,
                RuleName:   rule.Name,
                Category:   rule.Category,
                Severity:   rule.Severity,
                Confidence: result.Confidence,
            })
        }
    }

    return matches, nil
}

// 用户行为分析
func (ms *ModerationService) checkUserBehavior(ctx context.Context, userID string, context map[string]interface{}) []RuleMatch {
    var matches []RuleMatch

    // 检查用户历史违规记录
    violationCount, err := ms.recordRepo.CountUserViolations(ctx, userID, 30*24*time.Hour) // 30天内
    if err != nil {
        ms.logger.Errorf("Failed to get user violation count: %v", err)
        return matches
    }

    if violationCount >= 5 {
        matches = append(matches, RuleMatch{
            RuleID:     "behavior_frequent_violations",
            RuleName:   "频繁违规行为",
            Category:   "behavior",
            Severity:   3,
            Confidence: 0.7,
        })
    }

    // 检查发布频率
    if context != nil {
        if recentCount, ok := context["recent_posts"].(int); ok && recentCount > 10 {
            matches = append(matches, RuleMatch{
                RuleID:     "behavior_spam_posting",
                RuleName:   "疑似刷屏行为",
                Category:   "spam",
                Severity:   2,
                Confidence: 0.6,
            })
        }
    }

    return matches
}

// 确定审核结果
func (ms *ModerationService) determineResult(matches []RuleMatch, maxSeverity int, totalConfidence float64) (string, string) {
    if len(matches) == 0 {
        return "approved", "auto_approve"
    }

    avgConfidence := totalConfidence / float64(len(matches))

    // 高严重性或高置信度直接拒绝
    if maxSeverity >= 4 || avgConfidence >= 0.9 {
        return "rejected", "auto_reject"
    }

    // 中等严重性需要人工审核
    if maxSeverity >= 2 || avgConfidence >= 0.6 {
        return "pending", "manual_review"
    }

    // 低严重性标记但通过
    return "flagged", "auto_approve"
}

// 生成审核原因
func (ms *ModerationService) generateReason(matches []RuleMatch) string {
    if len(matches) == 0 {
        return "内容审核通过"
    }

    categories := make(map[string]bool)
    for _, match := range matches {
        categories[match.Category] = true
    }

    var reasons []string
    for category := range categories {
        switch category {
        case "spam":
            reasons = append(reasons, "疑似垃圾内容")
        case "abuse":
            reasons = append(reasons, "包含辱骂内容")
        case "adult":
            reasons = append(reasons, "包含成人内容")
        case "political":
            reasons = append(reasons, "包含敏感政治内容")
        case "behavior":
            reasons = append(reasons, "用户行为异常")
        }
    }

    if len(reasons) == 0 {
        return "触发审核规则"
    }

    return strings.Join(reasons, "、")
}
```

### 3. **关键词过滤器**

```go
type KeywordFilter struct {
    trie     *Trie
    keywords map[string]bool
    mutex    sync.RWMutex
}

func NewKeywordFilter() *KeywordFilter {
    return &KeywordFilter{
        trie:     NewTrie(),
        keywords: make(map[string]bool),
    }
}

// 加载关键词
func (kf *KeywordFilter) LoadKeywords(keywords []string) {
    kf.mutex.Lock()
    defer kf.mutex.Unlock()

    kf.trie = NewTrie()
    kf.keywords = make(map[string]bool)

    for _, keyword := range keywords {
        normalized := kf.normalizeText(keyword)
        kf.trie.Insert(normalized)
        kf.keywords[normalized] = true
    }
}

// 检查是否包含敏感词
func (kf *KeywordFilter) Contains(text, keyword string) bool {
    kf.mutex.RLock()
    defer kf.mutex.RUnlock()

    normalizedText := kf.normalizeText(text)
    normalizedKeyword := kf.normalizeText(keyword)

    return strings.Contains(normalizedText, normalizedKeyword)
}

// 查找敏感词位置
func (kf *KeywordFilter) FindPosition(text, keyword string) int {
    normalizedText := kf.normalizeText(text)
    normalizedKeyword := kf.normalizeText(keyword)

    return strings.Index(normalizedText, normalizedKeyword)
}

// 文本标准化
func (kf *KeywordFilter) normalizeText(text string) string {
    // 转小写
    text = strings.ToLower(text)
    
    // 移除空格和特殊字符
    text = regexp.MustCompile(`\s+`).ReplaceAllString(text, "")
    text = regexp.MustCompile(`[^\w\u4e00-\u9fff]`).ReplaceAllString(text, "")
    
    return text
}

// Trie树实现
type Trie struct {
    root *TrieNode
}

type TrieNode struct {
    children map[rune]*TrieNode
    isEnd    bool
}

func NewTrie() *Trie {
    return &Trie{
        root: &TrieNode{
            children: make(map[rune]*TrieNode),
            isEnd:    false,
        },
    }
}

func (t *Trie) Insert(word string) {
    node := t.root
    for _, char := range word {
        if _, exists := node.children[char]; !exists {
            node.children[char] = &TrieNode{
                children: make(map[rune]*TrieNode),
                isEnd:    false,
            }
        }
        node = node.children[char]
    }
    node.isEnd = true
}

func (t *Trie) Search(word string) bool {
    node := t.root
    for _, char := range word {
        if _, exists := node.children[char]; !exists {
            return false
        }
        node = node.children[char]
    }
    return node.isEnd
}
```

### 4. **人工审核管理**

```go
type ManualReviewService struct {
    recordRepo      ModerationRecordRepository
    userRepo        UserRepository
    notificationSvc NotificationService
    logger          *logrus.Logger
}

func NewManualReviewService(
    recordRepo ModerationRecordRepository,
    userRepo UserRepository,
    notificationSvc NotificationService,
) *ManualReviewService {
    return &ManualReviewService{
        recordRepo:      recordRepo,
        userRepo:        userRepo,
        notificationSvc: notificationSvc,
        logger:          logrus.New(),
    }
}

// 获取待审核内容
func (mrs *ManualReviewService) GetPendingReviews(ctx context.Context, reviewerID string, page, pageSize int) (*ReviewListResponse, error) {
    // 验证审核员权限
    reviewer, err := mrs.userRepo.FindByID(ctx, reviewerID)
    if err != nil {
        return nil, errors.New("审核员不存在")
    }

    if reviewer.Role != "admin" && reviewer.Role != "moderator" {
        return nil, errors.New("无审核权限")
    }

    // 获取待审核记录
    records, total, err := mrs.recordRepo.FindPendingReviews(ctx, page, pageSize)
    if err != nil {
        mrs.logger.Errorf("Failed to get pending reviews: %v", err)
        return nil, errors.New("获取待审核内容失败")
    }

    // 构建响应
    items := make([]ReviewItem, len(records))
    for i, record := range records {
        items[i] = ReviewItem{
            ID:          record.ID,
            ContentID:   record.ContentID,
            ContentType: record.ContentType,
            UserID:      record.UserID,
            Status:      record.Status,
            Confidence:  record.Confidence,
            RuleMatches: record.RuleMatches,
            Reason:      record.Reason,
            CreatedAt:   record.CreatedAt,
        }
    }

    return &ReviewListResponse{
        Success: true,
        Data:    items,
        Pagination: &PaginationInfo{
            CurrentPage: page,
            PageSize:    pageSize,
            TotalItems:  total,
            TotalPages:  int(math.Ceil(float64(total) / float64(pageSize))),
        },
    }, nil
}

// 审核决定
func (mrs *ManualReviewService) MakeReviewDecision(ctx context.Context, recordID, reviewerID, decision, note string) error {
    // 获取审核记录
    record, err := mrs.recordRepo.FindByID(ctx, recordID)
    if err != nil {
        return errors.New("审核记录不存在")
    }

    if record.Status != "pending" {
        return errors.New("该内容已被审核")
    }

    // 验证审核员权限
    reviewer, err := mrs.userRepo.FindByID(ctx, reviewerID)
    if err != nil {
        return errors.New("审核员不存在")
    }

    if reviewer.Role != "admin" && reviewer.Role != "moderator" {
        return errors.New("无审核权限")
    }

    // 更新审核记录
    now := time.Now()
    record.Status = decision
    record.ReviewerID = &reviewerID
    record.ReviewNote = note
    record.ReviewedAt = &now

    if err := mrs.recordRepo.Update(ctx, record); err != nil {
        mrs.logger.Errorf("Failed to update review record: %v", err)
        return errors.New("审核决定保存失败")
    }

    // 发送通知给内容作者
    go mrs.notifyUser(context.Background(), record, decision)

    mrs.logger.Infof("Manual review completed: record %s, decision %s, reviewer %s", 
        recordID, decision, reviewerID)

    return nil
}

// 通知用户审核结果
func (mrs *ManualReviewService) notifyUser(ctx context.Context, record *ModerationRecord, decision string) {
    var message string
    switch decision {
    case "approved":
        message = "您的内容已通过审核"
    case "rejected":
        message = "您的内容未通过审核：" + record.Reason
    default:
        return
    }

    notification := &Notification{
        UserID:  record.UserID,
        Type:    "moderation_result",
        Title:   "内容审核结果",
        Content: message,
        Data: map[string]interface{}{
            "content_id":   record.ContentID,
            "content_type": record.ContentType,
            "decision":     decision,
            "reason":       record.Reason,
        },
    }

    if err := mrs.notificationSvc.Send(ctx, notification); err != nil {
        mrs.logger.Errorf("Failed to send moderation notification: %v", err)
    }
}
```

## 📊 监控指标

### 1. **审核系统指标**

```go
type ModerationMetrics struct {
    moderationCount    *prometheus.CounterVec
    moderationDuration prometheus.Histogram
    violationCount     *prometheus.CounterVec
    reviewQueueSize    prometheus.Gauge
    falsePositiveRate  prometheus.Gauge
}

func NewModerationMetrics() *ModerationMetrics {
    return &ModerationMetrics{
        moderationCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "moderation_operations_total",
                Help: "Total number of moderation operations",
            },
            []string{"status", "action"},
        ),
        moderationDuration: prometheus.NewHistogram(
            prometheus.HistogramOpts{
                Name: "moderation_duration_seconds",
                Help: "Duration of moderation operations",
            },
        ),
        violationCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "content_violations_total",
                Help: "Total number of content violations",
            },
            []string{"category"},
        ),
        reviewQueueSize: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "manual_review_queue_size",
                Help: "Number of items in manual review queue",
            },
        ),
        falsePositiveRate: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "moderation_false_positive_rate",
                Help: "False positive rate of automatic moderation",
            },
        ),
    }
}
```

## 🔧 HTTP处理器

### 1. **审核管理API**

```go
func (mc *ModerationController) GetPendingReviews(c *gin.Context) {
    reviewerID := mc.getUserIDFromContext(c)
    if reviewerID == "" {
        c.JSON(401, gin.H{
            "success": false,
            "message": "请先登录",
        })
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

    response, err := mc.manualReviewService.GetPendingReviews(c.Request.Context(), reviewerID, page, pageSize)
    if err != nil {
        mc.logger.Errorf("Failed to get pending reviews: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, response)
}

func (mc *ModerationController) MakeReviewDecision(c *gin.Context) {
    var req struct {
        Decision string `json:"decision" binding:"required,oneof=approved rejected"`
        Note     string `json:"note" binding:"max=1000"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{
            "success": false,
            "message": "请求参数错误",
        })
        return
    }

    recordID := c.Param("record_id")
    reviewerID := mc.getUserIDFromContext(c)

    if err := mc.manualReviewService.MakeReviewDecision(c.Request.Context(), recordID, reviewerID, req.Decision, req.Note); err != nil {
        mc.logger.Errorf("Failed to make review decision: %v", err)
        c.JSON(500, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "审核决定已保存",
    })
}
```

## 📝 总结

评论审核机制为MovieInfo项目提供了完整的内容安全保障：

**核心功能**：
1. **多层审核**：关键词、正则、机器学习、行为分析
2. **智能分级**：根据严重性和置信度自动分级处理
3. **人工审核**：复杂情况的人工介入机制
4. **实时监控**：完整的审核过程监控和指标

**技术特性**：
- 高效的关键词过滤算法
- 灵活的规则配置系统
- 智能的机器学习集成
- 完善的审核流程管理

**安全保障**：
- 多维度内容检测
- 用户行为分析
- 审核结果追溯
- 申诉处理机制

下一步，我们将实现统计分析功能，为平台运营提供数据支持。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第40步：统计分析功能](40-analytics.md)
