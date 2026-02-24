/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PermissionDTO } from './PermissionDTO';
export type RoleDefinitionUpdateRequestDTO = {
    /**
     * * 角色描述
     */
    description?: string;
    /**
     * * 角色名称
     */
    name?: string;
    /**
     * * 权限列表
     */
    permissions?: Array<PermissionDTO>;
    /**
     * * 角色状态
     */
    status?: number;
};

