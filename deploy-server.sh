#!/bin/bash

# 服务器部署脚本 (Linux版本)
echo "开始服务器部署..."

# 设置环境变量
export DOCKER_BUILDKIT=1
export GOPROXY=https://goproxy.cn,direct

echo "环境变量设置完成"

# 检查Docker状态
echo "检查Docker状态..."
docker --version
if [ $? -ne 0 ]; then
    echo "❌ Docker未安装或未启动"
    exit 1
fi

# 停止并删除现有容器
echo "停止现有容器..."
docker stop yufung-backend yufung-frontend 2>/dev/null
docker rm yufung-backend yufung-frontend 2>/dev/null

# 删除旧镜像
echo "清理旧镜像..."
docker rmi yufung-backend:latest yufung-frontend:latest 2>/dev/null

# 构建后端
echo "开始构建后端Docker镜像..."
docker build \
    --build-arg GOPROXY=https://goproxy.cn,direct \
    --build-arg GOSUMDB=sum.golang.google.cn \
    --build-arg GO111MODULE=on \
    --progress=plain \
    --no-cache \
    -t yufung-backend:latest .

if [ $? -eq 0 ]; then
    echo "✅ 后端Docker镜像构建成功！"
else
    echo "❌ 后端Docker镜像构建失败！"
    exit 1
fi

# 构建前端
echo "开始构建前端Docker镜像..."
cd Yufung-admin-front

# 检查package-lock.json是否存在
if [ ! -f "package-lock.json" ]; then
    echo "❌ package-lock.json 文件不存在，正在生成..."
    npm install
fi

docker build \
    --progress=plain \
    --no-cache \
    -t yufung-frontend:latest .

if [ $? -eq 0 ]; then
    echo "✅ 前端Docker镜像构建成功！"
else
    echo "❌ 前端Docker镜像构建失败！"
    cd ..
    exit 1
fi

# 返回原目录
cd ..

# 运行容器
echo "启动应用容器..."

# 启动后端容器
docker run -d \
    --name yufung-backend \
    --restart unless-stopped \
    -p 8088:8088 \
    -v $(pwd)/logs:/app/logs \
    -v $(pwd)/uploads:/app/uploads \
    -e TZ=Asia/Shanghai \
    yufung-backend:latest

if [ $? -eq 0 ]; then
    echo "✅ 后端容器启动成功！"
else
    echo "❌ 后端容器启动失败！"
    exit 1
fi

# 启动前端容器
docker run -d \
    --name yufung-frontend \
    --restart unless-stopped \
    -p 80:80 \
    -e TZ=Asia/Shanghai \
    yufung-frontend:latest

if [ $? -eq 0 ]; then
    echo "✅ 前端容器启动成功！"
else
    echo "❌ 前端容器启动失败！"
    exit 1
fi

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查容器状态
echo "检查容器状态..."
docker ps | grep yufung

# 显示服务信息
echo ""
echo "🎉 部署完成！"
echo "📍 后端服务: http://localhost:8088"
echo "📍 前端服务: http://localhost:80"
echo "🔍 健康检查: http://localhost:8088/health"
echo ""
echo "查看日志命令:"
echo "docker logs yufung-backend"
echo "docker logs yufung-frontend" 