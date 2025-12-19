package common

import (
	identityRpcEnums "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
)

// RoleStatus 转换函数

// ConvertRoleStatusToHTTP 将 RPC 层的 RoleStatus 转换为 HTTP 层的 i32
func ConvertRoleStatusToHTTP(status identityRpcEnums.RoleStatus) int32 {
	return int32(status)
}

// ConvertRoleStatusPtrToHTTPPtr 将 RPC 层的 RoleStatus 指针转换为 HTTP 层的 i32 指针
func ConvertRoleStatusPtrToHTTPPtr(status *identityRpcEnums.RoleStatus) *int32 {
	if status == nil {
		return nil
	}

	result := int32(*status)

	return &result
}

// ConvertRoleStatusToRPC 将 HTTP 层的 i32 转换为 RPC 层的 RoleStatus
func ConvertRoleStatusToRPC(status int32) identityRpcEnums.RoleStatus {
	return identityRpcEnums.RoleStatus(status)
}

// ConvertRoleStatusToRPCPtr 将 HTTP 层的 i32 转换为 RPC 层的 RoleStatus 指针
// 注意：调用方需要确保 status 不为 nil，对于可能为 nil 的指针，请使用 ConvertRoleStatusPtrToRPCPtr
func ConvertRoleStatusToRPCPtr(status int32) *identityRpcEnums.RoleStatus {
	result := identityRpcEnums.RoleStatus(status)
	return &result
}

// ConvertRoleStatusPtrToRPCPtr 将 HTTP 层的 i32 指针转换为 RPC 层的 RoleStatus 指针
func ConvertRoleStatusPtrToRPCPtr(status *int32) *identityRpcEnums.RoleStatus {
	if status == nil {
		return nil
	}

	result := identityRpcEnums.RoleStatus(*status)

	return &result
}

// IdentityUserStatus 转换函数

// ConvertIdentityUserStatusToHTTP 将 Identity RPC 层的 UserStatus 转换为 HTTP 层的 i32
func ConvertIdentityUserStatusToHTTP(status identityRpcEnums.UserStatus) int32 {
	return int32(status)
}

// ConvertIdentityUserStatusPtrToHTTP 将 Identity RPC 层的 UserStatus 指针转换为 HTTP 层的 i32 指针
func ConvertIdentityUserStatusPtrToHTTP(status *identityRpcEnums.UserStatus) *int32 {
	if status == nil {
		return nil
	}

	result := int32(*status)

	return &result
}

// ConvertIdentityUserStatusToRPC 将 HTTP 层的 i32 转换为 RPC 层的 UserStatus
func ConvertIdentityUserStatusToRPC(status int32) identityRpcEnums.UserStatus {
	return identityRpcEnums.UserStatus(status)
}

// ConvertIdentityUserStatusToRPCPtr 将 HTTP 层的 i32 转换为 RPC 层的 UserStatus 指针
// 注意：调用方需要确保 status 不为 nil，对于可能为 nil 的指针，请使用 ConvertIdentityUserStatusPtrToRPCPtr
func ConvertIdentityUserStatusToRPCPtr(status int32) *identityRpcEnums.UserStatus {
	result := identityRpcEnums.UserStatus(status)
	return &result
}

// ConvertIdentityUserStatusPtrToRPCPtr 将 HTTP 层的 i32 指针转换为 RPC 层的 UserStatus 指针
func ConvertIdentityUserStatusPtrToRPCPtr(status *int32) *identityRpcEnums.UserStatus {
	if status == nil {
		return nil
	}

	result := identityRpcEnums.UserStatus(*status)

	return &result
}

// PermissionUserStatus 转换函数
// ConvertPermissionUserStatusToHTTP 将 RPC 层的 UserStatus 转换为 HTTP 层的 i32
func ConvertPermissionUserStatusToHTTP(status identityRpcEnums.UserStatus) int32 {
	return int32(status)
}

// ConvertPermissionUserStatusToRPC 将 HTTP 层的 i32 转换为 RPC 层的 UserStatus
func ConvertPermissionUserStatusToRPC(status int32) identityRpcEnums.UserStatus {
	return identityRpcEnums.UserStatus(status)
}

// ConvertPermissionUserStatusToRPCPtr 将 HTTP 层的 i32 转换为 RPC 层的 UserStatus 指针
// 注意：调用方需要确保 status 不为 nil，对于可能为 nil 的指针，请使用 ConvertPermissionUserStatusPtrToRPCPtr
func ConvertPermissionUserStatusToRPCPtr(status int32) *identityRpcEnums.UserStatus {
	result := identityRpcEnums.UserStatus(status)
	return &result
}

// ConvertPermissionUserStatusPtrToRPCPtr 将 HTTP 层的 i32 指针转换为 RPC 层的 UserStatus 指针
func ConvertPermissionUserStatusPtrToRPCPtr(status *int32) *identityRpcEnums.UserStatus {
	if status == nil {
		return nil
	}

	result := identityRpcEnums.UserStatus(*status)

	return &result
}

// Gender 转换函数

// ConvertGenderToHTTP 将 RPC 层的 Gender 转换为 HTTP 层的 i32
func ConvertGenderToHTTP(gender identityRpcEnums.Gender) int32 {
	return int32(gender)
}

// ConvertGenderPtrToHTTP 将 RPC 层的 Gender 指针转换为 HTTP 层的 i32 指针
func ConvertGenderPtrToHTTP(gender *identityRpcEnums.Gender) *int32 {
	if gender == nil {
		return nil
	}

	result := int32(*gender)

	return &result
}

// ConvertGenderToRPC 将 HTTP 层的 i32 转换为 RPC 层的 Gender
func ConvertGenderToRPC(gender int32) identityRpcEnums.Gender {
	return identityRpcEnums.Gender(gender)
}

// ConvertGenderToRPCPtr 将 HTTP 层的 i32 转换为 RPC 层的 Gender 指针
func ConvertGenderToRPCPtr(gender int32) *identityRpcEnums.Gender {
	result := identityRpcEnums.Gender(gender)
	return &result
}

// ConvertGenderPtrToRPCPtr 将 HTTP 层的 i32 指针转换为 RPC 层的 Gender 指针
func ConvertGenderPtrToRPCPtr(gender *int32) *identityRpcEnums.Gender {
	if gender == nil {
		return nil
	}

	result := identityRpcEnums.Gender(*gender)

	return &result
}
