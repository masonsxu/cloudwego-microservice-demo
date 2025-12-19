#!/usr/bin/env bash

# build.sh: 编译 Hertz API 网关服务。
# 此脚本遵循 shell 脚本的最佳实践，包括错误处理、
# 清晰的变量定义和信息丰富的输出。

# 如果命令以非零状态退出，立即退出。
set -eou pipefail

# --- 配置 ---
# 要构建的可执行文件的名称。
readonly RUN_NAME="hertz_service"
# 服务的根目录，根据脚本的位置相对确定。
readonly SERVICE_ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# 构建产物的输出目录。
readonly OUTPUT_DIR="${SERVICE_ROOT_DIR}/output"
# 最终可执行文件的路径。
readonly BINARY_PATH="${OUTPUT_DIR}/bin/${RUN_NAME}"

# --- 主要构建逻辑 ---

# 导航到服务的根目录以确保所有路径都正确。
cd "${SERVICE_ROOT_DIR}"

echo "INFO: 开始为 ${RUN_NAME} 构建..."
echo "INFO: 服务根目录: ${SERVICE_ROOT_DIR}"

# 步骤 1: 生成 Swagger 文档。
# 此命令从源代码注释生成 API 文档。
echo "INFO: 生成 Swagger 文档..."
if ! command -v swag &> /dev/null; then
    echo "ERROR: 未找到 'swag' 命令。请运行以下命令安装: go install github.com/swaggo/swag/cmd/swag@latest" >&2
    exit 1
fi
swag init --parseDependencyLevel 1 --useStructName true

# 步骤 2: 清理之前的构建产物并创建输出目录。
echo "INFO: 清理之前的构建并在 ${OUTPUT_DIR}/bin 创建输出目录..."
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}/bin"

# 步骤 3: 将必要的脚本复制到输出目录。
# 这将复制运行服务所需的引导脚本。
echo "INFO: 复制引导脚本..."
cp "${SERVICE_ROOT_DIR}/script/bootstrap.sh" "${OUTPUT_DIR}/"
chmod +x "${OUTPUT_DIR}/bootstrap.sh"

# 步骤 4: 为生产环境构建 Go 应用程序。
# -ldflags="-s -w" 会剥离调试信息，减小二进制文件的大小。
echo "INFO: 构建 Go 应用程序..."
go build -ldflags="-s -w" -o "${BINARY_PATH}" .

echo "SUCCESS: 构建成功完成。"
echo "INFO: 可执行文件位于 ${BINARY_PATH}"
