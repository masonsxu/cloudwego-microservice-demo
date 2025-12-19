package identity

import (
	"context"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	identityConv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// membershipServiceImpl 成员关系管理服务实现
type membershipServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      identityConv.Assembler
}

// NewMembershipService 创建新的成员关系管理服务实例
func NewMembershipService(
	identityClient identitycli.IdentityClient,
	assembler identityConv.Assembler,
	logger *hertzZerolog.Logger,
) MembershipService {
	return &membershipServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// =================================================================
// 3. 成员关系管理模块 (Membership Management)
// =================================================================

func (s *membershipServiceImpl) GetUserMemberships(
	ctx context.Context,
	req *identity.GetUserMembershipsRequestDTO,
) (*identity.GetUserMembershipsResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取用户成员关系",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Membership().ToRPCGetUserMembershipsRequest(req)
			return s.identityClient.GetUserMemberships(ctx, rpcReq)
		},
		"user_id", req.UserID,
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.GetUserMembershipsResponse)
	httpResp := s.assembler.Membership().ToHTTPGetUserMembershipsResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

// GetPrimaryMembership 获取用户主成员关系 - 获取用户在当前组织中的主成员关系
func (s *membershipServiceImpl) GetPrimaryMembership(
	ctx context.Context,
	req *identity.GetPrimaryMembershipRequestDTO,
) (*identity.UserMembershipResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取用户主成员关系",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.GetPrimaryMembership(ctx, *req.UserID)
		},
		"user_id", req.UserID,
	)
	if err != nil {
		return nil, err
	}

	rpcUserMembership := result.(*identity_srv.UserMembership)
	httpUserMembership := s.assembler.Membership().ToHTTPUserMembership(rpcUserMembership)

	httpResp := &identity.UserMembershipResponseDTO{
		BaseResp:   s.ResponseBuilder().BuildSuccessResponse(),
		Membership: httpUserMembership,
	}

	return httpResp, nil
}
