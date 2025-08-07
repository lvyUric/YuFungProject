@echo off
echo ====================================
echo   保险经纪管理系统 数据库初始化
echo ====================================
echo.

echo 正在初始化数据库...
mongo insurance_db scripts/init-mongo.js

if %errorlevel% equ 0 (
    echo.
    echo ✅ 数据库初始化完成！
    echo.
    echo =================================
    echo 默认管理员账户信息：
    echo 用户名: admin
    echo 密码: admin123
    echo =================================
    echo.
    echo 请在首次登录后修改默认密码！
) else (
    echo.
    echo ❌ 数据库初始化失败！
    echo 请检查MongoDB是否正常运行
)

echo.
pause 