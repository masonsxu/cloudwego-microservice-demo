---
name: codegen
description: 执行代码生成（Kitex/Hertz/Wire，支持自动检测）
argument-hint: "[kitex|hertz|wire|all]"
---

# 执行代码生成

统一入口执行 Kitex / Hertz / Wire 代码生成，支持自动检测或手动指定。

**用户输入**: $ARGUMENTS（可选：`kitex` / `hertz` / `wire` / `all`，为空则自动检测）

---

## 执行流程

### Step 1: 确定生成范围

根据 `$ARGUMENTS` 决定执行哪些生成步骤：

- **`kitex`**：仅 Kitex RPC 代码生成
- **`hertz`**：仅 Hertz HTTP 代码生成
- **`wire`**：仅 Wire 依赖注入生成
- **`all`**：依次执行 Kitex → Hertz → Wire
- **为空**：通过 `git diff --name-only` 自动检测变更文件，判断需要执行哪些步骤：
  - `idl/rpc/` 下有变更 → 执行 Kitex
  - `idl/http/` 下有变更 → 执行 Hertz
  - `wire/provider.go` 或 `wire/` 相关文件有变更 → 执行 Wire

### Step 2: 执行 Kitex 代码生成（如需）

```bash
cd rpc/identity_srv && ./script/gen_kitex_code.sh
```

确认：
- `kitex_gen/` 下文件已更新
- 无报错输出

### Step 3: 执行 Hertz 代码生成（如需）

```bash
cd gateway && ./script/gen_hertz_code.sh
```

或指定服务：
```bash
cd gateway && ./script/gen_hertz_code.sh identity
cd gateway && ./script/gen_hertz_code.sh permission
```

确认：
- `biz/handler/` 和 `biz/router/` 下文件已更新
- 无报错输出

### Step 4: 执行 Wire 生成（如需）

RPC 服务：
```bash
cd rpc/identity_srv/wire && wire
```

网关服务：
```bash
cd gateway/internal/wire && wire
```

确认 `wire_gen.go` 正确生成。

### Step 5: 后处理验证

依次执行：

```bash
# Go mod 整理
cd rpc/identity_srv && go mod tidy
cd gateway && go mod tidy

# 编译验证
cd rpc/identity_srv && go build ./...
cd gateway && go build ./...
```

### Step 6: 输出报告

汇总报告：
- 执行了哪些生成步骤
- 各步骤是否成功
- 编译验证结果
- 如有失败，给出修复建议（可参考 `docs/03-部署运维/故障排查.md`）

---

## 注意事项

- **禁止手动修改** `kitex_gen/` 目录下的任何文件
- 代码生成脚本需要 `kitex`、`thriftgo`、`hz` 命令行工具已安装
- 如果遇到代码生成错误，读取 `docs/03-部署运维/故障排查.md` 中的排查指南
