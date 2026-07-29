@echo off
chcp 437 >nul

echo ========================================
echo   QuoteGenerator Starting...
echo ========================================
echo.

set ROOT_DIR=%~dp0
set BACKEND_DIR=%ROOT_DIR%backend
set FRONTEND_DIR=%ROOT_DIR%frontend
set BACKEND_PORT=8080
set FRONTEND_PORT=3000

echo [1/2] Starting backend service...
start "QuoteGenerator-Backend" "%ComSpec%" /c "cd /d %BACKEND_DIR% && go run . --mode=http"

echo.
echo Checking backend service...
set BACKEND_OK=0
for /l %%i in (1,1,15) do (
    timeout /t 1 /nobreak >nul
    >nul 2>&1 powershell -NoProfile "try{$r=Invoke-WebRequest -Uri 'http://localhost:%BACKEND_PORT%' -UseBasicParsing -TimeoutSec 2; if($r.StatusCode -eq 200){exit 0}}catch{exit 1}" && (
        set BACKEND_OK=1
        goto :backend_ready
    )
)
:backend_ready

if "%BACKEND_OK%"=="1" (
    echo [OK] Backend service is running on http://localhost:%BACKEND_PORT%
) else (
    echo [WARN] Backend service may not have started properly!
    echo        Check the backend window for errors.
    echo        Expected: http://localhost:%BACKEND_PORT%
)

echo.
echo [2/2] Starting frontend service...
start "QuoteGenerator-Frontend" "%ComSpec%" /c "cd /d %FRONTEND_DIR% && npm run dev"

echo.
echo Checking frontend service...
set FRONTEND_OK=0
for /l %%i in (1,1,20) do (
    timeout /t 1 /nobreak >nul
    >nul 2>&1 powershell -NoProfile "try{$r=Invoke-WebRequest -Uri 'http://localhost:%FRONTEND_PORT%' -UseBasicParsing -TimeoutSec 2; if($r.StatusCode -eq 200){exit 0}}catch{exit 1}" && (
        set FRONTEND_OK=1
        goto :frontend_ready
    )
)
:frontend_ready

if "%FRONTEND_OK%"=="1" (
    echo [OK] Frontend service is running on http://localhost:%FRONTEND_PORT%
) else (
    echo [WARN] Frontend service may not have started properly!
    echo        Check the frontend window for errors.
    echo        Expected: http://localhost:%FRONTEND_PORT%
)

echo.
echo ========================================
if "%BACKEND_OK%"=="1" if "%FRONTEND_OK%"=="1" (
    echo   All services started successfully!
) else (
    echo   Some services may have issues - check above
)
echo   Frontend: http://localhost:%FRONTEND_PORT%
echo   Backend:  http://localhost:%BACKEND_PORT%
echo ========================================
echo.
echo Press any key to close this window (services keep running in background)
pause >nul