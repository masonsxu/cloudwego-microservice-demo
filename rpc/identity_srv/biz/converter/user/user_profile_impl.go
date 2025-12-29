package user

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/convutil"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/converter/enum"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl 用户转换器实现
type ConverterImpl struct {
	enumConverter enum.Converter
	baseConverter base.Converter
}

// NewConverter 创建一个新的用户转换器实例
func NewConverter(enumConverter enum.Converter) Converter {
	return &ConverterImpl{
		enumConverter: enumConverter,
		baseConverter: base.NewConverter(),
	}
}

// ============================================================================
// UserProfile 用户档案转换
// ============================================================================

// ModelUserProfileToThrift 将 models.UserProfile 转换为 identity_srv.UserProfile
func (c *ConverterImpl) ModelUserProfileToThrift(
	model *models.UserProfile,
) *identity_srv.UserProfile {
	if model == nil {
		return nil
	}

	id := model.ID.String()
	username := model.Username
	status := c.enumConverter.ModelUserStatusToThrift(model.Status)

	dto := &identity_srv.UserProfile{
		ID:        &id,
		Username:  &username,
		Status:    &status,
		CreatedAt: &model.CreatedAt,
		UpdatedAt: &model.UpdatedAt,
	}

	// 处理可选字符串字段
	if model.Email != "" {
		dto.Email = convutil.StringPtr(model.Email)
	}

	if model.Phone != "" {
		dto.Phone = convutil.StringPtr(model.Phone)
	}

	if model.FirstName != "" {
		dto.FirstName = convutil.StringPtr(model.FirstName)
	}

	if model.LastName != "" {
		dto.LastName = convutil.StringPtr(model.LastName)
	}

	if model.RealName != "" {
		dto.RealName = convutil.StringPtr(model.RealName)
	}

	// 处理性别字段
	if model.Gender != 0 {
		gender := c.enumConverter.ModelGenderToThrift(model.Gender)
		dto.Gender = &gender
	}

	if model.ProfessionalTitle != "" {
		dto.ProfessionalTitle = convutil.StringPtr(model.ProfessionalTitle)
	}

	if model.LicenseNumber != "" {
		dto.LicenseNumber = convutil.StringPtr(model.LicenseNumber)
	}

	// 处理 Specialties JSON 字段
	if model.Specialties != "" {
		specialtiesSlice := convutil.JSONToStringSlice(model.Specialties)
		if len(specialtiesSlice) > 0 {
			dto.Specialties = specialtiesSlice
		}
	}

	if model.EmployeeID != "" {
		dto.EmployeeID = convutil.StringPtr(model.EmployeeID)
	}

	// 处理可选数值字段
	if model.LoginAttempts > 0 {
		dto.LoginAttempts = model.LoginAttempts
	}

	dto.MustChangePassword = model.MustChangePassword

	// 处理可选时间戳字段
	if model.AccountExpiry != nil && *model.AccountExpiry > 0 {
		dto.AccountExpiry = model.AccountExpiry
	}

	if model.CreatedBy != nil {
		dto.CreatedBy = convutil.StringPtr(model.CreatedBy.String())
	}

	if model.UpdatedBy != nil {
		dto.UpdatedBy = convutil.StringPtr(model.UpdatedBy.String())
	}

	if model.LastLoginTime != nil && *model.LastLoginTime > 0 {
		dto.LastLoginTime = model.LastLoginTime
	}

	// 注意：roleIDs, primaryOrganizationID, primaryDepartmentID 需要在业务逻辑层填充
	// 这些字段不存储在 UserProfile 模型中，需要通过关联查询获取

	return dto
}

// ============================================================================
// Logic层需要的转换方法
// ============================================================================

// CreateUserRequestToModel 将创建用户请求转换为模型
func (c *ConverterImpl) CreateUserRequestToModel(
	req *identity_srv.CreateUserRequest,
) *models.UserProfile {
	if req == nil {
		return nil
	}

	model := &models.UserProfile{
		Username: *req.Username,
		Status:   models.UserStatusActive, // 默认激活状态
	}

	// 处理密码哈希
	if req.Password != nil {
		if hash, err := convutil.HashPassword(*req.Password); err == nil {
			model.PasswordHash = hash
		}
	}

	// 处理可选字段
	if req.Email != nil {
		model.Email = *req.Email
	}

	if req.Phone != nil {
		model.Phone = *req.Phone
	}

	if req.FirstName != nil {
		model.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		model.LastName = *req.LastName
	}

	if req.RealName != nil {
		model.RealName = *req.RealName
	}

	if req.Gender != nil {
		model.Gender = c.enumConverter.ThriftGenderToModel(*req.Gender)
	}

	if req.ProfessionalTitle != nil {
		model.ProfessionalTitle = *req.ProfessionalTitle
	}

	if req.LicenseNumber != nil {
		model.LicenseNumber = *req.LicenseNumber
	}

	if len(req.Specialties) > 0 {
		model.Specialties = convutil.StringSliceToJSON(req.Specialties)
	}

	if req.EmployeeID != nil {
		model.EmployeeID = *req.EmployeeID
	}

	if req.MustChangePassword != nil {
		model.MustChangePassword = *req.MustChangePassword
	}

	if req.AccountExpiry != nil {
		timestamp := *req.AccountExpiry
		model.AccountExpiry = &timestamp
	}

	return model
}

// ApplyUpdateUserToModel 将更新请求应用到现有模型
func (c *ConverterImpl) ApplyUpdateUserToModel(
	existing *models.UserProfile,
	req *identity_srv.UpdateUserRequest,
) *models.UserProfile {
	if existing == nil || req == nil {
		return existing
	}

	// 处理可选字段更新
	if req.Email != nil {
		existing.Email = *req.Email
	}

	if req.Phone != nil {
		existing.Phone = *req.Phone
	}

	if req.FirstName != nil {
		existing.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		existing.LastName = *req.LastName
	}

	if req.RealName != nil {
		existing.RealName = *req.RealName
	}

	if req.Gender != nil {
		existing.Gender = c.enumConverter.ThriftGenderToModel(*req.Gender)
	}

	if req.ProfessionalTitle != nil {
		existing.ProfessionalTitle = *req.ProfessionalTitle
	}

	if req.LicenseNumber != nil {
		existing.LicenseNumber = *req.LicenseNumber
	}

	if len(req.Specialties) > 0 {
		existing.Specialties = convutil.StringSliceToJSON(req.Specialties)
	}

	if req.EmployeeID != nil {
		existing.EmployeeID = *req.EmployeeID
	}

	if req.AccountExpiry != nil {
		timestamp := *req.AccountExpiry
		existing.AccountExpiry = &timestamp
	}

	return existing
}
