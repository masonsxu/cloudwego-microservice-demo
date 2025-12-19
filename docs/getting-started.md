# 快速开始

本文档帮助你快速启动和运行 CloudWeGo Scaffold 项目。

## 目录

- [前置要求](#前置要求)
- [方式一：Docker 快速启动](#方式一docker-快速启动推荐)
- [方式二：本地开发模式](#方式二本地开发模式)
- [验证安装](#验证安装)
- [API 示例](#api-示例)
- [脚手架使用](#脚手架使用)

---

## 前置要求

- **Go**: 1.24+ ([下载](https://go.dev/dl/))
- **Docker**: 20.10+ ([安装](https://docs.docker.com/get-docker/))
- **Docker Compose**: 2.0+ ([安装](https://docs.docker.com/compose/install/))

## 方式一：Docker 快速启动（推荐）

这是最简单的启动方式，适合快速体验和开发。

```bash
# 1. 克隆项目
git clone <repository-url>
cd cloudwego-microservice-demo

# 2. 进入 docker 目录
cd docker

# 3. 复制环境配置（可选，默认配置已优化）
cp .env.dev.example .env

# 4. 启动所有服务（基础设施 + 应用）
./deploy.sh dev up

# 5. 查看服务状态
./deploy.sh dev ps

# 6. 查看日志
./deploy.sh dev logs              # 所有日志
./deploy.sh follow identity_srv   # 实时跟踪特定服务
```

服务启动后访问：

| 服务 | 地址 |
|------|------|
| HTTP API | http://localhost:8080 |
| Swagger 文档 | http://localhost:8080/swagger/index.html |
| 健康检查 | http://localhost:8080/health |

## 方式二：本地开发模式

适合需要调试单个服务或修改代码的场景。

### 1. 安装开发工具

```bash
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest
go install github.com/cloudwego/hertz/cmd/hz@latest
go install github.com/google/wire/cmd/wire@latest
```

### 2. 启动基础设施

```bash
cd docker
./deploy.sh dev up-base
```

### 3. 配置并启动 RPC 服务

```bash
cd rpc/identity_srv
cp .env.example .env
vim .env  # 修改数据库连接等配置

# 构建并启动
sh build.sh && sh output/bootstrap.sh
```

### 4. 配置并启动网关

```bash
# 在新终端
cd gateway
cp .env.example .env
vim .env

# 构建并启动
sh build.sh && sh output/bootstrap.sh
```

## 验证安装

```bash
# 健康检查
curl http://localhost:8080/health

# 预期输出
{"status":"ok"}
```

## API 示例

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/identity/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

响应示例：

```json
{
    "base_resp": {
        "code": 0,
        "message": "success",
    },
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "token_type": "Bearer",
        "expires_in": 1800
    }
}
```

### 获取用户信息（需要认证）

```bash
curl -X GET http://localhost:8080/api/v1/identity/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

---

## 脚手架使用

如果你将此项目作为脚手架使用，需要先初始化项目。

### 初始化项目

```bash
# 1. 克隆项目
git clone <repository-url>
cd cloudwego-microservice-demo

# 2. 运行初始化脚本（替换为你的 module 路径）
./scripts/init.sh github.com/your-org/your-project https://github.com/masonsxu/cloudwego-microservice-demo.git

# 3. 更新依赖
go mod tidy

# 4. 重新生成代码
# RPC 服务
cd rpc/identity_srv
./script/gen_kitex_code.sh
cd wire && wire && cd ../..

# HTTP 网关
cd ../gateway
./script/gen_hertz_code.sh
cd internal/wire && wire && cd ../..
```

**为什么需要初始化？**

Go module 路径是项目的唯一标识符，必须与你的实际项目路径匹配。初始化脚本会自动替换：

- 所有 `go.mod` 文件中的 module 声明
- 所有 Go 源代码中的 import 路径
- 脚本文件中的 module 路径引用

### 从上游更新

如果脚手架本身有更新（bug 修复、新功能等），可以从上游仓库拉取：

```bash
# 方法一：使用更新脚本（推荐）
./scripts/update.sh main

# 方法二：手动更新
git remote add upstream https://github.com/masonsxu/cloudwego-microservice-demo.git
git fetch upstream main
git checkout -b update-from-scaffold
git merge upstream/main

# 解决冲突后
go mod tidy
cd rpc/identity_srv && ./script/gen_kitex_code.sh && cd wire && wire && cd ../../..
cd gateway && ./script/gen_hertz_code.sh && cd internal/wire && wire && cd ../..
```

### 更新时的常见冲突处理

1. **go.mod 冲突**：保留你的 module 路径，接受依赖版本更新
2. **配置文件冲突**：保留你的配置值，接受新的配置项
3. **业务代码冲突**：根据实际情况决定

### 更新检查清单

- [ ] 已提交或暂存当前更改
- [ ] 已设置 upstream 远程仓库
- [ ] 合并后检查是否有冲突
- [ ] 运行 `go mod tidy` 更新依赖
- [ ] 重新生成代码（如果 IDL 有更新）
- [ ] 重新生成 Wire 依赖注入代码
- [ ] 运行测试确保一切正常
- [ ] 检查构建是否成功

---

## 下一步

- [架构设计](architecture.md) - 了解项目整体架构
- [开发指南](development-guide.md) - 开始开发新功能
- [配置说明](configuration.md) - 详细的配置参考
