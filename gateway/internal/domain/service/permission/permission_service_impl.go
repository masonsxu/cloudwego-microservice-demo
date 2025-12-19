package permission

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
)

// permissionServiceImpl 权限管理聚合服务实现
// 实现所有子服务接口，提供统一的服务入口
type permissionServiceImpl struct {
	roleDefinitionService     RoleDefinitionService
	userRoleAssignmentService UserRoleAssignmentService
	// permissionSyncService     PermissionSyncService
	menuService MenuService
}

// NewService 创建权限管理聚合服务
func NewService(
	roleDefinitionService RoleDefinitionService,
	userRoleAssignmentService UserRoleAssignmentService,
	// permissionSyncService PermissionSyncService,
	menuService MenuService,
) Service {
	return &permissionServiceImpl{
		roleDefinitionService:     roleDefinitionService,
		userRoleAssignmentService: userRoleAssignmentService,
		// permissionSyncService:     permissionSyncService,
		menuService: menuService,
	}
}

// =================================================================
// RoleDefinitionService 接口实现 - 委托给 roleDefinitionService
// =================================================================

func (s *permissionServiceImpl) CreateRoleDefinition(
	ctx context.Context,
	req *permission.RoleDefinitionCreateRequestDTO,
) (*permission.RoleDefinitionCreateResponseDTO, error) {
	return s.roleDefinitionService.CreateRoleDefinition(ctx, req)
}

func (s *permissionServiceImpl) UpdateRoleDefinition(
	ctx context.Context,
	req *permission.RoleDefinitionUpdateRequestDTO,
) (*permission.RoleDefinitionUpdateResponseDTO, error) {
	return s.roleDefinitionService.UpdateRoleDefinition(ctx, req)
}

func (s *permissionServiceImpl) DeleteRoleDefinition(
	ctx context.Context,
	roleID string,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.roleDefinitionService.DeleteRoleDefinition(ctx, roleID)
}

func (s *permissionServiceImpl) GetRoleDefinition(
	ctx context.Context,
	req *permission.RoleDefinitionGetRequestDTO,
) (*permission.RoleDefinitionGetResponseDTO, error) {
	return s.roleDefinitionService.GetRoleDefinition(ctx, req)
}

func (s *permissionServiceImpl) ListRoleDefinitions(
	ctx context.Context,
	req *permission.RoleDefinitionQueryRequestDTO,
) (*permission.RoleDefinitionListResponseDTO, error) {
	return s.roleDefinitionService.ListRoleDefinitions(ctx, req)
}

// =================================================================
// UserRoleAssignmentService 接口实现 - 委托给 userRoleAssignmentService
// =================================================================

func (s *permissionServiceImpl) GetLastUserRoleAssignment(
	ctx context.Context,
	userID string,
) (*permission.AssignRoleToUserResponseDTO, error) {
	return s.userRoleAssignmentService.GetLastUserRoleAssignment(ctx, userID)
}

func (s *permissionServiceImpl) ListUserRoleAssignments(
	ctx context.Context,
	req *permission.UserRoleQueryRequestDTO,
) (*permission.UserRoleListResponseDTO, error) {
	return s.userRoleAssignmentService.ListUserRoleAssignments(ctx, req)
}

func (s *permissionServiceImpl) GetUsersByRole(
	ctx context.Context,
	req *permission.GetUsersByRoleRequestDTO,
) (*permission.GetUsersByRoleResponseDTO, error) {
	return s.userRoleAssignmentService.GetUsersByRole(ctx, req)
}

func (s *permissionServiceImpl) BatchBindUsersToRole(
	ctx context.Context,
	operatorID string,
	req *permission.BatchBindUsersToRoleRequestDTO,
) (*permission.BatchBindUsersToRoleResponseDTO, error) {
	return s.userRoleAssignmentService.BatchBindUsersToRole(ctx, operatorID, req)
}

// =================================================================
// PermissionSyncService 接口实现 - 委托给 permissionSyncService
// =================================================================

// func (s *permissionServiceImpl) SyncRoleToCasbin(
// 	ctx context.Context,
// 	roleID string,
// ) (*http_base.OperationStatusResponseDTO, error) {
// 	return s.permissionSyncService.SyncRoleToCasbin(ctx, roleID)
// }

// func (s *permissionServiceImpl) SyncAllRolesToCasbin(
// 	ctx context.Context,
// ) (*http_base.OperationStatusResponseDTO, error) {
// 	return s.permissionSyncService.SyncAllRolesToCasbin(ctx)
// }

// func (s *permissionServiceImpl) SyncAllUserRoles(
// 	ctx context.Context,
// ) (*http_base.OperationStatusResponseDTO, error) {
// 	return s.permissionSyncService.SyncAllUserRoles(ctx)
// }

// func (s *permissionServiceImpl) SyncUserRoles(
// 	ctx context.Context,
// 	userID string,
// ) (*http_base.OperationStatusResponseDTO, error) {
// 	return s.permissionSyncService.SyncUserRoles(ctx, userID)
// }

// =================================================================
// MenuService 接口实现 - 委托给 menuService
// =================================================================

func (s *permissionServiceImpl) UploadMenu(
	ctx context.Context,
	req *permission.UploadMenuRequestDTO,
) (*http_base.OperationStatusResponseDTO, error) {
	return s.menuService.UploadMenu(ctx, req)
}

func (s *permissionServiceImpl) GetMenuTree(
	ctx context.Context,
) (*permission.GetMenuTreeResponseDTO, error) {
	return s.menuService.GetMenuTree(ctx)
}

func (s *permissionServiceImpl) ConfigureRoleMenus(
	ctx context.Context,
	operatorID string,
	req *permission.ConfigureRoleMenusRequestDTO,
) (*permission.ConfigureRoleMenusResponseDTO, error) {
	return s.menuService.ConfigureRoleMenus(ctx, operatorID, req)
}

func (s *permissionServiceImpl) GetRoleMenuTree(
	ctx context.Context,
	req *permission.GetRoleMenuTreeRequestDTO,
) (*permission.GetRoleMenuTreeResponseDTO, error) {
	return s.menuService.GetRoleMenuTree(ctx, req)
}

func (s *permissionServiceImpl) GetUserMenuTree(
	ctx context.Context,
	req *permission.GetUserMenuTreeRequestDTO,
) (*permission.GetUserMenuTreeResponseDTO, error) {
	return s.menuService.GetUserMenuTree(ctx, req)
}

func (s *permissionServiceImpl) GetRoleMenuPermissions(
	ctx context.Context,
	req *permission.GetRoleMenuPermissionsRequestDTO,
) (*permission.GetRoleMenuPermissionsResponseDTO, error) {
	return s.menuService.GetRoleMenuPermissions(ctx, req)
}

func (s *permissionServiceImpl) HasMenuPermission(
	ctx context.Context,
	req *permission.HasMenuPermissionRequestDTO,
) (*permission.HasMenuPermissionResponseDTO, error) {
	return s.menuService.HasMenuPermission(ctx, req)
}
