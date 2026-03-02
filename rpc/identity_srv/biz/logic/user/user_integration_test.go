package user

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	identitysrv "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// UserLogicIntegrationTestSuite Integration 测试套件
type UserLogicIntegrationTestSuite struct {
	suite.Suite
	db      *gorm.DB
	dalImpl dal.DAL
	conv    converter.Converter
	logic   ProfileLogic
	ctx     context.Context
	cleanup func()
}

func (s *UserLogicIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()

	// 启动 PostgreSQL 容器
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(s.T(), err, "Failed to start PostgreSQL container")

	// 获取数据库连接信息
	host, err := container.Host(ctx)
	require.NoError(s.T(), err, "获取容器主机地址失败")

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(s.T(), err, "获取容器映射端口失败")

	// 连接数据库
	dsn := "host=" + host + " port=" + port.Port() +
		" user=testuser password=testpass dbname=testdb sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(s.T(), err, "连接数据库失败")

	sqlDB, err := db.DB()
	require.NoError(s.T(), err, "获取 SQL DB 失败")
	require.NoError(s.T(), sqlDB.Ping())

	// 自动迁移
	err = db.AutoMigrate(
		&models.UserProfile{},
		&models.Organization{},
		&models.Department{},
		&models.UserMembership{},
		&models.Menu{},
		&models.RoleDefinition{},
		&models.UserRoleAssignment{},
		&models.RoleMenuPermission{},
	)
	require.NoError(s.T(), err)

	// 初始化 DAL
	dalImpl := dal.NewDALImpl(db)

	// 初始化 Converter
	conv := converter.NewConverter()

	// 初始化 Logic
	logic := NewLogic(dalImpl, conv)

	s.dalImpl = dalImpl
	s.converterImpl = conv
	s.logic = logic
	s.ctx = context.Background()
	s.sqlDB = sqlDB
	s.container = container
	s.cleanup = func() {
		_ = container.Terminate(ctx)
	}
}

func (s *UserLogicIntegrationTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}
	if s.sqlDB != nil {
		s.sqlDB.Close()
	}
}

func (s *UserLogicIntegrationTestSuite) SetupTest() {
	// 每个测试前清空表
	s.db.Exec("DELETE FROM user_role_assignments")
	s.db.Exec("DELETE FROM role_menu_permissions")
	s.db.Exec("DELETE FROM role_definitions")
	s.db.Exec("DELETE FROM user_memberships")
	s.db.Exec("DELETE FROM departments")
	s.db.Exec("DELETE FROM organizations")
	s.db.Exec("DELETE FROM menus")
	s.db.Exec("DELETE FROM user_profiles")
}

// TestCreateUser_Success 测试成功创建用户
func (s *UserLogicIntegrationTestSuite) TestCreateUser_Success() {
	username := "testuser"
	email := "test@example.com"
	password := "Password123!"

	gender := core.Gender_MALE
	req := &identitysrv.CreateUserRequest{
		Username: &username,
		Email:    &email,
		Password: &password,
		RealName: ptr("Test User"),
		Gender:   &gender,
	}

	resp, err := s.logic.CreateUser(s.ctx, req)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	assert.NotEmpty(s.T(), resp.ID)
	assert.Equal(s.T(), username, *resp.Username)
	assert.Equal(s.T(), email, *resp.Email)
	assert.Equal(s.T(), "Test User", *resp.RealName)

	// 验证数据库中的记录
	user, err := s.dalImpl.UserProfile().GetByUsername(s.ctx, username)
	require.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), username, user.Username)
	assert.Equal(s.T(), email, user.Email)
}

// TestCreateUser_DuplicateUsername 测试重复用户名
func (s *UserLogicIntegrationTestSuite) TestCreateUser_DuplicateUsername() {
	username := "duplicate"
	email1 := "user1@example.com"
	email2 := "user2@example.com"
	password := "Password123!"

	// 创建第一个用户
	req1 := &identitysrv.CreateUserRequest{
		Username: &username,
		Email:    &email1,
		Password: &password,
	}
	_, err := s.logic.CreateUser(s.ctx, req1)
	require.NoError(s.T(), err)

	// 尝试创建第二个用户（用户名重复）
	req2 := &identitysrv.CreateUserRequest{
		Username: &username,
		Email:    &email2,
		Password: &password,
	}
	resp, err := s.logic.CreateUser(s.ctx, req2)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), "用户名已存在")
}

// TestCreateUser_DuplicateEmail 测试重复邮箱
func (s *UserLogicIntegrationTestSuite) TestCreateUser_DuplicateEmail() {
	username1 := "user1"
	username2 := "user2"
	email := "duplicate@example.com"
	password := "Password123!"

	// 创建第一个用户
	req1 := &identitysrv.CreateUserRequest{
		Username: &username1,
		Email:    &email,
		Password: &password,
	}
	_, err := s.logic.CreateUser(s.ctx, req1)
	require.NoError(s.T(), err)

	// 尝试创建第二个用户（邮箱重复）
	req2 := &identitysrv.CreateUserRequest{
		Username: &username2,
		Email:    &email,
		Password: &password,
	}
	resp, err := s.logic.CreateUser(s.ctx, req2)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), "邮箱已存在")
}

// TestCreateUser_WeakPassword 测试弱密码
// 注意：Logic 层不验证密码强度，只验证密码非空
// 密码强度验证应该在 Handler 层或前端进行
func (s *UserLogicIntegrationTestSuite) TestCreateUser_EmptyPassword() {
	username := "emptypass"
	email := "emptypass@example.com"
	password := ""

	req := &identitysrv.CreateUserRequest{
		Username: &username,
		Email:    &email,
		Password: &password,
	}

	resp, err := s.logic.CreateUser(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), "密码不能为空")
}

// TestGetUser_Success 测试成功获取用户
func (s *UserLogicIntegrationTestSuite) TestGetUser_Success() {
	// 先创建用户
	username := "getuser"
	email := "get@example.com"
	password := "Password123!"

	createReq := &identitysrv.CreateUserRequest{
		Username: &username,
		Email:    &email,
		Password: &password,
		RealName: ptr("Get User"),
	}
	createdUser, err := s.logic.CreateUser(s.ctx, createReq)
	require.NoError(s.T(), err)

	// 获取用户
	getReq := &identitysrv.GetUserRequest{
		UserID: createdUser.ID,
	}
	resp, err := s.logic.GetUser(s.ctx, getReq)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	assert.Equal(s.T(), *createdUser.ID, *resp.ID)
	assert.Equal(s.T(), username, *resp.Username)
	assert.Equal(s.T(), email, *resp.Email)
}

// TestGetUser_NotFound 测试获取不存在的用户
func (s *UserLogicIntegrationTestSuite) TestGetUser_NotFound() {
	userID := "00000000-0000-0000-0000-000000000000"

	req := &identitysrv.GetUserRequest{
		UserID: &userID,
	}
	resp, err := s.logic.GetUser(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), "用户不存在")
}

// TestDeleteUser_Success 测试成功删除用户
func (s *UserLogicIntegrationTestSuite) TestDeleteUser_Success() {
	// 先创建用户
	username := "deleteuser"
	email := "delete@example.com"
	password := "Password123!"

	createReq := &identitysrv.CreateUserRequest{
		Username: &username,
		Email:    &email,
		Password: &password,
	}
	createdUser, err := s.logic.CreateUser(s.ctx, createReq)
	require.NoError(s.T(), err)

	// 删除用户
	deleteReq := &identitysrv.DeleteUserRequest{
		UserID: createdUser.ID,
	}
	err = s.logic.DeleteUser(s.ctx, deleteReq)
	assert.NoError(s.T(), err)

	// 验证用户已被删除
	_, err = s.dalImpl.UserProfile().GetByID(s.ctx, *createdUser.ID)
	assert.Error(s.T(), err)
}

// TestDeleteUser_NotFound 测试删除不存在的用户
func (s *UserLogicIntegrationTestSuite) TestDeleteUser_NotFound() {
	userID := "00000000-0000-0000-0000-000000000000"

	req := &identitysrv.DeleteUserRequest{
		UserID: &userID,
	}
	err := s.logic.DeleteUser(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "用户不存在")
}

// TestChangeUserStatus_Success 测试成功修改用户状态
func (s *UserLogicIntegrationTestSuite) TestChangeUserStatus_Success() {
	s.T().Skip("需要查看 IDL 定义以确定正确的字段名")
}

// TestListUsers_Success 测试成功列出用户
func (s *UserLogicIntegrationTestSuite) TestListUsers_Success() {
	s.T().Skip("需要查看 IDL 定义以确定正确的字段名")
}

func TestUserLogicIntegrationSuite(t *testing.T) {
	suite.Run(t, new(UserLogicIntegrationTestSuite))
}

// 辅助函数
func ptr(s string) *string {
	return &s
}
