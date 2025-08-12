# GitHub Actions 测试指南

## 当前工作流配置

### 1. CI 工作流 (`.github/workflows/ci.yml`)
**触发条件:**
- 推送到 `main` 或 `master` 分支
- 创建 Pull Request

**执行内容:**
- 构建检查（前端 + 后端）
- Docker 镜像构建（仅推送到主分支时）

### 2. 构建和发布工作流 (`.github/workflows/build-and-release.yml`)
**触发条件:**
- 推送到 `main` 或 `master` 分支（仅 Docker 构建）
- 创建标签 `v*`（完整构建和发布）

**执行内容:**
- Docker 镜像构建和推送
- 多平台二进制文件构建（仅标签时）
- GitHub Release 创建（仅标签时）

### 3. 测试工作流 (`.github/workflows/test.yml`)
**触发条件:**
- 手动触发
- Pull Request

**执行内容:**
- 代码格式检查
- Go vet 检查
- 可选的单元测试

## 测试方法

### 1. 测试 CI 工作流
```bash
# 推送到主分支
git add .
git commit -m "test: trigger CI workflow"
git push origin main
```

### 2. 测试完整发布流程
```bash
# 创建标签
git tag v1.0.0
git push origin v1.0.0
```

### 3. 手动触发测试
在 GitHub 仓库页面：
1. 点击 "Actions" 标签
2. 选择 "Test" 工作流
3. 点击 "Run workflow"

## 预期结果

### CI 工作流成功后：
- ✅ 前端构建成功
- ✅ 后端构建成功
- ✅ Docker 镜像推送到 GHCR

### 发布工作流成功后：
- ✅ 多平台二进制文件构建
- ✅ Docker 镜像推送
- ✅ GitHub Release 创建
- ✅ 发布说明自动生成

## 故障排除

### 1. 工作流被跳过
- 检查分支名称是否为 `main` 或 `master`
- 检查是否有权限推送到仓库

### 2. Docker 推送失败
- 检查 GITHUB_TOKEN 权限
- 确认仓库设置允许 GitHub Packages

### 3. 二进制构建失败
- 检查 Go 版本兼容性
- 检查 CGO 依赖安装

## 当前状态

- ✅ 移除了测试依赖
- ✅ 支持多平台构建
- ✅ Docker 多架构支持
- ✅ 自动化发布流程
- ✅ 灵活的触发条件