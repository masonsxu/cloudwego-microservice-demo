/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type UserProfileDTO = {
    /**
     * * 账户过期时间
     */
    account_expiry?: number;
    /**
     * * 创建时间
     */
    created_at?: number;
    /**
     * * 创建者用户ID
     */
    created_by?: string;
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
     * * 用户唯一标识符
     */
    id?: string;
    /**
     * * 最后登录时间
     */
    last_login_time?: number;
    /**
     * * 姓
     */
    last_name?: string;
    /**
     * * 执业证书号
     */
    license_number?: string;
    /**
     * * 连续登录失败次数
     */
    login_attempts?: number;
    /**
     * * 是否必须修改密码
     */
    must_change_password?: boolean;
    /**
     * * 手机号码
     */
    phone?: string;
    /**
     * * 主部门ID
     */
    primary_department_id?: string;
    /**
     * * 主组织ID
     */
    primary_organization_id?: string;
    /**
     * * 职业头衔
     */
    professional_title?: string;
    /**
     * * 真实姓名
     */
    real_name?: string;
    /**
     * * 用户角色ID列表
     */
    role_ids?: Array<string>;
    /**
     * * 专业特长列表
     */
    specialties?: Array<string>;
    /**
     * * 用户状态
     */
    status?: number;
    /**
     * * 更新时间
     */
    updated_at?: number;
    /**
     * * 最后更新者用户ID
     */
    updated_by?: string;
    /**
     * * 用户名（登录名）
     */
    username?: string;
};

