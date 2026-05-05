# =============================================================================
# policy_srv PDP 服务镜像
# =============================================================================
FROM m.daocloud.io/docker.io/library/golang:1.25 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -trimpath -o policy_srv .

FROM m.daocloud.io/docker.io/library/alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl
ENV TZ=Asia/Shanghai

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /build/policy_srv .

RUN mkdir -p /app/logs && chown -R appuser:appuser /app

USER appuser

EXPOSE 8892 10001

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:10001/live || exit 1

ENTRYPOINT ["/app/policy_srv"]
