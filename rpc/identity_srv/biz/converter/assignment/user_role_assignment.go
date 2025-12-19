package assignment

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter defines the interface for converting user role assignment data.
type Converter interface {
	ModelToThrift(*models.UserRoleAssignment) *identity_srv.UserRoleAssignment
}
