package permission

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// mockEnumConverter 是枚举转换器的模拟实现
type mockEnumConverter struct{}

func (m *mockEnumConverter) ModelUserStatusToThrift(status models.UserStatus) core.UserStatus {
	switch status {
	case models.UserStatusActive:
		return core.UserStatus_ACTIVE
	case models.UserStatusInactive:
		return core.UserStatus_INACTIVE
	case models.UserStatusSuspended:
		return core.UserStatus_SUSPENDED
	default:
		return core.UserStatus_LOCKED
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
	default:
		return models.UserStatusLocked
	}
}

func (m *mockEnumConverter) ModelRoleStatusToThrift(status models.RoleStatus) core.RoleStatus {
	switch status {
	case models.RoleStatusActive:
		return core.RoleStatus_ACTIVE
	case models.RoleStatusInactive:
		return core.RoleStatus_INACTIVE
	default:
		return core.RoleStatus_DEPRECATED
	}
}

func (m *mockEnumConverter) ThriftRoleStatusToModel(status core.RoleStatus) models.RoleStatus {
	switch status {
	case core.RoleStatus_ACTIVE:
		return models.RoleStatusActive
	case core.RoleStatus_INACTIVE:
		return models.RoleStatusInactive
	default:
		return models.RoleStatusDeprecated
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

// TestNewConverter 测试转换器构造函数
func TestNewConverter(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	require.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
	assert.Equal(t, mockEnumConv, converter.(*ConverterImpl).enumConverter)
}

// TestConverterImpl_ModelToThrift 测试 ModelToThrift 方法
func TestConverterImpl_ModelToThrift(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("完整权限转换", func(t *testing.T) {
		model := &models.Permission{
			Resource:    "user:read",
			Action:      "read",
			Description: "读取用户信息权限",
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "user:read", *result.Resource)
		assert.Equal(t, "read", *result.Action)
		assert.Equal(t, "读取用户信息权限", *result.Description)
	})

	t.Run("无描述权限转换", func(t *testing.T) {
		model := &models.Permission{
			Resource: "admin:write",
			Action:   "write",
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "admin:write", *result.Resource)
		assert.Equal(t, "write", *result.Action)
		assert.Equal(t, "", *result.Description)
	})

	t.Run("空字符串字段转换", func(t *testing.T) {
		model := &models.Permission{
			Resource:    "",
			Action:      "",
			Description: "",
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "", *result.Resource)
		assert.Equal(t, "", *result.Action)
		assert.Equal(t, "", *result.Description)
	})

	t.Run("特殊字符权限转换", func(t *testing.T) {
		model := &models.Permission{
			Resource:    "api/v1/users:*",
			Action:      "execute",
			Description: "执行API操作（包含特殊字符：*、/、:）",
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "api/v1/users:*", *result.Resource)
		assert.Equal(t, "execute", *result.Action)
		assert.Equal(t, "执行API操作（包含特殊字符：*、/、:）", *result.Description)
	})
}

// TestConverterImpl_ThriftToModel 测试 ThriftToModel 方法
func TestConverterImpl_ThriftToModel(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ThriftToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整权限转换", func(t *testing.T) {
		resource := "role:read"
		action := "read"
		description := "读取角色信息权限"
		thrift := &identity_srv.Permission{
			Resource:    &resource,
			Action:      &action,
			Description: &description,
		}

		result := converter.ThriftToModel(thrift)

		require.NotNil(t, result)
		assert.Equal(t, "role:read", result.Resource)
		assert.Equal(t, "read", result.Action)
		assert.Equal(t, "读取角色信息权限", result.Description)
	})

	t.Run("无描述权限转换", func(t *testing.T) {
		resource := "system:write"
		action := "write"
		thrift := &identity_srv.Permission{
			Resource: &resource,
			Action:   &action,
		}

		result := converter.ThriftToModel(thrift)

		require.NotNil(t, result)
		assert.Equal(t, "system:write", result.Resource)
		assert.Equal(t, "write", result.Action)
		assert.Equal(t, "", result.Description)
	})

	t.Run("nil描述处理", func(t *testing.T) {
		resource := "menu:read"
		action := "read"
		thrift := &identity_srv.Permission{
			Resource:    &resource,
			Action:      &action,
			Description: nil,
		}

		result := converter.ThriftToModel(thrift)

		require.NotNil(t, result)
		assert.Equal(t, "menu:read", result.Resource)
		assert.Equal(t, "read", result.Action)
		assert.Equal(t, "", result.Description)
	})

	t.Run("空字符串字段转换", func(t *testing.T) {
		resource := ""
		action := ""
		description := ""
		thrift := &identity_srv.Permission{
			Resource:    &resource,
			Action:      &action,
			Description: &description,
		}

		result := converter.ThriftToModel(thrift)

		require.NotNil(t, result)
		assert.Equal(t, "", result.Resource)
		assert.Equal(t, "", result.Action)
		assert.Equal(t, "", result.Description)
	})

	t.Run("特殊字符权限转换", func(t *testing.T) {
		resource := "config:*"
		action := "manage"
		description := "管理配置（包含通配符：*）"
		thrift := &identity_srv.Permission{
			Resource:    &resource,
			Action:      &action,
			Description: &description,
		}

		result := converter.ThriftToModel(thrift)

		require.NotNil(t, result)
		assert.Equal(t, "config:*", result.Resource)
		assert.Equal(t, "manage", result.Action)
		assert.Equal(t, "管理配置（包含通配符：*）", result.Description)
	})
}

// TestConverterImpl_ModelSliceToThrift 测试 ModelSliceToThrift 方法
func TestConverterImpl_ModelSliceToThrift(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil切片转换", func(t *testing.T) {
		result := converter.ModelSliceToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("空切片转换", func(t *testing.T) {
		models := []*models.Permission{}
		result := converter.ModelSliceToThrift(models)

		require.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("多个权限转换", func(t *testing.T) {
		permissions := []*models.Permission{
			{
				Resource:    "user:read",
				Action:      "read",
				Description: "读取用户",
			},
			{
				Resource:    "user:write",
				Action:      "write",
				Description: "写入用户",
			},
			{
				Resource:    "admin:*",
				Action:      "all",
				Description: "管理员所有权限",
			},
		}

		result := converter.ModelSliceToThrift(permissions)

		require.NotNil(t, result)
		assert.Equal(t, 3, len(result))

		// 验证第一个权限
		assert.Equal(t, "user:read", *result[0].Resource)
		assert.Equal(t, "read", *result[0].Action)
		assert.Equal(t, "读取用户", *result[0].Description)

		// 验证第二个权限
		assert.Equal(t, "user:write", *result[1].Resource)
		assert.Equal(t, "write", *result[1].Action)
		assert.Equal(t, "写入用户", *result[1].Description)

		// 验证第三个权限
		assert.Equal(t, "admin:*", *result[2].Resource)
		assert.Equal(t, "all", *result[2].Action)
		assert.Equal(t, "管理员所有权限", *result[2].Description)
	})

	t.Run("包含nil元素的切片", func(t *testing.T) {
		permissions := []*models.Permission{
			{
				Resource: "valid:read",
				Action:   "read",
			},
			nil,
			{
				Resource: "valid:write",
				Action:   "write",
			},
		}

		result := converter.ModelSliceToThrift(permissions)

		require.NotNil(t, result)
		assert.Equal(t, 3, len(result))

		// 验证第一个权限
		assert.NotNil(t, result[0])
		assert.Equal(t, "valid:read", *result[0].Resource)

		// 验证第二个元素（nil转换结果）
		assert.Nil(t, result[1])

		// 验证第三个权限
		assert.NotNil(t, result[2])
		assert.Equal(t, "valid:write", *result[2].Resource)
	})

	t.Run("大量权限转换", func(t *testing.T) {
		permissions := make([]*models.Permission, 100)
		for i := 0; i < 100; i++ {
			permissions[i] = &models.Permission{
				Resource:    "resource:" + fmt.Sprintf("%d", i),
				Action:      "action:" + fmt.Sprintf("%d", i),
				Description: "权限描述 " + fmt.Sprintf("%d", i),
			}
		}

		result := converter.ModelSliceToThrift(permissions)

		require.NotNil(t, result)
		assert.Equal(t, 100, len(result))

		// 验证几个样本
		assert.Equal(t, "resource:65", *result[65].Resource)
		assert.Equal(t, "action:65", *result[65].Action)
		assert.Equal(t, "权限描述 65", *result[65].Description)
	})
}

// TestConverterImpl_ThriftSliceToModel 测试 ThriftSliceToModel 方法
func TestConverterImpl_ThriftSliceToModel(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil切片转换", func(t *testing.T) {
		result := converter.ThriftSliceToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("空切片转换", func(t *testing.T) {
		thrifts := []*identity_srv.Permission{}
		result := converter.ThriftSliceToModel(thrifts)

		require.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("多个权限转换", func(t *testing.T) {
		resource1 := "api:read"
		action1 := "read"
		desc1 := "读取API"

		resource2 := "api:write"
		action2 := "write"
		desc2 := "写入API"

		resource3 := "system:*"
		action3 := "all"
		desc3 := "系统所有权限"

		thrifts := []*identity_srv.Permission{
			{
				Resource:    &resource1,
				Action:      &action1,
				Description: &desc1,
			},
			{
				Resource:    &resource2,
				Action:      &action2,
				Description: &desc2,
			},
			{
				Resource:    &resource3,
				Action:      &action3,
				Description: &desc3,
			},
		}

		result := converter.ThriftSliceToModel(thrifts)

		require.NotNil(t, result)
		assert.Equal(t, 3, len(result))

		// 验证第一个权限
		assert.Equal(t, "api:read", result[0].Resource)
		assert.Equal(t, "read", result[0].Action)
		assert.Equal(t, "读取API", result[0].Description)

		// 验证第二个权限
		assert.Equal(t, "api:write", result[1].Resource)
		assert.Equal(t, "write", result[1].Action)
		assert.Equal(t, "写入API", result[1].Description)

		// 验证第三个权限
		assert.Equal(t, "system:*", result[2].Resource)
		assert.Equal(t, "all", result[2].Action)
		assert.Equal(t, "系统所有权限", result[2].Description)
	})

	t.Run("包含nil元素的切片", func(t *testing.T) {
		resource1 := "valid:read"
		action1 := "read"

		resource2 := "valid:write"
		action2 := "write"

		thrifts := []*identity_srv.Permission{
			{
				Resource: &resource1,
				Action:   &action1,
			},
			nil,
			{
				Resource: &resource2,
				Action:   &action2,
			},
		}

		result := converter.ThriftSliceToModel(thrifts)

		require.NotNil(t, result)
		assert.Equal(t, 3, len(result))

		// 验证第一个权限
		assert.NotNil(t, result[0])
		assert.Equal(t, "valid:read", result[0].Resource)

		// 验证第二个元素（nil转换结果）
		assert.Nil(t, result[1])

		// 验证第三个权限
		assert.NotNil(t, result[2])
		assert.Equal(t, "valid:write", result[2].Resource)
	})

	t.Run("大量权限转换", func(t *testing.T) {
		thrifts := make([]*identity_srv.Permission, 100)

		for i := 0; i < 100; i++ {
			resource := "resource:" + fmt.Sprintf("%d", i)
			action := "action:" + fmt.Sprintf("%d", i)
			description := "权限描述 " + fmt.Sprintf("%d", i)
			thrifts[i] = &identity_srv.Permission{
				Resource:    &resource,
				Action:      &action,
				Description: &description,
			}
		}

		result := converter.ThriftSliceToModel(thrifts)

		require.NotNil(t, result)
		assert.Equal(t, 100, len(result))

		// 验证几个样本
		assert.Equal(t, "resource:65", result[65].Resource)
		assert.Equal(t, "action:65", result[65].Action)
		assert.Equal(t, "权限描述 65", result[65].Description)
	})
}

// TestConverterImpl_CompleteRoundTrip 测试往返转换
func TestConverterImpl_CompleteRoundTrip(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("Model -> Thrift -> Model 往返转换", func(t *testing.T) {
		original := &models.Permission{
			Resource:    "user:profile:read",
			Action:      "read",
			Description: "读取用户档案信息",
		}

		// Model -> Thrift
		thrift := converter.ModelToThrift(original)
		require.NotNil(t, thrift)

		// Thrift -> Model
		result := converter.ThriftToModel(thrift)
		require.NotNil(t, result)

		// 验证往返转换的一致性
		assert.Equal(t, original.Resource, result.Resource)
		assert.Equal(t, original.Action, result.Action)
		assert.Equal(t, original.Description, result.Description)
	})

	t.Run("Thrift -> Model -> Thrift 往返转换", func(t *testing.T) {
		resource := "role:assignment:write"
		action := "write"
		description := "写入角色分配信息"
		original := &identity_srv.Permission{
			Resource:    &resource,
			Action:      &action,
			Description: &description,
		}

		// Thrift -> Model
		model := converter.ThriftToModel(original)
		require.NotNil(t, model)

		// Model -> Thrift
		result := converter.ModelToThrift(model)
		require.NotNil(t, result)

		// 验证往返转换的一致性
		assert.Equal(t, *original.Resource, *result.Resource)
		assert.Equal(t, *original.Action, *result.Action)
		assert.Equal(t, *original.Description, *result.Description)
	})

	t.Run("切片往返转换", func(t *testing.T) {
		original := []*models.Permission{
			{
				Resource:    "menu:read",
				Action:      "read",
				Description: "读取菜单",
			},
			{
				Resource:    "menu:write",
				Action:      "write",
				Description: "写入菜单",
			},
		}

		// Model -> Thrift
		thriftSlice := converter.ModelSliceToThrift(original)
		require.NotNil(t, thriftSlice)
		assert.Equal(t, 2, len(thriftSlice))

		// Thrift -> Model
		resultSlice := converter.ThriftSliceToModel(thriftSlice)
		require.NotNil(t, resultSlice)
		assert.Equal(t, 2, len(resultSlice))

		// 验证往返转换的一致性
		for i, perm := range original {
			assert.Equal(t, perm.Resource, resultSlice[i].Resource)
			assert.Equal(t, perm.Action, resultSlice[i].Action)
			assert.Equal(t, perm.Description, resultSlice[i].Description)
		}
	})
}

// TestConverterImpl_EdgeCases 测试边界情况
func TestConverterImpl_EdgeCases(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("空字符串字段处理", func(t *testing.T) {
		model := &models.Permission{
			Resource:    "",
			Action:      "",
			Description: "",
		}

		thrift := converter.ModelToThrift(model)
		require.NotNil(t, thrift)

		result := converter.ThriftToModel(thrift)
		require.NotNil(t, result)

		assert.Equal(t, "", result.Resource)
		assert.Equal(t, "", result.Action)
		assert.Equal(t, "", result.Description)
	})

	t.Run("nil字段处理", func(t *testing.T) {
		thrift := &identity_srv.Permission{
			Resource:    nil,
			Action:      nil,
			Description: nil,
		}

		// 这应该会导致panic，因为代码中直接解引用指针
		// 但在实际使用中，Thrift生成的代码通常不会创建nil指针
		// 所以这个测试用例主要验证nil输入的处理
		assert.Panics(t, func() {
			converter.ThriftToModel(thrift)
		})
	})

	t.Run("超长字符串处理", func(t *testing.T) {
		longResource := "resource:" + string(make([]byte, 1000))
		longAction := "action:" + string(make([]byte, 1000))
		longDescription := "description:" + string(make([]byte, 1000))

		model := &models.Permission{
			Resource:    longResource,
			Action:      longAction,
			Description: longDescription,
		}

		thrift := converter.ModelToThrift(model)
		require.NotNil(t, thrift)

		result := converter.ThriftToModel(thrift)
		require.NotNil(t, result)

		assert.Equal(t, longResource, result.Resource)
		assert.Equal(t, longAction, result.Action)
		assert.Equal(t, longDescription, result.Description)
	})

	t.Run("Unicode字符处理", func(t *testing.T) {
		model := &models.Permission{
			Resource:    "用户:管理:读取",
			Action:      "读取",
			Description: "读取用户管理权限（包含中文和特殊符号：🔒）",
		}

		thrift := converter.ModelToThrift(model)
		require.NotNil(t, thrift)

		result := converter.ThriftToModel(thrift)
		require.NotNil(t, result)

		assert.Equal(t, "用户:管理:读取", result.Resource)
		assert.Equal(t, "读取", result.Action)
		assert.Equal(t, "读取用户管理权限（包含中文和特殊符号：🔒）", result.Description)
	})
}

// BenchmarkConverterImpl_ModelToThrift 基准测试 ModelToThrift 方法
func BenchmarkConverterImpl_ModelToThrift(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	model := &models.Permission{
		Resource:    "benchmark:resource",
		Action:      "benchmark:action",
		Description: "基准测试权限描述",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelToThrift(model)
	}
}

// BenchmarkConverterImpl_ThriftToModel 基准测试 ThriftToModel 方法
func BenchmarkConverterImpl_ThriftToModel(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	resource := "benchmark:resource"
	action := "benchmark:action"
	description := "基准测试权限描述"
	thrift := &identity_srv.Permission{
		Resource:    &resource,
		Action:      &action,
		Description: &description,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ThriftToModel(thrift)
	}
}

// BenchmarkConverterImpl_ModelSliceToThrift 基准测试 ModelSliceToThrift 方法
func BenchmarkConverterImpl_ModelSliceToThrift(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	perms := make([]*models.Permission, 100)
	for i := 0; i < 100; i++ {
		perms[i] = &models.Permission{
			Resource:    "resource:" + fmt.Sprintf("%d", i),
			Action:      "action:" + fmt.Sprintf("%d", i),
			Description: "权限描述 " + fmt.Sprintf("%d", i),
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelSliceToThrift(perms)
	}
}

// BenchmarkConverterImpl_ThriftSliceToModel 基准测试 ThriftSliceToModel 方法
func BenchmarkConverterImpl_ThriftSliceToModel(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	thrifts := make([]*identity_srv.Permission, 100)

	for i := 0; i < 100; i++ {
		resource := "resource:" + fmt.Sprintf("%d", i)
		action := "action:" + fmt.Sprintf("%d", i)
		description := "权限描述 " + fmt.Sprintf("%d", i)
		thrifts[i] = &identity_srv.Permission{
			Resource:    &resource,
			Action:      &action,
			Description: &description,
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ThriftSliceToModel(thrifts)
	}
}
