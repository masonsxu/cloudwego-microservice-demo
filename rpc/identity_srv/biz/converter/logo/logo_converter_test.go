package logo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	assert.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
}

func TestConverterImpl_ModelToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("完整临时Logo转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		logoID := uuid.New()
		uploaderID := uuid.New()
		expiresAt := now.Add(7 * 24 * time.Hour).UnixMilli()

		model := &models.OrganizationLogo{
			BaseModel: models.BaseModel{
				ID:        logoID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Status:              models.LogoStatusTemporary,
			BoundOrganizationID: nil,
			FileID:              "organization-logos/test.png",
			FileName:            "test.png",
			FileSize:            1024000,
			MimeType:            "image/png",
			ExpiresAt:           &expiresAt,
			UploadedBy:          uploaderID,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, logoID.String(), *result.ID)
		assert.Equal(t, "organization-logos/test.png", *result.FileID)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_TEMPORARY, *result.Status)
		assert.Equal(t, "test.png", *result.FileName)
		assert.Equal(t, int64(1024000), *result.FileSize)
		assert.Equal(t, "image/png", *result.MimeType)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, uploaderID.String(), *result.UploadedBy)
		assert.Equal(t, expiresAt, *result.ExpiresAt)
		assert.Nil(t, result.BoundOrganizationID) // 临时Logo没有组织ID
	})

	t.Run("已绑定Logo转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		logoID := uuid.New()
		uploaderID := uuid.New()
		orgID := uuid.New()

		model := &models.OrganizationLogo{
			BaseModel: models.BaseModel{
				ID:        logoID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Status:              models.LogoStatusBound,
			BoundOrganizationID: &orgID,
			FileID:              "organization-logos/bound.jpg",
			FileName:            "bound.jpg",
			FileSize:            2048000,
			MimeType:            "image/jpeg",
			ExpiresAt:           nil, // 已绑定的Logo没有过期时间
			UploadedBy:          uploaderID,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, logoID.String(), *result.ID)
		assert.Equal(t, "organization-logos/bound.jpg", *result.FileID)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_BOUND, *result.Status)
		assert.Equal(t, "bound.jpg", *result.FileName)
		assert.Equal(t, int64(2048000), *result.FileSize)
		assert.Equal(t, "image/jpeg", *result.MimeType)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, uploaderID.String(), *result.UploadedBy)
		assert.Nil(t, result.ExpiresAt) // 已绑定的Logo没有过期时间
		require.NotNil(t, result.BoundOrganizationID)
		assert.Equal(t, orgID.String(), *result.BoundOrganizationID)
	})

	t.Run("已删除Logo转换", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		logoID := uuid.New()
		uploaderID := uuid.New()

		model := &models.OrganizationLogo{
			BaseModel: models.BaseModel{
				ID:        logoID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Status:              models.LogoStatusDeleted,
			BoundOrganizationID: nil,
			FileID:              "organization-logos/deleted.gif",
			FileName:            "deleted.gif",
			FileSize:            512000,
			MimeType:            "image/gif",
			ExpiresAt:           nil,
			UploadedBy:          uploaderID,
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, logoID.String(), *result.ID)
		assert.Equal(t, "organization-logos/deleted.gif", *result.FileID)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_DELETED, *result.Status)
		assert.Equal(t, "deleted.gif", *result.FileName)
		assert.Equal(t, int64(512000), *result.FileSize)
		assert.Equal(t, "image/gif", *result.MimeType)
		assert.Equal(t, nowMs, *result.CreatedAt)
		assert.Equal(t, nowMs, *result.UpdatedAt)
		assert.Equal(t, uploaderID.String(), *result.UploadedBy)
		assert.Nil(t, result.ExpiresAt)
		assert.Nil(t, result.BoundOrganizationID)
	})
}

func TestConverterImpl_UploadRequestToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.UploadRequestToModel(nil)
		assert.Nil(t, result)
	})

	t.Run("完整上传请求转换", func(t *testing.T) {
		uploaderID := uuid.New().String()
		fileName := "test-logo.png"
		mimeType := "image/png"
		fileContent := []byte("fake image content")

		req := &identity_srv.UploadTemporaryLogoRequest{
			FileContent: fileContent,
			FileName:    &fileName,
			MimeType:    &mimeType,
			UploadedBy:  &uploaderID,
		}

		result := converter.UploadRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, models.LogoStatusTemporary, result.Status)
		assert.Equal(t, fileName, result.FileName)
		assert.Equal(t, mimeType, result.MimeType)
		assert.Equal(t, int64(len(fileContent)), result.FileSize)
		assert.Equal(t, uploaderID, result.UploadedBy.String())
		require.NotNil(t, result.ExpiresAt)
		
		// 验证过期时间是大约7天后
		expectedExpiry := time.Now().Add(7 * 24 * time.Hour).UnixMilli()
		assert.InDelta(t, expectedExpiry, *result.ExpiresAt, 1000) // 允许1秒误差
	})

	t.Run("最小上传请求转换", func(t *testing.T) {
		req := &identity_srv.UploadTemporaryLogoRequest{}

		result := converter.UploadRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, models.LogoStatusTemporary, result.Status)
		assert.Equal(t, "", result.FileName) // 默认空值
		assert.Equal(t, "", result.MimeType) // 默认空值
		assert.Equal(t, int64(0), result.FileSize) // 没有文件内容
		assert.Equal(t, uuid.Nil, result.UploadedBy) // 无效UUID
		require.NotNil(t, result.ExpiresAt)
	})

	t.Run("无效上传者ID处理", func(t *testing.T) {
		invalidUploaderID := "invalid-uuid"
		req := &identity_srv.UploadTemporaryLogoRequest{
			UploadedBy: &invalidUploaderID,
		}

		result := converter.UploadRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, uuid.Nil, result.UploadedBy) // 无效UUID被设置为Nil
	})

	t.Run("空文件内容处理", func(t *testing.T) {
		emptyContent := []byte{}
		req := &identity_srv.UploadTemporaryLogoRequest{
			FileContent: emptyContent,
		}

		result := converter.UploadRequestToModel(req)

		require.NotNil(t, result)
		assert.Equal(t, int64(0), result.FileSize)
	})
}

func TestConverterImpl_BindRequestToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result, err := converter.BindRequestToModel(nil)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "绑定请求不能为空")
	})

	t.Run("空LogoID转换", func(t *testing.T) {
		req := &identity_srv.BindLogoToOrganizationRequest{
			LogoID:         nil,
			OrganizationID: func() *string { s := uuid.New().String(); return &s }(),
		}

		result, err := converter.BindRequestToModel(req)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "LogoID不能为空")
	})

	t.Run("空OrganizationID转换", func(t *testing.T) {
		req := &identity_srv.BindLogoToOrganizationRequest{
			LogoID:         func() *string { s := uuid.New().String(); return &s }(),
			OrganizationID: nil,
		}

		result, err := converter.BindRequestToModel(req)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OrganizationID不能为空")
	})

	t.Run("无效LogoID转换", func(t *testing.T) {
		req := &identity_srv.BindLogoToOrganizationRequest{
			LogoID:         func() *string { s := "invalid-uuid"; return &s }(),
			OrganizationID: func() *string { s := uuid.New().String(); return &s }(),
		}

		result, err := converter.BindRequestToModel(req)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的LogoID格式")
	})

	t.Run("无效OrganizationID转换", func(t *testing.T) {
		req := &identity_srv.BindLogoToOrganizationRequest{
			LogoID:         func() *string { s := uuid.New().String(); return &s }(),
			OrganizationID: func() *string { s := "invalid-uuid"; return &s }(),
		}

		result, err := converter.BindRequestToModel(req)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的OrganizationID格式")
	})

	t.Run("完整绑定请求转换", func(t *testing.T) {
		logoID := uuid.New()
		orgID := uuid.New()

		req := &identity_srv.BindLogoToOrganizationRequest{
			LogoID:         func() *string { s := logoID.String(); return &s }(),
			OrganizationID: func() *string { s := orgID.String(); return &s }(),
		}

		result, err := converter.BindRequestToModel(req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, logoID, result.ID)
		assert.Equal(t, models.LogoStatusBound, result.Status)
		require.NotNil(t, result.BoundOrganizationID)
		assert.Equal(t, orgID, *result.BoundOrganizationID)
		assert.Nil(t, result.ExpiresAt) // 绑定后清除过期时间
	})
}

func TestConverterImpl_StatusModelToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("临时状态转换", func(t *testing.T) {
		result := converter.StatusModelToThrift(models.LogoStatusTemporary)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_TEMPORARY, result)
	})

	t.Run("已绑定状态转换", func(t *testing.T) {
		result := converter.StatusModelToThrift(models.LogoStatusBound)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_BOUND, result)
	})

	t.Run("已删除状态转换", func(t *testing.T) {
		result := converter.StatusModelToThrift(models.LogoStatusDeleted)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_DELETED, result)
	})

	t.Run("未知状态转换", func(t *testing.T) {
		result := converter.StatusModelToThrift(models.OrganizationLogoStatus(99))
		assert.Equal(t, identity_srv.OrganizationLogoStatus_TEMPORARY, result) // 默认值
	})
}

func TestConverterImpl_StatusThriftToModel(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("临时状态转换", func(t *testing.T) {
		result := converter.StatusThriftToModel(identity_srv.OrganizationLogoStatus_TEMPORARY)
		assert.Equal(t, models.LogoStatusTemporary, result)
	})

	t.Run("已绑定状态转换", func(t *testing.T) {
		result := converter.StatusThriftToModel(identity_srv.OrganizationLogoStatus_BOUND)
		assert.Equal(t, models.LogoStatusBound, result)
	})

	t.Run("已删除状态转换", func(t *testing.T) {
		result := converter.StatusThriftToModel(identity_srv.OrganizationLogoStatus_DELETED)
		assert.Equal(t, models.LogoStatusDeleted, result)
	})

	t.Run("未知状态转换", func(t *testing.T) {
		result := converter.StatusThriftToModel(identity_srv.OrganizationLogoStatus(99))
		assert.Equal(t, models.LogoStatusTemporary, result) // 默认值
	})
}

// 状态往返转换测试
func TestConverterImpl_Status_RoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name     string
		model    models.OrganizationLogoStatus
		thrift   identity_srv.OrganizationLogoStatus
		expected models.OrganizationLogoStatus
	}{
		{
			name:     "临时状态往返转换",
			model:    models.LogoStatusTemporary,
			thrift:   identity_srv.OrganizationLogoStatus_TEMPORARY,
			expected: models.LogoStatusTemporary,
		},
		{
			name:     "已绑定状态往返转换",
			model:    models.LogoStatusBound,
			thrift:   identity_srv.OrganizationLogoStatus_BOUND,
			expected: models.LogoStatusBound,
		},
		{
			name:     "已删除状态往返转换",
			model:    models.LogoStatusDeleted,
			thrift:   identity_srv.OrganizationLogoStatus_DELETED,
			expected: models.LogoStatusDeleted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Model -> Thrift
			thriftResult := converter.StatusModelToThrift(tt.model)
			assert.Equal(t, tt.thrift, thriftResult)

			// Thrift -> Model
			modelResult := converter.StatusThriftToModel(tt.thrift)
			assert.Equal(t, tt.expected, modelResult)
		})
	}
}

// 完整往返转换测试
func TestConverterImpl_CompleteRoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("临时Logo完整往返", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		logoID := uuid.New()
		uploaderID := uuid.New()
		expiresAt := now.Add(7 * 24 * time.Hour).UnixMilli()

		originalModel := &models.OrganizationLogo{
			BaseModel: models.BaseModel{
				ID:        logoID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Status:              models.LogoStatusTemporary,
			BoundOrganizationID: nil,
			FileID:              "test.png",
			FileName:            "test.png",
			FileSize:            1024,
			MimeType:            "image/png",
			ExpiresAt:           &expiresAt,
			UploadedBy:          uploaderID,
		}

		// Model -> Thrift
		thriftLogo := converter.ModelToThrift(originalModel)
		require.NotNil(t, thriftLogo)

		// 验证关键字段
		assert.Equal(t, logoID.String(), *thriftLogo.ID)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_TEMPORARY, *thriftLogo.Status)
		assert.Equal(t, "test.png", *thriftLogo.FileName)
		assert.Equal(t, int64(1024), *thriftLogo.FileSize)
		assert.Equal(t, "image/png", *thriftLogo.MimeType)
		assert.Equal(t, uploaderID.String(), *thriftLogo.UploadedBy)
		assert.Equal(t, expiresAt, *thriftLogo.ExpiresAt)
		assert.Nil(t, thriftLogo.BoundOrganizationID)
	})

	t.Run("已绑定Logo完整往返", func(t *testing.T) {
		now := time.Now()
		nowMs := now.UnixMilli()
		logoID := uuid.New()
		uploaderID := uuid.New()
		orgID := uuid.New()

		originalModel := &models.OrganizationLogo{
			BaseModel: models.BaseModel{
				ID:        logoID,
				CreatedAt: nowMs,
				UpdatedAt: nowMs,
			},
			Status:              models.LogoStatusBound,
			BoundOrganizationID: &orgID,
			FileID:              "bound.jpg",
			FileName:            "bound.jpg",
			FileSize:            2048,
			MimeType:            "image/jpeg",
			ExpiresAt:           nil,
			UploadedBy:          uploaderID,
		}

		// Model -> Thrift
		thriftLogo := converter.ModelToThrift(originalModel)
		require.NotNil(t, thriftLogo)

		// 验证关键字段
		assert.Equal(t, logoID.String(), *thriftLogo.ID)
		assert.Equal(t, identity_srv.OrganizationLogoStatus_BOUND, *thriftLogo.Status)
		assert.Equal(t, "bound.jpg", *thriftLogo.FileName)
		assert.Equal(t, int64(2048), *thriftLogo.FileSize)
		assert.Equal(t, "image/jpeg", *thriftLogo.MimeType)
		assert.Equal(t, uploaderID.String(), *thriftLogo.UploadedBy)
		assert.Nil(t, thriftLogo.ExpiresAt)
		require.NotNil(t, thriftLogo.BoundOrganizationID)
		assert.Equal(t, orgID.String(), *thriftLogo.BoundOrganizationID)
	})
}

// 边界情况和错误处理测试
func TestConverterImpl_EdgeCases(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("空字符串ID处理", func(t *testing.T) {
		req := &identity_srv.BindLogoToOrganizationRequest{
			LogoID:         func() *string { s := ""; return &s }(),
			OrganizationID: func() *string { s := uuid.New().String(); return &s }(),
		}

		result, err := converter.BindRequestToModel(req)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "LogoID不能为空")
	})

	t.Run("状态默认值处理", func(t *testing.T) {
		thriftResult := converter.StatusModelToThrift(models.OrganizationLogoStatus(-1))
		assert.Equal(t, identity_srv.OrganizationLogoStatus_TEMPORARY, thriftResult)

		modelResult := converter.StatusThriftToModel(identity_srv.OrganizationLogoStatus(-1))
		assert.Equal(t, models.LogoStatusTemporary, modelResult)
	})
}

// 基准测试
func BenchmarkConverterImpl_ModelToThrift(b *testing.B) {
	converter := &ConverterImpl{}
	
	now := time.Now()
	nowMs := now.UnixMilli()
	logoID := uuid.New()
	uploaderID := uuid.New()

	model := &models.OrganizationLogo{
		BaseModel: models.BaseModel{
			ID:        logoID,
			CreatedAt: nowMs,
			UpdatedAt: nowMs,
		},
		Status:     models.LogoStatusTemporary,
		FileID:     "test.png",
		FileName:   "test.png",
		FileSize:   1024,
		MimeType:   "image/png",
		UploadedBy: uploaderID,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.ModelToThrift(model)
	}
}

func BenchmarkConverterImpl_UploadRequestToModel(b *testing.B) {
	converter := &ConverterImpl{}
	
	uploaderID := uuid.New().String()
	fileName := "test.png"
	mimeType := "image/png"
	fileContent := make([]byte, 1024) // 1KB

	req := &identity_srv.UploadTemporaryLogoRequest{
		FileContent: fileContent,
		FileName:    &fileName,
		MimeType:    &mimeType,
		UploadedBy:  &uploaderID,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.UploadRequestToModel(req)
	}
}

func BenchmarkConverterImpl_BindRequestToModel(b *testing.B) {
	converter := &ConverterImpl{}
	
	logoID := uuid.New().String()
	orgID := uuid.New().String()

	req := &identity_srv.BindLogoToOrganizationRequest{
		LogoID:         &logoID,
		OrganizationID: &orgID,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = converter.BindRequestToModel(req)
	}
}

func BenchmarkConverterImpl_StatusConversion(b *testing.B) {
	converter := &ConverterImpl{}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Model -> Thrift -> Model 往环
		thriftStatus := converter.StatusModelToThrift(models.LogoStatusTemporary)
		_ = converter.StatusThriftToModel(thriftStatus)
	}
}