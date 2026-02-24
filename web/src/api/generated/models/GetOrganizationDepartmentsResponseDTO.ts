/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { DepartmentDTO } from './DepartmentDTO';
import type { PageResponseDTO } from './PageResponseDTO';
export type GetOrganizationDepartmentsResponseDTO = {
    /**
     * * 基础响应信息
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 部门列表
     */
    departments?: Array<DepartmentDTO>;
    /**
     * * 分页信息
     */
    page?: PageResponseDTO;
};

