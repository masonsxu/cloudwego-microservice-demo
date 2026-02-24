/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { MenuNodeDTO } from './MenuNodeDTO';
import type { MenuPermissionDTO } from './MenuPermissionDTO';
export type GetUserMenuTreeResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 用户可访问的菜单树
     */
    menu_tree?: Array<MenuNodeDTO>;
    /**
     * * 用户菜单权限列表（扁平化格式，用于前端按钮级权限控制）
     */
    permissions?: Array<MenuPermissionDTO>;
    /**
     * * 用户拥有的角色列表
     */
    role_ids?: Array<string>;
    /**
     * * 用户ID
     */
    user_id?: string;
};

