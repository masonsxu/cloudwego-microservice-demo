/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type UserRoleAssignmentDTO = {
    /**
     * * 创建时间
     */
    created_at?: number;
    /**
     * * 创建者用户ID
     */
    created_by?: string;
    /**
     * * 分配记录的唯一ID
     */
    id?: string;
    /**
     * * 分配的角色ID (对应 RoleDefinition.id)
     */
    role_id?: string;
    /**
     * * 最后更新时间
     */
    updated_at?: number;
    /**
     * * 更新者用户ID
     */
    updated_by?: string;
    /**
     * * 被分配角色的用户ID
     */
    user_id?: string;
};

