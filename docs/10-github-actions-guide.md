# GitHub Actions CI/CD 使用指南

本项目配置了 GitHub Actions 自动化测试流水线，每次推送代码或创建 Pull Request 时自动运行测试。

## 📋 Workflows 概览

### 1. **Test Workflow** (`.github/workflows/test.yml`)

完整的测试流程，包括：
- ✅ RPC 服务测试（带 PostgreSQL、Redis、etcd）
- ✅ Gateway 服务测试
- ✅ 代码检查（golangci-lint）
- ✅ 构建验证

**触发条件**：
- Push 到 `main`, `master`, `develop` 分支
- 创建 Pull Request

**检查项**：
```
✓ 代码格式化检查
✓ 静态分析
✓ 单元测试
✅ 竞态条件检测
✓ 测试覆盖率
✓ 代码构建
```

### 2. **Coverage Workflow** (`.github/workflows/coverage.yml`)

生成详细的覆盖率报告：
- 📊 生成 HTML 覆盖率报告
- 📈 计算覆盖率百分比
- 💬 在 PR 中评论覆盖率
- 📦 上传覆盖率报告到 Artifacts
- 🎯 检查覆盖率阈值（30%）

**特性**：
- 合并 RPC 和 Gateway 的覆盖率
- 自动在 PR 中评论覆盖率结果
- 生成可视化 HTML 报告
- 上传到 Codecov

### 3. **CI Workflow** (`.github/workflows/ci.yml`)

快速检查 + 完整测试流程：
- 🔍 代码格式化检查
- 🔍 静态分析
- 🔍 TODO/FIXME/HACK 检查
- ✅ 完整测试（包含数据库）

**触发条件**：
- PR 或推送到 `main/master` 分支

## 🚀 使用方法

### 本地运行前自检

在推送代码前，建议先本地运行：

```bash
# 1. 格式化检查
golangci-lint run --disable-all --enable goimports,go vet

# 2. 运行测试（需要基础设施）
cd docker && podman-compose up -d
cd rpc/identity_srv && go test ./... -v
cd gateway && go test ./... -v

# 3. 生成覆盖率
./scripts/generate-coverage-report.sh
```

### 推送代码触发 CI

```bash
# 推送到主分支
git push origin main

# 或创建 PR
git push origin feature-branch
# 然后在 GitHub 上创建 PR
```

### 查看 CI 结果

1. **GitHub Actions 页面**
   - 进入仓库的 "Actions" 标签
   - 查看最近的 workflow 运行

2. **PR 检查状态**
   - PR 页面会显示所有检查项的状态
   - 必须所有检查通过才能合并

3. **覆盖率报告**
   - PR 评论会显示覆盖率百分比
   - 点击 Actions 运行查看详细报告
   - 下载 Artifacts 查看 HTML 报告

## 📊 Codecov 集成

项目集成了 Codecov 用于可视化覆盖率：

### 配置文件 (`.codecov.yml`)

```yaml
coverage:
  status:
    project:
      default:
        target: 70%    # 目标覆盖率
        threshold: 5%   # 允许下降 5%
    patch:
      default:
        target: 80%    # 新代码目标
        threshold: 10%  # 允许下降 10%
```

### 查看覆盖率

1. **Codecov Dashboard**
   ```
   https://codecov.io/github/YOUR_ORG/YOUR_REPO
   ```

2. **PR 评论**
   - 每次 PR 会自动评论覆盖率变化

3. **覆盖率徽章**
   ```markdown
   ![codecov](https://codecov.io/gh/YOUR_ORG/YOUR_REPO/branch/main/graph/badge.svg)
   ```

## 🔧 配置说明

### 环境变量

CI 使用以下环境变量（在 workflow 中配置）：

```bash
# 数据库
DB_DSN: host=localhost port=5432 user=test_user password=test_password...

# Redis
REDIS_ADDR: localhost:6379

# etcd
ETCD_ADDR: localhost:2379
```

### 服务依赖

CI 自动启动以下服务：
- PostgreSQL 15（端口 5432）
- Redis 7（端口 6379）
- etcd v3.5.9（端口 2379）

### Go 版本

- Go 1.24
- 使用缓存加速构建

## 📈 覆盖率目标

| 模块 | 当前目标 | 最终目标 |
|------|---------|----------|
| pkg 包 | 80% | 90% |
| Logic 层 | 30% | 70% |
| DAL 层 | 0% | 60% |
| Gateway 层 | 0% | 60% |
| **总体** | **30%** | **70%** |

## 🐛 常见问题

### Q: CI 测试失败，但本地测试通过？

**可能原因**：
1. 数据库版本不同（CI 使用 PostgreSQL 15）
2. 时区或环境变量不同
3. Go 版本不同

**解决方法**：
```bash
# 使用 CI 相同的数据库版本
docker run -d -p 5432:5432 \
  -e POSTGRES_DB=test_db \
  -e POSTGRES_USER=test \
  -e POSTGRES_PASSWORD=test \
  postgres:15-alpine
```

### Q: 覆盖率没有上传到 Codecov？

**检查项**：
1. Codecov token 是否配置
   - Settings → Secrets → Actions → `CODECOV_TOKEN`

2. workflow 是否成功运行
   - 查看 Actions 日志

3. 覆盖率文件是否生成
   - 检查 Artifacts

### Q: 如何跳过 CI？

**不推荐**，但如果必须：

```bash
# 在 commit message 中
git commit -m "feat: add feature [ci skip]"
git commit -m "feat: add feature [skip ci]"
```

### Q: 如何调试 CI 失败？

1. **启用调试日志**
   ```yaml
   - name: Run tests
     run: go test -v -race ./...
   ```

2. **使用 tmate 进行交互式调试**
   ```yaml
   - name: Setup tmate session
     uses: mxschmitt/action-tmate@v3
   ```

3. **下载 Artifacts**
   - CI 运行页面 → Artifacts
   - 下载日志和覆盖率文件

## 🎯 最佳实践

### 1. 保持测试快速

```go
// ✅ 好的测试 - 快速
func TestValidation(t *testing.T) {
    result := validateInput("test")
    assert.True(t, result)
}

// ❌ 差的测试 - 慢速（依赖数据库）
func TestValidationWithDB(t *testing.T) {
    db := setupDatabase()
    result := validateFromDB(db, "test")
    assert.True(t, result)
}
```

### 2. 使用 Table-Driven Tests

```go
func TestCalculate(t *testing.T) {
    tests := []struct {
        name   string
        input  int
        want   int
    }{
        {"case1", 1, 2},
        {"case2", 2, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, calculate(tt.input))
        })
    }
}
```

### 3. 提交前测试

```bash
# 使用 pre-commit hook
git pre-commit

# 或手动运行
go test ./... -race -cover
golangci-lint run
```

## 📝 代码审查检查清单

PR 合并前确保：
- [ ] 所有 CI 检查通过
- [ ] 覆盖率不低于当前分支
- [ ] 没有新增 TODO/FIXME
- [ ] 代码格式化正确
- [ ] 没有引入新的警告

## 🔗 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Codecov 文档](https://docs.codecov.com/)
- [golangci-lint 文档](https://golangci-lint.run/)
- [Go Testing 指南](docs/09-testing-guide.md)

---

**提示**: 本地运行 `./scripts/generate-coverage-report.sh` 可以预览 CI 中的覆盖率报告！
