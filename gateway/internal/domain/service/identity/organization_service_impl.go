package identity

import (
	"context"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	identityConv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// organizationServiceImpl 组织架构管理服务实现
type organizationServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      identityConv.Assembler
}

// NewOrganizationService 创建新的组织架构管理服务实例
func NewOrganizationService(
	identityClient identitycli.IdentityClient,
	assembler identityConv.Assembler,
	logger *hertzZerolog.Logger,
) OrganizationService {
	return &organizationServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// =================================================================
// 4. 组织架构管理模块 (Organization Management)
// =================================================================

func (s *organizationServiceImpl) CreateOrganization(
	ctx context.Context,
	req *identity.CreateOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "创建组织",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Organization().ToRPCCreateOrgRequest(req)
			return s.identityClient.CreateOrganization(ctx, rpcReq)
		},
		"name", req.Name, "type", req.FacilityType,
	)
	if err != nil {
		return nil, err
	}

	rpcOrg := result.(*identity_srv.Organization)
	httpOrg := s.assembler.Organization().ToHTTPOrganization(rpcOrg)

	httpResp := &identity.OrganizationResponseDTO{
		BaseResp:     s.ResponseBuilder().BuildSuccessResponse(),
		Organization: httpOrg,
	}

	return httpResp, nil
}

func (s *organizationServiceImpl) GetOrganization(
	ctx context.Context,
	req *identity.GetOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取组织信息",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Organization().ToRPCGetOrgRequest(req)
			return s.identityClient.GetOrganization(ctx, rpcReq)
		},
		"organization_id", req.OrganizationID,
	)
	if err != nil {
		return nil, err
	}

	rpcOrg := result.(*identity_srv.Organization)
	httpOrg := s.assembler.Organization().ToHTTPOrganization(rpcOrg)

	httpResp := &identity.OrganizationResponseDTO{
		BaseResp:     s.ResponseBuilder().BuildSuccessResponse(),
		Organization: httpOrg,
	}

	return httpResp, nil
}

func (s *organizationServiceImpl) UpdateOrganization(
	ctx context.Context,
	req *identity.UpdateOrganizationRequestDTO,
) (*identity.OrganizationResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "更新组织信息",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Organization().ToRPCUpdateOrgRequest(req)
			return s.identityClient.UpdateOrganization(ctx, rpcReq)
		},
		"organization_id", req.OrganizationID,
	)
	if err != nil {
		return nil, err
	}

	rpcOrg := result.(*identity_srv.Organization)
	httpOrg := s.assembler.Organization().ToHTTPOrganization(rpcOrg)

	httpResp := &identity.OrganizationResponseDTO{
		BaseResp:     s.ResponseBuilder().BuildSuccessResponse(),
		Organization: httpOrg,
	}

	return httpResp, nil
}

func (s *organizationServiceImpl) DeleteOrganization(
	ctx context.Context,
	req *identity.DeleteOrganizationRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	err := s.ProcessRPCVoidCall(ctx, "删除组织",
		func(ctx context.Context) error {
			return s.identityClient.DeleteOrganization(ctx, *req.OrganizationID)
		},
		"organization_id", req.OrganizationID,
	)
	if err != nil {
		return nil, err
	}

	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

func (s *organizationServiceImpl) ListOrganizations(
	ctx context.Context,
	req *identity.ListOrganizationsRequestDTO,
) (*identity.ListOrganizationsResponseDTO, error) {
	result, err := s.ProcessRPCCall(ctx, "获取组织列表",
		func(ctx context.Context) (interface{}, error) {
			rpcReq := s.assembler.Organization().ToRPCListOrgsRequest(req)
			return s.identityClient.ListOrganizations(ctx, rpcReq)
		},
		"parent_id", req.ParentID, "page", req.Page,
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.ListOrganizationsResponse)
	httpResp := s.assembler.Organization().ToHTTPListOrgsResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}
