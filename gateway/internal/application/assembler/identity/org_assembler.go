package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// Org Assembler
type orgAssembler struct{}

func NewOrgAssembler() IOrgAssembler {
	return &orgAssembler{}
}

// ToHTTPOrganization converts an RPC Organization to an HTTP OrganizationDTO.
func (a *orgAssembler) ToHTTPOrganization(
	rpc *identity_srv.Organization,
) *identity.OrganizationDTO {
	if rpc == nil {
		return nil
	}

	return &identity.OrganizationDTO{
		// 核心必填字段
		ID:   rpc.ID,
		Name: rpc.Name,

		// 可选字段
		Code:                common.CopyStringPtr(rpc.Code),
		ParentID:            common.CopyStringPtr(rpc.ParentID),
		FacilityType:        common.CopyStringPtr(rpc.FacilityType),
		AccreditationStatus: common.CopyStringPtr(rpc.AccreditationStatus),
		Logo:                common.CopyStringPtr(rpc.Logo),
		LogoID:              common.CopyStringPtr(rpc.LogoID),
		ProvinceCity:        common.CopyStringSlice(rpc.ProvinceCity),

		// 审计字段
		CreatedAt: common.CopyInt64Ptr(rpc.CreatedAt),
		UpdatedAt: common.CopyInt64Ptr(rpc.UpdatedAt),
	}
}

// ToHTTPOrganizations converts a slice of RPC Organizations to a slice of HTTP OrganizationDTOs.
func (a *orgAssembler) ToHTTPOrganizations(
	rpcOrgs []*identity_srv.Organization,
) []*identity.OrganizationDTO {
	if rpcOrgs == nil {
		return nil
	}

	httpOrgs := make([]*identity.OrganizationDTO, 0, len(rpcOrgs))
	for _, rpcOrg := range rpcOrgs {
		httpOrgs = append(httpOrgs, a.ToHTTPOrganization(rpcOrg))
	}

	return httpOrgs
}

func (a *orgAssembler) ToRPCCreateOrgRequest(
	dto *identity.CreateOrganizationRequestDTO,
) *identity_srv.CreateOrganizationRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.CreateOrganizationRequest{
		Name: dto.Name,
	}

	// 使用 ApplyIfSet 处理所有可选字段
	common.ApplyIfSet(dto.IsSetParentID, dto.ParentID, func(v *string) {
		req.ParentID = v
	})
	common.ApplyIfSet(dto.IsSetFacilityType, dto.FacilityType, func(v *string) {
		req.FacilityType = v
	})
	common.ApplyIfSet(dto.IsSetAccreditationStatus, dto.AccreditationStatus, func(v *string) {
		req.AccreditationStatus = v
	})
	common.ApplyIfSetSlice(dto.IsSetProvinceCity, dto.ProvinceCity, func(v []string) {
		req.ProvinceCity = v
	})

	return req
}

func (a *orgAssembler) ToRPCGetOrgRequest(
	dto *identity.GetOrganizationRequestDTO,
) *identity_srv.GetOrganizationRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.GetOrganizationRequest{
		OrganizationID: dto.OrganizationID,
	}
}

func (a *orgAssembler) ToRPCUpdateOrgRequest(
	dto *identity.UpdateOrganizationRequestDTO,
) *identity_srv.UpdateOrganizationRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.UpdateOrganizationRequest{
		OrganizationID: dto.OrganizationID,
	}

	// 使用 ApplyIfSet 处理所有可选字段
	common.ApplyIfSet(dto.IsSetName, dto.Name, func(v *string) {
		req.Name = v
	})
	common.ApplyIfSet(dto.IsSetParentID, dto.ParentID, func(v *string) {
		req.ParentID = v
	})
	common.ApplyIfSet(dto.IsSetFacilityType, dto.FacilityType, func(v *string) {
		req.FacilityType = v
	})
	common.ApplyIfSet(dto.IsSetAccreditationStatus, dto.AccreditationStatus, func(v *string) {
		req.AccreditationStatus = v
	})
	common.ApplyIfSetSlice(dto.IsSetProvinceCity, dto.ProvinceCity, func(v []string) {
		req.ProvinceCity = v
	})

	return req
}

func (a *orgAssembler) ToRPCListOrgsRequest(
	dto *identity.ListOrganizationsRequestDTO,
) *identity_srv.ListOrganizationsRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.ListOrganizationsRequest{
		ParentID: dto.ParentID,
		Page:     ToRPCPageRequest(dto.Page),
	}
}

func (a *orgAssembler) ToHTTPListOrgsResponse(
	rpc *identity_srv.ListOrganizationsResponse,
) *identity.ListOrganizationsResponseDTO {
	if rpc == nil {
		return nil
	}

	return &identity.ListOrganizationsResponseDTO{
		Organizations: a.ToHTTPOrganizations(rpc.Organizations),
		Page:          ToHTTPPageResponse(rpc.Page),
	}
}
