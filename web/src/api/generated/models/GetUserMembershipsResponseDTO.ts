/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { PageResponseDTO } from './PageResponseDTO';
import type { UserMembershipDTO } from './UserMembershipDTO';
export type GetUserMembershipsResponseDTO = {
    /**
     * * 基础响应信息
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 成员关系列表
     */
    memberships?: Array<UserMembershipDTO>;
    /**
     * * 分页信息
     */
    page?: PageResponseDTO;
};

