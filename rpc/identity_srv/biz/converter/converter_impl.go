package converter

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/authentication"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/definition"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/department"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/enum"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/logo"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/membership"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/menu"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/organization"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/user"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Impl 实现了 Converter 接口，聚合所有实体的转换器。
//
// 架构说明：
//   - 采用 Facade 模式，为业务层提供统一的转换服务入口
//   - 所有子转换器都是无状态的，可安全并发使用
//   - 初始化时一次性创建所有转换器实例，避免重复创建
//   - 通过组合模式复用 enumMapper 和 baseConverter
//
// 性能特性：
//   - 转换器实例化成本极低（纯函数对象）
//   - 无锁设计，支持高并发场景
//   - 内存占用小于 1KB
type Impl struct {
	// ============================================================================
	// 核心实体转换器 - 用户身份与权限
	// ============================================================================
	authenticationConverter authentication.Converter
	userProfileConverter    user.Converter
	membershipConverter     membership.Converter

	// ============================================================================
	// 组织架构转换器 - 机构管理
	// ============================================================================
	organizationConverter organization.Converter
	departmentConverter   department.Converter

	// ============================================================================
	// 业务领域转换器 - 专业功能模块
	// ============================================================================
	logoConverter               logo.Converter
	menuConverter               menu.Converter
	roleDefinitionConverter     definition.Converter
	userRoleAssignmentConverter assignment.Converter
	// ============================================================================
	// 基础设施转换器 - 通用工具
	// ============================================================================
	enumConverter enum.Converter
	baseConverter base.Converter
}

// NewConverter 创建转换器聚合实例
//
// 初始化流程：
//  1. 优先创建基础设施转换器（enum、base）
//  2. 将基础转换器注入到需要的实体转换器中（如 userProfileConverter）
//  3. 组装所有转换器形成完整的 Converter 实例
//
// 返回值：
//   - Converter: 可立即使用的转换器聚合接口
//
// 使用示例：
//
//	// 在业务逻辑层初始化
//	conv := converter.NewConverter()
//	userThrift := conv.UserProfile().ModelUserProfileToThrift(userModel)
func NewConverter() Converter {
	// 基础设施转换器（被其他转换器依赖）
	enumConverter := enum.NewConverter()
	baseConverter := base.NewConverter()

	return &Impl{
		// 核心实体转换器
		authenticationConverter: authentication.NewConverter(),
		userProfileConverter:    user.NewConverter(enumConverter),
		membershipConverter:     membership.NewConverter(enumConverter),

		// 组织架构转换器
		organizationConverter: organization.NewConverter(),
		departmentConverter:   department.NewConverter(),

		// 业务领域转换器
		logoConverter:               logo.NewConverter(),
		menuConverter:               menu.NewConverter(),
		roleDefinitionConverter:     definition.NewConverter(enumConverter),
		userRoleAssignmentConverter: assignment.NewConverter(),
		// 基础设施转换器
		enumConverter: enumConverter,
		baseConverter: baseConverter,
	}
}

// ============================================================================
// 子转换器访问方法 - 核心实体
// ============================================================================

// UserProfile 返回用户档案转换器
func (c *Impl) UserProfile() user.Converter {
	return c.userProfileConverter
}

// Membership 返回成员关系转换器
func (c *Impl) Membership() membership.Converter {
	return c.membershipConverter
}

// Authentication 返回认证相关转换器
func (c *Impl) Authentication() authentication.Converter {
	return c.authenticationConverter
}

// ============================================================================
// 子转换器访问方法 - 组织架构
// ============================================================================

// Organization 返回组织转换器
func (c *Impl) Organization() organization.Converter {
	return c.organizationConverter
}

// Department 返回部门转换器
func (c *Impl) Department() department.Converter {
	return c.departmentConverter
}

// ============================================================================
// 子转换器访问方法 - 业务领域
// ============================================================================

// Logo 返回组织Logo转换器
func (c *Impl) Logo() logo.Converter {
	return c.logoConverter
}

// Menu 返回菜单转换器
func (c *Impl) Menu() menu.Converter {
	return c.menuConverter
}

// RoleDefinition 返回角色定义转换器
func (c *Impl) RoleDefinition() definition.Converter {
	return c.roleDefinitionConverter
}

// UserRoleAssignment 返回用户角色分配转换器
func (c *Impl) UserRoleAssignment() assignment.Converter {
	return c.userRoleAssignmentConverter
}

// ============================================================================
// 子转换器访问方法 - 基础设施
// ============================================================================

// Enum 返回枚举映射器
func (c *Impl) Enum() enum.Converter {
	return c.enumConverter
}

// Base 返回基础转换器
func (c *Impl) Base() base.Converter {
	return c.baseConverter
}

// ============================================================================
// 复合转换方法 - 跨实体协作
// ============================================================================

// BuildLoginResponse 构建登录响应（用户档案 + 成员关系）
//
// 业务场景：
//   - 用户登录成功后，需要返回用户基本信息和所属组织/部门信息
//   - 涉及 UserProfile 和 UserMembership 两个实体的转换和组合
//
// 处理逻辑：
//  1. 转换用户档案（包含基本信息、认证状态等）
//  2. 转换成员关系列表（可能为空）
//  3. 组装为符合 IDL 定义的 LoginResponse 结构
//
// 参数：
//   - userProfile: 用户档案 Model（不能为 nil）
//   - memberships: 用户成员关系列表（可为空）
//
// 返回值：
//   - *identity_srv.LoginResponse: 完整的登录响应 DTO
//
// 注意事项：
//   - 如果 memberships 为空，返回空数组而非 nil，避免客户端空指针
func (c *Impl) BuildLoginResponse(
	userProfile *models.UserProfile,
	memberships []*models.UserMembership,
) *identity_srv.LoginResponse {
	resp := &identity_srv.LoginResponse{}

	// 转换用户档案
	if userProfile != nil {
		resp.UserProfile = c.userProfileConverter.ModelUserProfileToThrift(userProfile)
	}

	// 转换成员关系（确保返回空数组而非 nil）
	if len(memberships) > 0 {
		resp.Memberships = c.authenticationConverter.ModelUserMembershipsToThrift(memberships)
	} else {
		resp.Memberships = []*identity_srv.UserMembership{}
	}

	return resp
}
