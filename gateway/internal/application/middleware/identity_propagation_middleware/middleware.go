package middleware

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/masonsxu/cloudwego-microservice-demo/iamclient"
)

// IdentityPropagationImpl 身份传播中间件实现。无状态，可全局复用。
type IdentityPropagationImpl struct{}

// NewIdentityPropagation 创建中间件实例。
func NewIdentityPropagation() IdentityPropagationService {
	return &IdentityPropagationImpl{}
}

// MiddlewareFunc 把请求 header 中 6 个身份字段写入 metainfo 持久值，再 c.Next。
//
// metainfo key 与 HTTP header 镜像（见 iamclient/header.go），下游 RPC 服务
// 通过 metainfo.GetPersistentValue(ctx, iamclient.MetaUserID) 读取。
//
// 已存在持久值时不覆盖（兼容上游已写入的链路）。
func (m *IdentityPropagationImpl) MiddlewareFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ctx = setIfPresent(ctx, iamclient.MetaUserID, string(c.Request.Header.Peek(iamclient.HeaderUserID)))
		ctx = setIfPresent(ctx, iamclient.MetaUserName, string(c.Request.Header.Peek(iamclient.HeaderUserName)))
		ctx = setIfPresent(ctx, iamclient.MetaTenantID, string(c.Request.Header.Peek(iamclient.HeaderTenantID)))
		ctx = setIfPresent(ctx, iamclient.MetaUserRoles, string(c.Request.Header.Peek(iamclient.HeaderUserRoles)))
		ctx = setIfPresent(ctx, iamclient.MetaJTI, string(c.Request.Header.Peek(iamclient.HeaderJTI)))

		c.Next(ctx)
	}
}

// setIfPresent 仅在 value 非空时写入 metainfo，避免污染下游链路。
func setIfPresent(ctx context.Context, key, value string) context.Context {
	if value == "" {
		return ctx
	}

	return metainfo.WithPersistentValue(ctx, key, value)
}
