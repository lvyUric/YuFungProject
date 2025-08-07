@echo off
echo =====================================
echo MongoDB 公司集合字段升级工具
echo =====================================
echo.

REM 设置变量
set DB_NAME=insurance_db
set MONGO_HOST=106.52.172.124
set MONGO_PORT=27017

echo 请确认以下配置信息：
echo 数据库名称: %DB_NAME%
echo MongoDB 地址: %MONGO_HOST%:%MONGO_PORT%
echo.

set /p confirm="确认执行升级？(y/N): "
if /i not "%confirm%"=="y" (
    echo 操作已取消
    pause
    exit /b
)

echo.
echo 开始执行升级脚本...
echo.

REM 执行升级脚本
mongosh %DB_NAME% --host %MONGO_HOST% --port %MONGO_PORT% --file upgrade_company_schema.js

if %ERRORLEVEL% equ 0 (
    echo.
    echo ========================================
    echo 升级完成！
    echo ========================================
    echo.
    echo 请检查输出日志确认升级结果
    echo 如需回滚，请执行 rollback_company_schema.js
) else (
    echo.
    echo ========================================
    echo 升级失败！
    echo ========================================
    echo.
    echo 请检查错误信息并联系技术支持
)

echo.
pause
