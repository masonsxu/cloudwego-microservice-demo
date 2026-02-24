/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
export type GetUsersByRoleResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 角色ID
     */
    role_id?: string;
    /**
     * * 该角色下所有用户的ID列表
     */
    user_ids?: Array<string>;
};

