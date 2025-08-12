# Crontab Go

一个基于 Go 语言的定时任务管理系统，提供 Web 界面和 RESTful API 来管理定时任务和监控系统状态。

## 功能特性

- **任务管理**
  - 创建、查看、更新和删除定时任务
  - 支持命令执行和 HTTP 请求两种任务类型
  - 任务执行日志记录
  - 手动触发任务执行

- **系统监控**
  - 实时监控系统资源使用情况（CPU、内存、系统负载）
  - 自动收集和存储系统统计数据

- **Web 界面**
  - 直观的任务管理界面
  - 任务执行日志查看
  - 系统监控数据展示

## 技术栈

- **后端**: Go, Gin, GORM
- **数据库**: SQLite
- **前端**: HTML, CSS, JavaScript

## 项目结构

```
crontab_go/
├── cmd/                 # 应用程序入口
│   └── main.go
├── internal/            # 内部包
│   ├── application/     # 应用层逻辑
│   ├── domain/          # 领域层
│   │   ├── entity/      # 实体定义
│   │   └── service/     # 领域服务
│   ├── infrastructure/  # 基础设施层
│   │   ├── http/        # HTTP 服务器和处理器
│   │   └── persistence/ # 数据持久化
│   └── interfaces/      # 接口层
├── web/                 # Web 前端资源
│   ├── index.html
│   └── static/          # 静态资源
└── docs/                # 文档
    └── API.md           # API 文档
```

## 安装与运行

### 前提条件

- Go 1.16 或更高版本

### 安装依赖

```bash
go mod tidy
```

### 运行应用

```bash
go run cmd/main.go
```

应用将在 `http://localhost:8080` 启动。

## API 文档

详细的 API 文档请参考 [API.md](docs/API.md)。

## 数据模型

### 任务 (Task)

- **id**: 任务ID
- **name**: 任务名称
- **schedule**: Cron表达式，定义任务的执行计划
- **command**: 要执行的命令或URL
- **method**: HTTP请求方法 (可选，默认为GET)
- **headers**: JSON格式的请求头 (可选)
- **enabled**: 任务是否启用 (可选，默认为true)
- **description**: 任务描述 (可选)

### 任务日志 (TaskLog)

- **ID**: 日志ID
- **TaskID**: 关联的任务ID
- **TaskName**: 任务名称（冗余存储，便于查询）
- **StartTime**: 任务开始执行时间
- **EndTime**: 任务执行结束时间
- **Success**: 执行是否成功
- **Output**: 任务输出
- **Error**: 错误信息（如果有的话）

### 系统统计 (SystemStats)

- **ID**: 统计ID
- **CPUUsage**: CPU使用率
- **MemoryUsage**: 内存使用率
- **SystemLoad**: 系统负载
- **Timestamp**: 记录时间戳

## 贡献指南

欢迎提交问题和拉取请求来改进此项目。

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。
