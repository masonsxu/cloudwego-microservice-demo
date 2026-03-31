// OAuth2 客户端（MVP：authorization_code + refresh_token）
export type OAuth2ClientType = 'confidential' | 'public'
export type OAuth2GrantType = 'authorization_code' | 'refresh_token'

export interface OAuth2Client {
  id: string
  client_id: string
  client_name: string
  description?: string
  client_type: OAuth2ClientType
  grant_types: OAuth2GrantType[]
  redirect_uris?: string[]
  scopes?: string[]
  logo_uri?: string
  client_uri?: string
  access_token_lifespan: number
  refresh_token_lifespan: number
  is_active: boolean
  owner_id?: string
  created_at: number
  updated_at: number
}

export interface CreateOAuth2ClientRequest {
  client_name: string
  description?: string
  client_type: OAuth2ClientType
  grant_types: OAuth2GrantType[]
  redirect_uris?: string[]
  scopes?: string[]
  logo_uri?: string
  client_uri?: string
  access_token_lifespan?: number
  refresh_token_lifespan?: number
}

export interface UpdateOAuth2ClientRequest {
  client_name?: string
  description?: string
  grant_types?: OAuth2GrantType[]
  redirect_uris?: string[]
  scopes?: string[]
  logo_uri?: string
  client_uri?: string
  access_token_lifespan?: number
  refresh_token_lifespan?: number
  is_active?: boolean
}

export interface CreateOAuth2ClientResponse {
  client: OAuth2Client
  client_secret: string
}

// OAuth2 作用域
export interface OAuth2Scope {
  id: string
  name: string
  display_name: string
  description?: string
  is_default: boolean
  is_system: boolean
}

export interface OAuth2Config {
  enabled: boolean
  issuer: string
  access_token_lifespan: number
  refresh_token_lifespan: number
  auth_code_lifespan: number
  enforce_pkce: boolean
  consent_page_url: string
}

// OAuth2 授权同意
export interface OAuth2Consent {
  id: string
  client_id: string
  client_name: string
  scopes: string[]
  granted_at: number
  is_revoked: boolean
}

// 列表参数
export interface ListOAuth2ClientsParams {
  page?: number
  limit?: number
  is_active?: boolean
}
