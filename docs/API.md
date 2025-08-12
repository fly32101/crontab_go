# Crontab Go API 文档

## 概述

Crontab Go 是一个基于 Go 语言的定时任务管理系统，提供 RESTful API 用于管理定时任务和监控系统状态。

## API 端点

### 任务管理 API

所有任务相关的 API 都在 `/api/v1/tasks` 路径下。

#### 创建任务

- **URL**: `POST /api/v1/tasks`
- **描述**: 创建一个新的定时任务
- **请求体**:
  ```json
  {
    "name": "任务名称",
    "schedule": "cron表达式",
    "command": "要执行的命令或URL",
    "method": "HTTP请求方法 (可选，默认为GET)",
    "headers": "JSON格式的请求头 (可选)",
    "enabled": true,
    "description": "任务描述 (可选)"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "name": "任务名称",
    "schedule": "cron表达式",
    "command": "要执行的命令或URL",
    "method": "HTTP请求方法",
    "headers": "JSON格式的请求头",
    "enabled": true,
    "description": "任务描述"
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 请求体格式错误
  - 500: 服务器内部错误

#### 获取任务列表

- **URL**: `GET /api/v1/tasks`
- **描述**: 获取所有任务列表（不分页）

#### 分页获取任务列表

- **URL**: `GET /api/v1/tasks/paginated`
- **描述**: 分页获取任务列表
- **查询参数**:
  - `page`: 页码，从1开始 (可选，默认为1)
  - `page_size`: 每页大小 (可选，默认为10，最大100)
- **响应**:
  ```json
  {
    "page": 1,
    "page_size": 10,
    "total": 25,
    "total_pages": 3,
    "data": [
      {
        "id": 1,
        "name": "任务名称",
        "schedule": "cron表达式",
        "command": "要执行的命令或URL",
        "method": "HTTP请求方法",
        "headers": "JSON格式的请求头",
        "enabled": true,
        "description": "任务描述"
      }
    ]
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 请求参数错误
  - 500: 服务器内部错误
- **响应**:
  ```json
  [
    {
      "id": 1,
      "name": "任务名称",
      "schedule": "cron表达式",
      "command": "要执行的命令或URL",
      "method": "HTTP请求方法",
      "headers": "JSON格式的请求头",
      "enabled": true,
      "description": "任务描述"
    }
  ]
  ```
- **状态码**:
  - 200: 成功
  - 500: 服务器内部错误

#### 分页获取任务列表

- **URL**: `GET /api/v1/tasks/paginated`
- **描述**: 分页获取任务列表
- **查询参数**:
  - `page`: 页码，从1开始（可选，默认为1）
  - `page_size`: 每页大小（可选，默认为10，最大100）
- **响应**:
  ```json
  {
    "page": 1,
    "page_size": 10,
    "total": 25,
    "total_pages": 3,
    "data": [
      {
        "id": 1,
        "name": "任务名称",
        "schedule": "cron表达式",
        "command": "要执行的命令或URL",
        "method": "HTTP请求方法",
        "headers": "JSON格式的请求头",
        "enabled": true,
        "description": "任务描述"
      }
    ]
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 查询参数格式错误
  - 500: 服务器内部错误

#### 获取单个任务

- **URL**: `GET /api/v1/tasks/:id`
- **描述**: 获取指定ID的任务详情
- **参数**:
  - `id`: 任务ID (路径参数)
- **响应**:
  ```json
  {
    "id": 1,
    "name": "任务名称",
    "schedule": "cron表达式",
    "command": "要执行的命令或URL",
    "method": "HTTP请求方法",
    "headers": "JSON格式的请求头",
    "enabled": true,
    "description": "任务描述"
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 无效的任务ID
  - 404: 任务不存在
  - 500: 服务器内部错误

#### 更新任务

- **URL**: `PUT /api/v1/tasks/:id`
- **描述**: 更新指定ID的任务信息
- **参数**:
  - `id`: 任务ID (路径参数)
- **请求体**:
  ```json
  {
    "name": "任务名称",
    "schedule": "cron表达式",
    "command": "要执行的命令或URL",
    "method": "HTTP请求方法 (可选)",
    "headers": "JSON格式的请求头 (可选)",
    "enabled": true,
    "description": "任务描述 (可选)"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "name": "任务名称",
    "schedule": "cron表达式",
    "command": "要执行的命令或URL",
    "method": "HTTP请求方法",
    "headers": "JSON格式的请求头",
    "enabled": true,
    "description": "任务描述"
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 无效的任务ID或请求体格式错误
  - 500: 服务器内部错误

#### 删除任务

- **URL**: `DELETE /api/v1/tasks/:id`
- **描述**: 删除指定ID的任务
- **参数**:
  - `id`: 任务ID (路径参数)
- **响应**:
  ```json
  {
    "message": "Task deleted successfully"
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 无效的任务ID
  - 500: 服务器内部错误

#### 获取任务执行日志

- **URL**: `GET /api/v1/tasks/:id/logs`
- **描述**: 获取指定任务的执行日志（不分页）
- **参数**:
  - `id`: 任务ID (路径参数)
- **响应**:
  ```json
  [
    {
      "ID": 1,
      "TaskID": 1,
      "TaskName": "任务名称",
      "StartTime": "2023-01-01T12:00:00Z",
      "EndTime": "2023-01-01T12:00:05Z",
      "Success": true,
      "Output": "任务执行输出",
      "Error": ""
    }
  ]
  ```
- **状态码**:
  - 200: 成功
  - 400: 无效的任务ID
  - 500: 服务器内部错误

#### 分页获取任务执行日志

- **URL**: `GET /api/v1/tasks/:id/logs/paginated`
- **描述**: 分页获取指定任务的执行日志
- **参数**:
  - `id`: 任务ID (路径参数)
- **查询参数**:
  - `page`: 页码，从1开始（可选，默认为1）
  - `page_size`: 每页大小（可选，默认为10，最大100）
- **响应**:
  ```json
  {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_pages": 5,
    "data": [
      {
        "ID": 1,
        "TaskID": 1,
        "TaskName": "任务名称",
        "StartTime": "2023-01-01T12:00:00Z",
        "EndTime": "2023-01-01T12:00:05Z",
        "Success": true,
        "Output": "任务执行输出",
        "Error": ""
      }
    ]
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 无效的任务ID或查询参数格式错误
  - 500: 服务器内部错误

#### 立即执行任务

- **URL**: `POST /api/v1/tasks/:id/execute`
- **描述**: 立即执行指定ID的任务（不按照计划时间）
- **参数**:
  - `id`: 任务ID (路径参数)
- **响应**:
  ```json
  {
    "message": "Task executed successfully"
  }
  ```
- **状态码**:
  - 200: 成功
  - 400: 无效的任务ID
  - 500: 服务器内部错误

### 系统监控 API

所有系统监控相关的 API 都在 `/api/v1/system` 路径下。

#### 获取系统统计信息

- **URL**: `GET /api/v1/system/stats`
- **描述**: 获取最新的系统统计信息
- **响应**:
  ```json
  {
    "ID": 1,
    "CPUUsage": 25.5,
    "MemoryUsage": 45.2,
    "SystemLoad": 1.2,
    "Timestamp": "2023-01-01T12:00:00Z"
  }
  ```
- **状态码**:
  - 200: 成功
  - 500: 服务器内部错误

## 数据模型

### Task

| 字段 | 类型 | 描述 |
|------|------|------|
| id | int | 任务ID，主键 |
| name | string | 任务名称 |
| schedule | string | Cron表达式，定义任务的执行计划 |
| command | string | 要执行的命令或URL |
| method | string | HTTP请求方法 (可选，默认为GET) |
| headers | string | JSON格式的请求头 (可选) |
| enabled | bool | 任务是否启用 (可选，默认为true) |
| description | string | 任务描述 (可选) |

### TaskLog

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | uint | 日志ID，主键 |
| TaskID | int | 关联的任务ID |
| TaskName | string | 任务名称（冗余存储，便于查询） |
| StartTime | time.Time | 任务开始执行时间 |
| EndTime | time.Time | 任务执行结束时间 |
| Success | bool | 执行是否成功 |
| Output | string | 任务输出 |
| Error | string | 错误信息（如果有的话） |

### SystemStats

| 字段 | 类型 | 描述 |
|------|------|------|
| ID | uint | 统计ID，主键 |
| CPUUsage | float64 | CPU使用率 |
| MemoryUsage | float64 | 内存使用率 |
| SystemLoad | float64 | 系统负载 |
| Timestamp | time.Time | 记录时间戳 |
