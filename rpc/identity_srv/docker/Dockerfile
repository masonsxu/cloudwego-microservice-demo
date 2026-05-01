# =============================================================================
# identity_srv RPC 服务镜像
# =============================================================================
# 构建命令（必须在项目根目录执行）：
#   podman build -f rpc/identity_srv/docker/Dockerfile -t identity-srv:latest rpc/identity_srv
#
# 说明：
#   - identity_srv 自身 go.mod 不依赖外部 replace 路径，可以用子目录作为构建上下文
#   - 通过 GOWORK=off 显式禁用 go.work，避免容器内找不到兄弟模块
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

# 利用 Docker 层缓存：先单独 COPY go.mod/go.sum 下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并构建
COPY . .
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
COPY --from=builder /build/identity_srv .
COPY --from=builder /build/config/seed_menu.yaml ./config/seed_menu.yaml

# 准备日志目录并修正归属
RUN mkdir -p /app/logs && chown -R appuser:appuser /app

USER appuser

# 8891: Kitex RPC 端口；10000: 健康检查 HTTP 端口
EXPOSE 8891 10000

HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:10000/live || exit 1

ENTRYPOINT ["/app/identity_srv"]
