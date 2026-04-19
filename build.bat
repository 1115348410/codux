@echo off
echo ===========================================
echo   Codux Windows 一键构建脚本
echo ===========================================
echo.

cd /d "%~dp0Sources\DmuxWorkspace.WinUI"

echo [1/3] 正在还原依赖...
dotnet restore
if errorlevel 1 (
    echo 还原失败！按任意键退出...
    pause
    exit /b 1
)

echo.
echo [2/3] 正在构建并发布...
dotnet publish -c Release -r win-x64 --self-contained true -p:PublishSingleFile=true -p:IncludeNativeLibrariesForSelfExtract=true -o ./publish
if errorlevel 1 (
    echo 构建失败！按任意键退出...
    pause
    exit /b 1
)

echo.
echo [3/3] 构建完成！
echo.
echo ===========================================
echo   发布成功！
echo.
echo   exe 文件位于:
echo   %~dp0Sources\DmuxWorkspace.WinUI\publish\Codux.WinUI.exe
echo.
echo   双击 Codux.WinUI.exe 即可运行！
echo ===========================================
echo.
pause
