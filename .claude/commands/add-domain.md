# 创建新业务领域模块

为 RPC 服务创建完整的业务领域目录结构和骨架代码，并集成到三个聚合接口。

**用户输入**: $ARGUMENTS（领域名称，例如 `notification`、`audit`）

---

## 执行流程

### Step 1: 需求确认

根据用户输入 `$ARGUMENTS`，确认：
- 领域名称（小写，用于包名和目录名）
- 核心实体名称（PascalCase，用于接口和结构体命名）
- 实体主要字段
- 是否需要关联已有领域（如 user、organization）

### Step 2: 数据模型

创建 `rpc/identity_srv/models/<domain>.go`：
- 参考已有 Model 文件的 GORM 标签风格
- 包含基础字段（ID、创建时间、更新时间等）
- 使用 `gorm:"primaryKey"` 等标准标签

### Step 3: DAL 层

创建 `rpc/identity_srv/biz/dal/<domain>/` 目录：

1. **接口文件** `<entity>_interface.go`
   - 参考 `biz/dal/user/user_profile_interface.go`
   - 定义 Repository 接口（CRUD + 业务查询方法）
   - 接口命名：`<Entity>Repository`（无 `I` 前缀）

2. **实现文件** `<entity>_repository.go`
   - 参考 `biz/dal/user/user_profile_repository.go`
   - 实现命名：`<Entity>RepositoryImpl`
   - 使用 `*gorm.DB` 作为依赖
   - 错误处理使用 `errno.WrapDatabaseError()`

3. **聚合注册**：在 `biz/dal/dal.go` 中
   - 聚合接口添加 `<Entity>() <Entity>Repository` 方法
   - 在 `biz/dal/dal_impl.go` 中添加对应实现

### Step 4: Converter 层

创建 `rpc/identity_srv/biz/converter/<domain>/` 目录：

1. **接口文件** `<entity>.go`
   - 参考 `biz/converter/user/user_profile.go`
   - 定义转换接口（Model → DTO、DTO → Model）
   - 接口命名：`<Entity>Converter`

2. **实现文件** `<entity>_impl.go`
   - 参考 `biz/converter/user/user_profile_impl.go`
   - 纯函数实现，无副作用
   - 实现命名：`<Entity>ConverterImpl`

3. **聚合注册**：更新 `biz/converter/converter.go`

### Step 5: Logic 层

创建 `rpc/identity_srv/biz/logic/<domain>/` 目录：

1. **接口文件** `<entity>_logic.go`
   - 参考 `biz/logic/user/user_profile_logic.go`
   - 定义业务逻辑接口
   - 接口命名：`<Entity>Logic`

2. **实现文件** `<entity>_logic_impl.go`
   - 参考 `biz/logic/user/user_profile_logic_impl.go`
   - 依赖注入 DAL 和 Converter
   - 实现命名：`<Entity>LogicImpl`

3. **聚合注册**：更新 `biz/logic/logic.go`，嵌入新领域接口

### Step 6: Wire 依赖注入

1. 在 `rpc/identity_srv/wire/provider.go` 添加新领域的 Provider
2. 遵循规范：
   - 每个 Provider 函数只创建一个依赖
   - 返回类型必须是接口
3. 运行 `cd rpc/identity_srv/wire && wire`

### Step 7: 验证

```bash
cd rpc/identity_srv && go build ./...
cd rpc/identity_srv && go vet ./...
```

确认编译通过，三个聚合接口（`dal.go`、`converter.go`、`logic.go`）均已正确集成。

---

## 输出清单

创建完成后，列出所有新建和修改的文件：
- 新建文件（6-8 个）
- 修改文件（3-4 个聚合接口 + wire）

## 注意事项

- 先读取 `docs/02-architecture.md` 和 `docs/03-development.md` 了解架构和开发规范
- 命名规范：接口无 `I` 前缀，实现加 `Impl` 后缀
- 导入顺序：标准库 → 第三方库 → 项目内部包
- 错误码：按 `A-BB-CCC` 格式在 `pkg/errno/` 中定义新领域错误码
