package iamclient

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	policy "github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/kitex_gen/policy_srv"
)

// ErrPermissionDenied 表示 PDP 决策为拒绝。
//
// 业务侧应转成 403 Forbidden 返回。可以用 errors.Is(err, ErrPermissionDenied)
// 判断；具体拒绝原因（PDP 返回的 Reason）通过 PermissionDeniedError.Reason 读取。
var ErrPermissionDenied = errors.New("iamclient: permission denied")

// PermissionDeniedError 携带 PDP 决策的详细拒绝原因。
type PermissionDeniedError struct {
	Action   string
	Resource string
	Reason   string
}

func (e *PermissionDeniedError) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("%s: action=%s resource=%s", ErrPermissionDenied, e.Action, e.Resource)
	}

	return fmt.Sprintf("%s: action=%s resource=%s reason=%s",
		ErrPermissionDenied, e.Action, e.Resource, e.Reason)
}

// Unwrap 让 errors.Is(err, ErrPermissionDenied) 成立。
func (e *PermissionDeniedError) Unwrap() error { return ErrPermissionDenied }

// CheckOpt 是 Check / MustCheck 的可选参数。
type CheckOpt func(*checkOptions)

type checkOptions struct {
	resourceAttrs map[string]string
	skipCache     bool
}

// WithResourceAttr 给当前决策注入资源属性（owner_id、department_id 等），供
// PDP 做 ABAC 决策时使用。
func WithResourceAttr(key, value string) CheckOpt {
	return func(o *checkOptions) {
		if o.resourceAttrs == nil {
			o.resourceAttrs = make(map[string]string)
		}

		o.resourceAttrs[key] = value
	}
}

// WithoutCache 强制本次决策跳过本地缓存（如调试场景）。
func WithoutCache() CheckOpt {
	return func(o *checkOptions) { o.skipCache = true }
}

// Decision 是 PDP 决策结果。
type Decision struct {
	Allowed       bool
	Reason        string
	DataScopeHint string
}

// Check 询问 PDP "subject 能否对 resource 做 action"。
//
// 行为：
//   - 优先查本地 LRU 缓存（key = jti+action+resource+attrs，TTL 见 Config）。
//   - 未命中则调 policy_srv 的 Check RPC，并写回缓存。
//   - 网络/RPC 错误直接返回 err；不做"拒绝即默认放行"的危险兜底。
//
// 注意 Subject 必须由 Client.SubjectFromHeader / SubjectFromContext 创建，
// 否则 Subject.client 为 nil 会 panic（这是设计意图——禁止手工伪造 Subject 绕过）。
func (s *Subject) Check(ctx context.Context, action, resource string, opts ...CheckOpt) (*Decision, error) {
	if s == nil {
		return nil, errors.New("iamclient: nil subject")
	}

	if s.client == nil {
		return nil, errors.New("iamclient: subject was not created by iamclient.Client")
	}

	o := &checkOptions{}
	for _, opt := range opts {
		opt(o)
	}

	cacheKey := s.cacheKey(action, resource, o.resourceAttrs)

	if !o.skipCache && s.client.cache != nil {
		if d, ok := s.client.cache.get(cacheKey); ok {
			return d, nil
		}
	}

	req := &policy.CheckRequest{
		Subject: &policy.Subject{
			UserId: s.UserID,
			Tenant: s.TenantID,
			Roles:  s.Roles,
		},
		Action:             action,
		Resource:           resource,
		ResourceAttributes: o.resourceAttrs,
	}

	resp, err := s.client.policy.Check(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("iamclient: call policy_srv: %w", err)
	}

	d := &Decision{
		Allowed:       resp.GetAllowed(),
		Reason:        resp.GetReason(),
		DataScopeHint: resp.GetDataScopeHint(),
	}

	if !o.skipCache && s.client.cache != nil {
		s.client.cache.set(cacheKey, d)
	}

	return d, nil
}

// MustCheck 是 Check 的便捷封装：
//   - 决策允许返回 nil；
//   - 决策拒绝返回 *PermissionDeniedError（errors.Is(err, ErrPermissionDenied) 为真）；
//   - 网络/RPC 错误原样返回。
func (s *Subject) MustCheck(ctx context.Context, action, resource string, opts ...CheckOpt) error {
	d, err := s.Check(ctx, action, resource, opts...)
	if err != nil {
		return err
	}

	if !d.Allowed {
		return &PermissionDeniedError{
			Action:   action,
			Resource: resource,
			Reason:   d.Reason,
		}
	}

	return nil
}

// cacheKey 把决策维度拼成缓存 key。
//
// 设计：
//   - jti 不存在时 fallback 到 user_id + roles 拼装（roles 排序后参与，避免顺序差异打散缓存）；
//   - resourceAttrs 按 key 排序后参与，相同输入得到相同 key。
func (s *Subject) cacheKey(action, resource string, attrs map[string]string) string {
	var sb strings.Builder

	sb.Grow(64 + len(action) + len(resource))

	if s.Jti != "" {
		sb.WriteString(s.Jti)
	} else {
		sb.WriteString(s.UserID)
		sb.WriteByte('|')

		roles := make([]string, len(s.Roles))
		copy(roles, s.Roles)
		sort.Strings(roles)
		sb.WriteString(strings.Join(roles, ","))
	}

	sb.WriteByte('#')
	sb.WriteString(s.TenantID)
	sb.WriteByte('#')
	sb.WriteString(action)
	sb.WriteByte('#')
	sb.WriteString(resource)

	if len(attrs) > 0 {
		keys := make([]string, 0, len(attrs))
		for k := range attrs {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		sb.WriteByte('?')

		for i, k := range keys {
			if i > 0 {
				sb.WriteByte('&')
			}

			sb.WriteString(k)
			sb.WriteByte('=')
			sb.WriteString(attrs[k])
		}
	}

	return sb.String()
}
