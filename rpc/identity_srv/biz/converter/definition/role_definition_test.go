package definition

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// mockEnumConverter 模拟枚举转换器
type mockEnumConverter struct{}

func (m *mockEnumConverter) ModelUserStatusToThrift(status models.UserStatus) core.UserStatus {
	switch status {
	case models.UserStatusActive:
		return core.UserStatus_ACTIVE
	case models.UserStatusInactive:
		return core.UserStatus_INACTIVE
	case models.UserStatusSuspended:
		return core.UserStatus_SUSPENDED
	case models.UserStatusLocked:
		return core.UserStatus_LOCKED
	default:
		return core.UserStatus_INACTIVE
	}
}

func (m *mockEnumConverter) ThriftUserStatusToModel(status core.UserStatus) models.UserStatus {
	switch status {
	case core.UserStatus_ACTIVE:
		return models.UserStatusActive
	case core.UserStatus_INACTIVE:
		return models.UserStatusInactive
	case core.UserStatus_SUSPENDED:
		return models.UserStatusSuspended
	case core.UserStatus_LOCKED:
		return models.UserStatusLocked
	default:
		return models.UserStatusInactive
	}
}

func (m *mockEnumConverter) ModelRoleStatusToThrift(status models.RoleStatus) core.RoleStatus {
	switch status {
	case models.RoleStatusActive:
		return core.RoleStatus_ACTIVE
	case models.RoleStatusInactive:
		return core.RoleStatus_INACTIVE
	case models.RoleStatusDeprecated:
		return core.RoleStatus_DEPRECATED
	default:
		return core.RoleStatus_INACTIVE
	}
}

func (m *mockEnumConverter) ThriftRoleStatusToModel(status core.RoleStatus) models.RoleStatus {
	switch status {
	case core.RoleStatus_ACTIVE:
		return models.RoleStatusActive
	case core.RoleStatus_INACTIVE:
		return models.RoleStatusInactive
	case core.RoleStatus_DEPRECATED:
		return models.RoleStatusDeprecated
	default:
		return models.RoleStatusInactive
	}
}

func (m *mockEnumConverter) ModelGenderToThrift(gender models.Gender) core.Gender {
	switch gender {
	case models.GenderMale:
		return core.Gender_MALE
	case models.GenderFemale:
		return core.Gender_FEMALE
	default:
		return core.Gender_UNKNOWN
	}
}

func (m *mockEnumConverter) ThriftGenderToModel(gender core.Gender) models.Gender {
	switch gender {
	case core.Gender_MALE:
		return models.GenderMale
	case core.Gender_FEMALE:
		return models.GenderFemale
	default:
		return models.GenderUnknown
	}
}

func (m *mockEnumConverter) ModelDataScopeToThrift(scope models.DataScopeType) identity_srv.DataScope {
	switch scope {
	case models.DataScopeSelf:
		return identity_srv.DataScope_SELF
	case models.DataScopeDept:
		return identity_srv.DataScope_DEPT
	case models.DataScopeOrg:
		return identity_srv.DataScope_ORG
	default:
		return identity_srv.DataScope_SELF
	}
}

func (m *mockEnumConverter) ThriftDataScopeToModel(scope identity_srv.DataScope) models.DataScopeType {
	switch scope {
	case identity_srv.DataScope_SELF:
		return models.DataScopeSelf
	case identity_srv.DataScope_DEPT:
		return models.DataScopeDept
	case identity_srv.DataScope_ORG:
		return models.DataScopeOrg
	default:
		return models.DataScopeSelf
	}
}

func TestNewConverter(t *testing.T) {
	enumConverter := &mockEnumConverter{}
	converter := NewConverter(enumConverter)

	assert.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
}

func TestConverterImpl_ModelToThrift(t *testing.T) {
	enumConverter := &mockEnumConverter{}
	converter := NewConverter(enumConverter)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("完整角色定义转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()
		createdByID := uuid.New()
		updatedByID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "Admin",
			Description:  "Administrator role with full permissions",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: true,
			CreatedBy:    &createdByID,
			UpdatedBy:    &updatedByID,
			UserCount:    25,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "Admin", *result.Name)
		assert.Equal(t, "Administrator role with full permissions", *result.Description)
		assert.Equal(t, core.RoleStatus_ACTIVE, *result.Status)
		assert.Equal(t, true, result.IsSystemRole)
		assert.Equal(t, createdByID.String(), *result.CreatedBy)
		assert.Equal(t, updatedByID.String(), *result.UpdatedBy)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, int64(25), *result.UserCount)
		assert.Empty(t, result.Permissions) // 应该是空数组
	})

	t.Run("非系统角色转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "User",
			Description:  "Regular user role",
			Status:       models.RoleStatusInactive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			CreatedBy:    nil,
			UpdatedBy:    nil,
			UserCount:    0,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "User", *result.Name)
		assert.Equal(t, "Regular user role", *result.Description)
		assert.Equal(t, core.RoleStatus_INACTIVE, *result.Status)
		assert.Equal(t, false, result.IsSystemRole)
		assert.Nil(t, result.CreatedBy)
		assert.Nil(t, result.UpdatedBy)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, int64(0), *result.UserCount)
	})

	t.Run("已弃用角色转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "Legacy",
			Description:  "Deprecated legacy role",
			Status:       models.RoleStatusDeprecated,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			CreatedBy:    nil,
			UpdatedBy:    nil,
			UserCount:    5,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "Legacy", *result.Name)
		assert.Equal(t, "Deprecated legacy role", *result.Description)
		assert.Equal(t, core.RoleStatus_DEPRECATED, *result.Status)
		assert.Equal(t, false, result.IsSystemRole)
		assert.Nil(t, result.CreatedBy)
		assert.Nil(t, result.UpdatedBy)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, int64(5), *result.UserCount)
	})

	t.Run("空描述角色转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "Test",
			Description:  "",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			CreatedBy:    nil,
			UpdatedBy:    nil,
			UserCount:    0,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "Test", *result.Name)
		assert.Nil(t, result.Description) // 空描述应该被转换为nil
		assert.Equal(t, core.RoleStatus_ACTIVE, *result.Status)
		assert.Equal(t, false, result.IsSystemRole)
		assert.Nil(t, result.CreatedBy)
		assert.Nil(t, result.UpdatedBy)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, int64(0), *result.UserCount)
	})

	t.Run("仅有创建者的角色转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()
		createdByID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "CreatorOnly",
			Description:  "Role with only creator",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			CreatedBy:    &createdByID,
			UpdatedBy:    nil,
			UserCount:    10,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "CreatorOnly", *result.Name)
		assert.Equal(t, "Role with only creator", *result.Description)
		assert.Equal(t, core.RoleStatus_ACTIVE, *result.Status)
		assert.Equal(t, false, result.IsSystemRole)
		assert.Equal(t, createdByID.String(), *result.CreatedBy)
		assert.Nil(t, result.UpdatedBy)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, int64(10), *result.UserCount)
	})

	t.Run("仅有更新者的角色转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()
		updatedByID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "UpdaterOnly",
			Description:  "Role with only updater",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			CreatedBy:    nil,
			UpdatedBy:    &updatedByID,
			UserCount:    3,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "UpdaterOnly", *result.Name)
		assert.Equal(t, "Role with only updater", *result.Description)
		assert.Equal(t, core.RoleStatus_ACTIVE, *result.Status)
		assert.Equal(t, false, result.IsSystemRole)
		assert.Nil(t, result.CreatedBy)
		assert.Equal(t, updatedByID.String(), *result.UpdatedBy)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, int64(3), *result.UserCount)
	})

	t.Run("大用户数量角色转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "Popular",
			Description:  "Popular role with many users",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			CreatedBy:    nil,
			UpdatedBy:    nil,
			UserCount:    999999,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "Popular", *result.Name)
		assert.Equal(t, "Popular role with many users", *result.Description)
		assert.Equal(t, core.RoleStatus_ACTIVE, *result.Status)
		assert.Equal(t, false, result.IsSystemRole)
		assert.Nil(t, result.CreatedBy)
		assert.Nil(t, result.UpdatedBy)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, int64(999999), *result.UserCount)
	})
}

// 表格驱动测试
func TestConverterImpl_ModelToThrift_TableDriven(t *testing.T) {
	enumConverter := &mockEnumConverter{}
	converter := NewConverter(enumConverter)

	tests := []struct {
		name        string
		input       *models.RoleDefinition
		expectNil   bool
		description string
	}{
		{
			name:        "nil输入",
			input:       nil,
			expectNil:   true,
			description: "nil输入应该返回nil",
		},
		{
			name:        "空角色定义",
			input:       &models.RoleDefinition{},
			expectNil:   false,
			description: "空角色定义应该返回非nil结果",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ModelToThrift(tt.input)

			if tt.expectNil {
				assert.Nil(t, result, tt.description)
			} else {
				assert.NotNil(t, result, tt.description)
			}
		})
	}
}

// 基准测试
func BenchmarkConverterImpl_ModelToThrift(b *testing.B) {
	enumConverter := &mockEnumConverter{}
	converter := NewConverter(enumConverter)

	now := time.Now()
	nowMs := now.UnixMilli()
	roleID := uuid.New()
	createdByID := uuid.New()
	updatedByID := uuid.New()

	model := &models.RoleDefinition{
		BaseModel: models.BaseModel{
			ID:        roleID,
			CreatedAt: nowMs,
			UpdatedAt: nowMs,
		},
		Name:         "Benchmark Role",
		Description:  "Role for benchmark testing",
		Status:       models.RoleStatusActive,
		Permissions:  models.Permissions{},
		IsSystemRole: false,
		CreatedBy:    &createdByID,
		UpdatedBy:    &updatedByID,
		UserCount:    100,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelToThrift(model)
	}
}

// TestMockEnumConverter_DataScope 测试 mockEnumConverter 的 DataScope 转换方法
func TestMockEnumConverter_DataScope(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}

	t.Run("ModelDataScopeToThrift", func(t *testing.T) {
		tests := []struct {
			name     string
			input    models.DataScopeType
			expected identity_srv.DataScope
		}{
			{
				name:     "DataScopeSelf",
				input:    models.DataScopeSelf,
				expected: identity_srv.DataScope_SELF,
			},
			{
				name:     "DataScopeDept",
				input:    models.DataScopeDept,
				expected: identity_srv.DataScope_DEPT,
			},
			{
				name:     "DataScopeOrg",
				input:    models.DataScopeOrg,
				expected: identity_srv.DataScope_ORG,
			},
			{
				name:     "Invalid DataScope",
				input:    models.DataScopeType(0),
				expected: identity_srv.DataScope_SELF, // 默认值
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := mockEnumConv.ModelDataScopeToThrift(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("ThriftDataScopeToModel", func(t *testing.T) {
		tests := []struct {
			name     string
			input    identity_srv.DataScope
			expected models.DataScopeType
		}{
			{
				name:     "DataScope_SELF",
				input:    identity_srv.DataScope_SELF,
				expected: models.DataScopeSelf,
			},
			{
				name:     "DataScope_DEPT",
				input:    identity_srv.DataScope_DEPT,
				expected: models.DataScopeDept,
			},
			{
				name:     "DataScope_ORG",
				input:    identity_srv.DataScope_ORG,
				expected: models.DataScopeOrg,
			},
			{
				name:     "Invalid DataScope",
				input:    identity_srv.DataScope(0),
				expected: models.DataScopeSelf, // 默认值
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := mockEnumConv.ThriftDataScopeToModel(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("DataScope 往返转换", func(t *testing.T) {
		scopes := []models.DataScopeType{
			models.DataScopeSelf,
			models.DataScopeDept,
			models.DataScopeOrg,
		}

		for _, original := range scopes {
			// Model -> Thrift
			thrift := mockEnumConv.ModelDataScopeToThrift(original)

			// Thrift -> Model
			result := mockEnumConv.ThriftDataScopeToModel(thrift)

			// 验证往返转换的一致性
			assert.Equal(t, original, result, "往返转换应该保持一致")
		}
	})
}

// TestConverterImpl_ModelToThrift_WithDefaultScope 测试包含 DefaultScope 的角色定义转换
func TestConverterImpl_ModelToThrift_WithDefaultScope(t *testing.T) {
	enumConverter := &mockEnumConverter{}
	converter := NewConverter(enumConverter)

	t.Run("DataScopeSelf 转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "SelfScopeRole",
			Description:  "Role with self data scope",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			DefaultScope: models.DataScopeSelf,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "SelfScopeRole", *result.Name)
		assert.NotNil(t, result.DefaultScope)
		assert.Equal(t, identity_srv.DataScope_SELF, *result.DefaultScope)
	})

	t.Run("DataScopeDept 转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "DeptScopeRole",
			Description:  "Role with department data scope",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			DefaultScope: models.DataScopeDept,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "DeptScopeRole", *result.Name)
		assert.NotNil(t, result.DefaultScope)
		assert.Equal(t, identity_srv.DataScope_DEPT, *result.DefaultScope)
	})

	t.Run("DataScopeOrg 转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "OrgScopeRole",
			Description:  "Role with organization data scope",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			DefaultScope: models.DataScopeOrg,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "OrgScopeRole", *result.Name)
		assert.NotNil(t, result.DefaultScope)
		assert.Equal(t, identity_srv.DataScope_ORG, *result.DefaultScope)
	})

	t.Run("零值 DefaultScope 处理", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		roleID := uuid.New()

		model := &models.RoleDefinition{
			BaseModel: models.BaseModel{
				ID:        roleID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Name:         "ZeroScopeRole",
			Description:  "Role with zero default scope",
			Status:       models.RoleStatusActive,
			Permissions:  models.Permissions{},
			IsSystemRole: false,
			DefaultScope: models.DataScopeType(0), // 零值
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, roleID.String(), *result.Id)
		assert.Equal(t, "ZeroScopeRole", *result.Name)
		// 零值应该被转换为默认值 DataScope_SELF
		assert.NotNil(t, result.DefaultScope)
		assert.Equal(t, identity_srv.DataScope_SELF, *result.DefaultScope)
	})
}

// BenchmarkMockEnumConverter_DataScope 基准测试 DataScope 转换
func BenchmarkMockEnumConverter_DataScope(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}

	b.Run("ModelDataScopeToThrift", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = mockEnumConv.ModelDataScopeToThrift(models.DataScopeOrg)
		}
	})

	b.Run("ThriftDataScopeToModel", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = mockEnumConv.ThriftDataScopeToModel(identity_srv.DataScope_ORG)
		}
	})
}
