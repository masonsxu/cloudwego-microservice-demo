package identity

import (
	"context"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	identityassembler "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// logoServiceImpl 组织Logo管理服务实现
type logoServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      identityassembler.Assembler
}

// NewLogoService 创建新的Logo管理服务实例
func NewLogoService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) LogoService {
	return &logoServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// =================================================================
// 组织Logo管理模块 (Organization Logo Management)
// =================================================================

func (s *logoServiceImpl) UploadTemporaryLogo(
	ctx context.Context,
	req *identity.UploadTemporaryLogoRequestDTO,
	userID string,
) (*identity.OrganizationLogoResponseDTO, error) {
	// 转换请求
	rpcReq := s.assembler.Logo().ToRPCUploadTemporaryLogoRequest(req, userID)

	// 调用RPC服务
	result, err := s.ProcessRPCCall(ctx, "上传临时Logo",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.UploadTemporaryLogo(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.OrganizationLogo)

	// 转换响应
	httpResp := &identity.OrganizationLogoResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Logo:     s.assembler.Logo().ToHTTPOrganizationLogo(rpcResp),
	}

	return httpResp, nil
}

func (s *logoServiceImpl) GetOrganizationLogo(
	ctx context.Context,
	req *identity.GetOrganizationLogoRequestDTO,
) (*identity.OrganizationLogoResponseDTO, error) {
	// 转换请求
	rpcReq := s.assembler.Logo().ToRPCGetOrganizationLogoRequest(req)

	// 调用RPC服务
	result, err := s.ProcessRPCCall(ctx, "获取Logo信息",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.GetOrganizationLogo(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.OrganizationLogo)

	// 转换响应
	httpResp := &identity.OrganizationLogoResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
		Logo:     s.assembler.Logo().ToHTTPOrganizationLogo(rpcResp),
	}

	return httpResp, nil
}

func (s *logoServiceImpl) DeleteOrganizationLogo(
	ctx context.Context,
	req *identity.DeleteOrganizationLogoRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	// 转换请求
	rpcReq := s.assembler.Logo().ToRPCDeleteOrganizationLogoRequest(req)

	// 调用RPC服务
	_, err := s.ProcessRPCCall(ctx, "删除Logo",
		func(ctx context.Context) (interface{}, error) {
			return nil, s.identityClient.DeleteOrganizationLogo(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	// 构建响应
	httpResp := &http_base.OperationStatusResponseDTO{
		BaseResp: s.ResponseBuilder().BuildSuccessResponse(),
	}

	return httpResp, nil
}

func (s *logoServiceImpl) BindLogoToOrganization(
	ctx context.Context,
	req *identity.BindLogoToOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	// 转换请求
	rpcReq := s.assembler.Logo().ToRPCBindLogoToOrganizationRequest(req)

	// 调用RPC服务
	result, err := s.ProcessRPCCall(ctx, "绑定Logo到组织",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.BindLogoToOrganization(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcLogo := result.(*identity_srv.OrganizationLogo)

	// 如果绑定成功，获取更新后的组织信息
	orgResult, err := s.ProcessRPCCall(ctx, "获取组织信息",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.GetOrganization(ctx, &identity_srv.GetOrganizationRequest{
				OrganizationID: rpcLogo.BoundOrganizationID,
			})
		},
	)
	if err != nil {
		return nil, err
	}

	rpcOrg := orgResult.(*identity_srv.Organization)

	// 转换响应
	httpResp := &identity.OrganizationResponseDTO{
		BaseResp:     s.ResponseBuilder().BuildSuccessResponse(),
		Organization: s.assembler.Organization().ToHTTPOrganization(rpcOrg),
	}

	return httpResp, nil
}
