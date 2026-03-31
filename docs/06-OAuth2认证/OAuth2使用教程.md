# OAuth2 配置与使用教程（MVP）

本教程面向本项目当前 OAuth2 MVP 能力，覆盖两部分：

1. **配置**：如何开启 OAuth2、关键环境变量如何设置。
2. **使用**：如何创建客户端，并走通 `authorization_code + PKCE + refresh_token`。

> 当前 MVP 仅支持：`authorization_code`、`refresh_token`。
>
> 当前 MVP 不支持：`client_credentials`、`/oauth2/revoke`、`/oauth2/introspect`。

## 1. 功能边界（先看这个）

- 协议端点：
  - `GET /oauth2/authorize`
  - `POST /oauth2/token`
- 管理端点：
  - `POST /api/v1/oauth2/clients`
  - `GET /api/v1/oauth2/clients`
  - `GET /api/v1/oauth2/clients/:id`
  - `PUT /api/v1/oauth2/clients/:id`
  - `DELETE /api/v1/oauth2/clients/:id`
  - `POST /api/v1/oauth2/clients/:id/rotate-secret`
  - `GET /api/v1/oauth2/scopes`
  - `GET /api/v1/oauth2/consents`
  - `DELETE /api/v1/oauth2/consents/:clientID`

## 2. OAuth2 配置

OAuth2 配置位于网关配置（`gateway`）中，默认值可参考：
- `gateway/internal/infrastructure/config/defaults.go`
- `gateway/internal/infrastructure/config/types.go`

### 2.1 关键环境变量

| 环境变量 | 说明 | 默认值 |
|---|---|---|
| `OAUTH2_ENABLED` | 是否启用 OAuth2 | `false` |
| `OAUTH2_ISSUER` | 签发者（issuer），建议与网关外部访问地址一致 | `http://localhost:8080` |
| `OAUTH2_ACCESS_TOKEN_LIFESPAN` | Access Token 有效期 | `1h` |
| `OAUTH2_REFRESH_TOKEN_LIFESPAN` | Refresh Token 有效期 | `720h`（30天） |
| `OAUTH2_AUTH_CODE_LIFESPAN` | 授权码有效期 | `10m` |
| `OAUTH2_ENFORCE_PKCE` | 是否强制 PKCE | `true` |
| `OAUTH2_CONSENT_PAGE_URL` | 同意页 URL（当前默认） | `/oauth2/consent` |

### 2.2 本地开发建议配置

```bash
# gateway/.env（示例）
OAUTH2_ENABLED=true
OAUTH2_ISSUER=http://localhost:8080
OAUTH2_ACCESS_TOKEN_LIFESPAN=1h
OAUTH2_REFRESH_TOKEN_LIFESPAN=720h
OAUTH2_AUTH_CODE_LIFESPAN=10m
OAUTH2_ENFORCE_PKCE=true
OAUTH2_CONSENT_PAGE_URL=/oauth2/consent
```

> 注意：`/oauth2/token` 在默认 JWT 跳过路径中，属于 OAuth2 协议端点行为，非后台管理接口。

## 3. 启动服务

```bash
# 启动基础设施
cd docker && podman compose up -d

# 启动 RPC 服务
cd rpc/identity_srv && sh build.sh && sh output/bootstrap.sh

# 启动网关服务
cd gateway && sh build.sh && sh output/bootstrap.sh
```

## 4. 创建 OAuth2 客户端

可通过管理后台（系统设置 → OAuth2）或 API 创建。

### 4.1 创建 confidential 客户端（示例）

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "client_name": "Demo Web Backend",
    "description": "示例后端应用",
    "client_type": "confidential",
    "grant_types": ["authorization_code", "refresh_token"],
    "redirect_uris": ["http://localhost:3000/callback"],
    "scopes": ["openid", "profile", "email"]
  }'
```

响应中的 `client_secret` **只展示一次**，请立即保存。

### 4.2 创建 public 客户端（示例）

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "client_name": "Demo SPA",
    "description": "示例前端应用",
    "client_type": "public",
    "grant_types": ["authorization_code", "refresh_token"],
    "redirect_uris": ["http://localhost:5173/callback"],
    "scopes": ["openid", "profile"]
  }'
```

> `public` 客户端必须使用 PKCE（建议 `S256`）。

## 5. 使用流程（Authorization Code + PKCE）

下面是完整三步流程。

### 步骤 1：请求授权码

浏览器重定向到：

```text
GET /oauth2/authorize?
  response_type=code&
  client_id=CLIENT_ID&
  redirect_uri=REDIRECT_URI&
  scope=openid%20profile&
  state=RANDOM_STATE&
  code_challenge=CODE_CHALLENGE&
  code_challenge_method=S256
```

必要参数：
- `response_type=code`
- `client_id`
- `redirect_uri`
- `state`
- `code_challenge` + `code_challenge_method`

授权成功后会回跳：

```text
REDIRECT_URI?code=AUTH_CODE&state=RANDOM_STATE
```

### 步骤 2：用授权码换取 Token

```bash
curl -X POST http://localhost:8080/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -u "CLIENT_ID:CLIENT_SECRET" \
  -d "grant_type=authorization_code" \
  -d "code=AUTH_CODE" \
  -d "redirect_uri=REDIRECT_URI" \
  -d "code_verifier=CODE_VERIFIER"
```

返回示例：

```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "...",
  "scope": "openid profile"
}
```

> 对 `public` 客户端，不应在前端暴露 `client_secret`。认证方式按你客户端模型与服务配置执行，核心是必须提供 `code_verifier`。

### 步骤 3：用 refresh_token 刷新 access_token

```bash
curl -X POST http://localhost:8080/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -u "CLIENT_ID:CLIENT_SECRET" \
  -d "grant_type=refresh_token" \
  -d "refresh_token=REFRESH_TOKEN"
```

## 6. 客户端管理常用操作

### 6.1 查询客户端列表

```bash
curl -X GET "http://localhost:8080/api/v1/oauth2/clients?page=1&limit=20" \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 6.2 更新客户端

```bash
curl -X PUT http://localhost:8080/api/v1/oauth2/clients/CLIENT_UUID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "client_name": "Updated App Name",
    "redirect_uris": ["https://example.com/callback"],
    "is_active": true
  }'
```

### 6.3 轮换密钥

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/clients/CLIENT_UUID/rotate-secret \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 6.4 删除客户端

```bash
curl -X DELETE http://localhost:8080/api/v1/oauth2/clients/CLIENT_UUID \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## 7. 用户授权管理（consents）

### 7.1 查询我的授权记录

```bash
curl -X GET http://localhost:8080/api/v1/oauth2/consents \
  -H "Authorization: Bearer USER_TOKEN"
```

### 7.2 撤销我的授权

```bash
curl -X DELETE http://localhost:8080/api/v1/oauth2/consents/CLIENT_ID \
  -H "Authorization: Bearer USER_TOKEN"
```

## 8. 常见问题

### 8.1 `invalid_grant`

常见原因：
- 授权码过期（默认 10 分钟）
- `code_verifier` 与 `code_challenge` 不匹配
- `redirect_uri` 与创建客户端时不一致

### 8.2 `invalid_client`

常见原因：
- `client_id` / `client_secret` 错误
- 客户端被禁用或已删除
- 轮换密钥后仍使用旧密钥

### 8.3 为什么没有 revoke/introspect 端点？

因为当前是 MVP 收敛版本，显式不支持 `revoke` / `introspect`，后续如扩展会单独发布变更说明。

## 9. 参考

- [快速开始](../01-快速入门/快速开始.md)
- [配置参考](../01-快速入门/配置参考.md)
- [架构设计](../00-项目概览/架构设计.md)
- [RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749)
- [RFC 7636 (PKCE)](https://datatracker.ietf.org/doc/html/rfc7636)
