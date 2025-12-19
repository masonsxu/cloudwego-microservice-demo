package definition

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	definitionDal "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/definition"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// LogicImpl 角色定义业务逻辑实现
type LogicImpl struct {
	dal       dal.DAL
	converter converter.Converter
}

// NewLogic 创建角色定义业务逻辑实例
func NewLogic(dal dal.DAL, converter converter.Converter) RoleDefinitionLogic {
	return &LogicImpl{
		dal:       dal,
		converter: converter,
	}
}

// CreateRoleDefinition 创建一个新的角色定义
func (l *LogicImpl) CreateRoleDefinition(
	ctx context.Context,
	req *identity_srv.RoleDefinitionCreateRequest,
) (*identity_srv.RoleDefinition, error) {
	// 检查角色名称是否已存在
	exists, err := l.dal.RoleDefinition().CheckNameExists(ctx, *req.Name)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("检查角色名称是否存在失败: " + err.Error())
	}

	if exists {
		return nil, errno.ErrRoleNameAlreadyExists
	}

	// 转换权限列表
	identitys := make(models.Permissions, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		description := ""
		if p.Description != nil {
			description = *p.Description
		}

		identitys = append(identitys, &models.Permission{
			Resource:    *p.Resource,
			Action:      *p.Action,
			Description: description,
		})
	}

	// 创建角色定义记录
	roleDefinition := &models.RoleDefinition{
		Name:         *req.Name,
		Description:  *req.Description,
		Status:       models.RoleStatusInactive, // 默认未激活状态
		Permissions:  identitys,
		IsSystemRole: req.IsSystemRole,
	}

	// 保存到数据库
	if err := l.dal.RoleDefinition().Create(ctx, roleDefinition); err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("创建角色定义失败: " + err.Error())
	}

	// 转换为Thrift格式返回
	return l.converter.RoleDefinition().ModelToThrift(roleDefinition), nil
}

// UpdateRoleDefinition 更新一个已有的角色定义
func (l *LogicImpl) UpdateRoleDefinition(
	ctx context.Context,
	req *identity_srv.RoleDefinitionUpdateRequest,
) (*identity_srv.RoleDefinition, error) {
	// 参数验证
	if req.RoleDefinitionID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("角色定义ID不能为空")
	}

	roleID := *req.RoleDefinitionID

	// 查询现有角色定义
	role, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色定义失败: " + err.Error())
	}

	if role == nil {
		return nil, errno.ErrRoleDefinitionNotFound
	}

	// 检查是否为系统角色，系统角色不允许修改
	if role.IsSystemRole {
		return nil, errno.ErrSystemRoleCannotModify
	}

	// 更新字段
	if req.Description != nil {
		role.Description = *req.Description
	}

	if req.Status != nil {
		role.Status = models.RoleStatus(*req.Status)
	}

	// 更新角色名称（如果提供）
	if req.Name != nil {
		// 检查新名称是否与现有名称重复（排除自身）
		existingRole, err := l.dal.RoleDefinition().FindByName(ctx, *req.Name)
		if err != nil {
			return nil, errno.ErrOperationFailed.WithMessage("检查角色名称失败: " + err.Error())
		}

		// 如果找到同名角色且不是当前角色，则返回错误
		if existingRole != nil && existingRole.ID.String() != roleID {
			return nil, errno.ErrInvalidParams.WithMessage("角色名称已存在")
		}

		role.Name = *req.Name
	}

	if req.Permissions != nil {
		identitys := make(models.Permissions, 0, len(req.Permissions))
		for _, p := range req.Permissions {
			description := ""
			if p.Description != nil {
				description = *p.Description
			}

			identitys = append(identitys, &models.Permission{
				Resource:    *p.Resource,
				Action:      *p.Action,
				Description: description,
			})
		}

		role.Permissions = identitys
	}

	// 保存更新
	if err := l.dal.RoleDefinition().Update(ctx, role); err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("更新角色定义失败: " + err.Error())
	}

	// 转换为Thrift格式返回
	return l.converter.RoleDefinition().ModelToThrift(role), nil
}

// DeleteRoleDefinition 删除一个角色定义
func (l *LogicImpl) DeleteRoleDefinition(ctx context.Context, roleID string) error {
	// 参数验证
	if roleID == "" {
		return errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	// 查询角色定义
	role, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("查询角色定义失败: " + err.Error())
	}

	if role == nil {
		return errno.ErrRoleDefinitionNotFound
	}

	// 检查是否为系统角色，系统角色不允许删除
	if role.IsSystemRole {
		return errno.ErrSystemRoleCannotDelete
	}

	// 检查是否有用户正在使用该角色
	count, err := l.dal.UserRoleAssignment().CountByRoleID(ctx, roleID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("检查角色使用情况失败: " + err.Error())
	}

	if count > 0 {
		return errno.ErrRoleInUseCannotDelete
	}

	// 删除角色定义
	if err := l.dal.RoleDefinition().Delete(ctx, roleID); err != nil {
		return errno.ErrOperationFailed.WithMessage("删除角色定义失败: " + err.Error())
	}

	return nil
}

// GetRoleDefinition 根据ID获取角色定义
func (l *LogicImpl) GetRoleDefinition(
	ctx context.Context,
	roleID string,
) (*identity_srv.RoleDefinition, error) {
	// 参数验证
	if roleID == "" {
		return nil, errno.ErrInvalidParams.WithMessage("角色ID不能为空")
	}

	// 查询角色定义
	role, err := l.dal.RoleDefinition().GetByID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色定义失败: " + err.Error())
	}

	if role == nil {
		return nil, errno.ErrRoleDefinitionNotFound
	}

	// 查询该角色的用户数量
	userCount, err := l.dal.UserRoleAssignment().CountByRoleID(ctx, roleID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色用户数量失败: " + err.Error())
	}

	// 设置用户数量到模型（非持久化字段）
	role.UserCount = userCount

	// 转换为Thrift格式返回
	return l.converter.RoleDefinition().ModelToThrift(role), nil
}

// ListRoleDefinitions 分页列出角色定义
func (l *LogicImpl) ListRoleDefinitions(
	ctx context.Context,
	req *identity_srv.RoleDefinitionQueryRequest,
) (*identity_srv.RoleDefinitionListResponse, error) {
	// 使用 Base Converter 转换分页参数
	opts := l.converter.Base().PageRequestToQueryOptions(req.Page)

	// 构建查询条件
	conditions := &definitionDal.RoleDefinitionQueryConditions{
		Page: opts, // 设置分页参数
	}

	// 设置业务过滤条件
	if req != nil {
		if req.Name != nil {
			conditions.Name = req.Name
		}

		if req.Status != nil {
			status := models.RoleStatus(*req.Status)
			conditions.Status = &status
		}

		if req.IsSystemRole != nil {
			conditions.IsSystemRole = req.IsSystemRole
		}
	}

	// 查询角色定义列表
	roles, pageResult, err := l.dal.RoleDefinition().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("查询角色定义列表失败: " + err.Error())
	}

	// 批量查询所有角色的用户数量（避免 N+1 问题）
	// 收集所有角色ID
	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID.String())
	}

	// 创建角色ID到用户数量的映射
	roleUserCountMap := make(map[string]int64)

	if len(roleIDs) > 0 {
		// 遍历每个角色ID，查询用户数量
		for _, roleID := range roleIDs {
			count, err := l.dal.UserRoleAssignment().CountByRoleID(ctx, roleID)
			if err != nil {
				// 查询失败时记录错误但继续处理，设置为0
				roleUserCountMap[roleID] = 0
			} else {
				roleUserCountMap[roleID] = count
			}
		}
	}

	// 转换为Thrift格式，并设置用户数量
	thriftRoles := make([]*identity_srv.RoleDefinition, 0, len(roles))
	for _, role := range roles {
		// 设置用户数量到模型（非持久化字段）
		if userCount, exists := roleUserCountMap[role.ID.String()]; exists {
			role.UserCount = userCount
		}

		thriftRoles = append(thriftRoles, l.converter.RoleDefinition().ModelToThrift(role))
	}

	return &identity_srv.RoleDefinitionListResponse{
		Roles: thriftRoles,
		Page:  l.converter.Base().PageResponseToThrift(pageResult),
	}, nil
}
