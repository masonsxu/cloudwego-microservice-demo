/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type PageResponseDTO = {
    /**
     * 是否有下一页
     */
    has_next?: boolean;
    /**
     * 是否有上一页
     */
    has_prev?: boolean;
    /**
     * 每页数量
     */
    limit?: number;
    /**
     * 当前页码
     */
    page?: number;
    /**
     * 总记录数
     */
    total?: number;
    /**
     * 总页数
     */
    total_pages?: number;
};

