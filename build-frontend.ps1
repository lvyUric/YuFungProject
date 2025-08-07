# 前端Docker构建脚本
Write-Host "开始构建前端Docker镜像..." -ForegroundColor Green

# 设置工作目录
Set-Location Yufung-admin-front

# 检查package-lock.json是否存在
if (-not (Test-Path "package-lock.json")) {
    Write-Host "❌ package-lock.json 文件不存在，正在生成..." -ForegroundColor Yellow
    npm install
}

# 检查Docker状态
Write-Host "检查Docker状态..." -ForegroundColor Cyan
docker --version
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Docker未安装或未启动" -ForegroundColor Red
    exit 1
}

# 清理之前的构建
Write-Host "清理之前的前端构建..." -ForegroundColor Cyan
docker rmi yufung-frontend:latest -f 2>$null

# 开始构建
Write-Host "开始构建前端Docker镜像..." -ForegroundColor Green
docker build `
    --progress=plain `
    --no-cache `
    -t yufung-frontend:latest .

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 前端Docker镜像构建成功！" -ForegroundColor Green
    Write-Host "镜像名称: yufung-frontend:latest" -ForegroundColor Yellow
    
    # 显示镜像信息
    Write-Host "镜像信息:" -ForegroundColor Cyan
    docker images yufung-frontend:latest
    
    Write-Host ""
    Write-Host "运行容器命令:" -ForegroundColor Cyan
    Write-Host "docker run -d -p 80:80 --name yufung-frontend yufung-frontend:latest" -ForegroundColor White
} else {
    Write-Host "❌ 前端Docker镜像构建失败！" -ForegroundColor Red
    Write-Host "请检查错误信息并重试" -ForegroundColor Yellow
    exit 1
}

# 返回原目录
Set-Location .. 