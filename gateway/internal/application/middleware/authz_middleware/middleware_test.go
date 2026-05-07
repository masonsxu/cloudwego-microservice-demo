package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRules(t *testing.T) *Rules {
	t.Helper()

	return &Rules{
		Default: DefaultAllow,
		Public: []Endpoint{
			{Method: "GET", Path: "/healthz"},
			{Method: "POST", Path: "/api/v1/auth/login"},
		},
		Authenticated: []Endpoint{
			{Method: "GET", Path: "/api/v1/profile/me"},
		},
		Roles: []RolePrefix{
			{Prefix: "/api/v1/admin/", Require: []string{"admin", "superadmin"}},
			{Prefix: "/api/v1/identity/", Require: []string{"identity_admin"}},
		},
	}
}

func TestDecide_PublicAllowsWithoutUser(t *testing.T) {
	rules := newRules(t)

	d := Decide(rules, "GET", "/healthz", "", nil)
	assert.Equal(t, OutcomeAllow, d.Outcome)
	assert.Equal(t, "public", d.MatchedRule)
}

func TestDecide_NonPublicRequiresUser(t *testing.T) {
	rules := newRules(t)

	d := Decide(rules, "GET", "/api/v1/profile/me", "", nil)
	assert.Equal(t, OutcomeUnauthorized, d.Outcome)
	assert.Contains(t, d.Reason, "missing X-User-Id")
}

func TestDecide_AuthenticatedAllowsAnyRole(t *testing.T) {
	rules := newRules(t)

	d := Decide(rules, "GET", "/api/v1/profile/me", "user-1", []string{"any_role"})
	assert.Equal(t, OutcomeAllow, d.Outcome)
	assert.Equal(t, "authenticated", d.MatchedRule)
}

func TestDecide_RolesAllowsWhenMatched(t *testing.T) {
	rules := newRules(t)

	d := Decide(rules, "GET", "/api/v1/admin/dashboard", "user-1", []string{"admin"})
	assert.Equal(t, OutcomeAllow, d.Outcome)
	assert.Equal(t, "roles:/api/v1/admin/", d.MatchedRule)
}

func TestDecide_RolesForbidsWhenMissing(t *testing.T) {
	rules := newRules(t)

	d := Decide(rules, "GET", "/api/v1/admin/dashboard", "user-1", []string{"doctor"})
	assert.Equal(t, OutcomeForbidden, d.Outcome)
	assert.Equal(t, "roles:/api/v1/admin/", d.MatchedRule)
	assert.Equal(t, []string{"admin", "superadmin"}, d.RequiredRole)
}

func TestDecide_DefaultAllowFallsThrough(t *testing.T) {
	rules := newRules(t)

	d := Decide(rules, "GET", "/api/v1/menu/mine", "user-1", nil)
	assert.Equal(t, OutcomeAllow, d.Outcome)
	assert.Equal(t, "default:allow", d.MatchedRule)
}

func TestDecide_DefaultDenyForbidsUnknownPath(t *testing.T) {
	rules := newRules(t)
	rules.Default = DefaultDeny

	d := Decide(rules, "GET", "/api/v1/menu/mine", "user-1", nil)
	assert.Equal(t, OutcomeForbidden, d.Outcome)
	assert.Equal(t, "default:deny", d.MatchedRule)
}

func TestDecide_RolePrefixHasPriorityOverDefault(t *testing.T) {
	rules := newRules(t)
	rules.Default = DefaultDeny

	// 命中 admin 前缀 → 走 roles 决策（即便 default=deny 也以角色判定为准）
	d := Decide(rules, "DELETE", "/api/v1/admin/users/123", "user-1", []string{"superadmin"})
	assert.Equal(t, OutcomeAllow, d.Outcome)
	assert.Equal(t, "roles:/api/v1/admin/", d.MatchedRule)
}
