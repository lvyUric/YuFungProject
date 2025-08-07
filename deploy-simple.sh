#!/bin/bash

# 简化服务器部署脚本
echo "开始简化部署..."

# 检查Docker状态
echo "检查Docker状态..."
docker --version || { echo "❌ Docker未安装"; exit 1; }

# 停止并删除现有容器
echo "清理现有容器..."
docker stop yufung-backend yufung-frontend 2>/dev/null
docker rm yufung-backend yufung-frontend 2>/dev/null

# 删除旧镜像
echo "清理旧镜像..."
docker rmi yufung-backend:latest yufung-frontend:latest 2>/dev/null

# 构建后端
echo "构建后端镜像..."
docker build \
    --build-arg GOPROXY=https://goproxy.cn,direct \
    --build-arg GOSUMDB=sum.golang.google.cn \
    --build-arg GO111MODULE=on \
    --no-cache \
    -t yufung-backend:latest . || { echo "❌ 后端构建失败"; exit 1; }

echo "✅ 后端构建成功"

# 构建前端
echo "构建前端镜像..."
cd Yufung-admin-front
docker build --no-cache -t yufung-frontend:latest . || { echo "❌ 前端构建失败"; cd ..; exit 1; }
cd ..

echo "✅ 前端构建成功"

# 启动容器
echo "启动容器..."

# 启动后端
docker run -d \
    --name yufung-backend \
    --restart unless-stopped \
    -p 8088:8088 \
    -v $(pwd)/logs:/app/logs \
    -v $(pwd)/uploads:/app/uploads \
    -e TZ=Asia/Shanghai \
    yufung-backend:latest || { echo "❌ 后端启动失败"; exit 1; }

# 启动前端
docker run -d \
    --name yufung-frontend \
    --restart unless-stopped \
    -p 80:80 \
    -e TZ=Asia/Shanghai \
    yufung-frontend:latest || { echo "❌ 前端启动失败"; exit 1; }

echo "✅ 容器启动成功"

# 等待启动
echo "等待服务启动..."
sleep 15

# 检查状态
echo "检查服务状态..."
docker ps | grep yufung

echo ""
echo "🎉 部署完成！"
echo "📍 后端: http://localhost:8088"
echo "📍 前端: http://localhost:80"
echo "🔍 健康检查: http://localhost:8088/health" 