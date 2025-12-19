package permission

import (
	permissionModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// permissionAssembler 权限转换器实现
type permissionAssembler struct{}

// NewPermissionAssembler 创建权限转换器
func NewPermissionAssembler() IPermissionAssembler {
	return &permissionAssembler{}
}

// ToHTTPPermission 将RPC权限转换为HTTP权限DTO
func (a *permissionAssembler) ToHTTPPermission(
	rpc *identity_srv.Permission,
) *permissionModel.PermissionDTO {
	if rpc == nil {
		return nil
	}

	return &permissionModel.PermissionDTO{
		Resource:    rpc.Resource,
		Action:      rpc.Action,
		Description: common.CopyStringPtr(rpc.Description),
	}
}

// ToHTTPPermissions 将RPC权限列表转换为HTTP权限DTO列表
func (a *permissionAssembler) ToHTTPPermissions(
	rpc []*identity_srv.Permission,
) []*permissionModel.PermissionDTO {
	if rpc == nil {
		return nil
	}

	result := make([]*permissionModel.PermissionDTO, len(rpc))
	for i, p := range rpc {
		result[i] = a.ToHTTPPermission(p)
	}

	return result
}

// ToRPCPermission 将HTTP权限DTO转换为RPC权限
func (a *permissionAssembler) ToRPCPermission(
	http *permissionModel.PermissionDTO,
) *identity_srv.Permission {
	if http == nil {
		return nil
	}

	return &identity_srv.Permission{
		Resource:    http.Resource,
		Action:      http.Action,
		Description: http.Description,
	}
}

// ToRPCPermissions 将HTTP权限DTO列表转换为RPC权限列表
func (a *permissionAssembler) ToRPCPermissions(
	http []*permissionModel.PermissionDTO,
) []*identity_srv.Permission {
	if http == nil {
		return nil
	}

	result := make([]*identity_srv.Permission, len(http))
	for i, p := range http {
		result[i] = a.ToRPCPermission(p)
	}

	return result
}
