/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type RoleInfoDTO = {
    /**
     * * 角色编码（如：ROLE_ADMIN, ROLE_DOCTOR）
     */
    code?: string;
    /**
     * * 数据范围（1:仅自己, 2:本部门, 3:全组织）
     */
    data_scope?: number;
    /**
     * * 角色唯一标识符
     */
    id?: string;
    /**
     * * 角色名称（中文显示名称）
     */
    name?: string;
};

