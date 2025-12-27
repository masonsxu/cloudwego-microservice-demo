# 故障排查

本文档收集常见问题及其解决方案。

## 目录

- [服务启动问题](#服务启动问题)
- [代码生成问题](#代码生成问题)
- [运行时问题](#运行时问题)
- [Docker 环境问题](#docker-环境问题)
- [性能问题](#性能问题)
- [日志和调试](#日志和调试)
- [常见问题 FAQ](#常见问题-faq)

---

## 服务启动问题

### 端口已被占用

**错误信息**：

```
bind: address already in use
listen tcp :8891: bind: address already in use
```

**解决方法**：

```bash
# 查找并终止占用端口的进程
kill -9 $(lsof -t -i:8891)

# 或修改服务端口（在 .env 文件中）
SERVER_ADDRESS=:8893

# Docker 环境清理
cd docker && ./deploy.sh down && ./deploy.sh up
```

### 数据库连接失败

**错误信息**：

```
failed to connect to database
dial tcp 127.0.0.1:5432: connect: connection refused
```

**诊断方法**：

```bash
# 检查 PostgreSQL 是否运行
cd docker && ./deploy.sh ps

# 检查端口监听
lsof -i :5432

# 测试连接
psql -h localhost -p 5432 -U postgres -d identity_srv
```

**解决方法**：

1. 启动数据库：`cd docker && ./deploy.sh up-base`
2. 检查 `.env` 配置中的 `DB_HOST`、`DB_PORT`、`DB_PASSWORD`
3. 创建数据库（如果不存在）：
   ```sql
   CREATE DATABASE identity_srv;
   ```

### etcd 连接失败

**错误信息**：

```
failed to connect to etcd
context deadline exceeded
```

**诊断方法**：

```bash
# 检查 etcd 状态
cd docker && ./deploy.sh ps

# 测试连接
curl http://localhost:2379/version
```

**解决方法**：

1. 启动 etcd：`cd docker && ./deploy.sh up-base`
2. 检查 `ETCD_ADDRESS` 配置

---

## 代码生成问题

### kitex/hz 命令未找到

**错误信息**：

```
command not found: kitex
command not found: hz
```

**解决方法**：

```bash
# 安装工具
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
go install github.com/cloudwego/thriftgo@latest
go install github.com/cloudwego/hertz/cmd/hz@latest

# 确保 PATH 包含 GOPATH/bin
export PATH=$PATH:$(go env GOPATH)/bin

# 验证
kitex --version
hz --version
```

### IDL 文件未找到

**错误信息**：

```
open ../../idl/rpc/identity_srv/identity_service.thrift: no such file or directory
```

**解决方法**：

```bash
# 确认在正确目录执行
cd rpc/identity_srv
./script/gen_kitex_code.sh

# 检查 IDL 文件是否存在
ls -la ../../idl/rpc/identity_srv/
```

### Wire 生成失败

**错误信息**：

```
wire: no provider found for *gorm.DB
wire: cycle detected in provider set
```

**解决方法**：

1. 检查 Provider 函数参数和返回值类型
2. 确保所有依赖都有对应的 Provider
3. 重新生成：
   ```bash
   cd wire
   rm wire_gen.go
   wire
   ```

---

## 运行时问题

### Redis 连接失败

**错误信息**：

```
dial tcp 127.0.0.1:6379: connect: connection refused
```

**诊断方法**：

```bash
# 检查 Redis 是否运行
cd docker && ./deploy.sh ps

# 测试连接
redis-cli -h localhost -p 6379 ping
```

**解决方法**：

1. 启动 Redis：`cd docker && ./deploy.sh up-base`
2. 检查 `.env` 配置中的 `REDIS_HOST`、`REDIS_PORT`

### RPC 调用超时

**错误信息**：

```
rpc timeout: deadline exceeded
```

**诊断方法**：

```bash
# 检查服务是否运行
cd docker && ./deploy.sh ps

# 查看日志
docker logs -f backend-identity_srv-1
```

**解决方法**：

1. 增加超时时间（gateway `.env`）：
   ```env
   CLIENT_REQUEST_TIMEOUT=60s
   ```
2. 检查服务注册：
   ```bash
   docker exec etcd etcdctl get --prefix /kitex
   ```

### JWT 认证失败

**错误信息**：

```
401 Unauthorized
token is expired
```

**解决方法**：

1. 刷新 Token：
   ```bash
   curl -X POST http://localhost:8080/api/v1/identity/auth/refresh \
     -H "Authorization: Bearer YOUR_OLD_TOKEN"
   ```
2. 确保签名密钥一致
3. 检查 Token 格式：`Authorization: Bearer <token>`

---

## Docker 环境问题

### 容器启动后立即退出

**诊断方法**：

```bash
# 查看容器状态
docker ps -a | grep cloudwego

# 查看退出码
docker inspect identity-srv | grep ExitCode

# 查看日志
./deploy.sh logs identity_srv
```

**常见原因**：

1. 配置错误 - 检查 `.env` 文件
2. 依赖服务未就绪 - 先启动基础服务
3. 端口冲突 - 修改端口映射

### 无法连接容器内服务

**解决方法**：

```bash
# 检查网络
docker network ls

# 容器内部使用服务名
DB_HOST=postgres       # 而不是 127.0.0.1
ETCD_ADDRESS=etcd:2379 # 而不是 localhost:2379

# 宿主机使用 localhost
DB_HOST=127.0.0.1
```

---

## 性能问题

### 数据库查询慢

**诊断方法**：

```bash
# 启用 SQL 日志
LOG_LEVEL=debug
SERVER_DEBUG=true
```

**优化方法**：

1. 添加数据库索引
2. 使用 `Preload` 避免 N+1 查询
3. 调整连接池：
   ```env
   DB_MAX_IDLE_CONNS=20
   DB_MAX_OPEN_CONNS=200
   ```

### 内存占用过高

**诊断方法**：

```bash
# 查看容器资源
docker stats
```

**解决方法**：

1. 检查内存泄漏
2. 调整 GORM 连接池
3. 使用分页查询

---

## 日志和调试

### 查看实时日志

```bash
# 基础设施日志
cd docker && ./deploy.sh logs
./deploy.sh logs postgres       # PostgreSQL 日志
./deploy.sh logs redis          # Redis 日志

# 应用服务日志（查看终端输出）
# RPC 服务和网关服务日志直接输出到启动终端
```

### 启用详细日志

```env
LOG_LEVEL=debug
LOG_FORMAT=text
SERVER_DEBUG=true
```

### 追踪请求

```bash
# 使用 request_id 搜索日志
# 在应用服务终端中搜索包含特定 request_id 的日志

# 使用 Jaeger UI 查看链路追踪
# 访问 http://localhost:16686
```

---

## 常见问题 FAQ

### Q: 为什么不使用 config.yaml 文件？

A: 项目采用环境变量驱动配置，便于容器化部署和环境隔离。

### Q: Duration 类型环境变量如何设置？

A: 支持多种格式：`1h`、`30m`、`3600s` 或纯数字 `3600`。

### Q: 如何添加新的数据模型？

A:
1. 在 `models/` 目录创建 GORM 模型
2. 在 `config/database.go` 的 `AutoMigrate` 中添加模型
3. 重启服务

### Q: Wire 生成失败怎么办？

A:
1. 检查 Provider 函数签名
2. 确保所有依赖都有 Provider
3. 删除 `wire_gen.go` 后重新生成

### Q: JWT Token 如何安全存储？

A:
1. 使用 HttpOnly Cookie 防止 XSS
2. 生产环境启用 Secure Cookie
3. 设置合理的过期时间
4. 实现 Token 刷新机制

---

## 获取帮助

如果以上方案无法解决问题：

1. 检查 [GitHub Issues](https://github.com/masonsxu/cloudwego-microservice-demo/issues)
2. 提交新 Issue 并附上：
   - 错误信息
   - 复现步骤
   - 环境信息（Go 版本、Docker 版本等）
