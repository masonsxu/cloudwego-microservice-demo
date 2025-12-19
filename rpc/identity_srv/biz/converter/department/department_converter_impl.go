package department

import (
	"github.com/google/uuid"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl 部门转换器实现
type ConverterImpl struct{}

// NewConverter 创建一个新的部门转换器
func NewConverter() Converter {
	return &ConverterImpl{}
}

// ============================================================================
// Department 部门转换
// ============================================================================

// ModelDepartmentToThrift 将 models.Department 转换为 identity_srv.Department
func (c *ConverterImpl) ModelDepartmentToThrift(
	model *models.Department,
) *identity_srv.Department {
	if model == nil {
		return nil
	}

	id := model.ID.String()
	name := model.Name
	orgID := model.OrganizationID.String()

	dto := &identity_srv.Department{
		ID:             &id,
		Name:           &name,
		OrganizationID: &orgID,
		CreatedAt:      &model.CreatedAt,
		UpdatedAt:      &model.UpdatedAt,
	}

	// 处理可选的部门类型
	if model.DepartmentType != "" {
		dto.DepartmentType = convutil.StringPtr(model.DepartmentType)
	}

	// 处理可用设备列表 JSON 字段
	if model.AvailableEquipment != "" {
		equipmentSlice := convutil.JSONToULIDSlice(model.AvailableEquipment)
		if len(equipmentSlice) > 0 {
			// 将 []string 转换为 []core.ULID
			equipmentULIDs := make([]core.ULID, 0, len(equipmentSlice))
			equipmentULIDs = append(equipmentULIDs, equipmentSlice...)

			dto.AvailableEquipment = equipmentULIDs
		}
	}

	return dto
}

// ModelDepartmentsToThrift 将 models.Department 切片转换为 identity_srv.Department 切片
func (c *ConverterImpl) ModelDepartmentsToThrift(
	models []*models.Department,
) []*identity_srv.Department {
	if len(models) == 0 {
		return nil
	}

	dtos := make([]*identity_srv.Department, 0, len(models))
	for _, model := range models {
		if dto := c.ModelDepartmentToThrift(model); dto != nil {
			dtos = append(dtos, dto)
		}
	}

	return dtos
}

// ============================================================================
// Logic层需要的转换方法
// ============================================================================

// ModelToThrift Logic层使用的转换方法（别名）
func (c *ConverterImpl) ModelToThrift(dept *models.Department) *identity_srv.Department {
	return c.ModelDepartmentToThrift(dept)
}

// CreateRequestToModel 转换创建部门请求到模型
func (c *ConverterImpl) CreateRequestToModel(
	req *identity_srv.CreateDepartmentRequest,
) *models.Department {
	if req == nil {
		return nil
	}

	id := uuid.New()
	name := *req.Name
	orgID := uuid.MustParse(*req.OrganizationID)

	model := &models.Department{
		BaseModel: models.BaseModel{
			ID: id,
		},
		Name:           name,
		OrganizationID: orgID,
	}

	// 处理可选字段（只包含已定义的字段）
	if req.DepartmentType != nil {
		model.DepartmentType = *req.DepartmentType
	}

	return model
}

// ApplyUpdateToModel 应用更新请求到现有模型
func (c *ConverterImpl) ApplyUpdateToModel(
	existing *models.Department,
	req *identity_srv.UpdateDepartmentRequest,
) *models.Department {
	if existing == nil || req == nil {
		return existing
	}

	// 创建副本以避免修改原始模型
	updated := *existing

	// 应用更新字段（只包含已定义的字段）
	if req.Name != nil {
		updated.Name = *req.Name
	}

	if req.DepartmentType != nil {
		updated.DepartmentType = *req.DepartmentType
	}

	return &updated
}
