Write-Host "========================================" -ForegroundColor Cyan
Write-Host "MongoDB 测试数据初始化脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查Docker是否运行
Write-Host "正在检查Docker状态..." -ForegroundColor Yellow
try {
    docker version | Out-Null
    Write-Host "✓ Docker 运行正常" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker 未运行或未安装" -ForegroundColor Red
    exit 1
}

# 检查MongoDB容器是否存在
Write-Host "正在检查MongoDB容器..." -ForegroundColor Yellow
$mongoContainer = docker ps -a --format "table {{.Names}}" | Select-String "mongo"

if (-not $mongoContainer) {
    Write-Host "✗ 未找到MongoDB容器" -ForegroundColor Red
    Write-Host "请先启动MongoDB容器" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ 找到MongoDB容器" -ForegroundColor Green

# 将初始化脚本复制到容器中
Write-Host "正在复制初始化脚本到容器..." -ForegroundColor Yellow
$scriptPath = Join-Path (Get-Location) "scripts\mongodb\init_test_data.js"

if (-not (Test-Path $scriptPath)) {
    Write-Host "✗ 初始化脚本不存在: $scriptPath" -ForegroundColor Red
    exit 1
}

# 复制脚本到容器
docker cp $scriptPath yufung-mongo:/tmp/init_test_data.js

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ 脚本复制成功" -ForegroundColor Green
} else {
    Write-Host "✗ 脚本复制失败" -ForegroundColor Red
    exit 1
}

# 执行初始化脚本
Write-Host "正在执行初始化脚本..." -ForegroundColor Yellow
Write-Host ""

docker exec yufung-mongo mongosh insurance_db --eval "load('/tmp/init_test_data.js')"

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "测试数据初始化成功！" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "测试账号信息:" -ForegroundColor Cyan
    Write-Host "超级管理员: admin / secret" -ForegroundColor White
    Write-Host "平台管理员: platform_admin / secret" -ForegroundColor White
    Write-Host "平安公司管理员: pingan_admin / secret" -ForegroundColor White
    Write-Host "平安业务经理: zhang_manager / secret" -ForegroundColor White
    Write-Host "平安普通员工: li_employee / secret" -ForegroundColor White
    Write-Host "人寿公司管理员: chinalife_admin / secret" -ForegroundColor White
    Write-Host "人寿销售主管: wang_supervisor / secret" -ForegroundColor White
    Write-Host "========================================" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "测试数据初始化失败！" -ForegroundColor Red
    Write-Host "请检查MongoDB连接和容器状态" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
}

Write-Host ""
Write-Host "按任意键继续..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown") 