package department

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// DepartmentRepository 部门仓储接口
// 基于 models.Department 和 IDL 设计，管理组织内的部门结构
type DepartmentRepository interface {
	// 嵌入基础仓储接口
	base.BaseRepository[models.Department]

	// ============================================================================
	// 部门查询
	// ============================================================================

	// GetByName 根据名称获取部门（在指定组织内）
	GetByName(ctx context.Context, name string, organizationID string) (*models.Department, error)

	// ExistsByID 检查部门是否存在
	ExistsByID(ctx context.Context, departmentID string) (bool, error)

	// ============================================================================
	// 设备管理相关
	// ============================================================================

	// GetDepartmentEquipment 获取部门的所有可用设备ID列表
	GetDepartmentEquipment(ctx context.Context, departmentID string) ([]string, error)

	// AddEquipment 为部门添加设备
	AddEquipment(ctx context.Context, departmentID string, equipmentIDs []string) error

	// RemoveEquipment 从部门移除设备
	RemoveEquipment(ctx context.Context, departmentID string, equipmentIDs []string) error

	// UpdateEquipment 更新部门的设备列表
	UpdateEquipment(ctx context.Context, departmentID string, equipmentIDs []string) error

	// ============================================================================
	// 成员关系查询
	// ============================================================================

	// GetDepartmentMembers 获取部门成员ID列表
	GetDepartmentMembers(ctx context.Context, departmentID string) ([]string, error)

	// HasMembers 检查部门是否有成员
	HasMembers(ctx context.Context, departmentID string) (bool, error)

	// CountMembers 统计部门成员数量
	CountMembers(ctx context.Context, departmentID string) (int64, error)

	// ============================================================================
	// 数据完整性检查
	// ============================================================================

	// CheckNameExists 检查部门名称是否已存在（在指定组织内）
	CheckNameExists(
		ctx context.Context,
		name string,
		organizationID string,
		excludeID ...string,
	) (bool, error)

	// ValidateOrganizationExists 验证组织是否存在
	ValidateOrganizationExists(ctx context.Context, organizationID string) (bool, error)

	// ============================================================================
	// 批量操作
	// ============================================================================

	// BatchCreateDepartments 批量创建部门
	BatchCreateDepartments(ctx context.Context, departments []*models.Department) error

	// BatchUpdateOrganization 批量更新部门的组织归属
	BatchUpdateOrganization(
		ctx context.Context,
		departmentIDs []string,
		newOrganizationID string,
	) error

	// BatchDeleteByOrganization 批量删除指定组织的所有部门
	BatchDeleteByOrganization(ctx context.Context, organizationID string) error

	// ============================================================================
	// 统计分析
	// ============================================================================

	// CountByOrganization 统计指定组织的部门数量
	CountByOrganization(ctx context.Context, organizationID string) (int64, error)

	// CountByDepartmentType 统计指定类型的部门数量
	CountByDepartmentType(ctx context.Context, departmentType string) (int64, error)

	// GetDepartmentStatistics 获取部门统计信息
	GetDepartmentStatistics(ctx context.Context, departmentID string) (*DepartmentStatistics, error)

	// GetOrganizationDepartmentStatistics 获取组织的部门统计信息
	GetOrganizationDepartmentStatistics(
		ctx context.Context,
		organizationID string,
	) (*OrganizationDepartmentStatistics, error)

	// ============================================================================
	// 统一查询方法（新增）
	// ============================================================================

	// FindWithConditions 根据组合查询条件查询部门列表
	FindWithConditions(
		ctx context.Context,
		conditions *DepartmentQueryConditions,
	) ([]*models.Department, *models.PageResult, error)
}

// DepartmentQueryConditions 部门查询条件
// 支持多条件组合查询，提供灵活的查询能力
type DepartmentQueryConditions struct {
	Name           *string // 部门名称（模糊匹配）
	OrganizationID *string // 组织ID
	DepartmentType *string // 部门类型
	EquipmentID    *string // 设备ID（需JOIN department_equipment表）
	Page           *base.QueryOptions
}

// DepartmentStatistics 部门统计信息
type DepartmentStatistics struct {
	MembersCount   int64 `json:"members_count"`   // 成员数量
	EquipmentCount int64 `json:"equipment_count"` // 设备数量
	ActiveMembers  int64 `json:"active_members"`  // 活跃成员数量
}

// OrganizationDepartmentStatistics 组织部门统计信息
type OrganizationDepartmentStatistics struct {
	TotalDepartments int64            `json:"total_departments"` // 总部门数
	TypeDistribution map[string]int64 `json:"type_distribution"` // 部门类型分布
	TotalMembers     int64            `json:"total_members"`     // 总成员数
	TotalEquipment   int64            `json:"total_equipment"`   // 总设备数
	AverageMembers   float64          `json:"average_members"`   // 平均成员数
	AverageEquipment float64          `json:"average_equipment"` // 平均设备数
}
