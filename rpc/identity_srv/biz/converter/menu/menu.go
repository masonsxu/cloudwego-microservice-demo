package menu

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter defines the interface for converting Menu models to thrift objects.
type Converter interface {
	// ModelToThrift converts a menu model to its thrift representation, including its children recursively.
	ModelToThrift(model *models.Menu) *identity_srv.MenuNode

	// ModelsToThrift converts a slice of menu models to a slice of thrift representations.
	ModelsToThrift(models []*models.Menu) []*identity_srv.MenuNode
}
