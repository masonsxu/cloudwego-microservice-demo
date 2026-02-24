/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PermissionDTO } from './PermissionDTO';
export type RoleDefinitionDTO = {
    /**
     * * 创建时间
     */
    created_at?: number;
    /**
     * * 创建者用户ID
     */
    created_by?: string;
    /**
     * * 角色详细描述
     */
    description?: string;
    /**
     * * 角色唯一ID
     */
    id?: string;
    /**
     * * 是否为系统内置角色，不可删除
     */
    is_system_role?: boolean;
    /**
     * * 角色唯一名称 (中文，页面展示使用)
     */
    name?: string;
    /**
     * * 该角色拥有的权限列表
     */
    permissions?: Array<PermissionDTO>;
    /**
     * * 角色状态
     */
    status?: number;
    /**
     * * 最后更新时间
     */
    updated_at?: number;
    /**
     * * 更新者用户ID
     */
    updated_by?: string;
    /**
     * * 当前角色绑定的用户数量（非持久化字段，查询时动态计算）
     */
    user_count?: number;
};

