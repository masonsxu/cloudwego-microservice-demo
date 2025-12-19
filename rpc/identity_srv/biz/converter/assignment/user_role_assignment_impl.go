package assignment

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl implements the UserRoleAssignmentConverter interface.
type ConverterImpl struct{}

// NewConverter creates a new ConverterImpl.
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ModelToThrift converts a models.UserRoleAssignment to an identity_srv.UserRoleAssignment.
func (c *ConverterImpl) ModelToThrift(
	model *models.UserRoleAssignment,
) *identity_srv.UserRoleAssignment {
	if model == nil {
		return nil
	}

	// Convert UUID to string
	id := model.ID.String()
	userID := model.UserID.String()
	roleID := model.RoleID.String()

	result := &identity_srv.UserRoleAssignment{
		Id:        &id,
		UserID:    &userID,
		RoleID:    &roleID,
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	}

	// 安全处理可选的 CreatedBy 字段
	if model.CreatedBy != nil {
		createdBy := model.CreatedBy.String()
		result.CreatedBy = &createdBy
	}

	// 安全处理可选的 UpdatedBy 字段
	if model.UpdatedBy != nil {
		updatedBy := model.UpdatedBy.String()
		result.UpdatedBy = &updatedBy
	}

	return result
}
