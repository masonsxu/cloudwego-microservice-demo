# 部署指南

本文档介绍项目的部署方式和生产环境配置。

## 目录

- [容器部署](#容器部署)
- [部署命令参考](#部署命令参考)
- [生产环境配置](#生产环境配置)
- [生产环境检查清单](#生产环境检查清单)
- [服务端口参考](#服务端口参考)

---

## 容器部署

### 开发环境

```bash
# 1. 启动基础设施
cd docker && podman-compose up -d

# 2. 启动 RPC 服务（新终端）
cd rpc/identity_srv
cp .env.example .env
sh build.sh && sh output/bootstrap.sh

# 3. 启动网关服务（新终端）
cd gateway
cp .env.example .env
sh build.sh && sh output/bootstrap.sh

# 查看基础设施状态
podman-compose ps

# 查看日志
podman-compose logs -f
```

### 生产环境

生产环境建议使用容器编排工具（如 Kubernetes）部署应用服务，基础设施使用托管服务。

---

## 部署命令参考

### 基础设施管理

```bash
podman-compose up -d                    # 启动基础设施
podman-compose down                     # 停止基础设施
podman-compose restart                  # 重启基础设施
podman-compose ps                       # 查看服务状态
podman-compose down -v --remove-orphans # 清理环境（删除数据）
```

### 日志查看

```bash
podman-compose logs -f            # 所有日志
podman-compose logs -f postgres   # 特定服务日志
podman-compose logs -f redis      # Redis 日志
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
```

---

## 生产环境检查清单

### 安全配置

- [ ] 修改 `JWT_SIGNING_KEY` 为强随机密钥（至少 32 字符）
- [ ] 修改所有默认密码（数据库、Redis、对象存储）
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

## 服务端口参考

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

- [04-配置参考](04-configuration.md) - 详细配置参考
- [06-故障排查](06-troubleshooting.md) - 部署问题诊断
