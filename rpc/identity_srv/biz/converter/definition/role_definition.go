package definition

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter defines the interface for converting role definition data.
type Converter interface {
	ModelToThrift(*models.RoleDefinition) *identity_srv.RoleDefinition
	ModelsToThrift([]*models.RoleDefinition) []*identity_srv.RoleDefinition
}
