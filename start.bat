@echo off
chcp 65001 >nul

echo ========================================
echo   🚀 QuoteGenerator 启动中...
echo ========================================
echo.

set ROOT_DIR=%~dp0
set BACKEND_DIR=%ROOT_DIR%backend
set FRONTEND_DIR=%ROOT_DIR%frontend

echo 📦 启动后端服务...
start "QuoteGenerator-Backend" cmd /c "cd /d %BACKEND_DIR% && go run . --mode=http"

timeout /t 3 /nobreak >nul

echo 🎨 启动前端服务...
start "QuoteGenerator-Frontend" cmd /c "cd /d %FRONTEND_DIR% && npm run dev"

echo.
echo ========================================
echo   ✅ 服务启动完成！
echo   🌐 前端: http://localhost:3000
echo   🔧 后端: http://localhost:8080
echo ========================================
echo.
echo 关闭此窗口即可停止所有服务
pause