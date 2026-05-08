// Package middleware 把 jwt_middleware 注入的身份 HTTP header 镜像到 Kitex
// metainfo 持久值，让下游 RPC 服务（identity_srv 等）通过 TTHeader 自动收到，
// 业务侧再用 iamclient.Client.SubjectFromContext 还原 Subject。
//
// 设计要点：
//   - 必须排在 jwt_middleware 之后（依赖其注入的 X-User-* header）。
//   - 排在 authz_middleware 之前/之后均可，二者无依赖。
//   - 命中 public 路径（jwt 跳过验签）时 header 为空，本中间件什么也不做。
//   - 对 metainfo key 的命名约定见 iamclient/header.go（小写、与 HTTP header 镜像）。
package middleware

import "github.com/cloudwego/hertz/pkg/app"

// IdentityPropagationService 身份传播中间件接口
type IdentityPropagationService interface {
	MiddlewareFunc() app.HandlerFunc
}
