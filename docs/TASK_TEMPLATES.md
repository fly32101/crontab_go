# 任务模板功能

任务模板功能允许用户创建、管理和使用预定义的任务模板，提高任务创建的效率和标准化程度。

## 功能特性

- ✅ 创建和管理任务模板
- ✅ 模板分类管理
- ✅ 公共模板和私有模板
- ✅ 模板搜索和筛选
- ✅ 热门模板推荐
- ✅ 从模板快速创建任务
- ✅ 模板使用统计
- ✅ 标签系统
- ✅ 默认模板库

## 模板分类

系统预置了以下模板分类：

### 1. 通用任务 (general)
- 适用于各种通用场景的任务模板
- 图标：AppstoreOutlined
- 颜色：蓝色 (#1890ff)

### 2. 备份任务 (backup)
- 数据库备份、文件备份等相关模板
- 图标：DatabaseOutlined
- 颜色：绿色 (#52c41a)

### 3. 监控任务 (monitoring)
- 系统监控、健康检查等相关模板
- 图标：MonitorOutlined
- 颜色：橙色 (#faad14)

### 4. 清理任务 (cleanup)
- 日志清理、临时文件清理等相关模板
- 图标：DeleteOutlined
- 颜色：红色 (#f5222d)

### 5. 通知任务 (notification)
- 报告发送、通知推送等相关模板
- 图标：BellOutlined
- 颜色：紫色 (#722ed1)

### 6. API调用 (api)
- HTTP API 调用相关模板
- 图标：ApiOutlined
- 颜色：青色 (#13c2c2)

## 默认模板

系统预置了以下默认模板：

### 1. 数据库备份
- **分类**: backup
- **描述**: 定期备份数据库
- **Cron表达式**: `0 0 2 * * *` (每天凌晨2点)
- **命令**: `mysqldump -u root -p database_name > /backup/db_$(date +%Y%m%d).sql`
- **标签**: ["backup", "database", "mysql"]

### 2. 日志清理
- **分类**: cleanup
- **描述**: 清理7天前的日志文件
- **Cron表达式**: `0 0 3 * * *` (每天凌晨3点)
- **命令**: `find /var/log -name '*.log' -mtime +7 -delete`
- **标签**: ["cleanup", "logs"]

### 3. 系统监控
- **分类**: monitoring
- **描述**: 检查系统资源使用情况
- **Cron表达式**: `0 */10 * * * *` (每10分钟)
- **命令**: `df -h && free -m && uptime`
- **标签**: ["monitoring", "system"]

### 4. 健康检查
- **分类**: monitoring
- **描述**: 检查服务健康状态
- **Cron表达式**: `0 */5 * * * *` (每5分钟)
- **命令**: `https://api.example.com/health`
- **方法**: GET
- **标签**: ["monitoring", "health", "api"]

### 5. 每日报告
- **分类**: notification
- **描述**: 发送每日系统报告
- **Cron表达式**: `0 0 9 * * *` (每天上午9点)
- **命令**: `python /scripts/daily_report.py`
- **标签**: ["notification", "report"]
- **通知**: 成功和失败时都通知

## 使用指南

### 1. 浏览模板

1. 访问"任务模板"页面
2. 查看热门模板推荐
3. 使用搜索和筛选功能找到合适的模板
4. 按分类、标签、创建者等条件筛选

### 2. 创建模板

1. 点击"新建模板"按钮
2. 填写模板信息：
   - 模板名称（必填）
   - 模板描述
   - 分类选择
   - Cron表达式（必填）
   - 命令或URL（必填）
   - HTTP方法（如果是URL）
   - 请求头（如果是HTTP请求）
   - 标签
   - 是否为公共模板
   - 通知设置
3. 保存模板

### 3. 使用模板创建任务

1. 在模板列表中找到要使用的模板
2. 点击"使用"按钮
3. 填写任务名称
4. 选择是否立即启用
5. 确认创建任务

### 4. 管理模板

- **编辑模板**: 点击编辑按钮修改模板信息
- **删除模板**: 点击删除按钮移除模板
- **查看统计**: 查看模板使用次数和统计信息

## API 接口

### 模板管理

```http
# 创建模板
POST /api/v1/templates
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "模板名称",
  "description": "模板描述",
  "category": "general",
  "schedule": "0 */5 * * * *",
  "command": "echo 'Hello World'",
  "tags": "[\"tag1\", \"tag2\"]",
  "is_public": true
}

# 获取模板列表
GET /api/v1/templates
Authorization: Bearer <token>

# 搜索模板
GET /api/v1/templates/search?keyword=备份&category=backup&page=1&page_size=10
Authorization: Bearer <token>

# 获取热门模板
GET /api/v1/templates/popular?limit=10
Authorization: Bearer <token>

# 从模板创建任务
POST /api/v1/templates/create-task
Content-Type: application/json
Authorization: Bearer <token>

{
  "template_id": 1,
  "task_name": "我的备份任务",
  "enabled": true,
  "overrides": {
    "schedule": "0 0 1 * * *"
  }
}
```

### 分类管理

```http
# 获取分类列表
GET /api/v1/template-categories
Authorization: Bearer <token>

# 创建分类
POST /api/v1/template-categories
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "custom",
  "description": "自定义分类",
  "icon": "CustomOutlined",
  "color": "#1890ff",
  "sort_order": 10
}
```

## 数据模型

### TaskTemplate 任务模板

```go
type TaskTemplate struct {
    ID                 int       `json:"id"`
    Name               string    `json:"name"`                    // 模板名称
    Description        string    `json:"description"`             // 模板描述
    Category           string    `json:"category"`                // 模板分类
    Schedule           string    `json:"schedule"`                // Cron表达式
    Command            string    `json:"command"`                 // 命令或URL
    Method             string    `json:"method"`                  // HTTP请求方法
    Headers            string    `json:"headers"`                 // HTTP请求头
    NotifyOnSuccess    bool      `json:"notify_on_success"`       // 成功时是否通知
    NotifyOnFailure    bool      `json:"notify_on_failure"`       // 失败时是否通知
    NotificationTypes  string    `json:"notification_types"`     // 通知类型
    NotificationConfig string    `json:"notification_config"`    // 通知配置
    Tags               string    `json:"tags"`                    // 标签
    IsPublic           bool      `json:"is_public"`               // 是否为公共模板
    CreatedBy          int       `json:"created_by"`              // 创建者ID
    UsageCount         int       `json:"usage_count"`             // 使用次数
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}
```

### TaskTemplateCategory 模板分类

```go
type TaskTemplateCategory struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`        // 分类名称
    Description string `json:"description"` // 分类描述
    Icon        string `json:"icon"`        // 图标
    Color       string `json:"color"`       // 颜色
    SortOrder   int    `json:"sort_order"`  // 排序
}
```

## 最佳实践

### 1. 模板设计原则

- **通用性**: 设计通用的模板，适用于多种场景
- **参数化**: 使用环境变量或参数，提高模板的灵活性
- **文档化**: 提供清晰的描述和使用说明
- **标签化**: 使用合适的标签，便于搜索和分类

### 2. 命名规范

- **模板名称**: 使用描述性的名称，如"MySQL数据库备份"
- **标签**: 使用小写英文，如["backup", "mysql", "database"]
- **分类**: 使用预定义的分类名称

### 3. 安全考虑

- **敏感信息**: 不要在模板中包含密码、密钥等敏感信息
- **权限控制**: 合理设置公共模板和私有模板
- **命令验证**: 确保模板中的命令是安全的

### 4. 维护建议

- **定期更新**: 定期更新模板内容，保持最新状态
- **使用统计**: 关注模板使用统计，优化热门模板
- **用户反馈**: 收集用户反馈，改进模板质量

## 故障排除

### 常见问题

1. **模板创建失败**
   - 检查必填字段是否完整
   - 验证Cron表达式格式
   - 确认分类是否存在

2. **从模板创建任务失败**
   - 检查模板权限（私有模板只能创建者使用）
   - 验证任务名称是否重复
   - 确认用户有创建任务的权限

3. **模板搜索无结果**
   - 检查搜索关键词
   - 确认筛选条件
   - 验证模板是否为公共模板

### 日志查看

- 模板相关的操作日志会记录在系统日志中
- 可以通过API响应查看详细错误信息
- 数据库操作日志可以帮助诊断问题

## 扩展功能

### 未来计划

- [ ] 模板版本管理
- [ ] 模板导入/导出
- [ ] 模板评分和评论
- [ ] 模板使用分析
- [ ] 模板推荐算法
- [ ] 模板市场/商店

### 自定义扩展

- 可以通过API创建自定义分类
- 支持自定义标签系统
- 可以扩展模板字段
- 支持自定义模板验证规则

---

任务模板功能大大提高了任务创建的效率，通过预定义的模板，用户可以快速创建标准化的任务，减少重复工作，提高系统的可维护性。