---
name: diagnose
description: 根据错误信息进行分类诊断，定位根因并给出修复方案
argument-hint: "<error message>"
---

# 故障诊断

根据错误信息进行分类诊断，定位根因并给出修复方案。

**用户输入**: $ARGUMENTS（错误信息、日志片段或问题描述）

---

## 执行流程

### Step 1: 读取排查指南

先读取 `docs/03-部署运维/故障排查.md` 获取完整排查参考。

### Step 2: 错误分类

分析 `$ARGUMENTS` 中的错误信息，归类为以下类别：

#### 编译错误
- 特征：`cannot find package`、`undefined`、`type mismatch`
- 排查方向：
  1. 检查 `go.work` 中模块路径是否正确
  2. 检查 `go.mod` 依赖版本
  3. 是否需要重新运行代码生成
  4. `go mod tidy` 是否能解决

#### Wire 依赖注入错误
- 特征：`no provider found`、`multiple bindings`、`cycle detected`
- 排查方向：
  1. 检查 `wire/provider.go` 中 Provider 是否完整
  2. 检查返回类型是否为接口
  3. 确认 Provider Set 是否包含所有依赖
  4. 运行 `wire` 查看详细错误

#### 连接/运行时错误
- 特征：`connection refused`、`timeout`、`dial tcp`
- 排查方向：
  1. 检查基础设施是否启动：`cd docker && podman-compose ps`
  2. 检查端口：PostgreSQL(5432)、etcd(2379)、Redis(6379)、RPC(8891)、Gateway(8080)
  3. 检查环境变量配置（读取 `docs/01-快速入门/配置参考.md`）
  4. 检查服务注册发现（etcd）

#### Lint/格式错误
- 特征：`golangci-lint` 输出的告警
- 排查方向：
  1. 导入顺序：运行 `gci write .`
  2. 行长度：运行 `golines` 或 `golangci-lint format`
  3. 其他 lint 规则：检查 `.golangci.yml` 配置

#### 代码生成错误
- 特征：Kitex/Hertz/Wire 生成脚本报错
- 排查方向：
  1. IDL 语法检查
  2. 工具版本检查（`kitex --version`、`hz version`）
  3. 依赖 include 路径是否正确

#### 权限/认证错误
- 特征：`401 Unauthorized`、`403 Forbidden`、Casbin 相关
- 排查方向：
  1. 读取 `docs/04-权限管理/权限管理设计.md`
  2. 检查 JWT Token 配置
  3. 检查 Casbin 策略配置
  4. 检查中间件执行顺序

### Step 3: 深入排查

根据分类结果：
1. 读取相关代码文件
2. 检查配置和环境
3. 对照文档中的排查步骤逐一验证

### Step 4: 输出诊断报告

报告包含：
- **错误分类**：属于哪类问题
- **根因分析**：问题的根本原因
- **修复方案**：具体的修复步骤（代码改动或配置调整）
- **预防建议**：避免同类问题再次发生的建议

---

## 注意事项

- 不确定时多读文档，尤其是 `docs/03-部署运维/故障排查.md`
- 配置相关问题参考 `docs/01-快速入门/配置参考.md`
- 权限相关问题参考 `docs/04-权限管理/权限管理设计.md`
