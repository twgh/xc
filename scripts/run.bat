@echo off
REM xc 运行脚本

echo 正在编译并运行 xc 命令行工具...

REM 切换到项目根目录
cd /d %~dp0..

REM 编译
go build -o xc.exe cmd/main.go

if %errorlevel% == 0 (
    echo 编译成功!
) else (
    echo 编译失败!
    exit /b %errorlevel%
)

REM 运行
xc.exe %*

REM 清理临时文件
del xc.exe