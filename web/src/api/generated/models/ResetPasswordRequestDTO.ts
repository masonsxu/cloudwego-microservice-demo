/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type ResetPasswordRequestDTO = {
    /**
     * * 新密码（可选，为空时系统生成）
     */
    new_password?: string;
    /**
     * * 重置原因
     */
    reset_reason?: string;
    /**
     * * 目标用户ID
     */
    user_id?: string;
};

