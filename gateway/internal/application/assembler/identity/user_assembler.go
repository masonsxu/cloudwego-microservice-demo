package identity

import (
	identityModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// User Assembler
type userAssembler struct{}

func NewUserAssembler() IUserAssembler {
	return &userAssembler{}
}

// ToHTTPUserProfile converts an RPC UserProfile to an HTTP UserProfileDTO.
func (a *userAssembler) ToHTTPUserProfile(
	rpc *identity_srv.UserProfile,
) *identityModel.UserProfileDTO {
	if rpc == nil {
		return nil
	}

	return &identityModel.UserProfileDTO{
		// 核心必填字段
		ID:       rpc.ID,
		Username: rpc.Username,
		Status:   common.ConvertIdentityUserStatusPtrToHTTP(rpc.Status),
		Version:  &rpc.Version,

		// 基本信息（可选字段）
		Email:             common.CopyStringPtr(rpc.Email),
		Phone:             common.CopyStringPtr(rpc.Phone),
		FirstName:         common.CopyStringPtr(rpc.FirstName),
		LastName:          common.CopyStringPtr(rpc.LastName),
		RealName:          common.CopyStringPtr(rpc.RealName),
		Gender:            common.ConvertGenderPtrToHTTP(rpc.Gender),
		ProfessionalTitle: common.CopyStringPtr(rpc.ProfessionalTitle),
		LicenseNumber:     common.CopyStringPtr(rpc.LicenseNumber),
		Specialties:       common.CopyStringSlice(rpc.Specialties),
		EmployeeID:        common.CopyStringPtr(rpc.EmployeeID),

		// 状态与安全字段
		MustChangePassword: &rpc.MustChangePassword,
		AccountExpiry:      common.CopyInt64Ptr(rpc.AccountExpiry),
		LoginAttempts:      common.CopyInt32Ptr(&rpc.LoginAttempts),

		// 审计字段
		CreatedAt:     common.CopyInt64Ptr(rpc.CreatedAt),
		UpdatedAt:     common.CopyInt64Ptr(rpc.UpdatedAt),
		CreatedBy:     common.CopyStringPtr(rpc.CreatedBy),
		UpdatedBy:     common.CopyStringPtr(rpc.UpdatedBy),
		LastLoginTime: common.CopyInt64Ptr(rpc.LastLoginTime),

		// 关联信息字段
		RoleIDs:               common.CopyStringSlice(rpc.RoleIDs),
		PrimaryOrganizationID: common.CopyStringPtr(rpc.PrimaryOrganizationID),
		PrimaryDepartmentID:   common.CopyStringPtr(rpc.PrimaryDepartmentID),
	}
}

// ToHTTPUserProfiles converts a slice of RPC UserProfiles to a slice of HTTP UserProfileDTOs.
func (a *userAssembler) ToHTTPUserProfiles(
	rpcUsers []*identity_srv.UserProfile,
) []*identityModel.UserProfileDTO {
	if rpcUsers == nil {
		return nil
	}

	httpUsers := make([]*identityModel.UserProfileDTO, 0, len(rpcUsers))
	for _, rpcUser := range rpcUsers {
		httpUsers = append(httpUsers, a.ToHTTPUserProfile(rpcUser))
	}

	return httpUsers
}

// ToRPCCreateUserRequest converts an HTTP CreateUserRequestDTO to an RPC CreateUserRequest.
func (a *userAssembler) ToRPCCreateUserRequest(
	dto *identityModel.CreateUserRequestDTO,
) *identity_srv.CreateUserRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.CreateUserRequest{
		Username: dto.Username,
		Password: dto.Password,
	}

	// 使用 ApplyIfSet 处理所有可选字段
	common.ApplyIfSet(dto.IsSetEmail, dto.Email, func(v *string) { req.Email = v })
	common.ApplyIfSet(dto.IsSetPhone, dto.Phone, func(v *string) { req.Phone = v })
	common.ApplyIfSet(dto.IsSetFirstName, dto.FirstName, func(v *string) { req.FirstName = v })
	common.ApplyIfSet(dto.IsSetLastName, dto.LastName, func(v *string) { req.LastName = v })
	common.ApplyIfSet(dto.IsSetRealName, dto.RealName, func(v *string) { req.RealName = v })
	common.ApplyIfSet(dto.IsSetGender, dto.Gender, func(v *int32) {
		req.Gender = common.ConvertGenderToRPCPtr(*v)
	})
	common.ApplyIfSet(dto.IsSetProfessionalTitle, dto.ProfessionalTitle, func(v *string) {
		req.ProfessionalTitle = v
	})
	common.ApplyIfSet(dto.IsSetLicenseNumber, dto.LicenseNumber, func(v *string) {
		req.LicenseNumber = v
	})
	common.ApplyIfSetSlice(dto.IsSetSpecialties, dto.Specialties, func(v []string) {
		req.Specialties = v
	})
	common.ApplyIfSet(dto.IsSetEmployeeID, dto.EmployeeID, func(v *string) { req.EmployeeID = v })
	common.ApplyIfSet(dto.IsSetMustChangePassword, dto.MustChangePassword, func(v *bool) {
		req.MustChangePassword = v
	})
	common.ApplyIfSet(dto.IsSetAccountExpiry, dto.AccountExpiry, func(v *int64) {
		req.AccountExpiry = v
	})

	return req
}

func (a *userAssembler) ToRPCGetUserRequest(
	dto *identityModel.GetUserRequestDTO,
) *identity_srv.GetUserRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.GetUserRequest{
		UserID: dto.UserID,
	}
}

// ToRPCUpdateUserRequest converts an HTTP UpdateUserRequestDTO to an RPC UpdateUserRequest.
func (a *userAssembler) ToRPCUpdateUserRequest(
	dto *identityModel.UpdateUserRequestDTO,
) *identity_srv.UpdateUserRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.UpdateUserRequest{
		UserID:  dto.UserID,
		Version: dto.Version,
	}

	// 使用 ApplyIfSet 处理所有可选字段
	common.ApplyIfSet(dto.IsSetEmail, dto.Email, func(v *string) { req.Email = v })
	common.ApplyIfSet(dto.IsSetPhone, dto.Phone, func(v *string) { req.Phone = v })
	common.ApplyIfSet(dto.IsSetFirstName, dto.FirstName, func(v *string) { req.FirstName = v })
	common.ApplyIfSet(dto.IsSetLastName, dto.LastName, func(v *string) { req.LastName = v })
	common.ApplyIfSet(dto.IsSetRealName, dto.RealName, func(v *string) { req.RealName = v })
	common.ApplyIfSet(dto.IsSetGender, dto.Gender, func(v *int32) {
		req.Gender = common.ConvertGenderToRPCPtr(*v)
	})
	common.ApplyIfSet(dto.IsSetProfessionalTitle, dto.ProfessionalTitle, func(v *string) {
		req.ProfessionalTitle = v
	})
	common.ApplyIfSet(dto.IsSetLicenseNumber, dto.LicenseNumber, func(v *string) {
		req.LicenseNumber = v
	})
	common.ApplyIfSetSlice(dto.IsSetSpecialties, dto.Specialties, func(v []string) {
		req.Specialties = v
	})
	common.ApplyIfSet(dto.IsSetEmployeeID, dto.EmployeeID, func(v *string) { req.EmployeeID = v })
	common.ApplyIfSet(dto.IsSetAccountExpiry, dto.AccountExpiry, func(v *int64) {
		req.AccountExpiry = v
	})

	return req
}

// ToRPCUpdateMeRequest converts an HTTP UpdateMeRequestDTO to an RPC UpdateUserRequest.
// UserID will be set from the authentication context, not from the request.
func (a *userAssembler) ToRPCUpdateMeRequest(
	dto *identityModel.UpdateMeRequestDTO,
) *identity_srv.UpdateUserRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.UpdateUserRequest{
		Version: dto.Version,
		// UserID 将在 handler 层从上下文中获取并设置
	}

	// 使用 ApplyIfSet 处理所有可选字段
	common.ApplyIfSet(dto.IsSetEmail, dto.Email, func(v *string) { req.Email = v })
	common.ApplyIfSet(dto.IsSetPhone, dto.Phone, func(v *string) { req.Phone = v })
	common.ApplyIfSet(dto.IsSetFirstName, dto.FirstName, func(v *string) { req.FirstName = v })
	common.ApplyIfSet(dto.IsSetLastName, dto.LastName, func(v *string) { req.LastName = v })
	common.ApplyIfSet(dto.IsSetRealName, dto.RealName, func(v *string) { req.RealName = v })
	common.ApplyIfSet(dto.IsSetGender, dto.Gender, func(v *int32) {
		req.Gender = common.ConvertGenderToRPCPtr(*v)
	})
	common.ApplyIfSet(dto.IsSetProfessionalTitle, dto.ProfessionalTitle, func(v *string) {
		req.ProfessionalTitle = v
	})
	common.ApplyIfSet(dto.IsSetLicenseNumber, dto.LicenseNumber, func(v *string) {
		req.LicenseNumber = v
	})
	common.ApplyIfSetSlice(dto.IsSetSpecialties, dto.Specialties, func(v []string) {
		req.Specialties = v
	})
	common.ApplyIfSet(dto.IsSetEmployeeID, dto.EmployeeID, func(v *string) { req.EmployeeID = v })
	common.ApplyIfSet(dto.IsSetAccountExpiry, dto.AccountExpiry, func(v *int64) {
		req.AccountExpiry = v
	})

	return req
}

func (a *userAssembler) ToRPCDeleteUserRequest(
	dto *identityModel.DeleteUserRequestDTO,
) *identity_srv.DeleteUserRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.DeleteUserRequest{
		UserID: dto.UserID,
	}

	return req
}

func (a *userAssembler) ToRPCListUsersRequest(
	dto *identityModel.ListUsersRequestDTO,
) *identity_srv.ListUsersRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.ListUsersRequest{
		Page:           ToRPCPageRequest(dto.Page),
		OrganizationID: dto.OrganizationID,
		Status:         common.ConvertIdentityUserStatusPtrToRPCPtr(dto.Status),
	}
}

func (a *userAssembler) ToHTTPListUsersResponse(
	rpc *identity_srv.ListUsersResponse,
) *identityModel.ListUsersResponseDTO {
	if rpc == nil {
		return nil
	}

	return &identityModel.ListUsersResponseDTO{
		Users: a.ToHTTPUserProfiles(rpc.Users),
		Page:  ToHTTPPageResponse(rpc.Page),
	}
}

func (a *userAssembler) ToRPCSearchUsersRequest(
	rpc *identityModel.SearchUsersRequestDTO,
) *identity_srv.SearchUsersRequest {
	if rpc == nil {
		return nil
	}

	return &identity_srv.SearchUsersRequest{
		Page:           ToRPCPageRequest(rpc.Page),
		OrganizationID: rpc.OrganizationID,
	}
}

func (a *userAssembler) ToHTTPSearchUsersResponse(
	rpc *identity_srv.SearchUsersResponse,
) *identityModel.SearchUsersResponseDTO {
	if rpc == nil {
		return nil
	}

	return &identityModel.SearchUsersResponseDTO{
		Users: a.ToHTTPUserProfiles(rpc.Users),
		Page:  ToHTTPPageResponse(rpc.Page),
	}
}

func (a *userAssembler) ToRPCChangeUserStatusRequest(
	http *identityModel.ChangeUserStatusRequestDTO,
) *identity_srv.ChangeUserStatusRequest {
	if http == nil {
		return nil
	}

	return &identity_srv.ChangeUserStatusRequest{
		UserID:     http.UserID,
		NewStatus_: common.ConvertIdentityUserStatusToRPCPtr(*http.NewStatus),
	}
}

func (a *userAssembler) ToRPCUnlockUserRequest(
	http *identityModel.UnlockUserRequestDTO,
) *identity_srv.UnlockUserRequest {
	if http == nil {
		return nil
	}

	return &identity_srv.UnlockUserRequest{
		UserID: http.UserID,
	}
}
