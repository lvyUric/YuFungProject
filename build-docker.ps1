# 设置Docker构建参数
$env:DOCKER_BUILDKIT = "1"
$env:BUILDKIT_PROGRESS = "plain"

# 设置Go代理环境变量
$env:GOPROXY = "https://goproxy.cn,direct"
$env:GOSUMDB = "sum.golang.google.cn"
$env:GO111MODULE = "on"

Write-Host "开始构建Docker镜像..." -ForegroundColor Green

# 构建镜像，使用BuildKit和缓存
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
    Write-Host ""
    Write-Host "运行容器命令:" -ForegroundColor Cyan
    Write-Host "docker run -d -p 8088:8088 --name yufung-admin yufung-admin:latest" -ForegroundColor White
} else {
    Write-Host "❌ Docker镜像构建失败！" -ForegroundColor Red
    exit 1
} 