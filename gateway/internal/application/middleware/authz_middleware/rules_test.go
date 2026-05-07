package middleware

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRules_Defaults(t *testing.T) {
	rules, err := ParseRules([]byte(``))
	require.NoError(t, err)
	assert.Equal(t, DefaultAllow, rules.Default)
	assert.Empty(t, rules.Public)
	assert.Empty(t, rules.Authenticated)
	assert.Empty(t, rules.Roles)
}

func TestParseRules_FullExample(t *testing.T) {
	yaml := []byte(`
default: deny
public:
  - "GET  /healthz"
  - "POST /api/v1/auth/login"
  - "*    /.well-known/jwks.json"
authenticated:
  - "GET /api/v1/profile/me"
roles:
  - prefix: /api/v1/admin/
    require: [admin, superadmin]
  - prefix: /api/v1/identity/
    require: [identity_admin, superadmin]
`)
	rules, err := ParseRules(yaml)
	require.NoError(t, err)

	assert.Equal(t, DefaultDeny, rules.Default)

	require.Len(t, rules.Public, 3)
	assert.Equal(t, Endpoint{Method: "GET", Path: "/healthz"}, rules.Public[0])
	assert.Equal(t, Endpoint{Method: "POST", Path: "/api/v1/auth/login"}, rules.Public[1])
	assert.Equal(t, Endpoint{Method: "*", Path: "/.well-known/jwks.json"}, rules.Public[2])

	require.Len(t, rules.Authenticated, 1)
	assert.Equal(t, Endpoint{Method: "GET", Path: "/api/v1/profile/me"}, rules.Authenticated[0])

	require.Len(t, rules.Roles, 2)
	assert.Equal(t, "/api/v1/admin/", rules.Roles[0].Prefix)
	assert.Equal(t, []string{"admin", "superadmin"}, rules.Roles[0].Require)
	assert.Equal(t, "/api/v1/identity/", rules.Roles[1].Prefix)
}

func TestParseRules_InvalidDefault(t *testing.T) {
	_, err := ParseRules([]byte(`default: maybe`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "非法的 default 值")
}

func TestParseRules_BadEndpointFormat(t *testing.T) {
	cases := []struct {
		name string
		yaml string
	}{
		{"missing method", "public:\n  - \"/healthz\""},
		{"path without slash", "public:\n  - \"GET healthz\""},
		{"empty entry", "public:\n  - \"\""},
		{"too many fields", "public:\n  - \"GET /a /b\""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseRules([]byte(tc.yaml))
			require.Error(t, err)
			assert.Contains(t, err.Error(), "public 规则")
		})
	}
}

func TestParseRules_RoleRulesValidation(t *testing.T) {
	cases := []struct {
		name   string
		yaml   string
		errSub string
	}{
		{
			name:   "missing prefix",
			yaml:   "roles:\n  - require: [admin]",
			errSub: "缺少 prefix",
		},
		{
			name:   "prefix without slash",
			yaml:   "roles:\n  - prefix: api/v1/admin\n    require: [admin]",
			errSub: "必须以 / 开头",
		},
		{
			name:   "missing require",
			yaml:   "roles:\n  - prefix: /api/v1/admin/",
			errSub: "缺少 require",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseRules([]byte(tc.yaml))
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errSub)
		})
	}
}

func TestRules_MatchPublic(t *testing.T) {
	rules := &Rules{
		Public: []Endpoint{
			{Method: "GET", Path: "/healthz"},
			{Method: "*", Path: "/.well-known/jwks.json"},
		},
	}

	assert.True(t, rules.MatchPublic("GET", "/healthz"))
	assert.True(t, rules.MatchPublic("GET", "/.well-known/jwks.json"))
	assert.True(t, rules.MatchPublic("POST", "/.well-known/jwks.json"))
	assert.False(t, rules.MatchPublic("POST", "/healthz"))
	assert.False(t, rules.MatchPublic("GET", "/healthz/extra"))
}

func TestRules_MatchAuthenticated(t *testing.T) {
	rules := &Rules{
		Authenticated: []Endpoint{
			{Method: "GET", Path: "/api/v1/profile/me"},
		},
	}

	assert.True(t, rules.MatchAuthenticated("GET", "/api/v1/profile/me"))
	assert.False(t, rules.MatchAuthenticated("POST", "/api/v1/profile/me"))
}

func TestRules_MatchRolePrefix(t *testing.T) {
	rules := &Rules{
		Roles: []RolePrefix{
			{Prefix: "/api/v1/admin/users/", Require: []string{"superadmin"}},
			{Prefix: "/api/v1/admin/", Require: []string{"admin", "superadmin"}},
		},
	}

	rule, hit := rules.MatchRolePrefix("/api/v1/admin/users/123")
	require.True(t, hit)
	assert.Equal(t, "/api/v1/admin/users/", rule.Prefix)

	rule, hit = rules.MatchRolePrefix("/api/v1/admin/dashboard")
	require.True(t, hit)
	assert.Equal(t, "/api/v1/admin/", rule.Prefix)

	_, hit = rules.MatchRolePrefix("/api/v1/profile/me")
	assert.False(t, hit)
}

func TestRolePrefix_HasAnyRole(t *testing.T) {
	rule := RolePrefix{Require: []string{"admin", "superadmin"}}

	assert.True(t, rule.HasAnyRole([]string{"doctor", "admin"}))
	assert.True(t, rule.HasAnyRole([]string{"superadmin"}))
	assert.False(t, rule.HasAnyRole([]string{"doctor"}))
	assert.False(t, rule.HasAnyRole([]string{}))
	assert.False(t, rule.HasAnyRole(nil))

	// 角色名前后空白不影响匹配
	assert.True(t, rule.HasAnyRole([]string{"  admin  "}))
}

func TestSplitHeader(t *testing.T) {
	assert.Nil(t, splitHeader(""))
	assert.Equal(t, []string{"admin"}, splitHeader("admin"))
	assert.Equal(t, []string{"admin", "doctor"}, splitHeader("admin,doctor"))
	assert.Equal(t, []string{"admin", "doctor"}, splitHeader(" admin , doctor "))
	assert.Equal(t, []string{"admin"}, splitHeader("admin,,"))
	assert.Empty(t, splitHeader(",,,"))
}

// TestLoadRulesFromFile_DefaultRulesYAML 验证仓库内默认 authz_rules.yaml 能被正确解析。
//
// 防止 commit 时 YAML 语法 / 端点格式错误未被发现。
func TestLoadRulesFromFile_DefaultRulesYAML(t *testing.T) {
	// 测试在 internal/application/middleware/authz_middleware/，默认 yaml 在 gateway/config/
	cwd, err := os.Getwd()
	require.NoError(t, err)

	yamlPath := filepath.Join(cwd, "..", "..", "..", "..", "config", "authz_rules.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		t.Skipf("默认 authz_rules.yaml 不存在（%s），跳过：%v", yamlPath, err)
	}

	rules, err := LoadRulesFromFile(yamlPath)
	require.NoError(t, err)
	assert.Equal(t, DefaultAllow, rules.Default)
	assert.NotEmpty(t, rules.Public)
	// /healthz 应在 public 列表中
	assert.True(t, rules.MatchPublic("GET", "/healthz"))
	// 登录端点应公开
	assert.True(t, rules.MatchPublic("POST", "/api/v1/identity/auth/login"))
	// JWKS 应公开
	assert.True(t, rules.MatchPublic("GET", "/.well-known/jwks.json"))
}
