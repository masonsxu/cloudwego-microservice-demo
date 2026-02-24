/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { PageResponseDTO } from './PageResponseDTO';
import type { UserProfileDTO } from './UserProfileDTO';
export type ListUsersResponseDTO = {
    /**
     * * 基础响应信息
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 分页信息
     */
    page?: PageResponseDTO;
    /**
     * * 用户列表
     */
    users?: Array<UserProfileDTO>;
};

