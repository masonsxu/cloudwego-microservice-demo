package authentication

import (
	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl 认证逻辑转换实现
type ConverterImpl struct{}

// NewConverter 创建认证逻辑转换实现
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ============================================================================
// 成员关系转换
// ============================================================================

// ModelUserMembershipsToThrift 转换用户成员关系列表为 Thrift DTO
func (c *ConverterImpl) ModelUserMembershipsToThrift(
	memberships []*models.UserMembership,
) []*identity_srv.UserMembership {
	if len(memberships) == 0 {
		return []*identity_srv.UserMembership{}
	}

	result := make([]*identity_srv.UserMembership, len(memberships))
	for i, membership := range memberships {
		if membership != nil {
			id := membership.ID.String()
			userID := membership.UserID.String()
			orgID := membership.OrganizationID.String()

			// 转换基础字段
			thriftMembership := &identity_srv.UserMembership{
				ID:             &id,
				UserID:         &userID,
				OrganizationID: &orgID,
				CreatedAt:      &membership.CreatedAt,
				UpdatedAt:      &membership.UpdatedAt,
			}

			// 可选字段处理
			if membership.DepartmentID != uuid.Nil {
				deptID := membership.DepartmentID.String()
				thriftMembership.DepartmentID = &deptID
			}

			if membership.IsPrimary {
				thriftMembership.IsPrimary = &membership.IsPrimary
			}

			result[i] = thriftMembership
		}
	}

	return result
}
