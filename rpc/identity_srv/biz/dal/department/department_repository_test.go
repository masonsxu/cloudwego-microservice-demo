package department

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

type DepartmentRepositoryTestSuite struct {
	suite.Suite
	db        *gorm.DB
	repo      DepartmentRepository
	container testcontainers.Container
	cleanup   func()
}

func (s *DepartmentRepositoryTestSuite) SetupSuite() {
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

	err = db.AutoMigrate(&models.Department{}, &models.Organization{}, &models.UserMembership{})
	require.NoError(s.T(), err, "数据库迁移失败")

	s.db = db
	s.repo = NewDepartmentRepository(db)
	s.container = container
	s.cleanup = func() {
		_ = container.Terminate(ctx)
	}
}

func (s *DepartmentRepositoryTestSuite) TearDownSuite() {
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

func (s *DepartmentRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM user_memberships")
	s.db.Exec("DELETE FROM departments")
	s.db.Exec("DELETE FROM organizations")
}

func (s *DepartmentRepositoryTestSuite) TestCreate_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept := &models.Department{
		Name:           "测试科室",
		OrganizationID: org.ID,
		DepartmentType: "临床科室",
	}

	err = s.repo.Create(ctx, dept)

	require.NoError(s.T(), err)
	assert.NotZero(s.T(), dept.ID)
	assert.Equal(s.T(), "测试科室", dept.Name)
	assert.Equal(s.T(), "临床科室", dept.DepartmentType)
}

func (s *DepartmentRepositoryTestSuite) TestGetByID_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept := &models.Department{
		Name:           "测试科室",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByID(ctx, dept.ID.String())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), dept.ID, found.ID)
	assert.Equal(s.T(), "测试科室", found.Name)
}

func (s *DepartmentRepositoryTestSuite) TestGetByID_NotFound() {
	ctx := context.Background()

	found, err := s.repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *DepartmentRepositoryTestSuite) TestGetByName_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept := &models.Department{
		Name:           "测试科室",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept)
	require.NoError(s.T(), err)

	found, err := s.repo.GetByName(ctx, "测试科室", org.ID.String())

	require.NoError(s.T(), err)
	require.NotNil(s.T(), found)
	assert.Equal(s.T(), "测试科室", found.Name)
}

func (s *DepartmentRepositoryTestSuite) TestGetByName_NotFound() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	found, err := s.repo.GetByName(ctx, "不存在", org.ID.String())

	assert.Error(s.T(), err)
	assert.Nil(s.T(), found)
}

func (s *DepartmentRepositoryTestSuite) TestUpdate_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept := &models.Department{
		Name:           "旧名称",
		OrganizationID: org.ID,
		DepartmentType: "旧类型",
	}
	err = s.repo.Create(ctx, dept)
	require.NoError(s.T(), err)

	dept.Name = "新名称"
	dept.DepartmentType = "新类型"

	err = s.repo.Update(ctx, dept)

	require.NoError(s.T(), err)

	updated, err := s.repo.GetByID(ctx, dept.ID.String())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "新名称", updated.Name)
	assert.Equal(s.T(), "新类型", updated.DepartmentType)
}

func (s *DepartmentRepositoryTestSuite) TestSoftDelete_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept := &models.Department{
		Name:           "测试科室",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept)
	require.NoError(s.T(), err)

	err = s.repo.SoftDelete(ctx, dept.ID.String())

	require.NoError(s.T(), err)

	_, err = s.repo.GetByID(ctx, dept.ID.String())
	assert.Error(s.T(), err)
}

func (s *DepartmentRepositoryTestSuite) TestExistsByID_True() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept := &models.Department{
		Name:           "测试科室",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept)
	require.NoError(s.T(), err)

	exists, err := s.repo.ExistsByID(ctx, dept.ID.String())

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *DepartmentRepositoryTestSuite) TestExistsByID_False() {
	ctx := context.Background()

	exists, err := s.repo.ExistsByID(ctx, "00000000-0000-0000-0000-000000000000")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *DepartmentRepositoryTestSuite) TestCheckNameExists_True() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept1 := &models.Department{
		Name:           "冲突名称",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept1)
	require.NoError(s.T(), err)

	dept2 := &models.Department{
		Name:           "其他名称",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept2)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckNameExists(ctx, "冲突名称", org.ID.String(), dept2.ID.String())

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *DepartmentRepositoryTestSuite) TestCheckNameExists_False() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckNameExists(ctx, "不存在名称", org.ID.String(), "")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *DepartmentRepositoryTestSuite) TestCheckNameExists_ExcludeID() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept := &models.Department{
		Name:           "测试科室",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept)
	require.NoError(s.T(), err)

	exists, err := s.repo.CheckNameExists(ctx, "测试科室", org.ID.String(), dept.ID.String())

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *DepartmentRepositoryTestSuite) TestCountByOrganization_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept1 := &models.Department{
		Name:           "科室1",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept1)
	require.NoError(s.T(), err)

	dept2 := &models.Department{
		Name:           "科室2",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept2)
	require.NoError(s.T(), err)

	count, err := s.repo.CountByOrganization(ctx, org.ID.String())

	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), count)
}

func (s *DepartmentRepositoryTestSuite) TestCountByDepartmentType_Success() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept1 := &models.Department{
		Name:           "内科",
		OrganizationID: org.ID,
		DepartmentType: "临床科室",
	}
	err = s.repo.Create(ctx, dept1)
	require.NoError(s.T(), err)

	dept2 := &models.Department{
		Name:           "外科",
		OrganizationID: org.ID,
		DepartmentType: "临床科室",
	}
	err = s.repo.Create(ctx, dept2)
	require.NoError(s.T(), err)

	dept3 := &models.Department{
		Name:           "放射科",
		OrganizationID: org.ID,
		DepartmentType: "医技科室",
	}
	err = s.repo.Create(ctx, dept3)
	require.NoError(s.T(), err)

	count, err := s.repo.CountByDepartmentType(ctx, "临床科室")

	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), count)
}

func (s *DepartmentRepositoryTestSuite) TestFindWithConditions_NameFilter() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	dept1 := &models.Department{
		Name:           "内科",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept1)
	require.NoError(s.T(), err)

	dept2 := &models.Department{
		Name:           "外科",
		OrganizationID: org.ID,
	}
	err = s.repo.Create(ctx, dept2)
	require.NoError(s.T(), err)

	name := "内科"
	conditions := &DepartmentQueryConditions{
		Name: &name,
		Page: base.NewQueryOptions(),
	}

	depts, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), len(depts), 1)
	assert.Contains(s.T(), depts[0].Name, "内科")
}

func (s *DepartmentRepositoryTestSuite) TestFindWithConditions_OrganizationFilter() {
	ctx := context.Background()

	org1 := &models.Organization{
		Name: "组织1",
		Code: "ORG001",
	}
	err := s.db.Create(org1).Error
	require.NoError(s.T(), err)

	org2 := &models.Organization{
		Name: "组织2",
		Code: "ORG002",
	}
	err = s.db.Create(org2).Error
	require.NoError(s.T(), err)

	org1ID := org1.ID.String()

	dept1 := &models.Department{
		Name:           "科室1",
		OrganizationID: org1.ID,
	}
	err = s.repo.Create(ctx, dept1)
	require.NoError(s.T(), err)

	dept2 := &models.Department{
		Name:           "科室2",
		OrganizationID: org2.ID,
	}
	err = s.repo.Create(ctx, dept2)
	require.NoError(s.T(), err)

	conditions := &DepartmentQueryConditions{
		OrganizationID: &org1ID,
		Page:           base.NewQueryOptions(),
	}

	depts, _, err := s.repo.FindWithConditions(ctx, conditions)

	require.NoError(s.T(), err)
	assert.Len(s.T(), depts, 1)
	assert.Equal(s.T(), "科室1", depts[0].Name)
}

func (s *DepartmentRepositoryTestSuite) TestFindAll_Pagination() {
	ctx := context.Background()

	org := &models.Organization{
		Name: "测试组织",
		Code: "TEST001",
	}
	err := s.db.Create(org).Error
	require.NoError(s.T(), err)

	for i := 1; i <= 5; i++ {
		dept := &models.Department{
			Name:           "科室" + string(rune('0'+i)),
			OrganizationID: org.ID,
		}
		err := s.repo.Create(ctx, dept)
		require.NoError(s.T(), err)
	}

	opts := &base.QueryOptions{
		Page:     1,
		PageSize: 2,
	}

	depts, pageResult, err := s.repo.FindAll(ctx, opts)

	require.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(depts), 2)
	assert.GreaterOrEqual(s.T(), pageResult.Total, int32(5))
}

func TestDepartmentRepositorySuite(t *testing.T) {
	suite.Run(t, new(DepartmentRepositoryTestSuite))
}
