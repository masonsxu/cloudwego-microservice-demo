package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// Auth Assembler
type authAssembler struct{}

func NewAuthAssembler() IAuthAssembler {
	return &authAssembler{}
}

func (a *authAssembler) ToRPCLoginRequest(
	dto *identity.LoginRequestDTO,
) *identity_srv.LoginRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.LoginRequest{
		Username: dto.Username,
		Password: dto.Password,
	}
}

// ToHTTPLoginResponse converts an RPC LoginResponse to an HTTP LoginResponseDTO.
func (a *authAssembler) ToHTTPLoginResponse(
	rpc *identity_srv.LoginResponse,
) *identity.LoginResponseDTO {
	if rpc == nil {
		return nil
	}

	return &identity.LoginResponseDTO{
		UserProfile: NewUserAssembler().ToHTTPUserProfile(rpc.UserProfile),
		// TokenInfo 由 jwt_middleware 在签发后回填
	}
}

// ToRPCChangePasswordRequest converts an HTTP ChangePasswordRequestDTO to an RPC ChangePasswordRequest.
func (a *authAssembler) ToRPCChangePasswordRequest(
	dto *identity.ChangePasswordRequestDTO,
) *identity_srv.ChangePasswordRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.ChangePasswordRequest{
		OldPassword: dto.OldPassword,
		NewPassword: dto.NewPassword,
	}
}

// ToRPCResetPasswordRequest converts an HTTP ResetPasswordRequestDTO to an RPC ResetPasswordRequest.
func (a *authAssembler) ToRPCResetPasswordRequest(
	dto *identity.ResetPasswordRequestDTO,
) *identity_srv.ResetPasswordRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.ResetPasswordRequest{
		UserID:      dto.UserID,
		NewPassword: dto.NewPassword,
	}
}

// ToRPCForcePasswordChangeRequest converts an HTTP ForcePasswordChangeRequestDTO to an RPC ForcePasswordChangeRequest.
func (a *authAssembler) ToRPCForcePasswordChangeRequest(
	dto *identity.ForcePasswordChangeRequestDTO,
) *identity_srv.ForcePasswordChangeRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.ForcePasswordChangeRequest{
		UserID: dto.UserID,
	}
}

// ToHTTPRoleInfos converts RPC RoleDefinition array to HTTP RoleInfoDTO array
func (a *authAssembler) ToHTTPRoleInfos(
	rpcRoles []*identity_srv.RoleDefinition,
) []*identity.RoleInfoDTO {
	if rpcRoles == nil {
		return nil
	}

	result := make([]*identity.RoleInfoDTO, 0, len(rpcRoles))
	for _, rpcRole := range rpcRoles {
		if rpcRole != nil {
			httpRole := &identity.RoleInfoDTO{
				Id:   common.CopyStringPtr(rpcRole.Id),
				Code: common.CopyStringPtr(rpcRole.RoleCode),
				Name: common.CopyStringPtr(rpcRole.Name),
			}
			if rpcRole.DefaultScope != nil {
				dataScope := int32(*rpcRole.DefaultScope)
				httpRole.DataScope = &dataScope
			}

			result = append(result, httpRole)
		}
	}

	return result
}
