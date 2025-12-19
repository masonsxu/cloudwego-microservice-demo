package logo

import (
	"context"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// LogoLogic 组织Logo业务逻辑接口
type LogoLogic interface {
	// UploadTemporaryLogo 上传临时Logo（7天过期）
	UploadTemporaryLogo(
		ctx context.Context,
		req *identity_srv.UploadTemporaryLogoRequest,
	) (*identity_srv.OrganizationLogo, error)

	// GetOrganizationLogo 获取Logo详情
	GetOrganizationLogo(
		ctx context.Context,
		req *identity_srv.GetOrganizationLogoRequest,
	) (*identity_srv.OrganizationLogo, error)

	// DeleteOrganizationLogo 删除Logo（软删除+S3文件删除）
	DeleteOrganizationLogo(
		ctx context.Context,
		req *identity_srv.DeleteOrganizationLogoRequest,
	) error

	// BindLogoToOrganization 绑定Logo到组织（临时→永久）
	BindLogoToOrganization(
		ctx context.Context,
		req *identity_srv.BindLogoToOrganizationRequest,
	) (*identity_srv.OrganizationLogo, error)
}
