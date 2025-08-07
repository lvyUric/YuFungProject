#!/bin/bash

echo "🏥 保险经纪管理系统启动脚本"
echo "=================================="
echo

# 检查是否安装了必要的依赖
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.21+"
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js 16+"
    exit 1
fi

# 检查MongoDB
echo "🔍 检查MongoDB连接..."
if command -v mongo &> /dev/null; then
    if mongo --quiet --eval "db.adminCommand('ping')" &> /dev/null; then
        echo "✅ MongoDB连接正常"
    else
        echo "❌ MongoDB未运行，请先启动MongoDB服务"
        echo
        echo "方法1: 系统服务启动"
        echo "  sudo systemctl start mongod"
        echo
        echo "方法2: Docker启动"
        echo "  docker run --restart=always --name mongo-dev \\"
        echo "    -p 27017:27017 \\"
        echo "    -e TZ=Asia/Shanghai \\"
        echo "    --privileged=true \\"
        echo "    -e MONGO_INITDB_ROOT_USERNAME=admin \\"
        echo "    -e MONGO_INITDB_ROOT_PASSWORD=yf2025 \\"
        echo "    -d mongo"
        echo
        exit 1
    fi
else
    echo "⚠️  警告：mongo命令行工具未找到"
    echo "请确保MongoDB正在localhost:27017运行"
    echo
fi

echo "🔧 安装Go依赖..."
go mod tidy

echo "🔧 安装前端依赖..."
cd stelory-admin
if [ ! -d "node_modules" ]; then
    echo "正在安装前端依赖..."
    npm install
else
    echo "✅ 前端依赖已安装"
fi
cd ..

echo "🗄️  初始化数据库..."
if command -v mongo &> /dev/null; then
    mongo insurance_db scripts/init-mongo.js
    if [ $? -eq 0 ]; then
        echo "✅ 数据库初始化完成"
    else
        echo "⚠️  数据库初始化失败，但系统将继续启动"
    fi
else
    echo "⚠️  跳过数据库初始化（mongo命令行工具未找到）"
fi

echo
echo "🚀 启动服务..."
echo "后端将在 http://localhost:8080 启动"
echo "前端将在 http://localhost:3000 启动"
echo
echo "默认管理员账户："
echo "  用户名: admin"
echo "  密码: admin123"
echo
echo "按 Ctrl+C 停止服务"
echo

# 启动后端（后台）
echo "启动后端服务..."
go run cmd/main.go &
BACKEND_PID=$!

# 等待后端启动
sleep 3

# 启动前端
echo "启动前端服务..."
cd stelory-admin
npm start &
FRONTEND_PID=$!

# 等待用户中断
wait

# 清理进程
echo "正在停止服务..."
kill $BACKEND_PID 2>/dev/null
kill $FRONTEND_PID 2>/dev/null
echo "服务已停止" 