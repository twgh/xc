@echo off
REM xc 清理脚本

echo Cleaning build artifacts...

REM 切换到项目根目录
cd /d %~dp0..

REM 删除编译生成的文件
if exist "xc.exe" del "xc.exe"
if exist "dist" rmdir /s /q "dist"

echo Cleanup completed!