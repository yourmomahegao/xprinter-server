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

echo Copying files...

if not exist "%SOURCE%" (
    echo ERROR: xprinter.exe not found
    pause
    goto end
)

if not exist "%TARGET_DIR%" mkdir "%TARGET_DIR%"

copy /Y "%SOURCE%" "%TARGET%"
if errorlevel 1 (
    echo ERROR: Failed to copy file
    pause
    goto end
)

echo Checking service...

sc query "%SERVICE_NAME%" >nul 2>&1
if %errorlevel% neq 0 (
    echo Creating service...
    sc create "%SERVICE_NAME%" binPath= "\"%TARGET%\"" start= auto
    if errorlevel 1 (
        echo ERROR: Failed to create service
        pause
        goto end
    )
) else (
    echo Service already exists
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