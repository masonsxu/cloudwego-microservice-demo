# 添加新 RPC 接口方法

端到端引导新增 RPC 接口，覆盖从 IDL 定义到全层实现的完整流程。

**用户输入**: $ARGUMENTS（方法名称或功能描述，例如 `CreateDepartment`）

---

## 执行流程

### Step 1: 需求确认

根据用户输入 `$ARGUMENTS`，确认以下信息：
- 方法名称（PascalCase）
- 所属业务领域（user / organization / department / role / permission / menu / 新领域）
- 请求和响应字段设计

如果信息不够明确，通过提问澄清。

### Step 2: IDL 定义

在 `idl/rpc/identity_srv/` 下找到对应的 Thrift IDL 文件，添加：
1. 请求/响应结构体（`XxxRequest` / `XxxResponse`）
2. Service 方法定义

**参考已有 IDL 文件的命名和结构风格**，保持一致。

### Step 3: 代码生成

```bash
cd rpc/identity_srv && ./script/gen_kitex_code.sh
```

生成后确认 `kitex_gen/` 下新方法已生成。**禁止手动修改 `kitex_gen/` 目录**。

### Step 4: 分层实现

按以下顺序实现，每层参考对应模板：

#### 4.1 DAL 层（数据访问）

- **接口文件**：参考 `rpc/identity_srv/biz/dal/user/user_profile_interface.go`
  - 在 `biz/dal/<domain>/` 下创建或扩展接口
- **实现文件**：参考 `rpc/identity_srv/biz/dal/user/user_profile_repository.go`
  - 使用 GORM 实现数据访问
  - 错误处理使用 `errno.WrapDatabaseError()`
  - 检查记录不存在使用 `errno.IsRecordNotFound(err)`
- **聚合注册**：确保新接口方法已体现在 `biz/dal/dal.go` 的聚合接口中

#### 4.2 Converter 层（DTO 转换）

- **接口文件**：参考 `rpc/identity_srv/biz/converter/user/user_profile.go`
- **实现文件**：参考 `rpc/identity_srv/biz/converter/user/user_profile_impl.go`
  - 纯函数，Model ↔ Thrift DTO 转换
- **聚合注册**：更新 `biz/converter/converter.go`

#### 4.3 Logic 层（业务逻辑）

- **接口文件**：参考 `rpc/identity_srv/biz/logic/user/user_profile_logic.go`
- **实现文件**：参考 `rpc/identity_srv/biz/logic/user/user_profile_logic_impl.go`
  - 编排 DAL 调用和 Converter 转换
  - 返回 `errno.ErrNo` 类型错误
- **聚合注册**：更新 `biz/logic/logic.go`

#### 4.4 Handler 层

- **文件**：`rpc/identity_srv/handler.go`
- **模式**：统一使用 `s.logic.Xxx(ctx, req)` + `errno.ToKitexError(err)` 包装
- 参考已有方法的实现风格

### Step 5: Wire 依赖注入

如果新增了 Provider：
1. 更新 `rpc/identity_srv/wire/provider.go`
2. 运行 `cd rpc/identity_srv/wire && wire`
3. 确认 `wire_gen.go` 正确生成

### Step 6: 验证

```bash
cd rpc/identity_srv && go build ./...
cd rpc/identity_srv && go vet ./...
cd rpc/identity_srv && golangci-lint run
```

### Step 7: 编写测试

按 `docs/09-testing-guide.md` 规范为新代码编写测试：
- Converter 层：纯函数测试，参考 `biz/converter/user/user_profile_test.go`
- Logic 层：mock DAL 测试，参考 `biz/logic/user/user_logic_validation_test.go`
- 测试命名：`Test<InterfaceName>_<MethodName>`

---

## 注意事项

- 先读取 `docs/03-development.md` 了解完整开发规范
- 错误码格式 `A-BB-CCC`，参考 `rpc/identity_srv/pkg/errno/` 已有定义
- 导入顺序：标准库 → 第三方库 → 项目内部包
- 最大行长度 120 字符
- 接口无 `I` 前缀，实现加 `Impl` 后缀
