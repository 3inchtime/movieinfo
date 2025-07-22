# 2.1 开发环境准备

## 概述

开发环境是软件开发的基础设施，它直接影响开发效率、代码质量和团队协作。对于MovieInfo项目，我们需要搭建一个统一、稳定、高效的开发环境，确保所有开发人员都能在相同的环境下工作。

## 为什么开发环境如此重要？

### 1. **环境一致性**
- **减少环境差异**：统一的开发环境避免"在我机器上能运行"的问题
- **降低调试成本**：相同环境下的问题更容易复现和解决
- **提高协作效率**：团队成员可以快速上手和相互支持

### 2. **开发效率**
- **工具链优化**：合适的工具提升开发速度
- **自动化支持**：自动化工具减少重复性工作
- **快速反馈**：实时编译、热重载等功能加快开发迭代

### 3. **质量保证**
- **代码规范**：统一的代码格式化和检查工具
- **测试环境**：便于单元测试和集成测试
- **版本控制**：规范的Git工作流程

### 4. **部署一致性**
- **容器化**：开发环境与生产环境保持一致
- **依赖管理**：明确的依赖版本控制
- **配置管理**：环境配置的标准化

## 系统要求

### 1. **硬件要求**

#### 1.1 最低配置
- **CPU**: 双核 2.0GHz 以上
- **内存**: 8GB RAM
- **存储**: 50GB 可用空间
- **网络**: 稳定的互联网连接

#### 1.2 推荐配置
- **CPU**: 四核 3.0GHz 以上（支持虚拟化）
- **内存**: 16GB RAM 或更多
- **存储**: SSD 100GB+ 可用空间
- **网络**: 高速宽带连接

**说明**：
- **CPU虚拟化支持**：用于Docker容器运行
- **充足内存**：支持IDE、数据库、多个服务同时运行
- **SSD存储**：提升编译和数据库访问速度

### 2. **操作系统支持**

#### 2.1 主要支持
- **macOS**: 10.15+ (Catalina及以上)
- **Windows**: Windows 10/11 (支持WSL2)
- **Linux**: Ubuntu 20.04+, CentOS 8+, Debian 10+

#### 2.2 推荐系统
- **开发推荐**: macOS 或 Linux
- **Windows用户**: 建议使用WSL2 + Ubuntu

**原因**：
- **Unix-like系统**：更好的命令行工具支持
- **包管理器**：便于安装开发工具
- **Docker支持**：原生Docker支持更稳定

## 核心工具安装

### 1. **包管理器安装**

#### 1.1 macOS - Homebrew
```bash
# 安装Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 验证安装
brew --version

# 更新Homebrew
brew update
```

**Homebrew的优势**：
- **简化安装**：一条命令安装复杂软件
- **依赖管理**：自动处理软件依赖关系
- **版本管理**：支持多版本软件共存
- **社区支持**：丰富的软件包生态

#### 1.2 Windows - Chocolatey
```powershell
# 以管理员身份运行PowerShell
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 验证安装
choco --version
```

#### 1.3 Linux - 系统包管理器
```bash
# Ubuntu/Debian
sudo apt update
sudo apt upgrade

# CentOS/RHEL
sudo yum update
# 或者 (CentOS 8+)
sudo dnf update
```

### 2. **Git版本控制**

#### 2.1 安装Git
```bash
# macOS
brew install git

# Windows
choco install git

# Ubuntu/Debian
sudo apt install git

# CentOS/RHEL
sudo yum install git
```

#### 2.2 Git配置
```bash
# 配置用户信息
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# 配置默认分支名
git config --global init.defaultBranch main

# 配置编辑器
git config --global core.editor "code --wait"  # VS Code
# 或者
git config --global core.editor "vim"          # Vim

# 配置换行符处理
# Windows
git config --global core.autocrlf true
# macOS/Linux
git config --global core.autocrlf input

# 配置别名
git config --global alias.st status
git config --global alias.co checkout
git config --global alias.br branch
git config --global alias.ci commit
git config --global alias.lg "log --oneline --graph --decorate --all"
```

#### 2.3 SSH密钥配置
```bash
# 生成SSH密钥
ssh-keygen -t ed25519 -C "your.email@example.com"

# 启动ssh-agent
eval "$(ssh-agent -s)"

# 添加密钥到ssh-agent
ssh-add ~/.ssh/id_ed25519

# 查看公钥（添加到GitHub/GitLab）
cat ~/.ssh/id_ed25519.pub
```

### 3. **终端工具**

#### 3.1 现代终端
```bash
# macOS - iTerm2
brew install --cask iterm2

# Windows - Windows Terminal
choco install microsoft-windows-terminal

# Linux - 通常系统自带，或安装Terminator
sudo apt install terminator  # Ubuntu
```

#### 3.2 Shell增强
```bash
# 安装Oh My Zsh (macOS/Linux)
sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

# 安装有用的插件
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

# 编辑 ~/.zshrc
plugins=(git zsh-autosuggestions zsh-syntax-highlighting)
```

**Shell增强的好处**：
- **自动补全**：命令和路径自动补全
- **语法高亮**：命令语法错误提示
- **历史建议**：基于历史命令的智能建议
- **主题美化**：更美观的命令行界面

### 4. **代码编辑器**

#### 4.1 Visual Studio Code
```bash
# macOS
brew install --cask visual-studio-code

# Windows
choco install vscode

# Ubuntu/Debian
wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
sudo install -o root -g root -m 644 packages.microsoft.gpg /etc/apt/trusted.gpg.d/
sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/trusted.gpg.d/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
sudo apt update
sudo apt install code
```

#### 4.2 VS Code扩展安装
```bash
# Go语言支持
code --install-extension golang.go

# Git支持
code --install-extension eamodio.gitlens

# Docker支持
code --install-extension ms-azuretools.vscode-docker

# YAML支持
code --install-extension redhat.vscode-yaml

# Markdown支持
code --install-extension yzhang.markdown-all-in-one

# 代码格式化
code --install-extension esbenp.prettier-vscode

# 主题
code --install-extension dracula-theme.theme-dracula
```

#### 4.3 VS Code配置
```json
// settings.json
{
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    },
    "files.autoSave": "afterDelay",
    "files.autoSaveDelay": 1000,
    "terminal.integrated.shell.osx": "/bin/zsh",
    "workbench.colorTheme": "Dracula"
}
```

## 开发工具链

### 1. **命令行工具**

#### 1.1 基础工具
```bash
# 文件操作增强
brew install tree      # 目录树显示
brew install bat       # cat的增强版本
brew install exa       # ls的增强版本
brew install fd        # find的增强版本
brew install ripgrep   # grep的增强版本

# 网络工具
brew install curl      # HTTP客户端
brew install wget      # 文件下载工具
brew install httpie    # 用户友好的HTTP客户端

# 系统监控
brew install htop      # 进程监控
brew install iotop     # IO监控
```

#### 1.2 开发专用工具
```bash
# JSON处理
brew install jq        # JSON查询和处理

# 数据库客户端
brew install mysql-client
brew install redis

# API测试
brew install --cask postman
# 或者命令行版本
brew install httpie
```

### 2. **容器化工具**

#### 2.1 Docker安装
```bash
# macOS
brew install --cask docker

# Windows
choco install docker-desktop

# Ubuntu
sudo apt update
sudo apt install apt-transport-https ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io
```

#### 2.2 Docker配置
```bash
# 启动Docker服务
sudo systemctl start docker
sudo systemctl enable docker

# 将用户添加到docker组（Linux）
sudo usermod -aG docker $USER

# 验证安装
docker --version
docker run hello-world
```

#### 2.3 Docker Compose
```bash
# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.12.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
```

### 3. **数据库工具**

#### 3.1 数据库客户端
```bash
# 命令行客户端
brew install mysql-client
brew install redis

# GUI客户端
brew install --cask sequel-pro      # MySQL (macOS)
brew install --cask redis-pro       # Redis (macOS)
# Windows
choco install mysql.workbench
choco install redis-desktop-manager
```

#### 3.2 数据库管理工具
```bash
# 数据库迁移工具
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 数据库文档生成
go install github.com/k1LoW/tbls@latest
```

## 环境变量配置

### 1. **Shell配置文件**

#### 1.1 Zsh配置 (~/.zshrc)
```bash
# Go环境变量
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# 项目相关
export MOVIEINFO_ENV=development
export MOVIEINFO_CONFIG_PATH=$HOME/movieinfo/configs

# 数据库连接
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=movieinfo

# Redis连接
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=

# JWT密钥
export JWT_SECRET=your-secret-key

# 别名设置
alias ll='exa -la'
alias cat='bat'
alias grep='rg'
alias find='fd'

# 项目快捷方式
alias cdmovie='cd $HOME/movieinfo'
alias runweb='cd $HOME/movieinfo && go run cmd/web/main.go'
alias runuser='cd $HOME/movieinfo && go run cmd/user/main.go'
```

#### 2.2 环境变量加载
```bash
# 重新加载配置
source ~/.zshrc

# 验证环境变量
echo $GOPATH
echo $MOVIEINFO_ENV
```

## 项目目录结构

### 1. **工作目录组织**
```bash
# 创建工作目录
mkdir -p ~/workspace/movieinfo
cd ~/workspace/movieinfo

# 创建项目相关目录
mkdir -p {docs,scripts,configs,logs}
```

### 2. **目录结构说明**
```
~/workspace/
├── movieinfo/           # 主项目目录
│   ├── cmd/            # 应用程序入口
│   ├── internal/       # 内部包
│   ├── pkg/            # 公共包
│   ├── configs/        # 配置文件
│   ├── docs/           # 项目文档
│   ├── scripts/        # 脚本文件
│   └── logs/           # 日志文件
├── tools/              # 开发工具
└── temp/               # 临时文件
```

## 开发流程工具

### 1. **代码质量工具**

#### 1.1 代码格式化
```bash
# 安装goimports
go install golang.org/x/tools/cmd/goimports@latest

# 安装golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1
```

#### 1.2 预提交钩子
```bash
# 安装pre-commit
pip install pre-commit

# 创建.pre-commit-config.yaml
cat > .pre-commit-config.yaml << EOF
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
  - repo: local
    hooks:
      - id: go-fmt
        name: go fmt
        entry: gofmt
        language: system
        args: [-w]
        files: \.go$
      - id: go-imports
        name: go imports
        entry: goimports
        language: system
        args: [-w]
        files: \.go$
EOF

# 安装钩子
pre-commit install
```

### 2. **调试工具**

#### 2.1 Delve调试器
```bash
# 安装Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 验证安装
dlv version
```

#### 2.2 性能分析工具
```bash
# pprof工具（Go内置）
go tool pprof

# 可视化工具
brew install graphviz
```

## 验证环境

### 1. **环境检查脚本**
```bash
#!/bin/bash
# check-env.sh

echo "=== 开发环境检查 ==="

# 检查Go
if command -v go &> /dev/null; then
    echo "✅ Go: $(go version)"
else
    echo "❌ Go 未安装"
fi

# 检查Git
if command -v git &> /dev/null; then
    echo "✅ Git: $(git --version)"
else
    echo "❌ Git 未安装"
fi

# 检查Docker
if command -v docker &> /dev/null; then
    echo "✅ Docker: $(docker --version)"
else
    echo "❌ Docker 未安装"
fi

# 检查数据库连接
if command -v mysql &> /dev/null; then
    echo "✅ MySQL客户端已安装"
else
    echo "❌ MySQL客户端未安装"
fi

# 检查环境变量
if [ -n "$GOPATH" ]; then
    echo "✅ GOPATH: $GOPATH"
else
    echo "❌ GOPATH 未设置"
fi

echo "=== 检查完成 ==="
```

### 2. **运行检查**
```bash
# 使脚本可执行
chmod +x check-env.sh

# 运行检查
./check-env.sh
```

## 总结

开发环境的搭建为MovieInfo项目的开发奠定了坚实的基础。通过统一的工具链和配置，我们确保了团队协作的一致性和开发效率的最大化。

**关键配置要点**：
1. **工具统一**：使用相同的开发工具和版本
2. **环境一致**：通过容器化保证环境一致性
3. **自动化**：使用脚本和钩子自动化重复任务
4. **质量保证**：集成代码质量检查工具

**环境优势**：
- **高效开发**：优化的工具链提升开发速度
- **质量保证**：自动化检查确保代码质量
- **团队协作**：统一环境减少协作摩擦
- **问题定位**：一致的环境便于问题复现和解决

**下一步**：基于这个开发环境，我们将进行Go语言环境的详细配置，包括Go模块管理、依赖管理和开发工具的深度集成。
