/**
 * OIDC (OpenID Connect) 类型定义
 */

// ==================== OIDC Discovery ====================

export interface OIDCDiscoveryConfig {
  issuer: string
  authorization_endpoint: string
  token_endpoint: string
  userinfo_endpoint: string
  revocation_endpoint: string
  introspection_endpoint: string
  jwks_uri: string
  response_types_supported: string[]
  subject_types_supported: string[]
  id_token_signing_alg_values_supported: string[]
  scopes_supported: string[]
  token_endpoint_auth_methods_supported: string[]
}

// ==================== OIDC JWKS ====================

export interface OIDCJWKSKey {
  kty: string
  use: string
  kid: string
  alg: string
  n: string
  e: string
}

export interface OIDCJWKSResponse {
  keys: OIDCJWKSKey[]
}

// ==================== OIDC Client (OIDC 应用) ====================

export type OIDCClientType = 'confidential' | 'public'
export type OIDCGrantType = 'authorization_code' | 'refresh_token'

export interface OIDCClient {
  id: string
  client_id: string
  client_name: string
  client_type: OIDCClientType
  grant_types: OIDCGrantType[]
  redirect_uris: string[]
  scopes: string[]
  is_active: boolean
  access_token_lifespan: number
  refresh_token_lifespan: number
  description?: string
  created_at: number
  updated_at: number
}

export interface CreateOIDCClientRequest {
  client_name: string
  description?: string
  client_type: OIDCClientType
  grant_types: OIDCGrantType[]
  redirect_uris?: string[]
}

export interface CreateOIDCClientResponse {
  client: OIDCClient
  client_secret: string
}

export interface UpdateOIDCClientRequest {
  client_name?: string
  description?: string
  is_active?: boolean
  redirect_uris?: string[]
}

export interface ListOIDCClientsParams {
  page?: number
  limit?: number
}

export interface ListOIDCClientsResponse {
  clients: OIDCClient[]
  page: {
    total: number
    page: number
    size: number
  }
}

export interface GetOIDCClientResponse {
  client: OIDCClient
}

export interface RotateOIDCClientSecretResponse {
  client_secret: string
}

// ==================== OIDC Consent (授权同意) ====================

export interface OIDCConsent {
  id: string
  client_id: string
  client_name: string
  granted_at: number
  scopes: string[]
}

export interface ListOIDCConsentsResponse {
  consents: OIDCConsent[]
  page: {
    total: number
    page: number
    size: number
  }
}

// ==================== OIDC Config (运行时配置) ====================

export interface OIDCConfig {
  enabled: boolean
  enforce_pkce: boolean
  issuer: string
  access_token_lifespan: number
  refresh_token_lifespan: number
  auth_code_lifespan: number
  consent_page_url: string
}

export interface GetOIDCConfigResponse {
  config: OIDCConfig
}

// ==================== OIDC Scopes ====================

export interface OIDCScope {
  name: string
  description?: string
}

export interface ListOIDCScopesResponse {
  scopes: OIDCScope[]
}

// ==================== OIDC Token ====================

export interface OIDCTokenResponse {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token?: string
  id_token?: string
  scope?: string
}

// ==================== OIDC Userinfo ====================

export interface OIDCUserinfo {
  sub: string
  name?: string
  preferred_username?: string
  email?: string
  email_verified?: boolean
  picture?: string
}

// ==================== OIDC Introspect ====================

export interface OIDCIntrospectResponse {
  active: boolean
  client_id?: string
  username?: string
  scope?: string
  sub?: string
  exp?: number
  iat?: number
}
