/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { OrganizationDTO } from './OrganizationDTO';
export type DepartmentDTO = {
    /**
     * * 可用设备列表
     */
    available_equipment?: Array<string>;
    /**
     * * 部门代码
     */
    code?: string;
    /**
     * * 创建时间
     */
    created_at?: number;
    /**
     * * 部门类型
     */
    department_type?: string;
    /**
     * * 部门唯一标识符
     */
    id?: string;
    /**
     * * 成员数量
     */
    member_count?: number;
    /**
     * * 部门名称
     */
    name?: string;
    /**
     * 关联信息
     * * 所属组织信息
     */
    organization?: OrganizationDTO;
    /**
     * * 所属组织ID
     */
    organization_id?: string;
    /**
     * * 更新时间
     */
    updated_at?: number;
};

