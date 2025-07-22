# 第46步：页面布局设计

## 📋 概述

页面布局设计是前端开发的基础，定义了网站的整体视觉结构和用户交互模式。MovieInfo项目采用现代化的响应式设计，提供一致的用户体验和优雅的视觉呈现。

## 🎯 设计目标

### 1. **用户体验**
- 直观的导航结构
- 一致的视觉语言
- 流畅的交互体验
- 无障碍访问支持

### 2. **响应式设计**
- 移动优先策略
- 多设备适配
- 弹性布局系统
- 触摸友好界面

### 3. **性能优化**
- 快速加载速度
- 渐进式增强
- 资源优化
- 缓存策略

## 🎨 设计系统

### 1. **色彩系统**

```css
/* 主色调 */
:root {
  /* 品牌色 */
  --primary-color: #1a73e8;
  --primary-light: #4285f4;
  --primary-dark: #1557b0;
  
  /* 辅助色 */
  --secondary-color: #34a853;
  --accent-color: #fbbc04;
  --warning-color: #ea4335;
  
  /* 中性色 */
  --gray-50: #f8f9fa;
  --gray-100: #f1f3f4;
  --gray-200: #e8eaed;
  --gray-300: #dadce0;
  --gray-400: #bdc1c6;
  --gray-500: #9aa0a6;
  --gray-600: #80868b;
  --gray-700: #5f6368;
  --gray-800: #3c4043;
  --gray-900: #202124;
  
  /* 语义色 */
  --success-color: #137333;
  --error-color: #d93025;
  --warning-color: #f9ab00;
  --info-color: #1a73e8;
  
  /* 背景色 */
  --bg-primary: #ffffff;
  --bg-secondary: #f8f9fa;
  --bg-tertiary: #f1f3f4;
  --bg-dark: #202124;
  
  /* 文字色 */
  --text-primary: #202124;
  --text-secondary: #5f6368;
  --text-tertiary: #80868b;
  --text-inverse: #ffffff;
}

/* 深色主题 */
[data-theme="dark"] {
  --bg-primary: #202124;
  --bg-secondary: #303134;
  --bg-tertiary: #3c4043;
  
  --text-primary: #e8eaed;
  --text-secondary: #9aa0a6;
  --text-tertiary: #80868b;
}
```

### 2. **字体系统**

```css
/* 字体定义 */
:root {
  /* 字体族 */
  --font-family-primary: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  --font-family-secondary: 'Roboto', sans-serif;
  --font-family-mono: 'JetBrains Mono', 'Fira Code', monospace;
  
  /* 字体大小 */
  --font-size-xs: 0.75rem;    /* 12px */
  --font-size-sm: 0.875rem;   /* 14px */
  --font-size-base: 1rem;     /* 16px */
  --font-size-lg: 1.125rem;   /* 18px */
  --font-size-xl: 1.25rem;    /* 20px */
  --font-size-2xl: 1.5rem;    /* 24px */
  --font-size-3xl: 1.875rem;  /* 30px */
  --font-size-4xl: 2.25rem;   /* 36px */
  --font-size-5xl: 3rem;      /* 48px */
  
  /* 行高 */
  --line-height-tight: 1.25;
  --line-height-normal: 1.5;
  --line-height-relaxed: 1.75;
  
  /* 字重 */
  --font-weight-light: 300;
  --font-weight-normal: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;
}

/* 字体类 */
.text-xs { font-size: var(--font-size-xs); }
.text-sm { font-size: var(--font-size-sm); }
.text-base { font-size: var(--font-size-base); }
.text-lg { font-size: var(--font-size-lg); }
.text-xl { font-size: var(--font-size-xl); }
.text-2xl { font-size: var(--font-size-2xl); }
.text-3xl { font-size: var(--font-size-3xl); }
.text-4xl { font-size: var(--font-size-4xl); }
.text-5xl { font-size: var(--font-size-5xl); }

.font-light { font-weight: var(--font-weight-light); }
.font-normal { font-weight: var(--font-weight-normal); }
.font-medium { font-weight: var(--font-weight-medium); }
.font-semibold { font-weight: var(--font-weight-semibold); }
.font-bold { font-weight: var(--font-weight-bold); }
```

### 3. **间距系统**

```css
/* 间距定义 */
:root {
  --spacing-0: 0;
  --spacing-1: 0.25rem;   /* 4px */
  --spacing-2: 0.5rem;    /* 8px */
  --spacing-3: 0.75rem;   /* 12px */
  --spacing-4: 1rem;      /* 16px */
  --spacing-5: 1.25rem;   /* 20px */
  --spacing-6: 1.5rem;    /* 24px */
  --spacing-8: 2rem;      /* 32px */
  --spacing-10: 2.5rem;   /* 40px */
  --spacing-12: 3rem;     /* 48px */
  --spacing-16: 4rem;     /* 64px */
  --spacing-20: 5rem;     /* 80px */
  --spacing-24: 6rem;     /* 96px */
  --spacing-32: 8rem;     /* 128px */
}

/* 间距工具类 */
.m-0 { margin: var(--spacing-0); }
.m-1 { margin: var(--spacing-1); }
.m-2 { margin: var(--spacing-2); }
.m-3 { margin: var(--spacing-3); }
.m-4 { margin: var(--spacing-4); }
.m-5 { margin: var(--spacing-5); }
.m-6 { margin: var(--spacing-6); }
.m-8 { margin: var(--spacing-8); }

.p-0 { padding: var(--spacing-0); }
.p-1 { padding: var(--spacing-1); }
.p-2 { padding: var(--spacing-2); }
.p-3 { padding: var(--spacing-3); }
.p-4 { padding: var(--spacing-4); }
.p-5 { padding: var(--spacing-5); }
.p-6 { padding: var(--spacing-6); }
.p-8 { padding: var(--spacing-8); }
```

## 🏗️ 布局结构

### 1. **整体布局框架**

```html
<!DOCTYPE html>
<html lang="zh-CN" data-theme="light">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    
    <!-- 预加载关键资源 -->
    <link rel="preload" href="/static/css/critical.css" as="style">
    <link rel="preload" href="/static/fonts/inter-var.woff2" as="font" type="font/woff2" crossorigin>
    
    <!-- 关键CSS内联 -->
    <style>
        /* 关键路径CSS */
        body { margin: 0; font-family: var(--font-family-primary); }
        .header { position: sticky; top: 0; z-index: 1000; }
        .main-content { min-height: calc(100vh - 120px); }
        .loading { display: flex; justify-content: center; align-items: center; height: 200px; }
    </style>
    
    <!-- 非关键CSS异步加载 -->
    <link rel="preload" href="/static/css/main.css" as="style" onload="this.onload=null;this.rel='stylesheet'">
    <noscript><link rel="stylesheet" href="/static/css/main.css"></noscript>
</head>
<body class="{{.BodyClass}}">
    <!-- 跳转到主内容链接（无障碍） -->
    <a href="#main-content" class="skip-link">跳转到主内容</a>
    
    <!-- 页面头部 -->
    <header class="header" role="banner">
        {{template "header" .}}
    </header>
    
    <!-- 主要内容区域 -->
    <main id="main-content" class="main-content" role="main">
        {{block "content" .}}{{end}}
    </main>
    
    <!-- 页面底部 -->
    <footer class="footer" role="contentinfo">
        {{template "footer" .}}
    </footer>
    
    <!-- 返回顶部按钮 -->
    <button id="back-to-top" class="back-to-top" aria-label="返回顶部" style="display: none;">
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M7.41 15.41L12 10.83l4.59 4.58L18 14l-6-6-6 6z"/>
        </svg>
    </button>
    
    <!-- 加载指示器 -->
    <div id="loading-indicator" class="loading-indicator" style="display: none;">
        <div class="spinner"></div>
    </div>
    
    <!-- JavaScript -->
    <script src="/static/js/critical.js"></script>
    <script src="/static/js/main.js" defer></script>
</body>
</html>
```

### 2. **响应式网格系统**

```css
/* 容器 */
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-4);
}

.container-fluid {
  width: 100%;
  padding: 0 var(--spacing-4);
}

/* 网格系统 */
.row {
  display: flex;
  flex-wrap: wrap;
  margin: 0 calc(var(--spacing-3) * -1);
}

.col {
  flex: 1;
  padding: 0 var(--spacing-3);
}

/* 响应式列 */
.col-1 { flex: 0 0 8.333333%; }
.col-2 { flex: 0 0 16.666667%; }
.col-3 { flex: 0 0 25%; }
.col-4 { flex: 0 0 33.333333%; }
.col-5 { flex: 0 0 41.666667%; }
.col-6 { flex: 0 0 50%; }
.col-7 { flex: 0 0 58.333333%; }
.col-8 { flex: 0 0 66.666667%; }
.col-9 { flex: 0 0 75%; }
.col-10 { flex: 0 0 83.333333%; }
.col-11 { flex: 0 0 91.666667%; }
.col-12 { flex: 0 0 100%; }

/* 断点 */
@media (min-width: 576px) {
  .col-sm-1 { flex: 0 0 8.333333%; }
  .col-sm-2 { flex: 0 0 16.666667%; }
  .col-sm-3 { flex: 0 0 25%; }
  .col-sm-4 { flex: 0 0 33.333333%; }
  .col-sm-6 { flex: 0 0 50%; }
  .col-sm-12 { flex: 0 0 100%; }
}

@media (min-width: 768px) {
  .col-md-1 { flex: 0 0 8.333333%; }
  .col-md-2 { flex: 0 0 16.666667%; }
  .col-md-3 { flex: 0 0 25%; }
  .col-md-4 { flex: 0 0 33.333333%; }
  .col-md-6 { flex: 0 0 50%; }
  .col-md-8 { flex: 0 0 66.666667%; }
  .col-md-12 { flex: 0 0 100%; }
}

@media (min-width: 992px) {
  .col-lg-1 { flex: 0 0 8.333333%; }
  .col-lg-2 { flex: 0 0 16.666667%; }
  .col-lg-3 { flex: 0 0 25%; }
  .col-lg-4 { flex: 0 0 33.333333%; }
  .col-lg-6 { flex: 0 0 50%; }
  .col-lg-8 { flex: 0 0 66.666667%; }
  .col-lg-9 { flex: 0 0 75%; }
  .col-lg-12 { flex: 0 0 100%; }
}

@media (min-width: 1200px) {
  .col-xl-1 { flex: 0 0 8.333333%; }
  .col-xl-2 { flex: 0 0 16.666667%; }
  .col-xl-3 { flex: 0 0 25%; }
  .col-xl-4 { flex: 0 0 33.333333%; }
  .col-xl-6 { flex: 0 0 50%; }
  .col-xl-8 { flex: 0 0 66.666667%; }
  .col-xl-9 { flex: 0 0 75%; }
  .col-xl-12 { flex: 0 0 100%; }
}
```

### 3. **头部导航设计**

```html
<!-- 头部导航 -->
<header class="header">
  <nav class="navbar" role="navigation" aria-label="主导航">
    <div class="container">
      <div class="navbar-brand">
        <a href="/" class="brand-link" aria-label="MovieInfo 首页">
          <img src="/static/images/logo.svg" alt="MovieInfo" class="brand-logo">
          <span class="brand-text">MovieInfo</span>
        </a>
      </div>
      
      <!-- 移动端菜单按钮 -->
      <button class="navbar-toggle" 
              type="button" 
              aria-label="切换导航菜单"
              aria-expanded="false"
              aria-controls="navbar-menu">
        <span class="hamburger-line"></span>
        <span class="hamburger-line"></span>
        <span class="hamburger-line"></span>
      </button>
      
      <!-- 导航菜单 -->
      <div id="navbar-menu" class="navbar-menu">
        <ul class="navbar-nav" role="menubar">
          <li class="nav-item" role="none">
            <a href="/" class="nav-link" role="menuitem">首页</a>
          </li>
          <li class="nav-item dropdown" role="none">
            <a href="/movies" 
               class="nav-link dropdown-toggle" 
               role="menuitem"
               aria-haspopup="true"
               aria-expanded="false">
              电影
            </a>
            <ul class="dropdown-menu" role="menu">
              <li role="none"><a href="/movies/popular" class="dropdown-link" role="menuitem">热门电影</a></li>
              <li role="none"><a href="/movies/latest" class="dropdown-link" role="menuitem">最新电影</a></li>
              <li role="none"><a href="/movies/top-rated" class="dropdown-link" role="menuitem">高分电影</a></li>
              <li role="none"><a href="/categories" class="dropdown-link" role="menuitem">电影分类</a></li>
            </ul>
          </li>
          <li class="nav-item" role="none">
            <a href="/discover" class="nav-link" role="menuitem">发现</a>
          </li>
        </ul>
        
        <!-- 搜索框 -->
        <form class="navbar-search" role="search" action="/search" method="GET">
          <div class="search-input-group">
            <input type="search" 
                   name="q" 
                   class="search-input" 
                   placeholder="搜索电影、演员、导演..."
                   aria-label="搜索"
                   autocomplete="off">
            <button type="submit" class="search-button" aria-label="搜索">
              <svg class="icon" viewBox="0 0 24 24">
                <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
              </svg>
            </button>
          </div>
        </form>
        
        <!-- 用户菜单 -->
        <div class="navbar-user">
          {{if .User}}
          <div class="user-dropdown">
            <button class="user-button" 
                    aria-label="用户菜单"
                    aria-haspopup="true"
                    aria-expanded="false">
              <img src="{{.User.Avatar}}" alt="{{.User.Username}}" class="user-avatar">
              <span class="user-name">{{.User.Username}}</span>
            </button>
            <ul class="user-menu" role="menu">
              <li role="none"><a href="/profile" class="menu-link" role="menuitem">个人资料</a></li>
              <li role="none"><a href="/favorites" class="menu-link" role="menuitem">我的收藏</a></li>
              <li role="none"><a href="/watchlist" class="menu-link" role="menuitem">观看列表</a></li>
              <li role="none"><a href="/settings" class="menu-link" role="menuitem">设置</a></li>
              <li class="menu-divider" role="separator"></li>
              <li role="none"><a href="/logout" class="menu-link" role="menuitem">退出登录</a></li>
            </ul>
          </div>
          {{else}}
          <div class="auth-buttons">
            <a href="/login" class="btn btn-outline">登录</a>
            <a href="/register" class="btn btn-primary">注册</a>
          </div>
          {{end}}
        </div>
        
        <!-- 主题切换 -->
        <button class="theme-toggle" 
                aria-label="切换主题"
                title="切换深色/浅色主题">
          <svg class="icon icon-light" viewBox="0 0 24 24">
            <path d="M12 2.25a.75.75 0 01.75.75v2.25a.75.75 0 01-1.5 0V3a.75.75 0 01.75-.75zM7.5 12a4.5 4.5 0 119 0 4.5 4.5 0 01-9 0zM18.894 6.166a.75.75 0 00-1.06-1.06l-1.591 1.59a.75.75 0 101.06 1.061l1.591-1.59zM21.75 12a.75.75 0 01-.75.75h-2.25a.75.75 0 010-1.5H21a.75.75 0 01.75.75zM17.834 18.894a.75.75 0 001.06-1.06l-1.59-1.591a.75.75 0 10-1.061 1.06l1.59 1.591zM12 18a.75.75 0 01.75.75V21a.75.75 0 01-1.5 0v-2.25A.75.75 0 0112 18zM7.758 17.303a.75.75 0 00-1.061-1.06l-1.591 1.59a.75.75 0 001.06 1.061l1.591-1.59zM6 12a.75.75 0 01-.75.75H3a.75.75 0 010-1.5h2.25A.75.75 0 016 12zM6.697 7.757a.75.75 0 001.06-1.06l-1.59-1.591a.75.75 0 00-1.061 1.06l1.59 1.591z"/>
          </svg>
          <svg class="icon icon-dark" viewBox="0 0 24 24">
            <path d="M9.528 1.718a.75.75 0 01.162.819A8.97 8.97 0 009 6a9 9 0 009 9 8.97 8.97 0 003.463-.69.75.75 0 01.981.98 10.503 10.503 0 01-9.694 6.46c-5.799 0-10.5-4.701-10.5-10.5 0-4.368 2.667-8.112 6.46-9.694a.75.75 0 01.818.162z"/>
          </svg>
        </button>
      </div>
    </div>
  </nav>
</header>
```

### 4. **底部设计**

```html
<!-- 页面底部 -->
<footer class="footer">
  <div class="container">
    <div class="footer-content">
      <div class="row">
        <!-- 品牌信息 -->
        <div class="col-lg-4 col-md-6">
          <div class="footer-section">
            <div class="footer-brand">
              <img src="/static/images/logo.svg" alt="MovieInfo" class="footer-logo">
              <h3 class="footer-title">MovieInfo</h3>
            </div>
            <p class="footer-description">
              专业的电影信息平台，为您提供最新的电影资讯、影评和推荐。
            </p>
            <div class="social-links">
              <a href="#" class="social-link" aria-label="微博">
                <svg class="icon" viewBox="0 0 24 24">
                  <!-- 微博图标 -->
                </svg>
              </a>
              <a href="#" class="social-link" aria-label="微信">
                <svg class="icon" viewBox="0 0 24 24">
                  <!-- 微信图标 -->
                </svg>
              </a>
              <a href="#" class="social-link" aria-label="QQ">
                <svg class="icon" viewBox="0 0 24 24">
                  <!-- QQ图标 -->
                </svg>
              </a>
            </div>
          </div>
        </div>
        
        <!-- 快速链接 -->
        <div class="col-lg-2 col-md-6">
          <div class="footer-section">
            <h4 class="footer-heading">电影</h4>
            <ul class="footer-links">
              <li><a href="/movies/popular">热门电影</a></li>
              <li><a href="/movies/latest">最新电影</a></li>
              <li><a href="/movies/top-rated">高分电影</a></li>
              <li><a href="/categories">电影分类</a></li>
            </ul>
          </div>
        </div>
        
        <div class="col-lg-2 col-md-6">
          <div class="footer-section">
            <h4 class="footer-heading">发现</h4>
            <ul class="footer-links">
              <li><a href="/discover/trending">趋势</a></li>
              <li><a href="/discover/awards">获奖电影</a></li>
              <li><a href="/discover/festivals">电影节</a></li>
              <li><a href="/discover/collections">合集</a></li>
            </ul>
          </div>
        </div>
        
        <div class="col-lg-2 col-md-6">
          <div class="footer-section">
            <h4 class="footer-heading">帮助</h4>
            <ul class="footer-links">
              <li><a href="/help">帮助中心</a></li>
              <li><a href="/contact">联系我们</a></li>
              <li><a href="/feedback">意见反馈</a></li>
              <li><a href="/api">API文档</a></li>
            </ul>
          </div>
        </div>
        
        <div class="col-lg-2 col-md-6">
          <div class="footer-section">
            <h4 class="footer-heading">关于</h4>
            <ul class="footer-links">
              <li><a href="/about">关于我们</a></li>
              <li><a href="/privacy">隐私政策</a></li>
              <li><a href="/terms">服务条款</a></li>
              <li><a href="/careers">加入我们</a></li>
            </ul>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 版权信息 -->
    <div class="footer-bottom">
      <div class="row align-items-center">
        <div class="col-md-6">
          <p class="copyright">
            © 2024 MovieInfo. 保留所有权利。
          </p>
        </div>
        <div class="col-md-6">
          <div class="footer-meta">
            <span class="meta-item">ICP备案号：京ICP备12345678号</span>
            <span class="meta-item">京公网安备 11010802012345号</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</footer>
```

## 📱 响应式设计

### 1. **断点系统**

```css
/* 断点定义 */
:root {
  --breakpoint-xs: 0;
  --breakpoint-sm: 576px;
  --breakpoint-md: 768px;
  --breakpoint-lg: 992px;
  --breakpoint-xl: 1200px;
  --breakpoint-xxl: 1400px;
}

/* 媒体查询混合 */
@media (max-width: 575.98px) {
  /* 超小屏幕 */
  .container { padding: 0 var(--spacing-3); }
  .navbar-menu { 
    position: fixed;
    top: 60px;
    left: 0;
    right: 0;
    background: var(--bg-primary);
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    transform: translateY(-100%);
    transition: transform 0.3s ease;
  }
  .navbar-menu.active { transform: translateY(0); }
}

@media (min-width: 576px) and (max-width: 767.98px) {
  /* 小屏幕 */
  .container { max-width: 540px; }
}

@media (min-width: 768px) and (max-width: 991.98px) {
  /* 中等屏幕 */
  .container { max-width: 720px; }
}

@media (min-width: 992px) and (max-width: 1199.98px) {
  /* 大屏幕 */
  .container { max-width: 960px; }
}

@media (min-width: 1200px) {
  /* 超大屏幕 */
  .container { max-width: 1140px; }
}
```

### 2. **移动端优化**

```css
/* 移动端优化 */
@media (max-width: 767.98px) {
  /* 触摸友好的按钮大小 */
  .btn {
    min-height: 44px;
    padding: var(--spacing-3) var(--spacing-4);
  }
  
  /* 更大的点击区域 */
  .nav-link {
    padding: var(--spacing-4);
    display: block;
  }
  
  /* 移动端字体调整 */
  .text-4xl { font-size: var(--font-size-3xl); }
  .text-3xl { font-size: var(--font-size-2xl); }
  
  /* 移动端间距调整 */
  .section { padding: var(--spacing-8) 0; }
  .hero { padding: var(--spacing-12) 0; }
  
  /* 隐藏桌面端元素 */
  .d-none-mobile { display: none !important; }
}

/* 桌面端优化 */
@media (min-width: 768px) {
  /* 隐藏移动端元素 */
  .d-none-desktop { display: none !important; }
  
  /* 悬停效果 */
  .btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  }
  
  .card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(0,0,0,0.15);
  }
}
```

## 🎨 组件样式

### 1. **按钮组件**

```css
/* 按钮基础样式 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-3) var(--spacing-5);
  border: 1px solid transparent;
  border-radius: 6px;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  line-height: 1.5;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.btn:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* 按钮变体 */
.btn-primary {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background-color: var(--primary-dark);
  border-color: var(--primary-dark);
}

.btn-outline {
  background-color: transparent;
  border-color: var(--gray-300);
  color: var(--text-primary);
}

.btn-outline:hover {
  background-color: var(--gray-50);
  border-color: var(--gray-400);
}

.btn-ghost {
  background-color: transparent;
  border-color: transparent;
  color: var(--text-secondary);
}

.btn-ghost:hover {
  background-color: var(--gray-100);
  color: var(--text-primary);
}

/* 按钮大小 */
.btn-sm {
  padding: var(--spacing-2) var(--spacing-3);
  font-size: var(--font-size-xs);
}

.btn-lg {
  padding: var(--spacing-4) var(--spacing-6);
  font-size: var(--font-size-base);
}
```

### 2. **卡片组件**

```css
/* 卡片基础样式 */
.card {
  background-color: var(--bg-primary);
  border: 1px solid var(--gray-200);
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  overflow: hidden;
  transition: all 0.3s ease;
}

.card-header {
  padding: var(--spacing-4) var(--spacing-5);
  border-bottom: 1px solid var(--gray-200);
  background-color: var(--bg-secondary);
}

.card-body {
  padding: var(--spacing-5);
}

.card-footer {
  padding: var(--spacing-4) var(--spacing-5);
  border-top: 1px solid var(--gray-200);
  background-color: var(--bg-secondary);
}

/* 电影卡片 */
.movie-card {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.movie-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 30px rgba(0,0,0,0.2);
}

.movie-poster {
  width: 100%;
  aspect-ratio: 2/3;
  object-fit: cover;
}

.movie-info {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0,0,0,0.8));
  color: white;
  padding: var(--spacing-4);
}

.movie-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-2);
}

.movie-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  font-size: var(--font-size-sm);
  opacity: 0.9;
}
```

## 📝 总结

页面布局设计为MovieInfo项目提供了完整的视觉框架：

**核心特性**：
1. **设计系统**：统一的色彩、字体、间距系统
2. **响应式布局**：移动优先的自适应设计
3. **组件化设计**：可复用的UI组件库
4. **无障碍支持**：符合WCAG标准的可访问性设计

**技术特性**：
- CSS自定义属性系统
- 灵活的网格布局
- 现代化的CSS特性
- 性能优化的加载策略

**用户体验**：
- 直观的导航结构
- 一致的视觉语言
- 流畅的交互动画
- 优雅的错误处理

下一步，我们将基于这个布局系统实现具体的首页设计。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第47步：首页开发](47-homepage.md)
