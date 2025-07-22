# 2.5 IDE 配置与插件

## 概述

IDE（集成开发环境）是开发者的主要工作平台，合适的IDE配置和插件能够显著提升开发效率、代码质量和开发体验。本文档将详细介绍VS Code的配置优化、必要插件安装和开发工作流程设置。

## 为什么IDE配置如此重要？

### 1. **开发效率提升**
- **智能补全**：代码自动补全和智能提示
- **快速导航**：文件、函数、变量的快速跳转
- **重构支持**：安全的代码重构和重命名
- **调试集成**：可视化调试和断点管理

### 2. **代码质量保证**
- **语法检查**：实时语法错误检测
- **代码规范**：自动格式化和风格检查
- **静态分析**：潜在问题的提前发现
- **测试集成**：单元测试的运行和覆盖率显示

### 3. **团队协作**
- **统一配置**：团队成员使用相同的开发环境
- **版本控制集成**：Git操作的可视化界面
- **代码审查**：内置的代码审查工具
- **文档生成**：自动生成和维护文档

### 4. **项目管理**
- **多项目支持**：同时管理多个项目
- **任务管理**：TODO、FIXME等任务标记
- **终端集成**：内置终端和命令执行
- **扩展生态**：丰富的插件生态系统

## VS Code 基础配置

### 1. **全局设置配置**

#### 1.1 用户设置 (settings.json)
```json
{
    // 编辑器基础配置
    "editor.fontSize": 14,
    "editor.fontFamily": "'JetBrains Mono', 'Fira Code', 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace",
    "editor.fontLigatures": true,
    "editor.lineHeight": 1.5,
    "editor.tabSize": 4,
    "editor.insertSpaces": true,
    "editor.detectIndentation": true,
    "editor.wordWrap": "on",
    "editor.rulers": [80, 120],
    
    // 代码格式化
    "editor.formatOnSave": true,
    "editor.formatOnPaste": true,
    "editor.formatOnType": false,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true,
        "source.fixAll": true
    },
    
    // 文件配置
    "files.autoSave": "afterDelay",
    "files.autoSaveDelay": 1000,
    "files.trimTrailingWhitespace": true,
    "files.insertFinalNewline": true,
    "files.trimFinalNewlines": true,
    "files.encoding": "utf8",
    "files.eol": "\n",
    
    // 搜索配置
    "search.exclude": {
        "**/node_modules": true,
        "**/bower_components": true,
        "**/*.code-search": true,
        "**/vendor": true,
        "**/tmp": true,
        "**/logs": true,
        "**/.git": true
    },
    
    // 文件监视配置
    "files.watcherExclude": {
        "**/.git/objects/**": true,
        "**/.git/subtree-cache/**": true,
        "**/node_modules/*/**": true,
        "**/tmp/**": true,
        "**/logs/**": true
    },
    
    // 终端配置
    "terminal.integrated.fontSize": 13,
    "terminal.integrated.fontFamily": "'JetBrains Mono', 'Fira Code', monospace",
    "terminal.integrated.shell.osx": "/bin/zsh",
    "terminal.integrated.shell.linux": "/bin/bash",
    "terminal.integrated.shell.windows": "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
    
    // 工作台配置
    "workbench.colorTheme": "Dracula",
    "workbench.iconTheme": "material-icon-theme",
    "workbench.startupEditor": "newUntitledFile",
    "workbench.editor.enablePreview": false,
    "workbench.editor.enablePreviewFromQuickOpen": false,
    
    // 资源管理器配置
    "explorer.confirmDelete": false,
    "explorer.confirmDragAndDrop": false,
    "explorer.openEditors.visible": 0,
    
    // 扩展配置
    "extensions.autoUpdate": true,
    "extensions.autoCheckUpdates": true,
    
    // 遥测配置
    "telemetry.enableTelemetry": false,
    "telemetry.enableCrashReporter": false
}
```

#### 1.2 Go语言特定配置
```json
{
    // Go语言配置
    "go.useLanguageServer": true,
    "go.languageServerExperimentalFeatures": {
        "diagnostics": true,
        "documentLink": true
    },
    
    // Go工具配置
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.lintFlags": ["--fast"],
    "go.vetFlags": ["-all"],
    
    // Go测试配置
    "go.testFlags": ["-v", "-race"],
    "go.testTimeout": "30s",
    "go.coverOnSave": true,
    "go.coverOnSaveTimeout": "30s",
    "go.coverageDecorator": {
        "type": "gutter",
        "coveredHighlightColor": "rgba(64,128,128,0.5)",
        "uncoveredHighlightColor": "rgba(128,64,64,0.25)",
        "coveredGutterStyle": "blockgreen",
        "uncoveredGutterStyle": "blockred"
    },
    
    // Go构建配置
    "go.buildOnSave": "package",
    "go.buildFlags": [],
    "go.buildTags": "",
    "go.installDependenciesWhenBuilding": true,
    
    // Go工具管理
    "go.toolsManagement.autoUpdate": true,
    "go.toolsManagement.checkForUpdates": "local",
    
    // gopls配置
    "gopls": {
        "analyses": {
            "unusedparams": true,
            "shadow": true,
            "fieldalignment": false,
            "nilness": true,
            "unusedwrite": true,
            "useany": true
        },
        "staticcheck": true,
        "gofumpt": true,
        "codelenses": {
            "gc_details": false,
            "generate": true,
            "regenerate_cgo": true,
            "test": true,
            "tidy": true,
            "upgrade_dependency": true,
            "vendor": true
        },
        "usePlaceholders": true,
        "completeUnimported": true,
        "matcher": "fuzzy",
        "symbolMatcher": "fuzzy",
        "deepCompletion": true
    }
}
```

### 2. **工作区配置**

#### 2.1 项目工作区设置 (.vscode/settings.json)
```json
{
    // 项目特定配置
    "go.gopath": "${workspaceFolder}",
    "go.goroot": "/usr/local/go",
    
    // 环境变量
    "go.toolsEnvVars": {
        "GOPROXY": "https://goproxy.cn,direct",
        "GOSUMDB": "sum.golang.org",
        "GOPRIVATE": "github.com/yourcompany/*"
    },
    
    // 项目构建配置
    "go.buildTags": "integration",
    "go.testTags": "unit,integration",
    
    // 文件关联
    "files.associations": {
        "*.tmpl": "html",
        "*.tpl": "html",
        "Dockerfile*": "dockerfile",
        "docker-compose*.yml": "dockercompose",
        "*.proto": "proto3"
    },
    
    // 项目特定的排除规则
    "files.exclude": {
        "**/bin": true,
        "**/tmp": true,
        "**/*.log": true,
        "**/coverage.out": true,
        "**/coverage.html": true
    }
}
```

#### 2.2 调试配置 (.vscode/launch.json)
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Web Service",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/web",
            "env": {
                "MOVIEINFO_ENV": "development",
                "DB_HOST": "localhost",
                "DB_PORT": "3306",
                "DB_USER": "movieinfo",
                "DB_PASSWORD": "movieinfo123",
                "DB_NAME": "movieinfo",
                "REDIS_HOST": "localhost",
                "REDIS_PORT": "6379",
                "REDIS_PASSWORD": "movieinfo_redis_password"
            },
            "args": [],
            "showLog": true,
            "trace": "verbose"
        },
        {
            "name": "Launch User Service",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/user",
            "env": {
                "MOVIEINFO_ENV": "development",
                "SERVICE_PORT": "8081",
                "GRPC_PORT": "9081"
            },
            "args": []
        },
        {
            "name": "Launch Movie Service",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/movie",
            "env": {
                "MOVIEINFO_ENV": "development",
                "SERVICE_PORT": "8082",
                "GRPC_PORT": "9082"
            },
            "args": []
        },
        {
            "name": "Launch Comment Service",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/comment",
            "env": {
                "MOVIEINFO_ENV": "development",
                "SERVICE_PORT": "8083",
                "GRPC_PORT": "9083"
            },
            "args": []
        },
        {
            "name": "Debug Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}",
            "args": [
                "-test.v",
                "-test.run",
                "TestFunctionName"
            ]
        },
        {
            "name": "Attach to Process",
            "type": "go",
            "request": "attach",
            "mode": "local",
            "processId": 0
        },
        {
            "name": "Connect to Remote Delve",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "${workspaceFolder}",
            "port": 2345,
            "host": "127.0.0.1"
        }
    ]
}
```

#### 2.3 任务配置 (.vscode/tasks.json)
```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "go: build all services",
            "type": "shell",
            "command": "make",
            "args": ["build"],
            "group": {
                "kind": "build",
                "isDefault": true
            },
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            },
            "problemMatcher": "$go"
        },
        {
            "label": "go: test all",
            "type": "shell",
            "command": "go",
            "args": ["test", "-v", "./..."],
            "group": "test",
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            },
            "problemMatcher": "$go"
        },
        {
            "label": "go: test with coverage",
            "type": "shell",
            "command": "go",
            "args": [
                "test",
                "-v",
                "-coverprofile=coverage.out",
                "./..."
            ],
            "group": "test",
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            },
            "problemMatcher": "$go"
        },
        {
            "label": "go: generate coverage report",
            "type": "shell",
            "command": "go",
            "args": [
                "tool",
                "cover",
                "-html=coverage.out",
                "-o",
                "coverage.html"
            ],
            "group": "test",
            "dependsOn": "go: test with coverage"
        },
        {
            "label": "golangci-lint",
            "type": "shell",
            "command": "golangci-lint",
            "args": ["run"],
            "group": "test",
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            },
            "problemMatcher": [
                {
                    "owner": "golangci-lint",
                    "fileLocation": "absolute",
                    "pattern": {
                        "regexp": "^(.+?):(\\d+):(\\d+):\\s+(warning|error):\\s+(.+?)\\s+\\((.+?)\\)$",
                        "file": 1,
                        "line": 2,
                        "column": 3,
                        "severity": 4,
                        "message": 5,
                        "code": 6
                    }
                }
            ]
        },
        {
            "label": "docker: build and up",
            "type": "shell",
            "command": "docker-compose",
            "args": ["up", "-d", "--build"],
            "group": "build",
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            }
        },
        {
            "label": "docker: down",
            "type": "shell",
            "command": "docker-compose",
            "args": ["down"],
            "group": "build",
            "presentation": {
                "echo": true,
                "reveal": "always",
                "focus": false,
                "panel": "shared"
            }
        }
    ]
}
```

## 必要插件安装

### 1. **Go语言开发插件**

#### 1.1 核心插件
```bash
# Go语言支持
code --install-extension golang.go

# Go语言增强
code --install-extension golang.go-nightly

# Protocol Buffers支持
code --install-extension zxh404.vscode-proto3
```

**Go插件功能**：
- 语法高亮和智能补全
- 代码导航和重构
- 调试支持
- 测试运行和覆盖率
- 代码格式化和检查
- 包管理和依赖分析

#### 1.2 代码质量插件
```bash
# 代码检查和格式化
code --install-extension ms-vscode.vscode-go

# 错误检查增强
code --install-extension streetsidesoftware.code-spell-checker

# TODO高亮
code --install-extension wayou.vscode-todo-highlight

# 括号匹配
code --install-extension coenraads.bracket-pair-colorizer-2

# 缩进指示
code --install-extension oderwat.indent-rainbow
```

### 2. **版本控制插件**

#### 2.1 Git相关插件
```bash
# Git增强
code --install-extension eamodio.gitlens

# Git历史
code --install-extension donjayamanne.githistory

# Git图形化
code --install-extension mhutchie.git-graph

# Git忽略文件
code --install-extension codezombiech.gitignore
```

**GitLens配置**：
```json
{
    "gitlens.advanced.messages": {
        "suppressCommitHasNoPreviousCommitWarning": false,
        "suppressCommitNotFoundWarning": false,
        "suppressFileNotUnderSourceControlWarning": false,
        "suppressGitVersionWarning": false,
        "suppressLineUncommittedWarning": false,
        "suppressNoRepositoryWarning": false
    },
    "gitlens.blame.avatars": true,
    "gitlens.blame.compact": false,
    "gitlens.blame.dateFormat": "MMMM Do, YYYY h:mma",
    "gitlens.blame.format": "${message|50?} ${agoOrDate|14-}",
    "gitlens.blame.heatmap.enabled": true,
    "gitlens.blame.highlight.enabled": true,
    "gitlens.codeLens.authors.enabled": true,
    "gitlens.codeLens.recentChange.enabled": true,
    "gitlens.currentLine.enabled": true,
    "gitlens.hovers.currentLine.over": "line",
    "gitlens.statusBar.enabled": true
}
```

### 3. **Docker和容器插件**

#### 3.1 容器开发插件
```bash
# Docker支持
code --install-extension ms-azuretools.vscode-docker

# Docker Compose支持
code --install-extension ms-vscode-remote.remote-containers

# Kubernetes支持
code --install-extension ms-kubernetes-tools.vscode-kubernetes-tools
```

**Docker插件配置**：
```json
{
    "docker.showStartPage": false,
    "docker.dockerPath": "docker",
    "docker.dockerComposePath": "docker-compose",
    "docker.attachShellCommand.linuxContainer": "/bin/sh",
    "docker.attachShellCommand.windowsContainer": "powershell"
}
```

### 4. **数据库插件**

#### 4.1 数据库管理插件
```bash
# MySQL支持
code --install-extension formulahendry.vscode-mysql

# Redis支持
code --install-extension cweijan.vscode-redis-client

# 数据库客户端
code --install-extension cweijan.vscode-database-client2

# SQL格式化
code --install-extension bradymholt.pgformatter
```

### 5. **Web开发插件**

#### 5.1 前端开发插件
```bash
# HTML/CSS/JS支持
code --install-extension ms-vscode.vscode-html-css-support
code --install-extension bradlc.vscode-tailwindcss
code --install-extension esbenp.prettier-vscode

# 模板引擎支持
code --install-extension wholroyd.jinja
code --install-extension karunamurti.tera

# REST客户端
code --install-extension humao.rest-client
```

### 6. **效率提升插件**

#### 6.1 通用效率插件
```bash
# 文件图标
code --install-extension pkief.material-icon-theme

# 主题
code --install-extension dracula-theme.theme-dracula

# 路径智能补全
code --install-extension christian-kohler.path-intellisense

# 自动重命名标签
code --install-extension formulahendry.auto-rename-tag

# 多光标编辑
code --install-extension sleistner.vscode-fileutils

# 书签管理
code --install-extension alefragnani.bookmarks

# 代码片段
code --install-extension ms-vscode.vscode-snippet

# 文件搜索增强
code --install-extension ms-vscode.vscode-search-result
```

#### 6.2 代码片段配置
```json
// go.json (用户代码片段)
{
    "HTTP Handler": {
        "prefix": "handler",
        "body": [
            "func ${1:HandlerName}(c *gin.Context) {",
            "\t// TODO: Implement handler logic",
            "\t$0",
            "\tc.JSON(http.StatusOK, gin.H{",
            "\t\t\"message\": \"success\",",
            "\t})",
            "}"
        ],
        "description": "Create a Gin HTTP handler"
    },
    "Struct with JSON tags": {
        "prefix": "struct",
        "body": [
            "type ${1:StructName} struct {",
            "\t${2:Field} ${3:Type} `json:\"${4:field}\" db:\"${4:field}\"`",
            "\t$0",
            "}"
        ],
        "description": "Create a struct with JSON and DB tags"
    },
    "Error handling": {
        "prefix": "iferr",
        "body": [
            "if err != nil {",
            "\t${1:return err}",
            "}"
        ],
        "description": "Basic error handling"
    },
    "Test function": {
        "prefix": "test",
        "body": [
            "func Test${1:FunctionName}(t *testing.T) {",
            "\t// Arrange",
            "\t$2",
            "",
            "\t// Act",
            "\t$3",
            "",
            "\t// Assert",
            "\tassert.Equal(t, ${4:expected}, ${5:actual})",
            "}"
        ],
        "description": "Create a test function"
    }
}
```

## 开发工作流程配置

### 1. **快捷键配置**

#### 1.1 自定义快捷键 (keybindings.json)
```json
[
    {
        "key": "cmd+shift+t",
        "command": "go.test.package"
    },
    {
        "key": "cmd+shift+r",
        "command": "go.test.file"
    },
    {
        "key": "cmd+shift+c",
        "command": "go.test.coverage.toggle"
    },
    {
        "key": "cmd+shift+b",
        "command": "workbench.action.tasks.build"
    },
    {
        "key": "cmd+shift+d",
        "command": "workbench.view.debug"
    },
    {
        "key": "cmd+shift+g",
        "command": "workbench.view.scm"
    },
    {
        "key": "cmd+shift+e",
        "command": "workbench.view.explorer"
    },
    {
        "key": "cmd+shift+f",
        "command": "workbench.action.findInFiles"
    },
    {
        "key": "cmd+k cmd+f",
        "command": "editor.action.formatDocument"
    },
    {
        "key": "cmd+k cmd+o",
        "command": "editor.action.organizeImports"
    }
]
```

### 2. **项目模板配置**

#### 2.1 项目模板文件
```bash
# 创建项目模板目录
mkdir -p ~/.vscode/templates/go-microservice

# 复制配置文件到模板
cp .vscode/* ~/.vscode/templates/go-microservice/
```

#### 2.2 新项目初始化脚本
```bash
#!/bin/bash
# init-project.sh

PROJECT_NAME=$1
if [ -z "$PROJECT_NAME" ]; then
    echo "Usage: $0 <project-name>"
    exit 1
fi

# 创建项目目录
mkdir -p $PROJECT_NAME
cd $PROJECT_NAME

# 复制模板配置
cp -r ~/.vscode/templates/go-microservice/.vscode .

# 初始化Go模块
go mod init $PROJECT_NAME

# 创建基础目录结构
mkdir -p {cmd/{web,user,movie,comment},internal/{config,models,handlers,services,middleware},pkg/{database,redis,grpc},proto,templates,static,configs,docs,scripts,logs}

# 创建基础文件
touch {README.md,.gitignore,.golangci.yml,.air.toml,Makefile,docker-compose.yml}

echo "Project $PROJECT_NAME initialized successfully!"
```

## 总结

IDE配置与插件为MovieInfo项目提供了完整的开发环境支持。通过合理的配置和插件选择，我们建立了一个高效、智能的开发工作台。

**关键配置要点**：
1. **智能开发**：代码补全、语法检查、智能重构
2. **调试支持**：可视化调试、多服务调试配置
3. **版本控制**：Git集成、代码审查、历史追踪
4. **效率工具**：快捷键、代码片段、任务自动化

**开发优势**：
- **高效编码**：智能补全和快速导航
- **质量保证**：实时检查和自动格式化
- **调试便利**：可视化调试和性能分析
- **团队协作**：统一的开发环境和工具

**下一步**：基于这个完整的开发环境，我们将开始项目初始化，包括项目结构创建、Go模块初始化和基础配置设置。
