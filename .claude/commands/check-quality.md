# 代码质量检查

依次执行全套代码质量检查，输出综合报告和改进建议。

**用户输入**: $ARGUMENTS（可选：指定检查范围，例如 `rpc/identity_srv` 或 `gateway`，为空则检查全部）

---

## 执行流程

### Step 1: 确定检查范围

- 如果 `$ARGUMENTS` 指定了路径，仅检查该路径
- 如果为空，检查两个模块：`rpc/identity_srv` 和 `gateway`

### Step 2: 编译检查

```bash
cd rpc/identity_srv && go build ./...
cd gateway && go build ./...
```

记录编译是否通过，如有错误记录具体信息。

### Step 3: Lint 检查

```bash
cd rpc/identity_srv && golangci-lint run
cd gateway && golangci-lint run
```

记录各类告警数量和类型：
- 导入顺序问题
- 行长度超限
- 未使用的变量/导入
- 错误处理问题
- 其他 lint 规则违反

### Step 4: 测试覆盖率

```bash
cd rpc/identity_srv && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out | grep total
cd gateway && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out | grep total
```

先读取 `docs/09-testing-guide.md` 中的覆盖率目标，对比实际覆盖率。

### Step 5: 架构合规检查

检查以下合规性规则：

1. **分层约束**：
   - Handler 层不直接调用 DAL（应通过 Logic）
   - Logic 层不直接导入 `kitex_gen`（应通过 Converter）
   - DAL 层不导入 Logic 或 Handler 包

2. **禁止修改**：
   - `kitex_gen/` 目录下无手动修改（对比 git 状态）

3. **命名规范**：
   - 接口无 `I` 前缀
   - 实现加 `Impl` 后缀

4. **错误处理**：
   - DAL 层使用 `errno.WrapDatabaseError()`
   - 不直接比较 `gorm.ErrRecordNotFound`

5. **配置管理**：
   - 无 YAML 配置文件引入
   - 配置通过环境变量管理

### Step 6: 输出报告

汇总报告格式：

```
## 代码质量检查报告

### 编译状态
- rpc/identity_srv: ✅ / ❌
- gateway: ✅ / ❌

### Lint 检查
- rpc/identity_srv: X 个告警
- gateway: X 个告警
- [列出主要告警]

### 测试覆盖率
- rpc/identity_srv: XX.X%（目标: XX%）
- gateway: XX.X%（目标: XX%）

### 架构合规
- [列出违规项]

### 改进建议
1. [具体可操作的建议]
2. ...
```

---

## 注意事项

- 检查过程中如遇到编译错误，仍继续执行后续检查步骤
- 覆盖率目标以 `docs/09-testing-guide.md` 中的定义为准
- 架构合规检查通过代码搜索完成，不需要运行额外工具
