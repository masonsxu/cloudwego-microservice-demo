#!/bin/bash

# 测试覆盖率报告生成脚本
# 用途：生成项目的测试覆盖率报告和可视化

set -e

echo "=========================================="
echo "  测试覆盖率报告生成工具"
echo "=========================================="
echo

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 函数：打印带颜色的消息
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    print_error "Go 未安装，请先安装 Go"
    exit 1
fi

print_success "Go 已安装: $(go version)"

# 创建输出目录（使用绝对路径）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
OUTPUT_DIR="$PROJECT_ROOT/coverage-reports"
mkdir -p "$OUTPUT_DIR"
print_success "创建输出目录: $OUTPUT_DIR"

# 1. 运行 RPC 服务测试
echo
echo "=========================================="
echo "  测试 RPC 服务 (identity_srv)"
echo "=========================================="
cd "$PROJECT_ROOT/rpc/identity_srv"

echo "运行测试..."
if go test ./... -coverprofile="$OUTPUT_DIR/rpc-coverage.out" -v > "$OUTPUT_DIR/rpc-test.log" 2>&1; then
    print_success "RPC 服务测试通过"
else
    print_warning "RPC 服务测试有失败，查看日志: $OUTPUT_DIR/rpc-test.log"
fi

# 生成 RPC 覆盖率报告
echo "生成 RPC 覆盖率报告..."
if [ -f "$OUTPUT_DIR/rpc-coverage.out" ]; then
    go tool cover -func="$OUTPUT_DIR/rpc-coverage.out" > "$OUTPUT_DIR/rpc-coverage-func.txt"
    go tool cover -html="$OUTPUT_DIR/rpc-coverage.out" -o "$OUTPUT_DIR/rpc-coverage.html"

    # 提取总体覆盖率
    RPC_COVERAGE=$(go tool cover -func="$OUTPUT_DIR/rpc-coverage.out" | grep total | awk '{print $3}')
    print_success "RPC 服务总体覆盖率: $RPC_COVERAGE"

    # 列出覆盖率低于 80% 的函数
    echo ""
    echo "RPC 服务覆盖率低于 80% 的模块:"
    go tool cover -func="$OUTPUT_DIR/rpc-coverage.out" | awk '$3 != "100.0%" && $3 + 0 < 80.0 {print $1 " - " $3}' || echo "无"
else
    print_warning "RPC 覆盖率文件未生成，跳过报告生成"
fi

# 2. 运行 Gateway 服务测试
echo
echo "=========================================="
echo "  测试 Gateway 服务"
echo "=========================================="
cd "$PROJECT_ROOT/gateway"

echo "运行测试..."
if go test ./... -coverprofile="$OUTPUT_DIR/gateway-coverage.out" -v > "$OUTPUT_DIR/gateway-test.log" 2>&1; then
    print_success "Gateway 服务测试通过"
else
    print_warning "Gateway 服务测试有失败，查看日志: $OUTPUT_DIR/gateway-test.log"
fi

# 生成 Gateway 覆盖率报告
echo "生成 Gateway 覆盖率报告..."
if [ -f "$OUTPUT_DIR/gateway-coverage.out" ]; then
    go tool cover -func="$OUTPUT_DIR/gateway-coverage.out" > "$OUTPUT_DIR/gateway-coverage-func.txt"
    go tool cover -html="$OUTPUT_DIR/gateway-coverage.out" -o "$OUTPUT_DIR/gateway-coverage.html"

    # 提取总体覆盖率
    GATEWAY_COVERAGE=$(go tool cover -func="$OUTPUT_DIR/gateway-coverage.out" | grep total | awk '{print $3}')
    print_success "Gateway 服务总体覆盖率: $GATEWAY_COVERAGE"

    # 列出覆盖率低于 80% 的函数
    echo ""
    echo "Gateway 服务覆盖率低于 80% 的模块:"
    go tool cover -func="$OUTPUT_DIR/gateway-coverage.out" | awk '$3 != "100.0%" && $3 + 0 < 80.0 {print $1 " - " $3}' || echo "无"
else
    print_warning "Gateway 覆盖率文件未生成，跳过报告生成"
fi

# 3. 生成汇总报告
echo
echo "=========================================="
echo "  生成汇总报告"
echo "=========================================="

cat > "$OUTPUT_DIR/coverage-summary.md" << EOF
# 测试覆盖率报告

生成时间: $(date '+%Y-%m-%d %H:%M:%S')

## 覆盖率总览

| 服务 | 总体覆盖率 | 状态 |
|------|-----------|------|
| RPC 服务 | $RPC_COVERAGE | $(echo "$RPC_COVERAGE" | awk -F'%' '{if ($1+0 >= 80) print "✅"; else if ($1+0 >= 60) print "⚠️"; else print "❌"}') |
| Gateway 服务 | $GATEWAY_COVERAGE | $(echo "$GATEWAY_COVERAGE" | awk -F'%' '{if ($1+0 >= 80) print "✅"; else if ($1+0 >= 60) print "⚠️"; else print "❌"}') |

## 详细报告

### RPC 服务

- [函数级别覆盖率](./rpc-coverage-func.txt)
- [HTML 可视化报告](./rpc-coverage.html)
- [测试日志](./rpc-test.log)

### Gateway 服务

- [函数级别覆盖率](./gateway-coverage-func.txt)
- [HTML 可视化报告](./gateway-coverage.html)
- [测试日志](./gateway-test.log)

## 下一步建议

EOF

# 根据覆盖率给出建议
if [[ "$RPC_COVERAGE" < "60%" ]]; then
    echo "- ❌ RPC 服务覆盖率过低，优先补充核心模块测试" >> "$OUTPUT_DIR/coverage-summary.md"
elif [[ "$RPC_COVERAGE" < "80%" ]]; then
    echo "- ⚠️ RPC 服务覆盖率一般，建议补充未覆盖的分支" >> "$OUTPUT_DIR/coverage-summary.md"
else
    echo "- ✅ RPC 服务覆盖率良好，继续保持" >> "$OUTPUT_DIR/coverage-summary.md"
fi

if [[ "$GATEWAY_COVERAGE" < "60%" ]]; then
    echo "- ❌ Gateway 服务覆盖率过低，优先补充 Handler 和 Service 测试" >> "$OUTPUT_DIR/coverage-summary.md"
elif [[ "$GATEWAY_COVERAGE" < "80%" ]]; then
    echo "- ⚠️ Gateway 服务覆盖率一般，建议补充 Middleware 测试" >> "$OUTPUT_DIR/coverage-summary.md"
else
    echo "- ✅ Gateway 服务覆盖率良好，继续保持" >> "$OUTPUT_DIR/coverage-summary.md"
fi

cat >> "$OUTPUT_DIR/coverage-summary.md" << EOF

## 查看报告

在浏览器中打开 HTML 报告查看可视化覆盖率:

- RPC 服务: file://$OUTPUT_DIR/rpc-coverage.html
- Gateway 服务: file://$OUTPUT_DIR/gateway-coverage.html
EOF

# 4. 输出总结
echo
echo "=========================================="
echo "  测试覆盖率报告生成完成"
echo "=========================================="
echo
echo "报告文件位置: $OUTPUT_DIR/"
echo "  - coverage-summary.md          汇总报告"
echo "  - rpc-coverage.html            RPC 可视化报告"
echo "  - gateway-coverage.html        Gateway 可视化报告"
echo "  - rpc-coverage-func.txt        RPC 函数覆盖率"
echo "  - gateway-coverage-func.txt    Gateway 函数覆盖率"
echo

# 尝试在浏览器中打开报告
if [ -f "$OUTPUT_DIR/rpc-coverage.html" ]; then
    if command -v xdg-open &> /dev/null; then
        echo "尝试在浏览器中打开报告..."
        xdg-open "$OUTPUT_DIR/rpc-coverage.html" &> /dev/null &
    elif command -v open &> /dev/null; then
        echo "尝试在浏览器中打开报告..."
        open "$OUTPUT_DIR/rpc-coverage.html" &> /dev/null &
    fi
fi

print_success "完成！"
