package casbin

import (
	"context"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/config"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/log"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// CasbinManager Casbin 权限管理器
type CasbinManager struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
	logger   *zerolog.Logger
}

// NewCasbinManager 创建 Casbin 管理器
func NewCasbinManager(
	db *gorm.DB,
	config *config.CasbinConfig,
	logger *zerolog.Logger,
) (*CasbinManager, error) {
	// 1. 确保数据表存在
	err := db.AutoMigrate(&models.CasbinRule{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate casbin rules table: %w", err)
	}

	// 2. 创建 GORM Adapter
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm adapter: %w", err)
	}

	// 3. 创建 Casbin 执行器
	enforcer, err := casbin.NewEnforcer(config.ModelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// 4. 配置执行器
	enforcer.EnableAutoSave(true)
	enforcer.EnableAutoBuildRoleLinks(true)
	enforcer.EnableLog(config.EnableLog)

	// 5. 加载策略数据
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	logger.Info().
		Str("model_path", config.ModelPath).
		Bool("auto_save", true).
		Bool("auto_build_role_links", true).
		Msg("Casbin manager initialized successfully")

	return &CasbinManager{
		enforcer: enforcer,
		db:       db,
		logger:   logger,
	}, nil
}

// GetEnforcer 获取 Casbin Enforcer 实例
func (cm *CasbinManager) GetEnforcer() *casbin.Enforcer {
	return cm.enforcer
}

// AddUserRole 为用户添加角色
func (cm *CasbinManager) AddUserRole(userID, roleID string) error {
	added, err := cm.enforcer.AddRoleForUser(userID, roleID)
	if err != nil {
		return fmt.Errorf("添加用户角色失败: %w", err)
	}

	if !added {
		cm.logger.Debug().
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("用户角色关系已存在")
	} else {
		cm.logger.Info().
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("用户角色添加成功")
	}

	return nil
}

// AddUserRoleWithContext 为用户添加角色（带追踪ID的示例方法）
// 这是一个展示如何在业务逻辑中集成追踪 ID 的示例方法
func (cm *CasbinManager) AddUserRoleWithContext(ctx context.Context, userID, roleID string) error {
	// 使用 log.Ctx 获取带追踪信息的 logger
	logger := log.Ctx(ctx)

	added, err := cm.enforcer.AddRoleForUser(userID, roleID)
	if err != nil {
		logger.Error().
			Str("user_id", userID).
			Str("role_id", roleID).
			Err(err).
			Msg("添加用户角色失败")
		return fmt.Errorf("添加用户角色失败: %w", err)
	}

	if !added {
		logger.Debug().
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("用户角色关系已存在")
	} else {
		logger.Info().
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("用户角色添加成功")
	}

	return nil
}

// RemoveUserRole 移除用户角色
func (cm *CasbinManager) RemoveUserRole(userID, roleID string) error {
	removed, err := cm.enforcer.DeleteRoleForUser(userID, roleID)
	if err != nil {
		return fmt.Errorf("移除用户角色失败: %w", err)
	}

	if !removed {
		cm.logger.Debug().
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("用户角色关系不存在")
	} else {
		cm.logger.Info().
			Str("user_id", userID).
			Str("role_id", roleID).
			Msg("用户角色移除成功")
	}

	return nil
}

// GetUserRoles 获取用户的所有角色
func (cm *CasbinManager) GetUserRoles(userID string) ([]string, error) {
	roles, err := cm.enforcer.GetRolesForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	cm.logger.Debug().
		Str("user_id", userID).
		Interface("roles", roles).
		Msg("获取用户角色")
	return roles, nil
}

// AddRolePermission 为角色添加权限
func (cm *CasbinManager) AddRolePermission(roleID, resource, action string) error {
	added, err := cm.enforcer.AddPolicy(roleID, resource, action)
	if err != nil {
		return fmt.Errorf("添加角色权限失败: %w", err)
	}

	if !added {
		cm.logger.Debug().
			Str("role_id", roleID).
			Str("resource", resource).
			Str("action", action).
			Msg("角色权限已存在")
	} else {
		cm.logger.Info().
			Str("role_id", roleID).
			Str("resource", resource).
			Str("action", action).
			Msg("角色权限添加成功")
	}

	return nil
}

// RemoveRolePermission 移除角色权限
func (cm *CasbinManager) RemoveRolePermission(roleID, resource, action string) error {
	removed, err := cm.enforcer.RemovePolicy(roleID, resource, action)
	if err != nil {
		return fmt.Errorf("移除角色权限失败: %w", err)
	}

	if !removed {
		cm.logger.Debug().
			Str("role_id", roleID).
			Str("resource", resource).
			Str("action", action).
			Msg("角色权限不存在")
	} else {
		cm.logger.Info().
			Str("role_id", roleID).
			Str("resource", resource).
			Str("action", action).
			Msg("角色权限移除成功")
	}

	return nil
}

// HasPermission 检查用户是否具有指定权限
func (cm *CasbinManager) HasPermission(userID, resource, action string) (bool, error) {
	return cm.enforcer.Enforce(userID, resource, action)
}

// GetUserPermissions 获取用户的所有权限
func (cm *CasbinManager) GetUserPermissions(userID string) ([][]string, error) {
	permissions, err := cm.enforcer.GetImplicitPermissionsForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户权限失败: %w", err)
	}

	return permissions, nil
}

// AddRoleMenuMapping 为角色添加菜单映射
// roleID: 角色ID
// menuID: 语义化菜单ID (来自menu.yaml，版本间保持稳定)
// permission: 权限类型，必须是以下之一：
//   - view_own_organization: 查看所在组织
//   - view_all_organizations: 查看所有组织
func (cm *CasbinManager) AddRoleMenuMapping(roleID, menuID, permission string) error {
	// 验证权限类型是否有效
	if !models.IsValidMenuPermission(permission) {
		return fmt.Errorf("无效的权限类型: %s，有效值为: view_own_organization, view_all_organizations", permission)
	}

	// 使用 p2 策略类型存储：角色ID -> 语义化菜单ID -> 权限类型
	added, err := cm.enforcer.AddNamedPolicy("p2", roleID, menuID, permission)
	if err != nil {
		return fmt.Errorf("添加角色菜单映射失败: %w", err)
	}

	if !added {
		cm.logger.Debug().
			Str("role_id", roleID).
			Str("menu_id", menuID).
			Str("permission", permission).
			Msg("角色菜单映射已存在")
	} else {
		cm.logger.Info().
			Str("role_id", roleID).
			Str("menu_id", menuID).
			Str("permission", permission).
			Msg("角色菜单映射添加成功")
	}

	return nil
}

// RemoveRoleMenuMapping 移除角色菜单映射
func (cm *CasbinManager) RemoveRoleMenuMapping(roleID, menuID, permission string) error {
	removed, err := cm.enforcer.RemoveNamedPolicy("p2", roleID, menuID, permission)
	if err != nil {
		return fmt.Errorf("移除角色菜单映射失败: %w", err)
	}

	if !removed {
		cm.logger.Debug().
			Str("role_id", roleID).
			Str("menu_id", menuID).
			Str("permission", permission).
			Msg("角色菜单映射不存在")
	} else {
		cm.logger.Info().
			Str("role_id", roleID).
			Str("menu_id", menuID).
			Str("permission", permission).
			Msg("角色菜单映射移除成功")
	}

	return nil
}

// GetRoleMenuMappings 获取角色的菜单映射
// 返回的策略格式：[ptype, role_id, semantic_menu_id, permission]
// 其中 semantic_menu_id 是来自menu.yaml的语义化ID，版本间保持稳定
func (cm *CasbinManager) GetRoleMenuMappings(roleID string) ([][]string, error) {
	// 获取角色的所有 p2 类型策略：[ptype, role_id, semantic_menu_id, permission]
	policies, err := cm.enforcer.GetFilteredNamedPolicy("p2", 0, roleID)
	if err != nil {
		return nil, fmt.Errorf("获取角色菜单映射失败: %w", err)
	}
	return policies, nil
}

// ClearRoleMenuMappings 清空角色的所有菜单映射
func (cm *CasbinManager) ClearRoleMenuMappings(roleID string) error {
	removed, err := cm.enforcer.RemoveFilteredNamedPolicy("p2", 0, roleID)
	if err != nil {
		return fmt.Errorf("清空角色菜单映射失败: %w", err)
	}

	cm.logger.Info().
		Str("role_id", roleID).
		Bool("removed", removed).
		Msg("角色菜单映射清空成功")
	return nil
}

// SyncUserRoles 同步用户角色（替换现有角色）
func (cm *CasbinManager) SyncUserRoles(userID string, roleIDs []string) error {
	// 1. 清除用户的旧角色关系
	_, err := cm.enforcer.DeleteRolesForUser(userID)
	if err != nil {
		return fmt.Errorf("清除用户旧角色关系失败: %w", err)
	}

	// 2. 添加新的角色关系
	for _, roleID := range roleIDs {
		err = cm.AddUserRole(userID, roleID)
		if err != nil {
			cm.logger.Error().
				Str("user_id", userID).
				Str("role_id", roleID).
				Err(err).
				Msg("为用户添加角色失败")
			// 继续处理其他角色，不直接返回错误
		}
	}

	cm.logger.Info().
		Str("user_id", userID).
		Interface("role_ids", roleIDs).
		Msg("用户角色同步成功")

	return nil
}

// SavePolicy 保存策略到数据库
func (cm *CasbinManager) SavePolicy() error {
	return cm.enforcer.SavePolicy()
}

// LoadPolicy 从数据库加载策略
func (cm *CasbinManager) LoadPolicy() error {
	return cm.enforcer.LoadPolicy()
}
