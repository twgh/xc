@echo off
REM xc Windows平台编译脚本

echo 正在编译 xc 命令行工具 (Windows平台)...

REM 切换到项目根目录
cd /d %~dp0..

REM 创建输出目录
if not exist "dist" mkdir "dist"

REM Windows 64位
echo 编译 Windows 64位版本...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -o dist/xc-windows-amd64.exe cmd/main.go

if %errorlevel% == 0 (
    echo Windows 64位版本编译成功!
) else (
    echo Windows 64位版本编译失败!
    exit /b %errorlevel%
)

REM Windows 32位
echo 编译 Windows 32位版本...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386
go build -o dist/xc-windows-386.exe cmd/main.go

if %errorlevel% == 0 (
    echo Windows 32位版本编译成功!
) else (
    echo Windows 32位版本编译失败!
    exit /b %errorlevel%
)

echo.
echo Windows平台编译完成!
echo 输出文件位于 dist 目录中:
echo - xc-windows-amd64.exe (64位)
echo - xc-windows-386.exe (32位)