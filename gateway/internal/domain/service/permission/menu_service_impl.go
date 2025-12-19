package permission

import (
	"context"
	"fmt"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	permissionConv "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// menuServiceImpl implements the MenuService interface.
type menuServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      permissionConv.Assembler
}

// NewMenuService creates a new menu service instance.
func NewMenuService(
	identityClient identitycli.IdentityClient,
	assembler permissionConv.Assembler,
	logger *hertzZerolog.Logger,
) MenuService {
	return &menuServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// UploadMenu 上传菜单配置到权限服务
func (s *menuServiceImpl) UploadMenu(
	ctx context.Context,
	req *permission.UploadMenuRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	err := s.ProcessRPCVoidCall(ctx, "上传菜单配置",
		func(ctx context.Context) error {
			// 转换为RPC请求
			rpcReq := s.assembler.Menu().ToRPCUploadMenuRequest(req)
			if rpcReq == nil {
				return fmt.Errorf("转换上传请求失败")
			}

			// 调用RPC服务
			return s.identityClient.UploadMenu(ctx, rpcReq)
		},
		"yaml_size", len(req.MenuFile),
	)
	if err != nil {
		return nil, err
	}

	// 使用ResponseBuilder构建响应
	return s.ResponseBuilder().BuildOperationStatusResponse(), nil
}

// GetMenuTree 获取菜单树结构
func (s *menuServiceImpl) GetMenuTree(
	ctx context.Context,
) (*permission.GetMenuTreeResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "获取菜单树",
		func(ctx context.Context) (interface{}, error) {
			// 调用RPC服务
			return s.identityClient.GetMenuTree(ctx)
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应为HTTP响应
	rpcResp := result.(*identity_srv.GetMenuTreeResponse)

	httpResp := s.assembler.Menu().ToHTTPGetMenuTreeResponse(rpcResp)
	if httpResp == nil {
		return nil, fmt.Errorf("转换菜单树响应失败")
	}

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

// =================================================================
// 菜单权限管理服务方法实现
// =================================================================

// ConfigureRoleMenus 配置角色的菜单权限
func (s *menuServiceImpl) ConfigureRoleMenus(
	ctx context.Context,
	operatorID string,
	req *permission.ConfigureRoleMenusRequestDTO,
) (*permission.ConfigureRoleMenusResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "配置角色菜单权限",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Menu().ToRPCConfigureRoleMenusRequest(operatorID, req)
			if rpcReq == nil {
				return nil, fmt.Errorf("转换配置角色菜单权限请求失败")
			}

			// 调用RPC服务
			return s.identityClient.ConfigureRoleMenus(ctx, rpcReq)
		},
		"role_id", *req.RoleID,
		"menu_configs_count", len(req.MenuConfigs),
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应为HTTP响应
	rpcResp := result.(*identity_srv.ConfigureRoleMenusResponse)

	httpResp := s.assembler.Menu().ToHTTPConfigureRoleMenusResponse(rpcResp)
	if httpResp == nil {
		return nil, fmt.Errorf("转换配置角色菜单权限响应失败")
	}

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

// GetRoleMenuTree 获取角色的菜单树
func (s *menuServiceImpl) GetRoleMenuTree(
	ctx context.Context,
	req *permission.GetRoleMenuTreeRequestDTO,
) (*permission.GetRoleMenuTreeResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "获取角色菜单树",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Menu().ToRPCGetRoleMenuTreeRequest(req)
			if rpcReq == nil {
				return nil, fmt.Errorf("转换获取角色菜单树请求失败")
			}

			// 调用RPC服务
			return s.identityClient.GetRoleMenuTree(ctx, rpcReq)
		},
		"role_id", *req.RoleID,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应为HTTP响应
	rpcResp := result.(*identity_srv.GetRoleMenuTreeResponse)

	httpResp := s.assembler.Menu().ToHTTPGetRoleMenuTreeResponse(rpcResp)
	if httpResp == nil {
		return nil, fmt.Errorf("转换角色菜单树响应失败")
	}

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

// GetUserMenuTree 获取用户的菜单树
func (s *menuServiceImpl) GetUserMenuTree(
	ctx context.Context,
	req *permission.GetUserMenuTreeRequestDTO,
) (*permission.GetUserMenuTreeResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "获取用户菜单树",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Menu().ToRPCGetUserMenuTreeRequest(req)
			if rpcReq == nil {
				return nil, fmt.Errorf("转换获取用户菜单树请求失败")
			}

			// 调用RPC服务
			return s.identityClient.GetUserMenuTree(ctx, rpcReq)
		},
		"user_id", *req.UserID,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应为HTTP响应
	rpcResp := result.(*identity_srv.GetUserMenuTreeResponse)

	httpResp := s.assembler.Menu().ToHTTPGetUserMenuTreeResponse(rpcResp)
	if httpResp == nil {
		return nil, fmt.Errorf("转换用户菜单树响应失败")
	}

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

// GetRoleMenuPermissions 获取角色的菜单权限列表
func (s *menuServiceImpl) GetRoleMenuPermissions(
	ctx context.Context,
	req *permission.GetRoleMenuPermissionsRequestDTO,
) (*permission.GetRoleMenuPermissionsResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "获取角色菜单权限列表",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Menu().ToRPCGetRoleMenuPermissionsRequest(req)
			if rpcReq == nil {
				return nil, fmt.Errorf("转换获取角色菜单权限列表请求失败")
			}

			// 调用RPC服务
			return s.identityClient.GetRoleMenuPermissions(ctx, rpcReq)
		},
		"role_id", *req.RoleID,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应为HTTP响应
	rpcResp := result.(*identity_srv.GetRoleMenuPermissionsResponse)

	httpResp := s.assembler.Menu().ToHTTPGetRoleMenuPermissionsResponse(rpcResp)
	if httpResp == nil {
		return nil, fmt.Errorf("转换角色菜单权限列表响应失败")
	}

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}

// HasMenuPermission 检查角色是否具有指定菜单权限
func (s *menuServiceImpl) HasMenuPermission(
	ctx context.Context,
	req *permission.HasMenuPermissionRequestDTO,
) (*permission.HasMenuPermissionResponseDTO, error) {
	// 使用BaseService模板处理RPC调用
	result, err := s.ProcessRPCCall(ctx, "检查角色菜单权限",
		func(ctx context.Context) (interface{}, error) {
			// 转换为RPC请求
			rpcReq := s.assembler.Menu().ToRPCHasMenuPermissionRequest(req)
			if rpcReq == nil {
				return nil, fmt.Errorf("转换检查角色菜单权限请求失败")
			}

			// 调用RPC服务
			return s.identityClient.HasMenuPermission(ctx, rpcReq)
		},
		"role_id", *req.RoleID,
		"menu_id", req.MenuID,
		"permission", req.Permission,
	)
	if err != nil {
		return nil, err
	}

	// 转换RPC响应为HTTP响应
	rpcResp := result.(*identity_srv.HasMenuPermissionResponse)

	httpResp := s.assembler.Menu().ToHTTPHasMenuPermissionResponse(rpcResp)
	if httpResp == nil {
		return nil, fmt.Errorf("转换检查菜单权限响应失败")
	}

	// 设置成功的基础响应
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}
