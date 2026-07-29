@echo off
chcp 437 >nul

echo ========================================
echo   QuoteGenerator Starting...
echo ========================================
echo.

set ROOT_DIR=%~dp0
set BACKEND_DIR=%ROOT_DIR%backend
set FRONTEND_DIR=%ROOT_DIR%frontend

echo [1/2] Starting backend service...
start "QuoteGenerator-Backend" "%ComSpec%" /c "cd /d %BACKEND_DIR% && go run . --mode=http"

timeout /t 3 /nobreak >nul

echo [2/2] Starting frontend service...
start "QuoteGenerator-Frontend" "%ComSpec%" /c "cd /d %FRONTEND_DIR% && npm run dev"

echo.
echo ========================================
echo   All services started!
echo   Frontend: http://localhost:3000
echo   Backend:  http://localhost:8080
echo ========================================
echo.
echo Press any key to close this window (services keep running in background)
pause >nul