#!/bin/bash
# 测试辅助脚本 - 自动配置 podman 环境

set -e

# 检测 XDG_RUNTIME_DIR
if [ -z "$XDG_RUNTIME_DIR" ]; then
    echo "错误: XDG_RUNTIME_DIR 环境变量未设置"
    exit 1
fi

# 检测 podman socket
PODMAN_SOCKET="$XDG_RUNTIME_DIR/podman/podman.sock"
if [ ! -S "$PODMAN_SOCKET" ]; then
    echo "错误: podman socket 不存在: $PODMAN_SOCKET"
    echo "请确保 podman 正在运行"
    exit 1
fi

# 设置 DOCKER_HOST
export DOCKER_HOST="unix://$PODMAN_SOCKET"

# 进入脚本所在目录
cd "$(dirname "$0")"

# 默认运行 RPC 服务测试
if [ $# -eq 0 ]; then
    echo "🧪 运行 RPC 服务所有测试..."
    go test ./... -v
else
    echo "🧪 运行指定测试: $*"
    go test "$@"
fi
