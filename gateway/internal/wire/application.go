// Package wire 应用层依赖注入
package wire

import (
	"github.com/google/wire"
	identityassembler "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	permissionassembler "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/permission"
)

// ApplicationSet 应用层依赖注入集合
// 按照分层架构组织：子 assembler -> 聚合 assembler
var ApplicationSet = wire.NewSet(
	// 各业务领域 assembler
	identityassembler.NewAuthAssembler,
	identityassembler.NewUserAssembler,
	identityassembler.NewOrgAssembler,
	identityassembler.NewDepartmentAssembler,
	identityassembler.NewMembershipAssembler,
	identityassembler.NewLogoAssembler,

	// 权限相关 assembler
	permissionassembler.NewPermissionAssembler,
	permissionassembler.NewRoleAssembler,
	permissionassembler.NewUserRoleAssembler,
	permissionassembler.NewMenuAssembler,

	// 聚合 assembler
	identityassembler.NewIdentityAggregateAssembler,
	permissionassembler.NewPermissionAggregateAssembler,
)
