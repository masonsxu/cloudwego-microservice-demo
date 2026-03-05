package membership

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// mockEnumConverter 是枚举转换器的模拟实现
type mockEnumConverter struct{}

func (m *mockEnumConverter) ModelUserStatusToThrift(status models.UserStatus) core.UserStatus {
	return core.UserStatus_ACTIVE
}

func (m *mockEnumConverter) ThriftUserStatusToModel(status core.UserStatus) models.UserStatus {
	return models.UserStatusActive
}

func (m *mockEnumConverter) ModelRoleStatusToThrift(status models.RoleStatus) core.RoleStatus {
	return core.RoleStatus_ACTIVE
}

func (m *mockEnumConverter) ThriftRoleStatusToModel(status core.RoleStatus) models.RoleStatus {
	return models.RoleStatusActive
}

func (m *mockEnumConverter) ModelGenderToThrift(gender models.Gender) core.Gender {
	return core.Gender_UNKNOWN
}

func (m *mockEnumConverter) ThriftGenderToModel(gender core.Gender) models.Gender {
	return models.GenderUnknown
}

func (m *mockEnumConverter) ModelDataScopeToThrift(scope models.DataScopeType) identity_srv.DataScope {
	return identity_srv.DataScope_SELF
}

func (m *mockEnumConverter) ThriftDataScopeToModel(scope identity_srv.DataScope) models.DataScopeType {
	return models.DataScopeSelf
}

func TestNewConverter(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	require.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
	assert.Equal(t, mockEnumConv, converter.(*ConverterImpl).converter)
}

func TestConverterImpl_ModelUserMembershipToThrift(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelUserMembershipToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("完整成员关系转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()
		deptID := uuid.New()

		model := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   deptID,
			IsPrimary:      true,
			Status:         models.MembershipStatusActive,
		}

		result := converter.ModelUserMembershipToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, userID.String(), *result.ID)
		assert.Equal(t, userID.String(), *result.UserID)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.Equal(t, deptID.String(), *result.DepartmentID)
		assert.NotNil(t, result.IsPrimary)
		assert.True(t, *result.IsPrimary)
		assert.Equal(t, now, *result.CreatedAt)
		assert.Equal(t, now, *result.UpdatedAt)
	})

	t.Run("无部门ID的成员关系转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()

		model := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   uuid.Nil, // 零值
			IsPrimary:      false,
			Status:         models.MembershipStatusPending,
		}

		result := converter.ModelUserMembershipToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, userID.String(), *result.ID)
		assert.Equal(t, userID.String(), *result.UserID)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.Nil(t, result.DepartmentID) // 零值不设置
		assert.NotNil(t, result.IsPrimary)
		assert.False(t, *result.IsPrimary)
	})

	t.Run("非主要成员关系", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()

		model := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:         userID,
			OrganizationID: orgID,
			IsPrimary:      false,
			Status:         models.MembershipStatusActive,
		}

		result := converter.ModelUserMembershipToThrift(model)

		require.NotNil(t, result)
		assert.NotNil(t, result.IsPrimary)
		assert.False(t, *result.IsPrimary)
	})
}

func TestConverterImpl_ThriftUserMembershipToModel(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ThriftUserMembershipToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整Thrift对象转换", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()
		deptID := uuid.New()

		idStr := userID.String()
		userIDStr := userID.String()
		orgIDStr := orgID.String()
		deptIDStr := deptID.String()
		isPrimary := true

		dto := &identity_srv.UserMembership{
			ID:             &idStr,
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			DepartmentID:   &deptIDStr,
			IsPrimary:      &isPrimary,
		}

		result := converter.ThriftUserMembershipToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, deptID, result.DepartmentID)
		assert.True(t, result.IsPrimary)
	})

	t.Run("无部门ID的转换", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()

		idStr := userID.String()
		userIDStr := userID.String()
		orgIDStr := orgID.String()

		dto := &identity_srv.UserMembership{
			ID:             &idStr,
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			DepartmentID:   nil, // 可选字段为nil
			IsPrimary:      nil, // 可选字段为nil
		}

		result := converter.ThriftUserMembershipToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, uuid.Nil, result.DepartmentID) // nil转换为零值
		assert.False(t, result.IsPrimary)              // nil转换为false
	})

	t.Run("非主要成员关系", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()

		idStr := userID.String()
		userIDStr := userID.String()
		orgIDStr := orgID.String()
		isPrimary := false

		dto := &identity_srv.UserMembership{
			ID:             &idStr,
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			IsPrimary:      &isPrimary,
		}

		result := converter.ThriftUserMembershipToModel(dto)

		require.NotNil(t, result)
		assert.False(t, result.IsPrimary)
	})
}

func TestConverterImpl_ModelUserMembershipsToThrift(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("空切片转换", func(t *testing.T) {
		result := converter.ModelUserMembershipsToThrift([]*models.UserMembership{})
		assert.Nil(t, result)
	})

	t.Run("nil切片转换", func(t *testing.T) {
		result := converter.ModelUserMembershipsToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("多个成员关系批量转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID1 := uuid.New()
		userID2 := uuid.New()
		orgID1 := uuid.New()
		orgID2 := uuid.New()

		models := []*models.UserMembership{
			{
				BaseModel: models.BaseModel{
					ID:        userID1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:         userID1,
				OrganizationID: orgID1,
				IsPrimary:      true,
				Status:         models.MembershipStatusActive,
			},
			{
				BaseModel: models.BaseModel{
					ID:        userID2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:         userID2,
				OrganizationID: orgID2,
				IsPrimary:      false,
				Status:         models.MembershipStatusActive,
			},
		}

		result := converter.ModelUserMembershipsToThrift(models)

		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, userID1.String(), *result[0].ID)
		assert.Equal(t, userID2.String(), *result[1].ID)
	})

	t.Run("批量转换包含nil", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()

		models := []*models.UserMembership{
			{
				BaseModel: models.BaseModel{
					ID:        userID,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:         userID,
				OrganizationID: orgID,
				IsPrimary:      true,
				Status:         models.MembershipStatusActive,
			},
			nil, // 包含nil
		}

		result := converter.ModelUserMembershipsToThrift(models)

		require.NotNil(t, result)
		assert.Len(t, result, 1) // nil应该被过滤掉
		assert.Equal(t, userID.String(), *result[0].ID)
	})
}

func TestConverterImpl_AddMembershipRequestToModel(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.AddMembershipRequestToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整请求转换", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()
		deptID := uuid.New()

		userIDStr := userID.String()
		orgIDStr := orgID.String()
		deptIDStr := deptID.String()
		isPrimary := true

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			DepartmentID:   &deptIDStr,
			IsPrimary:      isPrimary,
		}

		result := converter.AddMembershipRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, deptID, result.DepartmentID)
		assert.True(t, result.IsPrimary)
		assert.Equal(t, models.MembershipStatusActive, result.Status) // 默认为活跃状态
	})

	t.Run("最小请求转换", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()

		userIDStr := userID.String()
		orgIDStr := orgID.String()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			IsPrimary:      false,
		}

		result := converter.AddMembershipRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, uuid.Nil, result.DepartmentID) // 未设置
		assert.False(t, result.IsPrimary)
		assert.Equal(t, models.MembershipStatusActive, result.Status)
	})

	t.Run("无部门ID的请求", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()

		userIDStr := userID.String()
		orgIDStr := orgID.String()

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			DepartmentID:   nil,
			IsPrimary:      true,
		}

		result := converter.AddMembershipRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, orgID, result.OrganizationID)
		assert.Equal(t, uuid.Nil, result.DepartmentID)
		assert.True(t, result.IsPrimary)
		assert.Equal(t, models.MembershipStatusActive, result.Status)
	})

	t.Run("部分更新应用", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()
		newDeptID := uuid.New()

		req := &identity_srv.UpdateMembershipRequest{
			DepartmentID: convutil.StringPtr(newDeptID.String()),
			// OrganizationID 和 IsPrimary 未设置
		}

		existing := &models.UserMembership{
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   uuid.Nil,
			IsPrimary:      false,
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, newDeptID, result.DepartmentID) // 应该更新
		assert.Equal(t, orgID, result.OrganizationID)   // 应该保持不变
		assert.False(t, result.IsPrimary)               // 应该保持不变
		assert.Equal(t, userID, result.UserID)          // 应该保持不变
	})

	t.Run("更新为nil部门ID", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()
		deptID := uuid.New()

		req := &identity_srv.UpdateMembershipRequest{
			DepartmentID: convutil.StringPtr(uuid.Nil.String()),
		}

		existing := &models.UserMembership{
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   deptID,
			IsPrimary:      false,
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, uuid.Nil, result.DepartmentID)
	})

	t.Run("切换主要状态", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()

		isPrimary := true

		req := &identity_srv.UpdateMembershipRequest{
			IsPrimary: &isPrimary,
		}

		existing := &models.UserMembership{
			UserID:         userID,
			OrganizationID: orgID,
			IsPrimary:      false,
		}

		result := converter.ApplyUpdateToModel(existing, req)

		require.NotNil(t, result)
		assert.True(t, result.IsPrimary)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, orgID, result.OrganizationID)
	})
}

func TestConverterImpl_ModelToThrift_Alias(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("ModelToThrift别名方法测试", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()

		model := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:         userID,
			OrganizationID: orgID,
			IsPrimary:      true,
			Status:         models.MembershipStatusActive,
		}

		// 使用别名方法
		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, userID.String(), *result.ID)
		assert.Equal(t, userID.String(), *result.UserID)
		assert.Equal(t, orgID.String(), *result.OrganizationID)
		assert.NotNil(t, result.IsPrimary)
		assert.True(t, *result.IsPrimary)
	})
}

func TestConverterImpl_CompleteRoundTrip(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("Model -> Thrift -> Model 往返转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()
		deptID := uuid.New()

		original := &models.UserMembership{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserID:         userID,
			OrganizationID: orgID,
			DepartmentID:   deptID,
			IsPrimary:      true,
			Status:         models.MembershipStatusActive,
		}

		// Model -> Thrift
		thrift := converter.ModelUserMembershipToThrift(original)
		require.NotNil(t, thrift)

		// Thrift -> Model
		model := converter.ThriftUserMembershipToModel(thrift)
		require.NotNil(t, model)

		// 验证关键字段
		assert.Equal(t, original.ID, model.ID)
		assert.Equal(t, original.UserID, model.UserID)
		assert.Equal(t, original.OrganizationID, model.OrganizationID)
		assert.Equal(t, original.DepartmentID, model.DepartmentID)
		assert.Equal(t, original.IsPrimary, model.IsPrimary)
	})

	t.Run("创建请求往返测试", func(t *testing.T) {
		userID := uuid.New()
		orgID := uuid.New()

		userIDStr := userID.String()
		orgIDStr := orgID.String()
		isPrimary := true

		req := &identity_srv.AddMembershipRequest{
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			IsPrimary:      isPrimary,
		}

		// Request -> Model
		model := converter.AddMembershipRequestToModel(req)
		require.NotNil(t, model)

		// Model -> Thrift
		thrift := converter.ModelUserMembershipToThrift(model)
		require.NotNil(t, thrift)

		// 验证关键字段
		assert.Equal(t, userIDStr, *thrift.UserID)
		assert.Equal(t, orgIDStr, *thrift.OrganizationID)
		assert.NotNil(t, thrift.IsPrimary)
		assert.True(t, *thrift.IsPrimary)
	})
}

func TestConverterImpl_EdgeCases(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("零值UUID处理", func(t *testing.T) {
		zeroUUID := uuid.Nil
		idStr := zeroUUID.String()
		userIDStr := zeroUUID.String()
		orgIDStr := zeroUUID.String()

		dto := &identity_srv.UserMembership{
			ID:             &idStr,
			UserID:         &userIDStr,
			OrganizationID: &orgIDStr,
			IsPrimary:      nil,
		}

		result := converter.ThriftUserMembershipToModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, uuid.Nil, result.ID)
		assert.Equal(t, uuid.Nil, result.UserID)
		assert.Equal(t, uuid.Nil, result.OrganizationID)
		assert.False(t, result.IsPrimary)
	})

	t.Run("相同用户和组织的多个成员关系", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		orgID := uuid.New()
		deptID1 := uuid.New()
		deptID2 := uuid.New()

		models := []*models.UserMembership{
			{
				BaseModel: models.BaseModel{
					ID:        deptID1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:         userID,
				OrganizationID: orgID,
				DepartmentID:   deptID1,
				IsPrimary:      true,
				Status:         models.MembershipStatusActive,
			},
			{
				BaseModel: models.BaseModel{
					ID:        deptID2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:         userID,
				OrganizationID: orgID,
				DepartmentID:   deptID2,
				IsPrimary:      false,
				Status:         models.MembershipStatusActive,
			},
		}

		result := converter.ModelUserMembershipsToThrift(models)

		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.NotNil(t, result[0].IsPrimary)
		assert.True(t, *result[0].IsPrimary)
		assert.NotNil(t, result[1].IsPrimary)
		assert.False(t, *result[1].IsPrimary)
	})
}
