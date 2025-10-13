@echo off
REM xc 测试脚本

echo 正在运行测试...

REM 切换到项目根目录
cd /d %~dp0..

REM 运行单元测试
go test -v ./...

if %errorlevel% == 0 (
    echo 所有测试通过!
) else (
    echo 测试失败!
    exit /b %errorlevel%
)