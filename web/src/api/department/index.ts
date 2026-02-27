import request from '../request'
import type { Department, CreateDepartmentRequest, UpdateDepartmentRequest } from '@/types/organization'
import type { PageParams } from '@/types/api'

export const departmentApi = {
  // 创建部门
  createDepartment: (data: CreateDepartmentRequest) =>
    request<{ department: Department }>({
      method: 'POST',
      url: '/api/v1/identity/departments',
      data,
    }),

  // 获取部门详情
  getDepartment: (departmentId: string) =>
    request<{ department: Department }>({
      method: 'GET',
      url: `/api/v1/identity/departments/${departmentId}`,
    }),

  // 更新部门
  updateDepartment: (departmentId: string, data: UpdateDepartmentRequest) =>
    request<{ department: Department }>({
      method: 'PUT',
      url: `/api/v1/identity/departments/${departmentId}`,
      data,
    }),

  // 删除部门
  deleteDepartment: (departmentId: string) =>
    request({
      method: 'DELETE',
      url: `/api/v1/identity/departments/${departmentId}`,
      data: {},
    }),

  // 获取组织下的部门列表
  getOrganizationDepartments: (organizationId: string, params?: PageParams) =>
    request<{ departments: Department[]; page: any }>({
      method: 'GET',
      url: `/api/v1/identity/organizations/${organizationId}/departments`,
      params,
    }),
}
