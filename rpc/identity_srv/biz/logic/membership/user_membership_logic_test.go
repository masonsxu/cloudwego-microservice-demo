package membership

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	membershipDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/membership"
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
// AddMembership 测试
// ============================================================================

func TestLogicImpl_AddMembership(t *testing.T) {
	userID := uuid.New().String()
	orgID := uuid.New().String()
	deptID := uuid.New().String()

	t.Run("成功添加成员关系-无部门", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		isPrimary := false
		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
			IsPrimary:      &isPrimary,
		}

		// 验证实体存在性
		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)

		// 冲突检查 - 无部门时按组织查询
		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(nil, gorm.ErrRecordNotFound)

		// 创建成员关系
		mocks.MembershipRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.AddMembership(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("成功添加成员关系-有部门", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		isPrimary := false
		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
			DepartmentID:   &deptID,
			IsPrimary:      &isPrimary,
		}

		// 验证实体存在性
		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(&models.Department{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(deptID)},
			OrganizationID: uuid.MustParse(orgID),
		}, nil)

		// 冲突检查 - 有部门时按部门查询
		mocks.MembershipRepo.EXPECT().
			GetByUserAndDepartment(ctx, userID, deptID).
			Return(nil, gorm.ErrRecordNotFound)

		// 创建成员关系
		mocks.MembershipRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.AddMembership(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("成功添加主要成员关系-先取消已有主要关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		isPrimary := true
		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
			IsPrimary:      &isPrimary,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(nil, gorm.ErrRecordNotFound)

		// 取消已有主要关系
		mocks.MembershipRepo.EXPECT().UnsetPrimaryByUserID(gomock.Any(), userID).Return(nil)
		// 创建新的成员关系
		mocks.MembershipRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.AddMembership(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyUserID := ""
		req := &identity_srv.AddMembershipRequest{
			UserID:         &emptyUserID,
			OrganizationID: &orgID,
		}

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("组织ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyOrgID := ""
		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &emptyOrgID,
		}

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("用户不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(false, nil)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrUserNotFound, err)
	})

	t.Run("组织不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(false, nil)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOrganizationNotFound, err)
	})

	t.Run("部门不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
			DepartmentID:   &deptID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrDepartmentNotFound, err)
	})

	t.Run("部门不属于指定组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		otherOrgID := uuid.New()
		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
			DepartmentID:   &deptID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(&models.Department{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(deptID)},
			OrganizationID: otherOrgID, // 不同的组织
		}, nil)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("成员关系冲突-组织级别", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)

		// 返回一个活跃的成员关系
		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(&models.UserMembership{
				Status: models.MembershipStatusActive,
			}, nil)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrMembershipAlreadyExists, err)
	})

	t.Run("成员关系冲突-部门级别", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
			DepartmentID:   &deptID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(&models.Department{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(deptID)},
			OrganizationID: uuid.MustParse(orgID),
		}, nil)

		// 按部门查询到已有活跃的成员关系
		mocks.MembershipRepo.EXPECT().
			GetByUserAndDepartment(ctx, userID, deptID).
			Return(&models.UserMembership{
				Status: models.MembershipStatusActive,
			}, nil)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrMembershipAlreadyExists, err)
	})

	t.Run("检查用户存在时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(false, gorm.ErrInvalidDB)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("创建成员关系数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.UserRepo.EXPECT().Exists(ctx, userID).Return(true, nil)
		mocks.OrgRepo.EXPECT().ExistsByID(ctx, orgID).Return(true, nil)
		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(nil, gorm.ErrRecordNotFound)
		mocks.MembershipRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.AddMembership(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("请求为nil", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.AddMembership(ctx, nil)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})
}

// ============================================================================
// UpdateMembership 测试
// ============================================================================

func TestLogicImpl_UpdateMembership(t *testing.T) {
	membershipID := uuid.New().String()
	orgID := uuid.New()
	deptID := uuid.New().String()

	t.Run("成功更新成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		isPrimary := false
		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &membershipID,
			IsPrimary:    &isPrimary,
		}

		existingMembership := &models.UserMembership{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(membershipID)},
			UserID:         uuid.New(),
			OrganizationID: orgID,
			IsPrimary:      false,
			Status:         models.MembershipStatusActive,
		}

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(existingMembership, nil)
		mocks.MembershipRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.UpdateMembership(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("成功更新部门-验证部门归属", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &membershipID,
			DepartmentID: &deptID,
		}

		existingMembership := &models.UserMembership{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(membershipID)},
			UserID:         uuid.New(),
			OrganizationID: orgID,
			IsPrimary:      false,
			Status:         models.MembershipStatusActive,
		}

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(existingMembership, nil)
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(&models.Department{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(deptID)},
			OrganizationID: orgID,
		}, nil)
		mocks.MembershipRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.UpdateMembership(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("更新为主要关系-取消已有主要关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		isPrimary := true
		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &membershipID,
			IsPrimary:    &isPrimary,
		}

		existingUserID := uuid.New()
		existingMembership := &models.UserMembership{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(membershipID)},
			UserID:         existingUserID,
			OrganizationID: orgID,
			IsPrimary:      false, // 当前不是主要关系
			Status:         models.MembershipStatusActive,
		}

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(existingMembership, nil)
		// 先取消已有主要关系
		mocks.MembershipRepo.EXPECT().
			UnsetPrimaryByUserID(gomock.Any(), existingUserID.String()).Return(nil)
		mocks.MembershipRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.UpdateMembership(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("成员关系不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &membershipID,
		}

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.UpdateMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrMembershipNotFound, err)
	})

	t.Run("成员关系ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyID := ""
		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &emptyID,
		}

		result, err := logic.UpdateMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("请求为nil", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.UpdateMembership(ctx, nil)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("部门不属于同一组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &membershipID,
			DepartmentID: &deptID,
		}

		existingMembership := &models.UserMembership{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(membershipID)},
			UserID:         uuid.New(),
			OrganizationID: orgID,
			Status:         models.MembershipStatusActive,
		}

		otherOrgID := uuid.New()

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(existingMembership, nil)
		mocks.DeptRepo.EXPECT().GetByID(ctx, deptID).Return(&models.Department{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(deptID)},
			OrganizationID: otherOrgID, // 不同的组织
		}, nil)

		result, err := logic.UpdateMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("获取现有成员关系数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &membershipID,
		}

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.UpdateMembership(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("更新事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateMembershipRequest{
			MembershipID: &membershipID,
		}

		existingMembership := &models.UserMembership{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(membershipID)},
			UserID:         uuid.New(),
			OrganizationID: orgID,
			IsPrimary:      false,
			Status:         models.MembershipStatusActive,
		}

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(existingMembership, nil)
		mocks.MembershipRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.UpdateMembership(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// RemoveMembership 测试
// ============================================================================

func TestLogicImpl_RemoveMembership(t *testing.T) {
	membershipID := uuid.New().String()

	t.Run("成功删除成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().ExistsByID(ctx, membershipID).Return(true, nil)
		mocks.MembershipRepo.EXPECT().Delete(gomock.Any(), membershipID).Return(nil)

		err := logic.RemoveMembership(ctx, membershipID)

		assert.NoError(t, err)
	})

	t.Run("成员关系不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().ExistsByID(ctx, membershipID).Return(false, nil)

		err := logic.RemoveMembership(ctx, membershipID)

		assertErrCode(t, errno.ErrMembershipNotFound, err)
	})

	t.Run("成员关系ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		err := logic.RemoveMembership(ctx, "")

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("检查存在时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().ExistsByID(ctx, membershipID).Return(false, gorm.ErrInvalidDB)

		err := logic.RemoveMembership(ctx, membershipID)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("删除事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().ExistsByID(ctx, membershipID).Return(true, nil)
		mocks.MembershipRepo.EXPECT().Delete(gomock.Any(), membershipID).Return(gorm.ErrInvalidDB)

		err := logic.RemoveMembership(ctx, membershipID)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// GetMembership 测试
// ============================================================================

func TestLogicImpl_GetMembership(t *testing.T) {
	membershipID := uuid.New().String()

	t.Run("成功获取成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		membership := &models.UserMembership{
			BaseModel:      models.BaseModel{ID: uuid.MustParse(membershipID)},
			UserID:         uuid.New(),
			OrganizationID: uuid.New(),
			IsPrimary:      true,
			Status:         models.MembershipStatusActive,
		}

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(membership, nil)

		result, err := logic.GetMembership(ctx, membershipID)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("成员关系不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.GetMembership(ctx, membershipID)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrMembershipNotFound, err)
	})

	t.Run("成员关系ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.GetMembership(ctx, "")

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().GetByID(ctx, membershipID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.GetMembership(ctx, membershipID)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// GetUserMemberships 测试
// ============================================================================

func TestLogicImpl_GetUserMemberships(t *testing.T) {
	userID := uuid.New().String()
	orgID := uuid.New().String()
	deptID := uuid.New().String()

	t.Run("成功按用户ID查询成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserMembershipsRequest{
			UserID: &userID,
		}

		memberships := []*models.UserMembership{
			{
				BaseModel:      models.BaseModel{ID: uuid.New()},
				UserID:         uuid.MustParse(userID),
				OrganizationID: uuid.New(),
				Status:         models.MembershipStatusActive,
			},
		}
		pageResult := &models.PageResult{Total: 1, Page: 1, Limit: 20, TotalPages: 1}

		mocks.MembershipRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *membershipDAL.UserMembershipQueryConditions,
			) ([]*models.UserMembership, *models.PageResult, error) {
				assert.Equal(t, userID, *cond.UserID)
				return memberships, pageResult, nil
			})

		result, err := logic.GetUserMemberships(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Memberships, 1)
	})

	t.Run("成功按组织ID查询成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserMembershipsRequest{
			OrganizationID: &orgID,
		}

		mocks.MembershipRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *membershipDAL.UserMembershipQueryConditions,
			) ([]*models.UserMembership, *models.PageResult, error) {
				assert.Equal(t, orgID, *cond.OrganizationID)

				return []*models.UserMembership{},
					&models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil
			})

		result, err := logic.GetUserMemberships(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Memberships)
	})

	t.Run("成功按部门ID查询成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserMembershipsRequest{
			DepartmentID: &deptID,
		}

		mocks.MembershipRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *membershipDAL.UserMembershipQueryConditions,
			) ([]*models.UserMembership, *models.PageResult, error) {
				assert.Equal(t, deptID, *cond.DepartmentID)

				return []*models.UserMembership{},
					&models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil
			})

		result, err := logic.GetUserMemberships(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Memberships)
	})

	t.Run("带分页参数查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		page := int32(2)
		limit := int32(10)
		req := &identity_srv.GetUserMembershipsRequest{
			UserID: &userID,
			Page:   &rpc_base.PageRequest{Page: &page, Limit: &limit},
		}

		mocks.MembershipRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.UserMembership{}, &models.PageResult{Total: 15, Page: 2, Limit: 10, TotalPages: 2}, nil)

		result, err := logic.GetUserMemberships(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result.Page)
	})

	t.Run("请求为nil", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.GetUserMemberships(ctx, nil)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("缺少所有查询条件", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserMembershipsRequest{}

		result, err := logic.GetUserMemberships(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserMembershipsRequest{
			UserID: &userID,
		}

		mocks.MembershipRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(nil, nil, gorm.ErrInvalidDB)

		result, err := logic.GetUserMemberships(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("用户ID为空字符串", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyUserID := ""
		req := &identity_srv.GetUserMembershipsRequest{
			UserID: &emptyUserID,
		}

		result, err := logic.GetUserMemberships(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})
}

// ============================================================================
// GetPrimaryMembership 测试
// ============================================================================

func TestLogicImpl_GetPrimaryMembership(t *testing.T) {
	userID := uuid.New().String()

	t.Run("成功获取主要成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		membership := &models.UserMembership{
			BaseModel:      models.BaseModel{ID: uuid.New()},
			UserID:         uuid.MustParse(userID),
			OrganizationID: uuid.New(),
			IsPrimary:      true,
			Status:         models.MembershipStatusActive,
		}

		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(ctx, userID).Return(membership, nil)

		result, err := logic.GetPrimaryMembership(ctx, userID)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("用户没有主要成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.GetPrimaryMembership(ctx, userID)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrMembershipNotFound, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.GetPrimaryMembership(ctx, "")

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(ctx, userID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.GetPrimaryMembership(ctx, userID)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// CheckMembership 测试
// ============================================================================

func TestLogicImpl_CheckMembership(t *testing.T) {
	userID := uuid.New().String()
	orgID := uuid.New().String()
	deptID := uuid.New().String()

	t.Run("成功检查组织成员关系-活跃", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CheckMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(&models.UserMembership{
				Status: models.MembershipStatusActive,
			}, nil)

		isMember, err := logic.CheckMembership(ctx, req)

		require.NoError(t, err)
		assert.True(t, isMember)
	})

	t.Run("成功检查组织成员关系-非活跃", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CheckMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(&models.UserMembership{
				Status: models.MembershipStatusSuspended,
			}, nil)

		isMember, err := logic.CheckMembership(ctx, req)

		require.NoError(t, err)
		assert.False(t, isMember)
	})

	t.Run("成功检查部门成员关系", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CheckMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
			DepartmentID:   &deptID,
		}

		mocks.MembershipRepo.EXPECT().
			GetByUserAndDepartment(ctx, userID, deptID).
			Return(&models.UserMembership{
				Status: models.MembershipStatusActive,
			}, nil)

		isMember, err := logic.CheckMembership(ctx, req)

		require.NoError(t, err)
		assert.True(t, isMember)
	})

	t.Run("成员关系不存在-返回false而非错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CheckMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(nil, gorm.ErrRecordNotFound)

		isMember, err := logic.CheckMembership(ctx, req)

		require.NoError(t, err)
		assert.False(t, isMember)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyUserID := ""
		req := &identity_srv.CheckMembershipRequest{
			UserID:         &emptyUserID,
			OrganizationID: &orgID,
		}

		isMember, err := logic.CheckMembership(ctx, req)

		assert.False(t, isMember)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("组织ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyOrgID := ""
		req := &identity_srv.CheckMembershipRequest{
			UserID:         &userID,
			OrganizationID: &emptyOrgID,
		}

		isMember, err := logic.CheckMembership(ctx, req)

		assert.False(t, isMember)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("请求为nil", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		isMember, err := logic.CheckMembership(ctx, nil)

		assert.False(t, isMember)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CheckMembershipRequest{
			UserID:         &userID,
			OrganizationID: &orgID,
		}

		mocks.MembershipRepo.EXPECT().
			GetByUserAndOrganization(ctx, userID, orgID).
			Return(nil, gorm.ErrInvalidDB)

		isMember, err := logic.CheckMembership(ctx, req)

		assert.False(t, isMember)
		assertErrCode(t, errno.ErrOperationFailed, err)
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
