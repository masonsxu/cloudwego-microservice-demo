package enum

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter is an interface for mapping enums between the models and Thrift layers.
type Converter interface {
	// UserStatus
	ModelUserStatusToThrift(models.UserStatus) core.UserStatus
	ThriftUserStatusToModel(core.UserStatus) models.UserStatus

	// RoleStatus
	ModelRoleStatusToThrift(models.RoleStatus) core.RoleStatus
	ThriftRoleStatusToModel(core.RoleStatus) models.RoleStatus

	// Gender
	ModelGenderToThrift(models.Gender) core.Gender
	ThriftGenderToModel(core.Gender) models.Gender
}
