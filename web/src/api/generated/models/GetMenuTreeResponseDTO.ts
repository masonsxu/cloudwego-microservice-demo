/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { MenuNodeDTO } from './MenuNodeDTO';
export type GetMenuTreeResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 完整的菜单树结构
     */
    menu_tree?: Array<MenuNodeDTO>;
};

