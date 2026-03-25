#!/bin/bash
# 脚本功能：Kitex代码生成工具
# 根据IDL定义生成Kitex项目代码

# 如果任何命令失败，立即退出
set -e

# --- 配置 ---
# 从 go.mod 文件自动获取模块名
MODULE_NAME=$(grep -m 1 "module" ./go.mod | awk '{print $2}')
KITEX_VERSION="v0.16.1"
KITEX_SERVICE_NAME="identity-service"
PROTO_FILE_PATH="../../idl/rpc/identity_srv/identity_service.proto"
IDL_INCLUDE_PATH="../../idl"
KITEX_CMD=""

# 颜色输出定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的信息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查kitex命令是否存在
check_kitex() {
    if command -v kitex >/dev/null 2>&1; then
        KITEX_CMD="kitex"
        return
    fi

    print_warn "未找到kitex命令，回退使用 go run github.com/cloudwego/kitex/tool/cmd/kitex@$KITEX_VERSION"
    KITEX_CMD="go run github.com/cloudwego/kitex/tool/cmd/kitex@$KITEX_VERSION"
}

# 检查 protobuf 代码生成工具是否存在
check_protobuf_tools() {
    if ! command -v protoc >/dev/null 2>&1; then
        print_error "未找到protoc命令，请先安装Protocol Buffers编译器"
        exit 1
    fi

    if ! command -v protoc-gen-go >/dev/null 2>&1; then
        print_error "未找到protoc-gen-go命令，请先安装protobuf Go插件"
        print_info "安装方法: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
        exit 1
    fi
}

# --- 执行 ---
print_info "当前工作目录: $(pwd)"
print_info "模块名: $MODULE_NAME"
print_info "开始生成 Kitex 代码..."

# 检查依赖工具
check_kitex
check_protobuf_tools

# 检查IDL文件是否存在
if [ ! -f "$PROTO_FILE_PATH" ]; then
    print_error "IDL文件不存在: $PROTO_FILE_PATH"
    exit 1
fi

# 使用 -v 参数输出更详细的日志，便于调试
$KITEX_CMD -v \
      -type protobuf \
      -module "$MODULE_NAME" \
      -service "$KITEX_SERVICE_NAME" \
      -I "$IDL_INCLUDE_PATH" \
      "$PROTO_FILE_PATH"

print_info "代码生成完毕。"