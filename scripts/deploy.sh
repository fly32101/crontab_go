#!/bin/bash

# 部署脚本
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

# 创建必要的目录
create_directories() {
    print_message "Creating necessary directories..."
    mkdir -p data
    mkdir -p logs
    mkdir -p ssl
}

# 生成自签名 SSL 证书（仅用于开发）
generate_ssl_cert() {
    if [ ! -f ssl/cert.pem ] || [ ! -f ssl/key.pem ]; then
        print_message "Generating self-signed SSL certificate..."
        openssl req -x509 -newkey rsa:4096 -keyout ssl/key.pem -out ssl/cert.pem -days 365 -nodes \
            -subj "/C=CN/ST=State/L=City/O=Organization/CN=localhost"
        print_warning "Using self-signed certificate. For production, please use a valid SSL certificate."
    fi
}

# 部署应用
deploy() {
    print_message "Starting deployment..."
    
    # 停止现有容器
    print_message "Stopping existing containers..."
    docker-compose down
    
    # 拉取最新镜像
    print_message "Pulling latest images..."
    docker-compose pull
    
    # 启动服务
    print_message "Starting services..."
    docker-compose up -d
    
    # 等待服务启动
    print_message "Waiting for services to start..."
    sleep 10
    
    # 检查服务状态
    if docker-compose ps | grep -q "Up"; then
        print_message "Deployment successful!"
        print_message "Application is running at:"
        print_message "  HTTP:  http://localhost"
        print_message "  HTTPS: https://localhost"
        print_message "  Direct: http://localhost:8080"
    else
        print_error "Deployment failed. Please check the logs:"
        docker-compose logs
        exit 1
    fi
}

# 显示日志
show_logs() {
    print_message "Showing application logs..."
    docker-compose logs -f
}

# 停止服务
stop_services() {
    print_message "Stopping services..."
    docker-compose down
}

# 清理
cleanup() {
    print_message "Cleaning up..."
    docker-compose down -v
    docker system prune -f
}

# 备份数据
backup_data() {
    print_message "Backing up data..."
    timestamp=$(date +%Y%m%d_%H%M%S)
    backup_dir="backup_${timestamp}"
    mkdir -p "${backup_dir}"
    
    if [ -d "data" ]; then
        cp -r data "${backup_dir}/"
        print_message "Data backed up to ${backup_dir}/"
    else
        print_warning "No data directory found to backup."
    fi
}

# 恢复数据
restore_data() {
    if [ -z "$1" ]; then
        print_error "Please specify backup directory: $0 restore <backup_directory>"
        exit 1
    fi
    
    backup_dir="$1"
    if [ ! -d "${backup_dir}" ]; then
        print_error "Backup directory ${backup_dir} not found."
        exit 1
    fi
    
    print_message "Restoring data from ${backup_dir}..."
    if [ -d "${backup_dir}/data" ]; then
        rm -rf data
        cp -r "${backup_dir}/data" .
        print_message "Data restored successfully."
    else
        print_error "No data directory found in backup."
        exit 1
    fi
}

# 显示帮助信息
show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  deploy    Deploy the application"
    echo "  logs      Show application logs"
    echo "  stop      Stop all services"
    echo "  cleanup   Stop services and clean up"
    echo "  backup    Backup application data"
    echo "  restore   Restore application data from backup"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 deploy"
    echo "  $0 logs"
    echo "  $0 backup"
    echo "  $0 restore backup_20231201_120000"
}

# 主函数
main() {
    case "${1:-deploy}" in
        "deploy")
            check_docker
            create_directories
            generate_ssl_cert
            deploy
            ;;
        "logs")
            show_logs
            ;;
        "stop")
            stop_services
            ;;
        "cleanup")
            cleanup
            ;;
        "backup")
            backup_data
            ;;
        "restore")
            restore_data "$2"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"