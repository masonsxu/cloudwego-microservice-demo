package department

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	departmentDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/department"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/mock"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// setupTest 初始化测试环境
func setupTest(t *testing.T) (*LogicImpl, *mock.TestMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)
	mocks := mock.NewTestMocks(ctrl)
	logic := &LogicImpl{
		dal:       mocks.DAL,
		converter: mocks.Converter,
	}

	return logic, mocks
}

// assertErrCode 断言错误码匹配
func assertErrCode(t *testing.T, expected errno.ErrNo, actual error) {
	t.Helper()

	errNo, ok := actual.(errno.ErrNo)
	require.True(t, ok, "expected errno.ErrNo, got %T: %v", actual, actual)
	assert.Equal(t, expected.ErrCode, errNo.ErrCode)
}

// ============================================================================
// CreateDepartment 测试
// ============================================================================

func TestLogicImpl_CreateDepartment(t *testing.T) {
	orgID := uuid.New().String()
	deptName := "急诊科"
	deptType := "clinical"

	t.Run("成功创建部门", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateDepartmentRequest{
			OrganizationID: &orgID,
			Name:           &deptName,
			DepartmentType: &deptType,
		}

		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.DeptRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.CreateDepartment(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, deptName, result.GetName())
	})

	t.Run("组织不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateDepartmentRequest{
			OrganizationID: &orgID,
			Name:           &deptName,
		}

		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(false, nil)

		result, err := logic.CreateDepartment(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOrganizationNotFound, err)
	})

	t.Run("检查组织存在时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateDepartmentRequest{
			OrganizationID: &orgID,
			Name:           &deptName,
		}

		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(false, gorm.ErrInvalidDB)

		result, err := logic.CreateDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("缺少部门名称", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateDepartmentRequest{
			OrganizationID: &orgID,
		}

		result, err := logic.CreateDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("缺少组织ID", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateDepartmentRequest{
			Name: &deptName,
		}

		result, err := logic.CreateDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("事务中创建失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateDepartmentRequest{
			OrganizationID: &orgID,
			Name:           &deptName,
		}

		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.DeptRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.CreateDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// GetDepartment 测试
// ============================================================================

func TestLogicImpl_GetDepartment(t *testing.T) {
	deptID := uuid.New().String()

	t.Run("成功获取部门", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		dept := &models.Department{
			BaseModel: models.BaseModel{ID: uuid.MustParse(deptID)},
			Name:      "急诊科",
		}

		req := &identity_srv.GetDepartmentRequest{DepartmentID: &deptID}
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(dept, nil)

		result, err := logic.GetDepartment(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "急诊科", result.GetName())
	})

	t.Run("部门不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetDepartmentRequest{DepartmentID: &deptID}
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.GetDepartment(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrDepartmentNotFound, err)
	})

	t.Run("部门ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetDepartmentRequest{}

		result, err := logic.GetDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetDepartmentRequest{DepartmentID: &deptID}
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.GetDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// UpdateDepartment 测试
// ============================================================================

func TestLogicImpl_UpdateDepartment(t *testing.T) {
	deptID := uuid.New().String()
	orgID := uuid.New()
	newName := "内科"

	t.Run("成功更新部门", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingDept := &models.Department{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(deptID)},
			Name:           "急诊科",
			OrganizationID: orgID,
		}

		req := &identity_srv.UpdateDepartmentRequest{
			DepartmentID: &deptID,
			Name:         &newName,
		}

		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(existingDept, nil)
		mocks.DeptRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.UpdateDepartment(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("部门不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateDepartmentRequest{
			DepartmentID: &deptID,
			Name:         &newName,
		}

		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.UpdateDepartment(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrDepartmentNotFound, err)
	})

	t.Run("部门ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateDepartmentRequest{
			Name: &newName,
		}

		result, err := logic.UpdateDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("更新事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingDept := &models.Department{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(deptID)},
			Name:           "急诊科",
			OrganizationID: orgID,
		}

		req := &identity_srv.UpdateDepartmentRequest{
			DepartmentID: &deptID,
			Name:         &newName,
		}

		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(existingDept, nil)
		mocks.DeptRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.UpdateDepartment(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// DeleteDepartment 测试
// ============================================================================

func TestLogicImpl_DeleteDepartment(t *testing.T) {
	deptID := uuid.New().String()

	t.Run("成功删除部门", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DeptRepo.EXPECT().ExistsByID(ctx, deptID).Return(true, nil)
		mocks.MembershipRepo.EXPECT().CountByDepartmentID(ctx, deptID).Return(int64(0), nil)
		mocks.DeptRepo.EXPECT().SoftDelete(gomock.Any(), deptID).Return(nil)

		err := logic.DeleteDepartment(ctx, deptID)

		assert.NoError(t, err)
	})

	t.Run("部门不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DeptRepo.EXPECT().ExistsByID(ctx, deptID).Return(false, nil)

		err := logic.DeleteDepartment(ctx, deptID)

		assertErrCode(t, errno.ErrDepartmentNotFound, err)
	})

	t.Run("部门有成员无法删除", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DeptRepo.EXPECT().ExistsByID(ctx, deptID).Return(true, nil)
		mocks.MembershipRepo.EXPECT().CountByDepartmentID(ctx, deptID).Return(int64(5), nil)

		err := logic.DeleteDepartment(ctx, deptID)

		assertErrCode(t, errno.ErrCannotDeleteDepartmentWithMembers, err)
	})

	t.Run("空部门ID", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		err := logic.DeleteDepartment(ctx, "")

		assert.Error(t, err)
	})

	t.Run("检查存在时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DeptRepo.EXPECT().ExistsByID(ctx, deptID).Return(false, gorm.ErrInvalidDB)

		err := logic.DeleteDepartment(ctx, deptID)

		assert.Error(t, err)
	})

	t.Run("查询成员数时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DeptRepo.EXPECT().ExistsByID(ctx, deptID).Return(true, nil)
		mocks.MembershipRepo.EXPECT().CountByDepartmentID(ctx, deptID).Return(int64(0), gorm.ErrInvalidDB)

		err := logic.DeleteDepartment(ctx, deptID)

		assert.Error(t, err)
	})
}

// ============================================================================
// GetDepartmentsByOrganization 测试
// ============================================================================

func TestLogicImpl_GetDepartmentsByOrganization(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("成功获取组织下部门列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		departments := []*models.Department{
			{BaseModel: models.BaseModel{ID: uuid.New()}, Name: "急诊科", OrganizationID: uuid.MustParse(orgID)},
			{BaseModel: models.BaseModel{ID: uuid.New()}, Name: "内科", OrganizationID: uuid.MustParse(orgID)},
		}
		pageResult := &models.PageResult{Total: 2, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.GetOrganizationDepartmentsRequest{
			OrganizationID: &orgID,
		}

		mocks.DeptRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *departmentDAL.DepartmentQueryConditions,
			) ([]*models.Department, *models.PageResult, error) {
				assert.Equal(t, orgID, *cond.OrganizationID)
				return departments, pageResult, nil
			})

		result, err := logic.GetDepartmentsByOrganization(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Departments, 2)
		assert.NotNil(t, result.Page)
	})

	t.Run("组织ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetOrganizationDepartmentsRequest{}

		result, err := logic.GetDepartmentsByOrganization(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("组织下无部门", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetOrganizationDepartmentsRequest{
			OrganizationID: &orgID,
		}

		mocks.DeptRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.Department{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 0}, nil)

		result, err := logic.GetDepartmentsByOrganization(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Departments)
	})

	t.Run("带分页参数查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		page := int32(2)
		limit := int32(10)
		req := &identity_srv.GetOrganizationDepartmentsRequest{
			OrganizationID: &orgID,
			Page:           &rpc_base.PageRequest{Page: &page, Limit: &limit},
		}

		mocks.DeptRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.Department{}, &models.PageResult{Total: 15, Page: 2, Limit: 10, TotalPages: 2}, nil)

		result, err := logic.GetDepartmentsByOrganization(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result.Page)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetOrganizationDepartmentsRequest{
			OrganizationID: &orgID,
		}

		mocks.DeptRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(nil, nil, gorm.ErrInvalidDB)

		result, err := logic.GetDepartmentsByOrganization(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// NewLogic 构造函数测试
// ============================================================================

func TestNewLogic(t *testing.T) {
	ctrl := gomock.NewController(t)
	mocks := mock.NewTestMocks(ctrl)

	logic := NewLogic(mocks.DAL, converter.NewConverter())

	assert.NotNil(t, logic)
}
