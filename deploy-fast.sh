#!/bin/bash

# 彻底清理部署脚本
echo "🚀 开始彻底清理部署..."

# 检查Docker状态
echo "📋 检查Docker状态..."
docker --version || { echo "❌ Docker未安装"; exit 1; }

# 停止所有相关容器
echo "🛑 停止所有相关容器..."
docker stop $(docker ps -aq --filter "name=yufung") 2>/dev/null || true
docker rm $(docker ps -aq --filter "name=yufung") 2>/dev/null || true

# 删除所有相关镜像
echo "🗑️ 删除所有相关镜像..."
docker rmi $(docker images --filter "reference=yufung*" -q) 2>/dev/null || true

# 彻底清理Docker缓存
echo "🧹 彻底清理Docker缓存..."
docker system prune -af --volumes
docker builder prune -af

# 删除Docker网络
echo "🌐 重建Docker网络..."
docker network rm yufung-network 2>/dev/null || true
docker network create yufung-network

# 清理前端所有构建产物和缓存
echo "🧽 彻底清理前端缓存..."
cd Yufung-admin-front
rm -rf dist
rm -rf build
rm -rf .umi
rm -rf .umi-production
rm -rf node_modules/.cache
rm -rf node_modules/.vite
rm -rf node_modules/.max
find . -name "*.cache" -type f -delete 2>/dev/null || true

# 重新安装依赖（确保没有缓存）
echo "📦 重新安装前端依赖..."
rm -rf node_modules package-lock.json
npm install --no-audit --no-fund

# 手动构建前端（不通过Docker）
echo "🔨 手动构建前端..."
echo "🔧 当前环境变量:"
echo "NODE_ENV: $NODE_ENV"
echo "UMI_ENV: $UMI_ENV"

# 设置环境变量确保使用正确的配置
export NODE_ENV=production
export UMI_ENV=dev

echo "🔧 设置后的环境变量:"
echo "NODE_ENV: $NODE_ENV"
echo "UMI_ENV: $UMI_ENV"

npm run build

# 检查构建结果
echo "🔍 检查构建结果..."
if [ ! -d "dist" ]; then
    echo "❌ 前端构建失败，dist目录不存在"
    exit 1
fi

echo "📁 dist目录内容:"
ls -la dist/

# 检查是否还有旧的API地址
echo "🔎 检查是否还有旧的API地址..."
if grep -r "proapi.azurewebsites.net" dist/ 2>/dev/null; then
    echo "❌ 构建结果中仍包含旧的API地址！"
    echo "📄 包含旧地址的文件："
    grep -r "proapi.azurewebsites.net" dist/ 2>/dev/null || true
    exit 1
else
    echo "✅ 构建结果检查通过，未发现旧的API地址"
fi

# 检查是否包含正确的API地址
echo "🔎 检查是否包含正确的API地址..."
if grep -r "106.52.172.124:8088" dist/ 2>/dev/null; then
    echo "✅ 找到正确的API地址："
    grep -r "106.52.172.124:8088" dist/ 2>/dev/null | head -5
else
    echo "⚠️ 未找到正确的API地址，这可能是个问题"
fi

cd ..

# 构建后端镜像
echo "🏗️ 构建后端镜像..."
docker build --no-cache --pull \
    --build-arg GOPROXY=https://goproxy.cn,direct \
    --build-arg GOSUMDB=sum.golang.google.cn \
    --build-arg GO111MODULE=on \
    -t yufung-backend:latest . || { echo "❌ 后端构建失败"; exit 1; }

echo "✅ 后端构建成功"

# 构建前端镜像
echo "🏗️ 构建前端镜像..."
cd Yufung-admin-front

# 显示构建上下文信息
echo "📋 构建上下文信息:"
echo "当前目录: $(pwd)"
echo "Dockerfile存在: $(test -f Dockerfile && echo '是' || echo '否')"
echo "dist目录存在: $(test -d dist && echo '是' || echo '否')"

# 构建镜像并显示详细输出
docker build --no-cache --pull --progress=plain -t yufung-frontend:latest . || { echo "❌ 前端构建失败"; cd ..; exit 1; }
cd ..

echo "✅ 前端构建成功"

# 启动后端容器
echo "🚀 启动后端容器..."
docker run -d \
    --name yufung-backend \
    --network yufung-network \
    --restart unless-stopped \
    -p 8088:8088 \
    -v $(pwd)/logs:/app/logs \
    -v $(pwd)/uploads:/app/uploads \
    -e TZ=Asia/Shanghai \
    yufung-backend:latest || { echo "❌ 后端启动失败"; exit 1; }

# 启动前端容器
echo "🚀 启动前端容器..."
docker run -d \
    --name yufung-frontend \
    --network yufung-network \
    --restart unless-stopped \
    -p 8080:8080 \
    -e TZ=Asia/Shanghai \
    yufung-frontend:latest || { echo "❌ 前端启动失败"; exit 1; }

echo "✅ 容器启动成功"

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 15

# 检查服务状态
echo "📊 检查服务状态..."
docker ps | grep yufung

# 检查前端容器日志
echo "📋 检查前端容器日志..."
docker logs yufung-frontend --tail 10

# 测试后端连通性
echo "🔗 测试后端连通性..."
curl -f http://localhost:8088/health || echo "⚠️ 后端健康检查失败"

echo ""
echo "🎉 彻底清理部署完成！"
echo "📍 后端: http://localhost:8088"
echo "📍 前端: http://localhost:8080"
echo "🔍 健康检查: http://localhost:8088/health"
echo ""
echo "🔧 如果仍有问题，请清除浏览器缓存并硬刷新页面！" 