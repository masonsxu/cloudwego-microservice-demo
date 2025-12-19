package casbin

import (
	"context"
	"fmt"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// MenuPermissionLogic 菜单权限管理逻辑接口
type MenuPermissionLogic interface {
	// ConfigureRoleMenus 配置角色菜单权限
	ConfigureRoleMenus(roleID string, menuConfigs []*MenuConfig) error

	// GetRoleMenus 获取角色的菜单列表
	GetRoleMenus(roleID string) ([]*identity_srv.MenuNode, error)

	// GetUserMenus 获取用户的菜单列表（基于角色合并）
	GetUserMenus(userID string) ([]*identity_srv.MenuNode, error)

	// HasMenuPermission 检查角色是否具有指定菜单权限
	HasMenuPermission(roleID, menuID, permission string) (bool, error)

	// GetRoleMenuPermissions 获取角色的所有菜单权限
	GetRoleMenuPermissions(roleID string) ([]*MenuPermission, error)

	// SyncRoleMenus 同步角色菜单（替换现有配置）
	SyncRoleMenus(roleID string, menuConfigs []*MenuConfig) error
}

// MenuPermissionLogicImpl 菜单权限管理逻辑实现
type MenuPermissionLogicImpl struct {
	casbinManager        *CasbinManager
	dal                  dal.DAL
	converter            converter.Converter
	userRoleAssignmentDA assignment.UserRoleAssignmentRepository
}

// MenuConfig 菜单配置
type MenuConfig struct {
	MenuID     string `json:"menu_id"`    // 语义化菜单ID (来自menu.yaml)
	Permission string `json:"permission"` // view_own_hospital, view_all_hospitals
}

// MenuPermission 菜单权限信息
type MenuPermission struct {
	MenuID     string `json:"menu_id"`    // 语义化菜单ID (来自menu.yaml)
	Permission string `json:"permission"` // view_own_hospital, view_all_hospitals
}

// NewMenuPermissionLogic 创建菜单权限管理逻辑实例
func NewMenuPermissionLogic(
	casbinManager *CasbinManager,
	dal dal.DAL,
	converter converter.Converter,
	userRoleAssignmentDA assignment.UserRoleAssignmentRepository,
) MenuPermissionLogic {
	return &MenuPermissionLogicImpl{
		casbinManager:        casbinManager,
		dal:                  dal,
		converter:            converter,
		userRoleAssignmentDA: userRoleAssignmentDA,
	}
}

// ConfigureRoleMenus 配置角色菜单权限
func (l *MenuPermissionLogicImpl) ConfigureRoleMenus(roleID string, menuConfigs []*MenuConfig) error {
	// 1. 验证角色是否存在
	if roleID == "" {
		return fmt.Errorf("角色ID不能为空")
	}

	// 2. 清除角色的旧菜单映射
	err := l.casbinManager.ClearRoleMenuMappings(roleID)
	if err != nil {
		return fmt.Errorf("清除角色旧菜单映射失败: %w", err)
	}

	// 3. 添加新的菜单映射
	for _, config := range menuConfigs {
		if config.MenuID == "" {
			continue // 跳过空的菜单ID
		}

		// 验证并添加有效的权限映射
		if config.Permission != "" {
			// AddRoleMenuMapping 方法会验证权限类型是否有效
			err = l.casbinManager.AddRoleMenuMapping(roleID, config.MenuID, config.Permission)
			if err != nil {
				return fmt.Errorf("添加角色菜单映射失败: %w", err)
			}
		}
	}

	return nil
}

// SyncRoleMenus 同步角色菜单（替换现有配置）
func (l *MenuPermissionLogicImpl) SyncRoleMenus(roleID string, menuConfigs []*MenuConfig) error {
	return l.ConfigureRoleMenus(roleID, menuConfigs)
}

// GetRoleMenus 获取角色的菜单列表
func (l *MenuPermissionLogicImpl) GetRoleMenus(roleID string) ([]*identity_srv.MenuNode, error) {
	// 1. 获取角色的菜单映射
	policies, err := l.casbinManager.GetRoleMenuMappings(roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色菜单映射失败: %w", err)
	}

	// 2. 提取菜单ID列表
	menuIDSet := make(map[string]bool)
	for _, policy := range policies {
		if len(policy) >= 2 {
			menuID := policy[1] // V1 是 menu_id
			menuIDSet[menuID] = true
		}
	}

	// 3. 转换为切片
	menuIDs := make([]string, 0, len(menuIDSet))
	for menuID := range menuIDSet {
		menuIDs = append(menuIDs, menuID)
	}

	// 4. 构建菜单树
	return l.buildMenuTree(menuIDs)
}

// GetUserMenus 获取用户的菜单列表（基于最新角色绑定）
func (l *MenuPermissionLogicImpl) GetUserMenus(userID string) ([]*identity_srv.MenuNode, error) {
	// 1. 获取用户的最新角色绑定
	ctx := context.TODO() // 临时使用，应该从上层传递
	lastAssignment, err := l.userRoleAssignmentDA.GetLastUserRoleAssignment(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户最新角色绑定失败: %w", err)
	}

	if lastAssignment == nil {
		// 用户没有角色绑定，返回空菜单
		return []*identity_srv.MenuNode{}, nil
	}

	// 2. 获取该角色的菜单映射
	roleID := lastAssignment.RoleID.String()
	policies, err := l.casbinManager.GetRoleMenuMappings(roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色菜单映射失败: %w", err)
	}

	// 3. 提取菜单ID列表
	menuIDSet := make(map[string]bool)
	for _, policy := range policies {
		if len(policy) >= 2 {
			menuID := policy[1] // V1 是 menu_id
			menuIDSet[menuID] = true
		}
	}

	// 4. 转换为切片
	menuIDs := make([]string, 0, len(menuIDSet))
	for menuID := range menuIDSet {
		menuIDs = append(menuIDs, menuID)
	}

	// 5. 构建菜单树
	return l.buildMenuTree(menuIDs)
}

// HasMenuPermission 检查角色是否具有指定菜单权限
func (l *MenuPermissionLogicImpl) HasMenuPermission(roleID, menuID, permission string) (bool, error) {
	policies, err := l.casbinManager.GetRoleMenuMappings(roleID)
	if err != nil {
		return false, err
	}

	for _, policy := range policies {
		if len(policy) >= 3 {
			if policy[1] == menuID && policy[2] == permission {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetRoleMenuPermissions 获取角色的所有菜单权限
func (l *MenuPermissionLogicImpl) GetRoleMenuPermissions(roleID string) ([]*MenuPermission, error) {
	policies, err := l.casbinManager.GetRoleMenuMappings(roleID)
	if err != nil {
		return nil, err
	}

	permissions := make([]*MenuPermission, 0, len(policies))
	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, &MenuPermission{
				MenuID:     policy[1], // V1 是 menu_id
				Permission: policy[2], // V2 是 permission
			})
		}
	}

	return permissions, nil
}

// buildMenuTree 构建菜单树
func (l *MenuPermissionLogicImpl) buildMenuTree(menuIDs []string) ([]*identity_srv.MenuNode, error) {
	if len(menuIDs) == 0 {
		return []*identity_srv.MenuNode{}, nil
	}

	// 直接从 DAL 层获取最新的菜单树
	ctx := context.TODO() // 临时使用，应该从上层传递
	menuModels, err := l.dal.Menu().GetLatestMenuTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取完整菜单树失败: %w", err)
	}

	// 使用 Converter 转换为 Thrift 结构
	fullMenuTree := l.converter.Menu().ModelsToThrift(menuModels)

	// 创建菜单ID集合，用于快速查找
	menuIDSet := make(map[string]bool)
	for _, id := range menuIDs {
		menuIDSet[id] = true
	}

	// 过滤并构建用户可访问的菜单树
	return l.filterMenuTree(fullMenuTree, menuIDSet), nil
}

// filterMenuTree 递归过滤菜单树，只保留用户有权限的菜单
func (l *MenuPermissionLogicImpl) filterMenuTree(menuNodes []*identity_srv.MenuNode, menuIDSet map[string]bool) []*identity_srv.MenuNode {
	var result []*identity_srv.MenuNode

	for _, node := range menuNodes {
		// 检查当前节点是否有权限
		if node.Id != nil && menuIDSet[*node.Id] {
			// 复制节点
			newNode := *node

			// 递归处理子菜单
			if len(node.Children) > 0 {
				newNode.Children = l.filterMenuTree(node.Children, menuIDSet)
			}

			result = append(result, &newNode)
		} else {
			// 即使当前节点没有权限，也要检查子菜单
			filteredChildren := l.filterMenuTree(node.Children, menuIDSet)
			if len(filteredChildren) > 0 {
				// 如果有子菜单有权限，则包含父菜单
				newNode := *node
				newNode.Children = filteredChildren
				result = append(result, &newNode)
			}
		}
	}

	return result
}
