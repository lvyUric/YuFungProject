# ========================================
# 保单管理表结构修改脚本执行器 (PowerShell)
# 使用方法: .\Run-PolicyUpdate.ps1
# ========================================

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "保单管理表结构修改脚本执行器" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 设置错误处理
$ErrorActionPreference = "Stop"

try {
    # 检查MongoDB是否安装
    Write-Host "检查MongoDB安装状态..." -ForegroundColor Yellow
    
    $mongoVersion = mongo --version 2>$null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 错误: 未找到MongoDB命令行工具" -ForegroundColor Red
        Write-Host "   请确保MongoDB已安装并添加到系统PATH中" -ForegroundColor Red
        Write-Host "   下载地址: https://www.mongodb.com/try/download/community" -ForegroundColor Red
        Read-Host "按任意键退出"
        exit 1
    }
    Write-Host "✅ MongoDB已安装" -ForegroundColor Green
    
    # 检查脚本文件是否存在
    if (-not (Test-Path "update-policy-structure.js")) {
        Write-Host "❌ 错误: 未找到脚本文件 update-policy-structure.js" -ForegroundColor Red
        Write-Host "   请确保脚本文件在当前目录下" -ForegroundColor Red
        Read-Host "按任意键退出"
        exit 1
    }
    Write-Host "✅ 脚本文件存在" -ForegroundColor Green
    
    # 提示用户确认
    Write-Host ""
    Write-Host "⚠️  重要提示:" -ForegroundColor Yellow
    Write-Host "   1. 请确保已备份数据库" -ForegroundColor Yellow
    Write-Host "   2. 请确保MongoDB服务正在运行" -ForegroundColor Yellow
    Write-Host "   3. 此操作将修改保单表结构" -ForegroundColor Yellow
    Write-Host ""
    
    $confirm = Read-Host "是否继续执行? (Y/N)"
    if ($confirm -notmatch '^[Yy]$') {
        Write-Host "操作已取消" -ForegroundColor Yellow
        Read-Host "按任意键退出"
        exit 0
    }
    
    # 执行前检查MongoDB服务
    Write-Host ""
    Write-Host "检查MongoDB服务状态..." -ForegroundColor Yellow
    
    # 尝试连接到MongoDB
    $testConnection = mongo --eval "db.runCommand({ping: 1})" --quiet 2>$null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 无法连接到MongoDB服务" -ForegroundColor Red
        Write-Host "   请确保MongoDB服务正在运行" -ForegroundColor Red
        Write-Host ""
        Write-Host "在Windows上启动MongoDB服务:" -ForegroundColor Cyan
        Write-Host "   net start MongoDB" -ForegroundColor Cyan
        Write-Host "   或者" -ForegroundColor Cyan
        Write-Host "   Start-Service MongoDB" -ForegroundColor Cyan
        Read-Host "按任意键退出"
        exit 1
    }
    Write-Host "✅ MongoDB服务运行正常" -ForegroundColor Green
    
    # 执行脚本
    Write-Host ""
    Write-Host "开始执行保单表结构修改..." -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    
    # 执行MongoDB脚本
    $result = mongo yufung_admin update-policy-structure.js
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Green
        Write-Host "✅ 脚本执行完成！" -ForegroundColor Green
        Write-Host "========================================" -ForegroundColor Green
        Write-Host ""
        Write-Host "建议检查以下内容:" -ForegroundColor Yellow
        Write-Host "1. 检查上方MongoDB输出是否有错误信息" -ForegroundColor Yellow
        Write-Host "2. 验证投保单号唯一性约束是否生效" -ForegroundColor Yellow
        Write-Host "3. 确认汇率字段精度是否正确" -ForegroundColor Yellow
        Write-Host "4. 测试应用程序功能是否正常" -ForegroundColor Yellow
        
        # 提供验证命令
        Write-Host ""
        Write-Host "验证命令:" -ForegroundColor Cyan
        Write-Host "mongo yufung_admin --eval `"db.policies.getIndexes().filter(i => i.name.includes('proposal')).forEach(i => print(JSON.stringify(i, null, 2)))`"" -ForegroundColor Gray
        
    } else {
        Write-Host ""
        Write-Host "========================================" -ForegroundColor Red
        Write-Host "❌ 脚本执行失败！" -ForegroundColor Red
        Write-Host "========================================" -ForegroundColor Red
        Write-Host ""
        Write-Host "可能的原因:" -ForegroundColor Yellow
        Write-Host "1. MongoDB服务未启动" -ForegroundColor Yellow
        Write-Host "2. 数据库连接失败" -ForegroundColor Yellow
        Write-Host "3. 权限不足" -ForegroundColor Yellow
        Write-Host "4. 数据冲突（如重复投保单号）" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "请检查上方错误信息并重试" -ForegroundColor Yellow
    }
    
} catch {
    Write-Host ""
    Write-Host "❌ 执行过程中发生错误:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host ""
    Write-Host "请检查错误信息并重试" -ForegroundColor Yellow
}

Write-Host ""
Read-Host "按任意键退出" 