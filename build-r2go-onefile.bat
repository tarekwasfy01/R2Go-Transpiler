@echo off
setlocal
cd /d "%~dp0"

set "GOEXE=go"
if exist "C:\Program Files\Go\bin\go.exe" (
  set "GOEXE=C:\Program Files\Go\bin\go.exe"
  set "PATH=C:\Program Files\Go\bin;%PATH%"
)
"%GOEXE%" version >nul 2>nul || (
  echo ERROR: Go was not found in PATH or C:\Program Files\Go\bin.
  exit /b 1
)

if not exist .cache\go mkdir .cache\go
if not exist .cache\go-build mkdir .cache\go-build
if not exist .cache\go-mod mkdir .cache\go-mod
if not exist dist mkdir dist

set "GOPATH=%CD%\.cache\go"
set "GOCACHE=%CD%\.cache\go-build"
set "GOMODCACHE=%CD%\.cache\go-mod"

echo [1/5] Resolving modules...
"%GOEXE%" mod download || exit /b 1

 echo [2/5] Running complete Go test suite...
 "%GOEXE%" test ./... || exit /b 1

echo [3/5] Embedding project and dependency licenses...
"%GOEXE%" generate ./internal/licenses || exit /b 1

echo [4/5] Preparing the Gio Windows onefile packager...
if not exist "%GOPATH%\bin\gogio.exe" (
  "%GOEXE%" install gioui.org/cmd/gogio@v0.10.0 || exit /b 1
)

echo [5/5] Building dist\r2go.exe with embedded icon...
"%GOPATH%\bin\gogio.exe" -target windows -buildmode exe -o dist\r2go.exe .\cmd\r2go || exit /b 1

echo DONE: %CD%\dist\r2go.exe
endlocal
