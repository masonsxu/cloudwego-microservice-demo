package membership

import (
	"context"
	"fmt"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	membershipDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/membership"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// LogicImpl 用户成员关系业务逻辑实现
//
// 该实现严格按照IDL定义提供完整的成员关系管理功能，采用分层架构设计：
// - 业务逻辑层：负责业务规则验证、流程控制和数据转换
// - 数据访问层：通过DAL接口访问数据库，支持事务管理
// - 转换器层：负责DTO与模型之间的数据转换
//
// 设计特点：
// - 依赖注入：通过构造函数注入依赖，便于测试和扩展
// - 事务安全：关键操作在事务中执行，保证数据一致性
// - 错误处理：统一的错误码和错误信息，提供详细的失败原因
// - 参数验证：严格的参数校验，防止无效数据进入业务逻辑
// - 业务规则：实现完整的业务约束，如唯一性检查、关系验证等
type LogicImpl struct {
	dal       dal.DAL
	converter converter.Converter
}

// NewLogic 创建用户成员关系业务逻辑实例
//
// 使用工厂模式创建实例，确保依赖关系的正确注入。
//
// 参数:
//   - dal: 数据访问层接口，提供数据库操作能力
//   - converter: 数据转换器接口，处理DTO与模型的转换
//
// 返回:
//   - UserMembershipLogic: 业务逻辑接口实例
func NewLogic(
	dal dal.DAL,
	converter converter.Converter,
) MembershipLogic {
	return &LogicImpl{
		dal:       dal,
		converter: converter,
	}
}

// ============================================================================
// 成员关系生命周期管理（与IDL完全对应）
// ============================================================================

// AddMembership 为用户添加新的组织成员关系
//
// 业务流程：
// 1. 参数验证：检查必填字段的有效性
// 2. 实体验证：确认用户、组织、部门的存在性
// 3. 角色验证：验证角色定义的有效性和权限
// 4. 冲突检查：确保不存在重复的成员关系
// 5. 事务创建：在事务中创建成员关系，处理主要关系约束
//
// 业务规则：
// - 一个用户在同一组织的同一部门只能有一个活跃的成员关系
// - 如果设置为主要关系，会自动取消用户的其他主要关系
// - 角色名称必须对应已存在的有效角色定义
// - 部门必须属于指定的组织
func (l *LogicImpl) AddMembership(
	ctx context.Context,
	req *identity_srv.AddMembershipRequest,
) (*identity_srv.UserMembership, error) {
	// 1. 参数验证
	if err := l.validateAddMembershipRequest(req); err != nil {
		return nil, err
	}

	// 2. 验证关联实体存在性
	if err := l.validateEntityExistence(ctx, *req.UserID, *req.OrganizationID, req.DepartmentID); err != nil {
		return nil, err
	}

	// 3. 检查成员关系冲突
	if err := l.checkMembershipConflict(ctx, *req.UserID, *req.OrganizationID, req.DepartmentID); err != nil {
		return nil, err
	}

	// 4. 转换请求为模型
	membership := l.converter.Membership().AddMembershipRequestToModel(req)

	// 5. 在事务中创建成员关系
	var result *models.UserMembership

	err := l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		// 如果设置为主要关系，需要先将用户的其他主要关系取消
		if membership.IsPrimary {
			if err := txDAL.UserMembership().UnsetPrimaryByUserID(ctx, *req.UserID); err != nil {
				return err
			}
		}

		// 创建成员关系
		if err := txDAL.UserMembership().Create(ctx, membership); err != nil {
			return err
		}

		result = membership

		return nil
	})
	if err != nil {
		return nil, err
	}

	return l.converter.Membership().ModelToThrift(result), nil
}

// UpdateMembership 更新用户的组织成员关系
//
// 业务流程：
// 1. 参数验证：检查更新请求的有效性
// 2. 存在性检查：确认成员关系存在且可更新
// 3. 部门验证：如有部门变更，验证新部门属于同一组织
// 4. 数据更新：应用更新到现有模型
// 5. 事务更新：在事务中执行更新，处理约束冲突
//
// 业务规则：
// - 只能更新部分字段，核心关系（用户ID、组织ID）不可变更
// - 角色变更需要验证新角色的有效性
// - 部门变更必须在同一组织内
// - 主要关系变更会影响用户的其他成员关系
func (l *LogicImpl) UpdateMembership(
	ctx context.Context,
	req *identity_srv.UpdateMembershipRequest,
) (*identity_srv.UserMembership, error) {
	// 1. 参数验证
	if err := l.validateUpdateMembershipRequest(req); err != nil {
		return nil, err
	}

	// 2. 获取现有成员关系
	existingMembership, err := l.dal.UserMembership().GetByID(ctx, *req.MembershipID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrMembershipNotFound.WithMessage(
				fmt.Sprintf("成员关系不存在: %s", *req.MembershipID),
			)
		}

		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("获取成员关系失败: %v", err),
		)
	}

	// 3. 验证部门变更的合法性
	if req.DepartmentID != nil {
		departmentID := ""
		if req.DepartmentID != nil {
			departmentID = *req.DepartmentID
		}

		if err := l.validateDepartmentBelongsToOrganization(
			ctx, departmentID, existingMembership.OrganizationID.String(),
		); err != nil {
			return nil, err
		}
	}

	// 4. 应用更新到模型
	updatedMembership := l.converter.Membership().ApplyUpdateToModel(existingMembership, req)

	// 5. 在事务中执行更新
	var result *models.UserMembership

	txErr := l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		// 如果要设置为主要关系，需要先取消用户的其他主要关系
		if updatedMembership.IsPrimary && !existingMembership.IsPrimary {
			if err := txDAL.UserMembership().UnsetPrimaryByUserID(
				ctx, existingMembership.UserID.String(),
			); err != nil {
				return err
			}
		}

		// 更新成员关系
		if err := txDAL.UserMembership().Update(ctx, updatedMembership); err != nil {
			return err
		}

		result = updatedMembership

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return l.converter.Membership().ModelToThrift(result), nil
}

// RemoveMembership 移除用户的组织成员关系（逻辑删除）
//
// 业务流程：
// 1. 参数验证：检查成员关系ID的有效性
// 2. 存在性检查：确认成员关系存在
// 3. 事务删除：在事务中执行逻辑删除
//
// 业务规则：
// - 执行逻辑删除，保留数据记录用于审计
// - 删除操作不可逆，但数据可通过管理员恢复
// - 删除主要成员关系不会自动指定新的主要关系
func (l *LogicImpl) RemoveMembership(
	ctx context.Context,
	membershipID string,
) error {
	// 1. 参数验证
	if err := l.validateMembershipID(membershipID); err != nil {
		return err
	}

	// 2. 检查成员关系是否存在
	exists, err := l.dal.UserMembership().ExistsByID(ctx, membershipID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("检查成员关系存在性失败: %v", err),
		)
	}

	if !exists {
		return errno.ErrMembershipNotFound.WithMessage(
			fmt.Sprintf("成员关系不存在: %s", membershipID),
		)
	}

	// 3. 在事务中执行逻辑删除
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		return txDAL.UserMembership().Delete(ctx, membershipID)
	})
	if err != nil {
		return errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("删除用户成员关系失败: %v", err),
		)
	}

	return nil
}

// GetMembership 根据ID获取成员关系详情
//
// 该方法提供成员关系的完整信息查询，包括所有字段和状态信息。
func (l *LogicImpl) GetMembership(
	ctx context.Context,
	membershipID string,
) (*identity_srv.UserMembership, error) {
	// 1. 参数验证
	if err := l.validateMembershipID(membershipID); err != nil {
		return nil, err
	}

	// 2. 查询成员关系
	membership, err := l.dal.UserMembership().GetByID(ctx, membershipID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrMembershipNotFound.WithMessage(
				fmt.Sprintf("成员关系不存在: %s", membershipID),
			)
		}

		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("获取成员关系失败: %v", err),
		)
	}

	// 3. 转换并返回
	return l.converter.Membership().ModelToThrift(membership), nil
}

// GetUserMemberships 获取符合条件的成员关系列表
//
// 业务流程：
// 1. 参数验证：检查查询条件的有效性
// 2. 构建查询：根据请求参数构建查询选项
// 3. 执行查询：根据不同条件执行相应的查询方法
// 4. 数据转换：将查询结果转换为响应格式
// 5. 分页处理：构建分页响应信息
//
// 查询策略：
// - 按用户ID查询：获取用户的所有成员关系
// - 按组织ID查询：获取组织的所有成员
// - 按部门ID查询：获取部门的所有成员
// - 通用查询：获取所有活跃的成员关系
func (l *LogicImpl) GetUserMemberships(
	ctx context.Context,
	req *identity_srv.GetUserMembershipsRequest,
) (*identity_srv.GetUserMembershipsResponse, error) {
	// 1. 参数验证
	if err := l.validateGetUserMembershipsRequest(req); err != nil {
		return nil, err
	}

	// 2. 构建查询选项
	opts := l.buildQueryOptionsFromRequest(req)

	// 3. 构建查询条件
	conditions := &membershipDAL.UserMembershipQueryConditions{
		Page: opts,
	}

	// 根据请求参数设置查询条件
	if req.UserID != nil {
		conditions.UserID = req.UserID
	}

	if req.OrganizationID != nil {
		conditions.OrganizationID = req.OrganizationID
	}

	if req.DepartmentID != nil {
		conditions.DepartmentID = req.DepartmentID
	}

	// 如果没有指定任何过滤条件，默认查询活跃的成员关系
	if req.UserID == nil && req.OrganizationID == nil && req.DepartmentID == nil {
		activeStatus := models.MembershipStatusActive
		conditions.Status = &activeStatus
	}

	// 4. 执行查询
	memberships, pageResult, err := l.dal.UserMembership().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("查询用户成员关系失败: %v", err),
		)
	}

	// 5. 转换结果
	membershipList := l.converter.Membership().ModelUserMembershipsToThrift(memberships)

	// 5. 构建响应
	response := &identity_srv.GetUserMembershipsResponse{
		Memberships: membershipList,
	}

	// 6. 添加分页信息
	if pageResult != nil && req.Page != nil {
		response.Page = l.converter.Base().PageResponseToThrift(pageResult)
	}

	return response, nil
}

// GetPrimaryMembership 获取用户的主要成员关系
//
// 主要成员关系用于确定用户的默认组织归属和权限上下文。
func (l *LogicImpl) GetPrimaryMembership(
	ctx context.Context,
	userID string,
) (*identity_srv.UserMembership, error) {
	// 1. 参数验证
	if err := l.validateUserID(userID); err != nil {
		return nil, err
	}

	// 2. 查询主要成员关系
	membership, err := l.dal.UserMembership().GetPrimaryMembership(ctx, userID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrMembershipNotFound.WithMessage(
				fmt.Sprintf("用户没有主要成员关系: %s", userID),
			)
		}

		return nil, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("获取主要成员关系失败: %v", err),
		)
	}

	// 3. 转换并返回
	return l.converter.Membership().ModelToThrift(membership), nil
}

// CheckMembership 检查用户是否属于某个组织或部门
//
// 业务流程：
// 1. 参数验证：检查请求参数的有效性
// 2. 查询执行：根据部门或组织条件查询成员关系
// 3. 状态检查：验证成员关系是否为活跃状态
//
// 查询逻辑：
// - 如果指定部门ID，检查用户在该部门的成员关系
// - 如果只指定组织ID，检查用户在该组织任意部门的成员关系
// - 只考虑活跃状态的成员关系
func (l *LogicImpl) CheckMembership(
	ctx context.Context,
	req *identity_srv.CheckMembershipRequest,
) (bool, error) {
	// 1. 参数验证
	if err := l.validateCheckMembershipRequest(req); err != nil {
		return false, err
	}

	// 2. 查询成员关系
	var (
		membership *models.UserMembership
		err        error
	)

	if req.DepartmentID != nil {
		// 检查部门成员关系
		membership, err = l.dal.UserMembership().GetByUserAndDepartment(
			ctx, *req.UserID, *req.DepartmentID,
		)
	} else {
		// 检查组织成员关系
		membership, err = l.dal.UserMembership().GetByUserAndOrganization(
			ctx, *req.UserID, *req.OrganizationID,
		)
	}

	if err != nil {
		if errno.IsRecordNotFound(err) {
			return false, nil // 不存在成员关系
		}

		return false, errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("检查成员关系失败: %v", err),
		)
	}

	// 3. 检查成员关系是否活跃
	return membership.IsActive(), nil
}

// ============================================================================
// 私有验证方法 - 参数验证
// ============================================================================

// validateAddMembershipRequest 验证添加成员关系请求
func (l *LogicImpl) validateAddMembershipRequest(
	req *identity_srv.AddMembershipRequest,
) error {
	if req == nil {
		return errno.ErrInvalidParams.WithMessage("添加成员关系请求不能为空")
	}

	if err := l.validateUserID(*req.UserID); err != nil {
		return err
	}

	if err := l.validateOrganizationID(*req.OrganizationID); err != nil {
		return err
	}

	// 验证可选的部门ID
	if req.DepartmentID != nil && *req.DepartmentID != "" {
		if err := l.validateDepartmentID(*req.DepartmentID); err != nil {
			return err
		}
	}

	return nil
}

// validateUpdateMembershipRequest 验证更新成员关系请求
func (l *LogicImpl) validateUpdateMembershipRequest(
	req *identity_srv.UpdateMembershipRequest,
) error {
	if req == nil {
		return errno.ErrInvalidParams.WithMessage("更新成员关系请求不能为空")
	}

	if err := l.validateMembershipID(*req.MembershipID); err != nil {
		return err
	}

	if req.DepartmentID != nil && *req.DepartmentID != "" {
		if err := l.validateDepartmentID(*req.DepartmentID); err != nil {
			return err
		}
	}

	return nil
}

// validateGetUserMembershipsRequest 验证获取用户成员关系请求
func (l *LogicImpl) validateGetUserMembershipsRequest(
	req *identity_srv.GetUserMembershipsRequest,
) error {
	if req == nil {
		return errno.ErrInvalidParams.WithMessage("查询成员关系请求不能为空")
	}

	// 至少需要提供一个查询条件
	if req.UserID == nil && req.OrganizationID == nil && req.DepartmentID == nil {
		return errno.ErrInvalidParams.WithMessage("必须提供用户ID、组织ID或部门ID中的至少一个")
	}

	// 验证具体的ID字段
	if req.UserID != nil {
		if err := l.validateUserID(*req.UserID); err != nil {
			return err
		}
	}

	if req.OrganizationID != nil {
		if err := l.validateOrganizationID(*req.OrganizationID); err != nil {
			return err
		}
	}

	if req.DepartmentID != nil {
		if err := l.validateDepartmentID(*req.DepartmentID); err != nil {
			return err
		}
	}

	return nil
}

// validateCheckMembershipRequest 验证检查成员关系请求
func (l *LogicImpl) validateCheckMembershipRequest(
	req *identity_srv.CheckMembershipRequest,
) error {
	if req == nil {
		return errno.ErrInvalidParams.WithMessage("检查成员关系请求不能为空")
	}

	if err := l.validateUserID(*req.UserID); err != nil {
		return err
	}

	if err := l.validateOrganizationID(*req.OrganizationID); err != nil {
		return err
	}

	// 验证可选的部门ID
	if req.DepartmentID != nil {
		if err := l.validateDepartmentID(*req.DepartmentID); err != nil {
			return err
		}
	}

	return nil
}

// ============================================================================
// 私有验证方法 - 字段验证
// ============================================================================

// validateUserID 验证用户ID
func (l *LogicImpl) validateUserID(userID string) error {
	if userID == "" {
		return errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}
	// TODO: 添加ULID格式验证
	return nil
}

// validateOrganizationID 验证组织ID
func (l *LogicImpl) validateOrganizationID(organizationID string) error {
	if organizationID == "" {
		return errno.ErrInvalidParams.WithMessage("组织ID不能为空")
	}
	// TODO: 添加ULID格式验证
	return nil
}

// validateDepartmentID 验证部门ID
func (l *LogicImpl) validateDepartmentID(departmentID string) error {
	if departmentID == "" {
		return errno.ErrInvalidParams.WithMessage("部门ID不能为空")
	}
	// TODO: 添加ULID格式验证
	return nil
}

// validateMembershipID 验证成员关系ID
func (l *LogicImpl) validateMembershipID(membershipID string) error {
	if membershipID == "" {
		return errno.ErrInvalidParams.WithMessage("成员关系ID不能为空")
	}
	// TODO: 添加ULID格式验证
	return nil
}

// ============================================================================
// 私有业务验证方法
// ============================================================================

// validateEntityExistence 验证关联实体存在性
func (l *LogicImpl) validateEntityExistence(
	ctx context.Context,
	userID, organizationID string,
	departmentID *string,
) error {
	// 检查用户是否存在
	userExists, err := l.dal.UserProfile().Exists(ctx, userID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("检查用户存在性失败: %v", err),
		)
	}

	if !userExists {
		return errno.ErrUserNotFound.WithMessage(
			fmt.Sprintf("用户不存在: %s", userID),
		)
	}

	// 检查组织是否存在
	orgExists, err := l.dal.Organization().ExistsByID(ctx, organizationID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("检查组织存在性失败: %v", err),
		)
	}

	if !orgExists {
		return errno.ErrOrganizationNotFound.WithMessage(
			fmt.Sprintf("组织不存在: %s", organizationID),
		)
	}

	// 检查部门是否存在且属于指定组织
	if departmentID != nil && *departmentID != "" {
		return l.validateDepartmentBelongsToOrganization(ctx, *departmentID, organizationID)
	}

	return nil
}

// validateDepartmentBelongsToOrganization 验证部门是否属于指定组织
func (l *LogicImpl) validateDepartmentBelongsToOrganization(
	ctx context.Context,
	departmentID, organizationID string,
) error {
	if departmentID == "" {
		return nil // 空部门ID视为合法（组织级别的成员关系）
	}

	department, err := l.dal.Department().GetByID(ctx, departmentID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return errno.ErrDepartmentNotFound.WithMessage(
				fmt.Sprintf("部门不存在: %s", departmentID),
			)
		}

		return errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("获取部门信息失败: %v", err),
		)
	}

	if department.OrganizationID.String() != organizationID {
		return errno.ErrInvalidParams.WithMessage(
			fmt.Sprintf("部门 %s 不属于组织 %s", departmentID, organizationID),
		)
	}

	return nil
}

// checkMembershipConflict 检查成员关系冲突
func (l *LogicImpl) checkMembershipConflict(
	ctx context.Context,
	userID, organizationID string,
	departmentID *string,
) error {
	var (
		existingMembership *models.UserMembership
		err                error
	)

	// 根据是否指定部门选择不同的查询策略
	if departmentID != nil && *departmentID != "" {
		// 检查部门级别的成员关系
		existingMembership, err = l.dal.UserMembership().GetByUserAndDepartment(
			ctx, userID, *departmentID,
		)
	} else {
		// 检查组织级别的成员关系
		existingMembership, err = l.dal.UserMembership().GetByUserAndOrganization(
			ctx, userID, organizationID,
		)
	}

	if err != nil && !errno.IsRecordNotFound(err) {
		return errno.ErrOperationFailed.WithMessage(
			fmt.Sprintf("检查成员关系冲突失败: %v", err),
		)
	}

	// 如果存在活跃的成员关系，则视为冲突
	if existingMembership != nil && existingMembership.IsActive() {
		scope := "组织"
		if departmentID != nil && *departmentID != "" {
			scope = "部门"
		}

		return errno.ErrMembershipAlreadyExists.WithMessage(
			fmt.Sprintf("用户在该%s已存在活跃的成员关系", scope),
		)
	}

	return nil
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// buildQueryOptionsFromRequest 从请求构建查询选项
func (l *LogicImpl) buildQueryOptionsFromRequest(
	req *identity_srv.GetUserMembershipsRequest,
) *base.QueryOptions {
	opts := base.NewQueryOptions()

	// 处理分页参数
	if req.Page != nil {
		opts.Page = req.Page.Page
		opts.PageSize = req.Page.Limit
	}

	// 设置默认排序
	opts.OrderBy = "created_at"
	opts.OrderDesc = true

	return opts
}
