#!/bin/bash

# 设置Docker构建参数
export DOCKER_BUILDKIT=1
export BUILDKIT_PROGRESS=plain

# 设置Go代理环境变量
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
export GO111MODULE=on

echo "开始构建Docker镜像..."

# 构建镜像，使用BuildKit和缓存
docker build \
    --build-arg GOPROXY=https://goproxy.cn,direct \
    --build-arg GOSUMDB=sum.golang.google.cn \
    --build-arg GO111MODULE=on \
    --progress=plain \
    --no-cache \
    -t yufung-admin:latest .

if [ $? -eq 0 ]; then
    echo "✅ Docker镜像构建成功！"
    echo "镜像名称: yufung-admin:latest"
    echo ""
    echo "运行容器命令:"
    echo "docker run -d -p 8088:8088 --name yufung-admin yufung-admin:latest"
else
    echo "❌ Docker镜像构建失败！"
    exit 1
fi 