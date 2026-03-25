package assignment

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/mock"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
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
// AssignRoleToUser 测试
// ============================================================================

func TestLogicImpl_AssignRoleToUser(t *testing.T) {
	userID := uuid.New().String()
	roleID := uuid.New().String()
	assignedBy := uuid.New().String()

	t.Run("成功分配角色", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AssignRoleToUserRequest{
			UserID:     &userID,
			RoleID:     &roleID,
			AssignedBy: &assignedBy,
		}

		mocks.AssignmentRepo.EXPECT().CheckUserRoleExists(ctx, userID, roleID).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.AssignRoleToUser(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.AssignmentID)
	})

	t.Run("成功分配角色_无assignedBy", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		emptyAssignedBy := ""
		req := &identity_srv.AssignRoleToUserRequest{
			UserID:     &userID,
			RoleID:     &roleID,
			AssignedBy: &emptyAssignedBy,
		}

		mocks.AssignmentRepo.EXPECT().CheckUserRoleExists(ctx, userID, roleID).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.AssignRoleToUser(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("角色已分配_重复分配", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AssignRoleToUserRequest{
			UserID:     &userID,
			RoleID:     &roleID,
			AssignedBy: &assignedBy,
		}

		mocks.AssignmentRepo.EXPECT().CheckUserRoleExists(ctx, userID, roleID).Return(true, nil)

		result, err := logic.AssignRoleToUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleAssignmentAlreadyExists, err)
	})

	t.Run("检查角色分配状态时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AssignRoleToUserRequest{
			UserID:     &userID,
			RoleID:     &roleID,
			AssignedBy: &assignedBy,
		}

		mocks.AssignmentRepo.EXPECT().CheckUserRoleExists(ctx, userID, roleID).Return(false, errors.New("db error"))

		result, err := logic.AssignRoleToUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("角色不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AssignRoleToUserRequest{
			UserID:     &userID,
			RoleID:     &roleID,
			AssignedBy: &assignedBy,
		}

		mocks.AssignmentRepo.EXPECT().CheckUserRoleExists(ctx, userID, roleID).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, nil)

		result, err := logic.AssignRoleToUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleDefinitionNotFound, err)
	})

	t.Run("查询角色信息时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AssignRoleToUserRequest{
			UserID:     &userID,
			RoleID:     &roleID,
			AssignedBy: &assignedBy,
		}

		mocks.AssignmentRepo.EXPECT().CheckUserRoleExists(ctx, userID, roleID).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, errors.New("db error"))

		result, err := logic.AssignRoleToUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("创建角色分配时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.AssignRoleToUserRequest{
			UserID:     &userID,
			RoleID:     &roleID,
			AssignedBy: &assignedBy,
		}

		mocks.AssignmentRepo.EXPECT().CheckUserRoleExists(ctx, userID, roleID).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		result, err := logic.AssignRoleToUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// UpdateUserRoleAssignment 测试
// ============================================================================

func TestLogicImpl_UpdateUserRoleAssignment(t *testing.T) {
	assignmentID := uuid.New().String()
	userID := uuid.New().String()
	roleID := uuid.New().String()
	updatedBy := uuid.New().String()

	t.Run("成功更新角色分配", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingAssignment := &models.UserRoleAssignment{
			BaseModel: models.BaseModel{ID: uuid.MustParse(assignmentID)},
			UserID:    uuid.New(),
			RoleID:    uuid.New(),
		}

		req := &identity_srv.UpdateUserRoleAssignmentRequest{
			AssignmentID: &assignmentID,
			UserID:       &userID,
			RoleID:       &roleID,
			UpdatedBy:    &updatedBy,
		}

		mocks.AssignmentRepo.EXPECT().GetByID(ctx, assignmentID).Return(existingAssignment, nil)
		mocks.AssignmentRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		err := logic.UpdateUserRoleAssignment(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("仅更新部分字段", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingAssignment := &models.UserRoleAssignment{
			BaseModel: models.BaseModel{ID: uuid.MustParse(assignmentID)},
			UserID:    uuid.New(),
			RoleID:    uuid.New(),
		}

		req := &identity_srv.UpdateUserRoleAssignmentRequest{
			AssignmentID: &assignmentID,
			RoleID:       &roleID,
		}

		mocks.AssignmentRepo.EXPECT().GetByID(ctx, assignmentID).Return(existingAssignment, nil)
		mocks.AssignmentRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		err := logic.UpdateUserRoleAssignment(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("分配ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateUserRoleAssignmentRequest{}

		err := logic.UpdateUserRoleAssignment(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("分配ID为空字符串", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyID := ""
		req := &identity_srv.UpdateUserRoleAssignmentRequest{
			AssignmentID: &emptyID,
		}

		err := logic.UpdateUserRoleAssignment(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色分配记录不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateUserRoleAssignmentRequest{
			AssignmentID: &assignmentID,
		}

		mocks.AssignmentRepo.EXPECT().GetByID(ctx, assignmentID).Return(nil, nil)

		err := logic.UpdateUserRoleAssignment(ctx, req)

		assertErrCode(t, errno.ErrRoleAssignmentNotFound, err)
	})

	t.Run("查询角色分配时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateUserRoleAssignmentRequest{
			AssignmentID: &assignmentID,
		}

		mocks.AssignmentRepo.EXPECT().GetByID(ctx, assignmentID).Return(nil, errors.New("db error"))

		err := logic.UpdateUserRoleAssignment(ctx, req)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("更新保存时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingAssignment := &models.UserRoleAssignment{
			BaseModel: models.BaseModel{ID: uuid.MustParse(assignmentID)},
			UserID:    uuid.New(),
			RoleID:    uuid.New(),
		}

		req := &identity_srv.UpdateUserRoleAssignmentRequest{
			AssignmentID: &assignmentID,
			UserID:       &userID,
		}

		mocks.AssignmentRepo.EXPECT().GetByID(ctx, assignmentID).Return(existingAssignment, nil)
		mocks.AssignmentRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		err := logic.UpdateUserRoleAssignment(ctx, req)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// RevokeRoleFromUser 测试
// ============================================================================

func TestLogicImpl_RevokeRoleFromUser(t *testing.T) {
	userID := uuid.New().String()
	roleID := uuid.New().String()

	t.Run("成功撤销角色", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		assignmentUUID := uuid.New()
		assignment := &models.UserRoleAssignment{
			BaseModel: models.BaseModel{ID: assignmentUUID},
			UserID:    uuid.MustParse(userID),
			RoleID:    uuid.MustParse(roleID),
		}

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(false, nil)
		mocks.AssignmentRepo.EXPECT().FindByUserAndRole(ctx, userID, roleID).Return(assignment, nil)
		mocks.AssignmentRepo.EXPECT().Delete(ctx, assignmentUUID.String()).Return(nil)

		err := logic.RevokeRoleFromUser(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("系统用户撤销非系统角色成功", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		assignmentUUID := uuid.New()
		assignment := &models.UserRoleAssignment{
			BaseModel: models.BaseModel{ID: assignmentUUID},
			UserID:    uuid.MustParse(userID),
			RoleID:    uuid.MustParse(roleID),
		}

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(true, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "custom_role",
			IsSystemRole: false,
		}, nil)
		mocks.AssignmentRepo.EXPECT().FindByUserAndRole(ctx, userID, roleID).Return(assignment, nil)
		mocks.AssignmentRepo.EXPECT().Delete(ctx, assignmentUUID.String()).Return(nil)

		err := logic.RevokeRoleFromUser(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("系统用户的系统角色不能被撤销", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(true, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "super_admin",
			IsSystemRole: true,
		}, nil)

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrSystemRoleCannotRevoke, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RevokeRoleFromUserRequest{
			RoleID: &roleID,
		}

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("用户ID为空字符串", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyUserID := ""
		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &emptyUserID,
			RoleID: &roleID,
		}

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
		}

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色ID为空字符串", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyRoleID := ""
		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &emptyRoleID,
		}

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("检查用户类型时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(false, errors.New("db error"))

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("系统用户查询角色定义失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(true, nil)
		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, errors.New("db error"))

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrRoleDefinitionNotFound, err)
	})

	t.Run("角色分配记录不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(false, nil)
		mocks.AssignmentRepo.EXPECT().FindByUserAndRole(ctx, userID, roleID).Return(nil, nil)

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrRoleAssignmentNotFound, err)
	})

	t.Run("查询角色分配记录时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(false, nil)
		mocks.AssignmentRepo.EXPECT().FindByUserAndRole(ctx, userID, roleID).Return(nil, errors.New("db error"))

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("删除角色分配记录时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		assignmentUUID := uuid.New()
		assignment := &models.UserRoleAssignment{
			BaseModel: models.BaseModel{ID: assignmentUUID},
			UserID:    uuid.MustParse(userID),
			RoleID:    uuid.MustParse(roleID),
		}

		req := &identity_srv.RevokeRoleFromUserRequest{
			UserID: &userID,
			RoleID: &roleID,
		}

		mocks.UserRepo.EXPECT().IsSystemUser(ctx, userID).Return(false, nil)
		mocks.AssignmentRepo.EXPECT().FindByUserAndRole(ctx, userID, roleID).Return(assignment, nil)
		mocks.AssignmentRepo.EXPECT().Delete(ctx, assignmentUUID.String()).Return(errors.New("db error"))

		err := logic.RevokeRoleFromUser(ctx, req)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// GetLastUserRoleAssignment 测试
// ============================================================================

func TestLogicImpl_GetLastUserRoleAssignment(t *testing.T) {
	userID := uuid.New().String()

	t.Run("成功获取最后角色分配", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		assignment := &models.UserRoleAssignment{
			BaseModel: models.BaseModel{ID: uuid.New()},
			UserID:    uuid.MustParse(userID),
			RoleID:    uuid.New(),
		}

		mocks.AssignmentRepo.EXPECT().GetLastUserRoleAssignment(ctx, userID).Return(assignment, nil)

		result, err := logic.GetLastUserRoleAssignment(ctx, userID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.GetUserID())
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.GetLastUserRoleAssignment(ctx, "")

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色分配记录不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.AssignmentRepo.EXPECT().GetLastUserRoleAssignment(ctx, userID).Return(nil, nil)

		result, err := logic.GetLastUserRoleAssignment(ctx, userID)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleAssignmentNotFound, err)
	})

	t.Run("查询数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.AssignmentRepo.EXPECT().GetLastUserRoleAssignment(ctx, userID).Return(nil, errors.New("db error"))

		result, err := logic.GetLastUserRoleAssignment(ctx, userID)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// ListUserRoleAssignments 测试
// ============================================================================

func TestLogicImpl_ListUserRoleAssignments(t *testing.T) {
	userID := uuid.New().String()
	roleID := uuid.New().String()

	t.Run("成功列出角色分配_按用户ID", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		assignments := []*models.UserRoleAssignment{
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				UserID:    uuid.MustParse(userID),
				RoleID:    uuid.New(),
			},
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				UserID:    uuid.MustParse(userID),
				RoleID:    uuid.New(),
			},
		}
		pageResult := &models.PageResult{Total: 2, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.UserRoleQueryRequest{
			UserID: &userID,
		}

		mocks.AssignmentRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(assignments, pageResult, nil)

		result, err := logic.ListUserRoleAssignments(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Assignments, 2)
		assert.NotNil(t, result.Page)
	})

	t.Run("成功列出角色分配_按角色ID", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		assignments := []*models.UserRoleAssignment{
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				UserID:    uuid.New(),
				RoleID:    uuid.MustParse(roleID),
			},
		}
		pageResult := &models.PageResult{Total: 1, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.UserRoleQueryRequest{
			RoleID: &roleID,
		}

		mocks.AssignmentRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(assignments, pageResult, nil)

		result, err := logic.ListUserRoleAssignments(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Assignments, 1)
	})

	t.Run("空条件查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UserRoleQueryRequest{}

		mocks.AssignmentRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.UserRoleAssignment{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil)

		result, err := logic.ListUserRoleAssignments(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Assignments)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UserRoleQueryRequest{
			UserID: &userID,
		}

		mocks.AssignmentRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(nil, nil, errors.New("db error"))

		result, err := logic.ListUserRoleAssignments(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// GetUsersByRole 测试
// ============================================================================

func TestLogicImpl_GetUsersByRole(t *testing.T) {
	roleID := uuid.New().String()

	t.Run("成功获取角色下用户列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		userIDs := []string{uuid.New().String(), uuid.New().String()}

		req := &identity_srv.GetUsersByRoleRequest{
			RoleID: &roleID,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().GetAllUserIDsByRoleID(ctx, roleID).Return(userIDs, nil)

		result, err := logic.GetUsersByRole(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, &roleID, result.RoleID)
		assert.Len(t, result.UserIDs, 2)
	})

	t.Run("角色下无用户", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUsersByRoleRequest{
			RoleID: &roleID,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().GetAllUserIDsByRoleID(ctx, roleID).Return([]string{}, nil)

		result, err := logic.GetUsersByRole(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.UserIDs)
	})

	t.Run("角色ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUsersByRoleRequest{}

		result, err := logic.GetUsersByRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色ID为空字符串", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyRoleID := ""
		req := &identity_srv.GetUsersByRoleRequest{
			RoleID: &emptyRoleID,
		}

		result, err := logic.GetUsersByRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUsersByRoleRequest{
			RoleID: &roleID,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, nil)

		result, err := logic.GetUsersByRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleDefinitionNotFound, err)
	})

	t.Run("查询角色信息时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUsersByRoleRequest{
			RoleID: &roleID,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, errors.New("db error"))

		result, err := logic.GetUsersByRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("查询角色用户列表时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUsersByRoleRequest{
			RoleID: &roleID,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().GetAllUserIDsByRoleID(ctx, roleID).Return(nil, errors.New("db error"))

		result, err := logic.GetUsersByRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// BatchBindUsersToRole 测试
// ============================================================================

func TestLogicImpl_BatchBindUsersToRole(t *testing.T) {
	roleID := uuid.New().String()
	userIDs := []string{uuid.New().String(), uuid.New().String(), uuid.New().String()}
	operatorID := uuid.New().String()

	t.Run("成功批量绑定用户到角色", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchBindUsersToRoleRequest{
			RoleID:     &roleID,
			UserIDs:    &core.StringListValue{Items: userIDs},
			OperatorID: &operatorID,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().ReplaceRoleUsers(ctx, roleID, userIDs, operatorID).Return(nil)

		result, err := logic.BatchBindUsersToRole(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.GetSuccess())
		assert.Equal(t, int32(3), result.GetSuccessCount())
	})

	t.Run("成功批量绑定_无操作者ID", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchBindUsersToRoleRequest{
			RoleID:  &roleID,
			UserIDs: &core.StringListValue{Items: userIDs},
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().ReplaceRoleUsers(ctx, roleID, userIDs, "").Return(nil)

		result, err := logic.BatchBindUsersToRole(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.GetSuccess())
	})

	t.Run("角色ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchBindUsersToRoleRequest{
			UserIDs: &core.StringListValue{Items: userIDs},
		}

		result, err := logic.BatchBindUsersToRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色ID为空字符串", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyRoleID := ""
		req := &identity_srv.BatchBindUsersToRoleRequest{
			RoleID:  &emptyRoleID,
			UserIDs: &core.StringListValue{Items: userIDs},
		}

		result, err := logic.BatchBindUsersToRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("用户ID列表为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchBindUsersToRoleRequest{
			RoleID: &roleID,
		}

		result, err := logic.BatchBindUsersToRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchBindUsersToRoleRequest{
			RoleID:  &roleID,
			UserIDs: &core.StringListValue{Items: userIDs},
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, nil)

		result, err := logic.BatchBindUsersToRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleDefinitionNotFound, err)
	})

	t.Run("查询角色信息时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchBindUsersToRoleRequest{
			RoleID:  &roleID,
			UserIDs: &core.StringListValue{Items: userIDs},
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, errors.New("db error"))

		result, err := logic.BatchBindUsersToRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("批量替换用户绑定时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchBindUsersToRoleRequest{
			RoleID:     &roleID,
			UserIDs:    &core.StringListValue{Items: userIDs},
			OperatorID: &operatorID,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(&models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "admin",
		}, nil)
		mocks.AssignmentRepo.EXPECT().ReplaceRoleUsers(ctx, roleID, userIDs, operatorID).Return(errors.New("db error"))

		result, err := logic.BatchBindUsersToRole(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// BatchGetUserRoles 测试
// ============================================================================

func TestLogicImpl_BatchGetUserRoles(t *testing.T) {
	userID1 := uuid.New().String()
	userID2 := uuid.New().String()
	roleID1 := uuid.New().String()
	roleID2 := uuid.New().String()

	t.Run("成功批量获取用户角色", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchGetUserRolesRequest{
			UserIDs: []string{userID1, userID2},
		}

		rolesMap := map[string][]string{
			userID1: {roleID1, roleID2},
			userID2: {roleID1},
		}

		mocks.AssignmentRepo.EXPECT().
			GetRolesByUserIDs(ctx, []string{userID1, userID2}).
			Return(rolesMap, nil)

		result, err := logic.BatchGetUserRoles(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.UserRoles, 2)
	})

	t.Run("空用户列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchGetUserRolesRequest{
			UserIDs: []string{},
		}

		mocks.AssignmentRepo.EXPECT().
			GetRolesByUserIDs(ctx, []string{}).
			Return(map[string][]string{}, nil)

		result, err := logic.BatchGetUserRoles(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.UserRoles)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchGetUserRolesRequest{
			UserIDs: []string{userID1},
		}

		mocks.AssignmentRepo.EXPECT().
			GetRolesByUserIDs(ctx, []string{userID1}).
			Return(nil, errors.New("db error"))

		result, err := logic.BatchGetUserRoles(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("部分用户无角色分配", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.BatchGetUserRolesRequest{
			UserIDs: []string{userID1, userID2},
		}

		// 只有 userID1 有角色，userID2 没有
		rolesMap := map[string][]string{
			userID1: {roleID1},
		}

		mocks.AssignmentRepo.EXPECT().
			GetRolesByUserIDs(ctx, []string{userID1, userID2}).
			Return(rolesMap, nil)

		result, err := logic.BatchGetUserRoles(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.UserRoles, 1)
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
