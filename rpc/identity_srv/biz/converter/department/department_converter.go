package department

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// DepartmentConverter 部门转换器接口，负责部门模型的转换
// 基于重构后的 IDL 设计，支持以下实体的双向转换：
// - Department: 部门实体
type Converter interface {
	// ============================================================================
	// Department 部门转换
	// ============================================================================

	// 部门：models -> Thrift
	ModelDepartmentToThrift(*models.Department) *identity_srv.Department
	ModelDepartmentsToThrift([]*models.Department) []*identity_srv.Department

	// 新增的Logic层需要的方法
	ModelToThrift(dept *models.Department) *identity_srv.Department
	CreateRequestToModel(req *identity_srv.CreateDepartmentRequest) *models.Department
	ApplyUpdateToModel(
		existing *models.Department,
		req *identity_srv.UpdateDepartmentRequest,
	) *models.Department
}
