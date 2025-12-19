package base

import (
	"strings"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

// ConverterImpl 基础转换器实现，负责通用类型的转换
// 支持丰富的分页参数转换和排序逻辑处理
type ConverterImpl struct{}

// NewConverter 创建一个新的基础转换器实例
func NewConverter() Converter {
	return &ConverterImpl{}
}

// PageResponseToThrift 将 models 分页结果转换为 Thrift 分页响应
func (c *ConverterImpl) PageResponseToThrift(p *models.PageResult) *rpc_base.PageResponse {
	if p == nil {
		return &rpc_base.PageResponse{
			Total:      &[]int32{0}[0],
			Page:       &[]int32{1}[0],
			Limit:      &[]int32{20}[0],
			TotalPages: &[]int32{0}[0],
			HasNext:    &[]bool{false}[0],
			HasPrev:    &[]bool{false}[0],
		}
	}

	return &rpc_base.PageResponse{
		Total:      &p.Total,
		Page:       &p.Page,
		Limit:      &p.Limit,
		TotalPages: &p.TotalPages,
		HasNext:    &p.HasNext,
		HasPrev:    &p.HasPrev,
	}
}

// PageRequestToQueryOptions 将 Thrift 分页请求转换为 QueryOptions
// 这是一个更完整的转换，支持搜索、过滤等高级特性
func (c *ConverterImpl) PageRequestToQueryOptions(
	req *rpc_base.PageRequest,
) *base.QueryOptions {
	opts := base.NewQueryOptions()

	if req == nil {
		return opts
	}

	// 设置分页参数
	page := req.GetPage()
	limit := req.GetLimit()
	opts.WithPage(page, limit)

	// 设置搜索关键词
	if req.Search != nil && *req.Search != "" {
		opts.WithSearch(strings.TrimSpace(*req.Search))
	}

	// 解析排序字段
	if req.Sort != nil && *req.Sort != "" {
		sortFields := strings.Split(strings.TrimSpace(*req.Sort), ",")
		// 取第一个排序字段作为主排序
		if len(sortFields) > 0 {
			sortField := strings.TrimSpace(sortFields[0])
			// 处理带 - 前缀的降序排序
			if field, found := strings.CutPrefix(sortField, "-"); found {
				opts.WithOrder(field, true)
			} else {
				opts.WithOrder(sortField, false)
			}
		}
	}

	// 处理过滤条件
	if req.Filter != nil {
		for key, value := range req.Filter {
			if strings.TrimSpace(value) != "" {
				opts.WithFilter(key, strings.TrimSpace(value))
			}
		}
	}

	// 处理 FetchAll 参数
	if req.FetchAll != nil && *req.FetchAll {
		opts.WithFetchAll(true)
	}

	return opts
}
