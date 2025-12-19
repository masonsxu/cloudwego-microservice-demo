#!/bin/bash

# =============================================================================
# CloudWeGo Scaffold Services - 部署脚本（精简版）
# =============================================================================
# 用途：简化 Docker Compose 部署操作
# 支持：开发环境部署
#
# 快速开始：
#   ./deploy.sh up              # 启动所有服务
#   ./deploy.sh up-base         # 启动基础设施
#   ./deploy.sh up-app          # 启动应用服务
#   ./deploy.sh down            # 停止服务
#   ./deploy.sh logs            # 查看日志
#   ./deploy.sh ps              # 查看状态
#   ./deploy.sh clean           # 清理环境
#   ./deploy.sh help            # 查看帮助
# =============================================================================

set -euo pipefail

# 脚本配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." &> /dev/null && pwd)"

# 项目配置
PROJECT_NAME="${PROJECT_NAME:-backend}"

# Compose 文件配置（固定使用开发环境）
COMPOSE_BASE="$SCRIPT_DIR/docker-compose.base.yml"
COMPOSE_DEV="$SCRIPT_DIR/docker-compose.dev.yml"
COMPOSE_FILES="-f $COMPOSE_BASE -f $COMPOSE_DEV"

# 服务分组定义
BASE_SERVICES="postgres etcd rustfs redis"
APP_SERVICES="identity_srv gateway"

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
${GREEN}CloudWeGo Scaffold Services - 部署脚本${NC}

${BLUE}用法：${NC}
    $(basename "$0") [操作] [选项]

${BLUE}操作：${NC}
    up              启动所有服务（基础设施 + 应用）
    up-base         仅启动基础设施（postgres、etcd、rustfs、redis）
    up-app          仅启动应用服务（identity_srv、gateway）
    down            停止所有服务
    logs [service]  查看服务日志（可选指定服务名）
    ps              查看服务状态
    clean           清理所有容器和卷
    help            显示此帮助信息

${BLUE}选项：${NC}
    -d, --detach    后台运行（适用于 up、up-base、up-app）
    --no-build      启动时不构建镜像
    --build-only    仅构建镜像，不启动服务
    -f, --force     强制执行（用于 clean）

${BLUE}示例：${NC}
    # 首次部署：启动所有服务
    $(basename "$0") up

    # 后台启动所有服务
    $(basename "$0") up -d

    # 只启动基础设施
    $(basename "$0") up-base

    # 只启动应用服务
    $(basename "$0") up-app -d

    # 仅构建镜像
    $(basename "$0") up --build-only

    # 查看所有日志
    $(basename "$0") logs

    # 查看特定服务日志
    $(basename "$0") logs gateway

    # 查看服务状态
    $(basename "$0") ps

    # 停止服务
    $(basename "$0") down

    # 清理环境（删除容器和卷）
    $(basename "$0") clean -f

${BLUE}服务列表：${NC}
    基础设施:
    - postgres          PostgreSQL 数据库
    - etcd              服务注册与发现
    - rustfs            S3 兼容对象存储
    - redis             Redis 缓存数据库

    应用服务:
    - identity_srv      身份认证服务
    - gateway           API 网关

${YELLOW}注意事项：${NC}
    1. 首次运行前请确保 .env 文件已正确配置
    2. 使用 Docker 原生命令获取更多功能：
       - docker logs -f <container>  # 实时跟踪日志
       - docker exec -it <container> sh  # 进入容器
    3. 清理操作将删除所有数据，请谨慎使用
EOF
}

# 检查依赖
check_dependencies() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    # 优先使用 docker compose 而不是 docker-compose
    if docker compose version &> /dev/null 2>&1; then
        COMPOSE_CMD="docker compose -p ${PROJECT_NAME}"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose -p ${PROJECT_NAME}"
        log_warn "检测到旧版 docker-compose，建议升级到 Docker Compose v2"
    else
        log_error "Docker Compose 未安装"
        exit 1
    fi
}

# 检查 .env 文件
check_env_file() {
    local env_example="$SCRIPT_DIR/.env.dev.example"

    if [ ! -f "$SCRIPT_DIR/.env" ]; then
        log_warn ".env 文件不存在"
        log_info "正在从 $(basename "$env_example") 创建 .env 文件..."
        cp "$env_example" "$SCRIPT_DIR/.env"
        log_info ".env 文件已创建，请根据需要修改配置"
        echo
    fi
}

# 启动服务
start_services() {
    local detach=""
    local build="--build"
    local build_only=false
    local services=""

    for arg in "$@"; do
        case $arg in
            -d|--detach)
                detach="-d"
                ;;
            --no-build)
                build=""
                ;;
            --build-only)
                build_only=true
                ;;
        esac
    done

    cd "$SCRIPT_DIR"

    # 如果仅构建镜像，不启动服务
    if [ "$build_only" = true ]; then
        log_info "仅构建镜像..."
        $COMPOSE_CMD $COMPOSE_FILES build $services
        log_info "镜像构建完成"
        return
    fi

    log_info "启动服务..."

    if [ -n "$build" ]; then
        log_info "检查并构建镜像..."
    fi

    $COMPOSE_CMD $COMPOSE_FILES up $detach $build $services

    if [ -n "$detach" ]; then
        echo
        log_info "服务已启动"
        echo
        show_service_info
    fi
}

# 启动基础服务
start_base_services() {
    local detach=""
    local build="--build"
    local build_only=false

    for arg in "$@"; do
        case $arg in
            -d|--detach)
                detach="-d"
                ;;
            --no-build)
                build=""
                ;;
            --build-only)
                build_only=true
                ;;
        esac
    done

    cd "$SCRIPT_DIR"

    # 如果仅构建镜像，不启动服务
    if [ "$build_only" = true ]; then
        log_info "仅构建基础服务镜像..."
        $COMPOSE_CMD $COMPOSE_FILES build $BASE_SERVICES
        log_info "基础服务镜像构建完成"
        return
    fi

    log_info "启动基础服务..."

    if [ -n "$build" ]; then
        log_info "检查基础服务镜像..."
    fi

    $COMPOSE_CMD $COMPOSE_FILES up $detach $build $BASE_SERVICES

    if [ -n "$detach" ]; then
        echo
        log_info "基础服务已启动"
        $COMPOSE_CMD $COMPOSE_FILES ps
    fi
}

# 启动应用服务
start_app_services() {
    local detach=""
    local build="--build"
    local build_only=false

    for arg in "$@"; do
        case $arg in
            -d|--detach)
                detach="-d"
                ;;
            --no-build)
                build=""
                ;;
            --build-only)
                build_only=true
                ;;
        esac
    done

    cd "$SCRIPT_DIR"

    # 如果仅构建镜像，不启动服务
    if [ "$build_only" = true ]; then
        log_info "仅构建应用服务镜像..."
        $COMPOSE_CMD $COMPOSE_FILES build $APP_SERVICES
        log_info "应用服务镜像构建完成"
        return
    fi

    log_info "启动应用服务..."

    # 检查基础服务是否运行
    local base_running=true
    for service in $BASE_SERVICES; do
        if ! docker ps --filter "name=${PROJECT_NAME}-${service}" --filter "status=running" --format "{{.Names}}" | grep -q "${PROJECT_NAME}-${service}"; then
            log_warn "基础服务 $service 未运行"
            base_running=false
        fi
    done

    if [ "$base_running" = false ]; then
        log_warn "部分基础服务未运行，建议先执行: $(basename "$0") up-base"
        log_info "继续启动应用服务..."
    fi

    if [ -n "$build" ]; then
        log_info "检查并构建应用镜像..."
    fi

    $COMPOSE_CMD $COMPOSE_FILES up $detach $build $APP_SERVICES

    if [ -n "$detach" ]; then
        echo
        log_info "应用服务已启动"
        $COMPOSE_CMD $COMPOSE_FILES ps
    fi
}

# 停止服务
stop_services() {
    cd "$SCRIPT_DIR"

    # 检查是否有运行的容器
    local running_containers
    running_containers=$($COMPOSE_CMD $COMPOSE_FILES ps -q 2>/dev/null || true)

    if [ -z "$running_containers" ]; then
        log_info "没有运行的容器，无需停止"
        return 0
    fi

    log_info "停止服务..."
    # 添加错误容忍，避免在资源不存在时报错
    $COMPOSE_CMD $COMPOSE_FILES down 2>&1 | grep -v "no container\|no pod\|not found" || true

    log_info "服务已停止"
}

# 查看日志
show_logs() {
    local service="${1:-}"
    cd "$SCRIPT_DIR"

    if [ -n "$service" ]; then
        log_info "查看 $service 服务日志（最近 100 条）..."
        $COMPOSE_CMD $COMPOSE_FILES logs --tail=100 "$service"
    else
        log_info "查看所有服务日志（最近 50 条）..."
        $COMPOSE_CMD $COMPOSE_FILES logs --tail=50
    fi
}

# 查看状态
show_status() {
    log_info "CloudWeGo Scaffold 服务状态:"
    cd "$SCRIPT_DIR"
    $COMPOSE_CMD $COMPOSE_FILES ps
}

# 清理环境
clean_environment() {
    local force="${1:-}"

    if [ "$force" != "-f" ] && [ "$force" != "--force" ]; then
        echo -e "${YELLOW}警告：此操作将删除所有容器、卷和数据！${NC}"
        read -p "确定要继续吗？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "操作已取消"
            exit 0
        fi
    fi

    log_warn "清理环境..."
    cd "$SCRIPT_DIR"

    # 停止并移除容器和卷
    # 添加错误容忍，避免在资源不存在时报错
    $COMPOSE_CMD $COMPOSE_FILES down -v --remove-orphans 2>&1 | grep -v "no container\|no pod\|not found" || true

    log_info "环境清理完成"
}

# 显示服务访问信息
show_service_info() {
    log_info "=== 服务访问地址 ==="
    echo -e "${GREEN}API Gateway:${NC}     http://localhost:8080"
    echo -e "${GREEN}Identity Service:${NC} http://localhost:8891 (RPC), http://localhost:10000 (Health)"
    echo
    echo -e "${GREEN}基础设施:${NC}"
    echo -e "${GREEN}PostgreSQL:${NC}      localhost:5432"
    echo -e "${GREEN}etcd:${NC}            localhost:2379"
    echo -e "${GREEN}Redis:${NC}           localhost:6379"
    echo -e "${GREEN}RustFS API:${NC}      http://localhost:9000"
    echo -e "${GREEN}RustFS Console:${NC}  http://localhost:9001"
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

    # 解析操作
    case "${1:-help}" in
        up)
            shift
            start_services "$@"
            ;;
        up-base)
            shift
            start_base_services "$@"
            ;;
        up-app)
            shift
            start_app_services "$@"
            ;;
        down)
            stop_services
            ;;
        logs)
            shift
            show_logs "$@"
            ;;
        ps)
            show_status
            ;;
        clean)
            shift
            clean_environment "$@"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知操作: $1"
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
