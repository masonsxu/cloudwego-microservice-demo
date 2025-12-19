# 安全政策

## 支持的版本

| 版本 | 支持状态 |
|------|----------|
| 0.1.x | :white_check_mark: 支持 |

## 报告漏洞

我们非常重视安全问题。如果你发现了安全漏洞，请按以下方式报告：

### 报告方式

**请勿在公开的 Issue 中报告安全漏洞。**

请通过以下方式私密报告：

1. **GitHub Security Advisories**（推荐）
   - 访问 [Security Advisories](https://github.com/masonsxu/cloudwego-scaffold/security/advisories/new)
   - 创建新的安全公告

2. **邮件报告**
   - 发送邮件至：[security@example.com]
   - 标题注明：[Security] CloudWeGo Scaffold 安全漏洞报告

### 报告内容

请在报告中包含：

- 漏洞描述
- 复现步骤
- 影响范围
- 可能的修复建议（如有）

### 响应时间

- **确认收到**：1-2 个工作日内
- **初步评估**：5 个工作日内
- **修复计划**：根据严重程度，7-30 天内

### 披露政策

- 我们会在修复发布后公开致谢（除非你希望保持匿名）
- 在修复发布前，请勿公开披露漏洞详情
- 我们会尽快发布安全更新

## 安全最佳实践

使用本项目时，请注意以下安全事项：

### 生产环境必须修改

- [ ] `JWT_SIGNING_KEY` - 使用强随机密钥（32+ 字符）
- [ ] `DB_PASSWORD` - 使用强密码
- [ ] `LOGO_STORAGE_SECRET_KEY` - 使用强密钥

### 推荐配置

```env
# 启用安全选项
JWT_COOKIE_SECURE_COOKIE=true   # 需要 HTTPS
JWT_COOKIE_HTTP_ONLY=true       # 防止 XSS
DB_SSLMODE=require              # 数据库 SSL
```

### 不要做的事

- 不要在版本控制中提交 `.env` 文件
- 不要使用默认密码
- 不要在日志中打印敏感信息
- 不要禁用安全中间件

## 安全更新

安全更新将通过以下渠道发布：

1. GitHub Releases
2. CHANGELOG.md
3. GitHub Security Advisories

建议关注项目的 Release 通知以获取安全更新。

## 致谢

感谢以下安全研究人员的贡献：

<!-- 将在此处列出负责任披露漏洞的研究人员 -->

---

感谢你帮助保护 CloudWeGo Scaffold 和其用户的安全！
