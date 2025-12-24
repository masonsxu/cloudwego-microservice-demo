package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/biz/dal/base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/kitex_gen/rpc_base"
	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
)

func TestNewConverter(t *testing.T) {
	converter := NewConverter()
	assert.NotNil(t, converter)
	assert.IsType(t, &ConverterImpl{}, converter)
}

func TestConverterImpl_PageResponseToThrift(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.PageResponseToThrift(nil)

		require.NotNil(t, result)
		assert.Equal(t, int32(0), *result.Total)
		assert.Equal(t, int32(1), *result.Page)
		assert.Equal(t, int32(20), *result.Limit)
		assert.Equal(t, int32(0), *result.TotalPages)
		assert.Equal(t, false, *result.HasNext)
		assert.Equal(t, false, *result.HasPrev)
	})

	t.Run("正常分页结果转换", func(t *testing.T) {
		pageResult := &models.PageResult{
			Total:      100,
			Page:       2,
			Limit:      10,
			TotalPages: 10,
			HasNext:    true,
			HasPrev:    true,
		}

		result := converter.PageResponseToThrift(pageResult)

		require.NotNil(t, result)
		assert.Equal(t, int32(100), *result.Total)
		assert.Equal(t, int32(2), *result.Page)
		assert.Equal(t, int32(10), *result.Limit)
		assert.Equal(t, int32(10), *result.TotalPages)
		assert.Equal(t, true, *result.HasNext)
		assert.Equal(t, true, *result.HasPrev)
	})

	t.Run("第一页分页结果转换", func(t *testing.T) {
		pageResult := &models.PageResult{
			Total:      25,
			Page:       1,
			Limit:      10,
			TotalPages: 3,
			HasNext:    true,
			HasPrev:    false,
		}

		result := converter.PageResponseToThrift(pageResult)

		require.NotNil(t, result)
		assert.Equal(t, int32(25), *result.Total)
		assert.Equal(t, int32(1), *result.Page)
		assert.Equal(t, int32(10), *result.Limit)
		assert.Equal(t, int32(3), *result.TotalPages)
		assert.Equal(t, true, *result.HasNext)
		assert.Equal(t, false, *result.HasPrev)
	})

	t.Run("最后一页分页结果转换", func(t *testing.T) {
		pageResult := &models.PageResult{
			Total:      25,
			Page:       3,
			Limit:      10,
			TotalPages: 3,
			HasNext:    false,
			HasPrev:    true,
		}

		result := converter.PageResponseToThrift(pageResult)

		require.NotNil(t, result)
		assert.Equal(t, int32(25), *result.Total)
		assert.Equal(t, int32(3), *result.Page)
		assert.Equal(t, int32(10), *result.Limit)
		assert.Equal(t, int32(3), *result.TotalPages)
		assert.Equal(t, false, *result.HasNext)
		assert.Equal(t, true, *result.HasPrev)
	})

	t.Run("空结果分页转换", func(t *testing.T) {
		pageResult := &models.PageResult{
			Total:      0,
			Page:       1,
			Limit:      20,
			TotalPages: 0,
			HasNext:    false,
			HasPrev:    false,
		}

		result := converter.PageResponseToThrift(pageResult)

		require.NotNil(t, result)
		assert.Equal(t, int32(0), *result.Total)
		assert.Equal(t, int32(1), *result.Page)
		assert.Equal(t, int32(20), *result.Limit)
		assert.Equal(t, int32(0), *result.TotalPages)
		assert.Equal(t, false, *result.HasNext)
		assert.Equal(t, false, *result.HasPrev)
	})

	t.Run("单页结果分页转换", func(t *testing.T) {
		pageResult := &models.PageResult{
			Total:      5,
			Page:       1,
			Limit:      20,
			TotalPages: 1,
			HasNext:    false,
			HasPrev:    false,
		}

		result := converter.PageResponseToThrift(pageResult)

		require.NotNil(t, result)
		assert.Equal(t, int32(5), *result.Total)
		assert.Equal(t, int32(1), *result.Page)
		assert.Equal(t, int32(20), *result.Limit)
		assert.Equal(t, int32(1), *result.TotalPages)
		assert.Equal(t, false, *result.HasNext)
		assert.Equal(t, false, *result.HasPrev)
	})
}

func TestConverterImpl_PageRequestToQueryOptions(t *testing.T) {
	converter := &ConverterImpl{}

	t.Run("nil输入转换", func(t *testing.T) {
		result := converter.PageRequestToQueryOptions(nil)

		require.NotNil(t, result)
		// 验证默认值
		assert.Equal(t, int32(1), result.Page)
		assert.Equal(t, int32(20), result.PageSize)
		assert.Equal(t, "created_at", result.OrderBy)
		assert.Equal(t, true, result.OrderDesc)
		assert.Equal(t, "", result.Search)
		assert.Empty(t, result.Filters)
		assert.Equal(t, false, result.FetchAll)
	})

	t.Run("基础分页参数转换", func(t *testing.T) {
		page := int32(2)
		limit := int32(15)
		req := &rpc_base.PageRequest{
			Page:  page,
			Limit: limit,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, int32(2), result.Page)
		assert.Equal(t, int32(15), result.PageSize)
	})

	t.Run("搜索关键词转换", func(t *testing.T) {
		search := "  test keyword  "
		req := &rpc_base.PageRequest{
			Search: &search,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "test keyword", result.Search)
	})

	t.Run("空搜索关键词处理", func(t *testing.T) {
		search := ""
		req := &rpc_base.PageRequest{
			Search: &search,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "", result.Search)
	})

	t.Run("空格搜索关键词处理", func(t *testing.T) {
		search := "   "
		req := &rpc_base.PageRequest{
			Search: &search,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "", result.Search)
	})

	t.Run("升序排序转换", func(t *testing.T) {
		sort := "name"
		req := &rpc_base.PageRequest{
			Sort: &sort,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "name", result.OrderBy)
		assert.Equal(t, false, result.OrderDesc)
	})

	t.Run("降序排序转换", func(t *testing.T) {
		sort := "-created_at"
		req := &rpc_base.PageRequest{
			Sort: &sort,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "created_at", result.OrderBy)
		assert.Equal(t, true, result.OrderDesc)
	})

	t.Run("多字段排序处理（取第一个）", func(t *testing.T) {
		sort := "name,-created_at,id"
		req := &rpc_base.PageRequest{
			Sort: &sort,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "name", result.OrderBy)
		assert.Equal(t, false, result.OrderDesc)
	})

	t.Run("带空格的排序字段处理", func(t *testing.T) {
		sort := "  -created_at  "
		req := &rpc_base.PageRequest{
			Sort: &sort,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "created_at", result.OrderBy)
		assert.Equal(t, true, result.OrderDesc)
	})

	t.Run("空排序字段处理", func(t *testing.T) {
		sort := ""
		req := &rpc_base.PageRequest{
			Sort: &sort,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, "created_at", result.OrderBy) // 保持默认值
		assert.Equal(t, true, result.OrderDesc)       // 保持默认值
	})

	t.Run("过滤条件转换", func(t *testing.T) {
		filter := map[string]string{
			"status": "active",
			"type":   "user",
		}
		req := &rpc_base.PageRequest{
			Filter: filter,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Len(t, result.Filters, 2)
		assert.Equal(t, "active", result.Filters["status"])
		assert.Equal(t, "user", result.Filters["type"])
	})

	t.Run("空值过滤条件处理", func(t *testing.T) {
		filter := map[string]string{
			"status": "active",
			"type":   "",      // 空值应该被忽略
			"empty":  "   ",   // 空格应该被忽略
			"name":   "test",
		}
		req := &rpc_base.PageRequest{
			Filter: filter,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Len(t, result.Filters, 2) // 只有非空值被保留
		assert.Equal(t, "active", result.Filters["status"])
		assert.Equal(t, "test", result.Filters["name"])
		assert.NotContains(t, result.Filters, "type")
		assert.NotContains(t, result.Filters, "empty")
	})

	t.Run("空过滤条件处理", func(t *testing.T) {
		filter := map[string]string{}
		req := &rpc_base.PageRequest{
			Filter: filter,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Empty(t, result.Filters)
	})

	t.Run("FetchAll参数转换", func(t *testing.T) {
		fetchAll := true
		req := &rpc_base.PageRequest{
			FetchAll: &fetchAll,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, true, result.FetchAll)
	})

	t.Run("FetchAll为false的处理", func(t *testing.T) {
		fetchAll := false
		req := &rpc_base.PageRequest{
			FetchAll: &fetchAll,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, false, result.FetchAll)
	})

	t.Run("完整参数转换", func(t *testing.T) {
		page := int32(3)
		limit := int32(25)
		search := "john"
		sort := "-created_at"
		filter := map[string]string{
			"status": "active",
			"role":   "admin",
		}
		fetchAll := false

		req := &rpc_base.PageRequest{
			Page:    page,
			Limit:   limit,
			Search:  &search,
			Sort:    &sort,
			Filter:  filter,
			FetchAll: &fetchAll,
		}

		result := converter.PageRequestToQueryOptions(req)

		require.NotNil(t, result)
		assert.Equal(t, int32(3), result.Page)
		assert.Equal(t, int32(25), result.PageSize)
		assert.Equal(t, "john", result.Search)
		assert.Equal(t, "created_at", result.OrderBy)
		assert.Equal(t, true, result.OrderDesc)
		assert.Len(t, result.Filters, 2)
		assert.Equal(t, "active", result.Filters["status"])
		assert.Equal(t, "admin", result.Filters["role"])
		assert.Equal(t, false, result.FetchAll)
	})
}

// 表格驱动测试
func TestConverterImpl_PageRequestToQueryOptions_TableDriven(t *testing.T) {
	converter := &ConverterImpl{}

	tests := []struct {
		name        string
		input       *rpc_base.PageRequest
		expected    *base.QueryOptions
		description string
	}{
		{
			name:        "nil输入",
			input:       nil,
			expected:    base.NewQueryOptions(),
			description: "nil输入应该返回默认QueryOptions",
		},
		{
			name: "最小参数",
			input: &rpc_base.PageRequest{},
			expected: func() *base.QueryOptions {
				opts := base.NewQueryOptions()
				return opts
			}(),
			description: "空PageRequest应该返回默认QueryOptions",
		},
		{
			name: "仅分页参数",
			input: func() *rpc_base.PageRequest {
				page := int32(2)
				limit := int32(10)
				return &rpc_base.PageRequest{
					Page:  page,
					Limit: limit,
				}
			}(),
			expected: func() *base.QueryOptions {
				opts := base.NewQueryOptions()
				opts.WithPage(2, 10)
				return opts
			}(),
			description: "仅设置分页参数",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.PageRequestToQueryOptions(tt.input)
			assert.Equal(t, tt.expected.Page, result.Page, tt.description)
			assert.Equal(t, tt.expected.PageSize, result.PageSize, tt.description)
		})
	}
}

// 基准测试
func BenchmarkConverterImpl_PageResponseToThrift(b *testing.B) {
	converter := &ConverterImpl{}
	pageResult := &models.PageResult{
		Total:      1000,
		Page:       5,
		Limit:      20,
		TotalPages: 50,
		HasNext:    true,
		HasPrev:    true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.PageResponseToThrift(pageResult)
	}
}

func BenchmarkConverterImpl_PageRequestToQueryOptions(b *testing.B) {
	converter := &ConverterImpl{}
	page := int32(5)
	limit := int32(20)
	search := "test search"
	sort := "-created_at"
	filter := map[string]string{
		"status": "active",
		"type":   "user",
	}
	fetchAll := false

	req := &rpc_base.PageRequest{
		Page:    page,
		Limit:   limit,
		Search:  &search,
		Sort:    &sort,
		Filter:  filter,
		FetchAll: &fetchAll,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.PageRequestToQueryOptions(req)
	}
}