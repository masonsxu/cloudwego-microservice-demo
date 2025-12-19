package authentication

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter 认证逻辑转换
type Converter interface {
	// ============================================================================
	// 成员关系转换
	// ============================================================================

	// 模型转换为 Thrift DTO
	ModelUserMembershipsToThrift([]*models.UserMembership) []*identity_srv.UserMembership
}
