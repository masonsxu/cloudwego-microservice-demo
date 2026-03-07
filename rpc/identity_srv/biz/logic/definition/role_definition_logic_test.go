package definition

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	definitionDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/definition"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/mock"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
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
// CreateRoleDefinition 测试
// ============================================================================

func TestLogicImpl_CreateRoleDefinition(t *testing.T) {
	roleName := "管理员"
	roleDesc := "系统管理员角色"
	resource := "user"
	action := "read"
	permDesc := "读取用户"

	t.Run("成功创建角色定义", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionCreateRequest{
			Name:        &roleName,
			Description: &roleDesc,
			Permissions: []*identity_srv.Permission{
				{Resource: &resource, Action: &action, Description: &permDesc},
			},
			IsSystemRole: false,
		}

		mocks.DefinitionRepo.EXPECT().CheckNameExists(ctx, roleName).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.CreateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, roleName, result.GetName())
	})

	t.Run("角色名称已存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionCreateRequest{
			Name:        &roleName,
			Description: &roleDesc,
		}

		mocks.DefinitionRepo.EXPECT().CheckNameExists(ctx, roleName).Return(true, nil)

		result, err := logic.CreateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleNameAlreadyExists, err)
	})

	t.Run("检查名称存在时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionCreateRequest{
			Name:        &roleName,
			Description: &roleDesc,
		}

		mocks.DefinitionRepo.EXPECT().CheckNameExists(ctx, roleName).Return(false, gorm.ErrInvalidDB)

		result, err := logic.CreateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("创建记录时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionCreateRequest{
			Name:        &roleName,
			Description: &roleDesc,
		}

		mocks.DefinitionRepo.EXPECT().CheckNameExists(ctx, roleName).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.CreateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("权限列表中Description为nil时使用空字符串", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionCreateRequest{
			Name:        &roleName,
			Description: &roleDesc,
			Permissions: []*identity_srv.Permission{
				{Resource: &resource, Action: &action, Description: nil},
			},
		}

		mocks.DefinitionRepo.EXPECT().CheckNameExists(ctx, roleName).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, role *models.RoleDefinition) error {
				require.Len(t, role.Permissions, 1)
				assert.Equal(t, "", role.Permissions[0].Description)

				return nil
			})

		result, err := logic.CreateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("创建系统角色", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionCreateRequest{
			Name:         &roleName,
			Description:  &roleDesc,
			IsSystemRole: true,
		}

		mocks.DefinitionRepo.EXPECT().CheckNameExists(ctx, roleName).Return(false, nil)
		mocks.DefinitionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, role *models.RoleDefinition) error {
				assert.True(t, role.IsSystemRole)
				assert.Equal(t, models.RoleStatusInactive, role.Status)

				return nil
			})

		result, err := logic.CreateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// ============================================================================
// UpdateRoleDefinition 测试
// ============================================================================

func TestLogicImpl_UpdateRoleDefinition(t *testing.T) {
	roleID := uuid.New().String()
	newName := "超级管理员"
	newDesc := "更新后的描述"
	activeStatus := core.RoleStatus(models.RoleStatusActive)

	t.Run("成功更新角色定义", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			Description:  "旧描述",
			IsSystemRole: false,
		}

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Description:      &newDesc,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("更新角色名称成功_无冲突", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			IsSystemRole: false,
		}

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Name:             &newName,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().FindByName(ctx, newName).Return(nil, nil)
		mocks.DefinitionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("更新角色名称_名称与自身相同", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			IsSystemRole: false,
		}

		// FindByName 返回的角色就是自身，ID相同，不应报错
		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Name:             &newName,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().FindByName(ctx, newName).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("更新角色名称_名称与其他角色冲突", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			IsSystemRole: false,
		}

		otherRoleID := uuid.New()
		otherRole := &models.RoleDefinition{
			BaseModel: models.BaseModel{ID: otherRoleID},
			Name:      newName,
		}

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Name:             &newName,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().FindByName(ctx, newName).Return(otherRole, nil)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("更新角色名称_检查名称时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			IsSystemRole: false,
		}

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Name:             &newName,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().FindByName(ctx, newName).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("更新角色状态", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			Status:       models.RoleStatusInactive,
			IsSystemRole: false,
		}

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Status:           &activeStatus,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, role *models.RoleDefinition) error {
				assert.Equal(t, models.RoleStatusActive, role.Status)
				return nil
			})

		result, err := logic.UpdateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("更新角色权限列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			IsSystemRole: false,
		}

		newResource := "organization"
		newAction := "write"

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Permissions: []*identity_srv.Permission{
				{Resource: &newResource, Action: &newAction},
			},
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, role *models.RoleDefinition) error {
				require.Len(t, role.Permissions, 1)
				assert.Equal(t, "organization", role.Permissions[0].Resource)
				assert.Equal(t, "write", role.Permissions[0].Action)

				return nil
			})

		result, err := logic.UpdateRoleDefinition(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("角色定义ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionUpdateRequest{
			Description: &newDesc,
		}

		result, err := logic.UpdateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色定义不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Description:      &newDesc,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, nil)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleDefinitionNotFound, err)
	})

	t.Run("查询角色定义时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Description:      &newDesc,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("系统角色不允许修改", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		systemRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "系统管理员",
			IsSystemRole: true,
		}

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Description:      &newDesc,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(systemRole, nil)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrSystemRoleCannotModify, err)
	})

	t.Run("更新保存时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "管理员",
			IsSystemRole: false,
		}

		req := &identity_srv.RoleDefinitionUpdateRequest{
			RoleDefinitionID: &roleID,
			Description:      &newDesc,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.DefinitionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.UpdateRoleDefinition(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// DeleteRoleDefinition 测试
// ============================================================================

func TestLogicImpl_DeleteRoleDefinition(t *testing.T) {
	roleID := uuid.New().String()

	t.Run("成功删除角色定义", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "测试角色",
			IsSystemRole: false,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID).Return(int64(0), nil)
		mocks.DefinitionRepo.EXPECT().Delete(ctx, roleID).Return(nil)

		err := logic.DeleteRoleDefinition(ctx, roleID)

		assert.NoError(t, err)
	})

	t.Run("角色ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		err := logic.DeleteRoleDefinition(ctx, "")

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色定义不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, nil)

		err := logic.DeleteRoleDefinition(ctx, roleID)

		assertErrCode(t, errno.ErrRoleDefinitionNotFound, err)
	})

	t.Run("查询角色定义时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, gorm.ErrInvalidDB)

		err := logic.DeleteRoleDefinition(ctx, roleID)

		assert.Error(t, err)
	})

	t.Run("系统角色不允许删除", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		systemRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "系统管理员",
			IsSystemRole: true,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(systemRole, nil)

		err := logic.DeleteRoleDefinition(ctx, roleID)

		assertErrCode(t, errno.ErrSystemRoleCannotDelete, err)
	})

	t.Run("角色正在使用中无法删除", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "测试角色",
			IsSystemRole: false,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID).Return(int64(3), nil)

		err := logic.DeleteRoleDefinition(ctx, roleID)

		assertErrCode(t, errno.ErrRoleInUseCannotDelete, err)
	})

	t.Run("检查角色使用情况时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "测试角色",
			IsSystemRole: false,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID).Return(int64(0), gorm.ErrInvalidDB)

		err := logic.DeleteRoleDefinition(ctx, roleID)

		assert.Error(t, err)
	})

	t.Run("删除操作时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:         "测试角色",
			IsSystemRole: false,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID).Return(int64(0), nil)
		mocks.DefinitionRepo.EXPECT().Delete(ctx, roleID).Return(gorm.ErrInvalidDB)

		err := logic.DeleteRoleDefinition(ctx, roleID)

		assert.Error(t, err)
	})
}

// ============================================================================
// GetRoleDefinition 测试
// ============================================================================

func TestLogicImpl_GetRoleDefinition(t *testing.T) {
	roleID := uuid.New().String()

	t.Run("成功获取角色定义", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel:   models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:        "管理员",
			Description: "管理员角色",
			Status:      models.RoleStatusActive,
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID).Return(int64(5), nil)

		result, err := logic.GetRoleDefinition(ctx, roleID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "管理员", result.GetName())
		assert.Equal(t, int64(5), result.GetUserCount())
	})

	t.Run("角色ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.GetRoleDefinition(ctx, "")

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("角色定义不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, nil)

		result, err := logic.GetRoleDefinition(ctx, roleID)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrRoleDefinitionNotFound, err)
	})

	t.Run("查询角色定义时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.GetRoleDefinition(ctx, roleID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("查询角色用户数量时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingRole := &models.RoleDefinition{
			BaseModel: models.BaseModel{ID: uuid.MustParse(roleID)},
			Name:      "管理员",
		}

		mocks.DefinitionRepo.EXPECT().GetByID(ctx, roleID).Return(existingRole, nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID).Return(int64(0), gorm.ErrInvalidDB)

		result, err := logic.GetRoleDefinition(ctx, roleID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// ListRoleDefinitions 测试
// ============================================================================

func TestLogicImpl_ListRoleDefinitions(t *testing.T) {
	t.Run("成功列出角色定义列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		roleID1 := uuid.New()
		roleID2 := uuid.New()

		roles := []*models.RoleDefinition{
			{BaseModel: models.BaseModel{ID: roleID1}, Name: "管理员"},
			{BaseModel: models.BaseModel{ID: roleID2}, Name: "医生"},
		}
		pageResult := &models.PageResult{Total: 2, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.RoleDefinitionQueryRequest{}

		mocks.DefinitionRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(roles, pageResult, nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID1.String()).Return(int64(3), nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID2.String()).Return(int64(1), nil)

		result, err := logic.ListRoleDefinitions(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Roles, 2)
		assert.NotNil(t, result.Page)
	})

	t.Run("带过滤条件查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		roleName := "管理员"
		activeStatus := core.RoleStatus(models.RoleStatusActive)
		isSystem := true

		req := &identity_srv.RoleDefinitionQueryRequest{
			Name:         &roleName,
			Status:       &activeStatus,
			IsSystemRole: &isSystem,
		}

		mocks.DefinitionRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *definitionDAL.RoleDefinitionQueryConditions,
			) ([]*models.RoleDefinition, *models.PageResult, error) {
				assert.NotNil(t, cond.Name)
				assert.Equal(t, roleName, *cond.Name)
				assert.NotNil(t, cond.Status)
				assert.Equal(t, models.RoleStatusActive, *cond.Status)
				assert.NotNil(t, cond.IsSystemRole)
				assert.True(t, *cond.IsSystemRole)

				emptyResult := &models.PageResult{
					Total: 0, Page: 1, Limit: 20, TotalPages: 1,
				}

				return []*models.RoleDefinition{}, emptyResult, nil
			})

		result, err := logic.ListRoleDefinitions(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Roles)
	})

	t.Run("带分页参数查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionQueryRequest{
			Page: &rpc_base.PageRequest{Page: 2, Limit: 10},
		}

		mocks.DefinitionRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.RoleDefinition{}, &models.PageResult{Total: 15, Page: 2, Limit: 10, TotalPages: 2}, nil)

		result, err := logic.ListRoleDefinitions(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result.Page)
	})

	t.Run("空结果列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionQueryRequest{}

		mocks.DefinitionRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.RoleDefinition{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil)

		result, err := logic.ListRoleDefinitions(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Roles)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.RoleDefinitionQueryRequest{}

		mocks.DefinitionRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(nil, nil, gorm.ErrInvalidDB)

		result, err := logic.ListRoleDefinitions(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("请求参数为nil", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.DefinitionRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.RoleDefinition{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil)

		result, err := logic.ListRoleDefinitions(ctx, nil)

		require.NoError(t, err)
		assert.Empty(t, result.Roles)
	})

	t.Run("查询用户数量部分失败时继续处理", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		roleID1 := uuid.New()
		roleID2 := uuid.New()

		roles := []*models.RoleDefinition{
			{BaseModel: models.BaseModel{ID: roleID1}, Name: "管理员"},
			{BaseModel: models.BaseModel{ID: roleID2}, Name: "医生"},
		}
		pageResult := &models.PageResult{Total: 2, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.RoleDefinitionQueryRequest{}

		mocks.DefinitionRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(roles, pageResult, nil)
		// 第一个角色查询成功，第二个失败
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID1.String()).Return(int64(3), nil)
		mocks.AssignmentRepo.EXPECT().CountByRoleID(ctx, roleID2.String()).Return(int64(0), gorm.ErrInvalidDB)

		result, err := logic.ListRoleDefinitions(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Roles, 2)
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
