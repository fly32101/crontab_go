# Crontab Go

[![Build and Release](https://github.com/your-username/crontab_go/actions/workflows/build-and-release.yml/badge.svg)](https://github.com/your-username/crontab_go/actions/workflows/build-and-release.yml)
[![Docker Image](https://ghcr-badge.deta.dev/your-username/crontab_go/latest_tag?trim=major&label=Docker%20Image)](https://github.com/your-username/crontab_go/pkgs/container/crontab_go)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/crontab_go)](https://goreportcard.com/report/github.com/your-username/crontab_go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个基于 Go 语言的现代化定时任务管理系统，提供美观的 Web 界面和完整的 RESTful API 来管理定时任务和监控系统状态。

## ✨ 功能特性

### 🔐 用户认证与权限管理

- JWT 基于令牌的身份验证
- 用户注册和登录系统
- 角色权限控制（管理员/普通用户）
- 安全的密码加密存储
- 自动登录状态检查

### 📋 任务管理

- 创建、查看、更新和删除定时任务
- 支持 Cron 表达式调度
- 支持命令执行和 HTTP 请求两种任务类型
- 任务执行日志记录与分页查看
- 手动触发任务执行
- 任务启用/禁用控制
- 实时任务状态监控

### 📊 系统监控

- 实时系统资源监控（CPU、内存、磁盘、系统负载）
- 美观的监控仪表盘
- 自动数据收集和存储
- 历史数据管理（保留最新100条记录）
- 响应式监控卡片设计

### 🎨 现代化界面

- 基于 Element Plus 的现代化 UI 设计
- 响应式布局，支持移动端
- Font Awesome 图标库
- 深色主题配色
- 直观的用户交互体验

### 🔍 高级功能

- 任务搜索和筛选
- 分页数据展示
- 实时数据更新
- 错误处理和用户反馈
- API 文档完整

## 🛠 技术栈

- **后端**: Go 1.21+, Gin Web Framework, GORM ORM
- **数据库**: SQLite
- **前端**: Vue 3, Element Plus, Font Awesome
- **认证**: JWT (JSON Web Tokens)
- **系统监控**: gopsutil
- **密码加密**: bcrypt

## 📁 项目结构

```
crontab_go/
├── cmd/                          # 应用程序入口
│   └── main.go                   # 主程序入口
├── internal/                     # 内部包（私有代码）
│   ├── application/              # 应用层逻辑
│   │   ├── auth/                 # 认证服务
│   │   ├── system/               # 系统监控服务
│   │   └── task/                 # 任务管理服务
│   ├── domain/                   # 领域层
│   │   ├── entity/               # 实体定义
│   │   │   ├── pagination.go     # 分页实体
│   │   │   ├── task.go          # 任务实体
│   │   │   ├── task_log.go      # 任务日志实体
│   │   │   ├── system.go        # 系统监控实体
│   │   │   └── user.go          # 用户实体
│   │   ├── repository/           # 仓库接口
│   │   └── service/             # 领域服务
│   ├── infrastructure/          # 基础设施层
│   │   └── persistence/         # 数据持久化
│   │       ├── sqlite.go        # SQLite 数据库配置
│   │       ├── task_repository.go
│   │       ├── user_repository.go
│   │       └── system_repository.go
│   └── interfaces/              # 接口层
│       └── http/                # HTTP 接口
│           ├── handler.go       # 请求处理器
│           ├── middleware.go    # 中间件
│           └── server.go        # 服务器配置
├── web/                         # Web 前端资源
│   ├── index.html              # 主页面
│   ├── login.html              # 登录页面
│   └── static/                 # 静态资源
│       └── app.js              # 前端 JavaScript
└── docs/                       # 文档
    └── API.md                  # API 文档
```

## 🚀 快速开始

### 📥 二进制部署（推荐）

1. **下载二进制文件**
   
   从 [Releases](https://github.com/your-username/crontab_go/releases) 页面下载对应平台的二进制文件：
   - Linux AMD64: `crontab-go-linux-amd64`
   - Linux ARM64: `crontab-go-linux-arm64`
   - Windows AMD64: `crontab-go-windows-amd64.exe`
   - macOS AMD64: `crontab-go-darwin-amd64`
   - macOS ARM64: `crontab-go-darwin-arm64`

2. **给执行权限**（Linux/macOS）
   ```bash
   chmod +x crontab-go-*
   ```

3. **运行应用**
   ```bash
   # Linux
   ./crontab-go-linux-amd64
   
   # macOS
   ./crontab-go-darwin-amd64
   
   # Windows
   crontab-go-windows-amd64.exe
   ```

4. **访问应用**
   - 应用将在 `http://localhost:8080` 启动
   - 首次访问会自动跳转到登录页面

### 📦 源码部署

#### 前提条件

- Go 1.21 或更高版本
- Node.js 18 或更高版本
- Git（用于克隆项目）

#### 安装与运行

1. **克隆项目**

   ```bash
   git clone <repository-url>
   cd crontab_go
   ```

2. **构建应用**

   ```bash
   # 使用构建脚本
   chmod +x scripts/build.sh
   ./scripts/build.sh
   
   # 或手动构建
   cd web && npm install && npm run build && cd ..
   go mod tidy
   go build -o crontab_go ./cmd
   ```

3. **运行应用**

   ```bash
   ./crontab_go
   ```

### 🐳 Docker 部署（可选）

如果你熟悉 Docker，也可以使用 Docker 部署：

#### 使用 Docker Compose

1. **下载配置文件**
   ```bash
   curl -O https://raw.githubusercontent.com/your-username/crontab_go/main/docker-compose.yml
   ```

2. **启动服务**
   ```bash
   docker-compose up -d
   ```

#### 自行构建镜像

```bash
# 克隆项目
git clone <repository-url>
cd crontab_go

# 构建镜像
docker build -t crontab-go:latest .

# 运行容器
docker run -d \
  --name crontab-go \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e TZ=Asia/Shanghai \
  crontab-go:latest
```

### 🔑 默认账户

系统会自动创建默认管理员账户：

- **用户名**: `admin`
- **密码**: `admin123`
- **角色**: 管理员

> ⚠️ **安全提示**: 首次登录后请及时修改默认密码

## 🐳 Docker 部署

### 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `DB_PATH` | `/app/data/crontab.db` | 数据库文件路径 |
| `JWT_SECRET` | 随机生成 | JWT 密钥（生产环境请设置） |
| `GIN_MODE` | `release` | Gin 运行模式 |
| `TZ` | `Asia/Shanghai` | 时区设置 |

### Docker Compose 配置示例

```yaml
version: '3.8'
services:
  crontab-go:
    image: ghcr.io/your-username/crontab_go:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
      - JWT_SECRET=your-secret-key
    restart: unless-stopped
```

### 数据持久化

- 数据库文件：`/app/data/crontab.db`
- 建议挂载 `/app/data` 目录到宿主机
- 支持定期备份和恢复

### 健康检查

容器内置健康检查，检查端点：`/api/v1/system/stats`

```bash
# 检查容器健康状态
docker ps
docker inspect crontab-go | grep Health
```

## 📖 使用指南

### 🔐 登录系统

1. 访问 `http://localhost:8080` 会自动跳转到登录页面
2. 使用默认账户 `admin/admin123` 登录
3. 也可以点击"立即注册"创建新账户
4. 登录成功后会跳转到主界面

### 📋 任务管理

1. **创建任务**
   - 点击"添加任务"按钮
   - 填写任务名称、Cron表达式、命令等信息
   - 支持系统命令和HTTP请求两种类型

2. **管理任务**
   - 查看任务列表和状态
   - 启用/禁用任务
   - 编辑任务配置
   - 手动执行任务
   - 查看执行日志

3. **Cron表达式示例**

   ```
   0 */5 * * * *    # 每5分钟执行一次
   0 0 12 * * *     # 每天12点执行
   0 0 0 * * 1      # 每周一执行
   0 30 9 * * 1-5   # 工作日9:30执行
   ```

### 📊 系统监控

- 实时查看CPU、内存、磁盘使用率
- 监控系统负载和运行时间
- 自动30秒刷新一次数据

### 🔍 高级功能

- **搜索**: 在任务列表中搜索任务名称或命令
- **筛选**: 按任务状态筛选（启用/禁用）
- **分页**: 支持大量任务的分页显示
- **日志**: 查看详细的任务执行日志

## 📚 API 文档

详细的 API 文档请参考 [API.md](docs/API.md)。

### 主要 API 端点

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `GET /api/v1/tasks/paginated` - 获取任务列表（分页）
- `POST /api/v1/tasks` - 创建任务
- `GET /api/v1/system/stats` - 获取系统监控数据

## 🗄 数据模型

### 👤 用户 (User)

- **ID**: 用户ID
- **Username**: 用户名（唯一）
- **Password**: 加密密码
- **Email**: 邮箱地址
- **Role**: 用户角色（admin/user）
- **IsActive**: 是否激活
- **CreatedAt/UpdatedAt**: 创建/更新时间

### 📋 任务 (Task)

- **ID**: 任务ID
- **Name**: 任务名称
- **Schedule**: Cron表达式，定义任务的执行计划
- **Command**: 要执行的命令或URL
- **Method**: HTTP请求方法（可选，默认为GET）
- **Headers**: JSON格式的请求头（可选）
- **Enabled**: 任务是否启用（默认为true）
- **Description**: 任务描述（可选）

### 📝 任务日志 (TaskLog)

- **ID**: 日志ID
- **TaskID**: 关联的任务ID
- **TaskName**: 任务名称（冗余存储，便于查询）
- **StartTime**: 任务开始执行时间
- **EndTime**: 任务执行结束时间
- **Success**: 执行是否成功
- **Output**: 任务输出
- **Error**: 错误信息（如果有的话）

### 📊 系统统计 (SystemStats)

- **ID**: 统计ID
- **CPUUsage**: CPU使用率
- **MemoryUsage**: 内存使用率
- **MemoryTotal/Used/Free**: 内存详细信息
- **DiskUsage**: 磁盘使用率
- **DiskTotal/Used/Free**: 磁盘详细信息
- **SystemLoad**: 系统负载
- **NetworkRxBytes/TxBytes**: 网络流量
- **ProcessCount**: 进程数量
- **GoroutineCount**: Goroutine数量
- **Uptime**: 系统运行时间
- **Timestamp**: 记录时间戳

## 🔧 配置说明

### 环境变量

- `JWT_SECRET`: JWT密钥（生产环境建议设置）
- `DB_PATH`: 数据库文件路径（默认：crontab.db）
- `PORT`: 服务端口（默认：8080）

### 数据库

- 使用 SQLite 作为默认数据库
- 自动创建和迁移数据库表
- 数据文件：`crontab.db`

### 系统监控

- 每10秒收集一次系统数据
- 每分钟清理一次旧数据
- 保留最新100条历史记录

## 🔒 安全特性

- **密码加密**: 使用 bcrypt 算法加密存储密码
- **JWT 认证**: 基于 JSON Web Token 的无状态认证
- **权限控制**: 基于角色的访问控制（RBAC）
- **自动登出**: Token 过期自动登出
- **CORS 支持**: 跨域请求支持

## 🚀 性能优化

- **分页查询**: 大数据量分页加载
- **实时监控**: 高效的系统资源监控
- **数据清理**: 自动清理历史数据，避免数据库膨胀
- **连接池**: 数据库连接池优化
- **静态资源**: CDN 加速的前端资源

## 🐛 故障排除

### 常见问题

1. **无法启动服务**
   - 检查端口8080是否被占用
   - 确认Go版本是否为1.21+

2. **登录失败**
   - 确认使用默认账户 admin/admin123
   - 检查数据库文件是否正常创建

3. **任务不执行**
   - 检查Cron表达式格式是否正确
   - 确认任务是否已启用
   - 查看任务执行日志

4. **系统监控数据异常**
   - 确认gopsutil库是否正常工作
   - 检查系统权限

### 日志查看

```bash
# 运行时会输出详细日志
go run cmd/main.go

# 查看任务执行日志
# 在Web界面中点击"日志"按钮查看
```

### 📢 通知功能

- 支持任务执行结果通知（邮件、钉钉、企业微信）
- 可配置通知时机（成功时通知、失败时通知）
- 支持多种通知方式同时使用
- 通知配置测试功能
- 详细的通知内容（任务信息、执行状态、输出结果等）

详细配置说明请参考 [通知功能文档](docs/NOTIFICATION.md)。

### 📊 统计报表

- 任务执行统计（成功率、执行次数、平均执行时间）
- 执行趋势分析（日趋势、小时分布）
- 性能指标监控（最短/最长/平均执行时间）
- 可视化图表展示（基于 ECharts）
- 多维度数据筛选和分析

### 📋 任务模板

- 预定义任务模板库
- 模板分类管理（备份、监控、清理、通知等）
- 公共模板和私有模板
- 模板搜索和筛选功能
- 热门模板推荐
- 从模板快速创建任务
- 模板使用统计和标签系统

详细使用说明请参考 [任务模板文档](docs/TASK_TEMPLATES.md)。

## 🛣 开发路线图

- [x] 任务执行结果通知（邮件、钉钉、企业微信）
- [x] 任务执行统计和报表
- [x] 任务模板功能
- [x] Docker 容器化部署
- [x] GitHub Actions CI/CD
- [ ] 任务依赖关系管理
- [ ] 更多系统监控指标
- [ ] 多用户权限细化
- [ ] 集群模式支持

## 🤝 贡献指南

欢迎提交问题和拉取请求来改进此项目！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [GORM](https://gorm.io/) - Go ORM 库
- [Element Plus](https://element-plus.org/) - Vue 3 UI 组件库
- [gopsutil](https://github.com/shirou/gopsutil) - 系统信息库
- [cron](https://github.com/robfig/cron) - Cron 表达式解析库

---

⭐ 如果这个项目对你有帮助，请给它一个星标！
