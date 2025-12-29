package definition

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
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
