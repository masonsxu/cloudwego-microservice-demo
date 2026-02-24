/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type OrganizationDTO = {
    /**
     * * 认证状态
     */
    accreditation_status?: string;
    /**
     * * 子组织列表
     */
    children?: Array<OrganizationDTO>;
    /**
     * * 组织代码
     */
    code?: string;
    /**
     * * 创建时间
     */
    created_at?: number;
    /**
     * * 部门数量
     */
    department_count?: number;
    /**
     * * 机构类型
     */
    facility_type?: string;
    /**
     * * 组织唯一标识符
     */
    id?: string;
    /**
     * * 组织Logo地址
     */
    logo?: string;
    /**
     * * 绑定的Logo ID
     */
    logo_id?: string;
    /**
     * * 成员数量
     */
    member_count?: number;
    /**
     * * 组织名称
     */
    name?: string;
    /**
     * 关联信息
     * * 父组织信息
     */
    parent?: OrganizationDTO;
    /**
     * * 父组织ID
     */
    parent_id?: string;
    /**
     * * 所在省市列表
     */
    province_city?: Array<string>;
    /**
     * * 更新时间
     */
    updated_at?: number;
};

