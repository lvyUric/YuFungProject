@echo off
echo 启动 Yufung Admin 前端...
echo.

REM 检查是否存在 node_modules 目录
if not exist "Yufung-admin-front\node_modules" (
    echo 检测到未安装依赖，正在安装...
    cd Yufung-admin-front
    npm install
    if errorlevel 1 (
        echo 依赖安装失败，请检查网络连接或Node.js环境
        pause
        exit /b 1
    )
    cd ..
)

REM 启动前端开发服务器
echo 启动前端开发服务器...
cd Yufung-admin-front
npm start

pause 