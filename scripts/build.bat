@echo off
REM xc 编译脚本

echo 正在编译 xc 命令行工具...

REM 切换到项目根目录
cd /d %~dp0..
REM 设置 Go 环境变量
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64

REM 编译主程序
go build -o xc.exe cmd/main.go

if %errorlevel% == 0 (
    echo 编译成功! xc.exe 已生成。
) else (
    echo 编译失败!
    exit /b %errorlevel%
)

echo.
echo 编译完成。