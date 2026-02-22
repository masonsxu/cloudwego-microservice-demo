package redis

import (
	"testing"
)

func TestGeneratePermissionCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		roleIDs  []string
		resource string
		action   string
		expected string
	}{
		{
			name:     "single role",
			userID:   "user123",
			roleIDs:  []string{"role1"},
			resource: "/api/v1/users",
			action:   "read",
			expected: "user123:[role1]:/api/v1/users:read",
		},
		{
			name:     "multiple roles",
			userID:   "user456",
			roleIDs:  []string{"role1", "role2", "role3"},
			resource: "/api/v1/orders",
			action:   "write",
			expected: "user456:[role1,role2,role3]:/api/v1/orders:write",
		},
		{
			name:     "empty roles",
			userID:   "user789",
			roleIDs:  []string{},
			resource: "/api/v1/test",
			action:   "manage",
			expected: "user789:[]:/api/v1/test:manage",
		},
		{
			name:     "uuid format",
			userID:   "550e8400-e29b-41d4-a716-446655440000",
			roleIDs:  []string{"550e8400-e29b-41d4-a716-446655440001"},
			resource: "menu:system",
			action:   "read",
			expected: "550e8400-e29b-41d4-a716-446655440000:[550e8400-e29b-41d4-a716-446655440001]:menu:system:read",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GeneratePermissionCacheKey(tt.userID, tt.roleIDs, tt.resource, tt.action)
			if result != tt.expected {
				t.Errorf("GeneratePermissionCacheKey() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestPermissionCacheResult_Struct(t *testing.T) {
	result := PermissionCacheResult{
		Allowed:   true,
		DataScope: "org",
		CachedAt:  1704067200, // 2024-01-01 00:00:00 UTC
	}

	if !result.Allowed {
		t.Error("Expected Allowed to be true")
	}

	if result.DataScope != "org" {
		t.Errorf("Expected DataScope to be 'org', got %s", result.DataScope)
	}

	if result.CachedAt != 1704067200 {
		t.Errorf("Expected CachedAt to be 1704067200, got %d", result.CachedAt)
	}
}

func TestPolicyCache_GetPermissionKey(t *testing.T) {
	pc := &PolicyCache{
		prefix: "gateway:permission:",
	}

	key := pc.getPermissionKey("user:123:resource:read")

	expected := "gateway:permission:user:123:resource:read"
	if key != expected {
		t.Errorf("getPermissionKey() = %s, want %s", key, expected)
	}
}

func TestPolicyCache_GetUserPermissionsPattern(t *testing.T) {
	pc := &PolicyCache{
		prefix: "gateway:permission:",
	}

	pattern := pc.getUserPermissionsPattern("user123")

	expected := "gateway:permission:user:user123:*"
	if pattern != expected {
		t.Errorf("getUserPermissionsPattern() = %s, want %s", pattern, expected)
	}
}

func TestPolicyCache_GetRolePermissionsPattern(t *testing.T) {
	pc := &PolicyCache{
		prefix: "gateway:permission:",
	}

	pattern := pc.getRolePermissionsPattern("role456")

	expected := "gateway:permission:role:role456:*"
	if pattern != expected {
		t.Errorf("getRolePermissionsPattern() = %s, want %s", pattern, expected)
	}
}

func TestPolicyCache_GetUserRolesKey(t *testing.T) {
	pc := &PolicyCache{
		prefix: "gateway:permission:",
	}

	key := pc.getUserRolesKey("user789")

	expected := "gateway:permission:user_roles:user789"
	if key != expected {
		t.Errorf("getUserRolesKey() = %s, want %s", key, expected)
	}
}

func TestPolicyCache_KeyPatterns(t *testing.T) {
	pc := &PolicyCache{
		prefix: "test:prefix:",
	}

	tests := []struct {
		name     string
		method   string
		input    string
		expected string
	}{
		{
			name:     "permission key",
			method:   "getPermissionKey",
			input:    "testKey",
			expected: "test:prefix:testKey",
		},
		{
			name:     "user permissions pattern",
			method:   "getUserPermissionsPattern",
			input:    "user1",
			expected: "test:prefix:user:user1:*",
		},
		{
			name:     "role permissions pattern",
			method:   "getRolePermissionsPattern",
			input:    "admin",
			expected: "test:prefix:role:admin:*",
		},
		{
			name:     "user roles key",
			method:   "getUserRolesKey",
			input:    "user2",
			expected: "test:prefix:user_roles:user2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string

			switch tt.method {
			case "getPermissionKey":
				result = pc.getPermissionKey(tt.input)
			case "getUserPermissionsPattern":
				result = pc.getUserPermissionsPattern(tt.input)
			case "getRolePermissionsPattern":
				result = pc.getRolePermissionsPattern(tt.input)
			case "getUserRolesKey":
				result = pc.getUserRolesKey(tt.input)
			}

			if result != tt.expected {
				t.Errorf("%s(%s) = %s, want %s", tt.method, tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewPolicyCache(t *testing.T) {
	cache := NewPolicyCache(nil, nil)
	if cache == nil {
		t.Error("NewPolicyCache() returned nil")
	}

	// 验证返回的是 PolicyCacheService 接口类型
	_ = cache
}
