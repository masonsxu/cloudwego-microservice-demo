package errno

// =================================================================
//
//	错误码规范
//
// =================================================================
// 错误码采用6位数字的结构，格式为：A-BB-CCC
// A: 错误级别 (1位)
//   - 1: 系统级错误
//   - 2: 业务级错误
//
// BB: 服务/模块编码 (2位)
//   - 00: 通用业务错误
//   - 01: 用户领域
//   - 02: 组织领域
//   - 03: 部门领域
//   - 04: 级联删除领域
//   - 05: 姿态资源领域
//   - 06: 组织Logo领域
//   - 07: 角色分配领域
//
// CCC: 具体错误编码 (3位)
// 错误码定义
const (
	// 通用业务错误 (200xxx)
	ErrorCodeInvalidParams   = 200100 // 参数错误
	ErrorCodeOperationFailed = 200101 // 操作失败（通用）

	// 用户相关错误 (201xxx)
	ErrorCodeUserNotFound           = 201001
	ErrorCodeUserAlreadyExists      = 201002
	ErrorCodeUsernameAlreadyExists  = 201003
	ErrorCodeEmailAlreadyExists     = 201004
	ErrorCodePhoneAlreadyExists     = 201005
	ErrorCodeInvalidPassword        = 201007
	ErrorCodeUserInactive           = 201012
	ErrorCodePhoneNumberAlreadyUsed = 201014
	ErrorCodeInvalidAccountType     = 201015
	ErrorCodeInvalidCredentials     = 201016
	ErrorCodeUserSuspended          = 201017
	ErrorCodeMustChangePassword     = 201018
	ErrorCodeSystemUserCannotDelete    = 201019 // 系统用户无法删除
	ErrorCodeSystemUserCannotModifyKey = 201020 // 系统用户关键属性无法修改

	// 组织相关错误 (202xxx)
	ErrorCodeOrganizationNotFound                    = 202001
	ErrorCodeOrganizationAlreadyExists               = 202002
	ErrorCodeOrganizationHasUsers                    = 202003
	ErrorCodeCannotDeleteOrganizationWithDepartments = 202005
	ErrorCodeParentOrganizationNotFound              = 202006

	// 部门相关错误 (203xxx)
	ErrorCodeDepartmentNotFound                = 203001
	ErrorCodeDepartmentAlreadyExists           = 203002
	ErrorCodeDepartmentCodeAlreadyExists       = 203003
	ErrorCodeCannotDeleteDepartmentWithMembers = 203005
	ErrorCodeDepartmentNameRequired            = 203006
	ErrorCodeDepartmentOrganizationRequired    = 203007

	// 级联删除和数据一致性相关错误 (204xxx)
	ErrorCodeUserNotInSameOrganization = 204002
	ErrorCodeDataInconsistency         = 204003
	ErrorCodeTransactionFailed         = 204004
	ErrorCodeMembershipNotFound        = 204005
	ErrorCodeMembershipAlreadyExists   = 204006

	// 组织Logo相关错误 (206xxx)
	ErrorCodeLogoNotFound      = 206001
	ErrorCodeLogoAlreadyBound  = 206002
	ErrorCodeLogoExpired       = 206003
	ErrorCodeLogoInvalidStatus = 206004
	ErrorCodeLogoBindingFailed = 206005
	ErrorCodeLogoAlreadyExists = 206006
	ErrorCodeInvalidFileType   = 206007
	ErrorCodeFileSizeExceeded  = 206008
	ErrorCodeFileUploadFailed  = 206009
	ErrorCodeFileDeleteFailed  = 206010

	// 角色定义相关错误 (207xxx)
	ErrorCodeRoleDefinitionNotFound      = 207001 // 角色定义不存在
	ErrorCodeRoleNameAlreadyExists       = 207002 // 角色名称已存在
	ErrorCodeSystemRoleCannotModify      = 207003 // 系统角色无法修改
	ErrorCodeSystemRoleCannotDelete      = 207004 // 系统角色无法删除
	ErrorCodeRoleInUseCannotDelete       = 207005 // 角色正在使用中，无法删除
	ErrorCodeRoleAssignmentNotFound      = 207006 // 角色分配不存在
	ErrorCodeRoleAssignmentAlreadyExists = 207007 // 角色分配已存在
	ErrorCodeRoleAssignmentConflict      = 207008 // 角色分配冲突
	ErrorCodeAssignerPermissionDenied    = 207009 // 分配者权限不足
	ErrorCodeMenuNotFound                = 207010 // 菜单不存在
	ErrorCodeMenuAlreadyExists           = 207011 // 菜单已存在
	ErrorCodeMenuConfigInvalid           = 207012 // 菜单配置无效
	ErrorCodeMenuYAMLParseFailed         = 207013 // 菜单YAML解析失败
	ErrorCodeMenuTreeInvalid             = 207014 // 菜单树结构无效
	ErrorCodeMenuPermissionDenied        = 207015 // 菜单权限不足
	ErrorCodeNoActiveRoles               = 207016 // 用户没有可用角色
	ErrorCodeSystemRoleCannotRevoke      = 207017 // 系统用户的系统角色无法撤销
)
