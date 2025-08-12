#!/bin/bash

# 构建脚本
set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    print_message "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed"
        exit 1
    fi
    
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed"
        exit 1
    fi
}

# 构建前端
build_frontend() {
    print_message "Building frontend..."
    cd web
    npm ci
    npm run build
    cd ..
}

# 构建后端
build_backend() {
    print_message "Building backend..."
    go mod tidy
    go build -ldflags="-s -w" -o crontab_go ./cmd
}

# 主函数
main() {
    check_dependencies
    build_frontend
    build_backend
    print_message "Build completed successfully!"
}

main "$@"