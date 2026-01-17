package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
	identityModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	permissionModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
)

// Assembler 身份管理聚合组装器接口 - 统一暴露给Service层
// 提供所有身份相关的协议转换能力，避免Service层直接依赖多个细分Assembler
type Assembler interface {
	// 获取各个业务领域的组装器
	Auth() IAuthAssembler
	User() IUserAssembler
	Organization() IOrgAssembler
	Department() IDepartmentAssembler
	Membership() IMembershipAssembler
	Logo() ILogoAssembler

	// 通用转换方法（避免重复代码）
	ToHTTPPageResponse(*rpc_base.PageResponse) *http_base.PageResponseDTO
	ToRPCPageRequest(*http_base.PageRequestDTO) *rpc_base.PageRequest
}

type IAuthAssembler interface {
	ToRPCLoginRequest(*identityModel.LoginRequestDTO) *identity_srv.LoginRequest
	ToHTTPLoginResponse(*identity_srv.LoginResponse) *identityModel.LoginResponseDTO
	ToRPCChangePasswordRequest(
		*identityModel.ChangePasswordRequestDTO,
	) *identity_srv.ChangePasswordRequest
	ToRPCResetPasswordRequest(
		*identityModel.ResetPasswordRequestDTO,
	) *identity_srv.ResetPasswordRequest
	ToRPCForcePasswordChangeRequest(
		*identityModel.ForcePasswordChangeRequestDTO,
	) *identity_srv.ForcePasswordChangeRequest
	ToHTTPMenuTree([]*identity_srv.MenuNode) []*permissionModel.MenuNodeDTO
	ToHTTPMenuPermissions([]*identity_srv.MenuPermission) []*permissionModel.MenuPermissionDTO
	ToHTTPRoleInfos([]*identity_srv.RoleDefinition) []*identityModel.RoleInfoDTO
}

type IDepartmentAssembler interface {
	ToHTTPDepartment(*identity_srv.Department) *identityModel.DepartmentDTO
	ToHTTPDepartments([]*identity_srv.Department) []*identityModel.DepartmentDTO
	ToRPCCreateDeptRequest(
		*identityModel.CreateDepartmentRequestDTO,
	) *identity_srv.CreateDepartmentRequest
	ToRPCUpdateDeptRequest(
		*identityModel.UpdateDepartmentRequestDTO,
	) *identity_srv.UpdateDepartmentRequest
	ToRPCGetDeptRequest(
		*identityModel.GetDepartmentRequestDTO,
	) *identity_srv.GetDepartmentRequest

	ToRPCListDeptsRequest(
		*identityModel.GetOrganizationDepartmentsRequestDTO,
	) *identity_srv.GetOrganizationDepartmentsRequest
	ToHTTPListDeptsResponse(
		*identity_srv.GetOrganizationDepartmentsResponse,
	) *identityModel.GetOrganizationDepartmentsResponseDTO
}

type IUserAssembler interface {
	ToHTTPUserProfile(*identity_srv.UserProfile) *identityModel.UserProfileDTO
	ToHTTPUserProfiles([]*identity_srv.UserProfile) []*identityModel.UserProfileDTO
	ToRPCCreateUserRequest(*identityModel.CreateUserRequestDTO) *identity_srv.CreateUserRequest
	ToRPCGetUserRequest(*identityModel.GetUserRequestDTO) *identity_srv.GetUserRequest
	ToRPCUpdateUserRequest(*identityModel.UpdateUserRequestDTO) *identity_srv.UpdateUserRequest
	ToRPCUpdateMeRequest(*identityModel.UpdateMeRequestDTO) *identity_srv.UpdateUserRequest
	ToRPCDeleteUserRequest(*identityModel.DeleteUserRequestDTO) *identity_srv.DeleteUserRequest
	ToRPCListUsersRequest(*identityModel.ListUsersRequestDTO) *identity_srv.ListUsersRequest
	ToHTTPListUsersResponse(*identity_srv.ListUsersResponse) *identityModel.ListUsersResponseDTO
	ToRPCSearchUsersRequest(*identityModel.SearchUsersRequestDTO) *identity_srv.SearchUsersRequest
	ToHTTPSearchUsersResponse(
		*identity_srv.SearchUsersResponse,
	) *identityModel.SearchUsersResponseDTO
	ToRPCChangeUserStatusRequest(
		*identityModel.ChangeUserStatusRequestDTO,
	) *identity_srv.ChangeUserStatusRequest
	ToRPCUnlockUserRequest(*identityModel.UnlockUserRequestDTO) *identity_srv.UnlockUserRequest
}

type IOrgAssembler interface {
	ToHTTPOrganization(*identity_srv.Organization) *identityModel.OrganizationDTO
	ToHTTPOrganizations([]*identity_srv.Organization) []*identityModel.OrganizationDTO
	ToRPCCreateOrgRequest(
		*identityModel.CreateOrganizationRequestDTO,
	) *identity_srv.CreateOrganizationRequest
	ToRPCUpdateOrgRequest(
		*identityModel.UpdateOrganizationRequestDTO,
	) *identity_srv.UpdateOrganizationRequest
	ToRPCGetOrgRequest(
		*identityModel.GetOrganizationRequestDTO,
	) *identity_srv.GetOrganizationRequest
	ToRPCListOrgsRequest(
		*identityModel.ListOrganizationsRequestDTO,
	) *identity_srv.ListOrganizationsRequest
	ToHTTPListOrgsResponse(
		*identity_srv.ListOrganizationsResponse,
	) *identityModel.ListOrganizationsResponseDTO
}

type IMembershipAssembler interface {
	ToHTTPUserMembership(*identity_srv.UserMembership) *identityModel.UserMembershipDTO
	ToHTTPUserMemberships([]*identity_srv.UserMembership) []*identityModel.UserMembershipDTO
	ToRPCGetUserMembershipsRequest(
		*identityModel.GetUserMembershipsRequestDTO,
	) *identity_srv.GetUserMembershipsRequest
	ToHTTPGetUserMembershipsResponse(
		*identity_srv.GetUserMembershipsResponse,
	) *identityModel.GetUserMembershipsResponseDTO
}

type ILogoAssembler interface {
	ToHTTPOrganizationLogo(*identity_srv.OrganizationLogo) *identityModel.OrganizationLogoDTO
	ToRPCUploadTemporaryLogoRequest(
		dto *identityModel.UploadTemporaryLogoRequestDTO,
		userID string,
	) *identity_srv.UploadTemporaryLogoRequest
	ToRPCGetOrganizationLogoRequest(
		*identityModel.GetOrganizationLogoRequestDTO,
	) *identity_srv.GetOrganizationLogoRequest
	ToRPCDeleteOrganizationLogoRequest(
		*identityModel.DeleteOrganizationLogoRequestDTO,
	) *identity_srv.DeleteOrganizationLogoRequest
	ToRPCBindLogoToOrganizationRequest(
		*identityModel.BindLogoToOrganizationRequestDTO,
	) *identity_srv.BindLogoToOrganizationRequest
}

// ToHTTPPageResponse is a generic function to convert RPC PageResponse to HTTP PageResponseDTO.
func ToHTTPPageResponse(rpc *rpc_base.PageResponse) *http_base.PageResponseDTO {
	if rpc == nil {
		return nil
	}

	return &http_base.PageResponseDTO{
		Total:      rpc.Total,
		Page:       rpc.Page,
		Limit:      rpc.Limit,
		TotalPages: rpc.TotalPages,
		HasNext:    rpc.HasNext,
		HasPrev:    rpc.HasPrev,
	}
}

// ToRPCPageRequest is a generic function to convert HTTP PageRequestDTO to RPC PageRequest.
func ToRPCPageRequest(http *http_base.PageRequestDTO) *rpc_base.PageRequest {
	if http == nil {
		return nil
	}

	return &rpc_base.PageRequest{
		Page:         http.Page,
		Limit:        http.Limit,
		Search:       http.Search,
		Filter:       http.Filter,
		Sort:         http.Sort,
		Fields:       http.Fields,
		IncludeTotal: http.IncludeTotal,
		FetchAll:     http.FetchAll,
	}
}
