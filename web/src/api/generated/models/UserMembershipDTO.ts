/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { DepartmentDTO } from './DepartmentDTO';
import type { OrganizationDTO } from './OrganizationDTO';
export type UserMembershipDTO = {
    /**
     * * 创建时间
     */
    created_at?: number;
    /**
     * * 所属部门信息
     */
    department?: DepartmentDTO;
    /**
     * * 部门ID（可选）
     */
    department_id?: string;
    /**
     * * 成员关系唯一标识符
     */
    id?: string;
    /**
     * * 是否为主要成员关系
     */
    is_primary?: boolean;
    /**
     * 关联信息
     * * 所属组织信息
     */
    organization?: OrganizationDTO;
    /**
     * * 组织ID
     */
    organization_id?: string;
    /**
     * * 更新时间
     */
    updated_at?: number;
    /**
     * * 用户ID
     */
    user_id?: string;
};

