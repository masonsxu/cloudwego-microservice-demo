package identity

import (
	identityModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// auditLogAssembler 审计日志组装器实现
type auditLogAssembler struct{}

// NewAuditLogAssembler 创建审计日志组装器
func NewAuditLogAssembler() IAuditLogAssembler {
	return &auditLogAssembler{}
}

// ToHTTPAuditLog 将 RPC AuditLog 转换为 HTTP AuditLogDTO
func (a *auditLogAssembler) ToHTTPAuditLog(rpc *identity_srv.AuditLog) *identityModel.AuditLogDTO {
	if rpc == nil {
		return nil
	}

	dto := &identityModel.AuditLogDTO{
		Id:             common.CopyStringPtr(rpc.Id),
		RequestID:      common.CopyStringPtr(rpc.RequestID),
		TraceID:        common.CopyStringPtr(rpc.TraceID),
		UserID:         common.CopyStringPtr(rpc.UserID),
		Username:       common.CopyStringPtr(rpc.Username),
		OrganizationID: common.CopyStringPtr(rpc.OrganizationID),
		Resource:       common.CopyStringPtr(rpc.Resource),
		ResourceID:     common.CopyStringPtr(rpc.ResourceID),
		StatusCode:     common.CopyInt32Ptr(rpc.StatusCode),
		Success:        common.CopyBoolPtr(rpc.Success),
		ClientIP:       common.CopyStringPtr(rpc.ClientIP),
		UserAgent:      common.CopyStringPtr(rpc.UserAgent),
		RequestBody:    common.CopyStringPtr(rpc.RequestBody),
		DurationMs:     common.CopyInt32Ptr(rpc.DurationMs),
		CreatedAt:      common.CopyInt64Ptr(rpc.CreatedAt),
	}

	// 转换枚举类型 Action 为 int32
	if rpc.Action != nil {
		action := int32(*rpc.Action)
		dto.Action = &action
	}

	return dto
}

// ToHTTPAuditLogs 批量转换 RPC AuditLog 到 HTTP AuditLogDTO
func (a *auditLogAssembler) ToHTTPAuditLogs(
	rpcs []*identity_srv.AuditLog,
) []*identityModel.AuditLogDTO {
	if len(rpcs) == 0 {
		return nil
	}

	dtos := make([]*identityModel.AuditLogDTO, 0, len(rpcs))
	for _, rpc := range rpcs {
		if dto := a.ToHTTPAuditLog(rpc); dto != nil {
			dtos = append(dtos, dto)
		}
	}

	return dtos
}

// ToRPCListAuditLogsRequest 将 HTTP ListAuditLogsRequestDTO 转换为 RPC ListAuditLogsRequest
func (a *auditLogAssembler) ToRPCListAuditLogsRequest(
	dto *identityModel.ListAuditLogsRequestDTO,
) *identity_srv.ListAuditLogsRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.ListAuditLogsRequest{
		Page:      ToRPCPageRequest(dto.Page),
		UserID:    dto.UserID,
		Resource:  dto.Resource,
		Success:   dto.Success,
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
	}

	// 转换 action 从 int32 到 AuditAction 枚举
	if dto.Action != nil {
		action := identity_srv.AuditAction(*dto.Action)
		req.Action = &action
	}

	return req
}

// ToHTTPListAuditLogsResponse 将 RPC ListAuditLogsResponse 转换为 HTTP ListAuditLogsResponseDTO
func (a *auditLogAssembler) ToHTTPListAuditLogsResponse(
	rpc *identity_srv.ListAuditLogsResponse,
) *identityModel.ListAuditLogsResponseDTO {
	if rpc == nil {
		return nil
	}

	resp := &identityModel.ListAuditLogsResponseDTO{
		AuditLogs: a.ToHTTPAuditLogs(rpc.AuditLogs),
		Page:      ToHTTPPageResponse(rpc.Page),
	}

	if rpc.Stats != nil {
		resp.Stats = &identityModel.AuditLogStatsDTO{
			TotalCount:    common.CopyInt64Ptr(rpc.Stats.TotalCount),
			SuccessCount:  common.CopyInt64Ptr(rpc.Stats.SuccessCount),
			AvgDurationMs: common.CopyInt32Ptr(rpc.Stats.AvgDurationMs),
		}
	}

	return resp
}
