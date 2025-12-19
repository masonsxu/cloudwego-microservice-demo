/*
Package errors 提供了 Gateway 的统一错误处理基础设施。

该包实现了一套完整的错误处理体系，支持从 RPC 层到 HTTP 响应层的错误转换和处理。
主要职责包括：错误类型定义、错误码规范、HTTP状态码映射、错误响应生成等。

# 错误码规范

采用5位数字结构的错误码系统：A-BB-CCC

  - A (错误级别): 1=系统级错误, 2=业务级错误
  - BB (服务模块): 00=通用错误, 01=用户服务, 02=组织服务等
  - CCC (具体错误): 唯一标识特定错误条件

# 错误类型

核心错误类型 APIError 提供了统一的错误表示：

	type APIError struct {
	    code    int32   // 业务错误码
	    message string  // 错误描述
	}

# 错误处理流程

1. RPC错误处理：ProcessRPCError() 将 RPC 错误转换为 APIError
2. HTTP错误处理：HandleHTTPStatusError() 处理 HTTP 状态码错误
3. 统一响应：AbortWithError() 生成标准化的错误响应
4. 服务层辅助：HandleServiceError() 简化服务层错误处理

# 预定义错误

包提供了常用错误的预定义变量：

	var (
	    ErrInvalidParams = NewAPIError(100001, "参数错误")
	    ErrUnauthorized  = NewAPIError(100002, "未授权访问")
	    ErrForbidden     = NewAPIError(100003, "权限不足")
	    ErrNotFound      = NewAPIError(100004, "资源不存在")
	    ErrInternal      = NewAPIError(100005, "系统繁忙，请稍后重试")
	    // JWT认证错误...
	    // 网关特有错误...
	)

# 使用示例

在处理程序中使用：

	func (h *Handler) SomeAPI(ctx context.Context, c *app.RequestContext, req *SomeRequest) {
	    resp, err := h.service.DoSomething(ctx, req)
	    if err != nil {
	        errors.HandleServiceError(c, err, "操作失败")
	        return
	    }
	    // 处理成功响应...
	}

在服务层处理 RPC 错误：

	func (s *Service) CallRPC(ctx context.Context) error {
	    resp, err := s.rpcClient.SomeMethod(ctx, req)
	    return errors.ProcessRPCError(err, "RPC调用失败")
	}

# HTTP状态码映射

错误码自动映射到相应的HTTP状态码：

  - 100001-100004: 4xx 客户端错误
  - 100005: 500 服务器内部错误
  - 101xxx: 401 认证相关错误
  - 110xxx: 5xx 网关特有错误
  - 20000-29999: 200 (RPC业务错误，HTTP层成功)

# 响应格式

所有错误响应都遵循统一的JSON格式：

	{
	    "code": 100001,
	    "message": "参数错误",
	    "request_id": "01234567-89ab-cdef",
	    "trace_id": "trace-12345",
	    "timestamp": 1640995200000
	}

该格式确保了前端能够统一处理所有类型的错误响应。
*/
package errors
