/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PermissionLevel } from './PermissionLevel';
export type MenuNodeDTO = {
    /**
     * * 子菜单列表 (可选)
     */
    children?: Array<MenuNodeDTO>;
    /**
     * * 前端组件路径 (可选)
     */
    component?: string;
    /**
     * * 是否有权限访问此菜单 (可选, 用于权限标记)
     */
    has_permission?: boolean;
    /**
     * * 菜单图标 (可选)
     */
    icon?: string;
    /**
     * * 菜单唯一标识符
     */
    id?: string;
    /**
     * * 菜单名称 (用于显示)
     */
    name?: string;
    /**
     * * 路由路径
     */
    path?: string;
    /**
     * * 权限级别
     */
    permission_level?: PermissionLevel;
};

