@echo off
REM xc 测试脚本

echo Running tests...

REM 切换到项目根目录
cd /d %~dp0..

REM 运行单元测试
go test -v ./...

if %errorlevel% == 0 (
    echo All tests passed!
) else (
    echo Tests failed!
    exit /b %errorlevel%
)