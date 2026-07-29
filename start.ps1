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

# 阻塞等待任意 Job 结束
$completedJob = Wait-Job -Job $backendJob, $frontendJob -Any

if ($completedJob.State -eq "Failed") {
    $reason = Receive-Job -Job $completedJob -ErrorAction SilentlyContinue
    Write-Host "❌ 服务异常退出: $reason" -ForegroundColor Red
} else {
    Write-Host "⚠️  服务已停止" -ForegroundColor Yellow
}

# 清理另一个仍在运行的 Job
$backendJob, $frontendJob | Where-Object { $_.State -eq "Running" } | Stop-Job
$backendJob, $frontendJob | Remove-Job -Force

Write-Host "👋 已退出" -ForegroundColor Gray
Start-Sleep -Seconds 2