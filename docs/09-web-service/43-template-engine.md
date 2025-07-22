# 第43步：模板引擎集成

## 📋 概述

模板引擎是Web应用的重要组件，负责将数据与HTML模板结合生成最终的页面内容。MovieInfo项目采用Go的html/template包，结合自定义函数和布局系统，提供灵活、高效的页面渲染能力。

## 🎯 设计目标

### 1. **开发效率**
- 模板继承和布局系统
- 组件化的模板设计
- 自动重载机制
- 丰富的模板函数

### 2. **性能优化**
- 模板预编译和缓存
- 部分模板更新
- 静态资源版本管理
- 压缩和优化

### 3. **安全性**
- XSS防护机制
- 输入数据转义
- CSRF保护
- 内容安全策略

## 🔧 模板引擎架构

### 1. **模板引擎结构**

```go
// 模板引擎配置
type TemplateEngineConfig struct {
    TemplateDir     string        `yaml:"template_dir"`
    StaticDir       string        `yaml:"static_dir"`
    EnableCache     bool          `yaml:"enable_cache"`
    EnableReload    bool          `yaml:"enable_reload"`
    EnableMinify    bool          `yaml:"enable_minify"`
    CacheTimeout    time.Duration `yaml:"cache_timeout"`
    FuncMap         template.FuncMap
}

// 模板引擎
type TemplateEngine struct {
    config      *TemplateEngineConfig
    templates   map[string]*template.Template
    layouts     map[string]*template.Template
    funcMap     template.FuncMap
    cache       *cache.Cache
    mutex       sync.RWMutex
    logger      *logrus.Logger
    watcher     *fsnotify.Watcher
}

func NewTemplateEngine(config *TemplateEngineConfig) *TemplateEngine {
    te := &TemplateEngine{
        config:    config,
        templates: make(map[string]*template.Template),
        layouts:   make(map[string]*template.Template),
        funcMap:   make(template.FuncMap),
        cache:     cache.New(config.CacheTimeout, config.CacheTimeout*2),
        logger:    logrus.New(),
    }
    
    // 注册默认函数
    te.registerDefaultFunctions()
    
    // 加载模板
    if err := te.loadTemplates(); err != nil {
        te.logger.Fatalf("Failed to load templates: %v", err)
    }
    
    // 启用文件监控（开发模式）
    if config.EnableReload {
        te.setupFileWatcher()
    }
    
    return te
}

// 渲染模板
func (te *TemplateEngine) Render() gin.HTMLRender {
    return &TemplateRender{
        engine: te,
    }
}

type TemplateRender struct {
    engine *TemplateEngine
}

func (tr *TemplateRender) Instance(name string, data interface{}) render.Render {
    return &TemplateInstance{
        engine:   tr.engine,
        template: name,
        data:     data,
    }
}

type TemplateInstance struct {
    engine   *TemplateEngine
    template string
    data     interface{}
}

func (ti *TemplateInstance) Render(w io.Writer) error {
    return ti.engine.ExecuteTemplate(w, ti.template, ti.data)
}

func (ti *TemplateInstance) WriteContentType(w http.ResponseWriter) {
    header := w.Header()
    if val := header["Content-Type"]; len(val) == 0 {
        header["Content-Type"] = []string{"text/html; charset=utf-8"}
    }
}
```

### 2. **模板加载和管理**

```go
// 加载所有模板
func (te *TemplateEngine) loadTemplates() error {
    te.mutex.Lock()
    defer te.mutex.Unlock()
    
    // 清空现有模板
    te.templates = make(map[string]*template.Template)
    te.layouts = make(map[string]*template.Template)
    
    // 加载布局模板
    if err := te.loadLayouts(); err != nil {
        return fmt.Errorf("failed to load layouts: %v", err)
    }
    
    // 加载页面模板
    if err := te.loadPageTemplates(); err != nil {
        return fmt.Errorf("failed to load page templates: %v", err)
    }
    
    // 加载组件模板
    if err := te.loadComponentTemplates(); err != nil {
        return fmt.Errorf("failed to load component templates: %v", err)
    }
    
    te.logger.Infof("Loaded %d templates and %d layouts", len(te.templates), len(te.layouts))
    return nil
}

// 加载布局模板
func (te *TemplateEngine) loadLayouts() error {
    layoutDir := filepath.Join(te.config.TemplateDir, "layouts")
    
    return filepath.Walk(layoutDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !strings.HasSuffix(path, ".html") {
            return nil
        }
        
        name := strings.TrimSuffix(filepath.Base(path), ".html")
        
        tmpl, err := template.New(name).Funcs(te.funcMap).ParseFiles(path)
        if err != nil {
            return fmt.Errorf("failed to parse layout %s: %v", name, err)
        }
        
        te.layouts[name] = tmpl
        te.logger.Debugf("Loaded layout: %s", name)
        
        return nil
    })
}

// 加载页面模板
func (te *TemplateEngine) loadPageTemplates() error {
    pageDir := filepath.Join(te.config.TemplateDir, "pages")
    
    return filepath.Walk(pageDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !strings.HasSuffix(path, ".html") {
            return nil
        }
        
        // 获取相对路径作为模板名
        relPath, err := filepath.Rel(pageDir, path)
        if err != nil {
            return err
        }
        
        name := strings.TrimSuffix(relPath, ".html")
        name = strings.ReplaceAll(name, "\\", "/") // 统一使用正斜杠
        
        // 解析模板及其依赖
        tmpl, err := te.parseTemplateWithDependencies(path, name)
        if err != nil {
            return fmt.Errorf("failed to parse template %s: %v", name, err)
        }
        
        te.templates[name] = tmpl
        te.logger.Debugf("Loaded template: %s", name)
        
        return nil
    })
}

// 解析模板及其依赖
func (te *TemplateEngine) parseTemplateWithDependencies(path, name string) (*template.Template, error) {
    // 读取模板内容
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    // 解析模板头部信息
    templateInfo := te.parseTemplateHeader(string(content))
    
    // 创建模板
    tmpl := template.New(name).Funcs(te.funcMap)
    
    // 如果指定了布局，先加载布局
    if templateInfo.Layout != "" {
        if layout, exists := te.layouts[templateInfo.Layout]; exists {
            tmpl, err = tmpl.AddParseTree(templateInfo.Layout, layout.Tree)
            if err != nil {
                return nil, err
            }
        }
    }
    
    // 加载组件依赖
    for _, component := range templateInfo.Components {
        componentPath := filepath.Join(te.config.TemplateDir, "components", component+".html")
        if _, err := os.Stat(componentPath); err == nil {
            tmpl, err = tmpl.ParseFiles(componentPath)
            if err != nil {
                return nil, err
            }
        }
    }
    
    // 解析主模板
    tmpl, err = tmpl.Parse(string(content))
    if err != nil {
        return nil, err
    }
    
    return tmpl, nil
}

// 模板信息结构
type TemplateInfo struct {
    Layout     string   `yaml:"layout"`
    Components []string `yaml:"components"`
    Title      string   `yaml:"title"`
    Meta       map[string]string `yaml:"meta"`
}

// 解析模板头部信息
func (te *TemplateEngine) parseTemplateHeader(content string) *TemplateInfo {
    info := &TemplateInfo{
        Meta: make(map[string]string),
    }
    
    // 查找YAML前置内容
    if strings.HasPrefix(content, "---\n") {
        endIndex := strings.Index(content[4:], "\n---\n")
        if endIndex != -1 {
            yamlContent := content[4 : endIndex+4]
            yaml.Unmarshal([]byte(yamlContent), info)
        }
    }
    
    return info
}
```

### 3. **模板函数注册**

```go
// 注册默认模板函数
func (te *TemplateEngine) registerDefaultFunctions() {
    te.funcMap = template.FuncMap{
        // 字符串处理
        "title":     strings.Title,
        "upper":     strings.ToUpper,
        "lower":     strings.ToLower,
        "trim":      strings.TrimSpace,
        "truncate":  te.truncateString,
        "slug":      te.generateSlug,
        
        // 数字处理
        "add":       te.add,
        "sub":       te.sub,
        "mul":       te.mul,
        "div":       te.div,
        "mod":       te.mod,
        "round":     te.round,
        
        // 时间处理
        "now":       time.Now,
        "formatTime": te.formatTime,
        "timeAgo":   te.timeAgo,
        "duration":  te.formatDuration,
        
        // 数组和切片
        "slice":     te.slice,
        "len":       te.length,
        "first":     te.first,
        "last":      te.last,
        "reverse":   te.reverse,
        "sort":      te.sort,
        
        // 条件判断
        "eq":        te.equal,
        "ne":        te.notEqual,
        "lt":        te.lessThan,
        "le":        te.lessEqual,
        "gt":        te.greaterThan,
        "ge":        te.greaterEqual,
        "in":        te.contains,
        
        // HTML和URL
        "safeHTML":  te.safeHTML,
        "safeCSS":   te.safeCSS,
        "safeJS":    te.safeJS,
        "url":       te.buildURL,
        "asset":     te.assetURL,
        
        // 业务相关
        "rating":    te.formatRating,
        "currency":  te.formatCurrency,
        "fileSize":  te.formatFileSize,
        "avatar":    te.avatarURL,
        
        // 模板包含
        "include":   te.includeTemplate,
        "partial":   te.renderPartial,
        
        // 国际化
        "t":         te.translate,
        "lang":      te.getCurrentLanguage,
        
        // 调试
        "debug":     te.debug,
        "dump":      te.dump,
    }
    
    // 合并用户自定义函数
    for name, fn := range te.config.FuncMap {
        te.funcMap[name] = fn
    }
}

// 字符串截断
func (te *TemplateEngine) truncateString(s string, length int) string {
    if len(s) <= length {
        return s
    }
    
    runes := []rune(s)
    if len(runes) <= length {
        return s
    }
    
    return string(runes[:length]) + "..."
}

// 生成URL友好的slug
func (te *TemplateEngine) generateSlug(s string) string {
    // 转小写
    s = strings.ToLower(s)
    
    // 替换空格和特殊字符
    reg := regexp.MustCompile(`[^a-z0-9\-]`)
    s = reg.ReplaceAllString(s, "-")
    
    // 移除多余的连字符
    reg = regexp.MustCompile(`-+`)
    s = reg.ReplaceAllString(s, "-")
    
    // 移除首尾连字符
    s = strings.Trim(s, "-")
    
    return s
}

// 时间格式化
func (te *TemplateEngine) formatTime(t time.Time, format string) string {
    switch format {
    case "date":
        return t.Format("2006-01-02")
    case "datetime":
        return t.Format("2006-01-02 15:04:05")
    case "time":
        return t.Format("15:04:05")
    case "iso":
        return t.Format(time.RFC3339)
    default:
        return t.Format(format)
    }
}

// 相对时间
func (te *TemplateEngine) timeAgo(t time.Time) string {
    duration := time.Since(t)
    
    if duration < time.Minute {
        return "刚刚"
    } else if duration < time.Hour {
        return fmt.Sprintf("%d分钟前", int(duration.Minutes()))
    } else if duration < 24*time.Hour {
        return fmt.Sprintf("%d小时前", int(duration.Hours()))
    } else if duration < 30*24*time.Hour {
        return fmt.Sprintf("%d天前", int(duration.Hours()/24))
    } else if duration < 365*24*time.Hour {
        return fmt.Sprintf("%d个月前", int(duration.Hours()/(24*30)))
    } else {
        return fmt.Sprintf("%d年前", int(duration.Hours()/(24*365)))
    }
}

// 安全HTML
func (te *TemplateEngine) safeHTML(s string) template.HTML {
    return template.HTML(s)
}

// 构建URL
func (te *TemplateEngine) buildURL(path string, params ...interface{}) string {
    if len(params) == 0 {
        return path
    }
    
    values := url.Values{}
    for i := 0; i < len(params); i += 2 {
        if i+1 < len(params) {
            key := fmt.Sprintf("%v", params[i])
            value := fmt.Sprintf("%v", params[i+1])
            values.Add(key, value)
        }
    }
    
    if len(values) > 0 {
        return path + "?" + values.Encode()
    }
    
    return path
}

// 静态资源URL
func (te *TemplateEngine) assetURL(path string) string {
    // 添加版本号或CDN前缀
    version := te.getAssetVersion(path)
    if version != "" {
        return fmt.Sprintf("/static/%s?v=%s", path, version)
    }
    return fmt.Sprintf("/static/%s", path)
}

// 获取资源版本
func (te *TemplateEngine) getAssetVersion(path string) string {
    // 可以基于文件修改时间或构建版本生成
    fullPath := filepath.Join(te.config.StaticDir, path)
    if info, err := os.Stat(fullPath); err == nil {
        return fmt.Sprintf("%d", info.ModTime().Unix())
    }
    return ""
}

// 评分格式化
func (te *TemplateEngine) formatRating(rating float64) string {
    return fmt.Sprintf("%.1f", rating)
}

// 头像URL
func (te *TemplateEngine) avatarURL(userID, avatar string) string {
    if avatar != "" {
        return avatar
    }
    
    // 生成默认头像
    hash := md5.Sum([]byte(userID))
    return fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=identicon&s=80", hash)
}
```

### 4. **模板布局系统**

```html
<!-- layouts/base.html -->
<!DOCTYPE html>
<html lang="{{.Lang | default "zh-CN"}}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{block "title" .}}MovieInfo - 电影信息平台{{end}}</title>
    
    <!-- SEO Meta -->
    {{block "meta" .}}
    <meta name="description" content="MovieInfo是一个专业的电影信息平台">
    <meta name="keywords" content="电影,影评,评分,推荐">
    {{end}}
    
    <!-- CSS -->
    <link rel="stylesheet" href="{{asset "css/bootstrap.min.css"}}">
    <link rel="stylesheet" href="{{asset "css/main.css"}}">
    {{block "css" .}}{{end}}
    
    <!-- Favicon -->
    <link rel="icon" type="image/x-icon" href="/favicon.ico">
</head>
<body class="{{block "body-class" .}}{{end}}">
    <!-- Header -->
    {{template "header" .}}
    
    <!-- Main Content -->
    <main class="main-content">
        {{block "content" .}}{{end}}
    </main>
    
    <!-- Footer -->
    {{template "footer" .}}
    
    <!-- JavaScript -->
    <script src="{{asset "js/jquery.min.js"}}"></script>
    <script src="{{asset "js/bootstrap.min.js"}}"></script>
    <script src="{{asset "js/main.js"}}"></script>
    {{block "js" .}}{{end}}
</body>
</html>

<!-- components/header.html -->
{{define "header"}}
<header class="header">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
        <div class="container">
            <a class="navbar-brand" href="/">
                <img src="{{asset "images/logo.png"}}" alt="MovieInfo" height="32">
            </a>
            
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
                <span class="navbar-toggler-icon"></span>
            </button>
            
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav me-auto">
                    <li class="nav-item">
                        <a class="nav-link" href="/">首页</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/movies">电影</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/movies/popular">热门</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/movies/top-rated">高分</a>
                    </li>
                </ul>
                
                <!-- Search Form -->
                <form class="d-flex me-3" action="/movies/search" method="GET">
                    <input class="form-control" type="search" name="q" placeholder="搜索电影..." value="{{.SearchQuery}}">
                    <button class="btn btn-outline-light" type="submit">搜索</button>
                </form>
                
                <!-- User Menu -->
                <ul class="navbar-nav">
                    {{if .User}}
                    <li class="nav-item dropdown">
                        <a class="nav-link dropdown-toggle" href="#" id="userDropdown" role="button" data-bs-toggle="dropdown">
                            <img src="{{avatar .User.ID .User.Avatar}}" alt="{{.User.Username}}" class="rounded-circle" width="24" height="24">
                            {{.User.Username}}
                        </a>
                        <ul class="dropdown-menu">
                            <li><a class="dropdown-item" href="/user/profile">个人资料</a></li>
                            <li><a class="dropdown-item" href="/user/favorites">我的收藏</a></li>
                            <li><a class="dropdown-item" href="/user/watchlist">观看列表</a></li>
                            <li><hr class="dropdown-divider"></li>
                            <li><a class="dropdown-item" href="/user/logout">退出登录</a></li>
                        </ul>
                    </li>
                    {{else}}
                    <li class="nav-item">
                        <a class="nav-link" href="/user/login">登录</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/user/register">注册</a>
                    </li>
                    {{end}}
                </ul>
            </div>
        </div>
    </nav>
</header>
{{end}}
```

### 5. **文件监控和热重载**

```go
// 设置文件监控
func (te *TemplateEngine) setupFileWatcher() {
    var err error
    te.watcher, err = fsnotify.NewWatcher()
    if err != nil {
        te.logger.Errorf("Failed to create file watcher: %v", err)
        return
    }
    
    // 监控模板目录
    err = filepath.Walk(te.config.TemplateDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return te.watcher.Add(path)
        }
        
        return nil
    })
    
    if err != nil {
        te.logger.Errorf("Failed to setup file watcher: %v", err)
        return
    }
    
    // 启动监控协程
    go te.watchFiles()
    
    te.logger.Info("Template file watcher started")
}

// 监控文件变化
func (te *TemplateEngine) watchFiles() {
    for {
        select {
        case event, ok := <-te.watcher.Events:
            if !ok {
                return
            }
            
            if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
                if strings.HasSuffix(event.Name, ".html") {
                    te.logger.Infof("Template file changed: %s", event.Name)
                    
                    // 延迟重载，避免频繁更新
                    time.AfterFunc(100*time.Millisecond, func() {
                        if err := te.loadTemplates(); err != nil {
                            te.logger.Errorf("Failed to reload templates: %v", err)
                        } else {
                            te.logger.Info("Templates reloaded successfully")
                        }
                    })
                }
            }
            
        case err, ok := <-te.watcher.Errors:
            if !ok {
                return
            }
            te.logger.Errorf("File watcher error: %v", err)
        }
    }
}

// 执行模板
func (te *TemplateEngine) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
    te.mutex.RLock()
    tmpl, exists := te.templates[name]
    te.mutex.RUnlock()
    
    if !exists {
        return fmt.Errorf("template %s not found", name)
    }
    
    // 如果启用缓存，检查缓存
    if te.config.EnableCache {
        cacheKey := fmt.Sprintf("template:%s:%s", name, te.generateDataHash(data))
        if cached, found := te.cache.Get(cacheKey); found {
            if content, ok := cached.([]byte); ok {
                _, err := w.Write(content)
                return err
            }
        }
        
        // 缓存未命中，渲染并缓存
        var buf bytes.Buffer
        if err := tmpl.Execute(&buf, data); err != nil {
            return err
        }
        
        content := buf.Bytes()
        te.cache.Set(cacheKey, content, te.config.CacheTimeout)
        
        _, err := w.Write(content)
        return err
    }
    
    // 直接渲染
    return tmpl.Execute(w, data)
}

// 生成数据哈希
func (te *TemplateEngine) generateDataHash(data interface{}) string {
    hash := md5.New()
    
    // 简单的数据序列化用于缓存键
    if data != nil {
        if jsonData, err := json.Marshal(data); err == nil {
            hash.Write(jsonData)
        }
    }
    
    return fmt.Sprintf("%x", hash.Sum(nil))[:8]
}
```

## 📝 总结

模板引擎集成为MovieInfo项目提供了强大的页面渲染能力：

**核心功能**：
1. **灵活的模板系统**：支持布局继承、组件化和模板函数
2. **高性能渲染**：模板缓存、预编译和优化机制
3. **开发友好**：热重载、错误提示和调试支持
4. **安全保障**：XSS防护、数据转义和内容安全

**技术特性**：
- 模块化的模板组织
- 丰富的模板函数库
- 智能的缓存策略
- 实时的文件监控

**用户体验**：
- 快速的页面加载
- 响应式设计支持
- SEO友好的结构
- 国际化支持

下一步，我们将实现静态资源管理，为网站提供高效的资源服务。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第44步：静态资源管理](44-static-assets.md)
