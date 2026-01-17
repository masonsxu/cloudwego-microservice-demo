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

	// 转换数据范围
	dataScope := c.converter.ModelDataScopeToThrift(model.DefaultScope)

	// 转换父角色ID和部门ID
	var parentRoleID, departmentID *string
	if model.ParentRoleID != nil {
		s := model.ParentRoleID.String()
		parentRoleID = &s
	}
	if model.DepartmentID != nil {
		s := model.DepartmentID.String()
		departmentID = &s
	}

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
		UserCount:    &model.UserCount,
		// Casbin 扩展字段
		RoleCode:     convutil.StringPtr(model.RoleCode),
		ParentRoleID: parentRoleID,
		DepartmentID: departmentID,
		DefaultScope: &dataScope,
	}
}

// ModelsToThrift converts a slice of models.RoleDefinition to a slice of identity_srv.RoleDefinition.
func (c *ConverterImpl) ModelsToThrift(
	models []*models.RoleDefinition,
) []*identity_srv.RoleDefinition {
	if models == nil {
		return nil
	}

	result := make([]*identity_srv.RoleDefinition, 0, len(models))
	for _, model := range models {
		if thrift := c.ModelToThrift(model); thrift != nil {
			result = append(result, thrift)
		}
	}
	return result
}
