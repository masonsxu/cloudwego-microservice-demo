package identity

import (
	"context"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	identityConv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// authServiceImpl 统一身份管理服务实现
type authServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      identityConv.Assembler
}

// NewAuthService 创建新的身份服务实例
func NewAuthService(
	identityClient identitycli.IdentityClient,
	assembler identityConv.Assembler,
	logger *hertzZerolog.Logger,
) AuthService {
	return &authServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// =================================================================
// 1. 身份认证模块 (Authentication)
// =================================================================

func (s *authServiceImpl) Login(
	ctx context.Context,
	req *identity.LoginRequestDTO,
) (*identity.LoginResponseDTO, Permission, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "用户登录",
		func(ctx context.Context) (interface{}, error) {
			// 转换HTTP请求到RPC请求
			rpcReq := s.assembler.Auth().ToRPCLoginRequest(req)

			// 调用RPC服务
			return s.identityClient.Login(ctx, rpcReq)
		},
		"username", req.Username,
	)
	if err != nil {
		return nil, "", err
	}

	// 转换RPC响应到HTTP响应
	rpcResp := result.(*identity_srv.LoginResponse)
	httpResp := s.assembler.Auth().ToHTTPLoginResponse(rpcResp)

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	// 处理用户成员关系信息
	if rpcResp.Memberships != nil {
		httpResp.Memberships = s.assembler.Membership().ToHTTPUserMemberships(rpcResp.Memberships)
	} else {
		httpResp.Memberships = []*identity.UserMembershipDTO{}
	}

	// 构建权限信息 - 从RPC响应中直接获取
	if rpcResp.MenuTree != nil {
		httpResp.MenuTree = s.assembler.Auth().ToHTTPMenuTree(rpcResp.MenuTree)
	} else {
		httpResp.MenuTree = []*permission.MenuNodeDTO{}
	}

	if len(rpcResp.RoleIDs) > 0 {
		httpResp.RoleIDs = make([]string, len(rpcResp.RoleIDs))
		copy(httpResp.RoleIDs, rpcResp.RoleIDs)
	} else {
		httpResp.RoleIDs = []string{}
	}

	// 处理菜单权限信息
	if rpcResp.Permissions != nil {
		httpResp.Permissions = s.assembler.Auth().ToHTTPMenuPermissions(rpcResp.Permissions)
	} else {
		httpResp.Permissions = []*permission.MenuPermissionDTO{}
	}

	// 处理角色详情信息
	if rpcResp.RoleDetails != nil {
		httpResp.Roles = s.assembler.Auth().ToHTTPRoleInfos(rpcResp.RoleDetails)
	} else {
		httpResp.Roles = []*identity.RoleInfoDTO{}
	}

	// 处理权限信息（现在是枚举类型）- 提取最高权限级别返回
	if len(rpcResp.Permissions) > 0 &&
		rpcResp.Permissions[0].Permission != nil {
		permLevel := *rpcResp.Permissions[0].Permission
		// 使用枚举的字符串表示
		return httpResp, Permission(permLevel.String()), nil
	}

	return httpResp, "", nil
}

func (s *authServiceImpl) ChangePassword(
	ctx context.Context,
	req *identity.ChangePasswordRequestDTO,
	userID string,
) (*http_base.OperationStatusResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	err := s.ProcessRPCVoidCall(ctx, "修改用户密码",
		func(ctx context.Context) error {
			// 转换为RPC请求
			rpcReq := s.assembler.Auth().ToRPCChangePasswordRequest(req)
			rpcReq.UserID = &userID

			// 调用RPC服务
			return s.identityClient.ChangePassword(ctx, rpcReq)
		},
		"user_id", userID,
	)
	if err != nil {
		return nil, err
	}

	// 使用ResponseBuilder构建响应
	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

func (s *authServiceImpl) ResetPassword(
	ctx context.Context,
	req *identity.ResetPasswordRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	err := s.ProcessRPCVoidCall(ctx, "重置密码",
		func(ctx context.Context) error {
			// 转换为RPC请求
			rpcReq := s.assembler.Auth().ToRPCResetPasswordRequest(req)

			// 调用RPC服务
			return s.identityClient.ResetPassword(ctx, rpcReq)
		},
		"user_id", req.UserID,
	)
	if err != nil {
		return nil, err
	}

	// 使用ResponseBuilder构建响应
	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

func (s *authServiceImpl) ForcePasswordChange(
	ctx context.Context,
	req *identity.ForcePasswordChangeRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	err := s.ProcessRPCVoidCall(ctx, "强制修改密码",
		func(ctx context.Context) error {
			// 转换为RPC请求
			rpcReq := s.assembler.Auth().ToRPCForcePasswordChangeRequest(req)

			// 调用RPC服务
			return s.identityClient.ForcePasswordChange(ctx, rpcReq)
		},
		"user_id", req.UserID, "reason", req.Reason,
	)
	if err != nil {
		return nil, err
	}

	// 使用ResponseBuilder构建响应
	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}
