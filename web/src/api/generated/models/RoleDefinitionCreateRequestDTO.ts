/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PermissionDTO } from './PermissionDTO';
export type RoleDefinitionCreateRequestDTO = {
    /**
     * * 角色描述
     */
    description?: string;
    /**
     * * 是否为系统内置角色
     */
    is_system_role?: boolean;
    /**
     * * 角色唯一名称 (英文，用于程序识别)
     */
    name?: string;
    /**
     * * 角色包含的权限列表
     */
    permissions?: Array<PermissionDTO>;
};

