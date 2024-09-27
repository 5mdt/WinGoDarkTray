@echo off
setlocal

rem Define variables
set APP_NAME=godark.exe
set RESOURCES=resources.rc
set ICON_FILE=icon.ico
set RESOURCE_SYSO=resources.syso

rem Clean previous builds
echo Cleaning up...
if exist %APP_NAME% del %APP_NAME%
if exist %RESOURCE_SYSO% del %RESOURCE_SYSO%

rem Compile the resources
echo Compiling resources...
windres %RESOURCES% -O coff -o %RESOURCE_SYSO%

rem Set architecture
set GOOS=windows
set GOARCH=amd64

rem Build the Go application
echo Building the application...
go build -ldflags=-H=windowsgui -o %APP_NAME%

echo Build complete.
endlocal
