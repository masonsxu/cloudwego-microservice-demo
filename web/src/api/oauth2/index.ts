import request from '../request'
import type {
  OAuth2Client,
  OAuth2Scope,
  OAuth2Consent,
  CreateOAuth2ClientRequest,
  CreateOAuth2ClientResponse,
  UpdateOAuth2ClientRequest,
  ListOAuth2ClientsParams,
} from '@/types/oauth2'

export const oauth2Api = {
  // 客户端管理
  getClients: (params?: ListOAuth2ClientsParams) =>
    request<{ clients: OAuth2Client[]; page: any }>({
      method: 'GET',
      url: '/api/v1/oauth2/clients',
      params,
    }),

  getClient: (clientId: string) =>
    request<{ client: OAuth2Client }>({
      method: 'GET',
      url: `/api/v1/oauth2/clients/${clientId}`,
    }),

  createClient: (data: CreateOAuth2ClientRequest) =>
    request<CreateOAuth2ClientResponse>({
      method: 'POST',
      url: '/api/v1/oauth2/clients',
      data,
    }),

  updateClient: (clientId: string, data: UpdateOAuth2ClientRequest) =>
    request<{ client: OAuth2Client }>({
      method: 'PUT',
      url: `/api/v1/oauth2/clients/${clientId}`,
      data,
    }),

  deleteClient: (clientId: string) =>
    request({
      method: 'DELETE',
      url: `/api/v1/oauth2/clients/${clientId}`,
    }),

  rotateClientSecret: (clientId: string) =>
    request<{ client_secret: string }>({
      method: 'POST',
      url: `/api/v1/oauth2/clients/${clientId}/rotate-secret`,
    }),

  // 作用域
  getScopes: () =>
    request<{ scopes: OAuth2Scope[] }>({
      method: 'GET',
      url: '/api/v1/oauth2/scopes',
    }),

  // 授权同意
  getMyConsents: (params?: { page?: number; limit?: number }) =>
    request<{ consents: OAuth2Consent[]; page: any }>({
      method: 'GET',
      url: '/api/v1/oauth2/consents',
      params,
    }),

  revokeMyConsent: (clientId: string) =>
    request({
      method: 'DELETE',
      url: `/api/v1/oauth2/consents/${clientId}`,
    }),
}
