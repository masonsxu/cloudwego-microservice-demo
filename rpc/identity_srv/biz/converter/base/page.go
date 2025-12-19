package base

import (
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// Converter 基础转换器接口，负责通用类型的转换
type Converter interface {
	// 分页响应转换：从 models PageResult 到 Thrift PageResponse
	PageResponseToThrift(*models.PageResult) *rpc_base.PageResponse

	// 分页请求转换：从 Thrift PageRequest 到 QueryOptions（完整转换，支持搜索、过滤等高级特性）
	PageRequestToQueryOptions(*rpc_base.PageRequest) *base.QueryOptions
}
