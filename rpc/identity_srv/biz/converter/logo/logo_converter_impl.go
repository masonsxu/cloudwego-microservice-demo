package logo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl 组织Logo转换器实现
type ConverterImpl struct{}

// NewConverter 创建组织Logo转换器
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ModelToThrift 将 models.OrganizationLogo 转换为 identity_srv.OrganizationLogo
func (c *ConverterImpl) ModelToThrift(
	model *models.OrganizationLogo,
) *identity_srv.OrganizationLogo {
	if model == nil {
		return nil
	}

	idStr := model.ID.String()
	status := c.StatusModelToThrift(model.Status)

	thrift := &identity_srv.OrganizationLogo{
		ID:        &idStr,
		FileID:    &model.FileID,
		Status:    &status,
		FileName:  &model.FileName,
		FileSize:  &model.FileSize,
		MimeType:  &model.MimeType,
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	}

	// BoundOrganizationID (可选)
	if model.BoundOrganizationID != nil {
		orgIDStr := model.BoundOrganizationID.String()
		thrift.BoundOrganizationID = &orgIDStr
	}

	// ExpiresAt (可选)
	if model.ExpiresAt != nil {
		thrift.ExpiresAt = model.ExpiresAt
	}

	// UploadedBy
	uploaderIDStr := model.UploadedBy.String()
	thrift.UploadedBy = &uploaderIDStr

	// DownloadUrl - 不存储在数据库中，由业务逻辑层生成
	// 这里留空，由Logic层调用RustFS生成URL后填充

	return thrift
}

// UploadRequestToModel 将上传请求转换为 Model
func (c *ConverterImpl) UploadRequestToModel(
	req *identity_srv.UploadTemporaryLogoRequest,
) *models.OrganizationLogo {
	if req == nil {
		return nil
	}

	logo := &models.OrganizationLogo{
		Status: models.LogoStatusTemporary,
	}

	// FileName
	if req.FileName != nil {
		logo.FileName = *req.FileName
	}

	// MimeType
	if req.MimeType != nil {
		logo.MimeType = *req.MimeType
	}

	// UploadedBy
	if req.UploadedBy != nil {
		uploaderID, err := uuid.Parse(*req.UploadedBy)
		if err == nil {
			logo.UploadedBy = uploaderID
		}
	}

	// FileSize
	if req.FileContent != nil {
		logo.FileSize = int64(len(req.FileContent))
	}

	// ExpiresAt - 临时Logo默认7天过期
	expiryTime := time.Now().Add(7 * 24 * time.Hour).UnixMilli()
	logo.ExpiresAt = &expiryTime

	return logo
}

// BindRequestToModel 将绑定请求转换为 Model（仅用于更新）
func (c *ConverterImpl) BindRequestToModel(
	req *identity_srv.BindLogoToOrganizationRequest,
) (*models.OrganizationLogo, error) {
	if req == nil {
		return nil, fmt.Errorf("绑定请求不能为空")
	}

	if req.LogoID == nil || *req.LogoID == "" {
		return nil, fmt.Errorf("LogoID不能为空")
	}

	if req.OrganizationID == nil || *req.OrganizationID == "" {
		return nil, fmt.Errorf("OrganizationID不能为空")
	}

	logoID, err := uuid.Parse(*req.LogoID)
	if err != nil {
		return nil, fmt.Errorf("无效的LogoID格式: %w", err)
	}

	orgID, err := uuid.Parse(*req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("无效的OrganizationID格式: %w", err)
	}

	logo := &models.OrganizationLogo{
		BaseModel: models.BaseModel{
			ID: logoID,
		},
		Status:              models.LogoStatusBound,
		BoundOrganizationID: &orgID,
		ExpiresAt:           nil, // 绑定后清除过期时间
	}

	return logo, nil
}

// StatusModelToThrift 将 Model 状态转换为 Thrift 状态
func (c *ConverterImpl) StatusModelToThrift(
	status models.OrganizationLogoStatus,
) identity_srv.OrganizationLogoStatus {
	switch status {
	case models.LogoStatusTemporary:
		return identity_srv.OrganizationLogoStatus_TEMPORARY
	case models.LogoStatusBound:
		return identity_srv.OrganizationLogoStatus_BOUND
	case models.LogoStatusDeleted:
		return identity_srv.OrganizationLogoStatus_DELETED
	default:
		return identity_srv.OrganizationLogoStatus_TEMPORARY
	}
}

// StatusThriftToModel 将 Thrift 状态转换为 Model 状态
func (c *ConverterImpl) StatusThriftToModel(
	status identity_srv.OrganizationLogoStatus,
) models.OrganizationLogoStatus {
	switch status {
	case identity_srv.OrganizationLogoStatus_TEMPORARY:
		return models.LogoStatusTemporary
	case identity_srv.OrganizationLogoStatus_BOUND:
		return models.LogoStatusBound
	case identity_srv.OrganizationLogoStatus_DELETED:
		return models.LogoStatusDeleted
	default:
		return models.LogoStatusTemporary
	}
}
