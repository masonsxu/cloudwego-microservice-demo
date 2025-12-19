#!/bin/bash

# 脚本功能：Hertz代码生成工具
# 根据IDL定义生成Hertz项目代码
# 使用方法:
#   ./gen_hertz_code.sh                 # 生成所有服务代码（详细模式）
#   ./gen_hertz_code.sh identity        # 仅生成identity服务代码（详细模式）
#   ./gen_hertz_code.sh permission      # 仅生成permission服务代码（详细模式）
#   ./gen_hertz_code.sh identity --silent  # 生成identity服务代码（静默模式）

set -e

# IDL根目录 (参考gen_kitex_code.sh的相对路径方式)
IDL_ROOT="../idl"
HTTP_IDL_ROOT="$IDL_ROOT/http"

# 颜色输出定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认使用verbose模式
VERBOSE_FLAG="--verbose"

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

# 检查hz命令是否存在
check_hz() {
    if ! command -v hz >/dev/null 2>&1; then
        print_error "未找到hz命令，请先安装hertz工具"
        print_info "安装方法: go install github.com/cloudwego/hertz/cmd/hz@latest"
        exit 1
    fi
}

# 解析命令行参数
parse_args() {
    if [[ $# -gt 0 ]]; then
        SERVICE_NAME="$1"
    fi
}

# 生成单个服务的代码
generate_service() {
    local service_name=$1
    local idl_file="$HTTP_IDL_ROOT/$service_name/${service_name}_service.thrift"

    # 检查IDL文件是否存在
    if [ ! -f "$idl_file" ]; then
        print_error "IDL文件不存在: $idl_file"
        print_info "当前工作目录: $(pwd)"
        return 1
    fi

    if [ "$SILENT_MODE" = false ]; then
        print_info "正在生成 $service_name 服务代码..."
    fi

    # 使用verbose模式执行hz命令
    hz $VERBOSE_FLAG update -I "$IDL_ROOT" -idl "$idl_file"

    if [ $? -eq 0 ]; then
        print_info "$service_name 服务代码生成完成"
    else
        print_error "$service_name 服务代码生成失败"
        return 1
    fi
}


# 生成所有服务的代码
generate_all() {
    print_info "开始生成所有服务代码..."

    # 检查所有IDL文件是否存在
    local services=("identity" "permission")

    for service in "${services[@]}"; do
        local idl_file="$HTTP_IDL_ROOT/$service/${service}_service.thrift"
        if [ ! -f "$idl_file" ]; then
            print_error "IDL文件不存在: $idl_file"
            print_info "当前工作目录: $(pwd)"
            exit 1
        fi
    done

    # 生成所有服务代码
    for service in "${services[@]}"; do
        generate_service "$service"
    done

    print_info "所有服务代码生成完成"
}


# 主逻辑
main() {
    check_hz

    # 解析命令行参数
    parse_args "$@"

    # 根据参数决定执行方式
    if [ -n "$SERVICE_NAME" ]; then
        case "$SERVICE_NAME" in
            identity|permission)
                generate_service "$SERVICE_NAME"
                ;;
            *)
                print_error "未知服务: $SERVICE_NAME"
                echo "使用方法: $0 [identity|permission]"
                exit 1
                ;;
        esac

    else
        generate_all
    fi
}


# 执行主逻辑
main "$@"
