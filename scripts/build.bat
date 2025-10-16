@echo off
REM xc 编译脚本

echo Building xc command line tool...

REM 切换到项目根目录
cd /d %~dp0..
REM 设置 Go 环境变量
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

REM 编译主程序 (资源文件将被自动包含)
go build -ldflags="-s -w" -trimpath -o xc.exe

if %errorlevel% == 0 (
    echo Build successful! xc.exe has been generated.
) else (
    echo Build failed!
    exit /b %errorlevel%
)

echo.
echo Build completed.