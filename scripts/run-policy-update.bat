@echo off
REM ========================================
REM 保单管理表结构修改脚本执行器 (Windows)
REM ========================================

echo ========================================
echo 保单管理表结构修改脚本执行器
echo ========================================
echo.

REM 检查MongoDB是否安装
echo 检查MongoDB安装状态...
mongo --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到MongoDB命令行工具
    echo    请确保MongoDB已安装并添加到系统PATH中
    echo    下载地址: https://www.mongodb.com/try/download/community
    pause
    exit /b 1
)
echo ✅ MongoDB已安装

REM 检查脚本文件是否存在
if not exist "update-policy-structure.js" (
    echo ❌ 错误: 未找到脚本文件 update-policy-structure.js
    echo    请确保脚本文件在当前目录下
    pause
    exit /b 1
)
echo ✅ 脚本文件存在

REM 提示用户确认
echo.
echo ⚠️  重要提示:
echo    1. 请确保已备份数据库
echo    2. 请确保MongoDB服务正在运行
echo    3. 此操作将修改保单表结构
echo.
set /p confirm="是否继续执行? (Y/N): "
if /i not "%confirm%"=="Y" (
    echo 操作已取消
    pause
    exit /b 0
)

REM 执行脚本
echo.
echo 开始执行保单表结构修改...
echo ========================================

REM 使用mongo命令执行JavaScript脚本
mongo yufung_admin update-policy-structure.js

REM 检查执行结果
if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo ✅ 脚本执行完成！
    echo ========================================
    echo.
    echo 建议检查以下内容:
    echo 1. 检查MongoDB输出是否有错误信息
    echo 2. 验证投保单号唯一性约束是否生效
    echo 3. 确认汇率字段精度是否正确
    echo 4. 测试应用程序功能是否正常
) else (
    echo.
    echo ========================================
    echo ❌ 脚本执行失败！
    echo ========================================
    echo.
    echo 可能的原因:
    echo 1. MongoDB服务未启动
    echo 2. 数据库连接失败
    echo 3. 权限不足
    echo 4. 数据冲突（如重复投保单号）
    echo.
    echo 请检查错误信息并重试
)

echo.
pause 