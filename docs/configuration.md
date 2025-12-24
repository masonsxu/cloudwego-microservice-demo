# 配置说明

本文档详细介绍 CloudWeGo Scaffold 的配置管理方式和所有可用配置项。

## 目录

- [配置管理方式](#配置管理方式)
- [配置文件位置](#配置文件位置)
- [通用配置](#通用配置)
- [数据库配置](#数据库配置)
- [服务注册发现](#服务注册发现)
- [JWT 认证配置](#jwt-认证配置)
- [文件存储配置](#文件存储配置)
- [环境差异对照](#环境差异对照)

---

## 配置管理方式

项目采用 **环境变量驱动配置** 模式，不使用 YAML 文件。

### 配置优先级（从高到低）

1. **系统环境变量**（最高优先级）
2. **`.env` 文件**（环境变量未设置时加载）
3. **`config/defaults.go` 默认值**（最低优先级）

### Duration 类型格式

Duration 类型支持多种格式：
- 带单位：`1h`、`30m`、`3600s`
- 纯数字：`3600`（默认按秒解析）

---

## 配置文件位置

```
项目根目录/
├── docker/
│   └── .env.dev.example      # Docker 开发环境配置模板
├── gateway/
│   └── .env.example          # HTTP 网关配置模板
└── rpc/identity_srv/
    └── .env.example          # 身份认证服务配置模板
```

### 首次配置

```bash
# Docker 环境
cd docker && cp .env.dev.example .env

# 本地开发 - identity_srv
cd rpc/identity_srv && cp .env.example .env

# 本地开发 - gateway
cd gateway && cp .env.example .env
```

---

## 通用配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `SERVER_DEBUG` | 调试模式 | `false` | `true` |
| `LOG_LEVEL` | 日志级别 | `info` | `debug`/`info`/`warn`/`error` |
| `LOG_FORMAT` | 日志格式 | `json` | `json`/`text` |
| `LOG_OUTPUT` | 日志输出 | `stdout` | `stdout`/`file` |

---

## 数据库配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `DB_HOST` | 数据库主机 | `127.0.0.1` | `postgres`（Docker） |
| `DB_PORT` | 数据库端口 | `5432` | `5432` |
| `DB_USERNAME` | 数据库用户名 | `postgres` | `postgres` |
| `DB_PASSWORD` | 数据库密码 | - | `your-password` |
| `DB_NAME` | 数据库名称 | - | `identity_srv` |
| `DB_SSLMODE` | SSL 模式 | `disable` | `disable`/`require`/`verify-full` |
| `DB_TIMEZONE` | 时区 | `Asia/Shanghai` | `Asia/Shanghai` |

### 连接池配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `DB_MAX_IDLE_CONNS` | 最大空闲连接数 | `10` | `10` |
| `DB_MAX_OPEN_CONNS` | 最大打开连接数 | `100` | `100` |
| `DB_CONN_MAX_LIFETIME` | 连接最大生命周期 | `1h` | `1h`/`60m`/`3600` |
| `DB_CONN_MAX_IDLE_TIME` | 连接最大空闲时间 | `5m` | `5m`/`300` |

---

## Redis 缓存配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `REDIS_HOST` | Redis 主机 | `127.0.0.1` | `redis`（Docker） |
| `REDIS_PORT` | Redis 端口 | `6379` | `6379` |
| `REDIS_PASSWORD` | Redis 密码 | - | `your-redis-password` |
| `REDIS_DB` | Redis 数据库 | `0` | `0` |
| `REDIS_POOL_SIZE` | 连接池大小 | `10` | `10` |
| `REDIS_MIN_IDLE_CONNS` | 最小空闲连接 | `5` | `5` |
| `REDIS_DIAL_TIMEOUT` | 连接超时 | `5s` | `5s` |
| `REDIS_READ_TIMEOUT` | 读取超时 | `3s` | `3s` |
| `REDIS_WRITE_TIMEOUT` | 写入超时 | `3s` | `3s` |

### Redis 使用场景

- **会话管理**: JWT Token 黑名单、用户会话状态
- **热点数据缓存**: 用户信息、权限数据
- **限流计数**: API 调用频率限制
- **分布式锁**: 防止重复操作

---

## 服务注册发现

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `ETCD_ADDRESS` | etcd 地址 | `127.0.0.1:2379` | `etcd:2379`（Docker） |
| `ETCD_USERNAME` | etcd 用户名 | - | - |
| `ETCD_PASSWORD` | etcd 密码 | - | - |
| `ETCD_TIMEOUT` | 连接超时 | `5s` | `5`/`5s` |

---

## JWT 认证配置

仅适用于 **gateway** 服务。

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `JWT_ENABLED` | 启用 JWT | `true` | `true` |
| `JWT_SIGNING_KEY` | 签名密钥 | - | `your-secret-key` |
| `JWT_TIMEOUT` | Token 有效期 | `30m` | `30m` |
| `JWT_MAX_REFRESH` | 最大刷新时间 | `168h` | `168h`（7天） |

### 跳过认证的路径

```env
JWT_SKIP_PATHS=/api/v1/identity/auth/login,/api/v1/identity/auth/refresh,/ping,/health
```

### Cookie 配置

| 变量名 | 说明 | 默认值 | 生产环境 |
|--------|------|--------|----------|
| `JWT_COOKIE_SEND_COOKIE` | 发送 Cookie | `true` | `true` |
| `JWT_COOKIE_HTTP_ONLY` | HttpOnly | `true` | `true` |
| `JWT_COOKIE_SECURE_COOKIE` | Secure（需 HTTPS） | `false` | `true` |
| `JWT_COOKIE_COOKIE_SAME_SITE` | SameSite | `lax` | `lax` |

---

## 文件存储配置

仅适用于 **identity_srv** 服务（组织 Logo 存储）。

### 双端点配置

为解决容器化部署中内外部访问问题，采用双端点配置：

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `LOGO_STORAGE_S3_ENDPOINT` | 内部端点（容器间通信） | `http://rustfs:9000` |
| `LOGO_STORAGE_S3_PUBLIC_ENDPOINT` | 公共端点（浏览器访问） | `http://localhost:9000` |

### 其他存储配置

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `LOGO_STORAGE_S3_REGION` | S3 区域 | `us-east-1` |
| `LOGO_STORAGE_S3_USE_SSL` | 使用 SSL | `false` |
| `LOGO_STORAGE_ACCESS_KEY` | 访问密钥 | `RustFSadmin` |
| `LOGO_STORAGE_SECRET_KEY` | 私钥 | - |
| `LOGO_STORAGE_MAX_FILE_SIZE` | 最大文件大小（字节） | `10485760`（10MB） |
| `LOGO_STORAGE_ALLOWED_FILE_TYPES` | 允许的文件类型 | `image/jpeg,image/png,...` |

---

## OpenTelemetry 链路追踪配置

| 变量名 | 说明 | 默认值 | 示例 |
|--------|------|--------|------|
| `OTEL_ENABLED` | 启用链路追踪 | `true` | `true` |
| `OTEL_SERVICE_NAME` | 服务名称 | - | `identity-srv` |
| `OTEL_EXPORTER_ENDPOINT` | Collector 端点 | `localhost:4317` | `jaeger:4317` |
| `OTEL_EXPORTER_INSECURE` | 跳过 TLS 验证 | `true` | `false`（生产环境） |
| `OTEL_SAMPLER_RATIO` | 采样率 | `1.0` | `0.1`（生产环境） |
| `OTEL_RESOURCE_ATTRIBUTES` | 资源属性 | - | `service.version=1.0.0` |

---

## 环境差异对照

| 配置项 | 开发环境 | 生产环境 |
|--------|----------|----------|
| `SERVER_DEBUG` | `true` | `false` |
| `LOG_LEVEL` | `debug` | `info`/`warn` |
| `DB_PASSWORD` | 简单密码 | **强密码** |
| `REDIS_PASSWORD` | 无或简单 | **强密码** |
| `JWT_SIGNING_KEY` | 简单字符串 | **32+ 字符强密钥** |
| `JWT_COOKIE_SECURE_COOKIE` | `false` | `true`（需 HTTPS） |
| `DB_SSLMODE` | `disable` | `require`/`verify-full` |
| `OTEL_SAMPLER_RATIO` | `1.0` | `0.1` |
| `OTEL_EXPORTER_INSECURE` | `true` | `false` |

---

## 添加新配置项

### 1. 定义配置结构

```go
// config/types.go
type Config struct {
    // ...
    NewFeature NewFeatureConfig `mapstructure:"new_feature"`
}

type NewFeatureConfig struct {
    Enabled bool   `mapstructure:"enabled"`
    Timeout int    `mapstructure:"timeout"`
}
```

### 2. 设置默认值

```go
// config/defaults.go
func setDefaults(v *viper.Viper) {
    // ...
    v.SetDefault("new_feature.enabled", false)
    v.SetDefault("new_feature.timeout", 30)
}
```

### 3. 添加环境变量映射

```go
// config/env.go
func loadEnvVariables(v *viper.Viper) {
    // ...
    mapToViper(v, "NEW_FEATURE_ENABLED", "new_feature.enabled", nil)
    mapToViper(v, "NEW_FEATURE_TIMEOUT", "new_feature.timeout", nil)
}
```

### 4. 更新 .env.example

```env
# 新功能配置
NEW_FEATURE_ENABLED=false
NEW_FEATURE_TIMEOUT=30
```

---

## 配置验证清单

部署前请确认：

- [ ] 数据库连接配置正确
- [ ] 每个服务的 `DB_NAME` 不同
- [ ] etcd 地址指向正确的服务注册中心
- [ ] 服务端口未被占用
- [ ] 生产环境已修改 `JWT_SIGNING_KEY` 为强密钥
- [ ] 生产环境已启用安全选项

---

## 下一步

- [部署指南](deployment.md) - 生产环境部署
- [故障排查](troubleshooting.md) - 配置问题诊断
