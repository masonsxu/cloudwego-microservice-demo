# GitHub Actions CI/CD 使用指南

本项目配置了 GitHub Actions 自动化 CI 流水线，每次推送代码或创建 Pull Request 时自动运行检查。

## Workflows 概览

项目共 4 个工作流文件：

| 工作流 | 文件 | 用途 |
|--------|------|------|
| CI | `ci.yml` | 统一的代码检查、测试、构建、覆盖率 |
| PR Review | `pr-review.yml` | AI 代码审查 |
| Issue Killer | `issue-killer.yml` | AI 自动实现 |
| Issue Triage | `issue-triage.yaml` | Issue 自动分类 |

## CI 工作流详解 (`ci.yml`)

### 触发条件

- Push 到 `main` 分支
- 创建/更新 PR 到 `main` 分支
- 自动排除：`**.md`、`docs/**`、`docker/**`、`LICENSE`、AI 工作流文件

### 并发控制

同一分支的多次 push 只保留最新一次运行，旧的运行自动取消：

```yaml
concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true
```

### Job 结构与依赖关系

```
lint ──────────────> build
test-rpc ──────┐
               ├──> coverage (仅 PR)
test-gateway ──┘
```

- `lint`、`test-rpc`、`test-gateway` 三个 job 并行运行，互不阻塞
- `build` 依赖 `lint` 通过后运行
- `coverage` 依赖两个 test job 完成后运行，且仅在 PR 时触发

### Job 1: Lint

无需服务容器，执行以下检查：

1. 生成 Swagger 文档（gateway 的 `docs/` 在 `.gitignore` 中，lint 前必须生成）
2. `golangci-lint run`：对 RPC 和 Gateway 分别执行静态分析
3. `golangci-lint fmt --diff`：检查代码格式一致性

使用 `golangci/golangci-lint-action@v7` 官方 action，锁定 `v2.1` 版本，自带 lint 结果缓存。

**格式化工具统一说明**：`.golangci.yml` 已配置 `gci` + `goimports` + `gofumpt` + `golines` + `swaggo`，通过 `golangci-lint fmt` 统一执行，无需单独安装 `goimports-reviser` 等工具。

### Job 2: Test RPC

需要服务容器（均配有 healthcheck）：

| 服务 | 镜像 | 端口 | 健康检查 |
|------|------|------|----------|
| PostgreSQL | `postgres:15-alpine` | 5432 | `pg_isready` |
| Redis | `redis:7-alpine` | 6379 | `redis-cli ping` |
| etcd | `quay.io/coreos/etcd:v3.5.9` | 2379 | `etcdctl endpoint health` |

执行 `go test -v -race -coverprofile -covermode=atomic`，并上传覆盖率产物。

### Job 3: Test Gateway

无需服务容器，先生成 Swagger 文档再运行测试，上传覆盖率产物。

### Job 4: Build

依赖 `lint` 通过，验证两个模块可以编译通过（输出到 `/dev/null`，仅做编译检查）。

### Job 5: Coverage（仅 PR）

依赖两个 test job 完成：

1. 下载两个模块的覆盖率产物
2. 使用 `gocovmerge` 合并覆盖率
3. 提取 RPC / Gateway / 总体覆盖率百分比
4. 检查阈值（30%），低于阈值则失败
5. 在 PR 上评论覆盖率表格（查找并更新已有评论，避免重复）
6. 上传合并后的覆盖率报告（保留 30 天）

PR 评论示例：

```
## 📊 测试覆盖率报告

| 模块 | 覆盖率 | 状态 |
|------|--------|------|
| RPC Service | xx% | 🟢/🟡/🟠/🔴 |
| Gateway Service | xx% | 🟢/🟡/🟠/🔴 |
| **总体** | **xx%** | 🟢/🟡/🟠/🔴 |

> 阈值: 30% | ✅ 达标 / ❌ 未达标
```

颜色规则：🟢 >= 80% / 🟡 >= 60% / 🟠 >= 30% / 🔴 < 30%

## 本地运行前自检

在推送代码前，建议先本地运行：

```bash
# 1. 静态分析
cd rpc/identity_srv && golangci-lint run --timeout=5m ./...
cd gateway && golangci-lint run --timeout=5m ./...

# 2. 格式化检查（查看差异）
cd rpc/identity_srv && golangci-lint fmt --diff
cd gateway && golangci-lint fmt --diff

# 3. 自动修复格式
cd rpc/identity_srv && golangci-lint fmt
cd gateway && golangci-lint fmt

# 4. 运行测试（需要基础设施）
cd docker && podman-compose up -d
cd rpc/identity_srv && go test ./... -v -race
cd gateway && go test ./... -v -race

# 5. 编译检查
cd rpc/identity_srv && go build -o /dev/null ./cmd/main.go
cd gateway && go build -o /dev/null ./cmd/main.go
```

## 配置说明

### 环境变量

CI 中测试使用以下环境变量：

```bash
# 数据库（仅 RPC 测试需要）
DB_DSN: host=localhost port=5432 user=test_user password=test_password dbname=identity_db_test sslmode=disable

# Redis（仅 RPC 测试需要）
REDIS_ADDR: localhost:6379

# etcd（仅 RPC 测试需要）
ETCD_ADDR: localhost:2379
```

### Go 版本

- Go 1.24
- 使用 `setup-go` 内置缓存，按 `go.sum` 文件做缓存 key

### 权限

CI 使用最小权限原则：

```yaml
permissions:
  contents: read        # 读取代码
  pull-requests: write  # 评论覆盖率
```

## 覆盖率目标

| 模块 | 当前目标 | 最终目标 |
|------|---------|----------|
| pkg 包 | 80% | 90% |
| Logic 层 | 30% | 70% |
| DAL 层 | 0% | 60% |
| Gateway 层 | 0% | 60% |
| **总体** | **30%** | **70%** |

## 常见问题

### Q: CI 测试失败，但本地测试通过？

**可能原因**：
1. 数据库版本不同（CI 使用 PostgreSQL 15）
2. 时区或环境变量不同
3. Go 版本不同
4. Swagger 文档未生成（Gateway 测试前需要 `swag init`）

**解决方法**：
```bash
# 使用 CI 相同的数据库版本
docker run -d -p 5432:5432 \
  -e POSTGRES_DB=identity_db_test \
  -e POSTGRES_USER=test_user \
  -e POSTGRES_PASSWORD=test_password \
  postgres:15-alpine
```

### Q: 修改了文档/Docker 配置后 CI 没有触发？

这是预期行为。`ci.yml` 配置了 `paths-ignore`，修改以下文件不触发 CI：
- `**.md`（所有 Markdown 文件）
- `docs/**`（文档目录）
- `docker/**`（Docker 配置）
- `LICENSE`

### Q: 同一 PR 多次 push，之前的运行怎么办？

CI 配置了并发控制（`cancel-in-progress: true`），同一分支的旧运行会被自动取消，只保留最新的运行。

### Q: 如何跳过 CI？

**不推荐**，但如果必须：

```bash
git commit -m "feat: add feature [ci skip]"
git commit -m "feat: add feature [skip ci]"
```

### Q: 如何调试 CI 失败？

1. **查看 Actions 日志**：进入仓库 Actions 标签，点击对应的 workflow run
2. **下载覆盖率报告**：在 run 页面的 Artifacts 区域下载
3. **使用 tmate 交互式调试**（临时添加到 workflow）：
   ```yaml
   - name: Setup tmate session
     uses: mxschmitt/action-tmate@v3
   ```

### Q: PR 上有多条覆盖率评论？

不会出现。coverage job 使用"查找并更新"策略，同一 PR 只维护一条覆盖率评论。

## 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [golangci-lint 文档](https://golangci-lint.run/)
- [测试指南](../02-开发规范/测试指南.md)
