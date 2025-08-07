#!/bin/bash

# 在容器内测试MongoDB连接的脚本

echo "🔍 测试从容器内连接MongoDB..."

# 测试网络连通性
echo "📡 测试网络连通性..."
if ping -c 3 106.52.172.124 > /dev/null 2>&1; then
    echo "✅ 网络连通性正常"
else
    echo "❌ 网络连通性异常"
fi

# 测试端口连通性
echo "🔌 测试MongoDB端口连通性..."
if nc -z 106.52.172.124 27017 2>/dev/null; then
    echo "✅ MongoDB端口27017可访问"
else
    echo "❌ MongoDB端口27017不可访问"
fi

# 测试MongoDB连接
echo "🗄️ 测试MongoDB连接..."
if mongosh "mongodb://admin:yf2025@106.52.172.124:27017/insurance_db?authSource=admin" --eval "db.runCommand('ping')" > /dev/null 2>&1; then
    echo "✅ MongoDB连接成功"
else
    echo "❌ MongoDB连接失败"
fi

echo "🎉 测试完成" 