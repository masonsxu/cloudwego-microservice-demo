/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { OrganizationDTO } from './OrganizationDTO';
import type { PageResponseDTO } from './PageResponseDTO';
export type ListOrganizationsResponseDTO = {
    /**
     * * 基础响应信息
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 组织列表
     */
    organizations?: Array<OrganizationDTO>;
    /**
     * * 分页信息
     */
    page?: PageResponseDTO;
};

