#!/usr/bin/env bash
# install.sh - Git hooks 安装/卸载脚本
#
# 用法：
#   ./scripts/git-hooks/install.sh             # 安装
#   ./scripts/git-hooks/install.sh --uninstall # 卸载

set -euo pipefail

# =========================
# 公共工具函数
# =========================

if [[ -t 2 ]] && [[ -z "${NO_COLOR:-}" ]]; then
    _RED='\033[0;31m'
    _YELLOW='\033[0;33m'
    _GREEN='\033[0;32m'
    _RESET='\033[0m'
else
    _RED='' _YELLOW='' _GREEN='' _RESET=''
fi

info()  { printf "${_GREEN}[hooks]${_RESET} %s\n" "$*" >&2; }
warn()  { printf "${_YELLOW}[hooks 警告]${_RESET} %s\n" "$*" >&2; }
error() { printf "${_RED}[hooks 错误]${_RESET} %s\n" "$*" >&2; }

# =========================
# 定位项目根目录
# =========================
repo_root=$(git rev-parse --show-toplevel 2>/dev/null)
if [[ -z "$repo_root" ]]; then
    error "当前不在 Git 仓库中，请在项目目录下运行此脚本。"
    exit 1
fi

hooks_source_dir="${repo_root}/scripts/git-hooks"
hooks_target_dir="${repo_root}/.git/hooks"

# 需要安装的 hook 列表
hooks=(
    "pre-commit"
)

# =========================
# 卸载模式
# =========================
if [[ "${1:-}" == "--uninstall" ]]; then
    info "开始卸载 Git hooks..."

    for hook in "${hooks[@]}"; do
        target="${hooks_target_dir}/${hook}"
        if [[ -L "$target" ]]; then
            rm -f "$target"
            info "已移除：${hook}"
        elif [[ -f "$target" ]]; then
            # 不是符号链接，可能是用户自定义的 hook，不删除
            warn "${hook} 不是符号链接（可能是自定义 hook），已跳过。"
        else
            info "${hook} 不存在，跳过。"
        fi
    done

    info "Git hooks 卸载完成。"
    exit 0
fi

# =========================
# 安装模式
# =========================
info "开始安装 Git hooks..."

# 确保目标目录存在
mkdir -p "$hooks_target_dir"

installed=0
skipped=0

for hook in "${hooks[@]}"; do
    source="${hooks_source_dir}/${hook}"
    target="${hooks_target_dir}/${hook}"

    # 检查源文件是否存在
    if [[ ! -f "$source" ]]; then
        warn "源文件不存在：${source}，跳过 ${hook}。"
        skipped=$((skipped + 1))
        continue
    fi

    # 设置源文件可执行权限
    chmod +x "$source"

    # 如果目标已存在且不是指向我们 source 的符号链接，先备份
    if [[ -f "$target" ]] && [[ ! -L "$target" ]]; then
        backup="${target}.backup.$(date +%Y%m%d%H%M%S)"
        mv "$target" "$backup"
        warn "${hook} 已有自定义内容，已备份到 ${backup##*/}"
    fi

    # 创建符号链接（-f 覆盖已有链接）
    ln -sf "$source" "$target"
    installed=$((installed + 1))
    info "已安装：${hook} → scripts/git-hooks/${hook}"
done

info "安装完成：${installed} 个 hook 已安装，${skipped} 个跳过。"

# 验证安装
info "验证安装结果："
for hook in "${hooks[@]}"; do
    target="${hooks_target_dir}/${hook}"
    if [[ -L "$target" ]] && [[ -x "$target" ]]; then
        printf "  ✓ %s\n" "$hook" >&2
    else
        printf "  ✗ %s（异常）\n" "$hook" >&2
    fi
done

exit 0
