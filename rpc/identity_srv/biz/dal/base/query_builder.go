package base

import (
	"context"
	"fmt"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"gorm.io/gorm"
)

// QueryBuilder 查询构建器，提供流式 API 构建复杂查询
// 支持链式调用、条件组合、软删除处理等功能
//
// 设计目标：
// - 消除 query 和 countQuery 的重复代码
// - 提供类型安全的泛型查询接口
// - 支持自定义搜索逻辑注入
// - 集成 base repository 的软删除机制
type QueryBuilder[T any] struct {
	db       *gorm.DB
	ctx      context.Context
	baseRepo *BaseRepositoryImpl[T]

	// 条件累积器：存储所有查询条件的闭包
	conditions []func(*gorm.DB) *gorm.DB

	// 自定义搜索逻辑（可选）
	searchFunc func(*gorm.DB, string) *gorm.DB

	// 预加载关联
	preloads []string

	// 排序选项
	orderBy   string
	orderDesc bool
}

// NewQueryBuilder 创建查询构建器实例
func (r *BaseRepositoryImpl[T]) NewQueryBuilder(ctx context.Context) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		db:         r.db,
		ctx:        ctx,
		baseRepo:   r,
		conditions: make([]func(*gorm.DB) *gorm.DB, 0),
		preloads:   make([]string, 0),
	}
}

// WhereEqual 添加精确匹配条件（仅当 value 非 nil 时生效）
//
// 示例：
//
//	qb.WhereEqual("user_id", conditions.UserID)  // 如果 UserID 为 nil，则跳过
//	qb.WhereEqual("status", &activeStatus)       // 自动解引用指针
func (qb *QueryBuilder[T]) WhereEqual(field string, value interface{}) *QueryBuilder[T] {
	// 处理 nil 指针（跳过条件）
	if isNilPointer(value) {
		return qb
	}

	// 解引用指针
	actualValue := dereferencePointer(value)

	qb.conditions = append(qb.conditions, func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", field), actualValue)
	})

	return qb
}

// WhereIn 添加 IN 条件
//
// 示例：
//
//	qb.WhereIn("id", []string{"id1", "id2", "id3"})
func (qb *QueryBuilder[T]) WhereIn(field string, values interface{}) *QueryBuilder[T] {
	if isNilPointer(values) {
		return qb
	}

	qb.conditions = append(qb.conditions, func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s IN ?", field), values)
	})

	return qb
}

// WhereCustom 添加自定义条件
//
// 示例：
//
//	qb.WhereCustom(func(db *gorm.DB) *gorm.DB {
//	    return db.Where("created_at > ?", startTime).Where("created_at < ?", endTime)
//	})
func (qb *QueryBuilder[T]) WhereCustom(conditionFunc func(*gorm.DB) *gorm.DB) *QueryBuilder[T] {
	if conditionFunc != nil {
		qb.conditions = append(qb.conditions, conditionFunc)
	}
	return qb
}

// WithSoftDelete 应用软删除过滤逻辑
// 使用 base repository 的 applySoftDeleteFilter 方法，确保逻辑一致
func (qb *QueryBuilder[T]) WithSoftDelete(opts *QueryOptions) *QueryBuilder[T] {
	if opts == nil {
		return qb
	}

	qb.conditions = append(qb.conditions, func(db *gorm.DB) *gorm.DB {
		return qb.baseRepo.applySoftDeleteFilter(db, opts)
	})

	return qb
}

// WithSearch 应用搜索条件（支持自定义搜索逻辑）
//
// 示例：
//
//	qb.WithSearch(opts.Search, func(db *gorm.DB, term string) *gorm.DB {
//	    return db.Where("username LIKE ? OR email LIKE ?", "%"+term+"%", "%"+term+"%")
//	})
func (qb *QueryBuilder[T]) WithSearch(searchTerm string, searchFunc func(*gorm.DB, string) *gorm.DB) *QueryBuilder[T] {
	if searchTerm != "" && searchFunc != nil {
		qb.searchFunc = searchFunc
		qb.conditions = append(qb.conditions, func(db *gorm.DB) *gorm.DB {
			return searchFunc(db, searchTerm)
		})
	}
	return qb
}

// WithPreload 添加预加载关联
//
// 示例：
//
//	qb.WithPreload("Organization", "Department")
func (qb *QueryBuilder[T]) WithPreload(associations ...string) *QueryBuilder[T] {
	qb.preloads = append(qb.preloads, associations...)
	return qb
}

// WithOrder 设置排序
func (qb *QueryBuilder[T]) WithOrder(opts *QueryOptions) *QueryBuilder[T] {
	if opts != nil && opts.OrderBy != "" {
		qb.orderBy = opts.OrderBy
		qb.orderDesc = opts.OrderDesc
	}
	return qb
}

// buildQuery 构建最终查询（包含所有条件、预加载、排序）
func (qb *QueryBuilder[T]) buildQuery() *gorm.DB {
	query := qb.db.WithContext(qb.ctx).Model(new(T))

	// 应用所有条件
	for _, condition := range qb.conditions {
		query = condition(query)
	}

	// 应用预加载
	for _, preload := range qb.preloads {
		query = query.Preload(preload)
	}

	// 应用排序
	if qb.orderBy != "" {
		orderClause := qb.orderBy
		if qb.orderDesc {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}
		query = query.Order(orderClause)
	}

	return query
}

// buildCountQuery 构建计数查询（仅包含过滤条件，不包含预加载和排序）
func (qb *QueryBuilder[T]) buildCountQuery() *gorm.DB {
	query := qb.db.WithContext(qb.ctx).Model(new(T))

	// 仅应用过滤条件（不包含预加载和排序，避免不必要的性能开销）
	for _, condition := range qb.conditions {
		query = condition(query)
	}

	return query
}

// FindWithPagination 执行分页查询
// 这是最常用的方法，自动处理计数和分页
func (qb *QueryBuilder[T]) FindWithPagination(opts *QueryOptions) ([]*T, *models.PageResult, error) {
	if opts == nil {
		opts = NewQueryOptions()
	}

	// 验证查询选项
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	// 计算总数（使用 buildCountQuery，避免不必要的 JOIN 和排序）
	var total int64
	countQuery := qb.buildCountQuery()
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// 执行查询（使用 buildQuery，包含预加载和排序）
	var entities []*T
	query := qb.buildQuery()

	if opts.FetchAll {
		// 不分页，返回所有数据
		if err := query.Find(&entities).Error; err != nil {
			return nil, nil, err
		}
	} else {
		// 分页查询
		offset := (opts.Page - 1) * opts.PageSize
		if err := query.Offset(int(offset)).Limit(int(opts.PageSize)).Find(&entities).Error; err != nil {
			return nil, nil, err
		}
	}

	// 构建分页结果
	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return entities, pageResult, nil
}

// Count 执行计数查询（不加载数据）
func (qb *QueryBuilder[T]) Count() (int64, error) {
	var count int64
	query := qb.buildCountQuery()
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Find 执行查询（不分页，返回所有匹配记录）
func (qb *QueryBuilder[T]) Find() ([]*T, error) {
	var entities []*T
	query := qb.buildQuery()
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// First 查询第一条记录
func (qb *QueryBuilder[T]) First() (*T, error) {
	var entity T
	query := qb.buildQuery()
	if err := query.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// ============================================================================
// 辅助函数
// ============================================================================

// isNilPointer 检查是否为 nil 指针
// 用于 WhereEqual 等方法跳过 nil 条件
func isNilPointer(value interface{}) bool {
	if value == nil {
		return true
	}

	// 使用类型断言检查常见的指针类型
	switch v := value.(type) {
	case *string:
		return v == nil
	case *int:
		return v == nil
	case *int32:
		return v == nil
	case *int64:
		return v == nil
	case *bool:
		return v == nil
	case *float32:
		return v == nil
	case *float64:
		return v == nil
	default:
		// 对于其他类型（如自定义枚举），假设非 nil（保守处理）
		return false
	}
}

// dereferencePointer 解引用指针
// 将指针类型的值转换为实际值，用于 GORM 查询
func dereferencePointer(value interface{}) interface{} {
	switch v := value.(type) {
	case *string:
		if v != nil {
			return *v
		}
	case *int:
		if v != nil {
			return *v
		}
	case *int32:
		if v != nil {
			return *v
		}
	case *int64:
		if v != nil {
			return *v
		}
	case *bool:
		if v != nil {
			return *v
		}
	case *float32:
		if v != nil {
			return *v
		}
	case *float64:
		if v != nil {
			return *v
		}
	default:
		// 其他类型直接返回（包括非指针类型和自定义类型）
		return value
	}
	return nil
}
