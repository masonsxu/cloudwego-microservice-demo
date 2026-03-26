package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	core "github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/core"
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
		Id:   rpc.Id,
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

	if dto.ParentID != nil {
		req.ParentID = dto.ParentID
	}
	if dto.FacilityType != nil {
		req.FacilityType = dto.FacilityType
	}
	if dto.AccreditationStatus != nil {
		req.AccreditationStatus = dto.AccreditationStatus
	}
	if len(dto.ProvinceCity) > 0 {
		req.ProvinceCity = &core.StringListValue{Items: dto.ProvinceCity}
	}

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

	if dto.Name != nil {
		req.Name = dto.Name
	}
	if dto.ParentID != nil {
		req.ParentID = dto.ParentID
	}
	if dto.FacilityType != nil {
		req.FacilityType = dto.FacilityType
	}
	if dto.AccreditationStatus != nil {
		req.AccreditationStatus = dto.AccreditationStatus
	}
	if dto.ProvinceCity != nil {
		items := make([]string, 0, len(dto.ProvinceCity.GetValues()))
		for _, value := range dto.ProvinceCity.GetValues() {
			if str := value.GetStringValue(); str != "" {
				items = append(items, str)
			}
		}
		if len(items) > 0 {
			req.ProvinceCity = &core.StringListValue{Items: items}
		}
	}

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
