package menu

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/menu"
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
	userRoleAssignmentDA assignment.UserRoleAssignmentRepository
	config               *config.Config
}

// NewLogic 创建菜单管理逻辑实现
func NewLogic(
	dal dal.DAL,
	converter converter.Converter,
	userRoleAssignmentDA assignment.UserRoleAssignmentRepository,
	config *config.Config,
) MenuLogic {
	return &LogicImpl{
		dal:                  dal,
		converter:            converter,
		userRoleAssignmentDA: userRoleAssignmentDA,
		config:               config,
	}
}

// =============================================================================
// 权限级别枚举转换辅助函数
// =============================================================================

// toThriftPermissionLevel 将内部权限类型转换为 Thrift 枚举
func toThriftPermissionLevel(p models.MenuPermissionType) identity_srv.PermissionLevel {
	switch p {
	case models.PermissionNone:
		return identity_srv.PermissionLevel_NONE
	case models.PermissionView:
		return identity_srv.PermissionLevel_READ
	case models.PermissionEdit, models.PermissionManage:
		return identity_srv.PermissionLevel_WRITE
	case models.PermissionFull:
		return identity_srv.PermissionLevel_FULL
	default:
		return identity_srv.PermissionLevel_NONE
	}
}

// toThriftPermissionLevelPtr 将内部权限类型转换为 Thrift 枚举指针
func toThriftPermissionLevelPtr(p models.MenuPermissionType) *identity_srv.PermissionLevel {
	level := toThriftPermissionLevel(p)
	return &level
}

// parsePermissionTypeFromThrift 从 Thrift 枚举解析权限类型
func parsePermissionTypeFromThrift(level identity_srv.PermissionLevel) models.MenuPermissionType {
	switch level {
	case identity_srv.PermissionLevel_NONE:
		return models.PermissionNone
	case identity_srv.PermissionLevel_READ:
		return models.PermissionView
	case identity_srv.PermissionLevel_WRITE:
		return models.PermissionEdit
	case identity_srv.PermissionLevel_FULL:
		return models.PermissionFull
	default:
		return models.PermissionNone
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

	productLine := menu.DefaultProductLine
	if req.ProductLine != nil && *req.ProductLine != "" {
		productLine = *req.ProductLine
	}

	// 使用parser模块解析YAML内容并转换为模型
	// 先用版本号0解析，获取哈希值
	menuModels, contentHash, err := parser.ParseAndFlattenMenu(*req.YamlContent, productLine, 0)
	if err != nil {
		return errno.ErrInvalidParams.WithMessage(fmt.Sprintf("解析菜单YAML失败: %s", err.Error()))
	}

	// 去重检查：查询该产品线最新版本的哈希
	latestHash, err := l.dal.Menu().GetLatestContentHash(ctx, productLine)
	if err == nil && latestHash == contentHash {
		// 内容相同，静默返回成功（不创建新版本）
		return nil
	}

	// 获取新版本号（自增）
	maxVersion, err := l.dal.Menu().GetMaxVersion(ctx, productLine)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage(fmt.Sprintf("获取版本号失败: %s", err.Error()))
	}

	newVersion := maxVersion + 1

	// 更新所有菜单的版本号
	for _, menu := range menuModels {
		menu.Version = newVersion
	}

	// 使用dal模块将菜单数据保存到数据库
	if err := l.dal.Menu().CreateMenuTree(ctx, menuModels); err != nil {
		return errno.ErrOperationFailed.WithMessage(fmt.Sprintf("保存菜单数据失败: %s", err.Error()))
	}

	return nil
}

// GetMenuTree 获取完整菜单树
func (l *LogicImpl) GetMenuTree(
	ctx context.Context,
) (*identity_srv.GetMenuTreeResponse, error) {
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx, menu.DefaultProductLine)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("获取菜单树失败: %s", err.Error()))
	}

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
	if req.RoleID == nil || *req.RoleID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	roleID, err := uuid.Parse(*req.RoleID)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的角色ID格式")
	}

	// 构建权限列表
	permissions := make([]*models.RoleMenuPermission, 0, len(req.MenuConfigs))
	for _, cfg := range req.MenuConfigs {
		if cfg.MenuID == nil || *cfg.MenuID == "" {
			continue
		}

		permission := &models.RoleMenuPermission{
			RoleID:         roleID,
			MenuID:         *cfg.MenuID,
			PermissionType: models.PermissionView,  // 默认查看权限
			DataScope:      models.DataScopeOwnOrg, // 默认所在组织
		}

		// 解析权限类型（从 Thrift 枚举）
		if cfg.Permission != nil {
			permission.PermissionType = parsePermissionTypeFromThrift(*cfg.Permission)
			// DataScope 根据权限类型设置
			if *cfg.Permission == identity_srv.PermissionLevel_FULL {
				permission.DataScope = models.DataScopeAllOrgs
			}
		}

		permissions = append(permissions, permission)
	}

	// 同步角色菜单权限
	if err := l.dal.RoleMenuPermission().SyncRoleMenus(ctx, roleID, permissions); err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("配置角色菜单权限失败: %s", err.Error()))
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

	roleID, err := uuid.Parse(*req.RoleID)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的角色ID格式")
	}

	// 构建带权限标记的完整菜单树
	menuNodes, err := l.buildMenuTreeWithPermissions(ctx, roleID)
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

	// 1. 获取用户的所有活跃角色ID
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
		return &identity_srv.GetUserMenuTreeResponse{
			MenuTree: []*identity_srv.MenuNode{},
			UserID:   req.UserID,
			RoleIDs:  []string{},
		}, nil
	}

	// 2. 检查是否有超管角色
	isSuperAdmin, err := l.checkSuperAdminRoles(ctx, roleIDs)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("检查超管角色失败: %s", err.Error()),
		)
	}

	var menuNodes []*identity_srv.MenuNode

	if isSuperAdmin {
		// 超管用户：返回完整菜单树
		menuNodes, err = l.getAllMenusWithoutPermissionMarks(ctx)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取完整菜单树失败: %s", err.Error()),
			)
		}
	} else {
		// 普通用户：获取合并权限并过滤菜单
		roleUUIDs := make([]uuid.UUID, len(roleIDs))
		for i, id := range roleIDs {
			roleUUIDs[i], _ = uuid.Parse(id)
		}

		// 获取合并后的权限
		mergedPermissions, err := l.dal.RoleMenuPermission().GetMergedPermissions(ctx, roleUUIDs)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取用户菜单权限失败: %s", err.Error()),
			)
		}

		// 构建权限映射
		permMap := make(map[string]models.MenuPermissionInfo)
		for _, perm := range mergedPermissions {
			permMap[perm.MenuID] = perm
		}

		// 获取完整菜单树
		menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx, menu.DefaultProductLine)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取完整菜单树失败: %s", err.Error()),
			)
		}

		fullMenuTree := l.converter.Menu().ModelsToThrift(menuModels)

		// 过滤已授权菜单
		menuNodes = l.filterAuthorizedMenus(fullMenuTree, permMap)
	}

	return &identity_srv.GetUserMenuTreeResponse{
		MenuTree: menuNodes,
		UserID:   req.UserID,
		RoleIDs:  roleIDs,
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

	roleID, err := uuid.Parse(*req.RoleID)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的角色ID格式")
	}

	// 检查是否为超管角色
	isSuperAdmin, err := l.isSuperAdminRole(ctx, *req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("检查超管角色失败: %s", err.Error()))
	}

	if isSuperAdmin {
		// 超管角色：返回所有菜单的完全控制权限
		return l.buildSuperAdminPermissionsResponse(ctx, req.RoleID)
	}

	// 普通角色：获取数据库配置的权限
	permissions, err := l.dal.RoleMenuPermission().GetByRoleID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("获取角色菜单权限失败: %s", err.Error()))
	}

	// 转换为 Thrift 结构
	thriftPermissions := make([]*identity_srv.MenuPermission, 0, len(permissions))
	for _, perm := range permissions {
		thriftPermissions = append(thriftPermissions, &identity_srv.MenuPermission{
			MenuID:     convutil.StringPtr(perm.MenuID),
			Permission: toThriftPermissionLevelPtr(perm.PermissionType),
		})
	}

	return &identity_srv.GetRoleMenuPermissionsResponse{
		Permissions: thriftPermissions,
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

	// 检查是否为超管角色
	isSuperAdmin, err := l.isSuperAdminRole(ctx, *req.RoleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("检查超管角色失败: %s", err.Error()))
	}

	if isSuperAdmin {
		return &identity_srv.HasMenuPermissionResponse{
			HasPermission: convutil.BoolPtr(true),
			RoleID:        req.RoleID,
			MenuID:        req.MenuID,
			Permission:    req.Permission,
		}, nil
	}

	roleID, err := uuid.Parse(*req.RoleID)
	if err != nil {
		return nil, errno.ErrInvalidParams.WithMessage("无效的角色ID格式")
	}

	// 解析请求的权限类型
	requiredPermType := models.PermissionView
	if req.Permission != nil {
		requiredPermType = parsePermissionTypeFromThrift(*req.Permission)
	}

	// 检查权限
	hasPermission, err := l.dal.RoleMenuPermission().HasPermission(ctx, roleID, *req.MenuID, requiredPermType)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(fmt.Sprintf("检查菜单权限失败: %s", err.Error()))
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
	if req.UserID == nil || *req.UserID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 获取用户的所有活跃角色ID
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
		return &identity_srv.GetUserMenuPermissionsResponse{
			Permissions: []*identity_srv.MenuPermission{},
			UserID:      req.UserID,
			RoleIDs:     []string{},
		}, nil
	}

	// 检查是否有超管角色
	isSuperAdmin, err := l.checkSuperAdminRoles(ctx, roleIDs)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("检查超管角色失败: %s", err.Error()),
		)
	}

	var thriftPermissions []*identity_srv.MenuPermission

	if isSuperAdmin {
		// 超管用户：返回所有菜单的完全控制权限
		thriftPermissions, err = l.buildSuperAdminPermissions(ctx)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("构建超管权限列表失败: %s", err.Error()),
			)
		}
	} else {
		// 普通用户：合并所有角色的权限
		roleUUIDs := make([]uuid.UUID, len(roleIDs))
		for i, id := range roleIDs {
			roleUUIDs[i], _ = uuid.Parse(id)
		}

		mergedPermissions, err := l.dal.RoleMenuPermission().GetMergedPermissions(ctx, roleUUIDs)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage(
				fmt.Sprintf("获取用户菜单权限失败: %s", err.Error()),
			)
		}

		thriftPermissions = make([]*identity_srv.MenuPermission, 0, len(mergedPermissions))
		for _, perm := range mergedPermissions {
			thriftPermissions = append(thriftPermissions, &identity_srv.MenuPermission{
				MenuID:     convutil.StringPtr(perm.MenuID),
				Permission: toThriftPermissionLevelPtr(perm.PermissionType),
			})
		}
	}

	return &identity_srv.GetUserMenuPermissionsResponse{
		Permissions: thriftPermissions,
		UserID:      req.UserID,
		RoleIDs:     roleIDs,
	}, nil
}

// =============================================================================
// 内部辅助方法
// =============================================================================

// isSuperAdminRole 检查角色是否为超管角色
func (l *LogicImpl) isSuperAdminRole(ctx context.Context, roleID string) (bool, error) {
	if len(l.config.SuperAdmin.RoleNames) == 0 {
		return false, nil
	}

	roleDefinition, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if roleDefinition == nil {
		return false, nil
	}

	superAdminRoles := make(map[string]bool)
	for _, roleName := range l.config.SuperAdmin.RoleNames {
		superAdminRoles[roleName] = true
	}

	return superAdminRoles[roleDefinition.Name], nil
}

// checkSuperAdminRoles 检查角色列表中是否有超管角色
func (l *LogicImpl) checkSuperAdminRoles(ctx context.Context, roleIDs []string) (bool, error) {
	for _, roleID := range roleIDs {
		isSuperAdmin, err := l.isSuperAdminRole(ctx, roleID)
		if err != nil {
			return false, err
		}
		if isSuperAdmin {
			return true, nil
		}
	}
	return false, nil
}

// buildMenuTreeWithPermissions 构建带权限标记的完整菜单树
func (l *LogicImpl) buildMenuTreeWithPermissions(
	ctx context.Context,
	roleID uuid.UUID,
) ([]*identity_srv.MenuNode, error) {
	// 检查是否为超管角色
	isSuperAdmin, err := l.isSuperAdminRole(ctx, roleID.String())
	if err != nil {
		return nil, err
	}

	if isSuperAdmin {
		return l.getAllMenusWithFullPermissions(ctx)
	}

	// 获取完整菜单树
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx, menu.DefaultProductLine)
	if err != nil {
		return nil, err
	}

	fullMenuTree := l.converter.Menu().ModelsToThrift(menuModels)

	// 获取角色的菜单权限
	permissions, err := l.dal.RoleMenuPermission().GetByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 构建权限映射
	permMap := make(map[string]models.MenuPermissionInfo)
	for _, perm := range permissions {
		permMap[perm.MenuID] = models.MenuPermissionInfo{
			MenuID:         perm.MenuID,
			PermissionType: perm.PermissionType,
			DataScope:      perm.DataScope,
		}
	}

	// 标记权限
	return l.markMenuPermissions(fullMenuTree, permMap), nil
}

// markMenuPermissions 递归标记菜单树中每个节点的权限状态
func (l *LogicImpl) markMenuPermissions(
	menuNodes []*identity_srv.MenuNode,
	permMap map[string]models.MenuPermissionInfo,
) []*identity_srv.MenuNode {
	var result []*identity_srv.MenuNode

	for _, node := range menuNodes {
		newNode := *node

		if perm, exists := permMap[*node.Id]; exists {
			newNode.HasPermission = convutil.BoolPtr(true)
			newNode.PermissionLevel = toThriftPermissionLevelPtr(perm.PermissionType)
		} else {
			newNode.HasPermission = convutil.BoolPtr(false)
			noneLevel := identity_srv.PermissionLevel_NONE
			newNode.PermissionLevel = &noneLevel
		}

		if len(node.Children) > 0 {
			newNode.Children = l.markMenuPermissions(node.Children, permMap)
		}

		result = append(result, &newNode)
	}

	return result
}

// filterAuthorizedMenus 递归过滤菜单树，只保留有权限的菜单节点
func (l *LogicImpl) filterAuthorizedMenus(
	menuNodes []*identity_srv.MenuNode,
	permMap map[string]models.MenuPermissionInfo,
) []*identity_srv.MenuNode {
	var authorizedMenus []*identity_srv.MenuNode

	for _, node := range menuNodes {
		_, hasDirectPermission := permMap[*node.Id]

		var authorizedChildren []*identity_srv.MenuNode
		if len(node.Children) > 0 {
			authorizedChildren = l.filterAuthorizedMenus(node.Children, permMap)
		}

		if hasDirectPermission || len(authorizedChildren) > 0 {
			authorizedNode := &identity_srv.MenuNode{
				Name:      node.Name,
				Id:        node.Id,
				Path:      node.Path,
				Icon:      node.Icon,
				Component: node.Component,
				Children:  authorizedChildren,
			}
			authorizedMenus = append(authorizedMenus, authorizedNode)
		}
	}

	return authorizedMenus
}

// getAllMenusWithFullPermissions 获取所有菜单并标记为完整权限
func (l *LogicImpl) getAllMenusWithFullPermissions(
	ctx context.Context,
) ([]*identity_srv.MenuNode, error) {
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx, menu.DefaultProductLine)
	if err != nil {
		return nil, err
	}

	fullMenuTree := l.converter.Menu().ModelsToThrift(menuModels)
	return l.markAllMenusWithFullPermission(fullMenuTree), nil
}

// getAllMenusWithoutPermissionMarks 获取所有菜单但不添加权限标记
func (l *LogicImpl) getAllMenusWithoutPermissionMarks(
	ctx context.Context,
) ([]*identity_srv.MenuNode, error) {
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx, menu.DefaultProductLine)
	if err != nil {
		return nil, err
	}
	return l.converter.Menu().ModelsToThrift(menuModels), nil
}

// markAllMenusWithFullPermission 递归标记所有菜单为完整权限
func (l *LogicImpl) markAllMenusWithFullPermission(
	menuNodes []*identity_srv.MenuNode,
) []*identity_srv.MenuNode {
	var result []*identity_srv.MenuNode

	for _, node := range menuNodes {
		newNode := *node
		newNode.HasPermission = convutil.BoolPtr(true)
		newNode.PermissionLevel = toThriftPermissionLevelPtr(models.PermissionFull)

		if len(node.Children) > 0 {
			newNode.Children = l.markAllMenusWithFullPermission(node.Children)
		}

		result = append(result, &newNode)
	}

	return result
}

// buildSuperAdminPermissions 构建超管用户的权限列表
func (l *LogicImpl) buildSuperAdminPermissions(
	ctx context.Context,
) ([]*identity_srv.MenuPermission, error) {
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx, menu.DefaultProductLine)
	if err != nil {
		return nil, err
	}

	permissions := make([]*identity_srv.MenuPermission, 0, len(menuModels))
	l.collectMenuPermissions(menuModels, &permissions)

	return permissions, nil
}

// collectMenuPermissions 递归收集菜单权限
func (l *LogicImpl) collectMenuPermissions(
	menus []*models.Menu,
	permissions *[]*identity_srv.MenuPermission,
) {
	for _, menu := range menus {
		*permissions = append(*permissions, &identity_srv.MenuPermission{
			MenuID:     convutil.StringPtr(menu.SemanticID),
			Permission: toThriftPermissionLevelPtr(models.PermissionFull),
		})
		if len(menu.Children) > 0 {
			l.collectMenuPermissions(menu.Children, permissions)
		}
	}
}

// buildSuperAdminPermissionsResponse 构建超管角色权限响应
func (l *LogicImpl) buildSuperAdminPermissionsResponse(
	ctx context.Context,
	roleID *string,
) (*identity_srv.GetRoleMenuPermissionsResponse, error) {
	permissions, err := l.buildSuperAdminPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return &identity_srv.GetRoleMenuPermissionsResponse{
		Permissions: permissions,
		RoleID:      roleID,
	}, nil
}
