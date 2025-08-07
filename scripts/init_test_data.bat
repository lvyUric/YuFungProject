@echo off
echo ========================================
echo MongoDB 测试数据初始化脚本
echo ========================================
echo.

REM 检查MongoDB连接
echo 正在检查MongoDB连接...

REM 使用Docker容器执行MongoDB脚本
docker exec yufung-mongo mongosh insurance_db --eval "load('/docker-entrypoint-initdb.d/init_test_data.js')"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo 测试数据初始化成功！
    echo ========================================
    echo.
    echo 测试账号信息:
    echo 超级管理员: admin / secret
    echo 平台管理员: platform_admin / secret  
    echo 平安公司管理员: pingan_admin / secret
    echo 平安业务经理: zhang_manager / secret
    echo 平安普通员工: li_employee / secret
    echo 人寿公司管理员: chinalife_admin / secret
    echo 人寿销售主管: wang_supervisor / secret
    echo ========================================
) else (
    echo.
    echo ========================================
    echo 测试数据初始化失败！
    echo 请检查MongoDB连接和容器状态
    echo ========================================
)

echo.
pause 