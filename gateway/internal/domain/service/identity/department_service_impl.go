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

// departmentServiceImpl 部门管理服务实现
type departmentServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      identityConv.Assembler
}

// NewDepartmentService 创建新的部门管理服务实例
func NewDepartmentService(
	identityClient identitycli.IdentityClient,
	assembler identityConv.Assembler,
	logger *hertzZerolog.Logger,
) DepartmentService {
	return &departmentServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// =================================================================
// 5. 部门管理模块 (Department Management)
// =================================================================

func (s *departmentServiceImpl) CreateDepartment(
	ctx context.Context,
	req *identity.CreateDepartmentRequestDTO,
) (*identity.DepartmentResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "创建部门",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Department().ToRPCCreateDeptRequest(req)

			// 调用RPC服务
			return s.identityClient.CreateDepartment(ctx, rpcReq)
		},
		"name", req.Name, "organization_id", req.OrganizationID,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应到HTTP响应
	rpcDept := result.(*identity_srv.Department)
	httpDept := s.assembler.Department().ToHTTPDepartment(rpcDept)

	// 设置成功的基础响应
	httpResp := &identity.DepartmentResponseDTO{
		BaseResp:   s.ResponseBuilder().BuildSuccessResponse(),
		Department: httpDept,
	}

	return httpResp, nil
}

func (s *departmentServiceImpl) GetDepartment(
	ctx context.Context,
	req *identity.GetDepartmentRequestDTO,
) (*identity.DepartmentResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "获取部门信息",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Department().ToRPCGetDeptRequest(req)

			// 调用RPC服务
			return s.identityClient.GetDepartment(ctx, rpcReq)
		},
		"department_id", req.DepartmentID,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应到HTTP响应
	rpcDept := result.(*identity_srv.Department)
	httpDept := s.assembler.Department().ToHTTPDepartment(rpcDept)

	// 设置成功的基础响应
	httpResp := &identity.DepartmentResponseDTO{
		BaseResp:   s.ResponseBuilder().BuildSuccessResponse(),
		Department: httpDept,
	}

	return httpResp, nil
}

func (s *departmentServiceImpl) UpdateDepartment(
	ctx context.Context,
	req *identity.UpdateDepartmentRequestDTO,
) (*identity.DepartmentResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "更新部门信息",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Department().ToRPCUpdateDeptRequest(req)

			// 调用RPC服务
			return s.identityClient.UpdateDepartment(ctx, rpcReq)
		},
		"department_id", req.DepartmentID,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应到HTTP响应
	rpcDept := result.(*identity_srv.Department)
	httpDept := s.assembler.Department().ToHTTPDepartment(rpcDept)

	// 设置成功的基础响应
	httpResp := &identity.DepartmentResponseDTO{
		BaseResp:   s.ResponseBuilder().BuildSuccessResponse(),
		Department: httpDept,
	}

	return httpResp, nil
}

func (s *departmentServiceImpl) DeleteDepartment(
	ctx context.Context,
	req *identity.DeleteDepartmentRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	err := s.ProcessRPCVoidCall(ctx, "删除部门",
		func(ctx context.Context) error {
			// 调用RPC服务
			return s.identityClient.DeleteDepartment(ctx, *req.DepartmentID)
		},
		"department_id", req.DepartmentID,
	)
	if err != nil {
		return nil, err
	}

	// 使用ResponseBuilder构建响应
	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

func (s *departmentServiceImpl) GetOrganizationDepartments(
	ctx context.Context,
	req *identity.GetOrganizationDepartmentsRequestDTO,
) (*identity.GetOrganizationDepartmentsResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "获取组织部门列表",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Department().ToRPCListDeptsRequest(req)

			// 调用RPC服务
			return s.identityClient.GetOrganizationDepartments(ctx, rpcReq)
		},
		"organization_id", req.OrganizationID,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应到HTTP响应
	rpcResp := result.(*identity_srv.GetOrganizationDepartmentsResponse)
	httpResp := s.assembler.Department().ToHTTPListDeptsResponse(rpcResp)

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}
