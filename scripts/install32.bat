@echo off
REM xc 32 位版本安装脚本

echo Installing xc 32-bit command line tool to GOPATH...

REM 切换到项目根目录
cd /d %~dp0..

REM 获取 GOPATH
for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i

REM 如果 %GOPATH%\bin 目录不存在则创建
if not exist "%GOPATH%\bin" mkdir "%GOPATH%\bin"

REM 编译 32 位版本
echo Building Windows 32-bit version...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386
go build -ldflags="-s -w" -trimpath -o "%GOPATH%\bin\xc.exe" .\cmd

if %errorlevel% == 0 (
    echo Windows 32-bit version built successfully!
) else (
    echo Windows 32-bit version build failed!
    exit /b %errorlevel%
)

echo.
echo Installation complete!
echo 32-bit version installed to %GOPATH%\bin\xc.exe
echo Make sure %GOPATH%\bin is in your system PATH