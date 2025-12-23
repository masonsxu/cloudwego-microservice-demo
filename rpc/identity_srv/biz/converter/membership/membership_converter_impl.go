package membership

import (
	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/enum"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl 成员关系转换器实现
type ConverterImpl struct {
	converter enum.Converter
}

// NewConverter 创建一个新的成员关系转换器实例
func NewConverter(converter enum.Converter) Converter {
	return &ConverterImpl{
		converter: converter,
	}
}

// ============================================================================
// UserMembership 成员关系转换
// ============================================================================

// ModelUserMembershipToThrift 将 models.UserMembership 转换为 identity_srv.UserMembership
func (c *ConverterImpl) ModelUserMembershipToThrift(
	model *models.UserMembership,
) *identity_srv.UserMembership {
	if model == nil {
		return nil
	}

	id := model.ID.String()
	userID := model.UserID.String()
	orgID := model.OrganizationID.String()

	dto := &identity_srv.UserMembership{
		ID:             &id,
		UserID:         &userID,
		OrganizationID: &orgID,
		IsPrimary:      convutil.BoolPtr(model.IsPrimary),
		CreatedAt:      &model.CreatedAt,
		UpdatedAt:      &model.UpdatedAt,
	}

	// 处理可选的部门ID
	if model.DepartmentID != uuid.Nil {
		dto.DepartmentID = convutil.StringPtr(model.DepartmentID.String())
	}

	return dto
}

// ThriftUserMembershipToModel 将 identity_srv.UserMembership 转换为 models.UserMembership
func (c *ConverterImpl) ThriftUserMembershipToModel(
	dto *identity_srv.UserMembership,
) *models.UserMembership {
	if dto == nil {
		return nil
	}

	id := uuid.MustParse(*dto.ID)
	userID := uuid.MustParse(*dto.UserID)
	orgID := uuid.MustParse(*dto.OrganizationID)

	model := &models.UserMembership{
		BaseModel: models.BaseModel{
			ID: id,
		},
		UserID:         userID,
		OrganizationID: orgID,
		IsPrimary:      convutil.BoolValue(dto.IsPrimary),
	}

	// 处理可选的部门ID
	if dto.DepartmentID != nil {
		model.DepartmentID = uuid.MustParse(*dto.DepartmentID)
	}

	return model
}

// ModelUserMembershipsToThrift 将 models.UserMembership 切片转换为 identity_srv.UserMembership 切片
func (c *ConverterImpl) ModelUserMembershipsToThrift(
	models []*models.UserMembership,
) []*identity_srv.UserMembership {
	if len(models) == 0 {
		return nil
	}

	dtos := make([]*identity_srv.UserMembership, 0, len(models))
	for _, model := range models {
		if dto := c.ModelUserMembershipToThrift(model); dto != nil {
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
	membership *models.UserMembership,
) *identity_srv.UserMembership {
	return c.ModelUserMembershipToThrift(membership)
}

// AddMembershipRequestToModel 转换添加成员关系请求到模型
func (c *ConverterImpl) AddMembershipRequestToModel(
	req *identity_srv.AddMembershipRequest,
) *models.UserMembership {
	if req == nil {
		return nil
	}

	userID := uuid.MustParse(*req.UserID)
	orgID := uuid.MustParse(*req.OrganizationID)

	membership := &models.UserMembership{
		UserID:         userID,
		OrganizationID: orgID,
		Status:         models.MembershipStatusActive, // 默认为活跃状态
	}

	// 处理可选字段
	if req.DepartmentID != nil {
		membership.DepartmentID = uuid.MustParse(*req.DepartmentID)
	}

	membership.IsPrimary = req.IsPrimary

	return membership
}

// ApplyUpdateToModel 应用更新请求到现有模型
func (c *ConverterImpl) ApplyUpdateToModel(
	existing *models.UserMembership,
	req *identity_srv.UpdateMembershipRequest,
) *models.UserMembership {
	if existing == nil || req == nil {
		return existing
	}

	// 创建副本以避免修改原始模型
	updated := *existing

	// 应用更新字段
	if req.OrganizationID != nil {
		updated.OrganizationID = uuid.MustParse(*req.OrganizationID)
	}

	if req.DepartmentID != nil {
		updated.DepartmentID = uuid.MustParse(*req.DepartmentID)
	}

	if req.IsPrimary != nil {
		updated.IsPrimary = *req.IsPrimary
	}

	return &updated
}
