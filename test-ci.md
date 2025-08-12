# GitHub Actions 测试指南

## 当前工作流配置

### 1. 发布工作流 (`.github/workflows/build-and-release.yml`)
**触发条件:**
- 创建 GitHub Release 时

**执行内容:**
- 多平台二进制文件构建
- 自动上传二进制文件到 Release

### 2. 简单测试工作流 (`.github/workflows/simple-test.yml`)
**触发条件:**
- 推送到 `main` 或 `master` 分支
- 手动触发

**执行内容:**
- 基础构建检查
- 帮助命令测试

### 3. 测试工作流 (`.github/workflows/test.yml`)
**触发条件:**
- 手动触发
- Pull Request

**执行内容:**
- 代码格式检查
- Go vet 检查
- 可选的单元测试

## 测试方法

### 1. 测试简单构建
```bash
# 推送到主分支
git add .
git commit -m "test: trigger simple test workflow"
git push origin master
```

### 2. 测试发布流程
在 GitHub 仓库页面：
1. 点击 "Releases" 标签
2. 点击 "Create a new release"
3. 填写标签版本（如 v1.0.0）
4. 填写发布说明
5. 点击 "Publish release"

### 3. 手动触发测试
在 GitHub 仓库页面：
1. 点击 "Actions" 标签
2. 选择 "Simple Test" 或 "Test" 工作流
3. 点击 "Run workflow"

## 预期结果

### 简单测试工作流成功后：
- ✅ Go 环境设置成功
- ✅ 后端构建成功
- ✅ 帮助命令测试通过

### 发布工作流成功后：
- ✅ 多平台二进制文件构建
- ✅ 二进制文件自动上传到 Release
- ✅ 用户可以直接下载使用

## 故障排除

### 1. 工作流被跳过
- 检查分支名称是否为 `main` 或 `master`
- 检查是否有权限推送到仓库

### 2. 二进制构建失败
- 检查 Go 版本兼容性
- 检查 CGO 依赖安装

## 当前状态

- ✅ 移除了测试依赖
- ✅ 移除了 Docker 自动构建
- ✅ 支持多平台二进制构建
- ✅ 自动化发布流程
- ✅ 简化的触发条件
- ✅ 只在创建标签时构建和发布