package casbin_middleware

// import (
// 	"context"
// 	"fmt"
// 	"log/slog"
// 	"strings"
// 	"time"

// 	"github.com/cloudwego/hertz/pkg/app"

// 	"github.com/masonsxu/cloudwego-microservice-demo/gateway/biz/model/http_base"
// 	casbinManager "github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/casbin"
// 	"github.com/masonsxu/cloudwego-microservice-demo/gateway/internal/infrastructure/errors"
// )

// // casbinMiddleware Casbin权限中间件实现
// type casbinMiddleware struct {
// 	casbinManager    *casbinManager.CasbinManager
// 	subjectExtractor *SubjectExtractor
// 	config           *PermissionConfig
// 	logger           *slog.Logger
// }

// // NewCasbinMiddleware 创建权限中间件实例
// func NewCasbinMiddleware(
// 	casbinMgr *casbinManager.CasbinManager,
// 	logger *slog.Logger,
// 	config *PermissionConfig,
// ) (CasbinMiddleware, error) {
// 	if casbinMgr == nil {
// 		return nil, fmt.Errorf("casbinManager不能为空")
// 	}
// 	if logger == nil {
// 		return nil, fmt.Errorf("logger不能为空")
// 	}

// 	if config == nil {
// 		config = DefaultPermissionConfig()
// 	}

// 	// 创建主体提取器
// 	subjectExtractor := NewSubjectExtractor(logger)

// 	impl := &casbinMiddleware{
// 		casbinManager:    casbinMgr,
// 		subjectExtractor: subjectExtractor,
// 		config:           config,
// 		logger:           logger,
// 	}

// 	// 配置自定义错误处理器
// 	impl.setupErrorHandlers()

// 	return impl, nil
// }

// // RequiresPermissions 实现CasbinMiddleware接口
// func (m *casbinMiddleware) RequiresPermissions(permission string) []app.HandlerFunc {
// 	return []app.HandlerFunc{
// 		m.createPermissionHandler(permission),
// 	}
// }

// // RequiresRoles 实现CasbinMiddleware接口
// func (m *casbinMiddleware) RequiresRoles(roles string) []app.HandlerFunc {
// 	return []app.HandlerFunc{
// 		m.createRoleHandler(roles),
// 	}
// }

// // RequiresAnyPermissions 实现CasbinMiddleware接口
// func (m *casbinMiddleware) RequiresAnyPermissions(permissions ...string) []app.HandlerFunc {
// 	return []app.HandlerFunc{
// 		m.createAnyPermissionsHandler(permissions),
// 	}
// }

// // RequiresAllPermissions 实现CasbinMiddleware接口
// func (m *casbinMiddleware) RequiresAllPermissions(permissions ...string) []app.HandlerFunc {
// 	return []app.HandlerFunc{
// 		m.createAllPermissionsHandler(permissions),
// 	}
// }

// // HasPermission 实现CasbinMiddleware接口
// func (m *casbinMiddleware) HasPermission(ctx context.Context, userID, domain, resource, action string) (bool, error) {
// 	return m.casbinManager.HasPermission(userID, domain, resource, action)
// }

// // HasRole 实现CasbinMiddleware接口
// func (m *casbinMiddleware) HasRole(ctx context.Context, userID, domain, role string) (bool, error) {
// 	subject := FormatSubject(userID, domain)
// 	hasRole, err := m.casbinManager.Enforcer.HasRoleForUser(subject, role)
// 	if err != nil {
// 		m.logger.Error("检查用户角色失败",
// 			slog.String("userID", userID),
// 			slog.String("domain", domain),
// 			slog.String("role", role),
// 			slog.String("error", err.Error()))
// 		return false, err
// 	}
// 	return hasRole, nil
// }

// // GetUserPermissions 实现CasbinMiddleware接口
// func (m *casbinMiddleware) GetUserPermissions(ctx context.Context, userID, domain string) ([]string, error) {
// 	permissions, err := m.casbinManager.GetUserPermissions(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	result := make([]string, len(permissions))
// 	for i, perm := range permissions {
// 		// permissions是[][]string类型，每个perm是[resource, action]格式
// 		if len(perm) >= 2 {
// 			result[i] = fmt.Sprintf("%s:%s", perm[0], perm[1])
// 		}
// 	}
// 	return result, nil
// }

// // GetUserRoles 实现CasbinMiddleware接口
// func (m *casbinMiddleware) GetUserRoles(ctx context.Context, userID, domain string) ([]string, error) {
// 	return m.casbinManager.GetUserRoles(userID)
// }

// // createPermissionHandler 创建权限检查处理器
// func (m *casbinMiddleware) createPermissionHandler(permission string) app.HandlerFunc {
// 	return func(ctx context.Context, c *app.RequestContext) {
// 		// 跳过路径检查
// 		if m.shouldSkipPath(string(c.Request.URI().Path())) {
// 			c.Next(ctx)
// 			return
// 		}

// 		// 提取主体和域
// 		subject := m.subjectExtractor.Extract(ctx, c)
// 		if subject == "" {
// 			m.handleUnauthorized(ctx, c, "未找到有效的用户身份")
// 			return
// 		}

// 		domain := m.subjectExtractor.ExtractDomain(ctx, c)

// 		// 解析权限字符串
// 		parts := strings.Split(permission, ":")
// 		if len(parts) != 2 {
// 			m.logger.Error("无效的权限字符串格式",
// 				slog.String("permission", permission))
// 			m.handleForbidden(ctx, c, "权限配置错误")
// 			return
// 		}
// 		resource, action := parts[0], parts[1]

// 		// 执行权限检查
// 		allowed, err := m.casbinManager.HasPermission(subject, domain, resource, action)
// 		if err != nil {
// 			m.logger.Error("权限检查失败",
// 				slog.String("subject", subject),
// 				slog.String("domain", domain),
// 				slog.String("resource", resource),
// 				slog.String("action", action),
// 				slog.String("error", err.Error()))
// 			m.handleForbidden(ctx, c, "权限检查失败")
// 			return
// 		}

// 		if !allowed {
// 			m.logger.Warn("权限检查失败：用户权限不足",
// 				slog.String("subject", subject),
// 				slog.String("domain", domain),
// 				slog.String("permission", permission),
// 				slog.String("path", string(c.Request.URI().Path())))
// 			m.handleForbidden(ctx, c, "权限不足")
// 			return
// 		}

// 		// 权限检查通过，继续处理请求
// 		m.logger.Debug("权限检查通过",
// 			slog.String("subject", subject),
// 			slog.String("domain", domain),
// 			slog.String("permission", permission))

// 		c.Next(ctx)
// 	}
// }

// // createRoleHandler 创建角色检查处理器
// func (m *casbinMiddleware) createRoleHandler(roles string) app.HandlerFunc {
// 	roleList := strings.Split(roles, ",")
// 	for i := range roleList {
// 		roleList[i] = strings.TrimSpace(roleList[i])
// 	}

// 	return func(ctx context.Context, c *app.RequestContext) {
// 		// 跳过路径检查
// 		if m.shouldSkipPath(string(c.Request.URI().Path())) {
// 			c.Next(ctx)
// 			return
// 		}

// 		// 提取主体和域
// 		subject := m.subjectExtractor.Extract(ctx, c)
// 		if subject == "" {
// 			m.handleUnauthorized(ctx, c, "未找到有效的用户身份")
// 			return
// 		}

// 		domain := m.subjectExtractor.ExtractDomain(ctx, c)

// 		// 检查用户是否拥有任一指定角色
// 		for _, role := range roleList {
// 			hasRole, err := m.HasRole(ctx, subject, domain, role)
// 			if err != nil {
// 				m.logger.Error("角色检查失败",
// 					slog.String("subject", subject),
// 					slog.String("role", role),
// 					slog.String("error", err.Error()))
// 				continue
// 			}

// 			if hasRole {
// 				m.logger.Debug("角色检查通过",
// 					slog.String("subject", subject),
// 					slog.String("role", role))
// 				c.Next(ctx)
// 				return
// 			}
// 		}

// 		m.logger.Warn("角色检查失败：用户无所需角色",
// 			slog.String("subject", subject),
// 			slog.String("required_roles", roles))
// 		m.handleForbidden(ctx, c, "角色权限不足")
// 	}
// }

// // createAnyPermissionsHandler 创建任一权限检查处理器
// func (m *casbinMiddleware) createAnyPermissionsHandler(permissions []string) app.HandlerFunc {
// 	return func(ctx context.Context, c *app.RequestContext) {
// 		// 跳过路径检查
// 		if m.shouldSkipPath(string(c.Request.URI().Path())) {
// 			c.Next(ctx)
// 			return
// 		}

// 		// 提取主体和域
// 		subject := m.subjectExtractor.Extract(ctx, c)
// 		if subject == "" {
// 			m.handleUnauthorized(ctx, c, "未找到有效的用户身份")
// 			return
// 		}

// 		// domain := m.subjectExtractor.ExtractDomain(ctx, c)

// 		// 检查用户是否拥有任一指定权限
// 		// for _, perm := range permissions {
// 		// 	resource, action, err := permission.ParsePermissionString(perm)
// 		// 	if err != nil {
// 		// 		m.logger.Error("无效的权限字符串格式",
// 		// 			slog.String("permission", perm),
// 		// 			slog.String("error", err.Error()))
// 		// 		continue
// 		// 	}

// 		// 	allowed, err := m.casbinManager.HasPermission(subject, domain, resource, action)
// 		// 	if err != nil {
// 		// 		m.logger.Error("权限检查失败",
// 		// 			slog.String("subject", subject),
// 		// 			slog.String("permission", perm),
// 		// 			slog.String("error", err.Error()))
// 		// 		continue
// 		// 	}

// 		// 	if allowed {
// 		// 		m.logger.Debug("任一权限检查通过",
// 		// 			slog.String("subject", subject),
// 		// 			slog.String("permission", perm))
// 		// 		c.Next(ctx)
// 		// 		return
// 		// 	}
// 		// }

// 		m.logger.Warn("任一权限检查失败：用户无所需权限",
// 			slog.String("subject", subject),
// 			slog.Any("required_permissions", permissions))
// 		m.handleForbidden(ctx, c, "权限不足")
// 	}
// }

// // createAllPermissionsHandler 创建所有权限检查处理器
// func (m *casbinMiddleware) createAllPermissionsHandler(permissions []string) app.HandlerFunc {
// 	return func(ctx context.Context, c *app.RequestContext) {
// 		// 跳过路径检查
// 		if m.shouldSkipPath(string(c.Request.URI().Path())) {
// 			c.Next(ctx)
// 			return
// 		}

// 		// 提取主体和域
// 		subject := m.subjectExtractor.Extract(ctx, c)
// 		if subject == "" {
// 			m.handleUnauthorized(ctx, c, "未找到有效的用户身份")
// 			return
// 		}

// 		// domain := m.subjectExtractor.ExtractDomain(ctx, c)

// 		// // 检查用户是否拥有所有指定权限
// 		// for _, perm := range permissions {
// 		// 	resource, action, err := permission.ParsePermissionString(perm)
// 		// 	if err != nil {
// 		// 		m.logger.Error("无效的权限字符串格式",
// 		// 			slog.String("permission", perm),
// 		// 			slog.String("error", err.Error()))
// 		// 		m.handleForbidden(ctx, c, "权限配置错误")
// 		// 		return
// 		// 	}

// 		// 	allowed, err := m.casbinManager.HasPermission(subject, domain, resource, action)
// 		// 	if err != nil {
// 		// 		m.logger.Error("权限检查失败",
// 		// 			slog.String("subject", subject),
// 		// 			slog.String("permission", perm),
// 		// 			slog.String("error", err.Error()))
// 		// 		m.handleForbidden(ctx, c, "权限检查失败")
// 		// 		return
// 		// 	}

// 		// 	if !allowed {
// 		// 		m.logger.Warn("所有权限检查失败：缺少必需权限",
// 		// 			slog.String("subject", subject),
// 		// 			slog.String("missing_permission", perm))
// 		// 		m.handleForbidden(ctx, c, "权限不足")
// 		// 		return
// 		// 	}
// 		// }

// 		m.logger.Debug("所有权限检查通过",
// 			slog.String("subject", subject),
// 			slog.Any("permissions", permissions))
// 		c.Next(ctx)
// 	}
// }

// // shouldSkipPath 检查是否应该跳过权限检查
// func (m *casbinMiddleware) shouldSkipPath(path string) bool {
// 	for _, skipPath := range m.config.SkipPaths {
// 		if path == skipPath {
// 			return true
// 		}
// 	}
// 	return false
// }

// // setupErrorHandlers 设置错误处理器
// func (m *casbinMiddleware) setupErrorHandlers() {
// 	if m.config.UnauthorizedHandler == nil {
// 		m.config.UnauthorizedHandler = m.defaultUnauthorizedHandler
// 	}
// 	if m.config.ForbiddenHandler == nil {
// 		m.config.ForbiddenHandler = m.defaultForbiddenHandler
// 	}
// }

// // handleUnauthorized 处理未认证错误
// func (m *casbinMiddleware) handleUnauthorized(ctx context.Context, c *app.RequestContext, message string) {
// 	m.logger.Warn("用户未认证",
// 		slog.String("path", string(c.Request.URI().Path())),
// 		slog.String("method", string(c.Request.Header.Method())),
// 		slog.String("message", message))

// 	if m.config.UnauthorizedHandler != nil {
// 		m.config.UnauthorizedHandler(ctx, c)
// 		return
// 	}

// 	m.defaultUnauthorizedHandler(ctx, c)
// }

// // handleForbidden 处理权限不足错误
// func (m *casbinMiddleware) handleForbidden(ctx context.Context, c *app.RequestContext, message string) {
// 	m.logger.Warn("权限不足",
// 		slog.String("path", string(c.Request.URI().Path())),
// 		slog.String("method", string(c.Request.Header.Method())),
// 		slog.String("message", message))

// 	if m.config.ForbiddenHandler != nil {
// 		m.config.ForbiddenHandler(ctx, c)
// 		return
// 	}

// 	m.defaultForbiddenHandler(ctx, c)
// }

// // defaultUnauthorizedHandler 默认未认证处理器
// func (m *casbinMiddleware) defaultUnauthorizedHandler(ctx context.Context, c *app.RequestContext) {
// 	// 生成标准化的错误响应
// 	// 权限不足异常
// 	apiError := errors.ErrUnauthorized

// 	httpStatus := errors.GetHTTPStatus(apiError.Code())
// 	requestID := errors.GenerateRequestID(c)
// 	traceID := errors.GenerateTraceID(c)
// 	timestamp := time.Now().UnixMilli()

// 	response := &http_base.OperationStatusResponseDTO{
// 		BaseResp: &http_base.BaseResponseDTO{
// 			Code:      apiError.Code(),
// 			Message:   apiError.Message(),
// 			RequestID: &requestID,
// 			TraceID:   &traceID,
// 			Timestamp: &timestamp,
// 		},
// 	}

// 	// 直接写入响应
// 	c.JSON(httpStatus, response)
// 	c.Abort()
// }

// // defaultForbiddenHandler 默认权限不足处理器
// func (m *casbinMiddleware) defaultForbiddenHandler(ctx context.Context, c *app.RequestContext) {
// 	// 生成标准化的错误响应
// 	// 权限不足异常
// 	apiError := errors.ErrForbidden

// 	httpStatus := errors.GetHTTPStatus(apiError.Code())
// 	requestID := errors.GenerateRequestID(c)
// 	traceID := errors.GenerateTraceID(c)
// 	timestamp := time.Now().UnixMilli()

// 	response := &http_base.OperationStatusResponseDTO{
// 		BaseResp: &http_base.BaseResponseDTO{
// 			Code:      apiError.Code(),
// 			Message:   apiError.Message(),
// 			RequestID: &requestID,
// 			TraceID:   &traceID,
// 			Timestamp: &timestamp,
// 		},
// 	}

// 	// 直接写入响应
// 	c.JSON(httpStatus, response)
// 	c.Abort()
// }
