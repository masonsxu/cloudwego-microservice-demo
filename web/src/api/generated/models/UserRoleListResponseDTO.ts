/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { PageResponseDTO } from './PageResponseDTO';
import type { UserRoleAssignmentDTO } from './UserRoleAssignmentDTO';
export type UserRoleListResponseDTO = {
    /**
     * * 用户角色分配列表
     */
    assignments?: Array<UserRoleAssignmentDTO>;
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 分页响应参数
     */
    page?: PageResponseDTO;
};

