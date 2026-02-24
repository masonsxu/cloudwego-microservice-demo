/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
export type BatchBindUsersToRoleResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 响应消息
     */
    message?: string;
    /**
     * * 操作是否成功
     */
    success?: boolean;
    /**
     * * 成功绑定的用户数量
     */
    success_count?: number;
};

