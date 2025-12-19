package membership

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// Logic 定义用户成员关系管理的业务逻辑接口
//
// 该接口负责用户与组织、部门之间的成员关系管理，包括：
// - 成员关系的创建、更新、删除
// - 成员关系的查询和检查
// - 主要成员关系的管理
//
// 所有方法严格按照 IDL identity_service.thrift 定义实现，保持接口契约的一致性。
// 成员关系表示一个用户在特定组织或部门中的角色和状态。
type MembershipLogic interface {
	// ============================================================================
	// 成员关系生命周期管理（与IDL identity_service.thrift 完全对应）
	// ============================================================================

	// AddMembership 为用户添加新的组织成员关系
	//
	// 该方法会：
	// 1. 验证用户、组织和部门的存在性
	// 2. 检查角色定义的有效性
	// 3. 确保不存在冲突的成员关系
	// 4. 在事务中创建成员关系，处理主要关系的唯一性约束
	//
	// 对应IDL: AddMembership(1: AddMembershipRequest req)
	//
	// 参数:
	//   - ctx: 上下文，用于传递请求ID、用户信息等
	//   - req: 添加成员关系的请求，包含用户ID、组织ID、角色等信息
	//
	// 返回:
	//   - *identity_srv.UserMembership: 成功创建的成员关系信息
	//   - error: 创建失败时的错误信息，包括参数验证、业务规则违反等
	AddMembership(
		ctx context.Context,
		req *identity_srv.AddMembershipRequest,
	) (*identity_srv.UserMembership, error)

	// UpdateMembership 更新用户的组织成员关系
	//
	// 该方法支持部分字段更新，包括：
	// - 角色变更（需验证新角色的有效性）
	// - 部门调动（需验证新部门属于同一组织）
	// - 状态变更（活跃、暂停等）
	// - 主要关系标记的变更
	// - 有效期的调整
	//
	// 对应IDL: UpdateMembership(1: UpdateMembershipRequest req)
	//
	// 参数:
	//   - ctx: 上下文信息
	//   - req: 更新请求，包含成员关系ID和需要更新的字段
	//
	// 返回:
	//   - *identity_srv.UserMembership: 更新后的成员关系信息
	//   - error: 更新失败时的错误信息
	UpdateMembership(
		ctx context.Context,
		req *identity_srv.UpdateMembershipRequest,
	) (*identity_srv.UserMembership, error)

	// RemoveMembership 移除用户的组织成员关系（逻辑删除）
	//
	// 该方法执行逻辑删除，不会物理删除数据记录，便于审计和数据恢复。
	// 删除操作会在事务中执行，确保数据一致性。
	//
	// 对应IDL: RemoveMembership(1: core.ULID membershipID)
	//
	// 参数:
	//   - ctx: 上下文信息
	//   - membershipID: 要删除的成员关系ID
	//
	// 返回:
	//   - error: 删除失败时的错误信息，如成员关系不存在等
	RemoveMembership(ctx context.Context, membershipID string) error

	// GetMembership 根据ID获取成员关系详情
	//
	// 该方法返回指定ID的成员关系完整信息，包括关联的用户、组织、部门信息。
	//
	// 对应IDL: GetMembership(1: core.ULID membershipID)
	//
	// 参数:
	//   - ctx: 上下文信息
	//   - membershipID: 成员关系的唯一标识符
	//
	// 返回:
	//   - *identity_srv.UserMembership: 查找到的成员关系信息
	//   - error: 查询失败或成员关系不存在时的错误信息
	GetMembership(ctx context.Context, membershipID string) (*identity_srv.UserMembership, error)

	// GetUserMemberships 获取符合条件的成员关系列表
	//
	// 该方法支持多种查询条件：
	// - 按用户ID查询：获取用户的所有成员关系
	// - 按组织ID查询：获取组织的所有成员
	// - 按部门ID查询：获取部门的所有成员
	// - 按状态过滤：只返回特定状态的成员关系
	// - 分页查询：支持分页参数，优化大数据量场景的性能
	//
	// 对应IDL: GetUserMemberships(1: GetUserMembershipsRequest req)
	//
	// 参数:
	//   - ctx: 上下文信息
	//   - req: 查询请求，包含过滤条件和分页参数
	//
	// 返回:
	//   - *identity_srv.GetUserMembershipsResponse: 包含成员关系列表和分页信息的响应
	//   - error: 查询失败时的错误信息
	GetUserMemberships(
		ctx context.Context,
		req *identity_srv.GetUserMembershipsRequest,
	) (*identity_srv.GetUserMembershipsResponse, error)

	// GetPrimaryMembership 获取用户的主要成员关系
	//
	// 每个用户可以有多个成员关系，但只能有一个主要关系。
	// 主要关系通常用于确定用户的默认组织归属和角色权限。
	//
	// 对应IDL: GetPrimaryMembership(1: core.ULID userID)
	//
	// 参数:
	//   - ctx: 上下文信息
	//   - userID: 用户的唯一标识符
	//
	// 返回:
	//   - *identity_srv.UserMembership: 用户的主要成员关系信息
	//   - error: 查询失败或用户没有主要关系时的错误信息
	GetPrimaryMembership(ctx context.Context, userID string) (*identity_srv.UserMembership, error)

	// CheckMembership 检查用户是否属于某个组织或部门
	//
	// 该方法用于权限验证和访问控制，快速判断用户是否有权访问特定资源。
	// 检查范围：
	// - 如果指定部门ID，则检查用户是否是该部门的成员
	// - 如果只指定组织ID，则检查用户是否是该组织的成员（任意部门）
	// - 只返回活跃状态的成员关系
	//
	// 对应IDL: CheckMembership(1: CheckMembershipRequest req)
	//
	// 参数:
	//   - ctx: 上下文信息
	//   - req: 检查请求，包含用户ID、组织ID和可选的部门ID
	//
	// 返回:
	//   - bool: true表示用户是该组织/部门的活跃成员，false表示不是
	//   - error: 检查过程中发生的错误
	CheckMembership(ctx context.Context, req *identity_srv.CheckMembershipRequest) (bool, error)
}

// Logic 接口设计原则：
//
// 1. 接口契约一致性：所有方法与 IDL 定义严格对应，确保 RPC 服务的实现一致性
// 2. 单一职责原则：接口专注于成员关系管理，不包含用户或组织的具体业务逻辑
// 3. 依赖反转原则：接口不依赖具体的数据访问实现，便于测试和扩展
// 4. 错误处理规范：使用项目统一的错误码和错误处理机制
// 5. 事务安全：涉及数据变更的操作都在事务中执行，保证数据一致性
// 6. 性能考虑：查询方法支持分页和条件过滤，避免大数据量查询的性能问题
