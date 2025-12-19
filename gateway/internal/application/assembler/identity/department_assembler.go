package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// Department Assembler
type departmentAssembler struct{}

func NewDepartmentAssembler() IDepartmentAssembler {
	return &departmentAssembler{}
}

// ToHTTPDepartment converts an RPC Department to an HTTP DepartmentDTO.
func (a *departmentAssembler) ToHTTPDepartment(
	rpc *identity_srv.Department,
) *identity.DepartmentDTO {
	if rpc == nil {
		return nil
	}

	return &identity.DepartmentDTO{
		// 核心必填字段
		ID:             rpc.ID,
		Name:           rpc.Name,
		OrganizationID: rpc.OrganizationID,

		// 可选字段
		Code:               common.CopyStringPtr(rpc.Code),
		DepartmentType:     common.CopyStringPtr(rpc.DepartmentType),
		AvailableEquipment: common.CopyStringSlice(rpc.AvailableEquipment),

		// 审计字段
		CreatedAt: common.CopyInt64Ptr(rpc.CreatedAt),
		UpdatedAt: common.CopyInt64Ptr(rpc.UpdatedAt),
		// Organization:       rpc.Organization,
		// MemberCount:        rpc.MemberCount,
	}
}

// ToHTTPDepartments converts a slice of RPC Departments to a slice of HTTP DepartmentDTOs.
func (a *departmentAssembler) ToHTTPDepartments(
	rpcDepts []*identity_srv.Department,
) []*identity.DepartmentDTO {
	if rpcDepts == nil {
		return nil
	}

	httpDepts := make([]*identity.DepartmentDTO, 0, len(rpcDepts))
	for _, rpcDept := range rpcDepts {
		httpDepts = append(httpDepts, a.ToHTTPDepartment(rpcDept))
	}

	return httpDepts
}

func (a *departmentAssembler) ToRPCCreateDeptRequest(
	dto *identity.CreateDepartmentRequestDTO,
) *identity_srv.CreateDepartmentRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.CreateDepartmentRequest{
		OrganizationID: dto.OrganizationID,
		Name:           dto.Name,
	}

	// 使用 ApplyIfSet 处理可选字段
	common.ApplyIfSet(dto.IsSetDepartmentType, dto.DepartmentType, func(v *string) {
		req.DepartmentType = v
	})

	return req
}

func (a *departmentAssembler) ToRPCUpdateDeptRequest(
	dto *identity.UpdateDepartmentRequestDTO,
) *identity_srv.UpdateDepartmentRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.UpdateDepartmentRequest{
		DepartmentID: dto.DepartmentID,
	}

	// 使用 ApplyIfSet 处理所有可选字段
	common.ApplyIfSet(dto.IsSetName, dto.Name, func(v *string) {
		req.Name = v
	})
	common.ApplyIfSet(dto.IsSetDepartmentType, dto.DepartmentType, func(v *string) {
		req.DepartmentType = v
	})

	return req
}

func (a *departmentAssembler) ToRPCGetDeptRequest(
	dto *identity.GetDepartmentRequestDTO,
) *identity_srv.GetDepartmentRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.GetDepartmentRequest{
		DepartmentID: dto.DepartmentID,
	}
}

func (a *departmentAssembler) ToRPCListDeptsRequest(
	dto *identity.GetOrganizationDepartmentsRequestDTO,
) *identity_srv.GetOrganizationDepartmentsRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.GetOrganizationDepartmentsRequest{
		OrganizationID: dto.OrganizationID,
		Page:           ToRPCPageRequest(dto.Page),
	}
}

func (a *departmentAssembler) ToHTTPListDeptsResponse(
	rpc *identity_srv.GetOrganizationDepartmentsResponse,
) *identity.GetOrganizationDepartmentsResponseDTO {
	if rpc == nil {
		return nil
	}

	return &identity.GetOrganizationDepartmentsResponseDTO{
		Departments: a.ToHTTPDepartments(rpc.Departments),
		Page:        ToHTTPPageResponse(rpc.Page),
	}
}
