package menu

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	assert.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
}

func TestConverterImpl_ModelToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.ModelToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("根菜单节点转换", func(t *testing.T) {
		now := int64(1640995200000) // 2022-01-01 00:00:00 UTC
		rootID := uuid.New()
		creatorID := uuid.New()

		model := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        rootID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "dashboard",
			Version:    "v1.0.0",
			Name:       "仪表盘",
			Path:       "/dashboard",
			Component:  "Dashboard",
			Icon:       "dashboard",
			ParentID:   nil,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{},
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "dashboard", *result.Id)
		assert.Equal(t, "仪表盘", *result.Name)
		assert.Equal(t, "/dashboard", *result.Path)
		assert.Equal(t, "Dashboard", *result.Component)
		assert.Equal(t, "dashboard", *result.Icon)
		assert.Empty(t, result.Children) // 空子菜单列表
	})

	t.Run("带子菜单的菜单转换", func(t *testing.T) {
		now := int64(1640995200000)
		parentID := uuid.New()
		childID := uuid.New()
		creatorID := uuid.New()

		childMenu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        childID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "dashboard-overview",
			Version:    "v1.0.0",
			Name:       "概览",
			Path:       "/dashboard/overview",
			Component:  "Dashboard/Overview",
			Icon:       "overview",
			ParentID:   &parentID,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{},
		}

		parentMenu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        parentID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "dashboard",
			Version:    "v1.0.0",
			Name:       "仪表盘",
			Path:       "/dashboard",
			Component:  "Dashboard",
			Icon:       "dashboard",
			ParentID:   nil,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{childMenu},
		}

		result := converter.ModelToThrift(parentMenu)

		require.NotNil(t, result)
		assert.Equal(t, "dashboard", *result.Id)
		assert.Equal(t, "仪表盘", *result.Name)
		assert.Equal(t, "/dashboard", *result.Path)
		assert.Equal(t, "Dashboard", *result.Component)
		assert.Equal(t, "dashboard", *result.Icon)

		// 验证子菜单
		require.Len(t, result.Children, 1)
		childResult := result.Children[0]
		assert.Equal(t, "dashboard-overview", *childResult.Id)
		assert.Equal(t, "概览", *childResult.Name)
		assert.Equal(t, "/dashboard/overview", *childResult.Path)
		assert.Equal(t, "Dashboard/Overview", *childResult.Component)
		assert.Equal(t, "overview", *childResult.Icon)
	})

	t.Run("多级嵌套菜单转换", func(t *testing.T) {
		now := int64(1640995200000)
		rootID := uuid.New()
		level1ID := uuid.New()
		level2ID := uuid.New()
		creatorID := uuid.New()

		// 第三级菜单
		level2Menu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        level2ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "system-settings-users",
			Version:    "v1.0.0",
			Name:       "用户管理",
			Path:       "/system/settings/users",
			Component:  "System/Settings/Users",
			Icon:       "users",
			ParentID:   &level1ID,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{},
		}

		// 第二级菜单
		level1Menu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        level1ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "system-settings",
			Version:    "v1.0.0",
			Name:       "系统设置",
			Path:       "/system/settings",
			Component:  "System/Settings",
			Icon:       "settings",
			ParentID:   &rootID,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{level2Menu},
		}

		// 第一级菜单
		rootMenu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        rootID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "system",
			Version:    "v1.0.0",
			Name:       "系统管理",
			Path:       "/system",
			Component:  "System",
			Icon:       "system",
			ParentID:   nil,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{level1Menu},
		}

		result := converter.ModelToThrift(rootMenu)

		require.NotNil(t, result)
		assert.Equal(t, "system", *result.Id)
		assert.Equal(t, "系统管理", *result.Name)

		// 验证第一级子菜单
		require.Len(t, result.Children, 1)
		level1Result := result.Children[0]
		assert.Equal(t, "system-settings", *level1Result.Id)
		assert.Equal(t, "系统设置", *level1Result.Name)

		// 验证第二级子菜单
		require.Len(t, level1Result.Children, 1)
		level2Result := level1Result.Children[0]
		assert.Equal(t, "system-settings-users", *level2Result.Id)
		assert.Equal(t, "用户管理", *level2Result.Name)
		assert.Equal(t, "/system/settings/users", *level2Result.Path)
	})

	t.Run("空字段处理", func(t *testing.T) {
		now := int64(1640995200000)
		menuID := uuid.New()
		creatorID := uuid.New()

		model := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        menuID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "empty-fields",
			Version:    "v1.0.0",
			Name:       "", // 空名称
			Path:       "", // 空路径
			Component:  "", // 空组件
			Icon:       "", // 空图标
			ParentID:   nil,
			Sort:       0,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{},
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		assert.Equal(t, "empty-fields", *result.Id)
		assert.Equal(t, "", *result.Name)      // 空值被保留
		assert.Equal(t, "", *result.Path)      // 空值被保留
		assert.Equal(t, "", *result.Component) // 空值被保留
		assert.Equal(t, "", *result.Icon)      // 空值被保留
	})
}

func TestConverterImpl_ModelsToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil切片转换", func(t *testing.T) {
		result := converter.ModelsToThrift(nil)
		assert.Nil(t, result)
	})

	t.Run("空切片转换", func(t *testing.T) {
		result := converter.ModelsToThrift([]*models.Menu{})
		assert.Nil(t, result)
	})

	t.Run("多个菜单转换", func(t *testing.T) {
		now := int64(1640995200000)
		menu1ID := uuid.New()
		menu2ID := uuid.New()
		creatorID := uuid.New()

		menu1 := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        menu1ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "dashboard",
			Version:    "v1.0.0",
			Name:       "仪表盘",
			Path:       "/dashboard",
			Component:  "Dashboard",
			Icon:       "dashboard",
			ParentID:   nil,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{},
		}

		menu2 := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        menu2ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "system",
			Version:    "v1.0.0",
			Name:       "系统管理",
			Path:       "/system",
			Component:  "System",
			Icon:       "system",
			ParentID:   nil,
			Sort:       2,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{},
		}

		models := []*models.Menu{menu1, menu2}
		result := converter.ModelsToThrift(models)

		require.NotNil(t, result)
		require.Len(t, result, 2)

		// 验证第一个菜单
		assert.Equal(t, "dashboard", *result[0].Id)
		assert.Equal(t, "仪表盘", *result[0].Name)
		assert.Equal(t, "/dashboard", *result[0].Path)

		// 验证第二个菜单
		assert.Equal(t, "system", *result[1].Id)
		assert.Equal(t, "系统管理", *result[1].Name)
		assert.Equal(t, "/system", *result[1].Path)
	})

	t.Run("包含nil元素的切片", func(t *testing.T) {
		now := int64(1640995200000)
		menuID := uuid.New()
		creatorID := uuid.New()

		menu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        menuID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "dashboard",
			Version:    "v1.0.0",
			Name:       "仪表盘",
			Path:       "/dashboard",
			Component:  "Dashboard",
			Icon:       "dashboard",
			ParentID:   nil,
			Sort:       1,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
			Children:   []*models.Menu{},
		}

		models := []*models.Menu{menu, nil, menu} // 包含nil元素
		result := converter.ModelsToThrift(models)

		require.NotNil(t, result)
		require.Len(t, result, 3)

		// 第一个菜单应该正常转换
		assert.Equal(t, "dashboard", *result[0].Id)
		assert.Equal(t, "仪表盘", *result[0].Name)

		// 第二个元素应该是nil
		assert.Nil(t, result[1])

		// 第三个菜单应该正常转换
		assert.Equal(t, "dashboard", *result[2].Id)
		assert.Equal(t, "仪表盘", *result[2].Name)
	})
}

// 递归菜单结构测试
func TestConverterImpl_RecursiveMenuStructure(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("复杂递归菜单结构", func(t *testing.T) {
		now := int64(1640995200000)

		// 创建一个复杂的菜单树结构
		// 系统
		//   ├── 用户管理
		//   │   ├── 用户列表
		//   │   └── 用户详情
		//   ├── 角色管理
		//   └── 权限管理
		//       ├── 权限列表
		//       └── 权限配置

		// 权限管理 -> 权限配置
		permissionConfig := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "permission-config",
			Version:    "v1.0.0",
			Name:       "权限配置",
			Path:       "/system/permission/config",
			Component:  "System/Permission/Config",
			Icon:       "config",
			Children:   []*models.Menu{},
		}

		// 权限管理 -> 权限列表
		permissionList := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "permission-list",
			Version:    "v1.0.0",
			Name:       "权限列表",
			Path:       "/system/permission/list",
			Component:  "System/Permission/List",
			Icon:       "list",
			Children:   []*models.Menu{},
		}

		// 权限管理
		permission := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "permission",
			Version:    "v1.0.0",
			Name:       "权限管理",
			Path:       "/system/permission",
			Component:  "System/Permission",
			Icon:       "permission",
			Children:   []*models.Menu{permissionList, permissionConfig},
		}

		// 用户管理 -> 用户详情
		userDetail := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "user-detail",
			Version:    "v1.0.0",
			Name:       "用户详情",
			Path:       "/system/user/detail",
			Component:  "System/User/Detail",
			Icon:       "detail",
			Children:   []*models.Menu{},
		}

		// 用户管理 -> 用户列表
		userList := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "user-list",
			Version:    "v1.0.0",
			Name:       "用户列表",
			Path:       "/system/user/list",
			Component:  "System/User/List",
			Icon:       "list",
			Children:   []*models.Menu{},
		}

		// 用户管理
		user := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "user",
			Version:    "v1.0.0",
			Name:       "用户管理",
			Path:       "/system/user",
			Component:  "System/User",
			Icon:       "user",
			Children:   []*models.Menu{userList, userDetail},
		}

		// 角色管理
		role := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "role",
			Version:    "v1.0.0",
			Name:       "角色管理",
			Path:       "/system/role",
			Component:  "System/Role",
			Icon:       "role",
			Children:   []*models.Menu{},
		}

		// 系统
		system := &models.Menu{
			BaseModel:  models.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now},
			SemanticID: "system",
			Version:    "v1.0.0",
			Name:       "系统管理",
			Path:       "/system",
			Component:  "System",
			Icon:       "system",
			Children:   []*models.Menu{user, role, permission},
		}

		result := converter.ModelToThrift(system)

		require.NotNil(t, result)
		assert.Equal(t, "system", *result.Id)
		assert.Equal(t, "系统管理", *result.Name)

		// 验证一级子菜单
		require.Len(t, result.Children, 3)

		// 用户管理
		userResult := result.Children[0]
		assert.Equal(t, "user", *userResult.Id)
		assert.Equal(t, "用户管理", *userResult.Name)
		require.Len(t, userResult.Children, 2)

		// 用户列表
		userListResult := userResult.Children[0]
		assert.Equal(t, "user-list", *userListResult.Id)
		assert.Equal(t, "用户列表", *userListResult.Name)
		assert.Empty(t, userListResult.Children)

		// 用户详情
		userDetailResult := userResult.Children[1]
		assert.Equal(t, "user-detail", *userDetailResult.Id)
		assert.Equal(t, "用户详情", *userDetailResult.Name)
		assert.Empty(t, userDetailResult.Children)

		// 角色管理
		roleResult := result.Children[1]
		assert.Equal(t, "role", *roleResult.Id)
		assert.Equal(t, "角色管理", *roleResult.Name)
		assert.Empty(t, roleResult.Children)

		// 权限管理
		permissionResult := result.Children[2]
		assert.Equal(t, "permission", *permissionResult.Id)
		assert.Equal(t, "权限管理", *permissionResult.Name)
		require.Len(t, permissionResult.Children, 2)

		// 权限列表
		permissionListResult := permissionResult.Children[0]
		assert.Equal(t, "permission-list", *permissionListResult.Id)
		assert.Equal(t, "权限列表", *permissionListResult.Name)
		assert.Empty(t, permissionListResult.Children)

		// 权限配置
		permissionConfigResult := permissionResult.Children[1]
		assert.Equal(t, "permission-config", *permissionConfigResult.Id)
		assert.Equal(t, "权限配置", *permissionConfigResult.Name)
		assert.Empty(t, permissionConfigResult.Children)
	})
}

// 边界情况和错误处理测试
func TestConverterImpl_EdgeCases(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("深层嵌套菜单", func(t *testing.T) {
		now := int64(1640995200000)

		// 创建5层深度的菜单结构
		var currentMenu *models.Menu

		for i := 0; i < 5; i++ {
			level := i + 1
			menu := &models.Menu{
				BaseModel: models.BaseModel{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				SemanticID: fmt.Sprintf("level-%d", level),
				Version:    "v1.0.0",
				Name:       fmt.Sprintf("级别 %d", level),
				Path:       fmt.Sprintf("/level/%d", level),
				Component:  fmt.Sprintf("Level%d", level),
				Icon:       fmt.Sprintf("level-%d", level),
				Children:   []*models.Menu{},
			}

			if currentMenu != nil {
				menu.Children = []*models.Menu{currentMenu}
			}

			currentMenu = menu
		}

		result := converter.ModelToThrift(currentMenu)

		require.NotNil(t, result)
		assert.Equal(t, "level-5", *result.Id)
		assert.Equal(t, "级别 5", *result.Name)

		// 验证5层深度
		current := result
		for i := 5; i >= 1; i-- {
			require.NotNil(t, current)

			expectedID := fmt.Sprintf("level-%d", i)
			expectedName := fmt.Sprintf("级别 %d", i)

			assert.Equal(t, expectedID, *current.Id)
			assert.Equal(t, expectedName, *current.Name)

			if i > 1 {
				require.Len(t, current.Children, 1)
				current = current.Children[0]
			} else {
				assert.Empty(t, current.Children)
			}
		}
	})

	t.Run("大量子菜单", func(t *testing.T) {
		now := int64(1640995200000)

		// 创建包含100个子菜单的父菜单
		var children []*models.Menu

		for i := 0; i < 100; i++ {
			child := &models.Menu{
				BaseModel: models.BaseModel{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				SemanticID: fmt.Sprintf("child-%d", i),
				Version:    "v1.0.0",
				Name:       fmt.Sprintf("子菜单 %d", i),
				Path:       fmt.Sprintf("/child/%d", i),
				Component:  fmt.Sprintf("Child%d", i),
				Icon:       fmt.Sprintf("child-%d", i),
				Children:   []*models.Menu{},
			}
			children = append(children, child)
		}

		parent := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "parent",
			Version:    "v1.0.0",
			Name:       "父菜单",
			Path:       "/parent",
			Component:  "Parent",
			Icon:       "parent",
			Children:   children,
		}

		result := converter.ModelToThrift(parent)

		require.NotNil(t, result)
		assert.Equal(t, "parent", *result.Id)
		assert.Equal(t, "父菜单", *result.Name)
		require.Len(t, result.Children, 100)

		// 验证前几个和后几个子菜单
		assert.Equal(t, "child-0", *result.Children[0].Id)
		assert.Equal(t, "子菜单 0", *result.Children[0].Name)
		assert.Equal(t, "child-99", *result.Children[99].Id)
		assert.Equal(t, "子菜单 99", *result.Children[99].Name)
	})
}

// 完整往返转换测试
func TestConverterImpl_CompleteRoundTrip(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("语义化ID验证", func(t *testing.T) {
		now := int64(1640995200000)
		menuID := uuid.New()

		// 使用语义化ID而非UUID
		model := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        menuID, // 这是UUID，但不应该被使用
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: "user-management-v2", // 语义化ID
			Version:    "v2.0.0",
			Name:       "用户管理 v2",
			Path:       "/users/v2",
			Component:  "Users/V2",
			Icon:       "users-v2",
			Children:   []*models.Menu{},
		}

		result := converter.ModelToThrift(model)

		require.NotNil(t, result)
		// 验证使用的是SemanticID而不是UUID
		assert.Equal(t, "user-management-v2", *result.Id)
		assert.NotEqual(t, menuID.String(), *result.Id)
		assert.Equal(t, "用户管理 v2", *result.Name)
		assert.Equal(t, "/users/v2", *result.Path)
		assert.Equal(t, "Users/V2", *result.Component)
		assert.Equal(t, "users-v2", *result.Icon)
	})
}

// 基准测试
func BenchmarkConverterImpl_ModelToThrift(b *testing.B) {
	converter := &ConverterImpl{}

	now := int64(1640995200000)
	menuID := uuid.New()

	model := &models.Menu{
		BaseModel: models.BaseModel{
			ID:        menuID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		SemanticID: "dashboard",
		Version:    "v1.0.0",
		Name:       "仪表盘",
		Path:       "/dashboard",
		Component:  "Dashboard",
		Icon:       "dashboard",
		Children:   []*models.Menu{},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelToThrift(model)
	}
}

func BenchmarkConverterImpl_ModelsToThrift(b *testing.B) {
	converter := &ConverterImpl{}

	now := int64(1640995200000)

	var menus []*models.Menu

	for i := 0; i < 10; i++ {
		menu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: fmt.Sprintf("menu-%d", i),
			Version:    "v1.0.0",
			Name:       fmt.Sprintf("菜单 %d", i),
			Path:       fmt.Sprintf("/menu/%d", i),
			Component:  fmt.Sprintf("Menu%d", i),
			Icon:       fmt.Sprintf("menu-%d", i),
			Children:   []*models.Menu{},
		}
		menus = append(menus, menu)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelsToThrift(menus)
	}
}

func BenchmarkConverterImpl_RecursiveMenu(b *testing.B) {
	converter := &ConverterImpl{}

	now := int64(1640995200000)

	// 创建3层深的菜单树
	var createMenu func(depth int) *models.Menu

	createMenu = func(depth int) *models.Menu {
		menu := &models.Menu{
			BaseModel: models.BaseModel{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			SemanticID: fmt.Sprintf("level-%d", depth),
			Version:    "v1.0.0",
			Name:       fmt.Sprintf("级别 %d", depth),
			Path:       fmt.Sprintf("/level/%d", depth),
			Component:  fmt.Sprintf("Level%d", depth),
			Icon:       fmt.Sprintf("level-%d", depth),
			Children:   []*models.Menu{},
		}

		if depth > 0 {
			menu.Children = []*models.Menu{createMenu(depth - 1)}
		}

		return menu
	}

	rootMenu := createMenu(3) // 4层深度

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = converter.ModelToThrift(rootMenu)
	}
}
