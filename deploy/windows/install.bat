@echo off
setlocal

net session >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Run as Administrator required
    pause
    exit /b
)

set SOURCE=%~dp0xprinter.exe
set TARGET_DIR=C:\XPrinter
set TARGET=%TARGET_DIR%\xprinter.exe
set SERVICE_NAME=XPrinter

echo Stopping service if running...

sc query "%SERVICE_NAME%" >nul 2>&1
if %errorlevel% == 0 (
    sc stop "%SERVICE_NAME%" >nul 2>&1
    timeout /t 2 >nul

    echo Deleting service...
    sc delete "%SERVICE_NAME%" >nul 2>&1
)

echo Removing old executable...

if exist "%TARGET%" (
    del /f /q "%TARGET%"
)

if exist "%TARGET_DIR%" (
    rmdir /s /q "%TARGET_DIR%"
)

echo Creating directory...
mkdir "%TARGET_DIR%"

echo Copying files...

if not exist "%SOURCE%" (
    echo ERROR: xprinter.exe not found
    pause
    goto end
)

copy /Y "%SOURCE%" "%TARGET%"
if errorlevel 1 (
    echo ERROR: Failed to copy file
    pause
    goto end
)

echo Creating service...

sc create "%SERVICE_NAME%" binPath= "\"%TARGET%\"" start= auto
if errorlevel 1 (
    echo ERROR: Failed to create service
    pause
    goto end
)

echo Starting service...

sc start "%SERVICE_NAME%"
if errorlevel 1 (
    echo ERROR: Failed to start service
    pause
    goto end
)

echo.
echo SUCCESS: XPrinter installed and running

:end
pause
endlocal