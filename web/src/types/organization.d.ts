import type { BaseResponse } from './api'

// 组织信息
export interface Organization {
  id: string
  code?: string
  name: string
  parent_id?: string
  logo?: string
  logo_id?: string
  province_city?: string[]
  facility_type?: string
  accreditation_status?: string
  member_count?: number
  department_count?: number
  created_at: number
  updated_at?: number
  parent?: Organization
  children?: Organization[]
}

// 部门信息
export interface Department {
  id: string
  code?: string
  name: string
  organization_id: string
  department_type?: string
  available_equipment?: string[]
  member_count?: number
  created_at: number
  updated_at?: number
  organization?: Organization
}

// 创建组织请求
export interface CreateOrganizationRequest {
  name: string
  parent_id?: string
  facility_type?: string
  accreditation_status?: string
  province_city?: string[]
}

// 更新组织请求
export interface UpdateOrganizationRequest {
  name?: string
  parent_id?: string
  facility_type?: string
  accreditation_status?: string
  province_city?: string[]
}

// 创建部门请求
export interface CreateDepartmentRequest {
  organization_id: string
  name: string
  department_type?: string
}

// 更新部门请求
export interface UpdateDepartmentRequest {
  name?: string
  department_type?: string
}
