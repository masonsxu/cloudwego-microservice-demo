package identity

import (
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/identity"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
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
		// PermissionInfo 将由上层服务设置
		// TokenInfo 将由上层服务设置
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
		OldPassword:  dto.OldPassword,
		NewPassword_: dto.NewPassword,
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
		UserID:       dto.UserID,
		NewPassword_: dto.NewPassword,
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

// ToHTTPMenuTree converts RPC MenuNode array to HTTP MenuNodeDTO array
func (a *authAssembler) ToHTTPMenuTree(
	rpcMenuNodes []*identity_srv.MenuNode,
) []*permission.MenuNodeDTO {
	if rpcMenuNodes == nil {
		return nil
	}

	result := make([]*permission.MenuNodeDTO, 0, len(rpcMenuNodes))
	for _, rpcNode := range rpcMenuNodes {
		if rpcNode != nil {
			httpNode := a.convertMenuNode(rpcNode)
			if httpNode != nil {
				result = append(result, httpNode)
			}
		}
	}

	return result
}

// convertMenuNode converts a single RPC MenuNode to HTTP MenuNodeDTO
func (a *authAssembler) convertMenuNode(rpcNode *identity_srv.MenuNode) *permission.MenuNodeDTO {
	if rpcNode == nil {
		return nil
	}

	httpNode := &permission.MenuNodeDTO{
		Name:      common.CopyStringPtr(rpcNode.Name),
		ID:        common.CopyStringPtr(rpcNode.Id),
		Path:      common.CopyStringPtr(rpcNode.Path),
		Icon:      common.CopyStringPtr(rpcNode.Icon),
		Component: common.CopyStringPtr(rpcNode.Component),
	}

	// 递归处理子菜单
	if len(rpcNode.Children) > 0 {
		httpNode.Children = make([]*permission.MenuNodeDTO, 0, len(rpcNode.Children))
		for _, childRpcNode := range rpcNode.Children {
			if childRpcNode != nil {
				childHttpNode := a.convertMenuNode(childRpcNode)
				if childHttpNode != nil {
					httpNode.Children = append(httpNode.Children, childHttpNode)
				}
			}
		}
	}

	return httpNode
}

// ToHTTPMenuPermissions converts RPC MenuPermission array to HTTP MenuPermissionDTO array
func (a *authAssembler) ToHTTPMenuPermissions(
	rpcPermissions []*identity_srv.MenuPermission,
) []*permission.MenuPermissionDTO {
	if rpcPermissions == nil {
		return nil
	}

	result := make([]*permission.MenuPermissionDTO, 0, len(rpcPermissions))
	for _, rpcPerm := range rpcPermissions {
		if rpcPerm != nil {
			httpPerm := &permission.MenuPermissionDTO{
				MenuID: common.CopyStringPtr(rpcPerm.MenuID),
			}
			if rpcPerm.Permission != nil {
				permLevel := permission.PermissionLevel(*rpcPerm.Permission)
				httpPerm.Permission = &permLevel
			}

			result = append(result, httpPerm)
		}
	}

	return result
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
				ID:   common.CopyStringPtr(rpcRole.Id),
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
