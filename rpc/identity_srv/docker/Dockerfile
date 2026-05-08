# =============================================================================
# identity_srv RPC 服务镜像
# =============================================================================
# 构建命令（必须在项目根目录执行，因为 identity_srv 通过 replace 引用同 monorepo
# 内的 iamclient 与 rpc/policy_srv 模块）：
#   podman build -f rpc/identity_srv/docker/Dockerfile -t identity-srv:latest .
#
# 说明：
#   - 构建上下文必须是项目根，才能同时拷贝 rpc/identity_srv/、iamclient/、rpc/policy_srv/
#   - 通过 GOWORK=off 显式禁用 go.work，依靠 go.mod 中的 replace 解析模块路径
# =============================================================================

# -----------------------------------------------------------------------------
# 构建阶段
# -----------------------------------------------------------------------------
FROM m.daocloud.io/docker.io/library/golang:1.25 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct \
    GOWORK=off

WORKDIR /build

# 利用层缓存：先 COPY 各 module 的 go.mod/go.sum 下载依赖（保留相对路径让 replace 能解析）
COPY rpc/identity_srv/go.mod rpc/identity_srv/go.sum ./rpc/identity_srv/
COPY rpc/policy_srv/go.mod rpc/policy_srv/go.sum ./rpc/policy_srv/
COPY iamclient/go.mod iamclient/go.sum ./iamclient/
RUN cd rpc/identity_srv && go mod download

# 复制全部源码（含 replace 指向的 iamclient、rpc/policy_srv）
COPY rpc/identity_srv/ ./rpc/identity_srv/
COPY rpc/policy_srv/ ./rpc/policy_srv/
COPY iamclient/ ./iamclient/

WORKDIR /build/rpc/identity_srv
RUN go build -ldflags="-s -w" -trimpath -o identity_srv .

# -----------------------------------------------------------------------------
# 运行阶段
# -----------------------------------------------------------------------------
FROM m.daocloud.io/docker.io/library/alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl
ENV TZ=Asia/Shanghai

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# 拷贝二进制和运行时需要的配置文件
COPY --from=builder /build/rpc/identity_srv/identity_srv .
COPY --from=builder /build/rpc/identity_srv/config/seed_menu.yaml ./config/seed_menu.yaml

# 准备日志目录并修正归属
RUN mkdir -p /app/logs && chown -R appuser:appuser /app

USER appuser

# 8891: Kitex RPC 端口；10000: 健康检查 HTTP 端口
EXPOSE 8891 10000

HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:10000/live || exit 1

ENTRYPOINT ["/app/identity_srv"]
