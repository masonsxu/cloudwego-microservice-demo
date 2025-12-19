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

// Converter 接口聚合了所有实体的转换接口，提供统一的模型与 Thrift DTO 转换入口。
//
// 设计理念：
//   - 统一入口：业务逻辑层只需依赖一个 Converter 接口，简化依赖管理和 Wire 注入
//   - 职责分离：每个子转换器专注于单一实体的转换逻辑
//   - 复合操作：提供跨实体的复合转换方法，避免业务层直接操作多个转换器
//   - 扩展友好：新增实体转换器只需添加对应的访问方法
//
// 使用示例：
//
//	func NewUserLogic(dal dal.DAL, conv converter.Converter) UserLogic {
//	    return &userLogicImpl{
//	        dal:  dal,
//	        conv: conv,  // 通过 conv.UserProfile() 访问用户转换器
//	    }
//	}
//
// 注意事项：
//   - 所有子转换器都是无状态的纯函数，初始化成本极低
//   - 如需独立使用某个转换器，可使用包级工厂函数（如 posture.NewPostureConverter）
//   - 复合转换方法应放在此接口中统一管理
type Converter interface {
	// ============================================================================
	// 核心实体转换器 - 用户身份与权限
	// ============================================================================

	// UserProfile 用户档案转换器
	// 负责 UserProfile Model ↔ Thrift DTO 的双向转换
	// 包含基本信息、认证状态、资质等字段
	UserProfile() user.Converter

	// Membership 用户成员关系转换器
	// 负责 UserMembership Model ↔ Thrift DTO 的转换
	// 处理用户与组织/部门的关联关系、角色权限分配
	Membership() membership.Converter

	// Authentication 认证相关转换器
	// 负责登录请求、密码修改等认证场景的数据转换
	// 包含登录响应的构建逻辑
	Authentication() authentication.Converter

	// ============================================================================
	// 组织架构转换器 - 机构管理
	// ============================================================================

	// Organization 组织/机构转换器
	// 负责 Organization Model ↔ Thrift DTO 的转换
	// 处理机构的层级结构、元数据等信息
	Organization() organization.Converter

	// Department 部门转换器
	// 负责 Department Model ↔ Thrift DTO 的转换
	// 处理机构内部门的层级关系和成员管理
	Department() department.Converter

	// ============================================================================
	// 业务领域转换器 - 专业功能模块
	// ============================================================================

	// Logo 组织Logo转换器
	// 负责 OrganizationLogo Model ↔ Thrift DTO 的转换
	// 处理组织Logo的上传、绑定、状态管理等场景
	Logo() logo.Converter
	Menu() menu.Converter
	RoleDefinition() definition.Converter
	UserRoleAssignment() assignment.Converter
	// ============================================================================
	// 基础设施转换器 - 通用工具
	// ============================================================================

	// Enum 枚举类型映射器
	// 负责枚举值在 Model 和 Thrift 之间的双向映射
	// 如用户状态、性别、角色类型等枚举
	Enum() enum.Converter

	// Base 基础类型转换器
	// 负责通用类型的转换，如分页参数、时间戳、UUID等
	// 提供可复用的基础转换函数
	Base() base.Converter

	// ============================================================================
	// 复合转换方法 - 跨实体协作
	// ============================================================================

	// BuildLoginResponse 构建登录响应
	// 需要协调 UserProfile 和 Membership 两个转换器
	// 将用户档案和成员关系组合为完整的登录响应
	//
	// 参数:
	//   - userProfile: 用户档案 Model
	//   - memberships: 用户成员关系列表
	//
	// 返回:
	//   - *identity_srv.LoginResponse: 完整的登录响应 DTO
	BuildLoginResponse(
		userProfile *models.UserProfile,
		memberships []*models.UserMembership,
	) *identity_srv.LoginResponse
}
