# OAuth2 使用教程

本指南介绍如何使用本项目的 OAuth2 功能，包括客户端管理、授权流程和用户授权管理。

## 目录

- [快速开始](#快速开始)
- [OAuth2 端点总览](#oauth2-端点总览)
- [创建 OAuth2 客户端](#创建-oauth2-客户端)
- [授权流程](#授权流程)
- [令牌管理](#令牌管理)
- [用户授权管理](#用户授权管理)
- [最佳实践](#最佳实践)

## 快速开始

### 1. 确保服务已启动

```bash
# 启动基础设施
cd docker && podman-compose up -d

# 启动 RPC 服务
cd rpc/identity_srv && sh build.sh && sh output/bootstrap.sh

# 启动网关服务
cd gateway && sh build.sh && sh output/bootstrap.sh
```

### 2. 创建第一个 OAuth2 客户端

使用以下命令创建一个测试客户端：

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "client_name": "My Test App",
    "description": "测试应用",
    "client_type": "confidential",
    "grant_types": ["authorization_code", "refresh_token"],
    "redirect_uris": ["http://localhost:3000/callback"],
    "scopes": ["openid", "profile", "email"]
  }'
```

**重要提示**：响应中的 `client_secret` 仅此一次可见，请妥善保存！

## OAuth2 端点总览

### 核心 OAuth2 协议端点

| 端点 | 方法 | 用途 | RFC 规范 |
|------|------|------|----------|
| `/oauth2/authorize` | GET | 授权端点（用户授权） | RFC 6749 |
| `/oauth2/token` | POST | 令牌端点（获取令牌） | RFC 6749 |
| `/oauth2/revoke` | POST | 令牌吊销 | RFC 7009 |
| `/oauth2/introspect` | POST | 令牌自省 | RFC 7662 |

### 客户端管理端点

| 端点 | 方法 | 用途 |
|------|------|------|
| `/api/v1/oauth2/clients` | POST | 创建客户端 |
| `/api/v1/oauth2/clients` | GET | 列出客户端 |
| `/api/v1/oauth2/clients/:id` | GET | 获取客户端详情 |
| `/api/v1/oauth2/clients/:id` | PUT | 更新客户端 |
| `/api/v1/oauth2/clients/:id` | DELETE | 删除客户端 |
| `/api/v1/oauth2/clients/:id/rotate-secret` | POST | 轮换客户端密钥 |

### 作用域和授权管理端点

| 端点 | 方法 | 用途 |
|------|------|------|
| `/api/v1/oauth2/scopes` | GET | 列出可用作用域 |
| `/api/v1/oauth2/consents` | GET | 查询我的授权记录 |
| `/api/v1/oauth2/consents/:clientID` | DELETE | 撤销我的授权 |

## 创建 OAuth2 客户端

### 客户端类型

本项目支持两种 OAuth2 客户端类型：

#### 1. Confidential Client（机密客户端）

适用于可以安全保管密钥的应用（如后端服务）：

- **特点**：拥有 `client_secret`，可以使用所有授权类型
- **适用场景**：Web 应用后端、移动应用后端、服务间调用
- **授权类型**：`authorization_code`、`client_credentials`、`refresh_token`

#### 2. Public Client（公共客户端）

适用于无法安全保管密钥的应用（如纯前端应用）：

- **特点**：无 `client_secret`，仅支持受限授权类型
- **适用场景**：单页应用（SPA）、移动应用（无后端）
- **授权类型**：通常使用 `authorization_code`（PKCE）

### 创建客户端示例

#### 创建机密客户端（Web 应用）

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "client_name": "E-commerce Backend",
    "description": "电商平台后端服务",
    "client_type": "confidential",
    "grant_types": ["authorization_code", "client_credentials", "refresh_token"],
    "redirect_uris": [
      "https://shop.example.com/auth/callback",
      "https://shop.example.com/auth/callback/native"
    ],
    "scopes": ["openid", "profile", "email", "orders.read", "orders.write"],
    "logo_uri": "https://shop.example.com/logo.png",
    "client_uri": "https://shop.example.com",
    "access_token_lifespan": 3600,
    "refresh_token_lifespan": 2592000
  }'
```

**响应示例**：

```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "client": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "client_id": "client_abc123xyz",
    "client_name": "E-commerce Backend",
    "description": "电商平台后端服务",
    "client_type": "confidential",
    "grant_types": ["authorization_code", "client_credentials", "refresh_token"],
    "redirect_uris": ["https://shop.example.com/auth/callback"],
    "scopes": ["openid", "profile", "email", "orders.read", "orders.write"],
    "is_active": true
  },
  "client_secret": "secret_xyz789abc"
}
```

#### 创建公共客户端（SPA）

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "client_name": "Task Manager SPA",
    "description": "任务管理单页应用",
    "client_type": "public",
    "grant_types": ["authorization_code"],
    "redirect_uris": [
      "http://localhost:5173/auth/callback",
      "https://tasks.example.com/auth/callback"
    ],
    "scopes": ["openid", "profile", "tasks.read", "tasks.write"]
  }'
```

### 客户端配置参数说明

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `client_name` | string | 是 | 客户端名称（1-128字符） |
| `description` | string | 否 | 客户端描述（最多512字符） |
| `client_type` | string | 是 | 客户端类型：`confidential` 或 `public` |
| `grant_types` | array | 是 | 允许的授权类型 |
| `redirect_uris` | array | 条件必填 | 授权码模式下必填 |
| `scopes` | array | 否 | 请求的权限范围 |
| `logo_uri` | string | 否 | 客户端 Logo URL |
| `client_uri` | string | 否 | 客户端主页 URL |
| `access_token_lifespan` | int32 | 否 | Access Token 有效期（秒），默认系统配置 |
| `refresh_token_lifespan` | int32 | 否 | Refresh Token 有效期（秒），默认系统配置 |

### 可用的授权类型

- `authorization_code`：授权码模式（推荐用于 Web 应用）
- `client_credentials`：客户端凭证模式（用于服务间调用）
- `refresh_token`：刷新令牌（获取新的访问令牌）
- `implicit`：隐式模式（不推荐，安全性较低）
- `password`：密码模式（不推荐，仅用于遗留系统）

### 可用的作用域（Scopes）

查询所有可用的作用域：

```bash
curl -X GET http://localhost:8080/api/v1/oauth2/scopes \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**示例响应**：

```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "scopes": [
    {
      "name": "openid",
      "description": "OpenID Connect 标识符"
    },
    {
      "name": "profile",
      "description": "访问用户基本信息"
    },
    {
      "name": "email",
      "description": "访问用户邮箱"
    }
  ]
}
```

## 授权流程

### 授权码模式（Authorization Code Flow）

适用于 Web 应用，是最安全的授权方式。

#### 步骤 1：重定向用户到授权端点

```
GET /oauth2/authorize?
    response_type=code&
    client_id=CLIENT_ID&
    redirect_uri=REDIRECT_URI&
    scope=openid+profile+email&
    state=RANDOM_STATE
```

**参数说明**：

- `response_type`：固定为 `code`
- `client_id`：客户端ID
- `redirect_uri`：回调地址（必须在客户端配置的 `redirect_uris` 中）
- `scope`：请求的权限范围（URL 编码）
- `state`：随机字符串，防止 CSRF 攻击

#### 步骤 2：用户授权

用户登录并确认授权后，会被重定向到：

```
REDIRECT_URI?
    code=AUTHORIZATION_CODE&
    state=RANDOM_STATE
```

#### 步骤 3：用授权码换取访问令牌

```bash
curl -X POST http://localhost:8080/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -u "CLIENT_ID:CLIENT_SECRET" \
  -d "grant_type=authorization_code" \
  -d "code=AUTHORIZATION_CODE" \
  -d "redirect_uri=REDIRECT_URI"
```

**响应示例**：

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
  "scope": "openid profile email"
}
```

### 客户端凭证模式（Client Credentials Flow）

适用于服务间调用（无用户上下文）。

```bash
curl -X POST http://localhost:8080/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -u "CLIENT_ID:CLIENT_SECRET" \
  -d "grant_type=client_credentials" \
  -d "scope=orders.read orders.write"
```

**响应示例**：

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "orders.read orders.write"
}
```

**注意**：此模式不返回 `refresh_token`，因为服务可以直接重新获取令牌。

### 刷新令牌（Refresh Token Flow）

使用 `refresh_token` 获取新的 `access_token`：

```bash
curl -X POST http://localhost:8080/oauth2/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -u "CLIENT_ID:CLIENT_SECRET" \
  -d "grant_type=refresh_token" \
  -d "refresh_token=REFRESH_TOKEN"
```

**响应示例**：

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "eyJhbGciOiJSUzI1NiIs...",
  "scope": "openid profile email"
}
```

## 令牌管理

### 令牌自省（Introspection）

验证令牌的有效性和权限：

```bash
curl -X POST http://localhost:8080/oauth2/introspect \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -u "CLIENT_ID:CLIENT_SECRET" \
  -d "token=ACCESS_TOKEN"
```

**响应示例**：

```json
{
  "active": true,
  "client_id": "client_abc123xyz",
  "token_type": "Bearer",
  "exp": 1672531200,
  "iat": 1672527600,
  "sub": "user_123",
  "scope": "openid profile email"
}
```

### 令牌吊销（Revocation）

主动吊销令牌（使其失效）：

```bash
curl -X POST http://localhost:8080/oauth2/revoke \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -u "CLIENT_ID:CLIENT_SECRET" \
  -d "token=ACCESS_OR_REFRESH_TOKEN" \
  -d "token_type_hint=refresh_token"
```

**注意**：
- `token_type_hint` 可选，值为 `access_token` 或 `refresh_token`
- 如果吊销 `refresh_token`，关联的 `access_token` 也会失效
- 如果吊销 `access_token`，不会影响 `refresh_token`

## 用户授权管理

### 查询授权记录

用户可以查看自己对所有客户端的授权记录：

```bash
curl -X GET http://localhost:8080/api/v1/oauth2/consents \
  -H "Authorization: Bearer USER_TOKEN"
```

**响应示例**：

```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "consents": [
    {
      "client_id": "client_abc123xyz",
      "client_name": "E-commerce App",
      "scopes": ["openid", "profile", "email"],
      "granted_at": 1672527600000
    }
  ]
}
```

### 撤销授权

用户可以撤销对特定客户端的授权：

```bash
curl -X DELETE http://localhost:8080/api/v1/oauth2/consents/client_abc123xyz \
  -H "Authorization: Bearer USER_TOKEN"
```

**效果**：
- 立即撤销该用户对客户端的授权
- 关联的 Access Token 和 Refresh Token 将失效
- 客户端需要重新引导用户完成授权流程

## 客户端管理

### 查看客户端列表

```bash
curl -X GET "http://localhost:8080/api/v1/oauth2/clients?is_active=true" \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 更新客户端配置

```bash
curl -X PUT http://localhost:8080/api/v1/oauth2/client_abc123xyz \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "client_name": "Updated App Name",
    "description": "新的描述",
    "redirect_uris": ["https://new.example.com/callback"],
    "scopes": ["openid", "profile"],
    "is_active": true
  }'
```

### 轮换客户端密钥

定期轮换密钥可以提高安全性：

```bash
curl -X POST http://localhost:8080/api/v1/oauth2/clients/client_abc123xyz/rotate-secret \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

**响应示例**：

```json
{
  "base_resp": {
    "code": 0,
    "message": "success"
  },
  "client_secret": "new_secret_xyz789"
}
```

**重要提示**：
- 旧密钥立即失效
- 新密钥仅此一次可见，请妥善保存
- 轮换后需更新所有使用该客户端的应用

### 删除客户端

**警告**：删除客户端不可逆！

```bash
curl -X DELETE http://localhost:8080/api/v1/oauth2/clients/client_abc123xyz \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## 最佳实践

### 安全建议

1. **使用 HTTPS**
   - 所有 OAuth2 通信必须通过 HTTPS
   - 授权码模式可以在 HTTP 环境测试，但生产环境必须用 HTTPS

2. **保护客户端密钥**
   - 机密客户端的 `client_secret` 必须安全存储
   - 不要在客户端代码（浏览器、移动应用）中硬编码密钥
   - 使用环境变量或密钥管理服务存储密钥

3. **使用 PKCE（推荐）**
   - 公共客户端应使用 PKCE（Proof Key for Code Exchange）
   - 防止授权码拦截攻击

4. **限制作用域**
   - 仅请求必要的最小权限
   - 不同环境使用不同的作用域配置

5. **验证 State 参数**
   - 授权流程中使用随机 `state` 参数
   - 防止 CSRF 攻击

6. **短期 Access Token**
   - Access Token 有效期建议 1 小时
   - 使用 Refresh Token 获取新令牌

### 令牌有效期建议

| 令牌类型 | 推荐有效期 | 说明 |
|---------|-----------|------|
| Access Token | 1 小时 | 短期有效，降低泄露风险 |
| Refresh Token | 30 天 | 长期有效，允许用户保持登录 |
| 授权码 | 10 分钟 | 即时使用，过期自动失效 |

### 客户端注册流程

1. **开发环境**
   - 创建测试客户端
   - 使用本地回调地址（如 `http://localhost:3000/callback`）

2. **生产环境**
   - 创建生产专用客户端
   - 使用 HTTPS 回调地址
   - 定期轮换密钥

3. **多环境隔离**
   - 为不同环境创建独立客户端
   - 避免跨环境使用客户端

### 错误处理

#### 常见错误码

| 错误 | HTTP 状态 | 说明 | 解决方案 |
|------|----------|------|----------|
| `invalid_client` | 401 | 客户端认证失败 | 检查 client_id 和 client_secret |
| `invalid_grant` | 400 | 授权码无效或过期 | 重新获取授权码 |
| `invalid_scope` | 400 | 作用域无效 | 检查请求的 scope |
| `access_denied` | 403 | 用户拒绝授权 | 提示用户重新授权 |
| `redirect_uri_mismatch` | 400 | 回调地址不匹配 | 确保回调地址在配置的 redirect_uris 中 |

#### 错误响应示例

```json
{
  "error": "invalid_client",
  "error_description": "Client authentication failed"
}
```

### 监控和审计

1. **日志记录**
   - 记录所有令牌颁发和吊销操作
   - 监控异常授权行为

2. **定期审计**
   - 审查客户端列表，删除不用的客户端
   - 检查用户授权记录
   - 轮换长期未换的密钥

## 故障排查

### 问题：授权码无效

**可能原因**：
- 授权码已过期（10 分钟有效期）
- 授权码已被使用
- 重定向地址不匹配

**解决方案**：
- 重新发起授权请求
- 确保 `redirect_uri` 与创建客户端时一致

### 问题：令牌自检返回 active=false

**可能原因**：
- 令牌已过期
- 令牌已被吊销
- 刷新令牌后，旧令牌失效

**解决方案**：
- 使用 `refresh_token` 获取新令牌
- 重新完成授权流程

### 问题：跨域请求失败

**可能原因**：
- 未配置 CORS
- 浏览器阻止跨域请求

**解决方案**：
- 在 OAuth2 服务端配置 CORS 头
- 使用代理服务器

## 相关文档

- [架构设计](../00-项目概览/架构设计.md) - 了解 OAuth2 在整体架构中的位置
- [权限管理设计](../04-权限管理/权限管理设计.md) - 了解 OAuth2 与 Casbin 的集成
- [API 网关日志规范](../02-开发规范/API网关日志规范.md) - OAuth2 相关日志编写规范

## 外部资源

- [RFC 6749 - OAuth 2.0 授权框架](https://datatracker.ietf.org/doc/html/rfc6749)
- [RFC 7009 - OAuth 2.0 令牌吊销](https://datatracker.ietf.org/doc/html/rfc7009)
- [RFC 7662 - OAuth 2.0 令牌自省](https://datatracker.ietf.org/doc/html/rfc7662)
- [OAuth 2.0 安全最佳实践](https://datatracker.ietf.org/doc/html/draft-ietf-oauth-security-topics)
- [Ory Fosite 文档](https://www.ory.sh/fosite/)
