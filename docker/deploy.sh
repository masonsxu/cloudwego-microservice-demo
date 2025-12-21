#!/bin/bash

# =============================================================================
# CloudWeGo 微服务基础设施 - 部署脚本
# =============================================================================
# 用途：快速启动/管理本地开发环境的基础设施服务
# 服务：PostgreSQL、etcd、Redis、RustFS、Jaeger
#
# 快速开始：
#   ./deploy.sh up      # 启动所有服务
#   ./deploy.sh down    # 停止所有服务
#   ./deploy.sh ps      # 查看服务状态
#   ./deploy.sh logs    # 查看日志
#   ./deploy.sh help    # 查看帮助
# =============================================================================

set -euo pipefail

# 脚本配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
PROJECT_NAME="${PROJECT_NAME:-backend}"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
${GREEN}CloudWeGo 微服务基础设施 - 部署脚本${NC}

${BLUE}用法：${NC}
    $(basename "$0") [命令]

${BLUE}命令：${NC}
    up              启动所有基础设施服务（后台运行）
    down            停止所有服务
    restart         重启所有服务
    logs [service]  查看服务日志（可选指定服务名）
    ps              查看服务状态
    clean           清理所有容器和数据卷（谨慎使用）
    help            显示此帮助信息

${BLUE}服务列表：${NC}
    - postgres      PostgreSQL 数据库
    - etcd          服务注册与发现
    - rustfs        S3 兼容对象存储
    - redis         Redis 缓存数据库
    - jaeger        链路追踪服务

${BLUE}示例：${NC}
    # 启动所有基础设施
    $(basename "$0") up

    # 查看所有日志
    $(basename "$0") logs

    # 查看特定服务日志
    $(basename "$0") logs postgres

    # 查看服务状态
    $(basename "$0") ps

    # 停止服务
    $(basename "$0") down

    # 重启服务
    $(basename "$0") restart

    # 清理环境（删除容器和数据）
    $(basename "$0") clean

${BLUE}服务访问地址：${NC}
    PostgreSQL:      localhost:5432
    etcd:            localhost:2379
    Redis:           localhost:6379
    RustFS API:      http://localhost:9000
    RustFS Console:  http://localhost:9001
    Jaeger UI:       http://localhost:16686

${YELLOW}注意事项：${NC}
    1. 首次运行前请确保 .env 文件已正确配置
    2. Go 应用服务请在本地启动调试
    3. 使用 Podman 原生命令获取更多功能：
       - podman logs -f <container>  # 实时跟踪日志
       - podman exec -it <container> sh  # 进入容器
    4. clean 命令将删除所有数据，请谨慎使用
EOF
}

# 检查依赖
check_dependencies() {
    if ! command -v podman &> /dev/null; then
        log_error "Podman 未安装，请先安装 Podman"
        log_info "安装方法：https://podman.io/getting-started/installation"
        exit 1
    fi

    if ! command -v podman-compose &> /dev/null; then
        log_error "podman-compose 未安装"
        log_info "安装方法: pip3 install podman-compose"
        exit 1
    fi

    COMPOSE_CMD="podman-compose -p ${PROJECT_NAME}"
}

# 检查 .env 文件
check_env_file() {
    local env_example="$SCRIPT_DIR/.env.dev.example"

    if [ ! -f "$SCRIPT_DIR/.env" ]; then
        log_warn ".env 文件不存在"
        if [ -f "$env_example" ]; then
            log_info "正在从 $(basename "$env_example") 创建 .env 文件..."
            cp "$env_example" "$SCRIPT_DIR/.env"
            log_info ".env 文件已创建，请根据需要修改配置"
            echo
        else
            log_error ".env.dev.example 文件不存在，无法自动创建配置"
            exit 1
        fi
    fi
}

# 启动服务
start_services() {
    cd "$SCRIPT_DIR"
    log_info "启动基础设施服务..."

    $COMPOSE_CMD up -d

    echo
    log_info "服务已启动"
    echo
    show_service_info
}

# 停止服务
stop_services() {
    cd "$SCRIPT_DIR"

    # 检查是否有运行的容器
    local running_containers
    running_containers=$($COMPOSE_CMD ps -q 2>/dev/null || true)

    if [ -z "$running_containers" ]; then
        log_info "没有运行的容器，无需停止"
        return 0
    fi

    log_info "停止服务..."
    $COMPOSE_CMD down 2>&1 | grep -v "no container\|no pod\|not found" || true

    log_info "服务已停止"
}

# 重启服务
restart_services() {
    cd "$SCRIPT_DIR"
    log_info "重启服务..."

    $COMPOSE_CMD restart

    log_info "服务已重启"
}

# 查看日志
show_logs() {
    local service="${1:-}"
    cd "$SCRIPT_DIR"

    if [ -n "$service" ]; then
        log_info "查看 $service 服务日志（最近 100 条）..."
        $COMPOSE_CMD logs --tail=100 -f "$service"
    else
        log_info "查看所有服务日志（最近 50 条）..."
        $COMPOSE_CMD logs --tail=50 -f
    fi
}

# 查看状态
show_status() {
    log_info "基础设施服务状态:"
    cd "$SCRIPT_DIR"
    $COMPOSE_CMD ps
}

# 清理环境
clean_environment() {
    echo -e "${YELLOW}警告：此操作将删除所有容器、卷和数据！${NC}"
    read -p "确定要继续吗？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "操作已取消"
        exit 0
    fi

    log_warn "清理环境..."
    cd "$SCRIPT_DIR"

    # 停止并移除容器和卷
    $COMPOSE_CMD down -v --remove-orphans 2>&1 | grep -v "no container\|no pod\|not found" || true

    log_info "环境清理完成"
}

# 显示服务访问信息
show_service_info() {
    echo -e "${GREEN}=== 服务访问地址 ===${NC}"
    echo -e "${GREEN}PostgreSQL:${NC}      localhost:5432"
    echo -e "${GREEN}etcd:${NC}            localhost:2379"
    echo -e "${GREEN}Redis:${NC}           localhost:6379"
    echo -e "${GREEN}RustFS API:${NC}      http://localhost:9000"
    echo -e "${GREEN}RustFS Console:${NC}  http://localhost:9001"
    echo -e "${GREEN}Jaeger UI:${NC}       http://localhost:16686"
    echo
    log_info "使用 '$(basename "$0") logs' 查看服务日志"
    log_info "使用 '$(basename "$0") ps' 查看服务状态"
}

# 解析参数
parse_args() {
    if [ $# -eq 0 ]; then
        show_help
        exit 0
    fi

    check_dependencies
    check_env_file

    # 解析命令
    case "${1:-help}" in
        up)
            start_services
            ;;
        down)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            shift
            show_logs "$@"
            ;;
        ps)
            show_status
            ;;
        clean)
            clean_environment
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $1"
            echo
            show_help
            exit 1
            ;;
    esac
}

# 主函数
main() {
    cd "$SCRIPT_DIR"
    parse_args "$@"
}

# 捕获信号
trap 'log_error "脚本被中断"; exit 1' INT TERM

# 执行主函数
main "$@"
