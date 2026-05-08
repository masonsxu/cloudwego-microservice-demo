package iamclient

import (
	"context"
	"errors"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

// Subject 是经过认证的调用主体。
//
// 由网关 jwt_middleware 验签 JWT 后，把身份字段以 HTTP header 形式注入下游
// 请求；业务侧通过 Client.SubjectFromHeader / SubjectFromContext 还原此结构，
// 再调 Subject.Check / MustCheck 做 PDP 鉴权。
//
// 严禁字段扩展（部门、数据范围、权限码列表等业务字段必须通过 PDP 决策实时计算，
// 不在 Subject 中承载——参考提案 §2 「铁律」）。
type Subject struct {
	UserID    string
	UserName  string
	TenantID  string
	Roles     []string
	RequestID string
	Jti       string

	client *Client
}

// HeaderGetter 是从 HTTP header 读取字段的最小接口。
//
// 标准库 net/http.Header 与 hertz protocol.RequestHeader 都满足此接口
// （都提供 Get(name string) string），iamclient 借此避免对具体框架的依赖。
type HeaderGetter interface {
	Get(name string) string
}

// ErrMissingUserID 表示 header / metadata 中没有有效的 user_id。
//
// 业务侧应把此错误转成 401 Unauthorized 返回（说明请求未经过网关 jwt_middleware
// 处理，或 token 缺失 sub 字段）。
var ErrMissingUserID = errors.New("iamclient: missing X-User-Id in incoming request")

// SubjectFromHeader 从 HTTP header 还原 Subject。
//
// 业务调用样态：
//
//	subject, err := cli.SubjectFromHeader(c.Request.Header)
//
// 必填字段：X-User-Id（缺失返回 ErrMissingUserID）。
// X-User-Roles 为逗号分隔字符串，会被切分；其余字段缺失视为空值。
func (c *Client) SubjectFromHeader(h HeaderGetter) (*Subject, error) {
	if h == nil {
		return nil, ErrMissingUserID
	}

	userID := h.Get(HeaderUserID)
	if userID == "" {
		return nil, ErrMissingUserID
	}

	return &Subject{
		UserID:    userID,
		UserName:  h.Get(HeaderUserName),
		TenantID:  h.Get(HeaderTenantID),
		Roles:     splitRoles(h.Get(HeaderUserRoles)),
		RequestID: h.Get(HeaderRequestID),
		Jti:       h.Get(HeaderJTI),
		client:    c,
	}, nil
}

// SubjectFromContext 从 Kitex metainfo 还原 Subject。
//
// 适用于下游 RPC 服务（如 identity_srv 内部 handler）；网关在调用下游前会把
// HTTP header 中的身份字段透传到 metainfo 持久值。
func (c *Client) SubjectFromContext(ctx context.Context) (*Subject, error) {
	userID, _ := metainfo.GetPersistentValue(ctx, MetaUserID)
	if userID == "" {
		return nil, ErrMissingUserID
	}

	userName, _ := metainfo.GetPersistentValue(ctx, MetaUserName)
	tenantID, _ := metainfo.GetPersistentValue(ctx, MetaTenantID)
	rolesRaw, _ := metainfo.GetPersistentValue(ctx, MetaUserRoles)
	requestID, _ := metainfo.GetPersistentValue(ctx, MetaRequestID)
	jti, _ := metainfo.GetPersistentValue(ctx, MetaJTI)

	return &Subject{
		UserID:    userID,
		UserName:  userName,
		TenantID:  tenantID,
		Roles:     splitRoles(rolesRaw),
		RequestID: requestID,
		Jti:       jti,
		client:    c,
	}, nil
}

// splitRoles 把 "doctor,admin" 切成 ["doctor", "admin"]，去除空项与首尾空白。
func splitRoles(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
