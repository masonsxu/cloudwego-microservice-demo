package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// Membership Assembler
type membershipAssembler struct{}

func NewMembershipAssembler() IMembershipAssembler {
	return &membershipAssembler{}
}

// ToHTTPUserMembership converts an RPC UserMembership to an HTTP UserMembershipDTO.
func (a *membershipAssembler) ToHTTPUserMembership(
	rpc *identity_srv.UserMembership,
) *identity.UserMembershipDTO {
	if rpc == nil {
		return nil
	}

	return &identity.UserMembershipDTO{
		// 核心必填字段
		ID:             rpc.ID,
		UserID:         rpc.UserID,
		OrganizationID: rpc.OrganizationID,

		// 可选字段
		DepartmentID: common.CopyStringPtr(rpc.DepartmentID),
		IsPrimary:    common.CopyBoolPtr(rpc.IsPrimary),

		// 审计字段
		CreatedAt: common.CopyInt64Ptr(rpc.CreatedAt),
		UpdatedAt: common.CopyInt64Ptr(rpc.UpdatedAt),
		// Organization:   NewOrgAssembler().ToHTTPOrganization(rpc.Organization),
		// Department:     NewOrgAssembler().ToHTTPDepartment(rpc.Department),
	}
}

// ToHTTPUserMemberships converts a slice of RPC UserMemberships to a slice of HTTP UserMembershipDTOs.
func (a *membershipAssembler) ToHTTPUserMemberships(
	rpcMemberships []*identity_srv.UserMembership,
) []*identity.UserMembershipDTO {
	if rpcMemberships == nil {
		return nil
	}

	httpMemberships := make([]*identity.UserMembershipDTO, 0, len(rpcMemberships))
	for _, rpcMembership := range rpcMemberships {
		httpMemberships = append(httpMemberships, a.ToHTTPUserMembership(rpcMembership))
	}

	return httpMemberships
}

func (a *membershipAssembler) ToRPCGetUserMembershipsRequest(
	dto *identity.GetUserMembershipsRequestDTO,
) *identity_srv.GetUserMembershipsRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.GetUserMembershipsRequest{
		UserID: dto.UserID,
		Page:   ToRPCPageRequest(dto.Page),
	}
}

func (a *membershipAssembler) ToHTTPGetUserMembershipsResponse(
	rpc *identity_srv.GetUserMembershipsResponse,
) *identity.GetUserMembershipsResponseDTO {
	if rpc == nil {
		return nil
	}

	return &identity.GetUserMembershipsResponseDTO{
		Memberships: a.ToHTTPUserMemberships(rpc.Memberships),
		Page:        ToHTTPPageResponse(rpc.Page),
	}
}
