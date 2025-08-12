# Makefile for Crontab Go

.PHONY: help build build-frontend build-backend clean run dev docker-build docker-run docker-stop

# 默认目标
help:
	@echo "Available targets:"
	@echo "  build          - Build the complete application"
	@echo "  build-frontend - Build frontend only"
	@echo "  build-backend  - Build backend only"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode with hot reload"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  docker-dev     - Run development environment with Docker"

# 构建完整应用
build: build-frontend build-backend

# 构建前端
build-frontend:
	@echo "Building frontend..."
	cd web && npm ci && npm run build

# 构建后端
build-backend:
	@echo "Building backend..."
	go mod tidy
	go build -ldflags="-s -w" -o crontab_go ./cmd

# 清理构建产物
clean:
	@echo "Cleaning build artifacts..."
	rm -f crontab_go crontab_go.exe
	rm -rf web/dist
	rm -rf tmp

# 运行应用
run: build
	@echo "Starting application..."
	./crontab_go

# 开发模式（需要安装 air）
dev:
	@echo "Starting development mode..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Docker 构建
docker-build:
	@echo "Building Docker image..."
	docker build -t crontab-go:latest .

# Docker 运行
docker-run:
	@echo "Starting with Docker Compose..."
	docker-compose up -d

# 停止 Docker 容器
docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down

# Docker 开发环境
docker-dev:
	@echo "Starting development environment with Docker..."
	docker-compose -f docker-compose.dev.yml up -d

# 检查依赖
check-deps:
	@echo "Checking dependencies..."
	@command -v go >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "Node.js is not installed"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "npm is not installed"; exit 1; }
	@echo "All dependencies are available"

# 格式化代码
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting frontend code..."
	cd web && npm run format 2>/dev/null || echo "Frontend format command not available"

# 代码检查
lint:
	@echo "Running Go vet..."
	go vet ./...
	@echo "Running frontend lint..."
	cd web && npm run lint 2>/dev/null || echo "Frontend lint command not available"

# 安装依赖
install-deps:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd web && npm install