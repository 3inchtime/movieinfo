# 第47步：首页开发

## 📋 概述

首页是MovieInfo项目的门户页面，承担着展示平台核心内容、引导用户探索和提升用户参与度的重要职责。一个优秀的首页需要平衡信息展示、用户体验和性能优化。

## 🎯 设计目标

### 1. **内容展示**
- 精选电影推荐
- 热门内容聚合
- 分类导航引导
- 实时数据展示

### 2. **用户引导**
- 清晰的功能入口
- 个性化推荐
- 搜索功能突出
- 注册转化优化

### 3. **性能体验**
- 快速首屏渲染
- 渐进式内容加载
- 图片懒加载
- 缓存策略优化

## 🏗️ 首页结构设计

### 1. **页面布局架构**

```html
<!-- 首页模板 -->
{{define "pages/home"}}
{{template "layouts/base" .}}

{{define "title"}}MovieInfo - 专业的电影信息平台{{end}}

{{define "meta"}}
<meta name="description" content="MovieInfo提供最新的电影信息、影评、评分和推荐，帮您发现好电影">
<meta name="keywords" content="电影,影评,评分,推荐,热门电影,最新电影">
<meta property="og:title" content="MovieInfo - 专业的电影信息平台">
<meta property="og:description" content="发现好电影，分享观影体验">
<meta property="og:image" content="/static/images/og-home.jpg">
<meta property="og:type" content="website">
{{end}}

{{define "css"}}
<link rel="stylesheet" href="{{asset "css/home.css"}}">
{{end}}

{{define "content"}}
<div class="homepage">
  <!-- 英雄区域 -->
  <section class="hero-section">
    {{template "components/hero-banner" .Data.FeaturedMovies}}
  </section>
  
  <!-- 快速导航 -->
  <section class="quick-nav-section">
    {{template "components/quick-navigation" .Data.Categories}}
  </section>
  
  <!-- 热门电影 -->
  <section class="popular-section">
    {{template "components/movie-section" dict "title" "热门电影" "movies" .Data.PopularMovies "viewAllUrl" "/movies/popular"}}
  </section>
  
  <!-- 最新电影 -->
  <section class="latest-section">
    {{template "components/movie-section" dict "title" "最新电影" "movies" .Data.LatestMovies "viewAllUrl" "/movies/latest"}}
  </section>
  
  <!-- 高分电影 -->
  <section class="top-rated-section">
    {{template "components/movie-section" dict "title" "高分电影" "movies" .Data.TopRatedMovies "viewAllUrl" "/movies/top-rated"}}
  </section>
  
  <!-- 统计信息 -->
  <section class="stats-section">
    {{template "components/platform-stats" .Data.Stats}}
  </section>
  
  <!-- 最新评论 -->
  <section class="recent-reviews-section">
    {{template "components/recent-reviews" .Data.RecentReviews}}
  </section>
</div>
{{end}}

{{define "js"}}
<script src="{{asset "js/home.js"}}"></script>
{{end}}
{{end}}
```

### 2. **英雄横幅组件**

```html
<!-- components/hero-banner.html -->
{{define "components/hero-banner"}}
<div class="hero-banner">
  <div class="hero-slider" id="heroSlider">
    {{range $index, $movie := .}}
    <div class="hero-slide {{if eq $index 0}}active{{end}}" data-slide="{{$index}}">
      <div class="hero-background">
        <img src="{{$movie.BackdropURL}}" 
             alt="{{$movie.Title}}" 
             class="hero-bg-image"
             loading="{{if eq $index 0}}eager{{else}}lazy{{end}}">
        <div class="hero-overlay"></div>
      </div>
      
      <div class="container">
        <div class="hero-content">
          <div class="row align-items-center">
            <div class="col-lg-8">
              <div class="hero-info">
                <h1 class="hero-title">{{$movie.Title}}</h1>
                <p class="hero-subtitle">{{$movie.Tagline}}</p>
                <p class="hero-description">{{truncate $movie.Overview 200}}</p>
                
                <div class="hero-meta">
                  <div class="meta-item">
                    <span class="meta-label">评分</span>
                    <div class="rating-display">
                      <span class="rating-value">{{formatRating $movie.Rating}}</span>
                      <div class="rating-stars">
                        {{range $i := seq 1 5}}
                        <svg class="star {{if le $i (div $movie.Rating 2)}}filled{{end}}" viewBox="0 0 24 24">
                          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                        </svg>
                        {{end}}
                      </div>
                    </div>
                  </div>
                  
                  <div class="meta-item">
                    <span class="meta-label">类型</span>
                    <div class="genre-tags">
                      {{range $movie.Genres}}
                      <span class="genre-tag">{{.Name}}</span>
                      {{end}}
                    </div>
                  </div>
                  
                  <div class="meta-item">
                    <span class="meta-label">上映时间</span>
                    <span class="meta-value">{{formatTime $movie.ReleaseDate "2006年1月2日"}}</span>
                  </div>
                </div>
                
                <div class="hero-actions">
                  <a href="/movies/{{$movie.ID}}" class="btn btn-primary btn-lg">
                    <svg class="icon" viewBox="0 0 24 24">
                      <path d="M8 5v14l11-7z"/>
                    </svg>
                    查看详情
                  </a>
                  
                  {{if $.User}}
                  <button class="btn btn-outline btn-lg" 
                          data-action="add-to-watchlist" 
                          data-movie-id="{{$movie.ID}}">
                    <svg class="icon" viewBox="0 0 24 24">
                      <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                    </svg>
                    加入观看列表
                  </button>
                  {{end}}
                  
                  <button class="btn btn-ghost btn-lg" 
                          data-action="share" 
                          data-movie-id="{{$movie.ID}}">
                    <svg class="icon" viewBox="0 0 24 24">
                      <path d="M18 16.08c-.76 0-1.44.3-1.96.77L8.91 12.7c.05-.23.09-.46.09-.7s-.04-.47-.09-.7l7.05-4.11c.54.5 1.25.81 2.04.81 1.66 0 3-1.34 3-3s-1.34-3-3-3-3 1.34-3 3c0 .24.04.47.09.7L8.04 9.81C7.5 9.31 6.79 9 6 9c-1.66 0-3 1.34-3 3s1.34 3 3 3c.79 0 1.50-.31 2.04-.81l7.12 4.16c-.05.21-.08.43-.08.65 0 1.61 1.31 2.92 2.92 2.92s2.92-1.31 2.92-2.92-1.31-2.92-2.92-2.92z"/>
                    </svg>
                    分享
                  </button>
                </div>
              </div>
            </div>
            
            <div class="col-lg-4">
              <div class="hero-poster">
                <img src="{{$movie.PosterURL}}" 
                     alt="{{$movie.Title}} 海报" 
                     class="poster-image"
                     loading="{{if eq $index 0}}eager{{else}}lazy{{end}}">
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    {{end}}
  </div>
  
  <!-- 轮播控制 -->
  <div class="hero-controls">
    <div class="container">
      <div class="slider-pagination">
        {{range $index, $movie := .}}
        <button class="pagination-dot {{if eq $index 0}}active{{end}}" 
                data-slide="{{$index}}"
                aria-label="切换到第{{add $index 1}}张幻灯片">
        </button>
        {{end}}
      </div>
      
      <div class="slider-navigation">
        <button class="nav-button nav-prev" aria-label="上一张">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z"/>
          </svg>
        </button>
        <button class="nav-button nav-next" aria-label="下一张">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</div>
{{end}}
```

### 3. **快速导航组件**

```html
<!-- components/quick-navigation.html -->
{{define "components/quick-navigation"}}
<div class="quick-navigation">
  <div class="container">
    <div class="quick-nav-grid">
      <a href="/movies/popular" class="quick-nav-item">
        <div class="nav-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        </div>
        <div class="nav-content">
          <h3 class="nav-title">热门电影</h3>
          <p class="nav-description">当前最受欢迎的电影</p>
        </div>
      </a>
      
      <a href="/movies/latest" class="quick-nav-item">
        <div class="nav-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        </div>
        <div class="nav-content">
          <h3 class="nav-title">最新电影</h3>
          <p class="nav-description">刚刚上映的新片</p>
        </div>
      </a>
      
      <a href="/movies/top-rated" class="quick-nav-item">
        <div class="nav-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        </div>
        <div class="nav-content">
          <h3 class="nav-title">高分电影</h3>
          <p class="nav-description">评分最高的佳作</p>
        </div>
      </a>
      
      <a href="/categories" class="quick-nav-item">
        <div class="nav-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M4 6H2v14c0 1.1.9 2 2 2h14v-2H4V6zm16-4H8c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-1 9H9V9h10v2zm-4 4H9v-2h6v2zm4-8H9V5h10v2z"/>
          </svg>
        </div>
        <div class="nav-content">
          <h3 class="nav-title">电影分类</h3>
          <p class="nav-description">按类型浏览电影</p>
        </div>
      </a>
      
      {{if .User}}
      <a href="/user/favorites" class="quick-nav-item">
        <div class="nav-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
          </svg>
        </div>
        <div class="nav-content">
          <h3 class="nav-title">我的收藏</h3>
          <p class="nav-description">收藏的电影列表</p>
        </div>
      </a>
      
      <a href="/user/watchlist" class="quick-nav-item">
        <div class="nav-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 2 2h8c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-5 14l-3-3 1.41-1.41L9 13.17l7.59-7.59L18 7l-9 9z"/>
          </svg>
        </div>
        <div class="nav-content">
          <h3 class="nav-title">观看列表</h3>
          <p class="nav-description">计划观看的电影</p>
        </div>
      </a>
      {{else}}
      <a href="/register" class="quick-nav-item quick-nav-cta">
        <div class="nav-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        </div>
        <div class="nav-content">
          <h3 class="nav-title">立即注册</h3>
          <p class="nav-description">加入MovieInfo社区</p>
        </div>
      </a>
      {{end}}
    </div>
  </div>
</div>
{{end}}
```

### 4. **电影区块组件**

```html
<!-- components/movie-section.html -->
{{define "components/movie-section"}}
<div class="movie-section">
  <div class="container">
    <div class="section-header">
      <h2 class="section-title">{{.title}}</h2>
      <a href="{{.viewAllUrl}}" class="view-all-link">
        查看全部
        <svg class="icon" viewBox="0 0 24 24">
          <path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"/>
        </svg>
      </a>
    </div>
    
    <div class="movie-grid" data-component="movie-grid">
      <div class="movie-slider" id="movieSlider-{{.title}}">
        {{range .movies}}
        <div class="movie-card-wrapper">
          <article class="movie-card" data-movie-id="{{.ID}}">
            <div class="movie-poster-container">
              <img src="{{.PosterURL}}" 
                   alt="{{.Title}} 海报" 
                   class="movie-poster"
                   loading="lazy"
                   onerror="this.src='/static/images/poster-placeholder.jpg'">
              
              <div class="movie-overlay">
                <div class="overlay-content">
                  <div class="movie-rating">
                    <svg class="icon star" viewBox="0 0 24 24">
                      <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                    </svg>
                    <span>{{formatRating .Rating}}</span>
                  </div>
                  
                  <div class="movie-actions">
                    <a href="/movies/{{.ID}}" class="action-btn primary">
                      <svg class="icon" viewBox="0 0 24 24">
                        <path d="M8 5v14l11-7z"/>
                      </svg>
                      详情
                    </a>
                    
                    {{if $.User}}
                    <button class="action-btn secondary" 
                            data-action="toggle-favorite" 
                            data-movie-id="{{.ID}}"
                            title="{{if .IsFavorite}}取消收藏{{else}}加入收藏{{end}}">
                      <svg class="icon" viewBox="0 0 24 24">
                        <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                      </svg>
                    </button>
                    
                    <button class="action-btn secondary" 
                            data-action="toggle-watchlist" 
                            data-movie-id="{{.ID}}"
                            title="{{if .IsInWatchlist}}从观看列表移除{{else}}加入观看列表{{end}}">
                      <svg class="icon" viewBox="0 0 24 24">
                        <path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 2 2h8c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-5 14l-3-3 1.41-1.41L9 13.17l7.59-7.59L18 7l-9 9z"/>
                      </svg>
                    </button>
                    {{end}}
                  </div>
                </div>
              </div>
            </div>
            
            <div class="movie-info">
              <h3 class="movie-title">
                <a href="/movies/{{.ID}}">{{.Title}}</a>
              </h3>
              <div class="movie-meta">
                <span class="release-year">{{formatTime .ReleaseDate "2006"}}</span>
                <span class="genre">{{index .Genres 0}}</span>
              </div>
            </div>
          </article>
        </div>
        {{end}}
      </div>
      
      <!-- 滑动控制 -->
      <button class="slider-nav slider-prev" 
              data-target="movieSlider-{{.title}}"
              aria-label="上一页">
        <svg class="icon" viewBox="0 0 24 24">
          <path d="M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z"/>
        </svg>
      </button>
      
      <button class="slider-nav slider-next" 
              data-target="movieSlider-{{.title}}"
              aria-label="下一页">
        <svg class="icon" viewBox="0 0 24 24">
          <path d="M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"/>
        </svg>
      </button>
    </div>
  </div>
</div>
{{end}}
```

### 5. **平台统计组件**

```html
<!-- components/platform-stats.html -->
{{define "components/platform-stats"}}
<div class="platform-stats">
  <div class="container">
    <div class="stats-grid">
      <div class="stat-item">
        <div class="stat-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4z"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-number" data-count="{{.TotalMovies}}">0</div>
          <div class="stat-label">部电影</div>
        </div>
      </div>
      
      <div class="stat-item">
        <div class="stat-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M16 4c0-1.11.89-2 2-2s2 .89 2 2-.89 2-2 2-2-.89-2-2zm4 18v-6h2.5l-2.54-7.63A1.5 1.5 0 0 0 18.5 8H16c-.8 0-1.5.7-1.5 1.5v6c0 .8.7 1.5 1.5 1.5h1v5h2zm-3.5-10.5c0-.8-.7-1.5-1.5-1.5H9c-.8 0-1.5.7-1.5 1.5v6c0 .8.7 1.5 1.5 1.5h1v5h2v-5h1c.8 0 1.5-.7 1.5-1.5v-6z"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-number" data-count="{{.TotalUsers}}">0</div>
          <div class="stat-label">位用户</div>
        </div>
      </div>
      
      <div class="stat-item">
        <div class="stat-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M20 2H4c-1.1 0-1.99.9-1.99 2L2 22l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-7 12h-2v-2h2v2zm0-4h-2V6h2v4z"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-number" data-count="{{.TotalReviews}}">0</div>
          <div class="stat-label">条评论</div>
        </div>
      </div>
      
      <div class="stat-item">
        <div class="stat-icon">
          <svg class="icon" viewBox="0 0 24 24">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-number" data-count="{{.TotalRatings}}">0</div>
          <div class="stat-label">次评分</div>
        </div>
      </div>
    </div>
  </div>
</div>
{{end}}
```

## 🎨 首页样式

### 1. **英雄区域样式**

```css
/* 英雄横幅 */
.hero-banner {
  position: relative;
  height: 70vh;
  min-height: 500px;
  overflow: hidden;
}

.hero-slide {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  transition: opacity 0.8s ease-in-out;
}

.hero-slide.active {
  opacity: 1;
}

.hero-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.hero-bg-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.hero-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    135deg,
    rgba(0,0,0,0.7) 0%,
    rgba(0,0,0,0.3) 50%,
    rgba(0,0,0,0.8) 100%
  );
}

.hero-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  align-items: center;
  color: white;
}

.hero-title {
  font-size: var(--font-size-5xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--spacing-4);
  text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
}

.hero-subtitle {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-medium);
  margin-bottom: var(--spacing-4);
  opacity: 0.9;
}

.hero-description {
  font-size: var(--font-size-lg);
  line-height: var(--line-height-relaxed);
  margin-bottom: var(--spacing-6);
  opacity: 0.8;
}

.hero-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-6);
  margin-bottom: var(--spacing-8);
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
}

.meta-label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  opacity: 0.7;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.rating-display {
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
}

.rating-value {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
}

.rating-stars {
  display: flex;
  gap: 2px;
}

.star {
  width: 16px;
  height: 16px;
  fill: var(--gray-400);
}

.star.filled {
  fill: var(--accent-color);
}

.genre-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-2);
}

.genre-tag {
  padding: var(--spacing-1) var(--spacing-3);
  background: rgba(255,255,255,0.2);
  border-radius: 20px;
  font-size: var(--font-size-sm);
  backdrop-filter: blur(10px);
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-4);
}

.hero-poster {
  text-align: center;
}

.poster-image {
  max-width: 300px;
  width: 100%;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.3);
}

/* 轮播控制 */
.hero-controls {
  position: absolute;
  bottom: var(--spacing-6);
  left: 0;
  right: 0;
  z-index: 3;
}

.slider-pagination {
  display: flex;
  justify-content: center;
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-4);
}

.pagination-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid rgba(255,255,255,0.5);
  background: transparent;
  cursor: pointer;
  transition: all 0.3s ease;
}

.pagination-dot.active {
  background: white;
  border-color: white;
}

.slider-navigation {
  display: flex;
  justify-content: space-between;
}

.nav-button {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(255,255,255,0.2);
  border: none;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.nav-button:hover {
  background: rgba(255,255,255,0.3);
  transform: scale(1.1);
}
```

### 2. **电影卡片样式**

```css
/* 电影网格 */
.movie-grid {
  position: relative;
}

.movie-slider {
  display: flex;
  gap: var(--spacing-4);
  overflow-x: auto;
  scroll-behavior: smooth;
  padding: var(--spacing-4) 0;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.movie-slider::-webkit-scrollbar {
  display: none;
}

.movie-card-wrapper {
  flex: 0 0 auto;
  width: 200px;
}

.movie-card {
  background: var(--bg-primary);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transition: all 0.3s ease;
  cursor: pointer;
}

.movie-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 30px rgba(0,0,0,0.2);
}

.movie-poster-container {
  position: relative;
  aspect-ratio: 2/3;
  overflow: hidden;
}

.movie-poster {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.movie-card:hover .movie-poster {
  transform: scale(1.05);
}

.movie-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    to bottom,
    transparent 0%,
    transparent 60%,
    rgba(0,0,0,0.8) 100%
  );
  opacity: 0;
  transition: opacity 0.3s ease;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: var(--spacing-4);
}

.movie-card:hover .movie-overlay {
  opacity: 1;
}

.movie-rating {
  display: flex;
  align-items: center;
  gap: var(--spacing-1);
  color: white;
  font-weight: var(--font-weight-medium);
  align-self: flex-end;
}

.movie-actions {
  display: flex;
  gap: var(--spacing-2);
  justify-content: center;
}

.action-btn {
  padding: var(--spacing-2);
  border-radius: 50%;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
}

.action-btn.primary {
  background: var(--primary-color);
  color: white;
}

.action-btn.secondary {
  background: rgba(255,255,255,0.2);
  color: white;
  backdrop-filter: blur(10px);
}

.action-btn:hover {
  transform: scale(1.1);
}

.movie-info {
  padding: var(--spacing-4);
}

.movie-title {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--spacing-2);
  line-height: 1.3;
}

.movie-title a {
  color: var(--text-primary);
  text-decoration: none;
}

.movie-title a:hover {
  color: var(--primary-color);
}

.movie-meta {
  display: flex;
  justify-content: space-between;
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

/* 滑动导航 */
.slider-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--bg-primary);
  border: 1px solid var(--gray-200);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 2;
}

.slider-nav:hover {
  background: var(--primary-color);
  color: white;
  transform: translateY(-50%) scale(1.1);
}

.slider-prev {
  left: -20px;
}

.slider-next {
  right: -20px;
}
```

## 📱 响应式适配

### 1. **移动端优化**

```css
/* 移动端适配 */
@media (max-width: 767.98px) {
  .hero-banner {
    height: 50vh;
    min-height: 400px;
  }
  
  .hero-title {
    font-size: var(--font-size-3xl);
  }
  
  .hero-subtitle {
    font-size: var(--font-size-lg);
  }
  
  .hero-meta {
    flex-direction: column;
    gap: var(--spacing-4);
  }
  
  .hero-actions {
    flex-direction: column;
  }
  
  .hero-poster {
    margin-top: var(--spacing-6);
  }
  
  .poster-image {
    max-width: 200px;
  }
  
  .quick-nav-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-3);
  }
  
  .movie-card-wrapper {
    width: 150px;
  }
  
  .slider-nav {
    display: none;
  }
}

@media (min-width: 768px) and (max-width: 991.98px) {
  .hero-title {
    font-size: var(--font-size-4xl);
  }
  
  .quick-nav-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .movie-card-wrapper {
    width: 180px;
  }
}
```

## 📝 总结

首页开发为MovieInfo项目提供了完整的门户体验：

**核心功能**：
1. **英雄横幅**：精选电影的视觉冲击展示
2. **快速导航**：直观的功能入口和分类导航
3. **内容聚合**：热门、最新、高分电影的有序展示
4. **用户引导**：个性化推荐和注册转化

**技术特性**：
- 响应式设计适配
- 图片懒加载优化
- 交互动画效果
- 无障碍访问支持

**用户体验**：
- 快速的首屏加载
- 流畅的滑动交互
- 清晰的信息层级
- 一致的视觉语言

下一步，我们将实现电影列表页，为用户提供完整的电影浏览体验。

---

**文档状态**: ✅ 已完成  
**最后更新**: 2025-07-22  
**下一步**: [第48步：电影列表页](48-movie-list.md)
