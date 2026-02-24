/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { PageResponseDTO } from './PageResponseDTO';
import type { RoleDefinitionDTO } from './RoleDefinitionDTO';
export type RoleDefinitionListResponseDTO = {
    /**
     * * 响应状态码
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 分页响应参数
     */
    page?: PageResponseDTO;
    /**
     * * 角色定义列表
     */
    roles?: Array<RoleDefinitionDTO>;
};

