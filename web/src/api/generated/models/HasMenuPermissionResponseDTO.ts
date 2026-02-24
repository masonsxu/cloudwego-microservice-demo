/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { PermissionLevel } from './PermissionLevel';
export type HasMenuPermissionResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 是否具有权限
     */
    has_permission?: boolean;
    /**
     * * 菜单ID
     */
    menu_id?: string;
    /**
     * * 权限级别
     */
    permission?: PermissionLevel;
    /**
     * * 角色ID
     */
    role_id?: string;
};

