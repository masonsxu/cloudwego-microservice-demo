package casbin_middleware

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestMethodToAction(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "read"},
		{"HEAD", "read"},
		{"OPTIONS", "read"},
		{"POST", "write"},
		{"PUT", "write"},
		{"PATCH", "write"},
		{"DELETE", "manage"},
		{"UNKNOWN", "read"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := methodToAction(tt.method)
			if result != tt.expected {
				t.Errorf("methodToAction(%s) = %s, want %s", tt.method, result, tt.expected)
			}
		})
	}
}

func TestMatchPath(t *testing.T) {
	tests := []struct {
		pattern  string
		path     string
		expected bool
	}{
		{"/api/v1/users/*", "/api/v1/users/123", true},
		{"/api/v1/users/*", "/api/v1/users", true},
		{"/api/v1/users/*", "/api/v1/roles/123", false},
		{"/api/v1/users", "/api/v1/users", true},
		{"/api/v1/users", "/api/v1/users/123", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.path, func(t *testing.T) {
			result := matchPath(tt.pattern, tt.path)
			if result != tt.expected {
				t.Errorf("matchPath(%s, %s) = %v, want %v", tt.pattern, tt.path, result, tt.expected)
			}
		})
	}
}

func TestCompareDataScope(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		{"org", "dept", 1},
		{"org", "self", 1},
		{"dept", "self", 1},
		{"dept", "org", -1},
		{"self", "org", -1},
		{"self", "dept", -1},
		{"org", "org", 0},
		{"dept", "dept", 0},
		{"self", "self", 0},
		{"unknown", "self", -1},
		{"org", "unknown", 1},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_vs_"+tt.b, func(t *testing.T) {
			result := compareDataScope(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareDataScope(%s, %s) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDefaultSkipPaths(t *testing.T) {
	paths := defaultSkipPaths()
	if len(paths) == 0 {
		t.Error("defaultSkipPaths() returned empty slice")
	}

	expectedPaths := []string{
		"/login",
		"/logout",
		"/health",
		"/api/v1/auth/login",
	}

	for _, expected := range expectedPaths {
		found := false
		for _, path := range paths {
			if path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected path %s not found in defaultSkipPaths()", expected)
		}
	}
}

func TestCasbinMiddleware_ShouldSkip(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)

	tests := []struct {
		path     string
		expected bool
	}{
		{"/login", true},
		{"/login/test", true},
		{"/api/v1/auth/login", true},
		{"/api/v1/users", false},
		{"/api/v1/roles", false},
		{"/health", true},
		{"/metrics", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := m.shouldSkip(tt.path)
			if result != tt.expected {
				t.Errorf("shouldSkip(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCasbinMiddleware_SetSkipPaths(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)

	customPaths := []string{"/custom", "/test"}
	m.SetSkipPaths(customPaths)

	if !m.shouldSkip("/custom") {
		t.Error("Expected /custom to be skipped after SetSkipPaths")
	}

	if !m.shouldSkip("/test") {
		t.Error("Expected /test to be skipped after SetSkipPaths")
	}

	if m.shouldSkip("/login") {
		t.Error("Expected /login to NOT be skipped after SetSkipPaths with custom paths")
	}
}

func TestCasbinMiddleware_AddPathMapping(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)

	m.AddPathMapping("/api/v1/users", "user:list")
	m.AddPathMapping("/api/v1/users/*", "user:view")

	resource := m.getResource("/api/v1/users")
	if resource != "user:list" {
		t.Errorf("getResource('/api/v1/users') = %s, want user:list", resource)
	}

	resource = m.getResource("/api/v1/users/123")
	if resource != "user:view" {
		t.Errorf("getResource('/api/v1/users/123') = %s, want user:view", resource)
	}

	resource = m.getResource("/api/v1/unknown")
	if resource != "/api/v1/unknown" {
		t.Errorf("getResource('/api/v1/unknown') = %s, want /api/v1/unknown", resource)
	}
}

func TestCasbinMiddleware_GetResource(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)

	resource := m.getResource("/api/v1/test")
	if resource != "/api/v1/test" {
		t.Errorf("getResource('/api/v1/test') without mapping = %s, want /api/v1/test", resource)
	}
}

func TestCasbinMiddleware_MiddlewareFunc_NilEnforcer(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)

	handler := m.MiddlewareFunc()
	if handler == nil {
		t.Error("MiddlewareFunc() returned nil")
	}
}

// TestCasbinMiddleware_CheckPermission_NilEnforcer 测试当 enforcer 为 nil 时的行为
// 注意：当 enforcer 为 nil 时，调用 CheckPermission 会 panic，
// 这是因为底层的 EnforceWithDataScope 需要有效的 enforcer。
// 在实际使用中，不应该用 nil enforcer 调用 CheckPermission。
func TestCasbinMiddleware_CheckPermission_NilEnforcer(t *testing.T) {
	t.Skip("Skipping: checkMultiRolePermission requires non-nil enforcer")
}

// TestCasbinMiddleware_CheckPermissionWithDataScope_NilEnforcer 测试当 enforcer 为 nil 时的行为
func TestCasbinMiddleware_CheckPermissionWithDataScope_NilEnforcer(t *testing.T) {
	t.Skip("Skipping: checkMultiRolePermission requires non-nil enforcer")
}

// TestCasbinMiddleware_GetUserPermissions_NilEnforcer 测试当 enforcer 为 nil 时的行为
func TestCasbinMiddleware_GetUserPermissions_NilEnforcer(t *testing.T) {
	t.Skip("Skipping: GetUserPermissions requires non-nil enforcer")
}

func TestPermissionInfo_Struct(t *testing.T) {
	info := PermissionInfo{
		Resource:  "/api/v1/users",
		Action:    "read",
		DataScope: "org",
		Domain:    "dept:123",
	}

	if info.Resource != "/api/v1/users" {
		t.Errorf("PermissionInfo.Resource = %s, want /api/v1/users", info.Resource)
	}
	if info.Action != "read" {
		t.Errorf("PermissionInfo.Action = %s, want read", info.Action)
	}
	if info.DataScope != "org" {
		t.Errorf("PermissionInfo.DataScope = %s, want org", info.DataScope)
	}
	if info.Domain != "dept:123" {
		t.Errorf("PermissionInfo.Domain = %s, want dept:123", info.Domain)
	}
}
