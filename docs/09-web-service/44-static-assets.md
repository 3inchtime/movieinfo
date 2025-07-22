# 第44步：静态资源管理

## 📋 概述

静态资源管理是Web应用性能优化的重要环节，包括CSS、JavaScript、图片、字体等文件的组织、压缩、缓存和分发。MovieInfo项目采用现代化的资源管理策略，提供高效的静态资源服务。

## 🎯 设计目标

### 1. **性能优化**
- 资源压缩和合并
- 浏览器缓存策略
- CDN分发支持
- 懒加载机制

### 2. **开发效率**
- 自动化构建流程
- 热重载支持
- 版本管理
- 依赖管理

### 3. **用户体验**
- 快速加载速度
- 渐进式加载
- 离线缓存支持
- 响应式资源

## 🔧 静态资源架构

### 1. **资源管理器结构**

```go
// 静态资源配置
type AssetManagerConfig struct {
    StaticDir       string            `yaml:"static_dir"`
    PublicPath      string            `yaml:"public_path"`
    CDNBaseURL      string            `yaml:"cdn_base_url"`
    EnableMinify    bool              `yaml:"enable_minify"`
    EnableGzip      bool              `yaml:"enable_gzip"`
    EnableVersioning bool             `yaml:"enable_versioning"`
    CacheMaxAge     time.Duration     `yaml:"cache_max_age"`
    Manifest        map[string]string `yaml:"manifest"`
}

// 资源管理器
type AssetManager struct {
    config      *AssetManagerConfig
    manifest    map[string]string
    cache       *cache.Cache
    compressor  *Compressor
    versioner   *Versioner
    logger      *logrus.Logger
    mutex       sync.RWMutex
}

func NewAssetManager(config *AssetManagerConfig) *AssetManager {
    am := &AssetManager{
        config:     config,
        manifest:   make(map[string]string),
        cache:      cache.New(1*time.Hour, 2*time.Hour),
        compressor: NewCompressor(),
        versioner:  NewVersioner(),
        logger:     logrus.New(),
    }
    
    // 加载资源清单
    if err := am.loadManifest(); err != nil {
        am.logger.Errorf("Failed to load asset manifest: %v", err)
    }
    
    // 初始化资源处理
    if err := am.initializeAssets(); err != nil {
        am.logger.Errorf("Failed to initialize assets: %v", err)
    }
    
    return am
}

// 获取资源URL
func (am *AssetManager) GetAssetURL(path string) string {
    am.mutex.RLock()
    defer am.mutex.RUnlock()
    
    // 检查清单中的版本化路径
    if versionedPath, exists := am.manifest[path]; exists {
        path = versionedPath
    }
    
    // 添加版本号
    if am.config.EnableVersioning {
        path = am.addVersionToPath(path)
    }
    
    // 使用CDN或本地路径
    if am.config.CDNBaseURL != "" {
        return am.config.CDNBaseURL + "/" + path
    }
    
    return am.config.PublicPath + "/" + path
}

// 处理静态资源请求
func (am *AssetManager) ServeAsset(c *gin.Context) {
    path := c.Param("filepath")
    
    // 安全检查
    if !am.isValidPath(path) {
        c.Status(404)
        return
    }
    
    fullPath := filepath.Join(am.config.StaticDir, path)
    
    // 检查文件是否存在
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        c.Status(404)
        return
    }
    
    // 设置缓存头
    am.setCacheHeaders(c, path)
    
    // 检查条件请求
    if am.handleConditionalRequest(c, fullPath) {
        return
    }
    
    // 处理压缩
    if am.shouldCompress(path) && am.acceptsGzip(c) {
        am.serveCompressed(c, fullPath)
        return
    }
    
    // 直接服务文件
    c.File(fullPath)
}

// 设置缓存头
func (am *AssetManager) setCacheHeaders(c *gin.Context, path string) {
    // 根据文件类型设置不同的缓存策略
    ext := filepath.Ext(path)
    
    switch ext {
    case ".css", ".js":
        // CSS和JS文件长期缓存
        c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(am.config.CacheMaxAge.Seconds())))
    case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
        // 图片文件长期缓存
        c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int((24*time.Hour).Seconds())))
    case ".woff", ".woff2", ".ttf", ".eot":
        // 字体文件超长期缓存
        c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int((365*24*time.Hour).Seconds())))
    default:
        // 其他文件短期缓存
        c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int((1*time.Hour).Seconds())))
    }
    
    // 设置ETag
    if etag := am.generateETag(path); etag != "" {
        c.Header("ETag", etag)
    }
}
```

### 2. **资源压缩器**

```go
// 压缩器
type Compressor struct {
    gzipPool sync.Pool
    logger   *logrus.Logger
}

func NewCompressor() *Compressor {
    return &Compressor{
        gzipPool: sync.Pool{
            New: func() interface{} {
                w, _ := gzip.NewWriterLevel(nil, gzip.BestCompression)
                return w
            },
        },
        logger: logrus.New(),
    }
}

// 压缩文件
func (c *Compressor) CompressFile(inputPath, outputPath string) error {
    input, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer input.Close()
    
    output, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer output.Close()
    
    // 获取gzip写入器
    gzipWriter := c.gzipPool.Get().(*gzip.Writer)
    defer c.gzipPool.Put(gzipWriter)
    
    gzipWriter.Reset(output)
    defer gzipWriter.Close()
    
    // 复制数据
    _, err = io.Copy(gzipWriter, input)
    return err
}

// 压缩CSS
func (c *Compressor) CompressCSS(content string) string {
    // 移除注释
    content = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(content, "")
    
    // 移除多余空白
    content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
    
    // 移除不必要的分号和空格
    content = regexp.MustCompile(`;\s*}`).ReplaceAllString(content, "}")
    content = regexp.MustCompile(`\s*{\s*`).ReplaceAllString(content, "{")
    content = regexp.MustCompile(`;\s*`).ReplaceAllString(content, ";")
    content = regexp.MustCompile(`:\s*`).ReplaceAllString(content, ":")
    
    return strings.TrimSpace(content)
}

// 压缩JavaScript
func (c *Compressor) CompressJS(content string) string {
    // 简单的JS压缩（生产环境建议使用专业工具如UglifyJS）
    
    // 移除单行注释
    content = regexp.MustCompile(`//.*$`).ReplaceAllString(content, "")
    
    // 移除多行注释
    content = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(content, "")
    
    // 移除多余空白
    content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
    
    // 移除行尾分号前的空格
    content = regexp.MustCompile(`\s*;\s*`).ReplaceAllString(content, ";")
    
    return strings.TrimSpace(content)
}
```

### 3. **版本管理器**

```go
// 版本管理器
type Versioner struct {
    versions map[string]string
    mutex    sync.RWMutex
    logger   *logrus.Logger
}

func NewVersioner() *Versioner {
    return &Versioner{
        versions: make(map[string]string),
        logger:   logrus.New(),
    }
}

// 生成文件版本
func (v *Versioner) GenerateVersion(filePath string) (string, error) {
    v.mutex.Lock()
    defer v.mutex.Unlock()
    
    // 检查缓存
    if version, exists := v.versions[filePath]; exists {
        return version, nil
    }
    
    // 计算文件哈希
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()
    
    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }
    
    version := fmt.Sprintf("%x", hash.Sum(nil))[:8]
    v.versions[filePath] = version
    
    return version, nil
}

// 添加版本到路径
func (am *AssetManager) addVersionToPath(path string) string {
    fullPath := filepath.Join(am.config.StaticDir, path)
    
    if version, err := am.versioner.GenerateVersion(fullPath); err == nil {
        ext := filepath.Ext(path)
        base := strings.TrimSuffix(path, ext)
        return fmt.Sprintf("%s.%s%s", base, version, ext)
    }
    
    return path
}

// 生成ETag
func (am *AssetManager) generateETag(path string) string {
    fullPath := filepath.Join(am.config.StaticDir, path)
    
    if info, err := os.Stat(fullPath); err == nil {
        return fmt.Sprintf(`"%x-%x"`, info.Size(), info.ModTime().Unix())
    }
    
    return ""
}
```

### 4. **资源构建系统**

```go
// 构建配置
type BuildConfig struct {
    SourceDir   string   `yaml:"source_dir"`
    OutputDir   string   `yaml:"output_dir"`
    EntryPoints []string `yaml:"entry_points"`
    EnableWatch bool     `yaml:"enable_watch"`
    Minify      bool     `yaml:"minify"`
    SourceMap   bool     `yaml:"source_map"`
}

// 资源构建器
type AssetBuilder struct {
    config     *BuildConfig
    compressor *Compressor
    bundler    *Bundler
    watcher    *fsnotify.Watcher
    logger     *logrus.Logger
}

func NewAssetBuilder(config *BuildConfig) *AssetBuilder {
    ab := &AssetBuilder{
        config:     config,
        compressor: NewCompressor(),
        bundler:    NewBundler(),
        logger:     logrus.New(),
    }
    
    if config.EnableWatch {
        ab.setupWatcher()
    }
    
    return ab
}

// 构建所有资源
func (ab *AssetBuilder) BuildAll() error {
    ab.logger.Info("Starting asset build...")
    
    // 清理输出目录
    if err := ab.cleanOutputDir(); err != nil {
        return err
    }
    
    // 构建CSS
    if err := ab.buildCSS(); err != nil {
        return fmt.Errorf("failed to build CSS: %v", err)
    }
    
    // 构建JavaScript
    if err := ab.buildJS(); err != nil {
        return fmt.Errorf("failed to build JavaScript: %v", err)
    }
    
    // 复制静态文件
    if err := ab.copyStaticFiles(); err != nil {
        return fmt.Errorf("failed to copy static files: %v", err)
    }
    
    // 生成清单文件
    if err := ab.generateManifest(); err != nil {
        return fmt.Errorf("failed to generate manifest: %v", err)
    }
    
    ab.logger.Info("Asset build completed successfully")
    return nil
}

// 构建CSS
func (ab *AssetBuilder) buildCSS() error {
    cssDir := filepath.Join(ab.config.SourceDir, "css")
    outputDir := filepath.Join(ab.config.OutputDir, "css")
    
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return err
    }
    
    return filepath.Walk(cssDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !strings.HasSuffix(path, ".css") {
            return nil
        }
        
        // 读取CSS文件
        content, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }
        
        // 处理CSS（压缩、自动前缀等）
        processedContent := string(content)
        if ab.config.Minify {
            processedContent = ab.compressor.CompressCSS(processedContent)
        }
        
        // 计算相对路径
        relPath, err := filepath.Rel(cssDir, path)
        if err != nil {
            return err
        }
        
        outputPath := filepath.Join(outputDir, relPath)
        
        // 确保输出目录存在
        if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
            return err
        }
        
        // 写入处理后的文件
        return ioutil.WriteFile(outputPath, []byte(processedContent), 0644)
    })
}

// 构建JavaScript
func (ab *AssetBuilder) buildJS() error {
    jsDir := filepath.Join(ab.config.SourceDir, "js")
    outputDir := filepath.Join(ab.config.OutputDir, "js")
    
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return err
    }
    
    // 处理入口点
    for _, entryPoint := range ab.config.EntryPoints {
        if err := ab.buildJSEntry(entryPoint, outputDir); err != nil {
            return err
        }
    }
    
    return nil
}

// 构建JavaScript入口点
func (ab *AssetBuilder) buildJSEntry(entryPoint, outputDir string) error {
    inputPath := filepath.Join(ab.config.SourceDir, "js", entryPoint)
    
    // 读取入口文件
    content, err := ioutil.ReadFile(inputPath)
    if err != nil {
        return err
    }
    
    // 解析依赖并打包
    bundledContent, err := ab.bundler.Bundle(string(content), filepath.Dir(inputPath))
    if err != nil {
        return err
    }
    
    // 压缩
    if ab.config.Minify {
        bundledContent = ab.compressor.CompressJS(bundledContent)
    }
    
    // 输出文件
    outputPath := filepath.Join(outputDir, entryPoint)
    if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
        return err
    }
    
    return ioutil.WriteFile(outputPath, []byte(bundledContent), 0644)
}

// 生成资源清单
func (ab *AssetBuilder) generateManifest() error {
    manifest := make(map[string]string)
    
    err := filepath.Walk(ab.config.OutputDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil
        }
        
        // 计算相对路径
        relPath, err := filepath.Rel(ab.config.OutputDir, path)
        if err != nil {
            return err
        }
        
        // 生成版本化路径
        hash := md5.New()
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()
        
        if _, err := io.Copy(hash, file); err != nil {
            return err
        }
        
        version := fmt.Sprintf("%x", hash.Sum(nil))[:8]
        ext := filepath.Ext(relPath)
        base := strings.TrimSuffix(relPath, ext)
        versionedPath := fmt.Sprintf("%s.%s%s", base, version, ext)
        
        manifest[relPath] = versionedPath
        
        return nil
    })
    
    if err != nil {
        return err
    }
    
    // 写入清单文件
    manifestPath := filepath.Join(ab.config.OutputDir, "manifest.json")
    manifestData, err := json.MarshalIndent(manifest, "", "  ")
    if err != nil {
        return err
    }
    
    return ioutil.WriteFile(manifestPath, manifestData, 0644)
}
```

### 5. **CDN集成**

```go
// CDN配置
type CDNConfig struct {
    Provider    string            `yaml:"provider"`    // aws, aliyun, qiniu
    BaseURL     string            `yaml:"base_url"`
    AccessKey   string            `yaml:"access_key"`
    SecretKey   string            `yaml:"secret_key"`
    Bucket      string            `yaml:"bucket"`
    Region      string            `yaml:"region"`
    Options     map[string]string `yaml:"options"`
}

// CDN管理器
type CDNManager struct {
    config   *CDNConfig
    uploader CDNUploader
    logger   *logrus.Logger
}

type CDNUploader interface {
    Upload(localPath, remotePath string) error
    Delete(remotePath string) error
    GetURL(remotePath string) string
}

func NewCDNManager(config *CDNConfig) *CDNManager {
    var uploader CDNUploader
    
    switch config.Provider {
    case "aws":
        uploader = NewAWSUploader(config)
    case "aliyun":
        uploader = NewAliyunUploader(config)
    case "qiniu":
        uploader = NewQiniuUploader(config)
    default:
        uploader = NewLocalUploader(config)
    }
    
    return &CDNManager{
        config:   config,
        uploader: uploader,
        logger:   logrus.New(),
    }
}

// 上传资源到CDN
func (cm *CDNManager) UploadAssets(localDir string) error {
    return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil
        }
        
        // 计算远程路径
        relPath, err := filepath.Rel(localDir, path)
        if err != nil {
            return err
        }
        
        remotePath := strings.ReplaceAll(relPath, "\\", "/")
        
        // 上传文件
        if err := cm.uploader.Upload(path, remotePath); err != nil {
            cm.logger.Errorf("Failed to upload %s: %v", path, err)
            return err
        }
        
        cm.logger.Infof("Uploaded: %s -> %s", path, remotePath)
        return nil
    })
}
```

## 📊 性能监控

### 1. **资源性能指标**

```go
type AssetMetrics struct {
    requestCount    *prometheus.CounterVec
    responseTime    *prometheus.HistogramVec
    cacheHitRate    *prometheus.CounterVec
    compressionRate prometheus.Gauge
    transferSize    *prometheus.HistogramVec
}

func NewAssetMetrics() *AssetMetrics {
    return &AssetMetrics{
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "asset_requests_total",
                Help: "Total number of asset requests",
            },
            []string{"type", "status"},
        ),
        responseTime: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "asset_response_time_seconds",
                Help: "Asset response time",
            },
            []string{"type"},
        ),
        cacheHitRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "asset_cache_operations_total",
                Help: "Total number of asset cache operations",
            },
            []string{"type"},
        ),
        compressionRate: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "asset_compression_ratio",
                Help: "Asset compression ratio",
            },
        ),
        transferSize: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "asset_transfer_size_bytes",
                Help: "Asset transfer size in bytes",
            },
            []string{"type", "compressed"},
        ),
    }
}
```

## 📝 总结

静态资源管理为MovieInfo项目提供了高效的资源服务能力：

**核心功能**：
1. **智能压缩**：CSS、JavaScript和图片的自动压缩优化
2. **版本管理**：基于内容哈希的版本控制和缓存破坏
3. **CDN集成**：支持多种CDN服务商的资源分发
4. **构建系统**：自动化的资源构建和打包流程

**性能优化**：
- 高效的缓存策略
- Gzip压缩支持
- 条件请求处理
- 资源懒加载

**开发体验**：
- 热重载支持
- 自动化构建
- 清单文件管理
- 错误监控

下一步，我们将实现页面渲染逻辑，将所有组件整合为完整的页面展示。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第45步：页面渲染逻辑](45-page-rendering.md)
