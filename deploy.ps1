# 昱丰保险经纪管理系统 Docker 部署脚本 (PowerShell版本)
# 作者: AI Assistant
# 日期: $(Get-Date -Format "yyyy-MM-dd")

param(
    [switch]$SkipBuild,
    [switch]$SkipPull
)

# 颜色定义
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"

# 服务器IP地址
$ServerIP = "106.52.172.124"

# 日志函数
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Red
}

function Write-Step {
    param([string]$Message)
    Write-Host "[STEP] $Message" -ForegroundColor $Blue
}

# 检查Docker是否安装
function Test-Docker {
    Write-Step "检查Docker环境..."
    
    try {
        $dockerVersion = docker --version
        Write-Info "Docker版本: $dockerVersion"
    }
    catch {
        Write-Error "Docker未安装或未启动，请先安装Docker Desktop"
        exit 1
    }
    
    try {
        $composeVersion = docker-compose --version
        Write-Info "Docker Compose版本: $composeVersion"
    }
    catch {
        Write-Error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    }
    
    Write-Info "Docker环境检查通过"
}

# 创建必要的目录
function New-Directories {
    Write-Step "创建必要的目录..."
    
    $directories = @("logs", "uploads")
    foreach ($dir in $directories) {
        if (!(Test-Path $dir)) {
            New-Item -ItemType Directory -Path $dir | Out-Null
            Write-Info "创建目录: $dir"
        }
    }
    
    Write-Info "目录创建完成"
}

# 停止现有容器
function Stop-Containers {
    Write-Step "停止现有容器..."
    
    try {
        docker-compose down --remove-orphans
        Write-Info "现有容器已停止"
    }
    catch {
        Write-Warn "停止容器时出现警告，继续执行..."
    }
}

# 拉取最新代码（如果需要）
function Update-Code {
    if (!$SkipPull) {
        Write-Step "拉取最新代码..."
        try {
            git pull
            Write-Info "代码更新完成"
        }
        catch {
            Write-Warn "代码更新失败，使用本地代码继续..."
        }
    }
}

# 构建镜像
function Build-Images {
    if (!$SkipBuild) {
        Write-Step "构建Docker镜像..."
        try {
            docker-compose build --no-cache
            Write-Info "镜像构建完成"
        }
        catch {
            Write-Error "镜像构建失败"
            exit 1
        }
    }
    else {
        Write-Info "跳过镜像构建"
    }
}

# 启动服务
function Start-Services {
    Write-Step "启动服务..."
    try {
        docker-compose up -d
        Write-Info "服务启动完成"
    }
    catch {
        Write-Error "服务启动失败"
        exit 1
    }
}

# 检查服务状态
function Test-Services {
    Write-Step "检查服务状态..."
    Start-Sleep -Seconds 10
    
    # 检查后端服务
    try {
        $response = Invoke-WebRequest -Uri "http://${ServerIP}:8088/health" -TimeoutSec 5 -UseBasicParsing
        if ($response.StatusCode -eq 200) {
            Write-Info "✅ 后端服务运行正常"
        }
    }
    catch {
        Write-Warn "⚠️ 后端服务可能未完全启动，请稍后检查"
    }
    
    # 检查前端服务
    try {
        $response = Invoke-WebRequest -Uri "http://${ServerIP}/health" -TimeoutSec 5 -UseBasicParsing
        if ($response.StatusCode -eq 200) {
            Write-Info "✅ 前端服务运行正常"
        }
    }
    catch {
        Write-Warn "⚠️ 前端服务可能未完全启动，请稍后检查"
    }
}

# 显示服务信息
function Show-Info {
    Write-Step "部署完成！服务信息如下："
    Write-Host ""
    Write-Host "🌐 前端访问地址: http://${ServerIP}" -ForegroundColor $Green
    Write-Host "🔧 后端API地址: http://${ServerIP}:8088" -ForegroundColor $Green
    Write-Host "📖 API文档地址: http://${ServerIP}:8088/swagger/index.html" -ForegroundColor $Green
    Write-Host "💾 数据库地址: mongodb://${ServerIP}:27017" -ForegroundColor $Green
    Write-Host ""
    Write-Host "📋 常用命令：" -ForegroundColor $Blue
    Write-Host "  查看服务状态: docker-compose ps"
    Write-Host "  查看服务日志: docker-compose logs -f"
    Write-Host "  停止服务: docker-compose down"
    Write-Host "  重启服务: docker-compose restart"
    Write-Host ""
}

# 主函数
function Main {
    Write-Host "🚀 开始部署昱丰保险经纪管理系统..." -ForegroundColor $Blue
    Write-Host ""
    
    Test-Docker
    New-Directories
    Update-Code
    Stop-Containers
    Build-Images
    Start-Services
    Test-Services
    Show-Info
    
    Write-Info "部署完成！"
}

# 执行主函数
try {
    Main
}
catch {
    Write-Error "部署过程中发生错误: $($_.Exception.Message)"
    exit 1
} 