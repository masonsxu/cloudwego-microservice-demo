package menu

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/casbin"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/parser"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// LogicImpl 菜单管理逻辑实现
type LogicImpl struct {
	dal                  dal.DAL
	converter            converter.Converter
	casbinManager        *casbin.CasbinManager
	userRoleAssignmentDA assignment.UserRoleAssignmentRepository
	config               *config.Config
}

// NewLogic 创建菜单管理逻辑实现
func NewLogic(
	dal dal.DAL,
	converter converter.Converter,
	casbinManager *casbin.CasbinManager,
	userRoleAssignmentDA assignment.UserRoleAssignmentRepository,
	config *config.Config,
) MenuLogic {
	return &LogicImpl{
		dal:                  dal,
		converter:            converter,
		casbinManager:        casbinManager,
		userRoleAssignmentDA: userRoleAssignmentDA,
		config:               config,
	}
}

// UploadMenu 上传并解析菜单配置文件 (menu.yaml)
func (l *LogicImpl) UploadMenu(
	ctx context.Context,
	req *identity_srv.UploadMenuRequest,
) error {
	if req.YamlContent == nil || *req.YamlContent == "" {
		return errno.ErrInvalidParams.WithMessage("YAML内容不能为空")
	}

	// 生成版本号，这里简单用当前时间戳
	version := fmt.Sprintf("v%d", getCurrentTimestamp())

	// 使用parser模块解析YAML内容并转换为模型
	menuModels, err := parser.ParseAndFlattenMenu(*req.YamlContent, version)
	if err != nil {
		return errno.ErrInvalidParams.WithMessage(fmt.Sprintf("解析菜单YAML失败: %s", err.Error()))
	}

	// 使用dal模块将菜单数据保存到数据库
	if err := l.dal.Menu().CreateMenuTree(ctx, menuModels); err != nil {
		return errno.ErrOperationFailed.WithMessage(fmt.Sprintf("保存菜单数据失败: %s", err.Error()))
	}

	return nil
}

// GetMenuTree 获取指定用户的菜单树
func (l *LogicImpl) GetMenuTree(
	ctx context.Context,
) (*identity_srv.GetMenuTreeResponse, error) {
	// 使用dal模块获取最新的菜单树
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("获取菜单树失败: %s", err.Error()))
	}

	// 使用converter模块将模型转换为Thrift对象
	menuNodes := l.converter.Menu().ModelsToThrift(menuModels)

	return &identity_srv.GetMenuTreeResponse{
		MenuTree: menuNodes,
	}, nil
}

// getCurrentTimestamp 获取当前时间戳（毫秒）
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// ConfigureRoleMenus 配置角色的菜单权限
func (l *LogicImpl) ConfigureRoleMenus(
	ctx context.Context,
	req *identity_srv.ConfigureRoleMenusRequest,
) (*identity_srv.ConfigureRoleMenusResponse, error) {
	// 1. 清除角色的旧菜单映射
	err := l.casbinManager.ClearRoleMenuMappings(*req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("清除角色旧菜单映射失败: %s", err.Error()),
		)
	}

	// 2. 添加新的菜单映射
	for _, config := range req.MenuConfigs {
		if config.MenuID == nil || *config.MenuID == "" {
			continue // 跳过空的菜单ID
		}

		// 验证并添加有效的权限映射
		if config.Permission != nil && *config.Permission != "" {
			// AddRoleMenuMapping 方法会验证权限类型是否有效
			err = l.casbinManager.AddRoleMenuMapping(
				*req.RoleID,
				*config.MenuID,
				*config.Permission,
			)
			if err != nil {
				return nil, errno.ErrOperationFailed.WithMessage(
					fmt.Sprintf("添加角色菜单映射失败: %s", err.Error()),
				)
			}
		}
	}

	successMsg := "菜单权限配置成功"

	return &identity_srv.ConfigureRoleMenusResponse{
		Success: convutil.BoolPtr(true),
		Message: &successMsg,
	}, nil
}

// GetRoleMenuTree 获取角色的菜单树
func (l *LogicImpl) GetRoleMenuTree(
	ctx context.Context,
	req *identity_srv.GetRoleMenuTreeRequest,
) (*identity_srv.GetRoleMenuTreeResponse, error) {
	if req.RoleID == nil || *req.RoleID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	// 使用新的方法构建带权限标记的完整菜单树
	menuNodes, err := l.buildMenuTreeWithPermissions(ctx, *req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("构建菜单树失败: %s", err.Error()))
	}

	return &identity_srv.GetRoleMenuTreeResponse{
		MenuTree: menuNodes,
		RoleID:   req.RoleID,
	}, nil
}

// GetUserMenuTree 获取用户的菜单树（基于所有活跃角色的权限合并）
func (l *LogicImpl) GetUserMenuTree(
	ctx context.Context,
	req *identity_srv.GetUserMenuTreeRequest,
) (*identity_srv.GetUserMenuTreeResponse, error) {
	if req.UserID == nil || *req.UserID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 1. 获取用户的所有活跃角色ID（只包含状态为 Active 的角色）
	roleIDs, err := l.userRoleAssignmentDA.GetActiveRoleIDsWithStatus(
		ctx,
		*req.UserID,
		models.RoleStatusActive,
	)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("获取用户角色列表失败: %s", err.Error()),
		)
	}

	if len(roleIDs) == 0 {
		// 用户没有活跃状态的角色绑定，返回空菜单树
		return &identity_srv.GetUserMenuTreeResponse{
			MenuTree: []*identity_srv.MenuNode{}, // 空菜单树
			UserID:   req.UserID,
			RoleIDs:  []string{},
		}, nil
	}

	// 2. 为每个角色获取菜单权限映射，并合并
	var (
		identityMaps []map[string]string
		isSuperAdmin bool
	)

	for _, roleID := range roleIDs {
		// 检查是否为超管角色
		isSuperAdminRole, err := l.isSuperAdminRole(ctx, roleID)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("检查超管角色失败: %s", err.Error()),
			)
		}

		if isSuperAdminRole {
			// 如果用户有超管角色，直接返回完整菜单树
			isSuperAdmin = true
			break
		}

		// 获取该角色的菜单权限映射
		policies, err := l.casbinManager.GetRoleMenuMappings(roleID)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取角色 %s 的菜单权限失败: %s", roleID, err.Error()),
			)
		}

		// 构建该角色的权限映射: menuID -> identity
		rolePermissionMap := make(map[string]string)

		for _, policy := range policies {
			if len(policy) >= 3 {
				menuID := policy[1]   // V1 是 menu_id
				identity := policy[2] // V2 是 identity
				rolePermissionMap[menuID] = identity
			}
		}

		identityMaps = append(identityMaps, rolePermissionMap)
	}

	// 3. 根据用户角色类型构建菜单树
	var menuNodes []*identity_srv.MenuNode

	if isSuperAdmin {
		// 超管用户：返回完整菜单树（无权限标记）
		menuNodes, err = l.getAllMenusWithoutPermissionMarks(ctx)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取完整菜单树失败: %s", err.Error()),
			)
		}
	} else {
		// 普通用户：合并所有角色的权限，构建授权菜单树
		// 3.1 合并权限映射（取最高权限）
		mergedPermissionMap := mergePermissionMaps(identityMaps)

		// 3.2 获取完整菜单树
		menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取完整菜单树失败: %s", err.Error()),
			)
		}

		// 3.3 转换为 Thrift 结构
		fullMenuTree := l.converter.Menu().ModelsToThrift(menuModels)

		// 3.4 过滤并返回授权菜单（只返回有权限的菜单）
		menuNodes = l.filterAuthorizedMenus(fullMenuTree, mergedPermissionMap)
	}

	return &identity_srv.GetUserMenuTreeResponse{
		MenuTree: menuNodes,
		UserID:   req.UserID,
		RoleIDs:  roleIDs, // 返回所有角色ID列表
	}, nil
}

// GetRoleMenuPermissions 获取角色的菜单权限列表
func (l *LogicImpl) GetRoleMenuPermissions(
	ctx context.Context,
	req *identity_srv.GetRoleMenuPermissionsRequest,
) (*identity_srv.GetRoleMenuPermissionsResponse, error) {
	if req.RoleID == nil || *req.RoleID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	// 1. 检查是否为超管角色
	isSuperAdmin, err := l.isSuperAdminRole(ctx, *req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("检查超管角色失败: %s", err.Error()))
	}

	if isSuperAdmin {
		// 超管角色：获取所有菜单并设置为 view_all_hospitals 权限
		menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取完整菜单树失败: %s", err.Error()),
			)
		}

		// 构建所有菜单的 view_all_hospitals 权限列表
		identitys := make([]*identity_srv.MenuPermission, 0)
		for _, menu := range menuModels {
			identitys = append(identitys, &identity_srv.MenuPermission{
				MenuID:     convutil.StringPtr(menu.SemanticID), // 使用语义ID
				Permission: convutil.StringPtr(models.MenuPermissionViewAllOrganizations),
			})
		}

		return &identity_srv.GetRoleMenuPermissionsResponse{
			Permissions: identitys,
			RoleID:      req.RoleID,
		}, nil
	}

	// 2. 普通角色：获取Casbin配置的权限
	// 获取角色的菜单映射
	policies, err := l.casbinManager.GetRoleMenuMappings(*req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("获取角色菜单权限失败: %s", err.Error()))
	}

	// 转换为 Thrift 结构
	identitys := make([]*identity_srv.MenuPermission, 0, len(policies))
	for _, policy := range policies {
		if len(policy) >= 3 {
			identitys = append(identitys, &identity_srv.MenuPermission{
				MenuID:     convutil.StringPtr(policy[1]), // V1 是 menu_id
				Permission: convutil.StringPtr(policy[2]), // V2 是 identity
			})
		}
	}

	return &identity_srv.GetRoleMenuPermissionsResponse{
		Permissions: identitys,
		RoleID:      req.RoleID,
	}, nil
}

// HasMenuPermission 检查角色是否具有指定菜单权限
func (l *LogicImpl) HasMenuPermission(
	ctx context.Context,
	req *identity_srv.HasMenuPermissionRequest,
) (*identity_srv.HasMenuPermissionResponse, error) {
	if req.RoleID == nil || *req.RoleID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	if req.MenuID == nil || *req.MenuID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("菜单ID不能为空")
	}

	if req.Permission == nil || *req.Permission == "" {
		return nil, errno.ErrInvalidParams.WithMessage("权限类型不能为空")
	}

	// 1. 检查是否为超管角色
	isSuperAdmin, err := l.isSuperAdminRole(ctx, *req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("检查超管角色失败: %s", err.Error()))
	}

	if isSuperAdmin {
		// 超管角色默认拥有所有菜单的所有权限
		return &identity_srv.HasMenuPermissionResponse{
			HasPermission: convutil.BoolPtr(true),
			RoleID:        req.RoleID,
			MenuID:        req.MenuID,
			Permission:    req.Permission,
		}, nil
	}

	// 2. 普通角色：检查具体权限配置
	// 获取角色的菜单映射
	policies, err := l.casbinManager.GetRoleMenuMappings(*req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("检查菜单权限失败: %s", err.Error()))
	}

	// 检查是否存在匹配的权限
	hasPermission := false

	for _, policy := range policies {
		if len(policy) >= 3 {
			if policy[1] == *req.MenuID && policy[2] == *req.Permission {
				hasPermission = true
				break
			}
		}
	}

	return &identity_srv.HasMenuPermissionResponse{
		HasPermission: &hasPermission,
		RoleID:        req.RoleID,
		MenuID:        req.MenuID,
		Permission:    req.Permission,
	}, nil
}

// GetUserMenuPermissions 获取用户的菜单权限列表（基于所有活跃角色合并）
func (l *LogicImpl) GetUserMenuPermissions(
	ctx context.Context,
	req *identity_srv.GetUserMenuPermissionsRequest,
) (*identity_srv.GetUserMenuPermissionsResponse, error) {
	// 1. 参数验证
	if req.UserID == nil || *req.UserID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 2. 获取用户的所有活跃角色ID
	roleIDs, err := l.userRoleAssignmentDA.GetActiveRoleIDsWithStatus(
		ctx,
		*req.UserID,
		models.RoleStatusActive,
	)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("获取用户角色列表失败: %s", err.Error()),
		)
	}

	// 3. 如果用户没有角色，返回空权限列表
	if len(roleIDs) == 0 {
		return &identity_srv.GetUserMenuPermissionsResponse{
			Permissions: []*identity_srv.MenuPermission{},
			UserID:      req.UserID,
			RoleIDs:     []string{},
		}, nil
	}

	// 4. 检查是否包含超管角色并收集权限映射
	var (
		permissionMaps []map[string]string
		isSuperAdmin   bool
	)

	for _, roleID := range roleIDs {
		// 检查是否为超管角色
		isSuperAdminRole, err := l.isSuperAdminRole(ctx, roleID)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("检查超管角色失败: %s", err.Error()),
			)
		}

		if isSuperAdminRole {
			isSuperAdmin = true
			break
		}

		// 获取该角色的菜单权限映射
		policies, err := l.casbinManager.GetRoleMenuMappings(roleID)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取角色 %s 的菜单权限失败: %s", roleID, err.Error()),
			)
		}

		// 构建该角色的权限映射
		rolePermissionMap := make(map[string]string)

		for _, policy := range policies {
			if len(policy) >= 3 {
				menuID := policy[1]     // V1 是 menu_id
				permission := policy[2] // V2 是 permission
				rolePermissionMap[menuID] = permission
			}
		}

		permissionMaps = append(permissionMaps, rolePermissionMap)
	}

	// 5. 构建权限列表
	var permissions []*identity_srv.MenuPermission

	if isSuperAdmin {
		// 超管用户：返回所有菜单的 view_all_hospitals 权限
		permissions, err = l.buildSuperAdminPermissions(ctx)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("构建超管权限列表失败: %s", err.Error()),
			)
		}
	} else {
		// 普通用户：合并所有角色的权限（取最高权限）
		mergedPermissionMap := mergePermissionMaps(permissionMaps)

		// 转换为 Thrift 结构
		permissions = make([]*identity_srv.MenuPermission, 0, len(mergedPermissionMap))
		for menuID, permission := range mergedPermissionMap {
			permissions = append(permissions, &identity_srv.MenuPermission{
				MenuID:     convutil.StringPtr(menuID),
				Permission: convutil.StringPtr(permission),
			})
		}
	}

	return &identity_srv.GetUserMenuPermissionsResponse{
		Permissions: permissions,
		UserID:      req.UserID,
		RoleIDs:     roleIDs,
	}, nil
}

// buildSuperAdminPermissions 构建超管用户的权限列表
// 返回所有菜单的 view_all_hospitals 权限
func (l *LogicImpl) buildSuperAdminPermissions(
	ctx context.Context,
) ([]*identity_srv.MenuPermission, error) {
	// 获取所有菜单
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("获取完整菜单树失败: %s", err.Error()),
		)
	}

	// 构建所有菜单的 view_all_hospitals 权限列表
	permissions := make([]*identity_srv.MenuPermission, 0, len(menuModels))
	for _, menu := range menuModels {
		permissions = append(permissions, &identity_srv.MenuPermission{
			MenuID:     convutil.StringPtr(menu.SemanticID),
			Permission: convutil.StringPtr(models.MenuPermissionViewAllOrganizations),
		})
	}

	return permissions, nil
}

// markMenuPermissions 递归标记菜单树中每个节点的权限状态，同时返回带权限映射的完整菜单树
func (l *LogicImpl) markMenuPermissions(
	menuNodes []*identity_srv.MenuNode,
	identityMap map[string]string, // menuID -> identity level
) []*identity_srv.MenuNode {
	var result []*identity_srv.MenuNode

	for _, node := range menuNodes {
		// 复制节点
		newNode := *node

		// 检查权限映射
		if identityLevel, exists := identityMap[*node.Id]; exists {
			// 有权限
			newNode.HasPermission = &[]bool{true}[0] // 创建bool指针
			newNode.PermissionLevel = &identityLevel
		} else {
			// 无权限
			newNode.HasPermission = &[]bool{false}[0] // 创建bool指针
			newNode.PermissionLevel = &[]string{"none"}[0]
		}

		// 递归处理子菜单
		if len(node.Children) > 0 {
			newNode.Children = l.markMenuPermissions(node.Children, identityMap)
		}

		result = append(result, &newNode)
	}

	return result
}

// filterAuthorizedMenus 递归过滤菜单树，只保留有权限的菜单节点
// 过滤规则：
// 1. 叶子节点：有权限则保留，无权限则移除
// 2. 父节点：如果有直接权限或有授权子节点则保留
// 3. 空父节点：所有子节点都无权限的父节点被移除
func (l *LogicImpl) filterAuthorizedMenus(
	menuNodes []*identity_srv.MenuNode,
	identityMap map[string]string, // menuID -> identity level
) []*identity_srv.MenuNode {
	var authorizedMenus []*identity_srv.MenuNode

	for _, node := range menuNodes {
		// 检查当前节点是否有权限
		_, hasDirectPermission := identityMap[*node.Id]

		// 递归处理子菜单，获取有权限的子菜单
		var authorizedChildren []*identity_srv.MenuNode
		if len(node.Children) > 0 {
			authorizedChildren = l.filterAuthorizedMenus(node.Children, identityMap)
		}

		// 决定是否保留当前节点：
		// 1. 有直接权限，或者
		// 2. 有授权的子菜单
		if hasDirectPermission || len(authorizedChildren) > 0 {
			// 创建新的菜单节点（不包含权限标记字段）
			authorizedNode := &identity_srv.MenuNode{
				Name:      node.Name,
				Id:        node.Id,
				Path:      node.Path,
				Icon:      node.Icon,
				Component: node.Component,
				Children:  authorizedChildren,
				// 注意：不设置 HasPermission 和 PermissionLevel 字段
			}

			authorizedMenus = append(authorizedMenus, authorizedNode)
		}
		// 如果既没有直接权限，也没有授权子菜单，则丢弃该节点
	}

	return authorizedMenus
}

// buildMenuTreeWithPermissions 构建带权限标记的完整菜单树（用于权限管理）
func (l *LogicImpl) buildMenuTreeWithPermissions(
	ctx context.Context,
	roleID string,
) ([]*identity_srv.MenuNode, error) {
	// 1. 检查是否为超管角色
	isSuperAdmin, err := l.isSuperAdminRole(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("检查超管角色失败: " + err.Error())
	}

	if isSuperAdmin {
		// 超管角色：获取所有菜单并标记为full权限
		return l.getAllMenusWithFullPermissions(ctx)
	}

	// 2. 普通角色：按现有逻辑处理
	// 获取完整菜单树
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取完整菜单树失败: " + err.Error())
	}

	// 3. 转换为 Thrift 结构
	fullMenuTree := l.converter.Menu().ModelsToThrift(menuModels)

	// 4. 获取角色的菜单权限映射
	policies, err := l.casbinManager.GetRoleMenuMappings(roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取角色菜单权限映射失败: " + err.Error())
	}

	// 5. 构建权限映射表: menuID -> identityLevel
	identityMap := make(map[string]string)

	for _, policy := range policies {
		if len(policy) >= 3 {
			menuID := policy[1]   // V1 是 menu_id
			identity := policy[2] // V2 是 identity
			identityMap[menuID] = identity
		}
	}

	// 6. 标记权限并返回完整菜单树
	return l.markMenuPermissions(fullMenuTree, identityMap), nil
}

// isSuperAdminRole 检查角色是否为超管角色
// 通过查询角色定义表，检查角色UUID对应的角色名称是否在配置的超管角色列表中
func (l *LogicImpl) isSuperAdminRole(ctx context.Context, roleID string) (bool, error) {
	// 1. 检查配置中是否定义了超管角色
	if len(l.config.SuperAdmin.RoleNames) == 0 {
		return false, nil
	}

	// 2. 通过 UUID 查询角色定义，获取角色名称
	roleDefinition, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 角色不存在，非超管
			return false, nil
		}

		return false, errno.ErrOperationFailed.WithMessage("查询角色定义失败: " + err.Error())
	}

	if roleDefinition == nil {
		// 角色不存在，非超管
		return false, nil
	}

	// 3. 创建超管角色名称的映射，用于快速查找
	superAdminRoles := make(map[string]bool)
	for _, roleName := range l.config.SuperAdmin.RoleNames {
		superAdminRoles[roleName] = true
	}

	// 4. 检查角色名称是否在超管角色列表中
	return superAdminRoles[roleDefinition.Name], nil
}

// getAllMenusWithFullPermissions 获取所有菜单并标记为完整权限
// 超管角色使用此方法获取完整菜单树，所有菜单都标记为 view_all_hospitals 权限
func (l *LogicImpl) getAllMenusWithFullPermissions(
	ctx context.Context,
) ([]*identity_srv.MenuNode, error) {
	// 1. 获取完整菜单树
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取完整菜单树失败: " + err.Error())
	}

	// 2. 转换为 Thrift 结构
	fullMenuTree := l.converter.Menu().ModelsToThrift(menuModels)

	// 3. 为所有菜单标记 view_all_hospitals 权限
	return l.markAllMenusWithFullPermission(fullMenuTree), nil
}

// getAllMenusWithoutPermissionMarks 获取所有菜单但不添加权限标记
// 超管用户菜单使用此方法，返回纯净的菜单结构
func (l *LogicImpl) getAllMenusWithoutPermissionMarks(
	ctx context.Context,
) ([]*identity_srv.MenuNode, error) {
	// 1. 获取完整菜单树
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取完整菜单树失败: " + err.Error())
	}

	// 2. 转换为 Thrift 结构（不添加权限标记）
	return l.converter.Menu().ModelsToThrift(menuModels), nil
}

// markAllMenusWithFullPermission 递归标记所有菜单为完整权限
// 超管专用方法，为菜单树中的所有节点标记 view_all_hospitals 权限
func (l *LogicImpl) markAllMenusWithFullPermission(
	menuNodes []*identity_srv.MenuNode,
) []*identity_srv.MenuNode {
	var result []*identity_srv.MenuNode

	for _, node := range menuNodes {
		// 复制节点
		newNode := *node

		// 标记为 view_all_hospitals 权限
		newNode.HasPermission = &[]bool{true}[0]
		newNode.PermissionLevel = &[]string{models.MenuPermissionViewAllOrganizations}[0]

		// 递归处理子菜单
		if len(node.Children) > 0 {
			newNode.Children = l.markAllMenusWithFullPermission(node.Children)
		}

		result = append(result, &newNode)
	}

	return result
}

// getHigherPermission 比较两个权限，返回更高权限等级
// 权限等级: view_all_hospitals > view_own_hospital > none
func getHigherPermission(perm1, perm2 string) string {
	// 如果任一权限为 view_all_hospitals，返回 view_all_hospitals
	if perm1 == models.MenuPermissionViewAllOrganizations ||
		perm2 == models.MenuPermissionViewAllOrganizations {
		return models.MenuPermissionViewAllOrganizations
	}

	// 如果任一权限为 view_own_hospital，返回 view_own_hospital
	if perm1 == models.MenuPermissionViewOwnOrganization ||
		perm2 == models.MenuPermissionViewOwnOrganization {
		return models.MenuPermissionViewOwnOrganization
	}

	// 都无权限，返回空字符串
	return ""
}

// mergePermissionMaps 合并多个权限映射，对每个菜单取最高权限等级
// 用于多角色权限合并场景
func mergePermissionMaps(maps []map[string]string) map[string]string {
	merged := make(map[string]string)

	for _, m := range maps {
		for menuID, identity := range m {
			// 如果菜单已存在，比较并取更高权限
			if existingPermission, exists := merged[menuID]; exists {
				merged[menuID] = getHigherPermission(existingPermission, identity)
			} else {
				// 菜单不存在，直接添加
				merged[menuID] = identity
			}
		}
	}

	return merged
}
