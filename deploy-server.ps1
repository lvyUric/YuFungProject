# 服务器部署脚本
Write-Host "开始服务器部署..." -ForegroundColor Green

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

# 停止并删除现有容器
Write-Host "停止现有容器..." -ForegroundColor Cyan
docker stop yufung-backend yufung-frontend 2>$null
docker rm yufung-backend yufung-frontend 2>$null

# 删除旧镜像
Write-Host "清理旧镜像..." -ForegroundColor Cyan
docker rmi yufung-backend:latest yufung-frontend:latest 2>$null

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

# 运行容器
Write-Host "启动应用容器..." -ForegroundColor Green

# 启动后端容器
docker run -d `
    --name yufung-backend `
    --restart unless-stopped `
    -p 8088:8088 `
    -v ${PWD}/logs:/app/logs `
    -v ${PWD}/uploads:/app/uploads `
    -e TZ=Asia/Shanghai `
    yufung-backend:latest

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 后端容器启动成功！" -ForegroundColor Green
} else {
    Write-Host "❌ 后端容器启动失败！" -ForegroundColor Red
    exit 1
}

# 启动前端容器
docker run -d `
    --name yufung-frontend `
    --restart unless-stopped `
    -p 80:80 `
    -e TZ=Asia/Shanghai `
    yufung-frontend:latest

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 前端容器启动成功！" -ForegroundColor Green
} else {
    Write-Host "❌ 前端容器启动失败！" -ForegroundColor Red
    exit 1
}

# 等待服务启动
Write-Host "等待服务启动..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# 检查容器状态
Write-Host "检查容器状态..." -ForegroundColor Cyan
docker ps | findstr yufung

# 显示服务信息
Write-Host ""
Write-Host "🎉 部署完成！" -ForegroundColor Green
Write-Host "📍 后端服务: http://localhost:8088" -ForegroundColor Cyan
Write-Host "📍 前端服务: http://localhost:80" -ForegroundColor Cyan
Write-Host "🔍 健康检查: http://localhost:8088/health" -ForegroundColor Cyan
Write-Host ""
Write-Host "查看日志命令:" -ForegroundColor Yellow
Write-Host "docker logs yufung-backend" -ForegroundColor White
Write-Host "docker logs yufung-frontend" -ForegroundColor White 