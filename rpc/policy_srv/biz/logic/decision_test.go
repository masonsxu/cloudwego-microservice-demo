package logic

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	fileadapter "github.com/casbin/casbin/v3/persist/file-adapter"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	pb "github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/kitex_gen/policy_srv"
)

func newTestEnforcer(t *testing.T, policy string) *EnforcerService {
	t.Helper()

	path := filepath.Join(t.TempDir(), "policy.csv")
	require.NoError(t, os.WriteFile(path, []byte(policy), 0o600))

	logger := zerolog.New(io.Discard)
	svc, err := newEnforcerWithAdapter(nil, &logger, fileadapter.NewAdapter(path))
	require.NoError(t, err)

	return svc
}

func TestDecisionService_Check_AllowsRolePolicy(t *testing.T) {
	enforcer := newTestEnforcer(t, "p, doctor, org-1, patient:123, read, dept\n")
	decision := NewDecisionService(enforcer)

	resp, err := decision.Check(t.Context(), &pb.CheckRequest{
		Subject:  &pb.Subject{UserId: "u-1", Tenant: "org-1", Roles: []string{"doctor"}},
		Action:   "read",
		Resource: "patient:123",
	})

	require.NoError(t, err)
	require.True(t, resp.GetAllowed())
	require.Equal(t, "dept", resp.GetDataScopeHint())
}

func TestDecisionService_Check_DeniesMissingPolicy(t *testing.T) {
	enforcer := newTestEnforcer(t, "p, doctor, org-1, patient:123, read, dept\n")
	decision := NewDecisionService(enforcer)

	resp, err := decision.Check(t.Context(), &pb.CheckRequest{
		Subject:  &pb.Subject{UserId: "u-1", Tenant: "org-1", Roles: []string{"doctor"}},
		Action:   "write",
		Resource: "patient:123",
	})

	require.NoError(t, err)
	require.False(t, resp.GetAllowed())
	require.NotEmpty(t, resp.GetReason())
}

func TestDecisionService_BatchCheck_ReturnsPerItemResults(t *testing.T) {
	enforcer := newTestEnforcer(t, "p, doctor, org-1, patient:123, read, dept\n")
	decision := NewDecisionService(enforcer)

	resp, err := decision.BatchCheck(t.Context(), &pb.BatchCheckRequest{
		Subject: &pb.Subject{UserId: "u-1", Tenant: "org-1", Roles: []string{"doctor"}},
		Items: []*pb.CheckItem{
			{Action: "read", Resource: "patient:123"},
			{Action: "write", Resource: "patient:123"},
		},
	})

	require.NoError(t, err)
	require.Len(t, resp.GetResults(), 2)
	require.True(t, resp.GetResults()[0].GetAllowed())
	require.False(t, resp.GetResults()[1].GetAllowed())
}

func TestDecisionService_ListPermissions_DeduplicatesRolePermissions(t *testing.T) {
	enforcer := newTestEnforcer(
		t,
		"p, doctor, org-1, patient:123, read, dept\np, doctor, org-1, patient:123, read, dept\n",
	)
	decision := NewDecisionService(enforcer)

	resp, err := decision.ListPermissions(t.Context(), &pb.ListPermissionsRequest{
		Subject: &pb.Subject{UserId: "u-1", Tenant: "org-1", Roles: []string{"doctor"}},
	})

	require.NoError(t, err)
	require.Len(t, resp.GetPermissions(), 1)
	require.Equal(t, "patient:123", resp.GetPermissions()[0].GetResource())
	require.Equal(t, "read", resp.GetPermissions()[0].GetAction())
}

func TestEnforcerService_AddAndRemovePolicyRule(t *testing.T) {
	enforcer := newTestEnforcer(t, "")

	added, err := enforcer.AddPolicyRule(PTypePolicy, []string{"doctor", "org-1", "patient:123", "read", "dept"})
	require.NoError(t, err)
	require.True(t, added)

	allowed, scope, err := enforcer.EnforceWithDataScope("doctor", "org-1", "patient:123", "read")
	require.NoError(t, err)
	require.True(t, allowed)
	require.Equal(t, "dept", scope)

	removed, err := enforcer.RemovePolicyRule(PTypePolicy, []string{"doctor", "org-1", "patient:123", "read", "dept"})
	require.NoError(t, err)
	require.True(t, removed)
}

func TestEnforcerService_AddPolicyRule_RejectsInvalidPType(t *testing.T) {
	enforcer := newTestEnforcer(t, "")

	added, err := enforcer.AddPolicyRule("x", []string{"doctor", "org-1", "patient:123", "read", "dept"})

	require.ErrorIs(t, err, ErrInvalidPType)
	require.False(t, added)
}

// TestEnforcerService_SeedDefaultsIfEmpty_WritesSuperadminWhenEmpty 验证空策略表会写入 superadmin 通配规则。
func TestEnforcerService_SeedDefaultsIfEmpty_WritesSuperadminWhenEmpty(t *testing.T) {
	enforcer := newTestEnforcer(t, "")

	require.NoError(t, enforcer.SeedDefaultsIfEmpty(t.Context()))

	pCount, _, _ := enforcer.GetPolicyCount()
	require.Equal(t, 1, pCount)

	allowed, scope, err := enforcer.EnforceWithDataScope("role:superadmin", "any-org", "user", "create")
	require.NoError(t, err)
	require.True(t, allowed)
	require.Equal(t, "all", scope)
}

// TestEnforcerService_SeedDefaultsIfEmpty_NoOpWhenNonEmpty 验证已有策略时不会重复写入。
func TestEnforcerService_SeedDefaultsIfEmpty_NoOpWhenNonEmpty(t *testing.T) {
	enforcer := newTestEnforcer(t, "p, doctor, org-1, patient:123, read, dept\n")

	require.NoError(t, enforcer.SeedDefaultsIfEmpty(t.Context()))

	pCount, _, _ := enforcer.GetPolicyCount()
	require.Equal(t, 1, pCount, "已有策略时不应写入 superadmin 种子")

	allowed, _, err := enforcer.EnforceWithDataScope("role:superadmin", "any-org", "user", "create")
	require.NoError(t, err)
	require.False(t, allowed, "未写入种子时 superadmin 不应被允许")
}

// TestEnforcerService_EnforceWithDataScope_WildcardMatching 验证 dom/obj/act 通配查询。
//
// 修复 Phase 5 前 enforcerEx 仅覆盖 act 通配，导致 superadmin 通配策略不生效。
func TestEnforcerService_EnforceWithDataScope_WildcardMatching(t *testing.T) {
	cases := []struct {
		name   string
		policy string
		sub    string
		dom    string
		obj    string
		act    string
		want   bool
		scope  string
	}{
		{
			name:   "全通配（dom * obj * act *）",
			policy: "p, superadmin, *, *, *, all\n",
			sub:    "superadmin",
			dom:    "org-1",
			obj:    "user",
			act:    "create",
			want:   true,
			scope:  "all",
		},
		{
			name:   "dom 通配 + obj 精确 + act 精确",
			policy: "p, superadmin, *, user, create, all\n",
			sub:    "superadmin",
			dom:    "org-1",
			obj:    "user",
			act:    "create",
			want:   true,
			scope:  "all",
		},
		{
			name:   "obj 通配但 act 不匹配",
			policy: "p, doctor, org-1, *, read, dept\n",
			sub:    "doctor",
			dom:    "org-1",
			obj:    "patient:123",
			act:    "write",
			want:   false,
			scope:  "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			enforcer := newTestEnforcer(t, c.policy)

			allowed, scope, err := enforcer.EnforceWithDataScope(c.sub, c.dom, c.obj, c.act)
			require.NoError(t, err)
			require.Equal(t, c.want, allowed)
			require.Equal(t, c.scope, scope)
		})
	}
}
