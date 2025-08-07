#!/bin/bash

echo "====================================="
echo "MongoDB 公司集合字段升级工具"
echo "====================================="
echo

# 设置变量
DB_NAME="yufung_admin"
MONGO_HOST="localhost"
MONGO_PORT="27017"

echo "请确认以下配置信息："
echo "数据库名称: $DB_NAME"
echo "MongoDB 地址: $MONGO_HOST:$MONGO_PORT"
echo

read -p "确认执行升级？(y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo "操作已取消"
    exit 0
fi

echo
echo "开始执行升级脚本..."
echo

# 执行升级脚本
mongosh $DB_NAME --host $MONGO_HOST --port $MONGO_PORT --file upgrade_company_schema.js

if [ $? -eq 0 ]; then
    echo
    echo "========================================"
    echo "升级完成！"
    echo "========================================"
    echo
    echo "请检查输出日志确认升级结果"
    echo "如需回滚，请执行 rollback_company_schema.js"
else
    echo
    echo "========================================"
    echo "升级失败！"
    echo "========================================"
    echo
    echo "请检查错误信息并联系技术支持"
fi

echo
read -p "按任意键继续..." 