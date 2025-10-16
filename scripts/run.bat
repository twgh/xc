@echo off
REM xc 运行脚本

echo Running xc command line tool...

REM 切换到项目根目录
cd /d %~dp0..

REM 运行
go run .
