// API 响应基础结构
export interface BaseResponse {
  base_resp: {
    code: number
    message: string
  }
}

// 分页参数
export interface PageParams {
  page?: number
  limit?: number
  search?: string
  sort?: string
}

// 分页响应
export interface PageResponse<T> {
  items?: T[]
  page?: {
    total: number
    page: number
    limit: number
    total_pages: number
    has_next: boolean
    has_prev: boolean
  }
}

// 用户状态枚举
export enum UserStatus {
  ACTIVE = 1,
  INACTIVE = 2,
  SUSPENDED = 3,
  LOCKED = 4
}

// 性别枚举
export enum Gender {
  UNKNOWN = 0,
  MALE = 1,
  FEMALE = 2
}

// 权限级别枚举
export enum PermissionLevel {
  NONE = 0,
  READ = 1,
  WRITE = 2,
  FULL = 3
}
