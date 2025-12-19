package errno

// 业务错误定义
var (
	// 通用业务错误
	ErrInvalidParams   = NewErrNo(ErrorCodeInvalidParams, "参数错误")
	ErrOperationFailed = NewErrNo(ErrorCodeOperationFailed, "操作失败")

	// 用户相关错误
	ErrUserNotFound           = NewErrNo(ErrorCodeUserNotFound, "用户不存在")
	ErrUserAlreadyExists      = NewErrNo(ErrorCodeUserAlreadyExists, "用户已存在")
	ErrUsernameAlreadyExists  = NewErrNo(ErrorCodeUsernameAlreadyExists, "用户名已存在")
	ErrEmailAlreadyExists     = NewErrNo(ErrorCodeEmailAlreadyExists, "邮箱已存在")
	ErrPhoneAlreadyExists     = NewErrNo(ErrorCodePhoneAlreadyExists, "手机号已存在")
	ErrInvalidPassword        = NewErrNo(ErrorCodeInvalidPassword, "密码不正确")
	ErrUserInactive           = NewErrNo(ErrorCodeUserInactive, "用户未激活")
	ErrPhoneNumberAlreadyUsed = NewErrNo(ErrorCodePhoneNumberAlreadyUsed, "手机号已被使用")
	ErrInvalidAccountType     = NewErrNo(ErrorCodeInvalidAccountType, "账户类型不符合要求")
	ErrInvalidCredentials     = NewErrNo(ErrorCodeInvalidCredentials, "用户名或密码错误")
	ErrUserSuspended          = NewErrNo(ErrorCodeUserSuspended, "用户已被暂停")
	ErrMustChangePassword     = NewErrNo(ErrorCodeMustChangePassword, "请先修改密码")

	// 系统用户保护相关错误
	ErrSystemUserCannotDelete    = NewErrNo(ErrorCodeSystemUserCannotDelete, "系统用户无法删除")
	ErrSystemUserCannotModifyKey = NewErrNo(ErrorCodeSystemUserCannotModifyKey, "系统用户关键属性无法修改")

	// 组织相关错误
	ErrOrganizationNotFound       = NewErrNo(ErrorCodeOrganizationNotFound, "组织不存在")
	ErrParentOrganizationNotFound = NewErrNo(ErrorCodeParentOrganizationNotFound, "引用的父组织不存在")
	ErrOrganizationAlreadyExists  = NewErrNo(ErrorCodeOrganizationAlreadyExists, "组织名称已存在")
	ErrOrganizationHasUsers       = NewErrNo(ErrorCodeOrganizationHasUsers, "组织下有关联用户，无法删除")

	// 部门相关错误
	ErrDepartmentNotFound          = NewErrNo(ErrorCodeDepartmentNotFound, "部门不存在")
	ErrDepartmentAlreadyExists     = NewErrNo(ErrorCodeDepartmentAlreadyExists, "部门已存在")
	ErrDepartmentCodeAlreadyExists = NewErrNo(
		ErrorCodeDepartmentCodeAlreadyExists,
		"部门代码已存在",
	)
	ErrCannotDeleteDepartmentWithMembers = NewErrNo(
		ErrorCodeCannotDeleteDepartmentWithMembers,
		"部门有成员，无法删除",
	)
	ErrDepartmentNameRequired         = NewErrNo(ErrorCodeDepartmentNameRequired, "部门名称不能为空")
	ErrDepartmentOrganizationRequired = NewErrNo(
		ErrorCodeDepartmentOrganizationRequired,
		"部门必须属于一个组织",
	)

	// 级联删除和数据一致性相关错误
	ErrUserNotInSameOrganization = NewErrNo(ErrorCodeUserNotInSameOrganization, "用户与团队不属于同一组织")
	ErrDataInconsistency         = NewErrNo(ErrorCodeDataInconsistency, "数据一致性错误")
	ErrTransactionFailed         = NewErrNo(ErrorCodeTransactionFailed, "事务执行失败")

	// 成员关系相关错误
	ErrMembershipNotFound      = NewErrNo(ErrorCodeMembershipNotFound, "成员关系不存在")
	ErrMembershipAlreadyExists = NewErrNo(ErrorCodeMembershipAlreadyExists, "成员关系已存在")

	// 组织Logo相关错误
	ErrLogoNotFound      = NewErrNo(ErrorCodeLogoNotFound, "Logo不存在")
	ErrLogoAlreadyBound  = NewErrNo(ErrorCodeLogoAlreadyBound, "Logo已被绑定")
	ErrLogoExpired       = NewErrNo(ErrorCodeLogoExpired, "Logo已过期")
	ErrLogoInvalidStatus = NewErrNo(ErrorCodeLogoInvalidStatus, "Logo状态无效")
	ErrLogoBindingFailed = NewErrNo(ErrorCodeLogoBindingFailed, "Logo绑定失败")
	ErrLogoAlreadyExists = NewErrNo(ErrorCodeLogoAlreadyExists, "Logo已存在")
	ErrInvalidFileType   = NewErrNo(ErrorCodeInvalidFileType, "不支持的文件类型")
	ErrFileSizeExceeded  = NewErrNo(ErrorCodeFileSizeExceeded, "文件大小超过限制")
	ErrFileUploadFailed  = NewErrNo(ErrorCodeFileUploadFailed, "文件上传失败")
	ErrFileDeleteFailed  = NewErrNo(ErrorCodeFileDeleteFailed, "文件删除失败")

	// 角色定义相关错误
	ErrRoleDefinitionNotFound = NewErrNo(ErrorCodeRoleDefinitionNotFound, "角色定义不存在")
	ErrRoleNameAlreadyExists  = NewErrNo(ErrorCodeRoleNameAlreadyExists, "角色名称已存在")
	ErrSystemRoleCannotModify = NewErrNo(ErrorCodeSystemRoleCannotModify, "系统角色无法修改")
	ErrSystemRoleCannotDelete = NewErrNo(ErrorCodeSystemRoleCannotDelete, "系统角色无法删除")
	ErrRoleInUseCannotDelete  = NewErrNo(ErrorCodeRoleInUseCannotDelete, "角色正在使用中，无法删除")

	// 用户角色分配相关错误
	ErrRoleAssignmentNotFound      = NewErrNo(ErrorCodeRoleAssignmentNotFound, "用户角色分配不存在")
	ErrRoleAssignmentAlreadyExists = NewErrNo(ErrorCodeRoleAssignmentAlreadyExists, "用户角色分配已存在")
	ErrRoleAssignmentConflict      = NewErrNo(ErrorCodeRoleAssignmentConflict, "用户角色分配冲突")
	ErrAssignerPermissionDenied    = NewErrNo(ErrorCodeAssignerPermissionDenied, "分配者权限不足")

	// 菜单权限相关错误
	ErrMenuNotFound         = NewErrNo(ErrorCodeMenuNotFound, "菜单不存在")
	ErrMenuAlreadyExists    = NewErrNo(ErrorCodeMenuAlreadyExists, "菜单已存在")
	ErrMenuConfigInvalid    = NewErrNo(ErrorCodeMenuConfigInvalid, "菜单配置无效")
	ErrMenuYAMLParseFailed  = NewErrNo(ErrorCodeMenuYAMLParseFailed, "菜单YAML解析失败")
	ErrMenuTreeInvalid      = NewErrNo(ErrorCodeMenuTreeInvalid, "菜单树结构无效")
	ErrMenuPermissionDenied = NewErrNo(ErrorCodeMenuPermissionDenied, "菜单权限不足")
	ErrNoActiveRoles        = NewErrNo(ErrorCodeNoActiveRoles, "用户没有可用角色")
	ErrSystemRoleCannotRevoke = NewErrNo(ErrorCodeSystemRoleCannotRevoke, "系统用户的系统角色无法撤销")
)
