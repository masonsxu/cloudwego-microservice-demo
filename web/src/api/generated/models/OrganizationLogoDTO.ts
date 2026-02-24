/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type OrganizationLogoDTO = {
    /**
     * * 绑定的组织ID（临时状态时为空）
     */
    bound_organization_id?: string;
    /**
     * * 创建时间
     */
    created_at?: number;
    /**
     * * 下载URL（预签名URL）
     */
    download_url?: string;
    /**
     * * 过期时间（临时状态）
     */
    expires_at?: number;
    /**
     * * 文件存储ID（S3路径: bucket/key）
     */
    file_id?: string;
    /**
     * * 原始文件名
     */
    file_name?: string;
    /**
     * * 文件大小（字节）
     */
    file_size?: number;
    /**
     * * Logo唯一标识符
     */
    id?: string;
    /**
     * * MIME类型
     */
    mime_type?: string;
    /**
     * * Logo状态（TEMPORARY=临时, BOUND=已绑定, DELETED=已删除）
     */
    status?: string;
    /**
     * * 更新时间
     */
    updated_at?: number;
    /**
     * * 上传者ID
     */
    uploaded_by?: string;
};

