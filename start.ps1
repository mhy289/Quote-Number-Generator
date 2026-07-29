param(
    [string]$BackendPort = "8080",
    [string]$FrontendPort = "3000"
)

$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BackendDir = Join-Path $RootDir "backend"
$FrontendDir = Join-Path $RootDir "frontend"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  🚀 QuoteGenerator 启动中..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 启动后端服务
Write-Host "📦 启动后端服务 (端口: $BackendPort)..." -ForegroundColor Yellow
$backendJob = Start-Job -ScriptBlock {
    param($dir, $port)
    Set-Location $dir
    go run . --mode=http --port=$port
} -ArgumentList $BackendDir, $BackendPort

# 等待后端启动
Start-Sleep -Seconds 3

# 启动前端服务
Write-Host "🎨 启动前端服务 (端口: $FrontendPort)..." -ForegroundColor Yellow
$frontendJob = Start-Job -ScriptBlock {
    param($dir, $port)
    Set-Location $dir
    npm run dev
} -ArgumentList $FrontendDir, $FrontendPort

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  ✅ 服务启动完成！" -ForegroundColor Green
Write-Host "  🌐 前端: http://localhost:$FrontendPort" -ForegroundColor Green
Write-Host "  🔧 后端: http://localhost:$BackendPort" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "按 Ctrl+C 停止所有服务" -ForegroundColor Gray

# 等待任意一个 Job 结束
while ($true) {
    $backendState = (Receive-Job -Job $backendJob -ErrorAction SilentlyContinue)
    $frontendState = (Receive-Job -Job $frontendJob -ErrorAction SilentlyContinue)

    if ($backendJob.State -eq "Failed") {
        Write-Host "❌ 后端服务异常退出" -ForegroundColor Red
        Write-Host $backendState -ForegroundColor Red
        break
    }
    if ($frontendJob.State -eq "Failed") {
        Write-Host "❌ 前端服务异常退出" -ForegroundColor Red
        Write-Host $frontendState -ForegroundColor Red
        break
    }

    Start-Sleep -Seconds 2
}

# 清理
Write-Host "🛑 正在停止所有服务..." -ForegroundColor Yellow
Stop-Job -Job $backendJob -ErrorAction SilentlyContinue
Stop-Job -Job $frontendJob -ErrorAction SilentlyContinue
Remove-Job -Job $backendJob -ErrorAction SilentlyContinue
Remove-Job -Job $frontendJob -ErrorAction SilentlyContinue
Write-Host "👋 已退出" -ForegroundColor Gray