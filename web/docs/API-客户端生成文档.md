# API 客户端自动生成文档

## 概述

本项目使用 **openapi-typescript-codegen** 从后端 Swagger 文档自动生成前端 API 客户端和 TypeScript 类型定义，确保前后端类型一致性。

## 原理

```
后端 Swagger 文档 (gateway/docs/swagger.yaml)
        ↓
openapi-typescript-codegen
        ↓
前端 API 客户端 (web/src/api/generated/)
```

## 优势

1. **类型安全**：自动从 Swagger 生成 TypeScript 类型，100% 类型安全
2. **同步更新**：后端 API 变更后，只需重新运行生成脚本即可同步
3. **前后端分离**：前端开发者不需要接触后端源码
4. **减少手写代码**：减少 API 调用代码的编写和维护成本

## 生成 API 客户端

### 前提条件

1. 后端服务已启动
2. Swagger 文档已生成到 `gateway/docs/swagger.yaml`
3. 前端项目已安装依赖 `openapi-typescript-codegen`

### 生成步骤

在项目根目录执行：

```bash
cd web
bash scripts/generate-api-from-swagger.sh
```

### 生成文件结构

```
src/api/generated/
├── core/               # 核心功能
│   ├── ApiError.ts
│   ├── ApiResult.ts
│   ├── request.ts     # HTTP 请求客户端
│   └── ...
├── models/            # TypeScript 类型定义
│   ├── LoginResponseDTO.ts
│   ├── UserProfileDTO.ts
│   └── ...
├── services/          # API 服务类
│   └── Service.ts
└── index.ts           # 导出文件
```

## 使用方法

### 1. 基本用法

```typescript
import { Service } from '@/api/generated'
import type { LoginRequestDTO, LoginResponseDTO } from '@/api/generated/models'

// 创建 API 实例
const api = new Service({
  baseURL: 'http://localhost:8080'
})

// 调用登录接口
async function login(username: string, password: string) {
  const response = await api.identity.login({
    requestBody: {
      username,
      password
    }
  })
  
  console.log('用户信息:', response.userProfile)
  console.log('访问令牌:', response.tokenInfo?.accessToken)
}
```

### 2. 与现有系统集成

#### 替代手写的 API 调用

**之前（手写）：**
```typescript
// src/api/auth.ts
export function login(data: LoginRequest) {
  return request<LoginResponseData>({
    url: '/api/v1/identity/auth/login',
    method: 'POST',
    data
  })
}
```

**现在（自动生成）：**
```typescript
import { Service } from '@/api/generated'

export function login(data: LoginRequestDTO) {
  const api = new Service({
    baseURL: import.meta.env.VITE_API_BASE_URL
  })
  
  return api.identity.login({
    requestBody: { data }
  })
}
```

#### 使用生成的类型

```typescript
import type {
  LoginResponseDTO,
  UserProfileDTO,
  TokenInfoDTO,
  MenuPermissionDTO
} from '@/api/generated/models'

// 类型安全的数据处理
function handleLogin(response: LoginResponseDTO) {
  const user: UserProfileDTO = response.userProfile
  const token: string = response.tokenInfo?.accessToken || ''
  const permissions: MenuPermissionDTO[] = response.permissions || []
  
  // TypeScript 会自动检查类型
  console.log(`用户 ${user.username} 登录成功`)
}
```

### 3. 配置拦截器

为生成的 API 客户端添加拦截器（如添加 Token）：

```typescript
// src/api/generated/request.ts
import { ApiError } from './core/ApiError'
import type { ApiRequestOptions } from './core/ApiRequestOptions'
import { useAuthStore } from '@/stores/auth'

export async function request<T>(
  method: string,
  path: string,
  options: ApiRequestOptions = {}
): Promise<T> {
  const authStore = useAuthStore()
  
  // 添加 Token
  if (authStore.token) {
    options.headers = {
      ...options.headers,
      Authorization: `Bearer ${authStore.token}`
    }
  }
  
  // 错误处理
  try {
    return await (window as any).openApiClient.request(method, path, options)
  } catch (error) {
    if (error instanceof ApiError) {
      console.error('API Error:', error.status, error.message)
    }
    throw error
  }
}
```

## 常见问题

### Q1: 重新生成会覆盖已修改的文件吗？

A: 是的，`src/api/generated/` 目录下的所有文件都是自动生成的，不应手动修改。如果需要自定义，请创建新的文件在 `src/api/` 目录下。

### Q2: 如何处理后端 API 变更？

A: 
1. 确保后端已更新 Swagger 文档
2. 重新运行 `bash scripts/generate-api-from-swagger.sh`
3. TypeScript 编译器会检查类型错误，根据错误提示修改调用代码

### Q3: 生成的 API 客户端路径与实际不符？

A: 
1. 检查 Swagger 文件的 `servers` 配置
2. 或在使用时传入 `baseURL`

```typescript
const api = new Service({
  baseURL: 'http://192.168.20.66:8088'  // 使用实际的地址
})
```

## 完整示例

```typescript
// src/api/auth.generated.ts
import { Service } from '@/api/generated'
import type { LoginRequestDTO, LoginResponseDTO } from '@/api/generated/models'

const api = new Service({
  baseURL: import.meta.env.VITE_API_BASE_URL
})

export async function login(credentials: LoginRequestDTO): Promise<LoginResponseDTO> {
  const response = await api.identity.login({
    requestBody: { credentials }
  })
  
  return response
}
```

## 维护建议

1. **定期同步**：后端 API 更新后及时重新生成
2. **版本控制**：将 `src/api/generated/` 添加到 `.gitignore`，让 CI/CD 自动生成
3. **文档优先**：后端修改 API 时，先更新 Swagger 文档

## 相关文档

- [openapi-typescript-codegen 官方文档](https://openapi-ts.dev/)
- [后端 Swagger 文档](../../gateway/docs/swagger.yaml)
- [API 开发指南](../02-开发规范/开发指南.md)
