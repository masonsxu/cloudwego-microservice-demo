package identity

import (
	"context"

	hertzZerolog "github.com/hertz-contrib/logger/zerolog"

	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	identityassembler "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/domain/common"
	identitycli "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/client/identity_cli"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// auditLogServiceImpl 审计日志查询服务实现
type auditLogServiceImpl struct {
	*common.BaseService
	identityClient identitycli.IdentityClient
	assembler      identityassembler.Assembler
}

// NewAuditLogService 创建新的审计日志查询服务实例
func NewAuditLogService(
	identityClient identitycli.IdentityClient,
	assembler identityassembler.Assembler,
	logger *hertzZerolog.Logger,
) AuditLogService {
	return &auditLogServiceImpl{
		BaseService:    common.NewBaseService(logger),
		identityClient: identityClient,
		assembler:      assembler,
	}
}

// =================================================================
// 审计日志模块 (Audit Log)
// =================================================================

func (s *auditLogServiceImpl) ListAuditLogs(
	ctx context.Context,
	req *identity.ListAuditLogsRequestDTO,
) (*identity.ListAuditLogsResponseDTO, error) {
	// 转换请求
	rpcReq := s.assembler.AuditLog().ToRPCListAuditLogsRequest(req)

	// 调用RPC服务
	result, err := s.ProcessRPCCall(ctx, "查询审计日志",
		func(ctx context.Context) (interface{}, error) {
			return s.identityClient.ListAuditLogs(ctx, rpcReq)
		},
	)
	if err != nil {
		return nil, err
	}

	rpcResp := result.(*identity_srv.ListAuditLogsResponse)

	// 转换响应
	httpResp := s.assembler.AuditLog().ToHTTPListAuditLogsResponse(rpcResp)
	httpResp.BaseResp = s.ResponseBuilder().BuildSuccessResponse()

	return httpResp, nil
}
