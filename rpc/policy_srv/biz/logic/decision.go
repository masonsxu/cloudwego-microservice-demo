package logic

import (
	"context"
	"fmt"

	pb "github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/kitex_gen/policy_srv"
)

// DecisionService 权限决策服务
type DecisionService struct {
	enforcer *EnforcerService
}

// NewDecisionService 创建决策服务
func NewDecisionService(enforcer *EnforcerService) *DecisionService {
	return &DecisionService{enforcer: enforcer}
}

// Check 单点权限检查
func (s *DecisionService) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	subject := req.GetSubject()
	if subject == nil {
		return &pb.CheckResponse{Allowed: false, Reason: "subject is nil"}, nil
	}

	userID := "user:" + subject.GetUserId()
	tenant := subject.GetTenant()

	allowed := false
	maxScope := ""

	// 按角色逐个检查（多角色取并集 — 任一角色允许即通过）
	for _, roleCode := range subject.GetRoles() {
		// 使用角色 code 作为 subject（Casbin model 中 p.sub 是 role code）
		ok, scope, err := s.enforcer.EnforceWithDataScope(roleCode, tenant, req.GetResource(), req.GetAction())
		if err != nil {
			return nil, fmt.Errorf("enforce failed: %w", err)
		}

		if ok {
			allowed = true
			if compareDataScope(scope, maxScope) > 0 {
				maxScope = scope
			}
		}
	}

	// 也检查是否有直接的 user-role 绑定（通过 g 关系）
	if !allowed {
		ok, scope, err := s.enforcer.EnforceWithDataScope(userID, tenant, req.GetResource(), req.GetAction())
		if err != nil {
			return nil, fmt.Errorf("enforce user binding failed: %w", err)
		}
		if ok {
			allowed = true
			maxScope = scope
		}
	}

	reason := ""
	if !allowed {
		reason = fmt.Sprintf("denied: user=%s roles=%v action=%s resource=%s tenant=%s",
			subject.GetUserId(), subject.GetRoles(), req.GetAction(), req.GetResource(), tenant)
	}

	return &pb.CheckResponse{
		Allowed:       allowed,
		Reason:        reason,
		DataScopeHint: maxScope,
	}, nil
}

// BatchCheck 批量权限检查
func (s *DecisionService) BatchCheck(ctx context.Context, req *pb.BatchCheckRequest) (*pb.BatchCheckResponse, error) {
	results := make([]*pb.CheckResult, 0, len(req.GetItems()))

	for _, item := range req.GetItems() {
		checkReq := &pb.CheckRequest{
			Subject:            req.GetSubject(),
			Action:             item.GetAction(),
			Resource:           item.GetResource(),
			ResourceAttributes: item.GetResourceAttributes(),
		}

		resp, err := s.Check(ctx, checkReq)
		if err != nil {
			return nil, err
		}

		results = append(results, &pb.CheckResult{
			Action:   item.GetAction(),
			Resource: item.GetResource(),
			Allowed:  resp.GetAllowed(),
			Reason:   resp.GetReason(),
		})
	}

	return &pb.BatchCheckResponse{Results: results}, nil
}

// ListPermissions 查询主体所有权限
func (s *DecisionService) ListPermissions(
	ctx context.Context,
	req *pb.ListPermissionsRequest,
) (*pb.ListPermissionsResponse, error) {
	subject := req.GetSubject()
	if subject == nil {
		return &pb.ListPermissionsResponse{}, nil
	}

	userID := "user:" + subject.GetUserId()
	domains := []string{subject.GetTenant(), "*"}

	perms, _ := s.enforcer.GetPermissionsForUser(userID, domains...)

	// 也获取通过角色继承的权限
	for _, roleCode := range subject.GetRoles() {
		rolePerms, _ := s.enforcer.GetPermissionsForUser(roleCode, domains...)
		perms = append(perms, rolePerms...)
	}

	// 去重
	seen := make(map[string]bool)
	items := make([]*pb.PermissionItem, 0, len(perms))

	for _, p := range perms {
		if len(p) < 5 {
			continue
		}
		// p = [sub, dom, obj, act, data_scope]
		key := p[2] + ":" + p[3]
		if seen[key] {
			continue
		}
		seen[key] = true

		items = append(items, &pb.PermissionItem{
			Resource:  p[2],
			Action:    p[3],
			DataScope: p[4],
		})
	}

	return &pb.ListPermissionsResponse{Permissions: items}, nil
}
