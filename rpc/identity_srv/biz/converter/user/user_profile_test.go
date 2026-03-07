package user

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

// TestNewConverter 测试转换器构造函数
func TestNewConverter(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	require.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
	assert.Equal(t, mockEnumConv, converter.(*ConverterImpl).enumConverter)
	assert.NotNil(t, converter.(*ConverterImpl).baseConverter)
}

// TestConverterImpl_ModelUserProfileToThrift 测试 ModelUserProfileToThrift 方法
func TestConverterImpl_ModelUserProfileToThrift(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelUserProfileToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("完整用户档案转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		createdBy := uuid.New()
		updatedBy := uuid.New()
		accountExpiry := now + 86400000 // 24小时后过期
		lastLoginTime := now - 3600000  // 1小时前登录

		model := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:           "testuser",
			PasswordHash:       "hashed_password",
			Email:              "test@example.com",
			Phone:              "13800138000",
			FirstName:          "John",
			LastName:           "Doe",
			RealName:           "John Doe",
			Gender:             models.GenderMale,
			ProfessionalTitle:  "Software Engineer",
			EmployeeID:         "EMP001",
			Status:             models.UserStatusActive,
			LoginAttempts:      2,
			MustChangePassword: true,
			AccountExpiry:      &accountExpiry,
			CreatedBy:          &createdBy,
			UpdatedBy:          &updatedBy,
			LastLoginTime:      &lastLoginTime,
		}

		result := converter.ModelUserProfileToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, userID.String(), *result.ID)
		assert.Equal(t, "testuser", *result.Username)
		assert.Equal(t, "test@example.com", *result.Email)
		assert.Equal(t, "13800138000", *result.Phone)
		assert.Equal(t, "John", *result.FirstName)
		assert.Equal(t, "Doe", *result.LastName)
		assert.Equal(t, "John Doe", *result.RealName)
		assert.Equal(t, core.Gender_MALE, *result.Gender)
		assert.Equal(t, "Software Engineer", *result.ProfessionalTitle)
		assert.Equal(t, "EMP001", *result.EmployeeID)
		assert.Equal(t, core.UserStatus_ACTIVE, *result.Status)
		assert.Equal(t, int32(2), result.LoginAttempts)
		assert.True(t, result.MustChangePassword)
		assert.Equal(t, &accountExpiry, result.AccountExpiry)
		assert.Equal(t, createdBy.String(), *result.CreatedBy)
		assert.Equal(t, updatedBy.String(), *result.UpdatedBy)
		assert.Equal(t, &lastLoginTime, result.LastLoginTime)
	})

	t.Run("最小用户档案转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()

		model := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username: "minimal",
			Status:   models.UserStatusInactive,
		}

		result := converter.ModelUserProfileToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, userID.String(), *result.ID)
		assert.Equal(t, "minimal", *result.Username)
		assert.Equal(t, core.UserStatus_INACTIVE, *result.Status)
		assert.Equal(t, &now, result.CreatedAt)
		assert.Equal(t, &now, result.UpdatedAt)

		// 验证可选字段为nil
		assert.Nil(t, result.Email)
		assert.Nil(t, result.Phone)
		assert.Nil(t, result.FirstName)
		assert.Nil(t, result.LastName)
		assert.Nil(t, result.RealName)
		assert.Nil(t, result.Gender)
		assert.Nil(t, result.ProfessionalTitle)
		assert.Nil(t, result.EmployeeID)
		assert.Equal(t, int32(0), result.LoginAttempts)
		assert.False(t, result.MustChangePassword)
		assert.Nil(t, result.AccountExpiry)
		assert.Nil(t, result.CreatedBy)
		assert.Nil(t, result.UpdatedBy)
		assert.Nil(t, result.LastLoginTime)
	})

	t.Run("空字符串字段处理", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()

		model := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:  "emptyfields",
			Email:     "",
			Phone:     "",
			FirstName: "",
			LastName:  "",
			RealName:  "",
			Status:    models.UserStatusActive,
		}

		result := converter.ModelUserProfileToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "emptyfields", *result.Username)

		// 空字符串字段应该被忽略（不设置指针）
		assert.Nil(t, result.Email)
		assert.Nil(t, result.Phone)
		assert.Nil(t, result.FirstName)
		assert.Nil(t, result.LastName)
		assert.Nil(t, result.RealName)
	})

	t.Run("零值字段处理", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()

		var zeroTime int64 = 0

		model := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:           "zerofields",
			Gender:             models.GenderUnknown, // 零值
			LoginAttempts:      0,                    // 零值
			MustChangePassword: false,                // 零值
			AccountExpiry:      &zeroTime,            // 零值时间戳
			LastLoginTime:      &zeroTime,            // 零值时间戳
			Status:             models.UserStatusActive,
		}

		result := converter.ModelUserProfileToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "zerofields", *result.Username)

		// 零值性别应该被忽略
		assert.Nil(t, result.Gender)

		// 零值登录尝试次数应该被忽略
		assert.Equal(t, int32(0), result.LoginAttempts)

		// 零值时间戳应该被忽略
		assert.Nil(t, result.AccountExpiry)
		assert.Nil(t, result.LastLoginTime)

		// 布尔值总是设置
		assert.False(t, result.MustChangePassword)
	})
}

// TestConverterImpl_CreateUserRequestToModel 测试 CreateUserRequestToModel 方法
func TestConverterImpl_CreateUserRequestToModel(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.CreateUserRequestToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整创建请求转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		accountExpiry := now + 86400000

		username := "newuser"
		password := "password123"
		email := "newuser@example.com"
		phone := "13900139000"
		firstName := "Jane"
		lastName := "Smith"
		realName := "Jane Smith"
		gender := core.Gender_FEMALE
		professionalTitle := "Product Manager"
		employeeID := "EMP002"
		mustChangePassword := true

		req := &identity_srv.CreateUserRequest{
			Username:           &username,
			Password:           &password,
			Email:              &email,
			Phone:              &phone,
			FirstName:          &firstName,
			LastName:           &lastName,
			RealName:           &realName,
			Gender:             &gender,
			ProfessionalTitle:  &professionalTitle,
			EmployeeID:         &employeeID,
			MustChangePassword: &mustChangePassword,
			AccountExpiry:      &accountExpiry,
		}

		result := converter.CreateUserRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, "newuser", result.Username)
		assert.NotEmpty(t, result.PasswordHash) // 密码应该被哈希
		assert.Equal(t, "newuser@example.com", result.Email)
		assert.Equal(t, "13900139000", result.Phone)
		assert.Equal(t, "Jane", result.FirstName)
		assert.Equal(t, "Smith", result.LastName)
		assert.Equal(t, "Jane Smith", result.RealName)
		assert.Equal(t, models.GenderFemale, result.Gender)
		assert.Equal(t, "Product Manager", result.ProfessionalTitle)
		assert.Equal(t, "EMP002", result.EmployeeID)
		assert.Equal(t, models.UserStatusActive, result.Status) // 默认激活状态
		assert.True(t, result.MustChangePassword)
		assert.Equal(t, &accountExpiry, result.AccountExpiry)
	})

	t.Run("最小创建请求转换", func(t *testing.T) {
		username := "minimal"
		password := "pass"

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
		}

		result := converter.CreateUserRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, "minimal", result.Username)
		assert.NotEmpty(t, result.PasswordHash)
		assert.Equal(t, models.UserStatusActive, result.Status) // 默认激活状态
		assert.False(t, result.MustChangePassword)              // 默认值
	})

	t.Run("密码哈希处理", func(t *testing.T) {
		username := "testuser"
		password := "" // 空密码

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
		}

		result := converter.CreateUserRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, "testuser", result.Username)
		assert.NotEmpty(t, result.PasswordHash) // 空密码也可以被哈希
	})
}

// TestConverterImpl_ApplyUpdateUserToModel 测试 ApplyUpdateUserToModel 方法
func TestConverterImpl_ApplyUpdateUserToModel(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("nil输入处理", func(t *testing.T) {
		// nil existing
		result := converter.ApplyUpdateUserToModel(nil, &identity_srv.UpdateUserRequest{})
		assert.Nil(t, result)

		// nil request
		existing := &models.UserProfile{Username: "test"}
		result = converter.ApplyUpdateUserToModel(existing, nil)
		assert.Equal(t, existing, result)
	})

	t.Run("完整更新应用", func(t *testing.T) {
		now := time.Now().UnixMilli()
		accountExpiry := now + 86400000

		email := "updated@example.com"
		phone := "13800138001"
		firstName := "Updated"
		lastName := "Name"
		realName := "Updated Name"
		gender := core.Gender_MALE
		professionalTitle := "Senior Engineer"
		employeeID := "EMP003"

		req := &identity_srv.UpdateUserRequest{
			Email:             &email,
			Phone:             &phone,
			FirstName:         &firstName,
			LastName:          &lastName,
			RealName:          &realName,
			Gender:            &gender,
			ProfessionalTitle: &professionalTitle,
			EmployeeID:        &employeeID,
			AccountExpiry:     &accountExpiry,
		}

		existing := &models.UserProfile{
			Username: "existinguser",
			Status:   models.UserStatusActive,
		}

		result := converter.ApplyUpdateUserToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, "updated@example.com", result.Email)
		assert.Equal(t, "13800138001", result.Phone)
		assert.Equal(t, "Updated", result.FirstName)
		assert.Equal(t, "Name", result.LastName)
		assert.Equal(t, "Updated Name", result.RealName)
		assert.Equal(t, models.GenderMale, result.Gender)
		assert.Equal(t, "Senior Engineer", result.ProfessionalTitle)
		assert.Equal(t, "EMP003", result.EmployeeID)
		assert.Equal(t, &accountExpiry, result.AccountExpiry)

		// 未更新的字段应该保持不变
		assert.Equal(t, "existinguser", result.Username)
		assert.Equal(t, models.UserStatusActive, result.Status)
	})

	t.Run("部分更新应用", func(t *testing.T) {
		newEmail := "partial@example.com"
		req := &identity_srv.UpdateUserRequest{
			Email: &newEmail,
		}

		existing := &models.UserProfile{
			Username: "partialuser",
			Email:    "old@example.com",
			Phone:    "13800138000",
			Status:   models.UserStatusActive,
		}

		result := converter.ApplyUpdateUserToModel(existing, req)

		require.NotNil(t, result)
		assert.Equal(t, "partial@example.com", result.Email)    // 应该更新
		assert.Equal(t, "13800138000", result.Phone)            // 应该保持不变
		assert.Equal(t, "partialuser", result.Username)         // 应该保持不变
		assert.Equal(t, models.UserStatusActive, result.Status) // 应该保持不变
	})
}

// TestConverterImpl_CompleteRoundTrip 测试往返转换
func TestConverterImpl_CompleteRoundTrip(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("完整用户档案往返转换", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()
		accountExpiry := now + 86400000

		original := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:           "roundtrip",
			Email:              "roundtrip@example.com",
			Phone:              "13800138000",
			FirstName:          "Round",
			LastName:           "Trip",
			RealName:           "Round Trip",
			Gender:             models.GenderFemale,
			ProfessionalTitle:  "Test Engineer",
			EmployeeID:         "EMP999",
			Status:             models.UserStatusActive,
			LoginAttempts:      1,
			MustChangePassword: false,
			AccountExpiry:      &accountExpiry,
		}

		// Model -> Thrift
		thrift := converter.ModelUserProfileToThrift(original)
		require.NotNil(t, thrift)

		// 注意：不能直接往返转换，因为 Thrift -> Model 的转换方法不存在
		// 这里主要验证转换的正确性和完整性
		assert.Equal(t, original.Username, *thrift.Username)
		assert.Equal(t, original.Email, *thrift.Email)
		assert.Equal(t, original.Phone, *thrift.Phone)
		assert.Equal(t, original.FirstName, *thrift.FirstName)
		assert.Equal(t, original.LastName, *thrift.LastName)
		assert.Equal(t, original.RealName, *thrift.RealName)
		assert.Equal(t, core.Gender_FEMALE, *thrift.Gender)
		assert.Equal(t, original.ProfessionalTitle, *thrift.ProfessionalTitle)
		assert.Equal(t, original.EmployeeID, *thrift.EmployeeID)
		assert.Equal(t, core.UserStatus_ACTIVE, *thrift.Status)
		assert.Equal(t, original.LoginAttempts, thrift.LoginAttempts)
		assert.Equal(t, original.MustChangePassword, thrift.MustChangePassword)
		assert.Equal(t, original.AccountExpiry, thrift.AccountExpiry)
	})

	t.Run("创建请求往返测试", func(t *testing.T) {
		username := "roundtripuser"
		password := "testpass"
		email := "rt@example.com"

		req := &identity_srv.CreateUserRequest{
			Username: &username,
			Password: &password,
			Email:    &email,
		}

		// Request -> Model
		model := converter.CreateUserRequestToModel(req)
		require.NotNil(t, model)

		// Model -> Thrift
		thrift := converter.ModelUserProfileToThrift(model)
		require.NotNil(t, thrift)

		// 验证关键字段
		assert.Equal(t, username, *thrift.Username)
		assert.Equal(t, email, *thrift.Email)
		assert.Equal(t, core.UserStatus_ACTIVE, *thrift.Status) // 默认激活状态
	})
}

// TestConverterImpl_EdgeCases 测试边界情况
func TestConverterImpl_EdgeCases(t *testing.T) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	t.Run("超长字符串处理", func(t *testing.T) {
		longString := string(make([]byte, 1000))
		now := time.Now().UnixMilli()
		userID := uuid.New()

		model := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:  "longuser",
			Email:     longString,
			Phone:     longString,
			FirstName: longString,
			LastName:  longString,
			RealName:  longString,
			Status:    models.UserStatusActive,
		}

		thrift := converter.ModelUserProfileToThrift(model)
		require.NotNil(t, thrift)

		assert.Equal(t, longString, *thrift.Email)
		assert.Equal(t, longString, *thrift.Phone)
		assert.Equal(t, longString, *thrift.FirstName)
		assert.Equal(t, longString, *thrift.LastName)
		assert.Equal(t, longString, *thrift.RealName)
	})

	t.Run("Unicode字符处理", func(t *testing.T) {
		now := time.Now().UnixMilli()
		userID := uuid.New()

		model := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:          "unicode用户",
			Email:             "用户@example.com",
			FirstName:         "名",
			LastName:          "姓",
			RealName:          "中文姓名",
			ProfessionalTitle: "高级工程师🔧",
			Status:            models.UserStatusActive,
		}

		thrift := converter.ModelUserProfileToThrift(model)
		require.NotNil(t, thrift)

		assert.Equal(t, "unicode用户", *thrift.Username)
		assert.Equal(t, "用户@example.com", *thrift.Email)
		assert.Equal(t, "名", *thrift.FirstName)
		assert.Equal(t, "姓", *thrift.LastName)
		assert.Equal(t, "中文姓名", *thrift.RealName)
		assert.Equal(t, "高级工程师🔧", *thrift.ProfessionalTitle)
	})

	t.Run("特殊字符处理", func(t *testing.T) {
		specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		now := time.Now().UnixMilli()
		userID := uuid.New()

		model := &models.UserProfile{
			BaseModel: models.BaseModel{
				ID:        userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:          "special",
			Email:             "special@example.com",
			ProfessionalTitle: specialChars,
			EmployeeID:        specialChars,
			Status:            models.UserStatusActive,
		}

		thrift := converter.ModelUserProfileToThrift(model)
		require.NotNil(t, thrift)

		assert.Equal(t, specialChars, *thrift.ProfessionalTitle)
		assert.Equal(t, specialChars, *thrift.EmployeeID)
	})
}

// BenchmarkConverterImpl_ModelUserProfileToThrift 基准测试 ModelUserProfileToThrift 方法
func BenchmarkConverterImpl_ModelUserProfileToThrift(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	now := time.Now().UnixMilli()
	userID := uuid.New()
	accountExpiry := now + 86400000

	model := &models.UserProfile{
		BaseModel: models.BaseModel{
			ID:        userID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:           "benchmark",
		Email:              "benchmark@example.com",
		Phone:              "13800138000",
		FirstName:          "Benchmark",
		LastName:           "User",
		RealName:           "Benchmark User",
		Gender:             models.GenderMale,
		ProfessionalTitle:  "Performance Engineer",
		EmployeeID:         "EMP001",
		Status:             models.UserStatusActive,
		LoginAttempts:      1,
		MustChangePassword: false,
		AccountExpiry:      &accountExpiry,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelUserProfileToThrift(model)
	}
}

// BenchmarkConverterImpl_CreateUserRequestToModel 基准测试 CreateUserRequestToModel 方法
func BenchmarkConverterImpl_CreateUserRequestToModel(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	username := "benchmarkuser"
	password := "benchmarkpass"
	email := "benchmark@example.com"
	phone := "13800138000"
	firstName := "Benchmark"
	lastName := "User"
	realName := "Benchmark User"
	gender := core.Gender_MALE
	professionalTitle := "Performance Engineer"
	employeeID := "EMP001"
	mustChangePassword := false
	accountExpiry := time.Now().UnixMilli() + 86400000

	req := &identity_srv.CreateUserRequest{
		Username:           &username,
		Password:           &password,
		Email:              &email,
		Phone:              &phone,
		FirstName:          &firstName,
		LastName:           &lastName,
		RealName:           &realName,
		Gender:             &gender,
		ProfessionalTitle:  &professionalTitle,
		EmployeeID:         &employeeID,
		MustChangePassword: &mustChangePassword,
		AccountExpiry:      &accountExpiry,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.CreateUserRequestToModel(req)
	}
}

// BenchmarkConverterImpl_ApplyUpdateUserToModel 基准测试 ApplyUpdateUserToModel 方法
func BenchmarkConverterImpl_ApplyUpdateUserToModel(b *testing.B) {
	mockEnumConv := &mockEnumConverter{}
	converter := NewConverter(mockEnumConv)

	email := "updated@example.com"
	phone := "13800138001"
	firstName := "Updated"
	lastName := "Name"
	realName := "Updated Name"
	gender := core.Gender_FEMALE
	professionalTitle := "Senior Engineer"
	employeeID := "EMP002"
	accountExpiry := time.Now().UnixMilli() + 86400000

	req := &identity_srv.UpdateUserRequest{
		Email:             &email,
		Phone:             &phone,
		FirstName:         &firstName,
		LastName:          &lastName,
		RealName:          &realName,
		Gender:            &gender,
		ProfessionalTitle: &professionalTitle,
		EmployeeID:        &employeeID,
		AccountExpiry:     &accountExpiry,
	}

	existing := &models.UserProfile{
		Username: "existinguser",
		Status:   models.UserStatusActive,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ApplyUpdateUserToModel(existing, req)
	}
}
