package permission

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter defines the interface for converting identity-related data.
type Converter interface {
	ModelToThrift(*models.Permission) *identity_srv.Permission
	ThriftToModel(*identity_srv.Permission) *models.Permission
	ModelSliceToThrift([]*models.Permission) []*identity_srv.Permission
	ThriftSliceToModel([]*identity_srv.Permission) []*models.Permission
}
