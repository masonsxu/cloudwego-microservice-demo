import request from './request'
import type {
  Organization,
  Department,
  CreateOrganizationRequest,
  UpdateOrganizationRequest,
  CreateDepartmentRequest,
  UpdateDepartmentRequest
} from '@/types/organization'
import type { PageParams, PageResponse } from '@/types/api'

// 获取组织列表
export function getOrganizationList(params?: PageParams & { parent_id?: string }) {
  return request<{ organizations: Organization[]; page: PageResponse<Organization> }>({
    url: '/api/v1/identity/organizations',
    method: 'GET',
    params
  })
}

// 获取组织详情
export function getOrganizationDetail(orgId: string) {
  return request<{ organization: Organization }>({
    url: `/api/v1/identity/organizations/${orgId}`,
    method: 'GET'
  })
}

// 创建组织
export function createOrganization(data: CreateOrganizationRequest) {
  return request<{ organization: Organization }>({
    url: '/api/v1/identity/organizations',
    method: 'POST',
    data
  })
}

// 更新组织
export function updateOrganization(orgId: string, data: UpdateOrganizationRequest) {
  return request<{ organization: Organization }>({
    url: `/api/v1/identity/organizations/${orgId}`,
    method: 'PUT',
    data
  })
}

// 删除组织
export function deleteOrganization(orgId: string) {
  return request({
    url: `/api/v1/identity/organizations/${orgId}`,
    method: 'DELETE'
  })
}

// 获取组织的部门列表
export function getDepartmentList(orgId: string, params?: PageParams) {
  return request<{ departments: Department[]; page: PageResponse<Department> }>({
    url: `/api/v1/identity/organizations/${orgId}/departments`,
    method: 'GET',
    params
  })
}

// 获取部门详情
export function getDepartmentDetail(deptId: string) {
  return request<{ department: Department }>({
    url: `/api/v1/identity/departments/${deptId}`,
    method: 'GET'
  })
}

// 创建部门
export function createDepartment(data: CreateDepartmentRequest) {
  return request<{ department: Department }>({
    url: '/api/v1/identity/departments',
    method: 'POST',
    data
  })
}

// 更新部门
export function updateDepartment(deptId: string, data: UpdateDepartmentRequest) {
  return request<{ department: Department }>({
    url: `/api/v1/identity/departments/${deptId}`,
    method: 'PUT',
    data
  })
}

// 删除部门
export function deleteDepartment(deptId: string) {
  return request({
    url: `/api/v1/identity/departments/${deptId}`,
    method: 'DELETE'
  })
}
