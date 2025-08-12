# Docker 部署指南

本文档详细介绍如何使用 Docker 部署 Crontab Go 管理系统。

## 🐳 Docker 镜像

### 官方镜像

```bash
# 最新版本
ghcr.io/your-username/crontab_go:latest

# 特定版本
ghcr.io/your-username/crontab_go:v1.0.0

# 开发版本
ghcr.io/your-username/crontab_go:main
```

### 支持的架构

- `linux/amd64` - x86_64 架构
- `linux/arm64` - ARM64 架构（如 Apple M1/M2）

## 🚀 快速开始

### 1. 使用 Docker Compose（推荐）

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'

services:
  crontab-go:
    image: ghcr.io/your-username/crontab_go:latest
    container_name: crontab-go
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - /var/log:/var/log:ro  # 可选：挂载日志目录
    environment:
      - TZ=Asia/Shanghai
      - DB_PATH=/app/data/crontab.db
      - JWT_SECRET=your-jwt-secret-key-change-this
      - GIN_MODE=release
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/api/v1/system/stats"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

启动服务：

```bash
docker-compose up -d
```

### 2. 使用 Docker 命令

```bash
# 创建数据目录
mkdir -p data

# 运行容器
docker run -d \
  --name crontab-go \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  -e JWT_SECRET=your-secret-key \
  --restart unless-stopped \
  ghcr.io/your-username/crontab_go:latest
```

## ⚙️ 配置选项

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DB_PATH` | `/app/data/crontab.db` | SQLite 数据库文件路径 |
| `JWT_SECRET` | 随机生成 | JWT 签名密钥（生产环境必须设置） |
| `GIN_MODE` | `release` | Gin 框架运行模式 (`debug`/`release`) |
| `TZ` | `Asia/Shanghai` | 容器时区 |
| `PORT` | `8080` | 服务监听端口 |

### 数据卷挂载

| 容器路径 | 说明 | 推荐挂载 |
|----------|------|----------|
| `/app/data` | 数据库和配置文件 | `./data:/app/data` |
| `/var/log` | 系统日志（只读） | `/var/log:/var/log:ro` |

### 端口映射

| 容器端口 | 说明 |
|----------|------|
| `8080` | HTTP 服务端口 |

## 🔧 高级配置

### 1. 使用 Nginx 反向代理

创建包含 Nginx 的 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  crontab-go:
    image: ghcr.io/your-username/crontab_go:latest
    container_name: crontab-go
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
      - JWT_SECRET=your-jwt-secret-key
    restart: unless-stopped
    networks:
      - crontab-network

  nginx:
    image: nginx:alpine
    container_name: crontab-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - crontab-go
    restart: unless-stopped
    networks:
      - crontab-network

networks:
  crontab-network:
    driver: bridge
```

### 2. SSL/TLS 配置

生成自签名证书（仅用于测试）：

```bash
mkdir -p ssl
openssl req -x509 -newkey rsa:4096 -keyout ssl/key.pem -out ssl/cert.pem -days 365 -nodes \
  -subj "/C=CN/ST=State/L=City/O=Organization/CN=localhost"
```

### 3. 日志配置

查看容器日志：

```bash
# 查看实时日志
docker-compose logs -f crontab-go

# 查看最近 100 行日志
docker-compose logs --tail=100 crontab-go
```

配置日志轮转：

```yaml
services:
  crontab-go:
    # ... 其他配置
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 🔒 安全配置

### 1. 非 root 用户运行

容器默认使用非 root 用户 `appuser` (UID: 1001) 运行，提高安全性。

### 2. 只读文件系统

```yaml
services:
  crontab-go:
    # ... 其他配置
    read_only: true
    tmpfs:
      - /tmp
    volumes:
      - ./data:/app/data
```

### 3. 资源限制

```yaml
services:
  crontab-go:
    # ... 其他配置
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

## 📊 监控和健康检查

### 健康检查

容器内置健康检查，检查应用是否正常运行：

```bash
# 检查健康状态
docker inspect crontab-go | grep -A 10 Health

# 手动执行健康检查
docker exec crontab-go wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/system/stats
```

### 监控指标

应用提供以下监控端点：

- `/api/v1/system/stats` - 系统统计信息
- `/api/v1/statistics/report` - 任务执行报表

## 🔄 备份和恢复

### 数据备份

```bash
# 停止容器
docker-compose stop crontab-go

# 备份数据
tar -czf backup-$(date +%Y%m%d_%H%M%S).tar.gz data/

# 启动容器
docker-compose start crontab-go
```

### 数据恢复

```bash
# 停止容器
docker-compose stop crontab-go

# 恢复数据
tar -xzf backup-20231201_120000.tar.gz

# 启动容器
docker-compose start crontab-go
```

### 自动备份脚本

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/crontab-backup-${DATE}.tar.gz"

# 创建备份目录
mkdir -p ${BACKUP_DIR}

# 备份数据
docker-compose exec -T crontab-go tar -czf - /app/data > ${BACKUP_FILE}

# 清理旧备份（保留最近 7 天）
find ${BACKUP_DIR} -name "crontab-backup-*.tar.gz" -mtime +7 -delete

echo "Backup completed: ${BACKUP_FILE}"
```

## 🚀 更新和升级

### 更新到最新版本

```bash
# 拉取最新镜像
docker-compose pull

# 重启服务
docker-compose up -d
```

### 更新到特定版本

```bash
# 修改 docker-compose.yml 中的镜像版本
# image: ghcr.io/your-username/crontab_go:v1.1.0

# 拉取指定版本
docker-compose pull

# 重启服务
docker-compose up -d
```

## 🐛 故障排除

### 常见问题

1. **容器启动失败**
   ```bash
   # 查看容器日志
   docker-compose logs crontab-go
   
   # 检查容器状态
   docker-compose ps
   ```

2. **数据库权限问题**
   ```bash
   # 检查数据目录权限
   ls -la data/
   
   # 修复权限
   sudo chown -R 1001:1001 data/
   ```

3. **端口冲突**
   ```bash
   # 检查端口占用
   netstat -tlnp | grep 8080
   
   # 修改端口映射
   # ports: - "8081:8080"
   ```

4. **内存不足**
   ```bash
   # 检查容器资源使用
   docker stats crontab-go
   
   # 增加内存限制
   # deploy.resources.limits.memory: 1G
   ```

### 调试模式

启用调试模式：

```yaml
services:
  crontab-go:
    # ... 其他配置
    environment:
      - GIN_MODE=debug
```

### 性能优化

1. **使用 SSD 存储**
2. **适当的内存分配**
3. **启用 Gzip 压缩**（通过 Nginx）
4. **配置适当的健康检查间隔**

## 📚 相关文档

- [安装指南](../README.md)
- [API 文档](API.md)
- [通知配置](NOTIFICATION.md)
- [任务模板](TASK_TEMPLATES.md)