# 测试Docker构建脚本
Write-Host "开始测试Docker构建..." -ForegroundColor Green

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
docker rmi yufung-admin:latest -f 2>$null

# 开始构建
Write-Host "开始构建Docker镜像..." -ForegroundColor Green
docker build `
    --build-arg GOPROXY=https://goproxy.cn,direct `
    --build-arg GOSUMDB=sum.golang.google.cn `
    --build-arg GO111MODULE=on `
    --progress=plain `
    --no-cache `
    -t yufung-admin:latest .

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Docker镜像构建成功！" -ForegroundColor Green
    Write-Host "镜像名称: yufung-admin:latest" -ForegroundColor Yellow
    
    # 显示镜像信息
    Write-Host "镜像信息:" -ForegroundColor Cyan
    docker images yufung-admin:latest
    
    Write-Host ""
    Write-Host "运行容器命令:" -ForegroundColor Cyan
    Write-Host "docker run -d -p 8088:8088 --name yufung-admin yufung-admin:latest" -ForegroundColor White
} else {
    Write-Host "❌ Docker镜像构建失败！" -ForegroundColor Red
    Write-Host "请检查错误信息并重试" -ForegroundColor Yellow
    exit 1
} 