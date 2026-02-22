# 测试指南

本文档提供项目的测试策略、测试编写规范和测试覆盖率指南。

## 目录

- [测试概览](#测试概览)
- [测试结构](#测试结构)
- [运行测试](#运行测试)
- [编写测试](#编写测试)
- [测试覆盖率](#测试覆盖率)
- [测试最佳实践](#测试最佳实践)

## 测试概览

项目采用分层测试策略，确保每个层级都有适当的测试覆盖：

| 层级 | 测试类型 | 覆盖率目标 | 状态 |
|------|---------|-----------|------|
| pkg 包 | 单元测试 | 80%+ | ✅ 已完成 |
| Middleware | 单元测试 + 集成测试 | 85%+ | ✅ 已完成 |
| DAL 层 | 单元测试（使用 mock） | 70%+ | ⚠️ 待补充 |
| Logic 层 | 单元测试（使用 mock） | 75%+ | ⚠️ 待补充 |
| Handler 层 | 单元测试 + 集成测试 | 60%+ | ❌ 未开始 |
| Gateway 层 | 单元测试 + 集成测试 | 60%+ | ❌ 未开始 |

### 当前测试覆盖率

```bash
# RPC 服务
pkg/errno              100.0%  ✅
pkg/log                 84.8%  ✅
pkg/password            83.3%  ✅
internal/middleware     89.1%  ✅
biz/converter           60.0%  ⚠️
biz/dal                  0.0%  ❌
biz/logic                0.0%  ❌

# Gateway 服务
internal/infrastructure/redis  5.1%  ⚠️
其他包                          0.0%  ❌
```

## 运行测试

### 运行所有测试

```bash
# 从项目根目录测试所有模块
go test ./... -v

# 测试 RPC 服务
cd rpc/identity_srv && go test ./... -v

# 测试网关服务
cd gateway && go test ./... -v
```

### 运行特定测试

```bash
# 运行单个包的测试
go test ./pkg/errno -v

# 运行单个测试函数
go test ./pkg/password -run TestHashPassword -v

# 运行匹配模式的测试
go test ./pkg/... -run "Test.*Error" -v
```

### 生成覆盖率报告

```bash
# 生成覆盖率报告
go test ./... -coverprofile=coverage.out

# 查看总体覆盖率
go tool cover -func=coverage.out | grep total

# 生成 HTML 覆盖率报告
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率低于 80% 的函数
go tool cover -func=coverage.out | awk '$3 != "100.0%" && $3 < "80.0"'
```

### 基准测试

```bash
# 运行所有基准测试
go test ./... -bench=. -benchmem

# 运行特定基准测试
go test ./pkg/password -bench BenchmarkHashPassword -benchmem
```

## 编写测试

### 测试文件组织

```
# 测试文件命名
package_test.go        # 包级测试文件
function_test.go       # 函数级测试文件

# 示例
rpc/identity_srv/pkg/errno/
  error.go             # 实现文件
  error_test.go        # 测试文件

rpc/identity_srv/biz/logic/user/
  user_logic.go        # 实现文件
  user_logic_test.go   # 测试文件
```

### 测试函数命名规范

```go
// 基本命名
func Test<FunctionName>(t *testing.T)

// 表格驱动测试
func Test<FunctionName>(t *testing.T) {
    tests := []struct {
        name string
        // 测试字段
    }{
        {"test case 1", ...},
        {"test case 2", ...},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}

// 示例
func TestHashPassword(t *testing.T) {
    t.Run("successfully hashes password", func(t *testing.T) {
        password := "MySecurePassword123!"
        hash, err := HashPassword(password)

        require.NoError(t, err)
        assert.NotEmpty(t, hash)
    })
}
```

### 各层测试示例

#### 1. pkg 包测试（无依赖）

```go
package password

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
    t.Run("successfully hashes password", func(t *testing.T) {
        password := "MySecurePassword123!"
        hash, err := HashPassword(password)

        require.NoError(t, err)
        assert.NotEmpty(t, hash)
        assert.NotEqual(t, password, hash)
        assert.Contains(t, hash, "$2a$")
    })
}
```

#### 2. DAL 层测试（使用 mock 或 testcontainers）

```go
package user

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserDALImpl_CreateUser(t *testing.T) {
    // 使用 testcontainers 或内存数据库
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    dal := NewUserDAL(db)

    t.Run("creates user successfully", func(t *testing.T) {
        user := &models.User{
            Username: "testuser",
            Email:    "test@example.com",
        }

        err := dal.CreateUser(context.Background(), user)

        require.NoError(t, err)
        assert.NotZero(t, user.ID)
    })

    t.Run("returns error for duplicate email", func(t *testing.T) {
        // 测试逻辑
    })
}
```

#### 3. Logic 层测试（使用 mock DAL）

```go
package user

import (
    "context"
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
)

func TestUserLogicImpl_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDAL := mocks.NewMockUserDAL(ctrl)
    logic := NewUserLogic(mockDAL)

    t.Run("creates user successfully", func(t *testing.T) {
        req := &CreateUserRequest{
            Username: "testuser",
            Email:    "test@example.com",
        }

        mockDAL.EXPECT().
            CreateUser(gomock.Any(), gomock.Any()).
            DoAndReturn(func(ctx context.Context, user *models.User) error {
                user.ID = 1
                return nil
            })

        err := logic.CreateUser(context.Background(), req)

        assert.NoError(t, err)
    })
}
```

#### 4. Middleware 测试

```go
package middleware

import (
    "context"
    "testing"
    "bytes"
    "github.com/rs/zerolog"
    "github.com/stretchr/testify/assert"
)

func TestMetaInfoMiddleware_ServerMiddleware(t *testing.T) {
    t.Run("ensures request_id exists", func(t *testing.T) {
        var logBuf bytes.Buffer
        logger := zerolog.New(&logBuf).With().Timestamp().Logger()

        middleware := NewMetaInfoMiddleware(&logger)

        ctx := context.Background()
        resultCtx := middleware.ensureRequestID(ctx)

        requestID := GetRequestID(resultCtx)
        assert.NotEmpty(t, requestID)
    })
}
```

### 使用 testfixtures 设置测试数据

```go
package user

import (
    "testing"
    "github.com/stretchr/testify/require"
    "gorm.io/gorm"
)

func setupTestUsers(t *testing.T, db *gorm.DB) {
    users := []models.User{
        {Username: "admin", Email: "admin@example.com"},
        {Username: "user", Email: "user@example.com"},
    }

    for _, user := range users {
        err := db.Create(&user).Error
        require.NoError(t, err)
    }
}

func TestUserDALImpl_GetUserByEmail(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    setupTestUsers(t, db)
    dal := NewUserDAL(db)

    user, err := dal.GetUserByEmail(context.Background(), "admin@example.com")

    require.NoError(t, err)
    assert.Equal(t, "admin", user.Username)
}
```

## 测试覆盖率

### 查看覆盖率

```bash
# 生成覆盖率报告
go test ./... -coverprofile=coverage.out

# 查看函数级别覆盖率
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html
```

### 覆盖率目标

| 组件类型 | 最低覆盖率 | 推荐覆盖率 |
|---------|----------|-----------|
| 核心业务逻辑 | 80% | 90%+ |
| 工具函数 | 85% | 95%+ |
| 数据访问层 | 70% | 80%+ |
| HTTP/RPC Handler | 60% | 75%+ |
| 中间件 | 75% | 85%+ |

### 提高覆盖率的技巧

1. **测试边界条件**：
   ```go
   t.Run("handles empty input", func(t *testing.T) {})
   t.Run("handles nil input", func(t *testing.T) {})
   t.Run("handles maximum values", func(t *testing.T) {})
   ```

2. **测试错误路径**：
   ```go
   t.Run("returns error on validation failure", func(t *testing.T) {})
   t.Run("returns error on database failure", func(t *testing.T) {})
   ```

3. **使用表格驱动测试覆盖多种场景**：
   ```go
   tests := []struct {
       name    string
       input   string
       wantErr bool
   }{
       {"valid input", "valid", false},
       {"empty input", "", true},
       {"invalid format", "invalid", true},
   }
   ```

## 测试最佳实践

### 1. 使用 testify 库

```go
import (
    "github.com/stretchr/testify/assert"  // 断言
    "github.com/stretchr/testify/require" // 必须通过的断言
    "github.com/stretchr/testify/suite"  // 测试套件
)

// assert - 失败后继续执行
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.True(t, condition)

// require - 失败后立即停止
require.NoError(t, err)
require.NotNil(t, result)
```

### 2. 表格驱动测试

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"missing domain", "user@", true},
        {"empty string", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3. 使用 setup/teardown

```go
func TestMain(m *testing.M) {
    // 全局 setup
    setupTestDatabase()

    // 运行测试
    code := m.Run()

    // 全局 teardown
    teardownTestDatabase()

    os.Exit(code)
}

func TestUserDAL(t *testing.T) {
    // 每个 test 的 setup
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    t.Run("test case 1", func(t *testing.T) {
        // 使用 db
    })
}
```

### 4. Mock 外部依赖

```go
// 使用 gomock 生成 mock
//go:generate mockgen -source=dal.go -destination=mocks/mock_dal.go -package=mocks

func TestUserLogic(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDAL := mocks.NewMockUserDAL(ctrl)

    mockDAL.EXPECT().
        GetUser(gomock.Any(), gomock.Eq("123")).
        Return(&models.User{ID: "123"}, nil)

    // 使用 mockDAL 测试 logic
}
```

### 5. 测试并发安全

```go
func TestConcurrentAccess(t *testing.T) {
    cache := NewCache()

    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            cache.Set("key", "value")
        }()
    }
    wg.Wait()

    // 验证结果
    assert.Equal(t, "value", cache.Get("key"))
}
```

### 6. 使用测试覆盖率标签

```go
// +build integration

package integration

import (
    "testing"
)

func TestIntegrationWithDatabase(t *testing.T) {
    // 集成测试
}
```

运行：
```bash
# 单元测试（不运行集成测试）
go test ./... -short

# 包含集成测试
go test ./...
```

## CI/CD 集成

### GitHub Actions 示例

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: test_db
          POSTGRES_USER: test_user
          POSTGRES_PASSWORD: test_pass
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Run tests
        run: |
          go test ./... -v -coverprofile=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

## 常见问题

### Q: 如何测试数据库操作？

A: 使用 testcontainers 或内存数据库（如 SQLite）进行隔离测试。

### Q: 如何测试 HTTP Handler？

A: 使用 httptest 包创建测试服务器：

```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
}
```

### Q: 如何模拟时间？

A: 使用可注入的时间接口：

```go
type TimeProvider interface {
    Now() time.Time
}

type TestTimeProvider struct {
    FixedTime time.Time
}

func (t *TestTimeProvider) Now() time.Time {
    return t.FixedTime
}
```

## 参考资料

- [Go Testing 官方文档](https://golang.org/pkg/testing/)
- [Testify 项目](https://github.com/stretchr/testify)
- [Go Mock 教程](https://github.com/golang/mock)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
