package permission

import (
	permissionModel "github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/permission"
	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/assembler/common"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/identity_srv"
)

// =============================================================================
// 枚举类型转换辅助函数
// =============================================================================

// rpcToHTTPPermissionLevel 将 RPC 枚举转换为 HTTP 枚举
func rpcToHTTPPermissionLevel(rpc *identity_srv.PermissionLevel) *permissionModel.PermissionLevel {
	if rpc == nil {
		return nil
	}

	httpLevel := permissionModel.PermissionLevel(*rpc)

	return &httpLevel
}

// httpToRPCPermissionLevel 将 HTTP 枚举转换为 RPC 枚举
func httpToRPCPermissionLevel(http *permissionModel.PermissionLevel) *identity_srv.PermissionLevel {
	if http == nil {
		return nil
	}

	rpcLevel := identity_srv.PermissionLevel(*http)

	return &rpcLevel
}

// menuAssembler implements the IMenuAssembler interface.
type menuAssembler struct{}

// NewMenuAssembler creates a new menu assembler.
func NewMenuAssembler() IMenuAssembler {
	return &menuAssembler{}
}

// ToHTTPMenuNode converts RPC MenuNode to HTTP MenuNodeDTO.
func (a *menuAssembler) ToHTTPMenuNode(
	rpcNode *identity_srv.MenuNode,
) *permissionModel.MenuNodeDTO {
	if rpcNode == nil {
		return nil
	}

	httpNode := &permissionModel.MenuNodeDTO{
		Name:            rpcNode.Name,
		ID:              rpcNode.Id,
		Path:            rpcNode.Path,
		Icon:            common.CopyStringPtr(rpcNode.Icon),
		Component:       common.CopyStringPtr(rpcNode.Component),
		HasPermission:   common.CopyBoolPtr(rpcNode.HasPermission),
		PermissionLevel: rpcToHTTPPermissionLevel(rpcNode.PermissionLevel),
	}

	// 递归处理子菜单
	if len(rpcNode.Children) > 0 {
		httpNode.Children = make([]*permissionModel.MenuNodeDTO, 0, len(rpcNode.Children))
		for _, rpcChild := range rpcNode.Children {
			if httpChild := a.ToHTTPMenuNode(rpcChild); httpChild != nil {
				httpNode.Children = append(httpNode.Children, httpChild)
			}
		}
	}

	return httpNode
}

// ToHTTPMenuNodes converts a slice of RPC MenuNode to HTTP MenuNodeDTO.
func (a *menuAssembler) ToHTTPMenuNodes(
	rpcNodes []*identity_srv.MenuNode,
) []*permissionModel.MenuNodeDTO {
	if rpcNodes == nil {
		return nil
	}

	result := make([]*permissionModel.MenuNodeDTO, 0, len(rpcNodes))
	for _, rpcNode := range rpcNodes {
		if httpNode := a.ToHTTPMenuNode(rpcNode); httpNode != nil {
			result = append(result, httpNode)
		}
	}

	return result
}

// ToRPCMenuNode converts a MenuNodeDTO to RPC MenuNode.
func (a *menuAssembler) ToRPCMenuNode(dto *permissionModel.MenuNodeDTO) *identity_srv.MenuNode {
	if dto == nil {
		return nil
	}

	rpcNode := &identity_srv.MenuNode{
		Name:            dto.Name,
		Id:              dto.ID,
		Path:            dto.Path,
		Icon:            dto.Icon,
		Component:       dto.Component,
		HasPermission:   dto.HasPermission,
		PermissionLevel: httpToRPCPermissionLevel(dto.PermissionLevel),
	}

	// 递归转换子菜单
	if len(dto.Children) > 0 {
		rpcNode.Children = a.ToRPCMenuNodes(dto.Children)
	}

	return rpcNode
}

// ToRPCMenuNodes converts a slice of MenuNodeDTO to RPC MenuNodes.
func (a *menuAssembler) ToRPCMenuNodes(
	dtos []*permissionModel.MenuNodeDTO,
) []*identity_srv.MenuNode {
	if dtos == nil {
		return nil
	}

	result := make([]*identity_srv.MenuNode, 0, len(dtos))
	for _, dto := range dtos {
		if rpcNode := a.ToRPCMenuNode(dto); rpcNode != nil {
			result = append(result, rpcNode)
		}
	}

	return result
}

// ToRPCUploadMenuRequest converts UploadMenuRequestDTO to RPC UploadMenuRequest.
func (a *menuAssembler) ToRPCUploadMenuRequest(
	dto *permissionModel.UploadMenuRequestDTO,
) *identity_srv.UploadMenuRequest {
	if dto == nil {
		return nil
	}

	// 将前端上传的form表单文件([]byte)转换为YAML内容字符串
	yamlContent := ""
	if dto.MenuFile != nil {
		yamlContent = string(dto.MenuFile)
	}

	return &identity_srv.UploadMenuRequest{
		YamlContent: &yamlContent,
	}
}

// ToHTTPGetMenuTreeResponse converts RPC GetMenuTreeResponse to HTTP GetMenuTreeResponseDTO.
func (a *menuAssembler) ToHTTPGetMenuTreeResponse(
	rpcResp *identity_srv.GetMenuTreeResponse,
) *permissionModel.GetMenuTreeResponseDTO {
	if rpcResp == nil {
		return nil
	}

	httpResp := &permissionModel.GetMenuTreeResponseDTO{
		MenuTree: make([]*permissionModel.MenuNodeDTO, 0, len(rpcResp.MenuTree)),
	}

	for _, rpcNode := range rpcResp.MenuTree {
		if httpNode := a.ToHTTPMenuNode(rpcNode); httpNode != nil {
			httpResp.MenuTree = append(httpResp.MenuTree, httpNode)
		}
	}

	return httpResp
}

// =================================================================
// 菜单权限管理相关转换方法
// =================================================================

// ToHTTPMenuConfig converts RPC MenuConfig to HTTP MenuConfigDTO.
func (a *menuAssembler) ToHTTPMenuConfig(
	rpcConfig *identity_srv.MenuConfig,
) *permissionModel.MenuConfigDTO {
	if rpcConfig == nil {
		return nil
	}

	return &permissionModel.MenuConfigDTO{
		MenuID:     rpcConfig.MenuID,
		Permission: rpcToHTTPPermissionLevel(rpcConfig.Permission),
	}
}

// ToHTTPMenuConfigs converts a slice of RPC MenuConfig to HTTP MenuConfigDTO.
func (a *menuAssembler) ToHTTPMenuConfigs(
	rpcConfigs []*identity_srv.MenuConfig,
) []*permissionModel.MenuConfigDTO {
	if rpcConfigs == nil {
		return nil
	}

	result := make([]*permissionModel.MenuConfigDTO, 0, len(rpcConfigs))
	for _, rpcConfig := range rpcConfigs {
		if httpConfig := a.ToHTTPMenuConfig(rpcConfig); httpConfig != nil {
			result = append(result, httpConfig)
		}
	}

	return result
}

// ToRPCMenuConfig converts HTTP MenuConfigDTO to RPC MenuConfig.
func (a *menuAssembler) ToRPCMenuConfig(
	dto *permissionModel.MenuConfigDTO,
) *identity_srv.MenuConfig {
	if dto == nil {
		return nil
	}

	return &identity_srv.MenuConfig{
		MenuID:     dto.MenuID,
		Permission: httpToRPCPermissionLevel(dto.Permission),
	}
}

// ToRPCMenuConfigs converts a slice of HTTP MenuConfigDTO to RPC MenuConfig.
func (a *menuAssembler) ToRPCMenuConfigs(
	dtos []*permissionModel.MenuConfigDTO,
) []*identity_srv.MenuConfig {
	if dtos == nil {
		return nil
	}

	result := make([]*identity_srv.MenuConfig, 0, len(dtos))
	for _, dto := range dtos {
		if rpcConfig := a.ToRPCMenuConfig(dto); rpcConfig != nil {
			result = append(result, rpcConfig)
		}
	}

	return result
}

// ToHTTPMenuPermission converts RPC MenuPermission to HTTP MenuPermissionDTO.
func (a *menuAssembler) ToHTTPMenuPermission(
	rpcPermission *identity_srv.MenuPermission,
) *permissionModel.MenuPermissionDTO {
	if rpcPermission == nil {
		return nil
	}

	return &permissionModel.MenuPermissionDTO{
		MenuID:     rpcPermission.MenuID,
		Permission: rpcToHTTPPermissionLevel(rpcPermission.Permission),
	}
}

// ToHTTPMenuPermissions converts a slice of RPC MenuPermission to HTTP MenuPermissionDTO.
func (a *menuAssembler) ToHTTPMenuPermissions(
	rpcPermissions []*identity_srv.MenuPermission,
) []*permissionModel.MenuPermissionDTO {
	if rpcPermissions == nil {
		return nil
	}

	result := make([]*permissionModel.MenuPermissionDTO, 0, len(rpcPermissions))
	for _, rpcPermission := range rpcPermissions {
		if httpPermission := a.ToHTTPMenuPermission(rpcPermission); httpPermission != nil {
			result = append(result, httpPermission)
		}
	}

	return result
}

// ToRPCConfigureRoleMenusRequest converts HTTP request to RPC ConfigureRoleMenusRequest.
func (a *menuAssembler) ToRPCConfigureRoleMenusRequest(
	operatorID string,
	dto *permissionModel.ConfigureRoleMenusRequestDTO,
) *identity_srv.ConfigureRoleMenusRequest {
	if dto == nil {
		return nil
	}

	rpcReq := &identity_srv.ConfigureRoleMenusRequest{
		RoleID:      dto.RoleID,
		MenuConfigs: a.ToRPCMenuConfigs(dto.MenuConfigs),
	}

	if operatorID != "" {
		rpcReq.OperatorID = &operatorID
	}

	return rpcReq
}

// ToHTTPConfigureRoleMenusResponse converts RPC ConfigureRoleMenusResponse to HTTP response.
func (a *menuAssembler) ToHTTPConfigureRoleMenusResponse(
	rpcResp *identity_srv.ConfigureRoleMenusResponse,
) *permissionModel.ConfigureRoleMenusResponseDTO {
	if rpcResp == nil {
		return nil
	}

	return &permissionModel.ConfigureRoleMenusResponseDTO{
		Success: rpcResp.Success,
		Message: common.CopyStringPtr(rpcResp.Message),
	}
}

// ToRPCGetRoleMenuTreeRequest converts HTTP request to RPC GetRoleMenuTreeRequest.
func (a *menuAssembler) ToRPCGetRoleMenuTreeRequest(
	dto *permissionModel.GetRoleMenuTreeRequestDTO,
) *identity_srv.GetRoleMenuTreeRequest {
	return &identity_srv.GetRoleMenuTreeRequest{
		RoleID: dto.RoleID,
	}
}

// ToHTTPGetRoleMenuTreeResponse converts RPC GetRoleMenuTreeResponse to HTTP response.
func (a *menuAssembler) ToHTTPGetRoleMenuTreeResponse(
	rpcResp *identity_srv.GetRoleMenuTreeResponse,
) *permissionModel.GetRoleMenuTreeResponseDTO {
	if rpcResp == nil {
		return nil
	}

	httpResp := &permissionModel.GetRoleMenuTreeResponseDTO{
		MenuTree: a.ToHTTPMenuNodes(rpcResp.MenuTree),
		RoleID:   rpcResp.RoleID,
	}

	return httpResp
}

// ToRPCGetUserMenuTreeRequest converts HTTP request to RPC GetUserMenuTreeRequest.
func (a *menuAssembler) ToRPCGetUserMenuTreeRequest(
	dto *permissionModel.GetUserMenuTreeRequestDTO,
) *identity_srv.GetUserMenuTreeRequest {
	return &identity_srv.GetUserMenuTreeRequest{
		UserID: dto.UserID,
	}
}

// ToHTTPGetUserMenuTreeResponse converts RPC GetUserMenuTreeResponse to HTTP response.
func (a *menuAssembler) ToHTTPGetUserMenuTreeResponse(
	rpcResp *identity_srv.GetUserMenuTreeResponse,
) *permissionModel.GetUserMenuTreeResponseDTO {
	if rpcResp == nil {
		return nil
	}

	httpResp := &permissionModel.GetUserMenuTreeResponseDTO{
		MenuTree: a.ToHTTPMenuNodes(rpcResp.MenuTree),
		UserID:   rpcResp.UserID,
		RoleIDs:  make([]string, len(rpcResp.RoleIDs)),
	}

	// 转换角色ID列表（从core.UUID到string）
	copy(httpResp.RoleIDs, rpcResp.RoleIDs)

	// 从菜单树中提取扁平化的权限列表
	httpResp.Permissions = a.extractPermissionsFromMenuTree(rpcResp.MenuTree)

	return httpResp
}

// extractPermissionsFromMenuTree 从 MenuNode 树中提取扁平化的权限列表
// 递归遍历所有菜单节点，收集具有权限的菜单，返回 MenuPermissionDTO 列表
func (a *menuAssembler) extractPermissionsFromMenuTree(
	menuNodes []*identity_srv.MenuNode,
) []*permissionModel.MenuPermissionDTO {
	if menuNodes == nil {
		return nil
	}

	permissions := make([]*permissionModel.MenuPermissionDTO, 0)

	for _, node := range menuNodes {
		if node == nil {
			continue
		}

		// 如果节点有权限级别，添加到权限列表
		if node.PermissionLevel != nil && node.Id != nil {
			perm := &permissionModel.MenuPermissionDTO{
				MenuID:     node.Id,
				Permission: rpcToHTTPPermissionLevel(node.PermissionLevel),
			}
			permissions = append(permissions, perm)
		}

		// 递归处理子菜单
		if len(node.Children) > 0 {
			childPermissions := a.extractPermissionsFromMenuTree(node.Children)
			permissions = append(permissions, childPermissions...)
		}
	}

	return permissions
}

// ToRPCGetRoleMenuPermissionsRequest converts HTTP request to RPC GetRoleMenuPermissionsRequest.
func (a *menuAssembler) ToRPCGetRoleMenuPermissionsRequest(
	dto *permissionModel.GetRoleMenuPermissionsRequestDTO,
) *identity_srv.GetRoleMenuPermissionsRequest {
	return &identity_srv.GetRoleMenuPermissionsRequest{
		RoleID: dto.RoleID,
	}
}

// ToHTTPGetRoleMenuPermissionsResponse converts RPC GetRoleMenuPermissionsResponse to HTTP response.
func (a *menuAssembler) ToHTTPGetRoleMenuPermissionsResponse(
	rpcResp *identity_srv.GetRoleMenuPermissionsResponse,
) *permissionModel.GetRoleMenuPermissionsResponseDTO {
	if rpcResp == nil {
		return nil
	}

	httpResp := &permissionModel.GetRoleMenuPermissionsResponseDTO{
		Permissions: a.ToHTTPMenuPermissions(rpcResp.Permissions),
		RoleID:      rpcResp.RoleID,
	}

	return httpResp
}

// ToRPCHasMenuPermissionRequest converts HTTP request to RPC HasMenuPermissionRequest.
func (a *menuAssembler) ToRPCHasMenuPermissionRequest(
	dto *permissionModel.HasMenuPermissionRequestDTO,
) *identity_srv.HasMenuPermissionRequest {
	if dto == nil {
		return nil
	}

	return &identity_srv.HasMenuPermissionRequest{
		RoleID:     dto.RoleID,
		MenuID:     dto.MenuID,
		Permission: httpToRPCPermissionLevel(dto.Permission),
	}
}

// ToHTTPHasMenuPermissionResponse converts RPC HasMenuPermissionResponse to HTTP response.
func (a *menuAssembler) ToHTTPHasMenuPermissionResponse(
	rpcResp *identity_srv.HasMenuPermissionResponse,
) *permissionModel.HasMenuPermissionResponseDTO {
	if rpcResp == nil {
		return nil
	}

	httpResp := &permissionModel.HasMenuPermissionResponseDTO{
		HasPermission: rpcResp.HasPermission,
		RoleID:        rpcResp.RoleID,
		MenuID:        rpcResp.MenuID,
		Permission:    rpcToHTTPPermissionLevel(rpcResp.Permission),
	}

	return httpResp
}
