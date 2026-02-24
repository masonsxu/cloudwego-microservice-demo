/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { MenuPermissionDTO } from './MenuPermissionDTO';
export type GetRoleMenuPermissionsResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 角色的菜单权限列表
     */
    permissions?: Array<MenuPermissionDTO>;
    /**
     * * 角色ID
     */
    role_id?: string;
};

