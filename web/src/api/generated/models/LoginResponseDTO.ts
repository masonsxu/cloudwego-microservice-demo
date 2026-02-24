/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseResponseDTO } from './BaseResponseDTO';
import type { MenuNodeDTO } from './MenuNodeDTO';
import type { MenuPermissionDTO } from './MenuPermissionDTO';
import type { RoleInfoDTO } from './RoleInfoDTO';
import type { TokenInfoDTO } from './TokenInfoDTO';
import type { UserMembershipDTO } from './UserMembershipDTO';
import type { UserProfileDTO } from './UserProfileDTO';
export type LoginResponseDTO = {
    /**
     * * 基础响应信息
     */
    base_resp?: BaseResponseDTO;
    /**
     * * 用户成员关系列表
     */
    memberships?: Array<UserMembershipDTO>;
    /**
     * * 用户菜单树（新增字段，推荐使用）
     */
    menu_tree?: Array<MenuNodeDTO>;
    /**
     * * 用户菜单权限列表（用于前端按钮级权限控制）
     */
    permissions?: Array<MenuPermissionDTO>;
    /**
     * * 用户角色ID列表
     */
    role_ids?: Array<string>;
    /**
     * * 用户角色详情列表（包含角色编码、名称和数据范围）
     */
    roles?: Array<RoleInfoDTO>;
    /**
     * * 访问令牌信息
     */
    token_info?: TokenInfoDTO;
    /**
     * * 用户个人信息
     */
    user_profile?: UserProfileDTO;
};

