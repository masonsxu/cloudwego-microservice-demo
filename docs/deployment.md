# 部署指南

本文档介绍 CloudWeGo Scaffold 的部署方式和生产环境配置。

## 目录

- [Docker 部署](#docker-部署)
- [部署命令参考](#部署命令参考)
- [生产环境配置](#生产环境配置)
- [生产环境检查清单](#生产环境检查清单)

---

## Docker 部署

### 开发环境

```bash
cd docker

# 启动所有服务
./deploy.sh dev up

# 查看状态
./deploy.sh dev ps

# 查看日志
./deploy.sh dev logs
```

### 生产环境

```bash
cd docker

# 1. 配置生产环境变量
cp .env.dev.example .env.prod
vim .env.prod  # 修改为生产配置

# 2. 构建生产镜像
./deploy.sh prod build

# 3. 启动生产环境
./deploy.sh prod up -d

# 4. 查看状态
./deploy.sh prod ps

# 5. 查看日志
./deploy.sh prod logs
```

---

## 部署命令参考

### 服务管理

```bash
# 启动所有服务
./deploy.sh dev up

# 仅启动基础设施（postgres, etcd, rustfs）
./deploy.sh dev up-base

# 仅启动应用服务（identity_srv, gateway）
./deploy.sh dev up-apps

# 后台启动
./deploy.sh dev up -d

# 停止所有服务
./deploy.sh dev down
```

### 日志查看

```bash
# 所有日志
./deploy.sh dev logs

# 特定服务日志
./deploy.sh dev logs identity_srv
./deploy.sh dev logs gateway

# 实时跟踪日志
./deploy.sh follow identity_srv
```

### 镜像构建

```bash
# 重新构建所有镜像
./deploy.sh dev rebuild

# 构建特定服务
./deploy.sh dev build identity_srv
./deploy.sh dev build gateway
```

---

## 生产环境配置

### 必须修改的配置

#### 1. JWT 签名密钥

```env
# 使用强随机密钥（至少 32 字符）
JWT_SIGNING_KEY=your-very-long-and-secure-random-string-here
```

生成强密钥：
```bash
openssl rand -base64 32
```

#### 2. 数据库密码

```env
DB_PASSWORD=your-secure-database-password
```

#### 3. 启用安全选项

```env
# 启用 HTTPS Cookie（需要 HTTPS）
JWT_COOKIE_SECURE_COOKIE=true

# 启用数据库 SSL
DB_SSLMODE=require
```

#### 4. 关闭调试模式

```env
SERVER_DEBUG=false
LOG_LEVEL=info
```

### 推荐的生产配置

```env
# 服务配置
SERVER_DEBUG=false
LOG_LEVEL=info

# 数据库配置
DB_HOST=postgres
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=your-secure-password
DB_SSLMODE=require
DB_MAX_IDLE_CONNS=20
DB_MAX_OPEN_CONNS=200
DB_CONN_MAX_LIFETIME=1h

# Redis 配置
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_POOL_SIZE=20
REDIS_MIN_IDLE_CONNS=5
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s

# JWT 配置
JWT_ENABLED=true
JWT_SIGNING_KEY=your-32-character-or-longer-secret-key
JWT_TIMEOUT=30m
JWT_MAX_REFRESH=168h
JWT_COOKIE_HTTP_ONLY=true
JWT_COOKIE_SECURE_COOKIE=true

# etcd 配置
ETCD_ADDRESS=etcd:2379
ETCD_TIMEOUT=5s

# OpenTelemetry 配置
OTEL_ENABLED=true
OTEL_EXPORTER_ENDPOINT=jaeger:4317
OTEL_EXPORTER_INSECURE=false
OTEL_SAMPLER_RATIO=0.1
OTEL_RESOURCE_ATTRIBUTES=service.version=1.0.0,deployment.environment=production
```

---

## 生产环境检查清单

### 安全配置

- [ ] 修改 `JWT_SIGNING_KEY` 为强随机密钥（至少 32 字符）
- [ ] 修改所有默认密码（数据库、对象存储等）
- [ ] 启用 `JWT_COOKIE_SECURE_COOKIE=true`（需要 HTTPS）
- [ ] 启用 `JWT_COOKIE_HTTP_ONLY=true`（防止 XSS）
- [ ] 设置 `DB_SSLMODE=require` 或 `verify-full`

### 性能配置

- [ ] 设置 `SERVER_DEBUG=false`
- [ ] 设置 `LOG_LEVEL=info` 或 `warn`
- [ ] 调整数据库连接池参数
- [ ] 配置适当的超时时间

### 运维配置

- [ ] 配置防火墙规则
- [ ] 设置数据库备份策略
- [ ] 配置监控和告警
- [ ] 设置日志轮转

### 网络配置

- [ ] 配置 HTTPS（推荐使用反向代理如 Nginx）
- [ ] 配置域名和 DNS
- [ ] 设置适当的 CORS 策略

---

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| gateway | 8080 | HTTP API |
| identity_srv | 8891 | RPC 服务 |
| identity_srv | 10000 | 健康检查 |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存 |
| etcd | 2379 | 服务发现 |
| RustFS | 9000 | 对象存储 |
| Jaeger | 16686 | 链路追踪 UI |
| Jaeger Collector | 4317 | OTLP 接收 |

---

## 下一步

- [配置说明](configuration.md) - 详细配置参考
- [故障排查](troubleshooting.md) - 部署问题诊断
