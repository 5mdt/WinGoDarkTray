@echo off
setlocal

rem Define variables
set APP_NAME=godark.exe

rem Clean previous builds
echo Cleaning up...
if exist %APP_NAME% del %APP_NAME%

rem Set architecture
set GOOS=windows
set GOARCH=amd64

rem Build the Go application
echo Building the application...
go build -ldflags=-H=windowsgui -o %APP_NAME%

echo Build complete.
endlocal
