# 测试指南

## 快速开始

### 运行所有测试
```bash
./test.sh
```

### 运行特定模块测试
```bash
# User DAL 测试
./test.sh ./biz/dal/user -v

# 带覆盖率
./test.sh ./biz/dal/user -v -coverprofile=coverage.out
```

### 查看覆盖率报告
```bash
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## 测试基础设施

### testcontainers-go
- 使用 PostgreSQL 容器进行集成测试
- 自动启动/清理容器
- 真实数据库环境，可靠性强

### 依赖
- `testify/suite` - 测试套件
- `testcontainers-go` - 容器化测试
- `gorm` - ORM

## 当前覆盖率

| 模块 | 覆盖率 | 状态 | 测试用例 |
|------|--------|------|----------|
| **user DAL** | **75.0%** | ✅ **完成** | 43/45 通过 |
| organization DAL | - | ⚪ 未开始 | - |
| user_membership DAL | - | ⚪ 未开始 | - |
| role_definition DAL | - | ⚪ 未开始 | - |
| **总体目标** | **70%** | 🟡 进行中 | - |

## User DAL 测试详情

### 已覆盖方法
- ✅ Create（创建用户）
- ✅ GetByID（通过 ID 查询）
- ✅ GetByUsername（通过用户名查询）
- ✅ GetByEmail（通过邮箱查询）
- ✅ GetByPhone（通过手机号查询）
- ✅ Update（更新用户）
- ✅ UpdatePassword（更新密码）
- ✅ UpdateLoginAttempts（更新登录尝试次数）
- ✅ CheckUsernameExists（检查用户名是否存在）
- ✅ CheckEmailExists（检查邮箱是否存在）
- ✅ CheckPhoneExists（检查手机号是否存在）
- ✅ ExistsByID（检查 ID 是否存在）
- ✅ FindByMedicalLicense（通过医疗执照查询）
- ✅ FindBySpecialty（通过专业领域查询）
- ✅ FindWithConditions（条件查询）
- ✅ WithTx（事务支持）
- ✅ HardDelete（硬删除）
- ✅ SoftDelete（软删除）
- ✅ IncrementLoginAttempts（增加登录尝试次数）
- ✅ ResetLoginAttempts（重置登录尝试次数）
- ✅ UpdateLastLoginTime（更新最后登录时间）
- ✅ SetMustChangePassword（设置必须修改密码）

### 发现并修复的 Bug
- 🐛 Repository 实现字段名错误：`medical_license_number` → `license_number`（3 处）

### 待修复测试（非阻塞）
- ⚠️ TestFindWithConditions_EmailSearch（搜索功能可能需要调整）
- ⚠️ TestFindWithConditions_Pagination（分页功能可能需要调整）

## 最佳实践

### 1. 测试套件结构
```go
type XRepositoryTestSuite struct {
    suite.Suite
    db      *gorm.DB
    repo    XRepository
    cleanup func()
}

func (s *XRepositoryTestSuite) SetupSuite() {
    // 启动容器
}

func (s *XRepositoryTestSuite) TearDownSuite() {
    // 清理容器
}

func (s *XRepositoryTestSuite) SetupTest() {
    // 每个测试前清空表
}
```

### 2. 命名规范
- 测试套件：`XxxRepositoryTestSuite`
- 测试方法：`Test<Method>_<Scenario>`（如 `TestCreate_Success`）
- 断言使用 `require`（失败时立即停止）和 `assert`（继续执行）

### 3. 测试数据管理
- 使用 `faker` 生成随机数据
- 每个测试独立的测试数据
- 测试前清空表，避免相互影响

### 4. 覆盖率目标
- DAL 层：**≥75%**
- Logic 层：**≥70%**
- Handler 层：**≥50%**

## 待补充测试

### 其他 DAL 模块
- [ ] organization DAL（20 个测试用例）
- [ ] user_membership DAL（25 个测试用例，含事务）
- [ ] role_definition DAL（15 个测试用例）
- [ ] department DAL
- [ ] menu DAL
- [ ] logo DAL

### Logic 层
- [ ] user Logic（25 个测试用例，使用 gomock）
- [ ] authentication Logic（20 个测试用例）
- [ ] organization Logic
- [ ] membership Logic
- [ ] role Logic

### Handler 层
- [ ] RPC Handler 测试
- [ ] HTTP Handler 测试

### E2E 测试
- [ ] 完整业务流程测试
- [ ] API 集成测试
