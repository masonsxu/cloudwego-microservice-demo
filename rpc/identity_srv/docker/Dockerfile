# 多阶段构建的 Dockerfile for identity-service
# 基于 Go 1.24 构建，遵循容器最佳实践

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

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o identity-service .

# 运行阶段
FROM m.daocloud.io/docker.io/library/alpine:latest

# 安装 ca-certificates 和 tzdata 以支持 HTTPS 请求和时区
RUN apk --no-cache add ca-certificates tzdata curl

# 设置时区环境变量
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件和配置文件
COPY --from=builder /app/identity-service .
COPY --from=builder /app/config/permission_model.conf ./config/
COPY --from=builder /app/script/ ./script/

# 更改文件所有者
RUN chown -R appuser:appuser /app

# 切换到非 root 用户
USER appuser
# 暴露服务端口
EXPOSE 8888
EXPOSE 10000

# 健康检查
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:10000/live || exit 1

# 启动命令
ENTRYPOINT ["./identity-service"]
