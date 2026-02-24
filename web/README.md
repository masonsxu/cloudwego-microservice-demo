# CloudWeGo 微服务管理平台 - 前端

基于 Vue 3 + TypeScript + Element Plus + Tailwind CSS 构建的现代化前端管理平台，采用奢华摩羯座配色主题。

## 🚀 快速开始

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:5173

### 生产构建

```bash
npm run build
```

## 📦 从 Swagger 文档生成 API 客户端

本项目采用 **API-First** 开发模式，从后端 Swagger 文档自动生成前端 API 客户端和类型定义，确保前后端类型一致性。

### 原理

```
后端 Swagger 文档 → openapi-typescript-codegen → 前端 API 客户端
```

### 生成 API 客户端

```bash
cd web
bash scripts/generate-api_from_swagger.sh
```

### 使用生成的 API

```typescript
import { Service } from '@/api/generated'
import type { LoginResponseDTO } from '@/api/generated/models'

const api = new Service({
  baseURL: import.meta.env.VITE_API_BASE_URL
})

const response = await api.identity.login({
  requestBody: { 
    username: 'admin', 
    password: 'password' 
  }
})
```

详细文档请参考：[API 客户端生成文档](./docs/API-客户端生成文档.md)

## 🎨 配色主题 - 奢华摩羯座

- **深岩灰** (#141416): 背景基础
- **香槟金** (#D4AF37): 核心高亮、图标
- **羊皮纸白** (#F2F0E4): 主标题、文本
- **矿石灰** (#8B9bb4): 副文本
- **青铜褐** (#2C2E33): 按钮、卡片底色

## 📁 项目结构

```
web/
├── src/
│   ├── api/                    # API 调用层
│   │   ├── generated/         # 从 Swagger 自动生成（勿手动修改）
│   │   ├── auth.ts           # 认证相关 API
│   │   ├── user.ts           # 用户管理 API
│   │   ├── organization.ts   # 组织管理 API
│   │   └── role.ts           # 角色权限 API
│   ├── assets/styles/         # 样式文件
│   ├── components/            # 公共组件
│   │   └── layout/            # 布局组件
│   ├── views/                 # 页面组件
│   ├── stores/                # Pinia 状态管理
│   ├── router/                # 路由配置
│   ├── types/                 # TypeScript 类型
│   ├── locales/               # 国际化资源
│   └── main.ts                # 入口文件
├── scripts/                  # 构建和生成脚本
├── docs/                     # 项目文档
└── package.json
```

## 🔧 开发脚本

```bash
# 从 Swagger 生成 API 客户端
bash scripts/generate-api-from-swagger.sh

# 开发
npm run dev

# 构建
npm run build

# 类型检查
npm run type-check

# 代码检查
npm run lint
```

## 📝 后端 Swagger 文档

- **位置**: `../gateway/docs/swagger.yaml`
- **在线访问**: http://192.168.20.66:8088/swagger/index.html

## 📚 相关文档

- [API 客户端生成文档](./docs/API-客户端生成文档.md)
- [项目架构设计](../docs/00-项目概览/架构设计.md)
- [开发指南](../docs/02-开发规范/开发指南.md)

## 🏗️ 技术栈

- **框架**: Vue 3.4+ (Composition API)
- **构建工具**: Vite 5.x
- **语言**: TypeScript 5.x
- **UI 组件**: Element Plus
- **样式**: Tailwind CSS 4.x + SCSS
- **路由**: Vue Router 4.x
- **状态管理**: Pinia 2.x
- **国际化**: Vue I18n 9.x
- **HTTP**: Axios
- **API 生成**: openapi-typescript-codegen

## 📄 许可证

MIT License
