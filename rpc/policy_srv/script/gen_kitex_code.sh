#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE_DIR="$(dirname "$SCRIPT_DIR")"
cd "$SERVICE_DIR"

MODULE_NAME=$(grep "^module " go.mod | awk '{print $2}')
KITEX_VERSION="v0.16.1"
KITEX_SERVICE_NAME="policy-service"
PROTO_FILE_PATH="../../idl/rpc/policy_srv/policy_service.proto"
IDL_INCLUDE_PATH="../../idl"

# 检查 kitex 命令
if ! command -v kitex &>/dev/null; then
    echo "kitex 未安装，使用 go run github.com/cloudwego/kitex/tool/cmd/kitex@${KITEX_VERSION}"
    KITEX_CMD="go run github.com/cloudwego/kitex/tool/cmd/kitex@${KITEX_VERSION}"
else
    KITEX_CMD="kitex"
fi

# 检查 protoc
if ! command -v protoc &>/dev/null; then
    echo "ERROR: protoc 未安装"
    exit 1
fi

echo "==> 开始生成 policy_srv Kitex 代码"
$KITEX_CMD -v -type protobuf \
    -module "$MODULE_NAME" \
    -service "$KITEX_SERVICE_NAME" \
    -I "$IDL_INCLUDE_PATH" \
    "$PROTO_FILE_PATH"

echo "==> 代码生成完成"
