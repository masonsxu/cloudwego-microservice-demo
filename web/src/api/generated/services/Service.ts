/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { AssignRoleToUserResponseDTO } from '../models/AssignRoleToUserResponseDTO';
import type { BatchBindUsersToRoleRequestDTO } from '../models/BatchBindUsersToRoleRequestDTO';
import type { BatchBindUsersToRoleResponseDTO } from '../models/BatchBindUsersToRoleResponseDTO';
import type { BindLogoToOrganizationRequestDTO } from '../models/BindLogoToOrganizationRequestDTO';
import type { ChangePasswordRequestDTO } from '../models/ChangePasswordRequestDTO';
import type { ChangeUserStatusRequestDTO } from '../models/ChangeUserStatusRequestDTO';
import type { ConfigureRoleMenusRequestDTO } from '../models/ConfigureRoleMenusRequestDTO';
import type { ConfigureRoleMenusResponseDTO } from '../models/ConfigureRoleMenusResponseDTO';
import type { CreateDepartmentRequestDTO } from '../models/CreateDepartmentRequestDTO';
import type { CreateOrganizationRequestDTO } from '../models/CreateOrganizationRequestDTO';
import type { CreateUserRequestDTO } from '../models/CreateUserRequestDTO';
import type { DeleteDepartmentRequestDTO } from '../models/DeleteDepartmentRequestDTO';
import type { DeleteOrganizationRequestDTO } from '../models/DeleteOrganizationRequestDTO';
import type { DeleteUserRequestDTO } from '../models/DeleteUserRequestDTO';
import type { DepartmentResponseDTO } from '../models/DepartmentResponseDTO';
import type { ForcePasswordChangeRequestDTO } from '../models/ForcePasswordChangeRequestDTO';
import type { GetMenuTreeResponseDTO } from '../models/GetMenuTreeResponseDTO';
import type { GetOrganizationDepartmentsResponseDTO } from '../models/GetOrganizationDepartmentsResponseDTO';
import type { GetPrimaryMembershipRequestDTO } from '../models/GetPrimaryMembershipRequestDTO';
import type { GetRoleMenuPermissionsResponseDTO } from '../models/GetRoleMenuPermissionsResponseDTO';
import type { GetRoleMenuTreeResponseDTO } from '../models/GetRoleMenuTreeResponseDTO';
import type { GetUserMembershipsResponseDTO } from '../models/GetUserMembershipsResponseDTO';
import type { GetUserMenuTreeResponseDTO } from '../models/GetUserMenuTreeResponseDTO';
import type { GetUsersByRoleResponseDTO } from '../models/GetUsersByRoleResponseDTO';
import type { HasMenuPermissionRequestDTO } from '../models/HasMenuPermissionRequestDTO';
import type { HasMenuPermissionResponseDTO } from '../models/HasMenuPermissionResponseDTO';
import type { ListOrganizationsResponseDTO } from '../models/ListOrganizationsResponseDTO';
import type { ListUsersResponseDTO } from '../models/ListUsersResponseDTO';
import type { LoginRequestDTO } from '../models/LoginRequestDTO';
import type { LoginResponseDTO } from '../models/LoginResponseDTO';
import type { LogoutRequestDTO } from '../models/LogoutRequestDTO';
import type { OperationStatusResponseDTO } from '../models/OperationStatusResponseDTO';
import type { OrganizationLogoResponseDTO } from '../models/OrganizationLogoResponseDTO';
import type { OrganizationResponseDTO } from '../models/OrganizationResponseDTO';
import type { RefreshTokenRequestDTO } from '../models/RefreshTokenRequestDTO';
import type { RefreshTokenResponseDTO } from '../models/RefreshTokenResponseDTO';
import type { ResetPasswordRequestDTO } from '../models/ResetPasswordRequestDTO';
import type { RoleDefinitionCreateRequestDTO } from '../models/RoleDefinitionCreateRequestDTO';
import type { RoleDefinitionCreateResponseDTO } from '../models/RoleDefinitionCreateResponseDTO';
import type { RoleDefinitionGetResponseDTO } from '../models/RoleDefinitionGetResponseDTO';
import type { RoleDefinitionListResponseDTO } from '../models/RoleDefinitionListResponseDTO';
import type { RoleDefinitionUpdateRequestDTO } from '../models/RoleDefinitionUpdateRequestDTO';
import type { RoleDefinitionUpdateResponseDTO } from '../models/RoleDefinitionUpdateResponseDTO';
import type { SearchUsersResponseDTO } from '../models/SearchUsersResponseDTO';
import type { UnlockUserRequestDTO } from '../models/UnlockUserRequestDTO';
import type { UpdateDepartmentRequestDTO } from '../models/UpdateDepartmentRequestDTO';
import type { UpdateMeRequestDTO } from '../models/UpdateMeRequestDTO';
import type { UpdateOrganizationRequestDTO } from '../models/UpdateOrganizationRequestDTO';
import type { UpdateUserRequestDTO } from '../models/UpdateUserRequestDTO';
import type { UserMembershipResponseDTO } from '../models/UserMembershipResponseDTO';
import type { UserProfileResponseDTO } from '../models/UserProfileResponseDTO';
import type { UserRoleListResponseDTO } from '../models/UserRoleListResponseDTO';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class Service {
    /**
     * 用户登录
     * 验证用户凭据并返回访问令牌和用户信息
     * @param req 请求体
     * @returns LoginResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityAuthLogin(
        req: LoginRequestDTO,
    ): CancelablePromise<LoginResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/auth/login',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 用户登出
     * 注销当前会话并使令牌失效
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityAuthLogout(
        req: LogoutRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/auth/logout',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 修改密码
     * 用户修改自己的密码（需要提供旧密码）
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityAuthPassword(
        req: ChangePasswordRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/auth/password',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 强制下次登录修改密码
     * 管理员标记用户需要在下次登录时强制修改密码
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityAuthPasswordForceChange(
        req: ForcePasswordChangeRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/auth/password/force-change',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 重置密码
     * 管理员重置用户密码（管理员权限）
     * @param userId 用户ID
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityAuthPasswordReset(
        userId: string,
        req: ResetPasswordRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/auth/password/reset',
            path: {
                'userID': userId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 刷新访问令牌
     * 使用刷新令牌获取新的访问令牌
     * @param req 请求体
     * @returns RefreshTokenResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityAuthRefresh(
        req: RefreshTokenRequestDTO,
    ): CancelablePromise<RefreshTokenResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/auth/refresh',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 创建部门
     * 在指定组织下创建新部门
     * @param req 请求体
     * @returns DepartmentResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityDepartments(
        req: CreateDepartmentRequestDTO,
    ): CancelablePromise<DepartmentResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/departments',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 删除部门
     * 软删除指定部门
     * @param departmentId 部门ID
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static deleteApiV1IdentityDepartments(
        departmentId: string,
        req: DeleteDepartmentRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/identity/departments/{departmentID}',
            path: {
                'departmentID': departmentId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `部门未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取部门信息
     * 根据部门ID获取部门详细信息
     * @param departmentId 部门ID
     * @returns DepartmentResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityDepartments(
        departmentId: string,
    ): CancelablePromise<DepartmentResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/departments/{departmentID}',
            path: {
                'departmentID': departmentId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `部门未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 更新部门信息
     * 更新指定部门的信息
     * @param departmentId 部门ID
     * @param req 请求体
     * @returns DepartmentResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityDepartments(
        departmentId: string,
        req: UpdateDepartmentRequestDTO,
    ): CancelablePromise<DepartmentResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/departments/{departmentID}',
            path: {
                'departmentID': departmentId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `部门未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 删除Logo
     * 删除Logo文件和数据库记录（软删除），同时删除S3存储的文件
     * @param logoId Logo ID
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static deleteApiV1IdentityOrganizationLogos(
        logoId: string,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/identity/organization-logos/{logoID}',
            path: {
                'logoID': logoId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `Logo未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取Logo信息
     * 根据Logo ID获取Logo元数据和预签名下载URL（有效期7天）
     * @param logoId Logo ID
     * @returns OrganizationLogoResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityOrganizationLogos(
        logoId: string,
    ): CancelablePromise<OrganizationLogoResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/organization-logos/{logoID}',
            path: {
                'logoID': logoId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `Logo未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 上传临时Logo
     * 上传组织Logo文件到临时存储（7天后自动过期），返回Logo元数据和预签名下载URL
     * @param fileName 文件名
     * @param fileContent 文件内容
     * @param mimeType MIME类型
     * @returns OrganizationLogoResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityOrganizationLogosTemporary(
        fileName: string,
        fileContent: Blob,
        mimeType?: string,
    ): CancelablePromise<OrganizationLogoResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/organization-logos/temporary',
            formData: {
                'file_name': fileName,
                'file_content': fileContent,
                'mime_type': mimeType,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                413: `文件过大`,
                415: `不支持的文件类型`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取组织列表
     * 分页查询组织列表，支持按父组织筛选
     * @param page 页码
     * @param limit 每页数量
     * @param search 搜索关键词
     * @param filter 字段级过滤
     * @param sort 排序规则
     * @param fields 指定返回字段
     * @param includeTotal 是否返回总数
     * @param parentId 按父组织ID筛选
     * @returns ListOrganizationsResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityOrganizations(
        page: number = 1,
        limit: number = 20,
        search?: string,
        filter?: string,
        sort?: string,
        fields?: string,
        includeTotal: boolean = false,
        parentId?: string,
    ): CancelablePromise<ListOrganizationsResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/organizations',
            query: {
                'page': page,
                'limit': limit,
                'search': search,
                'filter': filter,
                'sort': sort,
                'fields': fields,
                'include_total': includeTotal,
                'parent_id': parentId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 创建组织
     * 创建新的组织机构
     * @param req 请求体
     * @returns OrganizationResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityOrganizations(
        req: CreateOrganizationRequestDTO,
    ): CancelablePromise<OrganizationResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/organizations',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 删除组织
     * 软删除指定组织
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static deleteApiV1IdentityOrganizations(
        req: DeleteOrganizationRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/identity/organizations/{organizationID}',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `组织未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取组织信息
     * 根据组织ID获取组织详细信息
     * @param organizationId 组织ID
     * @returns OrganizationResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityOrganizations1(
        organizationId: string,
    ): CancelablePromise<OrganizationResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/organizations/{organizationID}',
            path: {
                'organizationID': organizationId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `组织未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 更新组织信息
     * 更新指定组织的信息
     * @param organizationId 组织ID
     * @param req 请求体
     * @returns OrganizationResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityOrganizations(
        organizationId: string,
        req: UpdateOrganizationRequestDTO,
    ): CancelablePromise<OrganizationResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/organizations/{organizationID}',
            path: {
                'organizationID': organizationId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `组织未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取组织部门列表
     * 获取指定组织下的所有部门
     * @param organizationId 组织ID
     * @param page 页码
     * @param limit 每页数量
     * @param search 搜索关键词
     * @param filter 字段级过滤
     * @param sort 排序规则
     * @param fields 指定返回字段
     * @param includeTotal 是否返回总数
     * @returns GetOrganizationDepartmentsResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityOrganizationsDepartments(
        organizationId: string,
        page: number = 1,
        limit: number = 20,
        search?: string,
        filter?: string,
        sort?: string,
        fields?: string,
        includeTotal: boolean = false,
    ): CancelablePromise<GetOrganizationDepartmentsResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/organizations/{organizationID}/departments',
            path: {
                'organizationID': organizationId,
            },
            query: {
                'page': page,
                'limit': limit,
                'search': search,
                'filter': filter,
                'sort': sort,
                'fields': fields,
                'include_total': includeTotal,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `组织未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 绑定Logo到组织
     * 将临时Logo绑定到组织，Logo状态从临时变更为永久保存，并更新组织的Logo信息
     * @param organizationId 组织ID
     * @param req 请求体
     * @returns OrganizationResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityOrganizationsLogo(
        organizationId: string,
        req: BindLogoToOrganizationRequestDTO,
    ): CancelablePromise<OrganizationResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/organizations/{organizationID}/logo',
            path: {
                'organizationID': organizationId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `Logo或组织未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取用户列表
     * 分页查询用户列表，支持按组织、状态等条件筛选
     * @param page 页码
     * @param limit 每页数量
     * @param search 搜索关键词
     * @param filter 字段级过滤
     * @param sort 排序规则
     * @param fields 指定返回字段
     * @param includeTotal 是否返回总数
     * @param organizationId 按组织ID筛选
     * @param status 按用户状态筛选
     * @param fetchAll 是否获取所有数据（不分页）
     * @returns ListUsersResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityUsers(
        page: number = 1,
        limit: number = 20,
        search?: string,
        filter?: string,
        sort?: string,
        fields?: string,
        includeTotal: boolean = false,
        organizationId?: string,
        status?: number,
        fetchAll: boolean = false,
    ): CancelablePromise<ListUsersResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/users',
            query: {
                'page': page,
                'limit': limit,
                'search': search,
                'filter': filter,
                'sort': sort,
                'fields': fields,
                'include_total': includeTotal,
                'organization_id': organizationId,
                'status': status,
                'fetch_all': fetchAll,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 创建用户
     * 管理员创建新用户账户
     * @param req 请求体
     * @returns UserProfileResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1IdentityUsers(
        req: CreateUserRequestDTO,
    ): CancelablePromise<UserProfileResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/identity/users',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 删除用户
     * 软删除指定用户（管理员权限）
     * @param userId 用户ID
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static deleteApiV1IdentityUsers(
        userId: string,
        req?: DeleteUserRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/identity/users/{userID}',
            path: {
                'userID': userId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `用户未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取用户信息
     * 根据用户ID获取用户详细信息
     * @param userId 用户ID
     * @returns UserProfileResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityUsers1(
        userId: string,
    ): CancelablePromise<UserProfileResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/users/{userID}',
            path: {
                'userID': userId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `用户未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 更新用户信息
     * 更新指定用户的基本信息
     * @param userId 用户ID
     * @param req 请求体
     * @returns UserProfileResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityUsers(
        userId: string,
        req: UpdateUserRequestDTO,
    ): CancelablePromise<UserProfileResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/users/{userID}',
            path: {
                'userID': userId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `用户未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取用户成员关系
     * 获取指定用户的所有组织成员关系
     * @param userId 用户ID
     * @param page 页码
     * @param limit 每页数量
     * @param search 搜索关键词
     * @param filter 字段级过滤
     * @param sort 排序规则
     * @param fields 指定返回字段
     * @param includeTotal 是否返回总数
     * @returns GetUserMembershipsResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityUsersMemberships(
        userId: string,
        page: number = 1,
        limit: number = 20,
        search?: string,
        filter?: string,
        sort?: string,
        fields?: string,
        includeTotal: boolean = false,
    ): CancelablePromise<GetUserMembershipsResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/users/{userID}/memberships',
            path: {
                'userID': userId,
            },
            query: {
                'page': page,
                'limit': limit,
                'search': search,
                'filter': filter,
                'sort': sort,
                'fields': fields,
                'include_total': includeTotal,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `用户未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取用户主成员关系
     * 获取指定用户的主成员关系
     * @param req 请求体
     * @returns UserMembershipResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityUsersPrimaryMembership(
        req: GetPrimaryMembershipRequestDTO,
    ): CancelablePromise<UserMembershipResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/users/{userID}/primary-membership',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `主成员关系未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 变更用户状态
     * 管理员变更用户状态（激活、停用、锁定等）
     * @param userId 用户ID
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityUsersStatus(
        userId: string,
        req: ChangeUserStatusRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/users/{userID}/status',
            path: {
                'userID': userId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `用户未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 解锁用户
     * 管理员解锁被锁定的用户
     * @param userId 用户ID
     * @param req 请求体
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityUsersUnlock(
        userId: string,
        req: UnlockUserRequestDTO,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/users/{userID}/unlock',
            path: {
                'userID': userId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                403: `权限不足`,
                404: `用户未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取当前用户信息
     * 获取当前登录用户的详细信息
     * @returns UserProfileResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityUsersMe(): CancelablePromise<UserProfileResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/users/me',
            errors: {
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 更新当前用户信息
     * 用户更新自己的基本信息（从认证上下文获取用户ID）
     * @param req 请求体
     * @returns UserProfileResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1IdentityUsersMe(
        req: UpdateMeRequestDTO,
    ): CancelablePromise<UserProfileResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/identity/users/me',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 搜索用户
     * 按关键词搜索用户
     * @param page 页码
     * @param limit 每页数量
     * @param search 搜索关键词
     * @param filter 字段级过滤
     * @param sort 排序规则
     * @param fields 指定返回字段
     * @param includeTotal 是否返回总数
     * @param organizationId 按组织ID筛选
     * @returns SearchUsersResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1IdentityUsersSearch(
        page: number = 1,
        limit: number = 20,
        search?: string,
        filter?: string,
        sort?: string,
        fields?: string,
        includeTotal: boolean = false,
        organizationId?: string,
    ): CancelablePromise<SearchUsersResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/identity/users/search',
            query: {
                'page': page,
                'limit': limit,
                'search': search,
                'filter': filter,
                'sort': sort,
                'fields': fields,
                'include_total': includeTotal,
                'organization_id': organizationId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取菜单树结构
     * 获取完整的菜单树结构，用于前端展示导航菜单预览
     * @returns GetMenuTreeResponseDTO 成功返回菜单树结构
     * @throws ApiError
     */
    public static getApiV1PermissionMenuTree(): CancelablePromise<GetMenuTreeResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/menu/tree',
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 上传菜单配置文件
     * 上传YAML格式的菜单配置文件，用于更新系统菜单
     * @param menuFile YAML格式的菜单配置文件
     * @returns GetMenuTreeResponseDTO 菜单上传成功
     * @throws ApiError
     */
    public static postApiV1PermissionMenuUpload(
        menuFile: Blob,
    ): CancelablePromise<GetMenuTreeResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/permission/menu/upload',
            formData: {
                'menu_file': menuFile,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 列出角色定义
     * 分页列出角色定义
     * @param page 页码
     * @param limit 每页数量
     * @param search 搜索关键词
     * @param filter 字段级过滤
     * @param sort 排序规则
     * @param fields 指定返回字段
     * @param includeTotal 是否返回总数
     * @param name 角色名称
     * @param status 角色状态
     * @param isSystemRole 是否为系统角色
     * @param fetchAll 是否获取所有数据（不分页）
     * @returns RoleDefinitionListResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionRoles(
        page: number = 1,
        limit: number = 20,
        search?: string,
        filter?: string,
        sort?: string,
        fields?: string,
        includeTotal: boolean = false,
        name?: string,
        status?: string,
        isSystemRole?: boolean,
        fetchAll: boolean = false,
    ): CancelablePromise<RoleDefinitionListResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/roles',
            query: {
                'page': page,
                'limit': limit,
                'search': search,
                'filter': filter,
                'sort': sort,
                'fields': fields,
                'include_total': includeTotal,
                'name': name,
                'status': status,
                'is_system_role': isSystemRole,
                'fetch_all': fetchAll,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 创建角色定义
     * 创建一个新的角色定义
     * @param req 请求体
     * @returns RoleDefinitionCreateResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1PermissionRoles(
        req: RoleDefinitionCreateRequestDTO,
    ): CancelablePromise<RoleDefinitionCreateResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/permission/roles',
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 更新角色定义
     * 更新指定角色定义的信息
     * @param roleDefinitionId 角色定义ID
     * @param req 请求体
     * @returns RoleDefinitionUpdateResponseDTO 成功
     * @throws ApiError
     */
    public static putApiV1PermissionRoles(
        roleDefinitionId: string,
        req: RoleDefinitionUpdateRequestDTO,
    ): CancelablePromise<RoleDefinitionUpdateResponseDTO> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/api/v1/permission/roles/{roleDefinitionID}',
            path: {
                'roleDefinitionID': roleDefinitionId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 删除角色定义
     * 删除一个角色定义
     * @param roleId 角色ID
     * @returns OperationStatusResponseDTO 成功
     * @throws ApiError
     */
    public static deleteApiV1PermissionRoles(
        roleId: string,
    ): CancelablePromise<OperationStatusResponseDTO> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/api/v1/permission/roles/{roleID}',
            path: {
                'roleID': roleId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `角色未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取角色定义
     * 根据ID获取角色定义
     * @param roleId 角色ID
     * @returns RoleDefinitionGetResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionRoles1(
        roleId: string,
    ): CancelablePromise<RoleDefinitionGetResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/roles/{roleID}',
            path: {
                'roleID': roleId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `角色未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 检查角色是否具有指定菜单权限
     * 检查指定角色是否具有对指定菜单的特定权限
     * @param roleId 角色ID
     * @param req 请求体
     * @returns HasMenuPermissionResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1PermissionRolesCheckMenuPermission(
        roleId: string,
        req: HasMenuPermissionRequestDTO,
    ): CancelablePromise<HasMenuPermissionResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/permission/roles/{roleID}/check-menu-permission',
            path: {
                'roleID': roleId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取角色的菜单权限列表
     * 获取指定角色的所有菜单权限配置列表
     * @param roleId 角色ID
     * @returns GetRoleMenuPermissionsResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionRolesMenuPermissions(
        roleId: string,
    ): CancelablePromise<GetRoleMenuPermissionsResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/roles/{roleID}/menu-permissions',
            path: {
                'roleID': roleId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取角色的菜单树
     * 根据角色的权限配置，返回该角色可访问的菜单树结构
     * @param roleId 角色ID
     * @returns GetRoleMenuTreeResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionRolesMenuTree(
        roleId: string,
    ): CancelablePromise<GetRoleMenuTreeResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/roles/{roleID}/menu-tree',
            path: {
                'roleID': roleId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 配置角色的菜单权限
     * 为指定角色配置菜单权限，会清除角色的旧菜单映射，然后添加新的映射
     * @param roleId 角色ID
     * @param req 请求体
     * @returns ConfigureRoleMenusResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1PermissionRolesMenus(
        roleId: string,
        req: ConfigureRoleMenusRequestDTO,
    ): CancelablePromise<ConfigureRoleMenusResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/permission/roles/{roleID}/menus',
            path: {
                'roleID': roleId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 根据角色ID获取所有用户
     * 获取指定角色下所有用户的ID列表（不分页）
     * @param roleId 角色ID
     * @returns GetUsersByRoleResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionRolesUsers(
        roleId: string,
    ): CancelablePromise<GetUsersByRoleResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/roles/{roleID}/users',
            path: {
                'roleID': roleId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `角色未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 批量绑定用户到角色
     * 批量为角色绑定用户，覆盖旧的绑定关系
     * @param roleId 角色ID
     * @param req 请求体
     * @returns BatchBindUsersToRoleResponseDTO 成功
     * @throws ApiError
     */
    public static postApiV1PermissionRolesUsersBatchBind(
        roleId: string,
        req: BatchBindUsersToRoleRequestDTO,
    ): CancelablePromise<BatchBindUsersToRoleResponseDTO> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/api/v1/permission/roles/{roleID}/users/batch-bind',
            path: {
                'roleID': roleId,
            },
            body: req,
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `角色未找到`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 列出用户的角色分配记录
     * 分页列出用户的角色分配记录
     * @param page 页码
     * @param limit 每页数量
     * @param search 搜索关键词
     * @param filter 字段级过滤
     * @param sort 排序规则
     * @param fields 指定返回字段
     * @param includeTotal 是否返回总数
     * @param userId 用户ID
     * @param roleId 角色ID
     * @returns UserRoleListResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionUserRoles(
        page: number = 1,
        limit: number = 20,
        search?: string,
        filter?: string,
        sort?: string,
        fields?: string,
        includeTotal: boolean = false,
        userId?: string,
        roleId?: string,
    ): CancelablePromise<UserRoleListResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/user-roles',
            query: {
                'page': page,
                'limit': limit,
                'search': search,
                'filter': filter,
                'sort': sort,
                'fields': fields,
                'include_total': includeTotal,
                'user_id': userId,
                'role_id': roleId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取用户的菜单树
     * 根据用户的最新角色绑定，返回该用户可访问的菜单树结构
     * @param userId 用户ID
     * @returns GetUserMenuTreeResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionUsersMenuTree(
        userId: string,
    ): CancelablePromise<GetUserMenuTreeResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/users/{userID}/menu-tree',
            path: {
                'userID': userId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                500: `内部错误`,
            },
        });
    }
    /**
     * 获取用户最后一次角色分配
     * 获取用户最后一次角色分配
     * @param userId 用户ID
     * @returns AssignRoleToUserResponseDTO 成功
     * @throws ApiError
     */
    public static getApiV1PermissionUsersRolesLatest(
        userId: string,
    ): CancelablePromise<AssignRoleToUserResponseDTO> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/api/v1/permission/users/{userID}/roles/latest',
            path: {
                'userID': userId,
            },
            errors: {
                400: `请求参数错误`,
                401: `认证失败`,
                404: `分配未找到`,
                500: `内部错误`,
            },
        });
    }
}
