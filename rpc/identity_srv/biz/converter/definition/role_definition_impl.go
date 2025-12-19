package definition

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/enum"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl implements the Converter interface.
type ConverterImpl struct{ converter enum.Converter }

// NewConverter creates a new ConverterImpl.
func NewConverter(converter enum.Converter) Converter {
	return &ConverterImpl{converter: converter}
}

// ModelToThrift converts a models.RoleDefinition to an permission_srv.RoleDefinition.
func (c *ConverterImpl) ModelToThrift(
	model *models.RoleDefinition,
) *identity_srv.RoleDefinition {
	if model == nil {
		return nil
	}

	var createdBy, updatedBy *string

	if model.CreatedBy != nil {
		s := model.CreatedBy.String()
		createdBy = &s
	}

	if model.UpdatedBy != nil {
		s := model.UpdatedBy.String()
		updatedBy = &s
	}

	status := c.converter.ModelRoleStatusToThrift(model.Status)

	// 简化实现，直接返回基本字段，权限转换暂时跳过
	return &identity_srv.RoleDefinition{
		Id:           convutil.StringPtr(model.ID.String()),
		Name:         convutil.StringPtr(model.Name),
		Description:  convutil.StringPtr(model.Description),
		Status:       &status,
		Permissions:  []*identity_srv.Permission{}, // 暂时返回空数组
		IsSystemRole: model.IsSystemRole,
		CreatedBy:    createdBy,
		UpdatedBy:    updatedBy,
		CreatedAt:    &model.CreatedAt,
		UpdatedAt:    &model.UpdatedAt,
		UserCount:    &model.UserCount, // 新增：用户数量
	}
}
