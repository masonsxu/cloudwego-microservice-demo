# 贡献指南

感谢你对 CloudWeGo Scaffold 的关注！我们欢迎各种形式的贡献。

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发环境设置](#开发环境设置)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)
- [代码规范](#代码规范)

---

## 行为准则

参与本项目即表示你同意遵守我们的 [行为准则](CODE_OF_CONDUCT.md)。

---

## 如何贡献

### 报告 Bug

1. 搜索 [现有 Issues](https://github.com/masonsxu/cloudwego-microservice-demo/issues) 确认问题尚未被报告
2. 创建新 Issue，包含：
   - 清晰的标题和描述
   - 复现步骤
   - 预期行为和实际行为
   - 环境信息（Go 版本、操作系统等）

### 提出新功能

1. 搜索现有 Issues 确认功能尚未被提出
2. 创建新 Issue 描述：
   - 功能的目的和价值
   - 预期的使用方式
   - 可能的实现方案

### 提交代码

1. Fork 项目
2. 创建特性分支
3. 编写代码和测试
4. 提交 Pull Request

---

## 开发环境设置

### 前置要求

- Go 1.24+
- Docker 20.10+
- Docker Compose 2.0+

### 安装开发工具

```bash
# Kitex 工具链
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest

# Hertz 工具
go install github.com/cloudwego/hertz/cmd/hz@latest

# Wire 依赖注入
go install github.com/google/wire/cmd/wire@latest

# golangci-lint 代码检查
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 启动开发环境

```bash
# 启动基础设施
cd docker && ./deploy.sh dev up-base

# 安装 Git Hooks
ln -s -f ../../scripts/git-hooks/pre-commit .git/hooks/pre-commit
```

---

## 提交规范

### Commit 消息格式

```
<type>: <subject>

<body>

<footer>
```

### Type 类型

| 类型 | 说明 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档更新 |
| `style` | 代码格式（不影响功能） |
| `refactor` | 重构（不是新功能或 Bug 修复） |
| `perf` | 性能优化 |
| `test` | 测试相关 |
| `chore` | 构建/工具链更新 |

### 示例

```
feat: 添加用户头像上传功能

- 支持 jpg/png/webp 格式
- 最大文件大小 5MB
- 自动生成缩略图

Closes #123
```

---

## Pull Request 流程

### 1. Fork 和克隆

```bash
# Fork 项目到你的 GitHub
# 然后克隆
git clone https://github.com/YOUR_USERNAME/cloudwego-microservice-demo.git
cd cloudwego-microservice-demo

# 添加上游仓库
git remote add upstream https://github.com/masonsxu/cloudwego-microservice-demo.git
```

### 2. 创建分支

```bash
# 同步最新代码
git fetch upstream
git checkout main
git merge upstream/main

# 创建特性分支
git checkout -b feature/your-feature-name
```

### 3. 开发和测试

```bash
# 编写代码...

# 运行测试
go test ./... -v

# 代码检查
golangci-lint run

# 如果修改了 IDL，重新生成代码
cd rpc/identity_srv && ./script/gen_kitex_code.sh
cd gateway && ./script/gen_hertz_code.sh

# 如果修改了依赖，重新生成 Wire
cd wire && wire
```

### 4. 提交和推送

```bash
git add .
git commit -m "feat: your feature description"
git push origin feature/your-feature-name
```

### 5. 创建 Pull Request

1. 访问你的 Fork 仓库
2. 点击 "New Pull Request"
3. 填写 PR 描述：
   - 变更说明
   - 关联的 Issue
   - 测试方法

---

## 代码规范

### Go 代码规范

- 遵循 [Effective Go](https://go.dev/doc/effective_go)
- 使用 `gofmt` 格式化代码
- 通过 `golangci-lint` 检查

### 项目特定规范

1. **IDL-First**：先修改 IDL，再生成代码
2. **分层架构**：遵循 Handler → Logic → DAL → Model 分层
3. **错误处理**：使用项目定义的错误码
4. **配置管理**：使用环境变量，不使用 YAML 文件

### 代码检查

```bash
# 运行所有检查
golangci-lint run

# 自动修复
golangci-lint run --fix
```

---

## 问题反馈

如有任何问题，欢迎：

- 提交 [Issue](https://github.com/masonsxu/cloudwego-microservice-demo/issues)
- 发起 [Discussion](https://github.com/masonsxu/cloudwego-microservice-demo/discussions)

感谢你的贡献！
