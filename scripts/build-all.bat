@echo off
REM xc Windows平台编译脚本

echo Building xc command line tool (Windows platform)...

REM 切换到项目根目录
cd /d %~dp0..

REM 创建输出目录
if not exist "dist" mkdir "dist"

REM Windows 64位
echo Building Windows 64-bit version...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -trimpath -o dist\xc-windows-amd64.exe

if %errorlevel% == 0 (
    echo Windows 64-bit version build successful!
) else (
    echo Windows 64-bit version build failed!
    exit /b %errorlevel%
)

REM Windows 32位
echo Building Windows 32-bit version...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386
go build -ldflags="-s -w" -trimpath -o dist\xc-windows-386.exe

if %errorlevel% == 0 (
    echo Windows 32-bit version build successful!
) else (
    echo Windows 32-bit version build failed!
    exit /b %errorlevel%
)

echo.
echo Windows platform build completed!
echo Output files are located in the dist directory:
echo - xc-windows-amd64.exe (64-bit)
echo - xc-windows-386.exe (32-bit)