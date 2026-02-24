/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { MenuNodeDTO } from './MenuNodeDTO';
export type GetRoleMenuTreeResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 角色可访问的菜单树
     */
    menu_tree?: Array<MenuNodeDTO>;
    /**
     * * 角色ID
     */
    role_id?: string;
};

