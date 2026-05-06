package main

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/biz/logic"
	pb "github.com/masonsxu/cloudwego-microservice-demo/rpc/policy-srv/kitex_gen/policy_srv"
)

// PolicyServiceImpl RPC 服务实现
type PolicyServiceImpl struct {
	decision *logic.DecisionService
	enforcer *logic.EnforcerService
}

// NewPolicyServiceImpl 创建 RPC 服务实现
func NewPolicyServiceImpl(decision *logic.DecisionService, enforcer *logic.EnforcerService) *PolicyServiceImpl {
	return &PolicyServiceImpl{
		decision: decision,
		enforcer: enforcer,
	}
}

// Check 单点权限决策
func (s *PolicyServiceImpl) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	return s.decision.Check(ctx, req)
}

// BatchCheck 批量权限决策
func (s *PolicyServiceImpl) BatchCheck(ctx context.Context, req *pb.BatchCheckRequest) (*pb.BatchCheckResponse, error) {
	return s.decision.BatchCheck(ctx, req)
}

// ListPermissions 查询主体所有权限
func (s *PolicyServiceImpl) ListPermissions(
	ctx context.Context,
	req *pb.ListPermissionsRequest,
) (*pb.ListPermissionsResponse, error) {
	return s.decision.ListPermissions(ctx, req)
}

// UpsertPolicy 策略创建/更新
func (s *PolicyServiceImpl) UpsertPolicy(
	ctx context.Context,
	req *pb.UpsertPolicyRequest,
) (*pb.UpsertPolicyResponse, error) {
	ok, err := s.enforcer.AddPolicyRule(req.GetPtype(), req.GetRule())
	if err != nil {
		return &pb.UpsertPolicyResponse{Success: false}, err
	}
	return &pb.UpsertPolicyResponse{Success: ok}, nil
}

// DeletePolicy 策略删除
func (s *PolicyServiceImpl) DeletePolicy(
	ctx context.Context,
	req *pb.DeletePolicyRequest,
) (*pb.DeletePolicyResponse, error) {
	ok, err := s.enforcer.RemovePolicyRule(req.GetPtype(), req.GetRule())
	if err != nil {
		return &pb.DeletePolicyResponse{Success: false}, err
	}
	return &pb.DeletePolicyResponse{Success: ok}, nil
}

// ReloadPolicies 策略重载
func (s *PolicyServiceImpl) ReloadPolicies(
	ctx context.Context,
	req *pb.ReloadPoliciesRequest,
) (*pb.ReloadPoliciesResponse, error) {
	if err := s.enforcer.ReloadPolicy(ctx); err != nil {
		return &pb.ReloadPoliciesResponse{Success: false}, err
	}

	p, g, g2 := s.enforcer.GetPolicyCount()
	return &pb.ReloadPoliciesResponse{
		Success:              true,
		PolicyCount:          int32(p),
		GroupingPolicyCount:  int32(g),
		RoleInheritanceCount: int32(g2),
	}, nil
}
