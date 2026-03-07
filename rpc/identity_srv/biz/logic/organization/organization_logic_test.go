package organization

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
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
		dal:               mocks.DAL,
		converter:         mocks.Converter,
		logoStorageClient: nil,
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
// CreateOrganization 测试
// ============================================================================

func TestLogicImpl_CreateOrganization(t *testing.T) {
	orgName := "测试组织"

	t.Run("成功创建组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateOrganizationRequest{
			Name: &orgName,
		}

		mocks.OrgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mocks.LogoRepo.EXPECT().GetByOrganizationID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.CreateOrganization(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orgName, result.GetName())
	})

	t.Run("成功创建带完整属性的组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		facilityType := "hospital"
		accreditationStatus := "accredited"

		req := &identity_srv.CreateOrganizationRequest{
			Name:                &orgName,
			FacilityType:        &facilityType,
			AccreditationStatus: &accreditationStatus,
			ProvinceCity:        []string{"北京", "海淀区"},
		}

		mocks.OrgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mocks.LogoRepo.EXPECT().GetByOrganizationID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.CreateOrganization(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orgName, result.GetName())
	})

	t.Run("缺少组织名称", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateOrganizationRequest{}

		result, err := logic.CreateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("事务中创建失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateOrganizationRequest{
			Name: &orgName,
		}

		mocks.OrgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.CreateOrganization(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("省市信息超过10个", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		provinceCity := make([]string, 11)
		for i := range provinceCity {
			provinceCity[i] = "城市"
		}

		req := &identity_srv.CreateOrganizationRequest{
			Name:         &orgName,
			ProvinceCity: provinceCity,
		}

		result, err := logic.CreateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("省市信息包含空值", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.CreateOrganizationRequest{
			Name:         &orgName,
			ProvinceCity: []string{"北京", ""},
		}

		result, err := logic.CreateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("成功创建带LogoID的组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		logoID := uuid.New().String()
		logo := &models.OrganizationLogo{
			BaseModel: models.BaseModel{ID: uuid.MustParse(logoID)},
			Status:    models.LogoStatusTemporary,
			FileID:    "test-file-id",
		}

		req := &identity_srv.CreateOrganizationRequest{
			Name:   &orgName,
			LogoID: &logoID,
		}

		mocks.OrgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mocks.LogoRepo.EXPECT().GetByID(gomock.Any(), logoID).Return(logo, nil)
		mocks.LogoRepo.EXPECT().BindToOrganization(gomock.Any(), uuid.MustParse(logoID), gomock.Any()).Return(nil)
		mocks.LogoRepo.EXPECT().GetByOrganizationID(gomock.Any(), gomock.Any()).Return(logo, nil)

		result, err := logic.CreateOrganization(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("创建带无效LogoID格式的组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		invalidLogoID := "invalid-uuid"

		req := &identity_srv.CreateOrganizationRequest{
			Name:   &orgName,
			LogoID: &invalidLogoID,
		}

		mocks.OrgRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		result, err := logic.CreateOrganization(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

// ============================================================================
// GetOrganization 测试
// ============================================================================

func TestLogicImpl_GetOrganization(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("成功获取组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		org := &models.Organization{
			BaseModel: models.BaseModel{ID: uuid.MustParse(orgID)},
			Name:      "测试组织",
		}

		req := &identity_srv.GetOrganizationRequest{OrganizationID: &orgID}
		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(org, nil)
		mocks.LogoRepo.EXPECT().
			GetByOrganizationID(gomock.Any(), uuid.MustParse(orgID)).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.GetOrganization(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试组织", result.GetName())
	})

	t.Run("组织不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetOrganizationRequest{OrganizationID: &orgID}
		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.GetOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOrganizationNotFound, err)
	})

	t.Run("组织ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetOrganizationRequest{}

		result, err := logic.GetOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.GetOrganizationRequest{OrganizationID: &orgID}
		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.GetOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("获取带Logo的组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		org := &models.Organization{
			BaseModel: models.BaseModel{ID: uuid.MustParse(orgID)},
			Name:      "测试组织",
		}

		logoID := uuid.New()
		logo := &models.OrganizationLogo{
			BaseModel: models.BaseModel{ID: logoID},
			Status:    models.LogoStatusBound,
			FileID:    "test-file-id",
		}

		req := &identity_srv.GetOrganizationRequest{OrganizationID: &orgID}
		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(org, nil)
		mocks.LogoRepo.EXPECT().GetByOrganizationID(gomock.Any(), uuid.MustParse(orgID)).Return(logo, nil)

		result, err := logic.GetOrganization(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
		// logoStorageClient 为 nil 时，generateLogoURL 返回 fileID 本身
		assert.NotNil(t, result.Logo)
		assert.Equal(t, "test-file-id", *result.Logo)
		assert.NotNil(t, result.LogoID)
		assert.Equal(t, logoID.String(), *result.LogoID)
	})
}

// ============================================================================
// UpdateOrganization 测试
// ============================================================================

func TestLogicImpl_UpdateOrganization(t *testing.T) {
	orgID := uuid.New().String()
	newName := "更新后的组织"

	t.Run("成功更新组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingOrg := &models.Organization{
			BaseModel: models.BaseModel{ID: uuid.MustParse(orgID)},
			Name:      "原组织名称",
		}

		req := &identity_srv.UpdateOrganizationRequest{
			OrganizationID: &orgID,
			Name:           &newName,
		}

		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(existingOrg, nil)
		mocks.OrgRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mocks.LogoRepo.EXPECT().
			GetByOrganizationID(gomock.Any(), uuid.MustParse(orgID)).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.UpdateOrganization(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("组织不存在", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateOrganizationRequest{
			OrganizationID: &orgID,
			Name:           &newName,
		}

		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(nil, gorm.ErrRecordNotFound)

		result, err := logic.UpdateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOrganizationNotFound, err)
	})

	t.Run("组织ID为空", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateOrganizationRequest{
			Name: &newName,
		}

		result, err := logic.UpdateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("更新事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		existingOrg := &models.Organization{
			BaseModel: models.BaseModel{ID: uuid.MustParse(orgID)},
			Name:      "原组织名称",
		}

		req := &identity_srv.UpdateOrganizationRequest{
			OrganizationID: &orgID,
			Name:           &newName,
		}

		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(existingOrg, nil)
		mocks.OrgRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

		result, err := logic.UpdateOrganization(ctx, req)

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("获取组织时数据库错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateOrganizationRequest{
			OrganizationID: &orgID,
			Name:           &newName,
		}

		mocks.OrgRepo.EXPECT().GetByID(ctx, orgID).Return(nil, gorm.ErrInvalidDB)

		result, err := logic.UpdateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})

	t.Run("更新省市信息包含空值", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.UpdateOrganizationRequest{
			OrganizationID: &orgID,
			ProvinceCity:   []string{"北京", ""},
		}

		result, err := logic.UpdateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("更新省市信息超过10个", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		provinceCity := make([]string, 11)
		for i := range provinceCity {
			provinceCity[i] = "城市"
		}

		req := &identity_srv.UpdateOrganizationRequest{
			OrganizationID: &orgID,
			ProvinceCity:   provinceCity,
		}

		result, err := logic.UpdateOrganization(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrInvalidParams, err)
	})
}

// ============================================================================
// DeleteOrganization 测试
// ============================================================================

func TestLogicImpl_DeleteOrganization(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("成功删除组织", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.OrgRepo.EXPECT().SoftDelete(gomock.Any(), orgID).Return(nil)

		err := logic.DeleteOrganization(ctx, orgID)

		assert.NoError(t, err)
	})

	t.Run("空组织ID", func(t *testing.T) {
		logic, _ := setupTest(t)
		ctx := context.Background()

		err := logic.DeleteOrganization(ctx, "")

		assertErrCode(t, errno.ErrInvalidParams, err)
	})

	t.Run("删除事务失败", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		mocks.OrgRepo.EXPECT().SoftDelete(gomock.Any(), orgID).Return(gorm.ErrInvalidDB)

		err := logic.DeleteOrganization(ctx, orgID)

		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// ListOrganizations 测试
// ============================================================================

func TestLogicImpl_ListOrganizations(t *testing.T) {
	t.Run("成功查询组织列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		organizations := []*models.Organization{
			{BaseModel: models.BaseModel{ID: uuid.New()}, Name: "组织A"},
			{BaseModel: models.BaseModel{ID: uuid.New()}, Name: "组织B"},
		}
		pageResult := &models.PageResult{Total: 2, Page: 1, Limit: 20, TotalPages: 1}

		req := &identity_srv.ListOrganizationsRequest{}

		mocks.OrgRepo.EXPECT().
			FindAll(ctx, gomock.Any()).
			Return(organizations, pageResult, nil)
		// 每个组织都会查询 Logo
		mocks.LogoRepo.EXPECT().GetByOrganizationID(gomock.Any(), gomock.Any()).
			Return(nil, gorm.ErrRecordNotFound).Times(2)

		result, err := logic.ListOrganizations(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Organizations, 2)
		assert.NotNil(t, result.Page)
	})

	t.Run("空列表", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.ListOrganizationsRequest{}

		mocks.OrgRepo.EXPECT().
			FindAll(ctx, gomock.Any()).
			Return([]*models.Organization{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil)

		result, err := logic.ListOrganizations(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Organizations)
	})

	t.Run("带分页参数查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.ListOrganizationsRequest{
			Page: &rpc_base.PageRequest{Page: 2, Limit: 10},
		}

		mocks.OrgRepo.EXPECT().
			FindAll(ctx, gomock.Any()).
			Return([]*models.Organization{}, &models.PageResult{Total: 15, Page: 2, Limit: 10, TotalPages: 2}, nil)

		result, err := logic.ListOrganizations(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, result.Page)
	})

	t.Run("按ParentID过滤查询", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		parentID := uuid.New().String()
		req := &identity_srv.ListOrganizationsRequest{
			ParentID: &parentID,
		}

		mocks.OrgRepo.EXPECT().
			FindAll(ctx, gomock.Any()).
			DoAndReturn(func(
				_ context.Context,
				opts *base.QueryOptions,
			) ([]*models.Organization, *models.PageResult, error) {
				// 验证 ParentID 过滤条件已添加
				assert.NotNil(t, opts.Filters["parent_id"])
				return []*models.Organization{}, &models.PageResult{Total: 0, Page: 1, Limit: 20, TotalPages: 1}, nil
			})

		result, err := logic.ListOrganizations(ctx, req)

		require.NoError(t, err)
		assert.Empty(t, result.Organizations)
	})

	t.Run("数据库查询错误", func(t *testing.T) {
		logic, mocks := setupTest(t)
		ctx := context.Background()

		req := &identity_srv.ListOrganizationsRequest{}

		mocks.OrgRepo.EXPECT().
			FindAll(ctx, gomock.Any()).
			Return(nil, nil, gorm.ErrInvalidDB)

		result, err := logic.ListOrganizations(ctx, req)

		assert.Nil(t, result)
		assertErrCode(t, errno.ErrOperationFailed, err)
	})
}

// ============================================================================
// NewLogic 构造函数测试
// ============================================================================

func TestNewLogic(t *testing.T) {
	ctrl := gomock.NewController(t)
	mocks := mock.NewTestMocks(ctrl)

	logic := NewLogic(mocks.DAL, converter.NewConverter(), nil)

	assert.NotNil(t, logic)
}
