#!/usr/bin/env bash
# =============================================================================
# 构建本地开发用的容器镜像（identity-srv、gateway）
# =============================================================================
# 用法：./scripts/build-images.sh [identity|gateway|all]
#
# 镜像构建后可被 podman kube play docker/pod.yml 直接引用
# （pod.yml 中已设置 imagePullPolicy: IfNotPresent 优先使用本地镜像）
# =============================================================================
set -eou pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."
readonly ROOT_DIR="$(pwd)"
readonly TARGET="${1:-all}"

build_identity() {
    echo "==> 构建 identity-srv:latest"
    podman build \
        -f "${ROOT_DIR}/rpc/identity_srv/docker/Dockerfile" \
        -t identity-srv:latest \
        "${ROOT_DIR}/rpc/identity_srv"
}

build_gateway() {
    echo "==> 构建 gateway:latest"
    podman build \
        -f "${ROOT_DIR}/gateway/docker/Dockerfile" \
        -t gateway:latest \
        "${ROOT_DIR}"
}

case "${TARGET}" in
    identity) build_identity ;;
    gateway)  build_gateway ;;
    all)      build_identity; build_gateway ;;
    *)
        echo "ERROR: 未知目标 '${TARGET}'，支持：identity | gateway | all" >&2
        exit 1
        ;;
esac

echo
echo "==> 完成，本地镜像列表："
podman images | grep -E "^(localhost/)?(identity-srv|gateway)\s" || true
