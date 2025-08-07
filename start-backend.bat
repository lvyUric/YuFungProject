@echo off
echo ====================================
echo   保险经纪管理系统 后端启动脚本
echo ====================================
echo.

echo [1/3] 检查MongoDB连接...
mongo --quiet --eval "db.adminCommand('ping')" > nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ MongoDB未运行，请先启动MongoDB服务
    echo.
    echo 方法1: 本地安装MongoDB
    echo   net start MongoDB
    echo.
    echo 方法2: 使用Docker启动MongoDB
    echo   docker run --restart=always --name mongo-dev ^
    echo     -p 27017:27017 ^
    echo     -e TZ=Asia/Shanghai ^
    echo     --privileged=true ^
    echo     -e MONGO_INITDB_ROOT_USERNAME=admin ^
    echo     -e MONGO_INITDB_ROOT_PASSWORD=yf2025 ^
    echo     -d mongo
    echo.
    pause
    exit /b 1
)
echo ✅ MongoDB连接正常

echo.
echo [2/3] 初始化数据库...
mongo insurance_db scripts/init-mongo.js
if %errorlevel% neq 0 (
    echo ⚠️  数据库初始化失败，但系统将继续启动
) else (
    echo ✅ 数据库初始化完成
)

echo.
echo [3/3] 启动后端服务...
echo 服务地址: http://localhost:8080
echo API文档: http://localhost:8080/health
echo.
echo 默认管理员账户:
echo   用户名: admin
echo   密码: admin123
echo.
echo 按 Ctrl+C 停止服务
echo ====================================
echo.

go run cmd/main.go
pause 