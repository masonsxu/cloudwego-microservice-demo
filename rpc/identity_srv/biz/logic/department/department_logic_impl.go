package department

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal"
	departmentDAL "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/department"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/pkg/errno"
)

// LogicImpl 部门管理业务逻辑实现
type LogicImpl struct {
	dal       dal.DAL
	converter converter.Converter
}

// NewLogic 创建部门管理业务逻辑实例
func NewLogic(dal dal.DAL, converter converter.Converter) DepartmentLogic {
	return &LogicImpl{
		dal:       dal,
		converter: converter,
	}
}

// ============================================================================
// 部门基础操作
// ============================================================================

// CreateDepartment 创建部门
func (l *LogicImpl) CreateDepartment(
	ctx context.Context,
	req *identity_srv.CreateDepartmentRequest,
) (*identity_srv.Department, error) {
	// 参数验证
	if err := l.validateCreateDepartmentRequest(req); err != nil {
		return nil, err
	}

	// 检查组织是否存在
	orgExists, err := l.dal.Organization().ExistsByID(ctx, *req.OrganizationID)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("检查组织是否存在失败: " + err.Error())
	}

	if !orgExists {
		return nil, errno.ErrOrganizationNotFound
	}

	// 转换请求为模型
	dept := l.converter.Department().CreateRequestToModel(req)

	// 在事务中创建部门
	var result *models.Department

	txErr := l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		if err := txDAL.Department().Create(ctx, dept); err != nil {
			return err
		}

		result = dept

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return l.converter.Department().ModelToThrift(result), nil
}

// GetDepartment 根据ID获取部门信息
func (l *LogicImpl) GetDepartment(
	ctx context.Context,
	req *identity_srv.GetDepartmentRequest,
) (*identity_srv.Department, error) {
	if req.DepartmentID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("部门ID不能为空")
	}

	dept, err := l.dal.Department().GetByID(ctx, *req.DepartmentID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrDepartmentNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取部门信息失败: " + err.Error())
	}

	return l.converter.Department().ModelToThrift(dept), nil
}

// UpdateDepartment 更新部门信息
func (l *LogicImpl) UpdateDepartment(
	ctx context.Context,
	req *identity_srv.UpdateDepartmentRequest,
) (*identity_srv.Department, error) {
	if req.DepartmentID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("部门ID不能为空")
	}

	// 获取现有部门
	existingDept, err := l.dal.Department().GetByID(ctx, *req.DepartmentID)
	if err != nil {
		if errno.IsRecordNotFound(err) {
			return nil, errno.ErrDepartmentNotFound
		}

		return nil, errno.ErrOperationFailed.WithMessage("获取部门信息失败: " + err.Error())
	}

	// 应用更新
	updatedDept := l.converter.Department().ApplyUpdateToModel(existingDept, req)

	// 在事务中更新
	var result *models.Department

	txErr := l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		if err := txDAL.Department().Update(ctx, updatedDept); err != nil {
			return err
		}

		result = updatedDept

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return l.converter.Department().ModelToThrift(result), nil
}

// DeleteDepartment 删除部门（软删除）
func (l *LogicImpl) DeleteDepartment(
	ctx context.Context,
	departmentID string,
) error {
	if departmentID == "" {
		return errno.ErrInvalidParams.WithMessage("部门ID不能为空")
	}

	// 检查部门是否存在
	exists, err := l.dal.Department().ExistsByID(ctx, departmentID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("检查部门是否存在失败: " + err.Error())
	}

	if !exists {
		return errno.ErrDepartmentNotFound
	}

	// 检查是否有成员
	memberCount, err := l.dal.UserMembership().CountByDepartmentID(ctx, departmentID)
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("检查部门成员失败: " + err.Error())
	}

	if memberCount > 0 {
		return errno.ErrCannotDeleteDepartmentWithMembers
	}

	// 软删除部门
	err = l.dal.WithTransaction(ctx, func(ctx context.Context, txDAL dal.DAL) error {
		return txDAL.Department().SoftDelete(ctx, departmentID)
	})
	if err != nil {
		return errno.ErrOperationFailed.WithMessage("删除部门失败: " + err.Error())
	}

	return nil
}

// GetDepartmentsByOrganization 获取组织下的所有部门
func (l *LogicImpl) GetDepartmentsByOrganization(
	ctx context.Context,
	req *identity_srv.GetOrganizationDepartmentsRequest,
) (*identity_srv.GetOrganizationDepartmentsResponse, error) {
	if req.OrganizationID == nil {
		return nil, errno.ErrInvalidParams.WithMessage("组织ID不能为空")
	}

	// 使用 Base Converter 转换分页参数
	opts := l.converter.Base().PageRequestToQueryOptions(req.Page)

	// 构建查询条件
	conditions := &departmentDAL.DepartmentQueryConditions{
		Page:           opts,
		OrganizationID: req.OrganizationID,
	}

	// 使用 FindWithConditions 查询
	departments, pageResult, err := l.dal.Department().FindWithConditions(ctx, conditions)
	if err != nil {
		return nil, errno.ErrOperationFailed.WithMessage("获取组织部门失败: " + err.Error())
	}

	// 转换结果
	departmentList := l.converter.Department().ModelDepartmentsToThrift(departments)

	return &identity_srv.GetOrganizationDepartmentsResponse{
		Departments: departmentList,
		Page:        l.converter.Base().PageResponseToThrift(pageResult),
	}, nil
}

// ============================================================================
// 私有辅助方法
// ============================================================================

// validateCreateDepartmentRequest 验证创建部门请求
func (l *LogicImpl) validateCreateDepartmentRequest(
	req *identity_srv.CreateDepartmentRequest,
) error {
	if req.Name == nil {
		return errno.ErrInvalidParams.WithMessage("部门名称不能为空")
	}

	if req.OrganizationID == nil {
		return errno.ErrInvalidParams.WithMessage("组织ID不能为空")
	}

	return nil
}
