# Iris Go 框架学习项目 Makefile

.PHONY: help build run test clean dev deps fmt lint

# 默认目标
help:
	@echo "可用的命令："
	@echo "  deps     - 安装依赖"
	@echo "  dev      - 开发模式运行（热重载）"
	@echo "  run      - 运行应用"
	@echo "  build    - 构建应用"
	@echo "  test     - 运行测试"
	@echo "  fmt      - 格式化代码"
	@echo "  lint     - 代码检查"
	@echo "  clean    - 清理构建文件"

# 安装依赖
deps:
	go mod download
	go mod tidy

# 开发模式运行（需要安装 air）
dev:
	@if ! command -v air &> /dev/null; then \
		echo "安装 air 工具..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# 运行应用
run:
	go run main.go

# 构建应用
build:
	go build -o bin/iris-sample main.go

# 运行测试
test:
	go test -v ./...

# 格式化代码
fmt:
	go fmt ./...

# 代码检查（需要安装 golangci-lint）
lint:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "安装 golangci-lint..."; \
		go install github.com/golang-org/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

# 清理构建文件
clean:
	rm -rf bin/
	rm -f *.log
	rm -f *.db
	rm -rf uploads/
	rm -rf temp/

# 生成文档（需要安装 swag）
docs:
	@if ! command -v swag &> /dev/null; then \
		echo "安装 swag 工具..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init

# 数据库迁移
migrate:
	go run main.go -migrate

# 创建示例数据
seed:
	go run main.go -seed

# 生产环境构建
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/iris-sample-linux main.go
	CGO_ENABLED=0 GOOS=windows go build -ldflags="-w -s" -o bin/iris-sample.exe main.go

# Docker 构建
docker-build:
	docker build -t iris-sample .

# Docker 运行
docker-run:
	docker run -p 8080:8080 iris-sample

# 安装开发工具
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang-org/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest