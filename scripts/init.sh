#!/bin/bash
# 初始化脚本：自动替换 Go module 路径
# 使用方法: ./scripts/init.sh <your-module-path> [scaffold-repo-url]
# 示例: ./scripts/init.sh github.com/your-org/your-project

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 原始 module 路径（脚手架中的占位符）
OLD_MODULE="github.com/masonsxu/cloudwego-scaffold"
OLD_MODULE_ESCAPED=$(echo "$OLD_MODULE" | sed 's/[\/&]/\\&/g')

# 检查参数
if [ -z "$1" ]; then
    echo -e "${RED}错误：请提供新的 module 路径${NC}"
    echo ""
    echo "使用方法:"
    echo "  $0 <your-module-path> [scaffold-repo-url]"
    echo ""
    echo "示例:"
    echo "  $0 github.com/your-org/your-project"
    echo "  $0 github.com/your-org/your-project https://github.com/masonsxu/cloudwego-scaffold.git"
    exit 1
fi

NEW_MODULE="$1"
NEW_MODULE_ESCAPED=$(echo "$NEW_MODULE" | sed 's/[\/&]/\\&/g')
SCAFFOLD_REPO="${2:-}"

# 验证 module 路径格式
if [[ ! "$NEW_MODULE" =~ ^[a-z0-9][a-z0-9._-]*/[a-z0-9][a-z0-9._-]*(/[a-z0-9][a-z0-9._-]*)*$ ]]; then
    echo -e "${YELLOW}警告：module 路径格式可能不正确${NC}"
    echo "  建议格式: github.com/org/project 或 example.com/project"
    read -p "是否继续? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Go Module 路径初始化工具${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "旧路径: ${YELLOW}$OLD_MODULE${NC}"
echo -e "新路径: ${GREEN}$NEW_MODULE${NC}"
echo ""

# 获取项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

# 统计需要替换的文件
echo -e "${BLUE}正在查找需要替换的文件...${NC}"

# 1. 替换 go.mod 文件
echo -e "${GREEN}[1/4] 替换 go.mod 文件...${NC}"
find . -name "go.mod" -type f | while read -r file; do
    if grep -q "$OLD_MODULE" "$file"; then
        echo "  处理: $file"
        # 替换 module 声明
        sed -i.tmp "s|^module $OLD_MODULE_ESCAPED|module $NEW_MODULE_ESCAPED|g" "$file"
        # 替换 require 中的依赖
        sed -i.tmp "s|$OLD_MODULE_ESCAPED|$NEW_MODULE_ESCAPED|g" "$file"
        # 替换 replace 指令
        sed -i.tmp "s|$OLD_MODULE_ESCAPED|$NEW_MODULE_ESCAPED|g" "$file"
        rm -f "$file.tmp"
    fi
done

# 2. 替换 Go 源代码中的 import 路径
echo -e "${GREEN}[2/4] 替换 Go 源代码中的 import 路径...${NC}"
find . -name "*.go" -type f \
    ! -path "./.git/*" \
    ! -path "./kitex_gen/*" \
    ! -path "./output/*" \
    ! -path "./vendor/*" \
    ! -path "./node_modules/*" | while read -r file; do
    if grep -q "$OLD_MODULE" "$file"; then
        echo "  处理: $file"
        sed -i.tmp "s|$OLD_MODULE_ESCAPED|$NEW_MODULE_ESCAPED|g" "$file"
        rm -f "$file.tmp"
    fi
done

# 3. 替换脚本文件中的 module 路径引用
echo -e "${GREEN}[3/4] 替换脚本文件中的 module 路径...${NC}"
find . -type f \( -name "*.sh" -o -name "*.yaml" -o -name "*.yml" \) \
    ! -path "./.git/*" \
    ! -path "./output/*" \
    ! -path "./vendor/*" | while read -r file; do
    if grep -q "$OLD_MODULE" "$file"; then
        echo "  处理: $file"
        sed -i.tmp "s|$OLD_MODULE_ESCAPED|$NEW_MODULE_ESCAPED|g" "$file"
        rm -f "$file.tmp"
    fi
done

# 4. 更新 go.work 文件（如果需要）
echo -e "${GREEN}[4/4] 检查 go.work 文件...${NC}"
if [ -f "go.work" ] && grep -q "$OLD_MODULE" "go.work"; then
    echo "  处理: go.work"
    sed -i.tmp "s|$OLD_MODULE_ESCAPED|$NEW_MODULE_ESCAPED|g" "go.work"
    rm -f "go.work.tmp"
fi

# 5. 设置 upstream 远程（如果提供了脚手架仓库 URL）
if [ -n "$SCAFFOLD_REPO" ]; then
    echo -e "${GREEN}[5/5] 设置 upstream 远程仓库...${NC}"
    if git remote | grep -q "^upstream$"; then
        echo "  upstream 远程已存在，更新 URL..."
        git remote set-url upstream "$SCAFFOLD_REPO"
    else
        echo "  添加 upstream 远程: $SCAFFOLD_REPO"
        git remote add upstream "$SCAFFOLD_REPO"
    fi
    echo -e "${GREEN}✓ upstream 远程已设置${NC}"
    echo ""
    echo -e "${YELLOW}提示：${NC} 以后可以使用 ./scripts/update.sh 从脚手架仓库拉取更新"
fi

echo ""
echo -e "${GREEN}✓ 所有文件已更新完成！${NC}"
echo ""
echo -e "${YELLOW}下一步操作：${NC}"
echo "  1. 运行 'go mod tidy' 更新依赖"
echo "  2. 如果使用了代码生成，需要重新生成代码："
echo "     - RPC 服务: cd rpc/identity_srv && ./script/gen_kitex_code.sh"
echo "     - HTTP 网关: cd gateway && ./script/gen_hertz_code.sh"
echo "  3. 重新生成 Wire 依赖注入代码："
echo "     - RPC 服务: cd rpc/identity_srv/wire && wire"
echo "     - HTTP 网关: cd gateway/internal/wire && wire"
echo ""
echo -e "${BLUE}提示：${NC}"
echo "  如果遇到问题，可以使用 git 恢复更改："
echo "    git checkout -- ."
echo ""
