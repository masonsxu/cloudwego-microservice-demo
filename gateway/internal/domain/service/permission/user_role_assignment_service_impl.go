package permission

import (
	"context"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	permissionconv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// userRoleAssignmentServiceImpl
type userRoleAssignmentServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      permissionconv.Assembler
}

// NewUserRoleAssignmentService
func NewUserRoleAssignmentService(
	identityClient identitycli.IdentityClient,
	assembler permissionconv.Assembler,
	logger *hertzZerolog.Logger,
) UserRoleAssignmentService {
	return &userRoleAssignmentServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

func (s *userRoleAssignmentServiceImpl) GetLastUserRoleAssignment(
	ctx context.Context,
	userID string,
) (*permission.AssignRoleToUserResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取最新用户角色分配",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.GetLastUserRoleAssignment(ctx, userID)
		},
		"userID", userID,
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.UserRoleAssignment)
	httpResp := s.assembler.UserRole().
		ToHTTPAssignRoleToUserResponse(&identity_srv.UserRoleAssignmentResponse{AssignmentID: rpcResp.Id})
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

func (s *userRoleAssignmentServiceImpl) ListUserRoleAssignments(
	ctx context.Context,
	req *permission.UserRoleQueryRequestDTO,
) (*permission.UserRoleListResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "列出用户角色分配",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.UserRole().ToRPCUserRoleQueryRequest(req)
			return s.identityClient.ListUserRoleAssignments(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.UserRoleListResponse)
	httpResp := s.assembler.UserRole().ToHTTPUserRoleListResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

func (s *userRoleAssignmentServiceImpl) GetUsersByRole(
	ctx context.Context,
	req *permission.GetUsersByRoleRequestDTO,
) (*permission.GetUsersByRoleResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取角色下所有用户",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.UserRole().ToRPCGetUsersByRoleRequest(req)
			return s.identityClient.GetUsersByRole(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.GetUsersByRoleResponse)
	httpResp := s.assembler.UserRole().ToHTTPGetUsersByRoleResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

func (s *userRoleAssignmentServiceImpl) BatchBindUsersToRole(
	ctx context.Context,
	operatorID string,
	req *permission.BatchBindUsersToRoleRequestDTO,
) (*permission.BatchBindUsersToRoleResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "批量绑定用户到角色",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.UserRole().ToRPCBatchBindUsersToRoleRequest(operatorID, req)
			return s.identityClient.BatchBindUsersToRole(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.BatchBindUsersToRoleResponse)
	httpResp := s.assembler.UserRole().ToHTTPBatchBindUsersToRoleResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}
