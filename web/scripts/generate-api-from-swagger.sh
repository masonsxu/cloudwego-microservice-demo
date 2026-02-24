#!/bin/bash

set -e

echo "========================================="
echo "  从 Swagger 文档生成前端 API 客户端"
echo "========================================="
echo ""

SWAGGER_FILE="../gateway/docs/swagger.yaml"
OUTPUT_DIR="src/api/generated"

# 检查 Swagger 文件是否存在
if [ ! -f "$SWAGGER_FILE" ]; then
    echo "❌ 错误: Swagger 文件不存在: $SWAGGER_FILE"
    echo "请先启动后端服务并确保已生成 Swagger 文档"
    exit 1
fi

echo "✓ 找到 Swagger 文件: $SWAGGER_FILE"
echo ""

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

echo "正在生成 TypeScript 类型定义和 API 客户端..."

# 使用 openapi-typescript-codegen 生成
npx openapi-typescript-codegen \
  --input "$SWAGGER_FILE" \
  --output "$OUTPUT_DIR" \
  --client axios

echo ""
echo "✓ 生成完成！"
echo "  - 类型定义: $OUTPUT_DIR/models"
echo "  - API 客户端: $OUTPUT_DIR/services"
echo ""

echo "正在后处理生成的代码..."
bash scripts/postprocess-generated-api.sh

echo ""
echo "========================================="
echo "  使用方法"
echo "========================================="
echo ""
echo "import { Service } from '@/api/generated'"
echo "import type { LoginResponseDTO } from '@/api/generated/models'"
echo ""
echo "const api = new Service({"
echo "  baseURL: import.meta.env.VITE_API_BASE_URL"
echo "})"
echo ""
echo "const response = await api.identity.login({"
echo "  requestBody: { username, password }"
echo "})"
echo ""
