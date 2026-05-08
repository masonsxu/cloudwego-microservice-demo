package iamclient

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	policy "github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/kitex_gen/policy_srv"
)

// fakePolicyClient 实现 policyservice.Client，用于在不启 RPC 真实服务的情况下
// 验证 iamclient 的请求构造、缓存命中、错误传播等逻辑。
type fakePolicyClient struct {
	checkFn      func(ctx context.Context, req *policy.CheckRequest) (*policy.CheckResponse, error)
	checkCalls   int
	lastCheckReq *policy.CheckRequest
	lastCheckCtx context.Context //nolint:containedctx // 测试桩需要回溯传入的 ctx
}

func (f *fakePolicyClient) Check(
	ctx context.Context,
	req *policy.CheckRequest,
	_ ...callopt.Option,
) (*policy.CheckResponse, error) {
	f.checkCalls++
	f.lastCheckReq = req
	f.lastCheckCtx = ctx

	return f.checkFn(ctx, req)
}

func (f *fakePolicyClient) BatchCheck(
	context.Context, *policy.BatchCheckRequest, ...callopt.Option,
) (*policy.BatchCheckResponse, error) {
	return nil, errors.New("not used")
}

func (f *fakePolicyClient) ListPermissions(
	context.Context, *policy.ListPermissionsRequest, ...callopt.Option,
) (*policy.ListPermissionsResponse, error) {
	return nil, errors.New("not used")
}

func (f *fakePolicyClient) UpsertPolicy(
	context.Context, *policy.UpsertPolicyRequest, ...callopt.Option,
) (*policy.UpsertPolicyResponse, error) {
	return nil, errors.New("not used")
}

func (f *fakePolicyClient) DeletePolicy(
	context.Context, *policy.DeletePolicyRequest, ...callopt.Option,
) (*policy.DeletePolicyResponse, error) {
	return nil, errors.New("not used")
}

func (f *fakePolicyClient) ReloadPolicies(
	context.Context, *policy.ReloadPoliciesRequest, ...callopt.Option,
) (*policy.ReloadPoliciesResponse, error) {
	return nil, errors.New("not used")
}

func newTestClient(t *testing.T, fake *fakePolicyClient) *Client {
	t.Helper()

	cache, err := newDecisionCache(100, 30*time.Second)
	require.NoError(t, err)

	return &Client{policy: fake, cache: cache}
}

func newTestSubject(c *Client) *Subject {
	return &Subject{
		UserID:   "u1",
		TenantID: "org-1",
		Roles:    []string{"doctor"},
		Jti:      "jti-1",
		client:   c,
	}
}

func TestCheck_AllowedRequestPayload(t *testing.T) {
	fake := &fakePolicyClient{
		checkFn: func(_ context.Context, _ *policy.CheckRequest) (*policy.CheckResponse, error) {
			return &policy.CheckResponse{Allowed: true, DataScopeHint: "self"}, nil
		},
	}
	c := newTestClient(t, fake)
	s := newTestSubject(c)

	d, err := s.Check(context.Background(), "read", "patient:7",
		WithResourceAttr("department_id", "d1"),
	)
	require.NoError(t, err)
	assert.True(t, d.Allowed)
	assert.Equal(t, "self", d.DataScopeHint)

	require.NotNil(t, fake.lastCheckReq)
	assert.Equal(t, "u1", fake.lastCheckReq.GetSubject().GetUserId())
	assert.Equal(t, "org-1", fake.lastCheckReq.GetSubject().GetTenant())
	assert.Equal(t, []string{"doctor"}, fake.lastCheckReq.GetSubject().GetRoles())
	assert.Equal(t, "read", fake.lastCheckReq.GetAction())
	assert.Equal(t, "patient:7", fake.lastCheckReq.GetResource())
	assert.Equal(t, map[string]string{"department_id": "d1"}, fake.lastCheckReq.GetResourceAttributes())
}

func TestCheck_DeniedReturnsDecision(t *testing.T) {
	fake := &fakePolicyClient{
		checkFn: func(_ context.Context, _ *policy.CheckRequest) (*policy.CheckResponse, error) {
			return &policy.CheckResponse{Allowed: false, Reason: "no role"}, nil
		},
	}
	c := newTestClient(t, fake)
	s := newTestSubject(c)

	d, err := s.Check(context.Background(), "delete", "patient:7")
	require.NoError(t, err)
	assert.False(t, d.Allowed)
	assert.Equal(t, "no role", d.Reason)
}

func TestCheck_RPCErrorPropagated(t *testing.T) {
	rpcErr := errors.New("transport closed")
	fake := &fakePolicyClient{
		checkFn: func(context.Context, *policy.CheckRequest) (*policy.CheckResponse, error) {
			return nil, rpcErr
		},
	}
	c := newTestClient(t, fake)
	s := newTestSubject(c)

	d, err := s.Check(context.Background(), "read", "patient:7")
	assert.Nil(t, d)
	assert.ErrorIs(t, err, rpcErr)
}

func TestCheck_CacheHitAvoidsSecondRPC(t *testing.T) {
	fake := &fakePolicyClient{
		checkFn: func(context.Context, *policy.CheckRequest) (*policy.CheckResponse, error) {
			return &policy.CheckResponse{Allowed: true}, nil
		},
	}
	c := newTestClient(t, fake)
	s := newTestSubject(c)

	for range 5 {
		d, err := s.Check(context.Background(), "read", "patient:7")
		require.NoError(t, err)
		assert.True(t, d.Allowed)
	}

	assert.Equal(t, 1, fake.checkCalls, "subsequent calls should hit cache")
}

func TestCheck_WithoutCacheBypasses(t *testing.T) {
	fake := &fakePolicyClient{
		checkFn: func(context.Context, *policy.CheckRequest) (*policy.CheckResponse, error) {
			return &policy.CheckResponse{Allowed: true}, nil
		},
	}
	c := newTestClient(t, fake)
	s := newTestSubject(c)

	for range 3 {
		_, err := s.Check(context.Background(), "read", "patient:7", WithoutCache())
		require.NoError(t, err)
	}

	assert.Equal(t, 3, fake.checkCalls)
}

func TestCheck_DifferentResourceProducesSeparateCacheEntries(t *testing.T) {
	fake := &fakePolicyClient{
		checkFn: func(_ context.Context, req *policy.CheckRequest) (*policy.CheckResponse, error) {
			return &policy.CheckResponse{Allowed: true, Reason: req.GetResource()}, nil
		},
	}
	c := newTestClient(t, fake)
	s := newTestSubject(c)

	d1, _ := s.Check(context.Background(), "read", "patient:1")
	d2, _ := s.Check(context.Background(), "read", "patient:2")

	assert.Equal(t, "patient:1", d1.Reason)
	assert.Equal(t, "patient:2", d2.Reason)
	assert.Equal(t, 2, fake.checkCalls)

	// Re-query: both served from cache.
	_, _ = s.Check(context.Background(), "read", "patient:1")
	_, _ = s.Check(context.Background(), "read", "patient:2")

	assert.Equal(t, 2, fake.checkCalls)
}

func TestMustCheck_Allowed(t *testing.T) {
	fake := &fakePolicyClient{
		checkFn: func(context.Context, *policy.CheckRequest) (*policy.CheckResponse, error) {
			return &policy.CheckResponse{Allowed: true}, nil
		},
	}
	s := newTestSubject(newTestClient(t, fake))

	err := s.MustCheck(context.Background(), "read", "patient:7")
	assert.NoError(t, err)
}

func TestMustCheck_DeniedReturnsTypedError(t *testing.T) {
	fake := &fakePolicyClient{
		checkFn: func(context.Context, *policy.CheckRequest) (*policy.CheckResponse, error) {
			return &policy.CheckResponse{Allowed: false, Reason: "no role"}, nil
		},
	}
	s := newTestSubject(newTestClient(t, fake))

	err := s.MustCheck(context.Background(), "delete", "patient:7")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPermissionDenied)

	var pde *PermissionDeniedError
	require.ErrorAs(t, err, &pde)
	assert.Equal(t, "delete", pde.Action)
	assert.Equal(t, "patient:7", pde.Resource)
	assert.Equal(t, "no role", pde.Reason)
}

func TestCheck_NilSubject(t *testing.T) {
	var s *Subject

	d, err := s.Check(context.Background(), "read", "p:1")
	assert.Nil(t, d)
	assert.Error(t, err)
}

func TestCheck_SubjectWithoutClient(t *testing.T) {
	s := &Subject{UserID: "u1"} // no client field set
	d, err := s.Check(context.Background(), "read", "p:1")
	assert.Nil(t, d)
	assert.Error(t, err)
}
