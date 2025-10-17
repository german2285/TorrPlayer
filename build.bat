@echo off
echo ================================================
echo   Building TorrPlayer Merged for Windows
echo ================================================
echo.

REM Set environment for Windows build
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

echo [1/3] Installing frontend dependencies...
cd frontend
call npm install
if errorlevel 1 (
    echo ERROR: npm install failed
    exit /b 1
)

echo.
echo [2/3] Building frontend...
call npm run build
if errorlevel 1 (
    echo ERROR: Frontend build failed
    exit /b 1
)
cd ..

echo.
echo [3/3] Building Go application with Wails...
wails build -clean -platform windows/amd64
if errorlevel 1 (
    echo ERROR: Wails build failed
    exit /b 1
)

echo.
echo ================================================
echo   Build completed successfully!
echo ================================================
echo.
echo Executable location: build\bin\torrplayer-merged.exe
echo.
echo IMPORTANT: Copy libmpv-2.dll to build\bin\ directory
echo.
pause
