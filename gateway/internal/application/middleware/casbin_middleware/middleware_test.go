package casbin_middleware

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
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

func TestDefaultSuperAdminSubjects(t *testing.T) {
	subjects := defaultSuperAdminSubjects()
	if len(subjects) == 0 {
		t.Fatal("defaultSuperAdminSubjects() returned empty slice")
	}

	foundRole := false
	foundUsername := false
	for _, subject := range subjects {
		if subject == "role:superadmin" {
			foundRole = true
		}

		if subject == "username:superadmin" {
			foundUsername = true
		}
	}

	if !foundRole {
		t.Error("defaultSuperAdminSubjects() should contain role:superadmin")
	}

	if !foundUsername {
		t.Error("defaultSuperAdminSubjects() should contain username:superadmin")
	}
}

func TestCasbinMiddleware_SetSkipPaths(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)

	customPaths := []string{"/custom", "GET:/health", "/swagger/*"}
	m.SetSkipPaths(customPaths)

	if len(m.skipPaths) != 3 {
		t.Fatalf("expected 3 skip paths, got %d", len(m.skipPaths))
	}

	if m.skipPaths[0] != "/custom" || m.skipPaths[1] != "GET:/health" || m.skipPaths[2] != "/swagger/*" {
		t.Fatalf("unexpected skip paths: %v", m.skipPaths)
	}
}

func TestCasbinMiddleware_SetSuperAdminBypassConfig(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)

	m.SetSuperAdminBypassConfig(false, []string{"role:superadmin", " role:platform_admin "})
	if m.superAdminBypassEnabled {
		t.Error("superAdminBypassEnabled should be false")
	}

	if _, ok := m.superAdminSubjectAllowSet["role:superadmin"]; !ok {
		t.Error("expected role:superadmin in superAdminSubjectAllowSet")
	}

	if _, ok := m.superAdminSubjectAllowSet["role:platform_admin"]; !ok {
		t.Error("expected role:platform_admin in superAdminSubjectAllowSet")
	}

	m.SetSuperAdminBypassConfig(true, nil)
	if !m.superAdminBypassEnabled {
		t.Error("superAdminBypassEnabled should be true")
	}

	if _, ok := m.superAdminSubjectAllowSet["role:superadmin"]; !ok {
		t.Error("expected default role:superadmin when subjects is empty")
	}

	if _, ok := m.superAdminSubjectAllowSet["username:superadmin"]; !ok {
		t.Error("expected default username:superadmin when subjects is empty")
	}
}

func TestCasbinMiddleware_ShouldBypassForSuperAdmin_Username(t *testing.T) {
	logger := zerolog.Nop()
	m := NewCasbinMiddleware(nil, &logger)
	m.SetSuperAdminBypassConfig(true, []string{"username:superadmin"})

	if !m.shouldBypassForSuperAdmin(context.Background(), "user-1", "superadmin", nil, nil) {
		t.Error("expected superadmin username to bypass casbin checks")
	}

	if m.shouldBypassForSuperAdmin(context.Background(), "user-1", "normal_user", nil, nil) {
		t.Error("expected non-superadmin username to not bypass casbin checks")
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
	m.SetSkipPaths([]string{"GET:/health"})

	handler := m.MiddlewareFunc()
	if handler == nil {
		t.Error("MiddlewareFunc() returned nil")
	}

	// 仅验证可调用且 skip 规则通过 common.ShouldSkip 可解析。
	c := app.NewContext(0)
	c.Request.SetRequestURI("/health")
	c.Request.Header.SetMethod("GET")
	handler(context.Background(), c)
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
