package middleware

import (
	"github.com/hertz-contrib/jwt"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
)

// createPayloadFromLoginData 从登录数据map创建JWT载荷
func createPayloadFromLoginData(data map[string]interface{}) jwt.MapClaims {
	claims := jwt.MapClaims{}

	// 设置标准JWT声明
	if userID, exists := data[IdentityKey]; exists && userID != nil {
		claims[IdentityKey] = userID
	}

	if username, exists := data[Username]; exists && username != nil {
		claims[Username] = username
	}

	if status, exists := data[Status]; exists && status != nil {
		claims[Status] = status
	}

	if roleID, exists := data[RoleID]; exists && roleID != nil {
		claims[RoleID] = roleID
	}

	if orgID, exists := data[OrganizationID]; exists && orgID != nil {
		claims[OrganizationID] = orgID
	}

	if deptID, exists := data[DepartmentID]; exists && deptID != nil {
		claims[DepartmentID] = deptID
	}

	if permission, exists := data[CorePermission]; exists && permission != nil {
		claims[CorePermission] = permission
	}

	return claims
}

// createPayloadFromJWTClaimsDTO 从JWTClaimsDTO创建JWT载荷
func createPayloadFromJWTClaimsDTO(user *http_base.JWTClaimsDTO) jwt.MapClaims {
	claims := jwt.MapClaims{}

	if user.UserProfileID != nil {
		claims[IdentityKey] = *user.UserProfileID
	}

	if user.Username != nil {
		claims[Username] = *user.Username
	}

	if user.Status != nil {
		claims[Status] = int(*user.Status)
	}

	if user.RoleID != nil {
		claims[RoleID] = *user.RoleID
	}

	if user.OrganizationID != nil {
		claims[OrganizationID] = *user.OrganizationID
	}

	if user.DepartmentID != nil {
		claims[DepartmentID] = *user.DepartmentID
	}

	if user.Permission != nil {
		claims[CorePermission] = user.Permission
	}

	// 设置JWT标准字段
	if user.Exp != nil {
		claims["exp"] = *user.Exp
	}

	if user.Iat != nil {
		claims["iat"] = *user.Iat
	}

	return claims
}

// payloadFunc 创建JWT载荷的函数
func payloadFunc(data interface{}) jwt.MapClaims {
	// 处理来自Authenticator的登录响应数据
	if loginData, ok := data.(map[string]interface{}); ok {
		return createPayloadFromLoginData(loginData)
	}

	// 处理直接传入的JWTClaimsDTO结构
	if user, ok := data.(*http_base.JWTClaimsDTO); ok {
		return createPayloadFromJWTClaimsDTO(user)
	}

	return jwt.MapClaims{}
}

// extractStringClaim 从claims中提取字符串值
func extractStringClaim(claims jwt.MapClaims, key string) (string, bool) {
	if value, exists := claims[key]; exists {
		if str, ok := value.(string); ok && str != "" {
			return str, true
		}
	}

	return "", false
}

// extractIntClaim 从claims中提取整数值（支持float64和int）
func extractIntClaim(claims jwt.MapClaims, key string) (int, bool) {
	if value, exists := claims[key]; exists {
		if intVal, ok := value.(float64); ok {
			return int(intVal), true
		} else if intVal, ok := value.(int); ok {
			return intVal, true
		}
	}

	return 0, false
}

// extractInt64Claim 从claims中提取int64值
func extractInt64Claim(claims jwt.MapClaims, key string) (int64, bool) {
	if value, exists := claims[key]; exists {
		if intVal, ok := value.(float64); ok {
			return int64(intVal), true
		}
	}

	return 0, false
}

// extractBasicUserInfo 提取用户基本信息（用户ID、用户名、状态）
// 该函数负责从用户个人信息中提取JWT所需的基本字段，包括用户ID、用户名和用户状态。
// 所有字段都进行了nil安全检查，确保只有非空值才会被添加到用户数据映射中。
//
// 参数:
//   - user: 用户个人信息对象，可能为nil
//   - userData: 目标数据映射，用于存储提取的用户信息
//
// 示例:
//
//	userData := make(map[string]interface{})
//	extractBasicUserInfo(loginResp.UserProfile, userData)
//	// userData 现在包含用户ID、用户名和状态（如果这些字段非空）
func extractBasicUserInfo(
	user *identity.UserProfileDTO,
	userData map[string]interface{},
	permission string,
) {
	// 检查用户对象是否为nil，如果为nil则直接返回
	if user == nil {
		return
	}

	// 提取用户ID，确保非空才添加到映射中
	if user.ID != nil {
		userData[IdentityKey] = *user.ID
	}

	// 提取用户名，确保非空才添加到映射中
	if user.Username != nil {
		userData[Username] = *user.Username
	}

	// 提取用户状态，确保非空才添加到映射中
	// 将int32类型转换为int类型以保持与现有代码的兼容性
	if user.Status != nil {
		userData[Status] = int(*user.Status)
	}

	// 提取权限，确保非空才添加到映射中
	if permission != "" {
		userData[CorePermission] = permission
	}
}

// extractMembershipInfo 提取成员关系信息（组织ID、部门ID）
// 该函数负责从成员关系列表中查找主要成员关系，并提取组织ID和部门ID。
// 函数会遍历所有成员关系，寻找标记为主要成员关系的记录，并从中提取组织和部门信息。
// 如果没有找到主要成员关系，则不会向用户数据映射中添加任何组织或部门信息。
//
// 参数:
//   - memberships: 成员关系列表，可能为nil或空切片
//   - userData: 目标数据映射，用于存储提取的成员关系信息
//
// 示例:
//
//	userData := make(map[string]interface{})
//	extractMembershipInfo(loginResp.Memberships, userData)
//	// userData 现在包含组织ID和部门ID（如果找到主要成员关系且这些字段非空）
func extractMembershipInfo(
	memberships []*identity.UserMembershipDTO,
	userData map[string]interface{},
) {
	// 检查成员关系列表是否为nil，如果为nil则直接返回
	if memberships == nil {
		return
	}

	// 遍历所有成员关系，寻找主要成员关系
	for _, membership := range memberships {
		// 检查成员关系对象是否为nil，以及是否标记为主要成员关系
		if membership != nil && membership.IsPrimary != nil && *membership.IsPrimary {
			// 提取组织ID，确保非空才添加到映射中
			if membership.OrganizationID != nil {
				userData[OrganizationID] = *membership.OrganizationID
			}

			// 提取部门ID，确保非空才添加到映射中
			if membership.DepartmentID != nil {
				userData[DepartmentID] = *membership.DepartmentID
			}

			// 找到主要成员关系后跳出循环，避免重复处理
			break
		}
	}
}

// extractRoleInfo 提取角色信息（角色ID）
// 该函数负责从角色ID列表中提取主要角色ID。根据业务逻辑，系统采用单角色模式，
// 因此从角色列表中选择第一个角色作为用户的主要角色。如果角色列表为空或nil，
// 则不会向用户数据映射中添加任何角色信息。
//
// 主要角色选择逻辑：
// - 如果角色列表不为空，选择第一个角色作为主要角色
// - 这种设计简化了权限管理，避免了多角色冲突的复杂性
// - 在实际业务场景中，用户通常只需要一个主要角色来确定其权限范围
//
// 参数:
//   - roleIDs: 角色ID列表，可能为nil或空切片
//   - userData: 目标数据映射，用于存储提取的角色信息
//
// 示例:
//
//	userData := make(map[string]interface{})
//	roleIDs := []string{"admin", "user", "viewer"}
//	extractRoleInfo(roleIDs, userData)
//	// userData 现在包含 roleID: "admin"（第一个角色）
//
//	// 处理空角色列表的情况
//	extractRoleInfo([]string{}, userData)
//	// userData 不会包含任何角色信息
func extractRoleInfo(roleIDs []string, userData map[string]interface{}) {
	// 检查角色列表是否为空，如果为空则直接返回
	// 这里使用len()检查而不是nil检查，因为空切片和nil切片都应该被视为无角色
	if len(roleIDs) == 0 {
		return
	}

	// 提取第一个角色作为主要角色
	// 根据业务需求，系统采用单角色模式，第一个角色被视为用户的主要角色
	userData[RoleID] = roleIDs[0]
}

// buildUserDataMap 构造用户信息map，供PayloadFunc使用
// 该函数作为主要的协调器，通过调用专门的辅助函数来提取JWT所需的所有字段。
// 采用模块化设计，将复杂的数据提取逻辑分解为独立的、可测试的组件：
// - extractBasicUserInfo: 提取用户基本信息（ID、用户名、状态）
// - extractMembershipInfo: 提取组织和部门信息
// - extractRoleInfo: 提取角色信息
//
// 这种设计降低了函数的复杂度，提高了代码的可读性和可维护性。
// 每个辅助函数都有单一职责，便于单独测试和复用。
//
// 参数:
//   - loginResp: 登录响应DTO，包含用户完整信息
//
// 返回值:
//   - map[string]interface{}: 包含JWT载荷所需字段的映射
//
// 示例:
//
//	userData := buildUserDataMap(loginResponse)
//	// userData 包含: userProfileID, username, status, organizationID, departmentID, roleID
func buildUserDataMap(
	loginResp *identity.LoginResponseDTO,
	permission string,
) map[string]interface{} {
	userData := map[string]interface{}{}

	if loginResp == nil || loginResp.UserProfile == nil {
		return userData
	}

	user := loginResp.UserProfile

	// 使用helper函数提取用户基本信息
	extractBasicUserInfo(user, userData, permission)

	// 使用helper函数提取成员关系信息
	extractMembershipInfo(loginResp.Memberships, userData)

	// 使用helper函数提取角色信息
	extractRoleInfo(loginResp.RoleIDs, userData)

	return userData
}
