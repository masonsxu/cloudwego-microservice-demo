/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type CreateOrganizationRequestDTO = {
    /**
     * * 认证状态
     */
    accreditation_status?: string;
    /**
     * * 机构类型
     */
    facility_type?: string;
    /**
     * * 组织名称
     */
    name?: string;
    /**
     * * 父组织ID（可选）
     */
    parent_id?: string;
    /**
     * * 所在省市列表
     */
    province_city?: Array<string>;
};

