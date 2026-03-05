package definition

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

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

type RoleDefinitionRepositoryTestSuite struct {
	suite.Suite
	db        *gorm.DB
	repo      RoleDefinitionRepository
	container testcontainers.Container
	sqlDB     interface{}
	cleanup   func()
}

func (s *RoleDefinitionRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "test_db",
				"POSTGRES_USER":     "test_user",
				"POSTGRES_PASSWORD": "test_pass",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10 * time.Second),
		},
		Started: true,
	})
	require.NoError(s.T(), err, "启动 PostgreSQL 容器失败")

	host, err := container.Host(ctx)
	require.NoError(s.T(), err, "获取容器主机地址失败")

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(s.T(), err, "获取容器映射端口失败")

	dsn := "host=" + host + " port=" + port.Port() +
		" user=test_user password=test_pass dbname=test_db sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(s.T(), err, "连接数据库失败")

	sqlDB, err := db.DB()
	require.NoError(s.T(), err, "获取 SQL DB 失败")

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = db.AutoMigrate(&models.RoleDefinition{})
	require.NoError(s.T(), err, "数据库迁移失败")

	s.db = db
	s.repo = NewRoleDefinitionRepository(db)
	s.container = container
	s.sqlDB = sqlDB
	s.cleanup = func() {
		_ = container.Terminate(ctx)
	}
}

func (s *RoleDefinitionRepositoryTestSuite) TearDownSuite() {
	if s.cleanup != nil {
		s.cleanup()
	}

	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func (s *RoleDefinitionRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM role_definitions")
}

func (s *RoleDefinitionRepositoryTestSuite) TestCreate_Success() {
	ctx := context.Background()

	role := &models.RoleDefinition{
		Name:         "测试角色",
		RoleCode:     "TEST_ROLE",
		Description:  "这是一个测试角色",
		Status:       models.RoleStatusActive,
		IsSystemRole: false,
	}

	err := s.repo.Create(ctx, role)

	require.NoError(s.T(), err)
	assert.NotZero(s.T(), role.ID)
	assert.Equal(s.T(), "TEST_ROLE", role.RoleCode)
	assert.Equal(s.T(), "测试角色", role.Name)
}

func (s *RoleDefinitionRepositoryTestSuite) TestGetByID_Success() {
	ctx := context.Background()

	role := &models.RoleDefinition{
		Name:     "测试角色",
		RoleCode: "TEST_ROLE",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByID(ctx, role.ID.String())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), role.ID, found.ID)
	assert.Equal(s.T(), "TEST_ROLE", found.RoleCode)
}

func (s *RoleDefinitionRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindByName_Success() {
	ctx := context.Background()

	role := &models.RoleDefinition{
		Name:     "管理员",
		RoleCode: "ADMIN",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role)
	require.NoError(s.T(), err)

	found, err := s.repo.FindByName(ctx, "管理员")

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), "管理员", found.Name)
	assert.Equal(s.T(), "ADMIN", found.RoleCode)
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindByName_NotFound() {
	ctx := context.Background()

	found, err := s.repo.FindByName(ctx, "不存在的角色")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *RoleDefinitionRepositoryTestSuite) TestCheckNameExists_True() {
	ctx := context.Background()

	role := &models.RoleDefinition{
		Name:     "已存在角色",
		RoleCode: "EXIST",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckNameExists(ctx, "已存在角色")

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *RoleDefinitionRepositoryTestSuite) TestCheckNameExists_False() {
	ctx := context.Background()

	exists, err := s.repo.CheckNameExists(ctx, "不存在的角色")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *RoleDefinitionRepositoryTestSuite) TestUpdate_Success() {
	ctx := context.Background()

	role := &models.RoleDefinition{
		Name:        "旧名称",
		RoleCode:    "ROLE",
		Description: "旧描述",
		Status:      models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role)
	require.NoError(s.T(), err)

	role.Name = "新名称"
	role.Description = "新描述"

	err = s.repo.Update(ctx, role)

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, role.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "新名称", updated.Name)
	assert.Equal(s.T(), "新描述", updated.Description)
}

func (s *RoleDefinitionRepositoryTestSuite) TestSoftDelete_Success() {
	ctx := context.Background()

	role := &models.RoleDefinition{
		Name:     "测试角色",
		RoleCode: "TEST",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role)
	require.NoError(s.T(), err)

	err = s.repo.SoftDelete(ctx, role.ID.String())

	require.NoError(s.T(), err)

	_, err = s.repo.GetByID(ctx, role.ID.String())
	assert.Error(s.T(), err)
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindByStatus_Success() {
	ctx := context.Background()

	role1 := &models.RoleDefinition{
		Name:     "活跃角色1",
		RoleCode: "ACTIVE1",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role1)
	require.NoError(s.T(), err)

	role2 := &models.RoleDefinition{
		Name:     "活跃角色2",
		RoleCode: "ACTIVE2",
		Status:   models.RoleStatusActive,
	}
	err = s.repo.Create(ctx, role2)
	require.NoError(s.T(), err)

	role3 := &models.RoleDefinition{
		Name:     "禁用角色",
		RoleCode: "DISABLED",
		Status:   models.RoleStatusInactive,
	}
	err = s.repo.Create(ctx, role3)
	require.NoError(s.T(), err)

	opts := &base.QueryOptions{
		Page:     1,
		PageSize: 10,
	}

	roles, pageResult, err := s.repo.FindByStatus(ctx, models.RoleStatusActive, opts)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(roles), 2)
	assert.GreaterOrEqual(s.T(), pageResult.Total, int32(2))

	for _, role := range roles {
		assert.Equal(s.T(), models.RoleStatusActive, role.Status)
	}
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindBySystemRole_SystemRoles() {
	ctx := context.Background()

	role1 := &models.RoleDefinition{
		Name:         "系统管理员",
		RoleCode:     "SYS_ADMIN",
		IsSystemRole: true,
		Status:       models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role1)
	require.NoError(s.T(), err)

	role2 := &models.RoleDefinition{
		Name:         "普通用户",
		RoleCode:     "USER",
		IsSystemRole: false,
		Status:       models.RoleStatusActive,
	}
	err = s.repo.Create(ctx, role2)
	require.NoError(s.T(), err)

	opts := &base.QueryOptions{
		Page:     1,
		PageSize: 10,
	}

	roles, pageResult, err := s.repo.FindBySystemRole(ctx, true, opts)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(roles), 1)
	assert.GreaterOrEqual(s.T(), pageResult.Total, int32(1))

	for _, role := range roles {
		assert.True(s.T(), role.IsSystemRole)
	}
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindWithConditions_NameFilter() {
	ctx := context.Background()

	role1 := &models.RoleDefinition{
		Name:     "管理员",
		RoleCode: "ADMIN",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role1)
	require.NoError(s.T(), err)

	role2 := &models.RoleDefinition{
		Name:     "普通用户",
		RoleCode: "USER",
		Status:   models.RoleStatusActive,
	}
	err = s.repo.Create(ctx, role2)
	require.NoError(s.T(), err)

	name := "管理员"
	conditions := &RoleDefinitionQueryConditions{
		Name: &name,
		Page: base.NewQueryOptions(),
	}

	roles, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.Len(s.T(), roles, 1)
	assert.Equal(s.T(), "管理员", roles[0].Name)
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindWithConditions_StatusFilter() {
	ctx := context.Background()

	status := models.RoleStatusActive
	conditions := &RoleDefinitionQueryConditions{
		Status: &status,
		Page:   base.NewQueryOptions(),
	}

	role := &models.RoleDefinition{
		Name:     "活跃角色",
		RoleCode: "ACTIVE",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role)
	require.NoError(s.T(), err)

	roles, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(roles), 1)
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindWithConditions_SystemRoleFilter() {
	ctx := context.Background()

	isSystem := true
	conditions := &RoleDefinitionQueryConditions{
		IsSystemRole: &isSystem,
		Page:         base.NewQueryOptions(),
	}

	role := &models.RoleDefinition{
		Name:         "系统角色",
		RoleCode:     "SYS",
		IsSystemRole: true,
		Status:       models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role)
	require.NoError(s.T(), err)

	roles, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(roles), 1)
	for _, role := range roles {


		assert.True(s.T(), role.IsSystemRole)
	}
}

func (s *RoleDefinitionRepositoryTestSuite) TestCountByStatus_Success() {
	ctx := context.Background()

	role1 := &models.RoleDefinition{
		Name:     "活跃角色1",
		RoleCode: "ACTIVE1",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role1)
	require.NoError(s.T(), err)

	role2 := &models.RoleDefinition{
		Name:     "活跃角色2",
		RoleCode: "ACTIVE2",
		Status:   models.RoleStatusActive,
	}
	err = s.repo.Create(ctx, role2)
	require.NoError(s.T(), err)

	count, err := s.repo.CountByStatus(ctx, models.RoleStatusActive)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), count, int64(2))
}

func (s *RoleDefinitionRepositoryTestSuite) TestListActiveRoles_Success() {
	ctx := context.Background()

	role1 := &models.RoleDefinition{
		Name:     "活跃角色",
		RoleCode: "ACTIVE",
		Status:   models.RoleStatusActive,
	}
	err := s.repo.Create(ctx, role1)
	require.NoError(s.T(), err)

	role2 := &models.RoleDefinition{
		Name:     "禁用角色",
		RoleCode: "INACTIVE",
		Status:   models.RoleStatusInactive,
	}
	err = s.repo.Create(ctx, role2)
	require.NoError(s.T(), err)

	opts := &base.QueryOptions{
		Page:     1,
		PageSize: 10,
	}

	roles, pageResult, err := s.repo.ListActiveRoles(ctx, opts)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(roles), 1)
	assert.GreaterOrEqual(s.T(), pageResult.Total, int32(1))

	for _, role := range roles {
		assert.Equal(s.T(), models.RoleStatusActive, role.Status)
	}
}

func (s *RoleDefinitionRepositoryTestSuite) TestFindAll_Pagination() {
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		role := &models.RoleDefinition{
			Name:     "角色" + string(rune('0'+i)),
			RoleCode: "ROLE00" + string(rune('0'+i)),
			Status:   models.RoleStatusActive,
		}
		err := s.repo.Create(ctx, role)
		require.NoError(s.T(), err)
	}

	opts := &base.QueryOptions{
		Page:     1,
		PageSize: 2,
	}

	roles, pageResult, err := s.repo.FindAll(ctx, opts)

	require.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(roles), 2)
	assert.GreaterOrEqual(s.T(), pageResult.Total, int32(5))
}

func TestRoleDefinitionRepositorySuite(t *testing.T) {
	suite.Run(t, new(RoleDefinitionRepositoryTestSuite))
}
