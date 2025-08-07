# 完整的前后端Docker构建脚本
Write-Host "开始构建完整的前后端Docker镜像..." -ForegroundColor Green

# 设置环境变量
$env:DOCKER_BUILDKIT = "1"
$env:GOPROXY = "https://goproxy.cn,direct"

Write-Host "环境变量设置完成" -ForegroundColor Yellow

# 检查Docker状态
Write-Host "检查Docker状态..." -ForegroundColor Cyan
docker --version
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Docker未安装或未启动" -ForegroundColor Red
    exit 1
}

# 清理之前的构建
Write-Host "清理之前的构建..." -ForegroundColor Cyan
docker rmi yufung-backend:latest -f 2>$null
docker rmi yufung-frontend:latest -f 2>$null

# 构建后端
Write-Host "开始构建后端Docker镜像..." -ForegroundColor Green
docker build `
    --build-arg GOPROXY=https://goproxy.cn,direct `
    --build-arg GOSUMDB=sum.golang.google.cn `
    --build-arg GO111MODULE=on `
    --progress=plain `
    --no-cache `
    -t yufung-backend:latest .

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 后端Docker镜像构建成功！" -ForegroundColor Green
} else {
    Write-Host "❌ 后端Docker镜像构建失败！" -ForegroundColor Red
    exit 1
}

# 构建前端
Write-Host "开始构建前端Docker镜像..." -ForegroundColor Green
Set-Location Yufung-admin-front

# 检查package-lock.json是否存在
if (-not (Test-Path "package-lock.json")) {
    Write-Host "❌ package-lock.json 文件不存在，正在生成..." -ForegroundColor Yellow
    npm install
}

docker build `
    --progress=plain `
    --no-cache `
    -t yufung-frontend:latest .

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 前端Docker镜像构建成功！" -ForegroundColor Green
} else {
    Write-Host "❌ 前端Docker镜像构建失败！" -ForegroundColor Red
    Set-Location ..
    exit 1
}

# 返回原目录
Set-Location ..

# 显示所有镜像信息
Write-Host ""
Write-Host "所有镜像信息:" -ForegroundColor Cyan
docker images | findstr yufung

Write-Host ""
Write-Host "运行完整应用命令:" -ForegroundColor Cyan
Write-Host "docker-compose up -d" -ForegroundColor White
Write-Host ""
Write-Host "或者单独运行:" -ForegroundColor Cyan
Write-Host "docker run -d -p 8088:8088 --name yufung-backend yufung-backend:latest" -ForegroundColor White
Write-Host "docker run -d -p 80:80 --name yufung-frontend yufung-frontend:latest" -ForegroundColor White 