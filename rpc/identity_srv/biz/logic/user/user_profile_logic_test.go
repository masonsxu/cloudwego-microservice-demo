package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	userDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/user"
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
// CreateUser 测试
// ============================================================================

func TestLogicImpl_CreateUser(t *testing.T) {
	username := "testuser"
	password := "password123"
	email := "test@example.com"
	phone := "+1234567890"

	t.Run("成功创建用户_仅用户名和密码", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, nil)
		mocks.UserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(gomock.Any(), gomock.Any()).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.CreateUser(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("成功创建用户_含邮箱和手机号", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
			Email:    &email,
			Phone:    &phone,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, nil)
		mocks.UserRepo.EXPECT().CheckEmailExists(ctx, email).Return(false, nil)
		mocks.UserRepo.EXPECT().CheckPhoneExists(ctx, phone).Return(false, nil)
		mocks.UserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(gomock.Any(), gomock.Any()).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.CreateUser(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("用户名为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Password: &password,
		}

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("密码为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
		}

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(true, nil)

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrUsernameAlreadyExists, err)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
			Email:    &email,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, nil)
		mocks.UserRepo.EXPECT().CheckEmailExists(ctx, email).Return(true, nil)

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrEmailAlreadyExists, err)
	})

	t.Run("手机号已存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
			Phone:    &phone,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, nil)
		mocks.UserRepo.EXPECT().CheckPhoneExists(ctx, phone).Return(true, nil)

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrPhoneAlreadyExists, err)
	})

	t.Run("检查用户名时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, gorm.ErrInvalidDB)

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("检查邮箱时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
			Email:    &email,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, nil)
		mocks.UserRepo.EXPECT().CheckEmailExists(ctx, email).Return(false, gorm.ErrInvalidDB)

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("检查手机号时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
			Phone:    &phone,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, nil)
		mocks.UserRepo.EXPECT().CheckPhoneExists(ctx, phone).Return(false, gorm.ErrInvalidDB)

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("事务中创建失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
		}

		mocks.UserRepo.EXPECT().CheckUsernameExists(ctx, username).Return(false, nil)
		mocks.UserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.CreateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// GetUser 测试
// ============================================================================

func TestLogicImpl_GetUser(t *testing.T) {
	userID := uuid.New().String()

	t.Run("成功获取用户", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
			Email:     "test@example.com",
		}

		req := &identity_srv.GetUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(profile, nil)
		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(gomock.Any(), userID).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.GetUser(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "testuser", result.GetUsername())
	})

	t.Run("用户不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.GetUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrUserNotFound, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserRequest{}

		result, err := logic.GetUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.GetUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// UpdateUser 测试
// ============================================================================

func TestLogicImpl_UpdateUser(t *testing.T) {
	userID := uuid.New().String()
	newEmail := "new@example.com"
	newPhone := "+9876543210"
	newRealName := "张三"

	t.Run("成功更新用户", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingProfile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
			Email:     "old@example.com",
		}

		req := &identity_srv.UpdateUserRequest{
			UserID:   &userID,
			Email:    &newEmail,
			RealName: &newRealName,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(existingProfile, nil)
		mocks.UserRepo.EXPECT().CheckEmailExists(ctx, newEmail, userID).Return(false, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(gomock.Any(), userID).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.UpdateUser(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("用户不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateUserRequest{
			UserID: &userID,
			Email:  &newEmail,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrUserNotFound, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateUserRequest{
			Email: &newEmail,
		}

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("邮箱已被其他用户使用", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingProfile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
		}

		req := &identity_srv.UpdateUserRequest{
			UserID: &userID,
			Email:  &newEmail,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(existingProfile, nil)
		mocks.UserRepo.EXPECT().CheckEmailExists(ctx, newEmail, userID).Return(true, nil)

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrEmailAlreadyExists, err)
	})

	t.Run("手机号已被其他用户使用", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingProfile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
		}

		req := &identity_srv.UpdateUserRequest{
			UserID: &userID,
			Phone:  &newPhone,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(existingProfile, nil)
		mocks.UserRepo.EXPECT().CheckPhoneExists(ctx, newPhone, userID).Return(true, nil)

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrPhoneAlreadyExists, err)
	})

	t.Run("检查邮箱唯一性时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingProfile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
		}

		req := &identity_srv.UpdateUserRequest{
			UserID: &userID,
			Email:  &newEmail,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(existingProfile, nil)
		mocks.UserRepo.EXPECT().CheckEmailExists(ctx, newEmail, userID).Return(false, gorm.ErrInvalidDB)

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("检查手机号唯一性时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingProfile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
		}

		req := &identity_srv.UpdateUserRequest{
			UserID: &userID,
			Phone:  &newPhone,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(existingProfile, nil)
		mocks.UserRepo.EXPECT().CheckPhoneExists(ctx, newPhone, userID).Return(false, gorm.ErrInvalidDB)

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("获取用户档案时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateUserRequest{
			UserID: &userID,
			Email:  &newEmail,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("更新事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingProfile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
		}

		req := &identity_srv.UpdateUserRequest{
			UserID:   &userID,
			RealName: &newRealName,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(existingProfile, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.UpdateUser(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("更新系统用户信息_允许非关键属性", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingProfile := &models.UserProfile{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(userID)},
			Username:     "admin",
			IsSystemUser: true,
		}

		req := &identity_srv.UpdateUserRequest{
			UserID:   &userID,
			RealName: &newRealName,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(existingProfile, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mocks.MembershipRepo.EXPECT().GetPrimaryMembership(gomock.Any(), userID).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.UpdateUser(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// ============================================================================
// DeleteUser 测试
// ============================================================================

func TestLogicImpl_DeleteUser(t *testing.T) {
	userID := uuid.New().String()

	t.Run("成功删除普通用户", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		user := &models.UserProfile{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(userID)},
			Username:     "testuser",
			IsSystemUser: false,
		}

		req := &identity_srv.DeleteUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
		mocks.UserRepo.EXPECT().SoftDelete(gomock.Any(), userID).Return(nil)

		err := logic.DeleteUser(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("用户不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.DeleteUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		err := logic.DeleteUser(ctx, req)

		assertErrCode(t, errno.ErrUserNotFound, err)
	})

	t.Run("系统用户无法删除", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		user := &models.UserProfile{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(userID)},
			Username:     "admin",
			IsSystemUser: true,
		}

		req := &identity_srv.DeleteUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)

		err := logic.DeleteUser(ctx, req)

		assertErrCode(t, errno.ErrSystemUserCannotDelete, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.DeleteUserRequest{}

		err := logic.DeleteUser(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("获取用户信息时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.DeleteUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrInvalidDB)

		err := logic.DeleteUser(ctx, req)

		assert.Error(t, err)
	})

	t.Run("软删除事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		user := &models.UserProfile{
			BaseModel:    models.BaseModel{ID: uuid.MustParse(userID)},
			Username:     "testuser",
			IsSystemUser: false,
		}

		req := &identity_srv.DeleteUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(user, nil)
		mocks.UserRepo.EXPECT().SoftDelete(gomock.Any(), userID).Return(gorm.ErrInvalidDB)

		err := logic.DeleteUser(ctx, req)

		assert.Error(t, err)
	})
}

// ============================================================================
// ListUsers 测试
// ============================================================================

func TestLogicImpl_ListUsers(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("成功查询用户列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profiles := []*models.UserProfile{
			{BaseModel: models.BaseModel{ID: uuid.New()}, Username: "user1"},
			{BaseModel: models.BaseModel{ID: uuid.New()}, Username: "user2"},
		}
		pageResult := &models.PageResult{Total: 2, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.ListUsersRequest{}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(profiles, pageResult, nil)
		mocks.MembershipRepo.EXPECT().
			GetPrimaryMembershipsByUserIDs(gomock.Any(), gomock.Any()).
			Return(map[string]*models.UserMembership{}, nil)

		result, err := logic.ListUsers(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Users, 2)
		assert.NotNil(t, result.Page)
	})

	t.Run("带状态过滤查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		activeStatus := core.UserStatus_ACTIVE
		req := &identity_srv.ListUsersRequest{
			Status: &activeStatus,
		}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *userDAL.UserProfileQueryConditions,
			) ([]*models.UserProfile, *models.PageResult, error) {
				require.NotNil(t, cond.Status)
				assert.Equal(t, models.UserStatus(core.UserStatus_ACTIVE), *cond.Status)

				return []*models.UserProfile{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil
			})
		// 空列表不会触发批量关联查询

		result, err := logic.ListUsers(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Users)
	})

	t.Run("带组织ID过滤查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.ListUsersRequest{
			OrganizationID: &orgID,
		}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *userDAL.UserProfileQueryConditions,
			) ([]*models.UserProfile, *models.PageResult, error) {
				require.NotNil(t, cond.OrgID)
				assert.Equal(t, orgID, *cond.OrgID)

				return []*models.UserProfile{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil
			})
		// 空列表不会触发批量关联查询

		result, err := logic.ListUsers(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Users)
	})

	t.Run("带分页参数查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.ListUsersRequest{
			Page: &rpc_base.PageRequest{Page: 2, Limit: 10},
		}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.UserProfile{}, &models.PageResult{Total: 15, Page: 2, Limit: 10, TotalPages: 2}, nil)
		// 空列表不会触发批量关联查询

		result, err := logic.ListUsers(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result.Page)
	})

	t.Run("空请求参数", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.UserProfile{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil)
		// 空列表不会触发批量关联查询

		result, err := logic.ListUsers(ctx, nil)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.ListUsersRequest{}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(nil, nil, gorm.ErrInvalidDB)

		result, err := logic.ListUsers(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// SearchUsers 测试
// ============================================================================

func TestLogicImpl_SearchUsers(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("成功搜索用户", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profiles := []*models.UserProfile{
			{BaseModel: models.BaseModel{ID: uuid.New()}, Username: "user1"},
		}
		pageResult := &models.PageResult{Total: 1, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.SearchUsersRequest{}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(profiles, pageResult, nil)
		mocks.MembershipRepo.EXPECT().
			GetPrimaryMembershipsByUserIDs(gomock.Any(), gomock.Any()).
			Return(map[string]*models.UserMembership{}, nil)

		result, err := logic.SearchUsers(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Users, 1)
		assert.NotNil(t, result.Page)
	})

	t.Run("带组织ID搜索", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.SearchUsersRequest{
			OrganizationID: &orgID,
		}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				cond *userDAL.UserProfileQueryConditions,
			) ([]*models.UserProfile, *models.PageResult, error) {
				require.NotNil(t, cond.OrgID)
				assert.Equal(t, orgID, *cond.OrgID)

				return []*models.UserProfile{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil
			})
		// 空列表不会触发批量关联查询

		result, err := logic.SearchUsers(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Users)
	})

	t.Run("请求参数为nil", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		result, err := logic.SearchUsers(ctx, nil)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("带分页参数搜索", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.SearchUsersRequest{
			Page: &rpc_base.PageRequest{Page: 1, Limit: 5},
		}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return([]*models.UserProfile{}, &models.PageResult{Total: 0, Page: 1, Limit: 5, TotalPages: 1}, nil)
		// 空列表不会触发批量关联查询

		result, err := logic.SearchUsers(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result.Page)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.SearchUsersRequest{}

		mocks.UserRepo.EXPECT().
			FindWithConditions(ctx, gomock.Any()).
			Return(nil, nil, gorm.ErrInvalidDB)

		result, err := logic.SearchUsers(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// ChangeUserStatus 测试
// ============================================================================

func TestLogicImpl_ChangeUserStatus(t *testing.T) {
	userID := uuid.New().String()

	t.Run("成功更改用户状态为停用", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
			Status:    models.UserStatusActive,
		}

		newStatus := core.UserStatus_SUSPENDED
		req := &identity_srv.ChangeUserStatusRequest{
			UserID:     &userID,
			NewStatus_: &newStatus,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(profile, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, p *models.UserProfile) error {
				assert.Equal(t, models.UserStatusSuspended, p.Status)
				return nil
			})

		err := logic.ChangeUserStatus(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("成功激活用户", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
			Status:    models.UserStatusInactive,
		}

		newStatus := core.UserStatus_ACTIVE
		req := &identity_srv.ChangeUserStatusRequest{
			UserID:     &userID,
			NewStatus_: &newStatus,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(profile, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, p *models.UserProfile) error {
				assert.Equal(t, models.UserStatusActive, p.Status)
				return nil
			})

		err := logic.ChangeUserStatus(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("用户不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		newStatus := core.UserStatus_ACTIVE
		req := &identity_srv.ChangeUserStatusRequest{
			UserID:     &userID,
			NewStatus_: &newStatus,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		err := logic.ChangeUserStatus(ctx, req)

		assertErrCode(t, errno.ErrUserNotFound, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyID := ""
		newStatus := core.UserStatus_ACTIVE
		req := &identity_srv.ChangeUserStatusRequest{
			UserID:     &emptyID,
			NewStatus_: &newStatus,
		}

		err := logic.ChangeUserStatus(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		newStatus := core.UserStatus_ACTIVE
		req := &identity_srv.ChangeUserStatusRequest{
			UserID:     &userID,
			NewStatus_: &newStatus,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrInvalidDB)

		err := logic.ChangeUserStatus(ctx, req)

		assert.Error(t, err)
	})

	t.Run("更新事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
			Status:    models.UserStatusActive,
		}

		newStatus := core.UserStatus_SUSPENDED
		req := &identity_srv.ChangeUserStatusRequest{
			UserID:     &userID,
			NewStatus_: &newStatus,
		}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(profile, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		err := logic.ChangeUserStatus(ctx, req)

		assert.Error(t, err)
	})
}

// ============================================================================
// UnlockUser 测试
// ============================================================================

func TestLogicImpl_UnlockUser(t *testing.T) {
	userID := uuid.New().String()

	t.Run("成功解锁用户", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
			Status:    models.UserStatusLocked,
		}

		req := &identity_srv.UnlockUserRequest{UserID: &userID}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(profile, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, p *models.UserProfile) error {
				assert.Equal(t, models.UserStatusActive, p.Status)
				return nil
			})

		err := logic.UnlockUser(ctx, req)

		assert.NoError(t, err)
	})

	t.Run("用户不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UnlockUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrRecordNotFound)

		err := logic.UnlockUser(ctx, req)

		assertErrCode(t, errno.ErrUserNotFound, err)
	})

	t.Run("用户ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		emptyID := ""
		req := &identity_srv.UnlockUserRequest{UserID: &emptyID}

		err := logic.UnlockUser(ctx, req)

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UnlockUserRequest{UserID: &userID}
		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(nil, gorm.ErrInvalidDB)

		err := logic.UnlockUser(ctx, req)

		assert.Error(t, err)
	})

	t.Run("更新事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		profile := &models.UserProfile{
			BaseModel: models.BaseModel{ID: uuid.MustParse(userID)},
			Username:  "testuser",
			Status:    models.UserStatusLocked,
		}

		req := &identity_srv.UnlockUserRequest{UserID: &userID}

		mocks.UserRepo.EXPECT().GetByID(ctx, userID).Return(profile, nil)
		mocks.UserRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		err := logic.UnlockUser(ctx, req)

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
