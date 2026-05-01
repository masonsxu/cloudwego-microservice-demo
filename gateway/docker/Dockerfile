# =============================================================================
# gateway HTTP 网关镜像
# =============================================================================
# 构建命令（必须在项目根目录执行，因为 gateway 通过 replace 引用 ../rpc/identity_srv）：
#   podman build -f gateway/docker/Dockerfile -t gateway:latest .
#
# 说明：
#   - 构建上下文必须是项目根，才能同时拷贝 gateway/ 和 rpc/identity_srv/
#   - GOWORK=off 禁用 workspace，由 go.mod 中的 replace 指令解析模块路径
#   - swag init 在容器内执行，宿主机不需要预先生成 docs/
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

# 利用层缓存：先拷 go.mod/go.sum 下载依赖（保留相对路径，让 replace 能解析）
COPY rpc/identity_srv/go.mod rpc/identity_srv/go.sum ./rpc/identity_srv/
COPY gateway/go.mod gateway/go.sum ./gateway/
RUN cd gateway && go mod download

# 安装 swag（构建期工具，不进入运行镜像）
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制全部源码（含 replace 指向的 rpc/identity_srv）
COPY rpc/identity_srv/ ./rpc/identity_srv/
COPY gateway/ ./gateway/

WORKDIR /build/gateway

# 生成 Swagger 文档
RUN swag init --parseDependencyLevel 1 --useStructName true

# 构建二进制
RUN go build -ldflags="-s -w" -trimpath -o hertz_service .

# -----------------------------------------------------------------------------
# 运行阶段
# -----------------------------------------------------------------------------
FROM m.daocloud.io/docker.io/library/alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl
ENV TZ=Asia/Shanghai

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# 拷贝二进制 + 运行时所需配置
COPY --from=builder /build/gateway/hertz_service .
COPY --from=builder /build/gateway/config/ ./config/
# CASBIN_MENU_MAPPING_FILE 默认相对路径 menu.yaml（位于项目根，直接从构建上下文拷贝）
COPY menu.yaml ./menu.yaml

RUN mkdir -p /app/logs && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ping || exit 1

ENTRYPOINT ["/app/hertz_service"]
