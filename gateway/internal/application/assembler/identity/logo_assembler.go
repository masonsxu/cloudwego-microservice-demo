package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// Logo Assembler
type logoAssembler struct{}

func NewLogoAssembler() ILogoAssembler {
	return &logoAssembler{}
}

// ToHTTPOrganizationLogo converts an RPC OrganizationLogo to an HTTP OrganizationLogoDTO.
func (a *logoAssembler) ToHTTPOrganizationLogo(
	rpc *identity_srv.OrganizationLogo,
) *identity.OrganizationLogoDTO {
	if rpc == nil {
		return nil
	}

	dto := &identity.OrganizationLogoDTO{
		// 核心字段
		ID:     common.CopyStringPtr(rpc.ID),
		FileID: common.CopyStringPtr(rpc.FileID),

		// 绑定信息
		BoundOrganizationID: common.CopyStringPtr(rpc.BoundOrganizationID),

		// 文件元信息
		FileName: common.CopyStringPtr(rpc.FileName),
		FileSize: common.CopyInt64Ptr(rpc.FileSize),
		MimeType: common.CopyStringPtr(rpc.MimeType),

		// 时间字段
		ExpiresAt: common.CopyInt64Ptr(rpc.ExpiresAt),
		CreatedAt: common.CopyInt64Ptr(rpc.CreatedAt),
		UpdatedAt: common.CopyInt64Ptr(rpc.UpdatedAt),

		// 下载URL（预签名）
		DownloadUrl: common.CopyStringPtr(rpc.DownloadUrl),
	}

	// 转换枚举类型 Status 为字符串
	if rpc.Status != nil {
		statusStr := rpc.Status.String()
		dto.Status = &statusStr
	}

	return dto
}

// ToRPCUploadTemporaryLogoRequest converts HTTP UploadTemporaryLogoRequestDTO to RPC request.
func (a *logoAssembler) ToRPCUploadTemporaryLogoRequest(
	dto *identity.UploadTemporaryLogoRequestDTO,
	userID string,
) *identity_srv.UploadTemporaryLogoRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.UploadTemporaryLogoRequest{
		FileContent: dto.FileContent,
		FileName:    dto.FileName,
		MimeType:    dto.MimeType,
		UploadedBy:  &userID,
	}

	return req
}

// ToRPCGetOrganizationLogoRequest converts HTTP GetOrganizationLogoRequestDTO to RPC request.
func (a *logoAssembler) ToRPCGetOrganizationLogoRequest(
	dto *identity.GetOrganizationLogoRequestDTO,
) *identity_srv.GetOrganizationLogoRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.GetOrganizationLogoRequest{
		LogoID: dto.LogoID,
	}

	return req
}

// ToRPCDeleteOrganizationLogoRequest converts HTTP DeleteOrganizationLogoRequestDTO to RPC request.
func (a *logoAssembler) ToRPCDeleteOrganizationLogoRequest(
	dto *identity.DeleteOrganizationLogoRequestDTO,
) *identity_srv.DeleteOrganizationLogoRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.DeleteOrganizationLogoRequest{
		LogoID: dto.LogoID,
	}

	return req
}

// ToRPCBindLogoToOrganizationRequest converts HTTP BindLogoToOrganizationRequestDTO to RPC request.
func (a *logoAssembler) ToRPCBindLogoToOrganizationRequest(
	dto *identity.BindLogoToOrganizationRequestDTO,
) *identity_srv.BindLogoToOrganizationRequest {
	if dto == nil {
		return nil
	}

	req := &identity_srv.BindLogoToOrganizationRequest{
		LogoID:         dto.LogoID,
		OrganizationID: dto.OrganizationID,
	}

	return req
}
