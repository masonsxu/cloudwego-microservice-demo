/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type CreateUserRequestDTO = {
    /**
     * * 账户过期时间
     */
    account_expiry?: number;
    /**
     * * 邮箱地址
     */
    email?: string;
    /**
     * * 员工工号
     */
    employee_id?: string;
    /**
     * * 名
     */
    first_name?: string;
    /**
     * * 性别（非必填，0:未知, 1:男, 2:女）
     */
    gender?: number;
    /**
     * * 姓
     */
    last_name?: string;
    /**
     * * 执业证书号
     */
    license_number?: string;
    /**
     * * 是否必须在下次登录时修改密码
     */
    must_change_password?: boolean;
    /**
     * * 组织ID
     */
    organization_id?: string;
    /**
     * * 密码
     */
    password?: string;
    /**
     * * 手机号码
     */
    phone?: string;
    /**
     * * 职业头衔
     */
    professional_title?: string;
    /**
     * * 真实姓名
     */
    real_name?: string;
    /**
     * * 角色ID列表
     */
    role_ids?: Array<string>;
    /**
     * * 专业特长列表
     */
    specialties?: Array<string>;
    /**
     * * 用户名
     */
    username?: string;
};

