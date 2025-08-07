@echo off
echo ====================================
echo   保险经纪管理系统 前端启动脚本
echo ====================================
echo.

echo [1/2] 检查依赖...
cd stelory-admin
if not exist "node_modules" (
    echo 正在安装前端依赖...
    npm install
) else (
    echo ✅ 依赖已安装
)

echo.
echo [2/2] 启动前端服务...
echo 前端地址: http://localhost:3000
echo 后端地址: http://localhost:8080
echo.
echo 默认管理员账户:
echo   用户名: admin
echo   密码: admin123
echo.
echo 按 Ctrl+C 停止服务
echo ====================================
echo.

npm start
pause 