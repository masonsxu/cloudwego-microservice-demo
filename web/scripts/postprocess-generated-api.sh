#!/bin/bash

set -e

echo "========================================="
echo "  后处理生成的 API 客户端代码"
echo "========================================="
echo ""

GENERATED_DIR="src/api/generated"

echo "正在修复 TypeScript 编译问题..."

# 1. 修复 PermissionLevel enum 问题
echo "  - 修复 PermissionLevel enum..."
cat > "$GENERATED_DIR/models/PermissionLevel.ts" << 'EOF'
/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

// Permission levels (使用对象替代 enum 以支持 erasableSyntaxOnly)
export const PermissionLevel = {
  NONE: 0,
  READ: 1,
  WRITE: 2,
  FULL: 3,
} as const

export type PermissionLevel = typeof PermissionLevel[keyof typeof PermissionLevel]
EOF

# 2. 修复其他可能的 enum 文件
echo "  - 修复其他 enum..."
find "$GENERATED_DIR/models" -name "*.ts" -exec grep -l "export enum" {} \; | while read file; do
  echo "    修复: $file"
  sed -i.bak 's/export enum \([A-Za-z]+\)/export const \1 = {/' "$file"
  sed -i.bak 's/} \([A-Z_][A-Z_0-9]*\)/} as const  # \1_type = typeof \1[keyof typeof \1]/' "$file"
  rm -f "${file}.bak"
done

echo ""
echo "✓ 后处理完成！"
echo "  - 生成的代码已适配项目 TypeScript 配置"
echo ""
