// Package errors 提供了 Hertz 网关层的统一错误处理机制
// 本文件提供统一的 JSON 响应发送功能
package errors

import (
	"github.com/cloudwego/hertz/pkg/app"
)

// JSON 发送 JSON 响应
// 注意：
// - RequestID 通过 HTTP Header (X-Request-ID) 传递，由 requestid 中间件自动处理
// - Date 响应头由 ResponseHeaderMiddleware 自动添加
//
// 使用示例:
//
//	resp, err := service.GetData(ctx, req)
//	if err != nil {
//	    errors.HandleServiceError(c, err, "获取数据失败")
//	    return
//	}
//	errors.JSON(c, consts.StatusOK, resp)
func JSON(c *app.RequestContext, code int, obj interface{}) {
	c.JSON(code, obj)
}
