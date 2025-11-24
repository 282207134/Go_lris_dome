#!/bin/bash

# Iris Go 框架学习项目启动脚本

echo "=== Iris Go 框架学习项目启动脚本 ==="

# 检查 Go 版本
echo "检查 Go 版本..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo "Go 版本: $GO_VERSION"
else
    echo "错误: 未找到 Go，请先安装 Go 1.21 或更高版本"
    exit 1
fi

# 检查项目目录
if [ ! -f "go.mod" ]; then
    echo "错误: 请在项目根目录运行此脚本"
    exit 1
fi

# 安装依赖
echo "安装项目依赖..."
go mod download
if [ $? -ne 0 ]; then
    echo "错误: 依赖安装失败"
    exit 1
fi

# 创建必要的目录
echo "创建必要的目录..."
mkdir -p static/uploads
mkdir -p data
mkdir -p logs
mkdir -p temp

# 设置环境变量
export GIN_MODE=debug
export SERVER_PORT=8080
export DB_DRIVER=sqlite
export DB_DATABASE=./data/iris_sample.db
export JWT_SECRET=your-secret-key-change-this-in-production
export JWT_EXPIRATION_TIME=86400
export LOG_LEVEL=info

echo "环境变量设置完成"
echo "  GIN_MODE=$GIN_MODE"
echo "  SERVER_PORT=$SERVER_PORT"
echo "  DB_DRIVER=$DB_DRIVER"
echo "  DB_DATABASE=$DB_DATABASE"
echo "  LOG_LEVEL=$LOG_LEVEL"

# 检查端口是否被占用
if lsof -Pi :$SERVER_PORT -sTCP:LISTEN -t >/dev/null ; then
    echo "警告: 端口 $SERVER_PORT 已被占用"
    echo "请检查是否有其他程序在使用该端口"
fi

# 启动应用
echo ""
echo "启动 Iris Go 框架学习项目..."
echo "访问地址: http://localhost:$SERVER_PORT"
echo "API 文档: http://localhost:$SERVER_PORT/api/docs"
echo "健康检查: http://localhost:$SERVER_PORT/api/health"
echo ""
echo "按 Ctrl+C 停止应用"
echo ""

go run main.go