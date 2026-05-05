package middleware

import (
	"github.com/hertz-contrib/jwt"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
)

// createPayloadFromLoginData 从登录数据 map 创建 JWT 载荷（目标态最小 schema）
func createPayloadFromLoginData(data map[string]interface{}) jwt.MapClaims {
	claims := jwt.MapClaims{}

	if userID, exists := data[IdentityKey]; exists && userID != nil {
		claims[IdentityKey] = userID
	}

	if username, exists := data[Username]; exists && username != nil {
		claims[Username] = username
	}

	if tenant, exists := data[Tenant]; exists && tenant != nil {
		claims[Tenant] = tenant
	}

	if roles, exists := data[Roles]; exists && roles != nil {
		claims[Roles] = roles
	}

	return claims
}

// createPayloadFromJWTClaimsDTO 从 JWTClaimsDTO 创建 JWT 载荷
func createPayloadFromJWTClaimsDTO(user *http_base.JWTClaimsDTO) jwt.MapClaims {
	claims := jwt.MapClaims{}

	if user.UserProfileID != nil {
		claims[IdentityKey] = *user.UserProfileID
	}

	if user.Username != nil {
		claims[Username] = *user.Username
	}

	return claims
}

// payloadFunc 创建 JWT 载荷的函数
func payloadFunc(data interface{}) jwt.MapClaims {
	if loginData, ok := data.(map[string]interface{}); ok {
		return createPayloadFromLoginData(loginData)
	}

	if user, ok := data.(*http_base.JWTClaimsDTO); ok {
		return createPayloadFromJWTClaimsDTO(user)
	}

	return jwt.MapClaims{}
}

// extractStringClaim 从 claims 中提取字符串值
func extractStringClaim(claims jwt.MapClaims, key string) (string, bool) {
	if value, exists := claims[key]; exists {
		if str, ok := value.(string); ok && str != "" {
			return str, true
		}
	}

	return "", false
}

// extractInt64Claim 从 claims 中提取 int64 值
func extractInt64Claim(claims jwt.MapClaims, key string) (int64, bool) {
	if value, exists := claims[key]; exists {
		if intVal, ok := value.(float64); ok {
			return int64(intVal), true
		}
	}

	return 0, false
}

// extractStringSliceClaim 从 claims 中提取字符串切片值
func extractStringSliceClaim(claims jwt.MapClaims, key string) ([]string, bool) {
	value, exists := claims[key]
	if !exists {
		return nil, false
	}

	if strSlice, ok := value.([]string); ok && len(strSlice) > 0 {
		return strSlice, true
	}

	ifaceSlice, ok := value.([]interface{})
	if !ok || len(ifaceSlice) == 0 {
		return nil, false
	}

	result := make([]string, 0, len(ifaceSlice))
	for _, item := range ifaceSlice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}

	if len(result) > 0 {
		return result, true
	}

	return nil, false
}

// buildUserDataMap 构造用户信息 map，供 PayloadFunc 使用（目标态最小字段集）
func buildUserDataMap(loginResp *identity.LoginResponseDTO) map[string]interface{} {
	userData := map[string]interface{}{}

	if loginResp == nil || loginResp.UserProfile == nil {
		return userData
	}

	user := loginResp.UserProfile

	if user.Id != nil {
		userData[IdentityKey] = *user.Id
	}

	if user.Username != nil {
		userData[Username] = *user.Username
	}

	// 从主成员关系提取 tenant（OrganizationID）
	if loginResp.Memberships != nil {
		for _, m := range loginResp.Memberships {
			if m != nil && m.IsPrimary != nil && *m.IsPrimary && m.OrganizationID != nil {
				userData[Tenant] = *m.OrganizationID
				break
			}
		}
	}

	// 存储角色 code 列表（当前 loginResp.RoleIDs 是 UUID，Phase 4 改为 codes）
	if len(loginResp.RoleIDs) > 0 {
		userData[Roles] = loginResp.RoleIDs
	}

	return userData
}
