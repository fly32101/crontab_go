## 🚀 Release v{VERSION}

### ✨ 新功能
- 

### 🐛 Bug 修复
- 

### 🔧 改进
- 

### 📦 Docker 镜像

```bash
# 拉取最新镜像
docker pull ghcr.io/your-username/crontab_go:{VERSION}

# 使用 Docker Compose
curl -O https://raw.githubusercontent.com/your-username/crontab_go/{VERSION}/docker-compose.yml
docker-compose up -d
```

### 📥 二进制下载

| 平台 | 架构 | 下载链接 |
|------|------|----------|
| Linux | AMD64 | [crontab-go-linux-amd64](https://github.com/your-username/crontab_go/releases/download/{VERSION}/crontab-go-linux-amd64) |
| Linux | ARM64 | [crontab-go-linux-arm64](https://github.com/your-username/crontab_go/releases/download/{VERSION}/crontab-go-linux-arm64) |
| Windows | AMD64 | [crontab-go-windows-amd64.exe](https://github.com/your-username/crontab_go/releases/download/{VERSION}/crontab-go-windows-amd64.exe) |
| macOS | AMD64 | [crontab-go-darwin-amd64](https://github.com/your-username/crontab_go/releases/download/{VERSION}/crontab-go-darwin-amd64) |
| macOS | ARM64 | [crontab-go-darwin-arm64](https://github.com/your-username/crontab_go/releases/download/{VERSION}/crontab-go-darwin-arm64) |

### 🔄 升级指南

#### Docker 用户

```bash
# 停止现有服务
docker-compose down

# 更新镜像版本
sed -i 's/ghcr.io\/your-username\/crontab_go:.*/ghcr.io\/your-username\/crontab_go:{VERSION}/' docker-compose.yml

# 启动新版本
docker-compose up -d
```

#### 二进制用户

1. 下载新版本二进制文件
2. 停止现有服务
3. 替换二进制文件
4. 启动新版本

### 📖 文档

- [安装指南](https://github.com/your-username/crontab_go/blob/{VERSION}/README.md)
- [Docker 部署](https://github.com/your-username/crontab_go/blob/{VERSION}/docs/DOCKER.md)
- [API 文档](https://github.com/your-username/crontab_go/blob/{VERSION}/docs/API.md)
- [通知配置](https://github.com/your-username/crontab_go/blob/{VERSION}/docs/NOTIFICATION.md)
- [任务模板](https://github.com/your-username/crontab_go/blob/{VERSION}/docs/TASK_TEMPLATES.md)

### 🔗 相关链接

- [完整更新日志](https://github.com/your-username/crontab_go/blob/{VERSION}/CHANGELOG.md)
- [问题反馈](https://github.com/your-username/crontab_go/issues)
- [讨论区](https://github.com/your-username/crontab_go/discussions)