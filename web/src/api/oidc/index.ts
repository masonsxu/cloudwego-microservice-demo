import request from '../request'
import type {
  OIDCDiscoveryConfig,
  OIDCJWKSResponse,
  OIDCUserinfo,
  OIDCIntrospectResponse,
} from '@/types/oidc'

export const oidcProviderApi = {
  getDiscovery: () =>
    request<OIDCDiscoveryConfig>({
      method: 'GET',
      url: '/.well-known/openid-configuration',
    }),

  getJWKS: () =>
    request<OIDCJWKSResponse>({
      method: 'GET',
      url: '/keys',
    }),

  getUserinfo: () =>
    request<OIDCUserinfo>({
      method: 'GET',
      url: '/userinfo',
    }),

  introspectToken: (data: { token: string; token_type_hint?: string }) =>
    request<OIDCIntrospectResponse>({
      method: 'POST',
      url: '/oauth/introspect',
      data,
    }),
}
