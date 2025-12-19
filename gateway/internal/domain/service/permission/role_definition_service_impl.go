package permission

import (
	"context"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	permissionconv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// roleDefinitionServiceImpl
type roleDefinitionServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      permissionconv.Assembler
}

// NewRoleDefinitionService
func NewRoleDefinitionService(
	identityClient identitycli.IdentityClient,
	assembler permissionconv.Assembler,
	logger *hertzZerolog.Logger,
) RoleDefinitionService {
	return &roleDefinitionServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

func (s *roleDefinitionServiceImpl) CreateRoleDefinition(
	ctx context.Context,
	req *permission.RoleDefinitionCreateRequestDTO,
) (*permission.RoleDefinitionCreateResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "创建角色定义",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Role().ToRPCRoleDefinitionCreateRequest(req)
			return s.identityClient.CreateRoleDefinition(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.RoleDefinition)

	httpRoleDefinition := s.assembler.Role().ToHTTPRoleDefinition(rpcResp)

	httpResp := &permission.RoleDefinitionCreateResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Role:     httpRoleDefinition,
	}

	return httpResp, nil
}

func (s *roleDefinitionServiceImpl) UpdateRoleDefinition(
	ctx context.Context,
	req *permission.RoleDefinitionUpdateRequestDTO,
) (*permission.RoleDefinitionUpdateResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "更新角色定义",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Role().ToRPCRoleDefinitionUpdateRequest(req)
			return s.identityClient.UpdateRoleDefinition(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.RoleDefinition)

	httpRoleDefinition := s.assembler.Role().ToHTTPRoleDefinition(rpcResp)

	httpResp := &permission.RoleDefinitionUpdateResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Role:     httpRoleDefinition,
	}

	return httpResp, nil
}

func (s *roleDefinitionServiceImpl) DeleteRoleDefinition(
	ctx context.Context,
	roleID string,
) (*http_base.OperationStatusResponseDTO, error) {
	err := s.ProcessRPCVoidCall(ctx, "删除角色定义",
		func(ctx context.Context) error {
			return s.identityClient.DeleteRoleDefinition(ctx, roleID)
		},
		"roleID", roleID,
	)
	if err != nil {
		return nil, err
	}

	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

func (s *roleDefinitionServiceImpl) GetRoleDefinition(
	ctx context.Context,
	req *permission.RoleDefinitionGetRequestDTO,
) (*permission.RoleDefinitionGetResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取角色定义",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.GetRoleDefinition(ctx, *req.RoleID)
		},
		"roleID", *req.RoleID,
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.RoleDefinition)

	httpRoleDefinition := s.assembler.Role().ToHTTPRoleDefinition(rpcResp)

	httpResp := &permission.RoleDefinitionGetResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Role:     httpRoleDefinition,
	}

	return httpResp, nil
}

func (s *roleDefinitionServiceImpl) ListRoleDefinitions(
	ctx context.Context,
	req *permission.RoleDefinitionQueryRequestDTO,
) (*permission.RoleDefinitionListResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "列出角色定义",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Role().ToRPCRoleDefinitionQueryRequest(req)
			return s.identityClient.ListRoleDefinitions(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.RoleDefinitionListResponse)

	httpResp := s.assembler.Role().ToHTTPRoleDefinitionListResponse(rpcResp)
	if httpResp == nil {
		return nil, nil
	}

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}
