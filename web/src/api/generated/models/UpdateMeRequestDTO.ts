/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type UpdateMeRequestDTO = {
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
     * * 性别
     */
    gender?: number;
    /**
     * * 姓
     */
    last_name?: string;
    /**
     * * 许可证号
     */
    medical_license_number?: string;
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
     * * 专业特长列表
     */
    specialties?: Array<string>;
};

