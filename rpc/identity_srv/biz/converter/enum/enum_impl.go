package enum

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl implements the EnumMapper interface.
type ConverterImpl struct{}

// NewConverter creates a new Converter.
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ModelUserStatusToThrift converts models.UserStatus to core.UserStatus.
func (m *ConverterImpl) ModelUserStatusToThrift(status models.UserStatus) core.UserStatus {
	switch status {
	case models.UserStatusActive:
		return core.UserStatus_ACTIVE
	case models.UserStatusInactive:
		return core.UserStatus_INACTIVE
	case models.UserStatusSuspended:
		return core.UserStatus_SUSPENDED
	case models.UserStatusLocked:
		return core.UserStatus_LOCKED
	default:
		return core.UserStatus_INACTIVE
	}
}

// ThriftUserStatusToModel converts core.UserStatus to models.UserStatus.
func (m *ConverterImpl) ThriftUserStatusToModel(status core.UserStatus) models.UserStatus {
	switch status {
	case core.UserStatus_ACTIVE:
		return models.UserStatusActive
	case core.UserStatus_INACTIVE:
		return models.UserStatusInactive
	case core.UserStatus_SUSPENDED:
		return models.UserStatusSuspended
	case core.UserStatus_LOCKED:
		return models.UserStatusLocked
	default:
		return models.UserStatusInactive
	}
}

// ModelRoleStatusToThrift converts models.RoleStatus to core.RoleStatus.
func (m *ConverterImpl) ModelRoleStatusToThrift(status models.RoleStatus) core.RoleStatus {
	switch status {
	case models.RoleStatusActive:
		return core.RoleStatus_ACTIVE
	case models.RoleStatusInactive:
		return core.RoleStatus_INACTIVE
	case models.RoleStatusDeprecated:
		return core.RoleStatus_DEPRECATED
	default:
		return core.RoleStatus_INACTIVE
	}
}

// ThriftRoleStatusToModel converts core.RoleStatus to models.RoleStatus.
func (m *ConverterImpl) ThriftRoleStatusToModel(status core.RoleStatus) models.RoleStatus {
	switch status {
	case core.RoleStatus_ACTIVE:
		return models.RoleStatusActive
	case core.RoleStatus_INACTIVE:
		return models.RoleStatusInactive
	case core.RoleStatus_DEPRECATED:
		return models.RoleStatusDeprecated
	default:
		return models.RoleStatusInactive
	}
}

// ModelGenderToThrift converts models.Gender to core.Gender.
func (m *ConverterImpl) ModelGenderToThrift(gender models.Gender) core.Gender {
	switch gender {
	case models.GenderMale:
		return core.Gender_MALE
	case models.GenderFemale:
		return core.Gender_FEMALE
	case models.GenderUnknown:
		return core.Gender_UNKNOWN
	default:
		return core.Gender_UNKNOWN
	}
}

// ThriftGenderToModel converts core.Gender to models.Gender.
func (m *ConverterImpl) ThriftGenderToModel(gender core.Gender) models.Gender {
	switch gender {
	case core.Gender_MALE:
		return models.GenderMale
	case core.Gender_FEMALE:
		return models.GenderFemale
	case core.Gender_UNKNOWN:
		return models.GenderUnknown
	default:
		return models.GenderUnknown
	}
}
