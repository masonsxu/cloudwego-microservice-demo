# 多阶段构建的 Dockerfile for gateway
# 基于 Go 1.24 构建，遵循容器最佳实践
# 必须在项目根目录下构建: podman build -f gateway/docker/Dockerfile -t gateway .

# 构建阶段
FROM m.daocloud.io/docker.io/library/golang:1.24 AS builder

# 设置 Go 环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 go.work 和 go.work.sum (如果有)
# COPY go.work go.work.sum ./

# 复制依赖的 rpc 服务 go.mod 和 go.sum
# 注意：这里假设构建上下文是项目根目录
COPY rpc/identity_srv/go.mod rpc/identity_srv/go.sum ./rpc/identity_srv/

# 复制 api/api_gateway 的 go.mod 和 go.sum
COPY api/api_gateway/go.mod api/api_gateway/go.sum ./api/api_gateway/

# 下载依赖
# 切换到 api 目录下载依赖，这样 go.mod 中的相对路径 replace 才能正确解析
WORKDIR /app/api/api_gateway
RUN go mod download

# 安装 swag 工具
RUN go install github.com/swaggo/swag/cmd/swag@latest

# 复制源代码
# 复制 rpc 服务源码 (如果 api 网关依赖了 rpc 的代码而不仅仅是接口定义，通常依赖 idl 生成的代码)
# 如果依赖的是 kitex 生成的代码，通常在 rpc/xxx_srv/kitex_gen 下，或者 api 网关自己生成了一份
# 根据 go.mod replace 指令，api 网关直接依赖 rpc 服务的模块，所以需要复制 rpc 服务源码
WORKDIR /app
COPY rpc/identity_srv/ ./rpc/identity_srv/
# 复制 api 网关源码
COPY api/api_gateway/ ./api/api_gateway/

# 构建应用
WORKDIR /app/api/api_gateway
RUN chmod +x ./build.sh && ./build.sh

# 运行阶段
FROM m.daocloud.io/docker.io/library/alpine:latest

# 安装 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata curl

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D -s /bin/sh gateway

# 设置工作目录
WORKDIR /app

# 复制构建产物
COPY --from=builder /app/api/api_gateway/output/ ./output/

# 更改文件所有者
RUN chown -R gateway:gateway /app

# 切换到非 root 用户
USER gateway

# 暴露服务端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ping || exit 1

# 启动命令
ENTRYPOINT ["/app/output/bootstrap.sh"]
