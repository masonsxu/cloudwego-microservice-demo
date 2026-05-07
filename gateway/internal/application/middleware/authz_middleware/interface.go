// Package middleware 提供网关路由级 ACL（粗粒度授权）中间件实现。
//
// 仅支持「path 前缀 + 角色 OR」风格的策略，不持有领域字段（部门、数据范围等）。
// 任何超出此粒度的需求请走 PDP（policy_srv）。
package middleware

import "github.com/cloudwego/hertz/pkg/app"

// AuthZMiddlewareService 路由级 ACL 中间件接口
type AuthZMiddlewareService interface {
	MiddlewareFunc() app.HandlerFunc
}
