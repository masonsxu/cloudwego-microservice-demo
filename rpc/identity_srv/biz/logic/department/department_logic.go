package department

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// DepartmentLogic 部门管理业务逻辑接口
// 负责机构内部门的创建、管理和成员关系维护
type DepartmentLogic interface {
	// ============================================================================
	// 部门基础操作
	// ============================================================================

	// CreateDepartment 创建部门
	CreateDepartment(
		ctx context.Context,
		req *identity_srv.CreateDepartmentRequest,
	) (*identity_srv.Department, error)

	// GetDepartment 根据ID获取部门信息
	GetDepartment(
		ctx context.Context,
		req *identity_srv.GetDepartmentRequest,
	) (*identity_srv.Department, error)

	// UpdateDepartment 更新部门信息
	UpdateDepartment(
		ctx context.Context,
		req *identity_srv.UpdateDepartmentRequest,
	) (*identity_srv.Department, error)

	// DeleteDepartment 删除部门（软删除）
	DeleteDepartment(ctx context.Context, departmentID string) error

	// ============================================================================
	// 部门查询操作
	// ============================================================================

	// GetDepartmentsByOrganization 获取组织下的所有部门
	GetDepartmentsByOrganization(
		ctx context.Context,
		req *identity_srv.GetOrganizationDepartmentsRequest,
	) (*identity_srv.GetOrganizationDepartmentsResponse, error)
}
