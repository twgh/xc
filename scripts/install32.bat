@echo off
REM xc 32位版本安装脚本

echo Installing xc 32-bit command line tool to GOPATH...

REM Change to project root directory
cd /d %~dp0..

REM Get GOPATH
for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i

REM Create GOPATH bin directory if it doesn't exist
if not exist "%GOPATH%\bin" mkdir "%GOPATH%\bin"

REM Build 32-bit version
echo Building Windows 32-bit version...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=386
go build -o "%GOPATH%\bin\xc.exe" cmd/main.go

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