package models

// PageOptions 通用分页入参，供 DAL/Service 使用
// 注意：这是领域内的简单值对象，与 Thrift 的 PageRequest 通过 converter 适配
// Page 从 1 开始，Limit > 0
type PageOptions struct {
	Page     int32  `json:"page"`      // 页码，从 1 开始
	Limit    int32  `json:"limit"`     // 每页大小
	SortBy   string `json:"sort_by"`   // 排序字段
	SortDesc bool   `json:"sort_desc"` // 是否降序
}

// PageResult 分页结果元信息
type PageResult struct {
	Total      int32 `json:"total"`       // 总记录数
	Page       int32 `json:"page"`        // 当前页码
	Limit      int32 `json:"limit"`       // 每页大小
	TotalPages int32 `json:"total_pages"` // 总页数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
	HasPrev    bool  `json:"has_prev"`    // 是否有上一页
}

// NewPageOptions 创建分页选项
func NewPageOptions(page, limit int32, sortBy string, sortDesc bool) *PageOptions {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100 // 防止过大的分页请求
	}

	if sortBy == "" {
		sortBy = "created_at"
	}

	return &PageOptions{
		Page:     page,
		Limit:    limit,
		SortBy:   sortBy,
		SortDesc: sortDesc,
	}
}

// GetOffset 计算数据库查询的 offset
func (p *PageOptions) GetOffset() int32 {
	return (p.Page - 1) * p.Limit
}

// NewPageResult 创建分页结果
func NewPageResult(total, page, limit int32) *PageResult {
	totalPages := (total + limit - 1) / limit // 向上取整
	if totalPages == 0 {
		totalPages = 1
	}

	return &PageResult{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
