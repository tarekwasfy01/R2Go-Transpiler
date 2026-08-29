@echo off
setlocal
where go >nul 2>nul || (
  echo Go was not found in PATH.
  exit /b 1
)
go test ./...
if errorlevel 1 exit /b %errorlevel%
echo.
echo Go build/test passed.
endlocal
