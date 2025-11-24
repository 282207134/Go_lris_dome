@echo off
REM Iris Go 框架学习项目启动脚本 (Windows)

echo === Iris Go 框架学习项目启动脚本 ===

REM 检查 Go 版本
echo 检查 Go 版本...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到 Go，请先安装 Go 1.21 或更高版本
    pause
    exit /b 1
)

REM 检查项目目录
if not exist "go.mod" (
    echo 错误: 请在项目根目录运行此脚本
    pause
    exit /b 1
)

REM 安装依赖
echo 安装项目依赖...
go mod download
if %errorlevel% neq 0 (
    echo 错误: 依赖安装失败
    pause
    exit /b 1
)

REM 创建必要的目录
echo 创建必要的目录...
if not exist "static\uploads" mkdir static\uploads
if not exist "data" mkdir data
if not exist "logs" mkdir logs
if not exist "temp" mkdir temp

REM 设置环境变量
set GIN_MODE=debug
set SERVER_PORT=8080
set DB_DRIVER=sqlite
set DB_DATABASE=./data/iris_sample.db
set JWT_SECRET=your-secret-key-change-this-in-production
set JWT_EXPIRATION_TIME=86400
set LOG_LEVEL=info

echo 环境变量设置完成
echo   GIN_MODE=%GIN_MODE%
echo   SERVER_PORT=%SERVER_PORT%
echo   DB_DRIVER=%DB_DRIVER%
echo   DB_DATABASE=%DB_DATABASE%
echo   LOG_LEVEL=%LOG_LEVEL%

REM 启动应用
echo.
echo 启动 Iris Go 框架学习项目...
echo 访问地址: http://localhost:%SERVER_PORT%
echo API 文档: http://localhost:%SERVER_PORT%/api/docs
echo 健康检查: http://localhost:%SERVER_PORT%/api/health
echo.
echo 按 Ctrl+C 停止应用
echo.

go run main.go

pause