@echo off
setlocal enabledelayedexpansion

:: 部署脚本 (Windows 版本)
echo [INFO] Starting deployment...

:: 检查 Docker 是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker is not installed. Please install Docker first.
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Docker Compose is not installed. Please install Docker Compose first.
    exit /b 1
)

:: 创建必要的目录
echo [INFO] Creating necessary directories...
if not exist "data" mkdir data
if not exist "logs" mkdir logs
if not exist "ssl" mkdir ssl

:: 停止现有容器
echo [INFO] Stopping existing containers...
docker-compose down

:: 拉取最新镜像
echo [INFO] Pulling latest images...
docker-compose pull

:: 启动服务
echo [INFO] Starting services...
docker-compose up -d

:: 等待服务启动
echo [INFO] Waiting for services to start...
timeout /t 10 /nobreak >nul

:: 检查服务状态
docker-compose ps | findstr "Up" >nul
if errorlevel 1 (
    echo [ERROR] Deployment failed. Please check the logs:
    docker-compose logs
    exit /b 1
) else (
    echo [INFO] Deployment successful!
    echo [INFO] Application is running at:
    echo [INFO]   HTTP:  http://localhost
    echo [INFO]   HTTPS: https://localhost
    echo [INFO]   Direct: http://localhost:8080
)

endlocal