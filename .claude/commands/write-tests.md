# 按项目规范编写测试

为指定代码路径按照项目测试规范编写测试用例。

**用户输入**: $ARGUMENTS（目标代码路径或文件，例如 `rpc/identity_srv/biz/converter/department/`）

---

## 执行流程

### Step 1: 读取测试规范

先读取 `docs/09-testing-guide.md`，了解：
- 分层测试策略
- 覆盖率目标
- 测试命名规范
- 测试结构要求

### Step 2: 分析目标代码

读取 `$ARGUMENTS` 指定的代码，确定：
- 代码所在层级（Converter / DAL / Logic / Handler / Middleware）
- 需要测试的函数/方法列表
- 依赖关系（需要 mock 的接口）

### Step 3: 选择测试策略

根据代码层级选择对应策略和参考模板：

#### Converter 层 → 纯函数测试
- **参考模板**：`rpc/identity_srv/biz/converter/user/user_profile_test.go`
- 策略：直接输入/输出验证，无需 mock
- 覆盖：正常转换、空值处理、边界值

#### Logic 层 → Mock DAL 测试
- **参考模板**：`rpc/identity_srv/biz/logic/user/user_logic_validation_test.go`
- 策略：mock DAL 接口，验证业务逻辑
- 覆盖：正常流程、验证失败、DAL 错误传播

#### Middleware 层 → HTTP 测试
- **参考模板**：`gateway/internal/application/middleware/casbin_middleware/middleware_test.go`
- 策略：构造 HTTP 请求上下文，验证中间件行为
- 覆盖：放行、拦截、异常处理

#### DAL 层 → 集成测试（可选）
- 策略：使用测试数据库或 mock GORM
- 覆盖：CRUD 操作、约束违反、记录不存在

### Step 4: 编写测试

遵循以下规范：
- **文件命名**：`<source>_test.go`，与被测文件同目录
- **函数命名**：`Test<InterfaceName>_<MethodName>` 或 `Test<InterfaceName>_<MethodName>_<Scenario>`
- **测试结构**：使用 table-driven tests（`tests := []struct{...}`）
- **子测试**：使用 `t.Run(tt.name, func(t *testing.T) {...})`
- **断言**：使用 `testify/assert` 或 `testify/require`

### Step 5: 运行验证

```bash
# 运行新编写的测试
go test <target_package> -v -count=1

# 检查覆盖率
go test <target_package> -coverprofile=coverage.out
go tool cover -func=coverage.out
```

确认测试全部通过，覆盖率达到 `docs/09-testing-guide.md` 中的目标。

### Step 6: 代码检查

```bash
golangci-lint run <target_path>
```

确保测试代码也通过 lint 检查。

---

## 注意事项

- 导入顺序：标准库 → 第三方库（testify 等）→ 项目内部包
- 最大行长度 120 字符
- 不要为 `kitex_gen/` 下的生成代码编写测试
- Mock 接口应放在测试文件内或 `_test.go` 同包中
