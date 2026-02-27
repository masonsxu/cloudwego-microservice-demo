import request from '../request'
import type { Organization, CreateOrganizationRequest, UpdateOrganizationRequest } from '@/types/organization'
import type { PageParams } from '@/types/api'

// 供 view 直接使用的类型别名（保持命名一致性）
export type OrganizationDTO = Organization
export type CreateOrganizationRequestDTO = CreateOrganizationRequest
export type UpdateOrganizationRequestDTO = UpdateOrganizationRequest

export interface OrganizationLogo {
  id: string
  file_name: string
  file_size: number
  mime_type: string
  download_url: string
  expires_at: number
  status: string
  bound_organization_id: string
  uploaded_by: string
  created_at: number
  updated_at: number
}

export interface ListOrganizationsParams extends PageParams {
  parent_id?: string
  include_total?: boolean
}

export const organizationApi = {
  // 获取组织列表
  listOrganizations: (params?: ListOrganizationsParams) =>
    request<{ organizations: Organization[]; page: any }>({
      method: 'GET',
      url: '/api/v1/identity/organizations',
      params,
    }),

  // 获取组织详情
  getOrganization: (organizationId: string) =>
    request<{ organization: Organization }>({
      method: 'GET',
      url: `/api/v1/identity/organizations/${organizationId}`,
    }),

  // 创建组织
  createOrganization: (data: CreateOrganizationRequest) =>
    request<{ organization: Organization }>({
      method: 'POST',
      url: '/api/v1/identity/organizations',
      data,
    }),

  // 更新组织
  updateOrganization: (organizationId: string, data: UpdateOrganizationRequest) =>
    request<{ organization: Organization }>({
      method: 'PUT',
      url: `/api/v1/identity/organizations/${organizationId}`,
      data,
    }),

  // 删除组织
  deleteOrganization: (organizationId: string) =>
    request({
      method: 'DELETE',
      url: `/api/v1/identity/organizations/${organizationId}`,
      data: {},
    }),

  // 绑定 Logo 到组织
  bindLogoToOrganization: (organizationId: string, logoId: string) =>
    request<{ organization: Organization }>({
      method: 'PUT',
      url: `/api/v1/identity/organizations/${organizationId}/logo`,
      data: { logo_id: logoId },
    }),

  // 上传临时 Logo
  uploadTemporaryLogo: (fileName: string, fileContent: Blob, mimeType?: string) => {
    const formData = new FormData()
    formData.append('file_name', fileName)
    formData.append('file_content', fileContent)
    if (mimeType) formData.append('mime_type', mimeType)
    return request<{ logo: OrganizationLogo }>({
      method: 'POST',
      url: '/api/v1/identity/organization-logos/temporary',
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  // 获取 Logo 详情
  getOrganizationLogo: (logoId: string) =>
    request<{ logo: OrganizationLogo }>({
      method: 'GET',
      url: `/api/v1/identity/organization-logos/${logoId}`,
    }),

  // 删除 Logo
  deleteOrganizationLogo: (logoId: string) =>
    request({
      method: 'DELETE',
      url: `/api/v1/identity/organization-logos/${logoId}`,
    }),
}

// 命名导出（供用户相关 view 使用）
export function getOrganizationList(params?: ListOrganizationsParams) {
  return organizationApi.listOrganizations(params)
}
