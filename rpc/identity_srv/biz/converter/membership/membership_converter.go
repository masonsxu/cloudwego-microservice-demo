package membership

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter is an interface for converting membership-related data.
type Converter interface {
	ModelUserMembershipToThrift(*models.UserMembership) *identity_srv.UserMembership
	ThriftUserMembershipToModel(*identity_srv.UserMembership) *models.UserMembership
	ModelUserMembershipsToThrift([]*models.UserMembership) []*identity_srv.UserMembership

	// Alias for ModelUserMembershipToThrift
	ModelToThrift(*models.UserMembership) *identity_srv.UserMembership

	AddMembershipRequestToModel(*identity_srv.AddMembershipRequest) *models.UserMembership
	ApplyUpdateToModel(
		*models.UserMembership,
		*identity_srv.UpdateMembershipRequest,
	) *models.UserMembership
}
