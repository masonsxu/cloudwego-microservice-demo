package organization

import (
	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl 组织转换器实现
type ConverterImpl struct{}

// NewConverter 创建一个新的组织转换器实例
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ============================================================================
// Organization 组织转换
// ============================================================================

// ModelOrganizationToThrift 将 models.Organization 转换为 identity_srv.Organization
func (c *ConverterImpl) ModelOrganizationToThrift(
	model *models.Organization,
) *identity_srv.Organization {
	if model == nil {
		return nil
	}

	id := model.ID.String()
	name := model.Name

	dto := &identity_srv.Organization{
		ID:        &id,
		Name:      &name,
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	}

	// 处理可选的组织代码
	if model.Code != "" {
		dto.Code = &model.Code
	}

	// 处理可选的父组织ID
	if model.ParentID != uuid.Nil {
		dto.ParentID = convutil.StringPtr(model.ParentID.String())
	}

	// 处理可选的机构类型
	if model.FacilityType != "" {
		dto.FacilityType = convutil.StringPtr(model.FacilityType)
	}

	// 处理可选的认证状态
	if model.AccreditationStatus != "" {
		dto.AccreditationStatus = convutil.StringPtr(model.AccreditationStatus)
	}

	// 处理可选的省市信息
	if len(model.ProvinceCity) > 0 {
		dto.ProvinceCity = []string(model.ProvinceCity)
	}

	return dto
}

// ThriftOrganizationToModel 将 identity_srv.Organization 转换为 models.Organization
func (c *ConverterImpl) ThriftOrganizationToModel(
	dto *identity_srv.Organization,
) *models.Organization {
	if dto == nil {
		return nil
	}

	id := uuid.MustParse(*dto.ID)
	name := *dto.Name

	model := &models.Organization{
		BaseModel: models.BaseModel{
			ID: id,
		},
		Code:                convutil.StringValue(dto.Code),
		Name:                name,
		FacilityType:        convutil.StringValue(dto.FacilityType),
		AccreditationStatus: convutil.StringValue(dto.AccreditationStatus),
		ProvinceCity:        models.StringSlice(dto.ProvinceCity),
	}

	// 处理可选的父组织ID
	if dto.ParentID != nil {
		model.ParentID = uuid.MustParse(*dto.ParentID)
	}

	return model
}

// ModelOrganizationsToThrift 将 models.Organization 切片转换为 identity_srv.Organization 切片
func (c *ConverterImpl) ModelOrganizationsToThrift(
	models []*models.Organization,
) []*identity_srv.Organization {
	if len(models) == 0 {
		return nil
	}

	dtos := make([]*identity_srv.Organization, 0, len(models))
	for _, model := range models {
		if dto := c.ModelOrganizationToThrift(model); dto != nil {
			dtos = append(dtos, dto)
		}
	}

	return dtos
}

// ============================================================================
// Logic层需要的转换方法
// ============================================================================

// ModelToThrift Logic层使用的转换方法（别名）
func (c *ConverterImpl) ModelToThrift(
	org *models.Organization,
) *identity_srv.Organization {
	return c.ModelOrganizationToThrift(org)
}

// CreateRequestToModel 转换创建组织请求到模型
func (c *ConverterImpl) CreateRequestToModel(
	req *identity_srv.CreateOrganizationRequest,
) *models.Organization {
	if req == nil {
		return nil
	}

	name := *req.Name

	model := &models.Organization{
		Name: name,
		// Code将在repository层的BeforeCreate中自动生成或基于name生成
		Code: "", // 后续可以基于name生成唯一code
	}

	// 处理可选字段
	if req.ParentID != nil {
		model.ParentID = uuid.MustParse(*req.ParentID)
	}

	if req.FacilityType != nil {
		model.FacilityType = *req.FacilityType
	}

	if req.AccreditationStatus != nil {
		model.AccreditationStatus = *req.AccreditationStatus
	}

	if req.ProvinceCity != nil {
		model.ProvinceCity = models.StringSlice(req.ProvinceCity)
	}

	return model
}

// ApplyUpdateToModel 应用更新请求到现有模型
func (c *ConverterImpl) ApplyUpdateToModel(
	existing *models.Organization,
	req *identity_srv.UpdateOrganizationRequest,
) *models.Organization {
	if existing == nil || req == nil {
		return existing
	}

	// 创建副本以避免修改原始模型
	updated := *existing

	// 应用更新字段
	if req.Name != nil {
		updated.Name = *req.Name
	}

	if req.ParentID != nil {
		updated.ParentID = uuid.MustParse(*req.ParentID)
	}

	if req.FacilityType != nil {
		updated.FacilityType = *req.FacilityType
	}

	if req.AccreditationStatus != nil {
		updated.AccreditationStatus = *req.AccreditationStatus
	}

	if req.ProvinceCity != nil {
		updated.ProvinceCity = models.StringSlice(req.ProvinceCity)
	}

	return &updated
}
