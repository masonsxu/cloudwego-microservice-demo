package organization

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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

type OrganizationRepositoryTestSuite struct {
	suite.Suite
	db        *gorm.DB
	repo      OrganizationRepository
	container testcontainers.Container
	sqlDB     interface{}
	cleanup   func()
}

func (s *OrganizationRepositoryTestSuite) SetupSuite() {
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

	err = db.AutoMigrate(&models.Organization{})
	require.NoError(s.T(), err, "数据库迁移失败")

	s.db = db
	s.repo = NewOrganizationRepository(db)
	s.container = container
	s.sqlDB = sqlDB
	s.cleanup = func() {
		_ = container.Terminate(ctx)
	}
}

func (s *OrganizationRepositoryTestSuite) TearDownSuite() {
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

func (s *OrganizationRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM organizations")
}

func (s *OrganizationRepositoryTestSuite) TestCreate_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name:         "测试组织",
		Code:         "TEST001",
		FacilityType: "医院",
	}

	err := s.repo.Create(ctx, org)

	require.NoError(s.T(), err)
	assert.NotZero(s.T(), org.ID)
	assert.Equal(s.T(), "TEST001", org.Code)
	assert.Equal(s.T(), "测试组织", org.Name)
}

func (s *OrganizationRepositoryTestSuite) TestGetByID_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name:         "测试组织",
		Code:         "TEST001",
		FacilityType: "医院",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByID(ctx, org.ID.String())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), org.ID, found.ID)
	assert.Equal(s.T(), "TEST001", found.Code)
}

func (s *OrganizationRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *OrganizationRepositoryTestSuite) TestGetByCode_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name:         "测试组织",
		Code:         "TEST001",
		FacilityType: "医院",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByCode(ctx, "TEST001")

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), "TEST001", found.Code)
	assert.Equal(s.T(), "测试组织", found.Name)
}

func (s *OrganizationRepositoryTestSuite) TestGetByCode_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByCode(ctx, "NONEXISTENT")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *OrganizationRepositoryTestSuite) TestUpdate_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name:         "旧名称",
		Code:         "TEST001",
		FacilityType: "医院",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	org.Name = "新名称"
	org.FacilityType = "专科医院"

	err = s.repo.Update(ctx, org)

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, org.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "新名称", updated.Name)
}

func (s *OrganizationRepositoryTestSuite) TestSoftDelete_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	err = s.repo.SoftDelete(ctx, org.ID.String())

	require.NoError(s.T(), err)

	_, err = s.repo.GetByID(ctx, org.ID.String())
	assert.Error(s.T(), err)
}

func (s *OrganizationRepositoryTestSuite) TestUpdateParent_Success() {
	ctx := context.Background()

	parentOrg := &models.Organization{
		Name: "父组织",
		Code: "PARENT",
	}
	err := s.repo.Create(ctx, parentOrg)
	require.NoError(s.T(), err)

	childOrg := &models.Organization{
		Name: "子组织",
		Code: "CHILD",
	}
	err = s.repo.Create(ctx, childOrg)
	require.NoError(s.T(), err)

	// 创建时不设置ParentID
	err = s.repo.UpdateParent(ctx, childOrg.ID.String(), parentOrg.ID.String())
	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, childOrg.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), parentOrg.ID, updated.ParentID)
}

func (s *OrganizationRepositoryTestSuite) TestUpdateParent_TwoLevelLimit() {
	ctx := context.Background()

	grandparent := &models.Organization{
		Name: "根组织",
		Code: "ROOT",
	}
	err := s.repo.Create(ctx, grandparent)
	require.NoError(s.T(), err)

	parent := &models.Organization{
		Name:     "父组织",
		Code:     "PARENT",
		ParentID: grandparent.ID,
	}
	err = s.repo.Create(ctx, parent)
	require.NoError(s.T(), err)

	child := &models.Organization{
		Name: "子组织",
		Code: "CHILD",
	}
	err = s.repo.Create(ctx, child)
	require.NoError(s.T(), err)

	err = s.repo.UpdateParent(ctx, child.ID.String(), parent.ID.String())

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "不支持超过2级的组织层级")
}

func (s *OrganizationRepositoryTestSuite) TestUpdateParent_NonExistentParent() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	err = s.repo.UpdateParent(ctx, org.ID.String(), "00000000-0000-0000-0000-000000000000")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "新父组织不存在")
}

func (s *OrganizationRepositoryTestSuite) TestHasChildren_True() {
	ctx := context.Background()

	parent := &models.Organization{
		Name: "父组织",
		Code: "PARENT",
	}
	err := s.repo.Create(ctx, parent)
	require.NoError(s.T(), err)

	child := &models.Organization{
		Name:     "子组织",
		Code:     "CHILD",
		ParentID: parent.ID,
	}
	err = s.repo.Create(ctx, child)
	require.NoError(s.T(), err)

	hasChildren, err := s.repo.HasChildren(ctx, parent.ID.String())

	require.NoError(s.T(), err)
	assert.True(s.T(), hasChildren)
}

func (s *OrganizationRepositoryTestSuite) TestHasChildren_False() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	hasChildren, err := s.repo.HasChildren(ctx, org.ID.String())

	require.NoError(s.T(), err)
	assert.False(s.T(), hasChildren)
}

func (s *OrganizationRepositoryTestSuite) TestCountChildren_Success() {
	ctx := context.Background()

	parent := &models.Organization{
		Name: "父组织",
		Code: "PARENT",
	}
	err := s.repo.Create(ctx, parent)
	require.NoError(s.T(), err)

	child1 := &models.Organization{
		Name:     "子组织1",
		Code:     "CHILD1",
		ParentID: parent.ID,
	}
	err = s.repo.Create(ctx, child1)
	require.NoError(s.T(), err)

	child2 := &models.Organization{
		Name:     "子组织2",
		Code:     "CHILD2",
		ParentID: parent.ID,
	}
	err = s.repo.Create(ctx, child2)
	require.NoError(s.T(), err)

	count, err := s.repo.CountChildren(ctx, parent.ID.String())

	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), count)
}

func (s *OrganizationRepositoryTestSuite) TestExistsByID_True() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	exists, err := s.repo.ExistsByID(ctx, org.ID.String())

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *OrganizationRepositoryTestSuite) TestExistsByID_False() {
	ctx := context.Background()

	exists, err := s.repo.ExistsByID(ctx, "00000000-0000-0000-0000-000000000000")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *OrganizationRepositoryTestSuite) TestCheckNameConflict_Conflict() {
	ctx := context.Background()

	org1 := &models.Organization{
		Name: "冲突名称",
		Code: "ORG1",
	}
	err := s.repo.Create(ctx, org1)
	require.NoError(s.T(), err)

	org2 := &models.Organization{
		Name: "其他名称",
		Code: "ORG2",
	}
	err = s.repo.Create(ctx, org2)
	require.NoError(s.T(), err)

	conflict, err := s.repo.CheckNameConflict(ctx, "冲突名称", uuid.Nil.String(), org2.ID.String())

	require.NoError(s.T(), err)
	assert.True(s.T(), conflict)
}

func (s *OrganizationRepositoryTestSuite) TestCheckNameConflict_NoConflict() {
	ctx := context.Background()

	conflict, err := s.repo.CheckNameConflict(ctx, "不存在名称", uuid.Nil.String())

	require.NoError(s.T(), err)
	assert.False(s.T(), conflict)
}

func (s *OrganizationRepositoryTestSuite) TestCheckNameConflict_ExcludeID() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST",
	}
	err := s.repo.Create(ctx, org)
	require.NoError(s.T(), err)

	conflict, err := s.repo.CheckNameConflict(ctx, "测试组织", uuid.Nil.String(), org.ID.String())

	require.NoError(s.T(), err)
	assert.False(s.T(), conflict)
}

func (s *OrganizationRepositoryTestSuite) TestFindWithConditions_CodeFilter() {
	ctx := context.Background()

	org1 := &models.Organization{
		Name: "组织1",
		Code: "CODE001",
	}
	err := s.repo.Create(ctx, org1)
	require.NoError(s.T(), err)

	org2 := &models.Organization{
		Name: "组织2",
		Code: "CODE002",
	}
	err = s.repo.Create(ctx, org2)
	require.NoError(s.T(), err)

	code := "CODE001"
	conditions := &OrganizationQueryConditions{
		Code: &code,
		Page: base.NewQueryOptions(),
	}

	orgs, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.Len(s.T(), orgs, 1)
	assert.Equal(s.T(), "CODE001", orgs[0].Code)
}

func (s *OrganizationRepositoryTestSuite) TestFindWithConditions_NameFilter() {
	ctx := context.Background()

	org1 := &models.Organization{
		Name: "北京医院",
		Code: "BJ001",
	}
	err := s.repo.Create(ctx, org1)
	require.NoError(s.T(), err)

	org2 := &models.Organization{
		Name: "上海医院",
		Code: "SH001",
	}
	err = s.repo.Create(ctx, org2)
	require.NoError(s.T(), err)

	name := "北京"
	conditions := &OrganizationQueryConditions{
		Name: &name,
		Page: base.NewQueryOptions(),
	}

	orgs, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(orgs), 1)
	assert.Contains(s.T(), orgs[0].Name, "北京")
}

func (s *OrganizationRepositoryTestSuite) TestFindWithConditions_RootOrganizations() {
	ctx := context.Background()

	rootOrg := &models.Organization{
		Name: "根组织",
		Code: "ROOT",
	}
	err := s.repo.Create(ctx, rootOrg)
	require.NoError(s.T(), err)

	childOrg := &models.Organization{
		Name:     "子组织",
		Code:     "CHILD",
		ParentID: rootOrg.ID,
	}
	err = s.repo.Create(ctx, childOrg)
	require.NoError(s.T(), err)

	emptyParentID := ""
	conditions := &OrganizationQueryConditions{
		ParentID: &emptyParentID,
		Page:     base.NewQueryOptions(),
	}

	orgs, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(orgs), 1)

	foundRoot := false

	for _, org := range orgs {
		if org.ID == rootOrg.ID {
			foundRoot = true
			break
		}
	}

	assert.True(s.T(), foundRoot)
}

func (s *OrganizationRepositoryTestSuite) TestFindWithConditions_Search() {
	ctx := context.Background()

	org1 := &models.Organization{
		Name:         "测试医院",
		Code:         "TEST001",
		FacilityType: "综合医院",
	}
	err := s.repo.Create(ctx, org1)
	require.NoError(s.T(), err)

	org2 := &models.Organization{
		Name:         "诊所",
		Code:         "CLINIC001",
		FacilityType: "诊所",
	}
	err = s.repo.Create(ctx, org2)
	require.NoError(s.T(), err)

	conditions := &OrganizationQueryConditions{
		Page: &base.QueryOptions{
			Search:   "医院",
			Page:     1,
			PageSize: 10,
		},
	}

	orgs, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(orgs), 1)
	assert.Contains(s.T(), orgs[0].FacilityType, "医院")
}

func (s *OrganizationRepositoryTestSuite) TestFindAll_Pagination() {
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		org := &models.Organization{
			Name: "组织" + string(rune('0'+i)),
			Code: "ORG00" + string(rune('0'+i)),
		}
		err := s.repo.Create(ctx, org)
		require.NoError(s.T(), err)
	}

	opts := &base.QueryOptions{
		Page:     1,
		PageSize: 2,
	}

	orgs, pageResult, err := s.repo.FindAll(ctx, opts)

	require.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(orgs), 2)
	assert.GreaterOrEqual(s.T(), pageResult.Total, int32(5))
}

func TestOrganizationRepositorySuite(t *testing.T) {
	suite.Run(t, new(OrganizationRepositoryTestSuite))
}
