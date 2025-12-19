package logo

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter 组织Logo转换器接口
type Converter interface {
	// Model -> Thrift 转换
	ModelToThrift(*models.OrganizationLogo) *identity_srv.OrganizationLogo

	// Request -> Model 转换
	UploadRequestToModel(*identity_srv.UploadTemporaryLogoRequest) *models.OrganizationLogo
	BindRequestToModel(
		*identity_srv.BindLogoToOrganizationRequest,
	) (*models.OrganizationLogo, error)

	// Status 转换
	StatusModelToThrift(models.OrganizationLogoStatus) identity_srv.OrganizationLogoStatus
	StatusThriftToModel(identity_srv.OrganizationLogoStatus) models.OrganizationLogoStatus
}
