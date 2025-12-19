package models

// UserStatus 用户状态枚举
type UserStatus int32

const (
	UserStatusActive    UserStatus = 1 // 活跃状态​​：用户账号正常可用，拥有所有权限，可以登录和操作。
	UserStatusInactive  UserStatus = 2 // 未激活状态​​：需要由管理员进行激活
	UserStatusSuspended UserStatus = 3 // 暂停状态​​：离职
	UserStatusLocked    UserStatus = 4 // 锁定状态​​：用户账号被锁定，无法登录和操作。
)

// RoleStatus 角色定义状态枚举
type RoleStatus int32

const (
	RoleStatusActive     RoleStatus = 1 // 活跃状态​​：角色正常可用。
	RoleStatusInactive   RoleStatus = 2 // 未激活状态​​：需要由管理员进行激活
	RoleStatusDeprecated RoleStatus = 3 // 已弃用
)

// Gender 性别枚举
type Gender int32

const (
	GenderUnknown Gender = 0 // 未知
	GenderMale    Gender = 1 // 男性
	GenderFemale  Gender = 2 // 女性
)

// MembershipStatus 成员关系状态枚举
type MembershipStatus int32

const (
	MembershipStatusActive    MembershipStatus = 1 // 活跃的成员关系
	MembershipStatusPending   MembershipStatus = 2 // 待处理/待接受
	MembershipStatusSuspended MembershipStatus = 3 // 暂停状态​​：用户账号被暂停，无法登录和操作。
	MembershipStatusEnded     MembershipStatus = 4 // 已结束 (例如，调离、项目结束)
)
