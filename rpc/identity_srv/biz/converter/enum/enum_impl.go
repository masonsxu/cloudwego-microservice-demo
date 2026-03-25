package enum

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
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
		return core.UserStatus_USER_STATUS_ACTIVE
	case models.UserStatusInactive:
		return core.UserStatus_USER_STATUS_INACTIVE
	case models.UserStatusSuspended:
		return core.UserStatus_USER_STATUS_SUSPENDED
	case models.UserStatusLocked:
		return core.UserStatus_USER_STATUS_LOCKED
	default:
		return core.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

// ThriftUserStatusToModel converts core.UserStatus to models.UserStatus.
func (m *ConverterImpl) ThriftUserStatusToModel(status core.UserStatus) models.UserStatus {
	switch status {
	case core.UserStatus_USER_STATUS_ACTIVE:
		return models.UserStatusActive
	case core.UserStatus_USER_STATUS_INACTIVE:
		return models.UserStatusInactive
	case core.UserStatus_USER_STATUS_SUSPENDED:
		return models.UserStatusSuspended
	case core.UserStatus_USER_STATUS_LOCKED:
		return models.UserStatusLocked
	default:
		return models.UserStatusInactive
	}
}

// ModelRoleStatusToThrift converts models.RoleStatus to core.RoleStatus.
func (m *ConverterImpl) ModelRoleStatusToThrift(status models.RoleStatus) core.RoleStatus {
	switch status {
	case models.RoleStatusActive:
		return core.RoleStatus_ROLE_STATUS_ACTIVE
	case models.RoleStatusInactive:
		return core.RoleStatus_ROLE_STATUS_INACTIVE
	case models.RoleStatusDeprecated:
		return core.RoleStatus_ROLE_STATUS_DEPRECATED
	default:
		return core.RoleStatus_ROLE_STATUS_UNSPECIFIED
	}
}

// ThriftRoleStatusToModel converts core.RoleStatus to models.RoleStatus.
func (m *ConverterImpl) ThriftRoleStatusToModel(status core.RoleStatus) models.RoleStatus {
	switch status {
	case core.RoleStatus_ROLE_STATUS_ACTIVE:
		return models.RoleStatusActive
	case core.RoleStatus_ROLE_STATUS_INACTIVE:
		return models.RoleStatusInactive
	case core.RoleStatus_ROLE_STATUS_DEPRECATED:
		return models.RoleStatusDeprecated
	default:
		return models.RoleStatusInactive
	}
}

// ModelGenderToThrift converts models.Gender to core.Gender.
func (m *ConverterImpl) ModelGenderToThrift(gender models.Gender) core.Gender {
	switch gender {
	case models.GenderMale:
		return core.Gender_GENDER_MALE
	case models.GenderFemale:
		return core.Gender_GENDER_FEMALE
	case models.GenderUnknown:
		return core.Gender_GENDER_UNSPECIFIED
	default:
		return core.Gender_GENDER_UNSPECIFIED
	}
}

// ThriftGenderToModel converts core.Gender to models.Gender.
func (m *ConverterImpl) ThriftGenderToModel(gender core.Gender) models.Gender {
	switch gender {
	case core.Gender_GENDER_MALE:
		return models.GenderMale
	case core.Gender_GENDER_FEMALE:
		return models.GenderFemale
	case core.Gender_GENDER_UNSPECIFIED:
		return models.GenderUnknown
	default:
		return models.GenderUnknown
	}
}

// ModelDataScopeToThrift converts models.DataScopeType to identity_srv.DataScope.
func (m *ConverterImpl) ModelDataScopeToThrift(scope models.DataScopeType) identity_srv.DataScope {
	switch scope {
	case models.DataScopeSelf:
		return identity_srv.DataScope_DATA_SCOPE_SELF
	case models.DataScopeDept:
		return identity_srv.DataScope_DATA_SCOPE_DEPT
	case models.DataScopeOrg:
		return identity_srv.DataScope_DATA_SCOPE_ORG
	default:
		return identity_srv.DataScope_DATA_SCOPE_UNSPECIFIED
	}
}

// ThriftDataScopeToModel converts identity_srv.DataScope to models.DataScopeType.
func (m *ConverterImpl) ThriftDataScopeToModel(scope identity_srv.DataScope) models.DataScopeType {
	switch scope {
	case identity_srv.DataScope_DATA_SCOPE_SELF:
		return models.DataScopeSelf
	case identity_srv.DataScope_DATA_SCOPE_DEPT:
		return models.DataScopeDept
	case identity_srv.DataScope_DATA_SCOPE_ORG:
		return models.DataScopeOrg
	default:
		return models.DataScopeSelf
	}
}
