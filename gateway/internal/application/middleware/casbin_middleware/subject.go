package casbin_middleware

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	authctx "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/application/context/auth_context"
)

// SubjectExtractor 主体提取器
// 从Hertz请求上下文中提取Casbin主体标识符
type SubjectExtractor struct {
	logger *slog.Logger
}

// NewSubjectExtractor 创建主体提取器
func NewSubjectExtractor(logger *slog.Logger) *SubjectExtractor {
	return &SubjectExtractor{
		logger: logger,
	}
}

// Extract 实现casbin.LookupHandler接口
// 从认证上下文中提取用户名作为Casbin主体
func (e *SubjectExtractor) Extract(ctx context.Context, c *app.RequestContext) string {
	// 获取认证上下文
	authContext, exists := authctx.GetAuthContext(c)
	if !exists {
		e.logger.Debug("未找到认证上下文，返回空主体",
			slog.String("path", string(c.Request.URI().Path())),
			slog.String("method", string(c.Request.Header.Method())))
		return ""
	}

	// 提取用户名
	username, hasUsername := authContext.GetUsername()
	if !hasUsername || username == "" {
		e.logger.Debug("认证上下文中未找到用户名，返回空主体",
			slog.String("path", string(c.Request.URI().Path())),
			slog.String("method", string(c.Request.Header.Method())))
		return ""
	}

	// 验证用户名格式
	if !e.isValidUsername(username) {
		e.logger.Warn("无效的用户名格式",
			slog.String("username", username),
			slog.String("path", string(c.Request.URI().Path())))
		return ""
	}

	// 记录主体提取成功
	e.logger.Debug("成功提取Casbin主体",
		slog.String("username", username),
		slog.String("path", string(c.Request.URI().Path())),
		slog.String("method", string(c.Request.Header.Method())))

	return username
}

// ExtractDomain 从认证上下文中提取域信息
// 用于多租户权限控制
func (e *SubjectExtractor) ExtractDomain(ctx context.Context, c *app.RequestContext) string {
	// 获取认证上下文
	// authContext, exists := authctx.GetAuthContext(c)
	// if !exists {
	// 	return permission.WildcardDomain // 默认使用通配符域
	// }

	// // 提取组织ID
	// orgID, hasOrgID := authContext.GetOrganizationID()
	// if !hasOrgID || orgID == "" {
	// 	return permission.WildcardDomain
	// }

	// // 检查是否为超级管理员
	// username, hasUsername := authContext.GetUsername()
	// if hasUsername && username == permission.SuperAdminUsername {
	// 	return permission.WildcardDomain // 超级管理员使用通配符域
	// }

	return "*" // 默认返回通配符域
}

// isValidUsername 验证用户名格式
func (e *SubjectExtractor) isValidUsername(username string) bool {
	if username == "" {
		return false
	}

	// // 超级管理员特殊处理
	// if username == permission.SuperAdminUsername {
	// 	return true
	// }

	// 一般用户名格式验证：不能包含特殊字符
	if strings.ContainsAny(username, " \t\n\r\v\f") {
		return false
	}

	// 用户名长度检查
	if len(username) < 1 || len(username) > 100 {
		return false
	}

	return true
}

// FormatSubject 格式化主体标识符
// 为多租户场景提供统一的主体格式化
func FormatSubject(username, domain string) string {
	// if username == permission.SuperAdminUsername {
	// 	return username // 超级管理员不需要域后缀
	// }

	// if domain == "" || domain == permission.WildcardDomain {
	// 	return username
	// }

	return fmt.Sprintf("%s@%s", username, domain)
}

// ParseSubject 解析主体标识符
// 从格式化的主体标识符中解析出用户名和域
func ParseSubject(subject string) (username, domain string) {
	// if subject == permission.SuperAdminUsername {
	// 	return subject, permission.WildcardDomain
	// }

	parts := strings.SplitN(subject, "@", 2)
	// if len(parts) == 1 {
	// 	return parts[0], permission.WildcardDomain
	// }

	return parts[0], parts[1]
}
