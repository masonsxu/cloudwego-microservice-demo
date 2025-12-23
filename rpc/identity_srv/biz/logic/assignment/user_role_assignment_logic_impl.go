package assignment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	assignmentDal "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/assignment"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
	tracelog "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/log"
)

// LogicImpl 用户角色分配业务逻辑实现
type LogicImpl struct {
	dal       dal.DAL
	converter converter.Converter
}

// NewUserRoleAssignmentLogic 创建用户角色分配业务逻辑实例
func NewLogic(
	dal dal.DAL,
	converter converter.Converter,
) RoleAssignmentLogic {
	return &LogicImpl{
		dal:       dal,
		converter: converter,
	}
}

// AssignRoleToUser 为用户分配一个角色
func (l *LogicImpl) AssignRoleToUser(
	ctx context.Context,
	req *identity_srv.AssignRoleToUserRequest,
) (*identity_srv.UserRoleAssignmentResponse, error) {
	userID := *req.UserID
	roleID := *req.RoleID
	assignedByID := *req.AssignedBy

	// 检查用户是否已分配该角色，避免重复分配
	exists, err := l.dal.UserRoleAssignment().CheckUserRoleExists(ctx, userID, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("检查角色分配状态失败: " + err.Error())
	}

	if exists {
		return nil, errno.ErrRoleAssignmentAlreadyExists
	}

	// 验证角色是否存在
	role, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色信息失败: " + err.Error())
	}

	if role == nil {
		return nil, errno.ErrRoleDefinitionNotFound
	}

	// 创建角色分配记录
	assignment := &models.UserRoleAssignment{
		UserID: uuid.MustParse(userID),
		RoleID: uuid.MustParse(roleID),
	}

	if assignedByID != "" {
		assignedByUUID := uuid.MustParse(assignedByID)
		assignment.CreatedBy = &assignedByUUID
	}

	// 保存到数据库
	if err := l.dal.UserRoleAssignment().Create(ctx, assignment); err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("创建角色分配失败: " + err.Error())
	}

	return &identity_srv.UserRoleAssignmentResponse{
		AssignmentID: convutil.StringPtr(assignment.ID.String()),
	}, nil
}

// UpdateUserRoleAssignment 更新用户的角色分配信息
func (l *LogicImpl) UpdateUserRoleAssignment(
	ctx context.Context,
	req *identity_srv.UpdateUserRoleAssignmentRequest,
) error {
	// 参数验证
	if req.AssignmentID == nil || *req.AssignmentID == "" {
		return errno.ErrInvalidParams.WithMessage("分配ID不能为空")
	}

	assignmentID := *req.AssignmentID

	// 查询现有分配记录
	assignment, err := l.dal.UserRoleAssignment().GetByID(ctx, assignmentID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("查询角色分配失败: " + err.Error())
	}

	if assignment == nil {
		return errno.ErrRoleAssignmentNotFound
	}

	// 更新字段
	if req.UserID != nil {
		assignment.UserID = uuid.MustParse(*req.UserID)
	}

	if req.RoleID != nil {
		assignment.RoleID = uuid.MustParse(*req.RoleID)
	}

	if req.UpdatedBy != nil {
		updatedByUUID := uuid.MustParse(*req.UpdatedBy)
		assignment.UpdatedBy = &updatedByUUID
	}

	// 保存更新
	if err := l.dal.UserRoleAssignment().Update(ctx, assignment); err != nil {
		return errno.ErrOperationFailed.WithMessage("更新角色分配失败: " + err.Error())
	}

	return nil
}

// RevokeRoleFromUser 撤销用户的角色分配
func (l *LogicImpl) RevokeRoleFromUser(
	ctx context.Context,
	req *identity_srv.RevokeRoleFromUserRequest,
) error {
	// 参数验证
	if req.UserID == nil || *req.UserID == "" {
		return errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	if req.RoleID == nil || *req.RoleID == "" {
		return errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	userID := *req.UserID
	roleID := *req.RoleID

	// 1. 检查用户是否为系统用户
	isSystemUser, err := l.dal.UserProfile().IsSystemUser(ctx, userID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("检查用户类型失败: " + err.Error())
	}

	// 2. 如果是系统用户，检查要撤销的角色是否为系统角色
	if isSystemUser {
		role, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
		if err != nil {
			return errno.ErrRoleDefinitionNotFound
		}

		if role.IsSystemRole {
			tracelog.Ctx(ctx).Warn().
				Str("user_id", userID).
				Str("role_id", roleID).
				Str("role_name", role.Name).
				Msg("尝试撤销系统用户的系统角色被拒绝")

			return errno.ErrSystemRoleCannotRevoke.WithMessage(
				fmt.Sprintf("系统用户的系统角色 '%s' 不能被撤销", role.Name),
			)
		}
	}

	// 3. 查找用户和角色的分配记录
	assignment, err := l.dal.UserRoleAssignment().FindByUserAndRole(ctx, userID, roleID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("查询角色分配失败: " + err.Error())
	}

	if assignment == nil {
		return errno.ErrRoleAssignmentNotFound.WithMessage("未找到该用户的角色分配记录")
	}

	// 4. 删除角色分配记录
	if err := l.dal.UserRoleAssignment().Delete(ctx, assignment.ID.String()); err != nil {
		return errno.ErrOperationFailed.WithMessage("撤销角色分配失败: " + err.Error())
	}

	// 5. 审计日志
	tracelog.Ctx(ctx).Info().
		Str("user_id", userID).
		Str("role_id", roleID).
		Msg("角色撤销成功")

	return nil
}

// GetLastUserRoleAssignment 获取用户最后一次的角色分配信息
func (l *LogicImpl) GetLastUserRoleAssignment(
	ctx context.Context,
	userID string,
) (*identity_srv.UserRoleAssignment, error) {
	// 参数验证
	if userID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID不能为空")
	}

	// 查询用户最后的角色分配
	assignment, err := l.dal.UserRoleAssignment().GetLastUserRoleAssignment(ctx, userID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询用户角色分配失败: " + err.Error())
	}

	if assignment == nil {
		return nil, errno.ErrRoleAssignmentNotFound
	}

	// 转换为Thrift格式返回
	return l.converter.UserRoleAssignment().ModelToThrift(assignment), nil
}

// ListUserRoleAssignments 列出用户的角色分配记录
func (l *LogicImpl) ListUserRoleAssignments(
	ctx context.Context,
	req *identity_srv.UserRoleQueryRequest,
) (*identity_srv.UserRoleListResponse, error) {
	// 构建查询条件
	conditions := &assignmentDal.UserRoleAssignmentQueryConditions{}

	if req.UserID != nil {
		userID := *req.UserID
		conditions.UserID = &userID
	}

	if req.RoleID != nil {
		roleID := *req.RoleID
		conditions.RoleID = &roleID
	}

	// 查询角色分配记录
	assignments, pageResult, err := l.dal.UserRoleAssignment().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色分配列表失败: " + err.Error())
	}

	// 转换为Thrift格式
	thriftAssignments := make([]*identity_srv.UserRoleAssignment, 0, len(assignments))
	for _, assignment := range assignments {
		thriftAssignments = append(
			thriftAssignments,
			l.converter.UserRoleAssignment().ModelToThrift(assignment),
		)
	}

	return &identity_srv.UserRoleListResponse{
		Assignments: thriftAssignments,
		Page:        l.converter.Base().PageResponseToThrift(pageResult),
	}, nil
}

// GetUsersByRole 根据角色ID获取该角色下所有用户
func (l *LogicImpl) GetUsersByRole(
	ctx context.Context,
	req *identity_srv.GetUsersByRoleRequest,
) (*identity_srv.GetUsersByRoleResponse, error) {
	// 参数验证
	if req.RoleID == nil || *req.RoleID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	roleID := *req.RoleID

	// 验证角色是否存在
	role, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色信息失败: " + err.Error())
	}

	if role == nil {
		return nil, errno.ErrRoleDefinitionNotFound
	}

	// 获取该角色下所有用户ID
	userIDs, err := l.dal.UserRoleAssignment().GetAllUserIDsByRoleID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色用户列表失败: " + err.Error())
	}

	return &identity_srv.GetUsersByRoleResponse{
		RoleID:  req.RoleID,
		UserIDs: userIDs,
	}, nil
}

// BatchBindUsersToRole 批量绑定用户到角色
func (l *LogicImpl) BatchBindUsersToRole(
	ctx context.Context,
	req *identity_srv.BatchBindUsersToRoleRequest,
) (*identity_srv.BatchBindUsersToRoleResponse, error) {
	// 参数验证
	if req.RoleID == nil || *req.RoleID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	if req.UserIDs == nil {
		return nil, errno.ErrInvalidParams.WithMessage("用户ID列表不能为空")
	}

	roleID := *req.RoleID
	userIDs := req.UserIDs

	operatorID := ""
	if req.OperatorID != nil {
		operatorID = *req.OperatorID
	}

	// 验证角色是否存在
	role, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色信息失败: " + err.Error())
	}

	if role == nil {
		return nil, errno.ErrRoleDefinitionNotFound
	}

	// 批量替换角色的用户绑定（事务操作）
	err = l.dal.UserRoleAssignment().ReplaceRoleUsers(ctx, roleID, userIDs, operatorID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("批量绑定用户到角色失败: " + err.Error())
	}

	successCount := int32(len(userIDs))
	message := "批量绑定成功"

	return &identity_srv.BatchBindUsersToRoleResponse{
		Success:      convutil.BoolPtr(true),
		SuccessCount: &successCount,
		Message:      &message,
	}, nil
}

// BatchGetUserRoles 批量获取多个用户的角色分配
func (l *LogicImpl) BatchGetUserRoles(
	ctx context.Context,
	req *identity_srv.BatchGetUserRolesRequest,
) (*identity_srv.BatchGetUserRolesResponse, error) { // 批量查询
	rolesMap, err := l.dal.UserRoleAssignment().GetRolesByUserIDs(ctx, req.UserIDs)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("批量查询用户角色失败: " + err.Error())
	}

	// 转换为响应格式
	userRoles := make([]*identity_srv.UserRoles, 0, len(rolesMap))
	for userID, roleIDs := range rolesMap {
		userRoles = append(userRoles, &identity_srv.UserRoles{
			UserID:  convutil.StringPtr(userID),
			RoleIDs: roleIDs,
		})
	}

	return &identity_srv.BatchGetUserRolesResponse{
		UserRoles: userRoles,
	}, nil
}
