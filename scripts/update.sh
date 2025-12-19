#!/bin/bash
# 更新脚本：从脚手架仓库拉取更新
# 使用方法: ./scripts/update.sh [branch]
# 示例: ./scripts/update.sh main

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

# 检查是否在 Git 仓库中
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}错误：当前目录不是 Git 仓库${NC}"
    exit 1
fi

# 检查是否有未提交的更改
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}警告：检测到未提交的更改${NC}"
    echo "  建议先提交或暂存当前更改，然后再更新"
    read -p "是否继续? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# 检查 upstream 远程是否存在
if ! git remote | grep -q "^upstream$"; then
    echo -e "${RED}错误：未找到 upstream 远程仓库${NC}"
    echo ""
    echo "请先设置 upstream 远程："
    echo "  git remote add upstream <scaffold-repo-url>"
    echo ""
    echo "或者在初始化时提供脚手架仓库 URL："
    echo "  ./scripts/init.sh <your-module-path> <scaffold-repo-url>"
    exit 1
fi

UPSTREAM_URL=$(git remote get-url upstream)
BRANCH="${1:-main}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  从脚手架仓库拉取更新${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "Upstream: ${GREEN}$UPSTREAM_URL${NC}"
echo -e "分支: ${GREEN}$BRANCH${NC}"
echo ""

# 获取当前分支
CURRENT_BRANCH=$(git branch --show-current)
echo -e "${BLUE}当前分支: ${CURRENT_BRANCH}${NC}"
echo ""

# 1. 获取 upstream 的最新更改
echo -e "${GREEN}[1/4] 获取 upstream 最新更改...${NC}"
git fetch upstream "$BRANCH"

# 2. 检查是否有更新
LOCAL=$(git rev-parse HEAD)
REMOTE=$(git rev-parse upstream/$BRANCH)

if [ "$LOCAL" = "$REMOTE" ]; then
    echo -e "${GREEN}✓ 已经是最新版本，无需更新${NC}"
    exit 0
fi

echo -e "${YELLOW}发现新版本：${NC}"
echo "  本地: $LOCAL"
echo "  远程: $REMOTE"
echo ""

# 3. 创建更新分支（可选，用于安全合并）
read -p "是否创建临时分支进行合并? (Y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Nn]$ ]]; then
    UPDATE_BRANCH="update-from-scaffold-$(date +%Y%m%d-%H%M%S)"
    echo -e "${BLUE}创建临时分支: $UPDATE_BRANCH${NC}"
    git checkout -b "$UPDATE_BRANCH"
    echo ""
fi

# 4. 合并 upstream 更改
echo -e "${GREEN}[2/4] 合并 upstream 更改...${NC}"
echo -e "${YELLOW}注意：如果出现冲突，脚本会暂停让你手动解决${NC}"
echo ""

# 尝试合并
if git merge upstream/$BRANCH --no-edit; then
    echo -e "${GREEN}✓ 合并成功，无冲突${NC}"
else
    # 检查是否有冲突
    if git diff --check --diff-filter=U 2>/dev/null || [ -n "$(git ls-files -u)" ]; then
        echo ""
        echo -e "${RED}========================================${NC}"
        echo -e "${RED}  检测到合并冲突${NC}"
        echo -e "${RED}========================================${NC}"
        echo ""
        echo -e "${YELLOW}冲突文件列表：${NC}"
        git diff --name-only --diff-filter=U 2>/dev/null || git ls-files -u | awk '{print $4}' | sort -u
        echo ""
        echo -e "${YELLOW}解决冲突指南：${NC}"
        echo "  1. 打开冲突文件，查找 <<<<<<< HEAD 标记"
        echo "  2. 保留你的更改（通常是你的 module 路径）"
        echo "  3. 删除冲突标记（<<<<<<< HEAD, =======, >>>>>>> upstream/main）"
        echo "  4. 保存文件后运行: git add <file>"
        echo "  5. 完成所有冲突解决后运行: git commit"
        echo ""
        echo -e "${BLUE}常见冲突类型：${NC}"
        echo "  - go.mod 文件：保留你的 module 路径，但接受依赖版本更新"
        echo "  - 配置文件：通常保留你的配置，但接受新的配置项"
        echo "  - 业务代码：根据实际情况决定保留哪个版本"
        echo ""
        echo -e "${YELLOW}提示：${NC}"
        echo "  如果冲突太多，可以中止合并："
        echo "    git merge --abort"
        echo ""
        echo "  然后手动选择需要的更改："
        echo "    git checkout upstream/$BRANCH -- <file>"
        echo ""
        exit 1
    else
        echo -e "${GREEN}✓ 合并完成${NC}"
    fi
fi

# 5. 更新依赖
echo -e "${GREEN}[3/4] 更新 Go 依赖...${NC}"
if command -v go >/dev/null 2>&1; then
    go mod tidy
    echo -e "${GREEN}✓ 依赖已更新${NC}"
else
    echo -e "${YELLOW}警告：未找到 go 命令，跳过依赖更新${NC}"
fi

# 6. 提示重新生成代码
echo -e "${GREEN}[4/4] 检查是否需要重新生成代码...${NC}"
echo ""
echo -e "${YELLOW}重要：如果脚手架更新了 IDL 或代码生成配置，需要重新生成代码${NC}"
echo ""
echo "建议执行以下操作："
echo "  1. 重新生成 RPC 代码（如果 IDL 有更新）："
echo "     cd rpc/identity_srv && ./script/gen_kitex_code.sh"
echo ""
echo "  2. 重新生成 HTTP 代码（如果 IDL 有更新）："
echo "     cd gateway && ./script/gen_hertz_code.sh"
echo ""
echo "  3. 重新生成 Wire 依赖注入代码："
echo "     cd rpc/identity_srv/wire && wire"
echo "     cd ../../gateway/internal/wire && wire"
echo ""

# 如果创建了临时分支，提示合并到主分支
if [ -n "$UPDATE_BRANCH" ] && [ "$UPDATE_BRANCH" != "$CURRENT_BRANCH" ]; then
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  更新完成${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
    echo -e "当前在临时分支: ${GREEN}$UPDATE_BRANCH${NC}"
    echo ""
    echo "如果一切正常，可以合并到主分支："
    echo "  git checkout $CURRENT_BRANCH"
    echo "  git merge $UPDATE_BRANCH"
    echo "  git branch -d $UPDATE_BRANCH"
    echo ""
    echo "如果有问题，可以删除临时分支："
    echo "  git checkout $CURRENT_BRANCH"
    echo "  git branch -D $UPDATE_BRANCH"
    echo ""
else
    echo -e "${GREEN}✓ 更新完成！${NC}"
fi
