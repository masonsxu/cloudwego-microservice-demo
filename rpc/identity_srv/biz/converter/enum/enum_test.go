package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	assert.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
}

func TestConverterImpl_ModelUserStatusToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("活跃状态转换", func(t *testing.T) {
		result := converter.ModelUserStatusToThrift(models.UserStatusActive)
		assert.Equal(t, core.UserStatus_ACTIVE, result)
	})

	t.Run("未激活状态转换", func(t *testing.T) {
		result := converter.ModelUserStatusToThrift(models.UserStatusInactive)
		assert.Equal(t, core.UserStatus_INACTIVE, result)
	})

	t.Run("暂停状态转换", func(t *testing.T) {
		result := converter.ModelUserStatusToThrift(models.UserStatusSuspended)
		assert.Equal(t, core.UserStatus_SUSPENDED, result)
	})

	t.Run("锁定状态转换", func(t *testing.T) {
		result := converter.ModelUserStatusToThrift(models.UserStatusLocked)
		assert.Equal(t, core.UserStatus_LOCKED, result)
	})

	t.Run("未知状态转换", func(t *testing.T) {
		result := converter.ModelUserStatusToThrift(models.UserStatus(999))
		assert.Equal(t, core.UserStatus_INACTIVE, result) // 默认值
	})

	t.Run("零值状态转换", func(t *testing.T) {
		result := converter.ModelUserStatusToThrift(models.UserStatus(0))
		assert.Equal(t, core.UserStatus_INACTIVE, result) // 默认值
	})
}

func TestConverterImpl_ThriftUserStatusToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("活跃状态转换", func(t *testing.T) {
		result := converter.ThriftUserStatusToModel(core.UserStatus_ACTIVE)
		assert.Equal(t, models.UserStatusActive, result)
	})

	t.Run("未激活状态转换", func(t *testing.T) {
		result := converter.ThriftUserStatusToModel(core.UserStatus_INACTIVE)
		assert.Equal(t, models.UserStatusInactive, result)
	})

	t.Run("暂停状态转换", func(t *testing.T) {
		result := converter.ThriftUserStatusToModel(core.UserStatus_SUSPENDED)
		assert.Equal(t, models.UserStatusSuspended, result)
	})

	t.Run("锁定状态转换", func(t *testing.T) {
		result := converter.ThriftUserStatusToModel(core.UserStatus_LOCKED)
		assert.Equal(t, models.UserStatusLocked, result)
	})

	t.Run("未知状态转换", func(t *testing.T) {
		result := converter.ThriftUserStatusToModel(core.UserStatus(999))
		assert.Equal(t, models.UserStatusInactive, result) // 默认值
	})

	t.Run("零值状态转换", func(t *testing.T) {
		result := converter.ThriftUserStatusToModel(core.UserStatus(0))
		assert.Equal(t, models.UserStatusInactive, result) // 默认值
	})
}

func TestConverterImpl_ModelRoleStatusToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("活跃状态转换", func(t *testing.T) {
		result := converter.ModelRoleStatusToThrift(models.RoleStatusActive)
		assert.Equal(t, core.RoleStatus_ACTIVE, result)
	})

	t.Run("未激活状态转换", func(t *testing.T) {
		result := converter.ModelRoleStatusToThrift(models.RoleStatusInactive)
		assert.Equal(t, core.RoleStatus_INACTIVE, result)
	})

	t.Run("已弃用状态转换", func(t *testing.T) {
		result := converter.ModelRoleStatusToThrift(models.RoleStatusDeprecated)
		assert.Equal(t, core.RoleStatus_DEPRECATED, result)
	})

	t.Run("未知状态转换", func(t *testing.T) {
		result := converter.ModelRoleStatusToThrift(models.RoleStatus(999))
		assert.Equal(t, core.RoleStatus_INACTIVE, result) // 默认值
	})

	t.Run("零值状态转换", func(t *testing.T) {
		result := converter.ModelRoleStatusToThrift(models.RoleStatus(0))
		assert.Equal(t, core.RoleStatus_INACTIVE, result) // 默认值
	})
}

func TestConverterImpl_ThriftRoleStatusToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("活跃状态转换", func(t *testing.T) {
		result := converter.ThriftRoleStatusToModel(core.RoleStatus_ACTIVE)
		assert.Equal(t, models.RoleStatusActive, result)
	})

	t.Run("未激活状态转换", func(t *testing.T) {
		result := converter.ThriftRoleStatusToModel(core.RoleStatus_INACTIVE)
		assert.Equal(t, models.RoleStatusInactive, result)
	})

	t.Run("已弃用状态转换", func(t *testing.T) {
		result := converter.ThriftRoleStatusToModel(core.RoleStatus_DEPRECATED)
		assert.Equal(t, models.RoleStatusDeprecated, result)
	})

	t.Run("未知状态转换", func(t *testing.T) {
		result := converter.ThriftRoleStatusToModel(core.RoleStatus(999))
		assert.Equal(t, models.RoleStatusInactive, result) // 默认值
	})

	t.Run("零值状态转换", func(t *testing.T) {
		result := converter.ThriftRoleStatusToModel(core.RoleStatus(0))
		assert.Equal(t, models.RoleStatusInactive, result) // 默认值
	})
}

func TestConverterImpl_ModelGenderToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("男性性别转换", func(t *testing.T) {
		result := converter.ModelGenderToThrift(models.GenderMale)
		assert.Equal(t, core.Gender_MALE, result)
	})

	t.Run("女性性别转换", func(t *testing.T) {
		result := converter.ModelGenderToThrift(models.GenderFemale)
		assert.Equal(t, core.Gender_FEMALE, result)
	})

	t.Run("未知性别转换", func(t *testing.T) {
		result := converter.ModelGenderToThrift(models.GenderUnknown)
		assert.Equal(t, core.Gender_UNKNOWN, result)
	})

	t.Run("未知性别值转换", func(t *testing.T) {
		result := converter.ModelGenderToThrift(models.Gender(999))
		assert.Equal(t, core.Gender_UNKNOWN, result) // 默认值
	})

	t.Run("零值性别转换", func(t *testing.T) {
		result := converter.ModelGenderToThrift(models.Gender(0))
		assert.Equal(t, core.Gender_UNKNOWN, result) // 默认值
	})
}

func TestConverterImpl_ThriftGenderToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("男性性别转换", func(t *testing.T) {
		result := converter.ThriftGenderToModel(core.Gender_MALE)
		assert.Equal(t, models.GenderMale, result)
	})

	t.Run("女性性别转换", func(t *testing.T) {
		result := converter.ThriftGenderToModel(core.Gender_FEMALE)
		assert.Equal(t, models.GenderFemale, result)
	})

	t.Run("未知性别转换", func(t *testing.T) {
		result := converter.ThriftGenderToModel(core.Gender_UNKNOWN)
		assert.Equal(t, models.GenderUnknown, result)
	})

	t.Run("未知性别值转换", func(t *testing.T) {
		result := converter.ThriftGenderToModel(core.Gender(999))
		assert.Equal(t, models.GenderUnknown, result) // 默认值
	})

	t.Run("零值性别转换", func(t *testing.T) {
		result := converter.ThriftGenderToModel(core.Gender(0))
		assert.Equal(t, models.GenderUnknown, result) // 默认值
	})
}

// 表格驱动测试 - UserStatus 双向转换
func TestConverterImpl_UserStatus_RoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name     string
		model    models.UserStatus
		thrift   core.UserStatus
		expected models.UserStatus
	}{
		{
			name:     "活跃状态往返转换",
			model:    models.UserStatusActive,
			thrift:   core.UserStatus_ACTIVE,
			expected: models.UserStatusActive,
		},
		{
			name:     "未激活状态往返转换",
			model:    models.UserStatusInactive,
			thrift:   core.UserStatus_INACTIVE,
			expected: models.UserStatusInactive,
		},
		{
			name:     "暂停状态往返转换",
			model:    models.UserStatusSuspended,
			thrift:   core.UserStatus_SUSPENDED,
			expected: models.UserStatusSuspended,
		},
		{
			name:     "锁定状态往返转换",
			model:    models.UserStatusLocked,
			thrift:   core.UserStatus_LOCKED,
			expected: models.UserStatusLocked,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Model -> Thrift
			thriftResult := converter.ModelUserStatusToThrift(tt.model)
			assert.Equal(t, tt.thrift, thriftResult)

			// Thrift -> Model
			modelResult := converter.ThriftUserStatusToModel(tt.thrift)
			assert.Equal(t, tt.expected, modelResult)
		})
	}
}

// 表格驱动测试 - RoleStatus 双向转换
func TestConverterImpl_RoleStatus_RoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name     string
		model    models.RoleStatus
		thrift   core.RoleStatus
		expected models.RoleStatus
	}{
		{
			name:     "活跃状态往返转换",
			model:    models.RoleStatusActive,
			thrift:   core.RoleStatus_ACTIVE,
			expected: models.RoleStatusActive,
		},
		{
			name:     "未激活状态往返转换",
			model:    models.RoleStatusInactive,
			thrift:   core.RoleStatus_INACTIVE,
			expected: models.RoleStatusInactive,
		},
		{
			name:     "已弃用状态往返转换",
			model:    models.RoleStatusDeprecated,
			thrift:   core.RoleStatus_DEPRECATED,
			expected: models.RoleStatusDeprecated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Model -> Thrift
			thriftResult := converter.ModelRoleStatusToThrift(tt.model)
			assert.Equal(t, tt.thrift, thriftResult)

			// Thrift -> Model
			modelResult := converter.ThriftRoleStatusToModel(tt.thrift)
			assert.Equal(t, tt.expected, modelResult)
		})
	}
}

// 表格驱动测试 - Gender 双向转换
func TestConverterImpl_Gender_RoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name     string
		model    models.Gender
		thrift   core.Gender
		expected models.Gender
	}{
		{
			name:     "男性性别往返转换",
			model:    models.GenderMale,
			thrift:   core.Gender_MALE,
			expected: models.GenderMale,
		},
		{
			name:     "女性性别往返转换",
			model:    models.GenderFemale,
			thrift:   core.Gender_FEMALE,
			expected: models.GenderFemale,
		},
		{
			name:     "未知性别往返转换",
			model:    models.GenderUnknown,
			thrift:   core.Gender_UNKNOWN,
			expected: models.GenderUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Model -> Thrift
			thriftResult := converter.ModelGenderToThrift(tt.model)
			assert.Equal(t, tt.thrift, thriftResult)

			// Thrift -> Model
			modelResult := converter.ThriftGenderToModel(tt.thrift)
			assert.Equal(t, tt.expected, modelResult)
		})
	}
}

// 边界情况和默认值测试
func TestConverterImpl_DefaultValues(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("UserStatus 默认值", func(t *testing.T) {
		// 未知模型值 -> Thrift
		thriftResult := converter.ModelUserStatusToThrift(models.UserStatus(999))
		assert.Equal(t, core.UserStatus_INACTIVE, thriftResult)

		// 未知 Thrift 值 -> Model
		modelResult := converter.ThriftUserStatusToModel(core.UserStatus(999))
		assert.Equal(t, models.UserStatusInactive, modelResult)
	})

	t.Run("RoleStatus 默认值", func(t *testing.T) {
		// 未知模型值 -> Thrift
		thriftResult := converter.ModelRoleStatusToThrift(models.RoleStatus(999))
		assert.Equal(t, core.RoleStatus_INACTIVE, thriftResult)

		// 未知 Thrift 值 -> Model
		modelResult := converter.ThriftRoleStatusToModel(core.RoleStatus(999))
		assert.Equal(t, models.RoleStatusInactive, modelResult)
	})

	t.Run("Gender 默认值", func(t *testing.T) {
		// 未知模型值 -> Thrift
		thriftResult := converter.ModelGenderToThrift(models.Gender(999))
		assert.Equal(t, core.Gender_UNKNOWN, thriftResult)

		// 未知 Thrift 值 -> Model
		modelResult := converter.ThriftGenderToModel(core.Gender(999))
		assert.Equal(t, models.GenderUnknown, modelResult)
	})
}

// 完整往返转换测试
func TestConverterImpl_CompleteRoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("UserStatus 完整往返", func(t *testing.T) {
		originalStatus := models.UserStatusActive

		// Model -> Thrift
		thriftStatus := converter.ModelUserStatusToThrift(originalStatus)

		// Thrift -> Model
		finalStatus := converter.ThriftUserStatusToModel(thriftStatus)

		assert.Equal(t, originalStatus, finalStatus)
	})

	t.Run("RoleStatus 完整往返", func(t *testing.T) {
		originalStatus := models.RoleStatusActive

		// Model -> Thrift
		thriftStatus := converter.ModelRoleStatusToThrift(originalStatus)

		// Thrift -> Model
		finalStatus := converter.ThriftRoleStatusToModel(thriftStatus)

		assert.Equal(t, originalStatus, finalStatus)
	})

	t.Run("Gender 完整往返", func(t *testing.T) {
		originalGender := models.GenderMale

		// Model -> Thrift
		thriftGender := converter.ModelGenderToThrift(originalGender)

		// Thrift -> Model
		finalGender := converter.ThriftGenderToModel(thriftGender)

		assert.Equal(t, originalGender, finalGender)
	})
}

// 基准测试
func BenchmarkConverterImpl_ModelUserStatusToThrift(b *testing.B) {
	converter := &ConverterImpl{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelUserStatusToThrift(models.UserStatusActive)
	}
}

func BenchmarkConverterImpl_ThriftUserStatusToModel(b *testing.B) {
	converter := &ConverterImpl{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ThriftUserStatusToModel(core.UserStatus_ACTIVE)
	}
}

func BenchmarkConverterImpl_ModelRoleStatusToThrift(b *testing.B) {
	converter := &ConverterImpl{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelRoleStatusToThrift(models.RoleStatusActive)
	}
}

func BenchmarkConverterImpl_ThriftRoleStatusToModel(b *testing.B) {
	converter := &ConverterImpl{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ThriftRoleStatusToModel(core.RoleStatus_ACTIVE)
	}
}

func BenchmarkConverterImpl_ModelGenderToThrift(b *testing.B) {
	converter := &ConverterImpl{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelGenderToThrift(models.GenderMale)
	}
}

func BenchmarkConverterImpl_ThriftGenderToModel(b *testing.B) {
	converter := &ConverterImpl{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ThriftGenderToModel(core.Gender_MALE)
	}
}

// 综合基准测试 - 模拟真实使用场景
func BenchmarkConverterImpl_CompleteConversion(b *testing.B) {
	converter := &ConverterImpl{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 模拟完整的枚举转换流程
		userStatus := converter.ModelUserStatusToThrift(models.UserStatusActive)
		_ = converter.ThriftUserStatusToModel(userStatus)

		roleStatus := converter.ModelRoleStatusToThrift(models.RoleStatusActive)
		_ = converter.ThriftRoleStatusToModel(roleStatus)

		gender := converter.ModelGenderToThrift(models.GenderMale)
		_ = converter.ThriftGenderToModel(gender)
	}
}
