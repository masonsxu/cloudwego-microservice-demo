package user

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter 用户转换器接口，负责用户档案的模型转换
// 基于重构后的 IDL 设计，支持以下实体的双向转换：
// - UserProfile: 用户档案实体（包含个人信息、专业信息、状态管理等）
type Converter interface {
	// ============================================================================
	// UserProfile 用户档案转换
	// ============================================================================

	// 用户档案：models -> Thrift
	ModelUserProfileToThrift(*models.UserProfile) *identity_srv.UserProfile

	// ============================================================================
	// Logic层需要的转换方法
	// ============================================================================

	// 创建用户请求转换
	CreateUserRequestToModel(*identity_srv.CreateUserRequest) *models.UserProfile

	// 更新用户请求转换
	ApplyUpdateUserToModel(*models.UserProfile, *identity_srv.UpdateUserRequest) *models.UserProfile
}
