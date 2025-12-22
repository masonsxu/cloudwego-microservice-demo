package base

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/masonsxu/cloudwego-microservice-demo/rpc/identity-srv/models"
	"gorm.io/gorm"
)

// ============================================================================
// 错误定义
// ============================================================================

// ErrConflictingQueryOptions 查询选项冲突错误
var ErrConflictingQueryOptions = errors.New(
	"conflicting query options: DeletedOnly and IncludeDeleted cannot both be true",
)

// ============================================================================
// Schema 缓存机制
// ============================================================================

// modelSchema 模型元信息缓存
// 存储模型的软删除相关信息，避免重复反射解析
type modelSchema struct {
	hasSoftDelete   bool   // 是否支持软删除
	deletedAtColumn string // DeletedAt 字段的实际列名
	tableName       string // 表名
}

// modelSchemaCache 全局模型 Schema 缓存
// Key: reflect.Type, Value: *modelSchema
var modelSchemaCache sync.Map

// parseModelSchema 解析模型 Schema 信息并缓存
// 使用 GORM 的 Schema API 动态获取字段信息
func parseModelSchema[T any](db *gorm.DB) *modelSchema {
	var model T

	modelType := reflect.TypeOf(model)

	// 处理指针类型
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 尝试从缓存获取
	if cached, ok := modelSchemaCache.Load(modelType); ok {
		return cached.(*modelSchema)
	}

	// 使用 GORM Schema Parse API 解析模型
	s := &modelSchema{
		hasSoftDelete:   false,
		deletedAtColumn: "",
		tableName:       "",
	}

	// 解析 Schema
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(&model); err == nil {
		s.tableName = stmt.Schema.Table

		// 查找 DeletedAt 字段
		for _, field := range stmt.Schema.Fields {
			if field.Name == "DeletedAt" {
				s.hasSoftDelete = true
				s.deletedAtColumn = field.DBName

				break
			}
		}
	} else {
		// 如果 Schema 解析失败，回退到反射检查
		if modelType.Kind() == reflect.Struct {
			if _, ok := modelType.FieldByName("DeletedAt"); ok {
				s.hasSoftDelete = true
				s.deletedAtColumn = "deleted_at" // 使用默认列名
			}
		}
	}

	// 存储到缓存
	modelSchemaCache.Store(modelType, s)

	return s
}

// BaseRepository 基础仓储接口，定义所有实体的通用操作
// 采用泛型设计，提供类型安全的CRUD操作
type BaseRepository[T any] interface {
	// 基础CRUD操作
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error

	// 软删除操作
	SoftDelete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error

	// 批量操作
	BatchCreate(ctx context.Context, entities []*T) error
	BatchGetByIDs(ctx context.Context, ids []string) ([]*T, error)
	BatchUpdate(ctx context.Context, entities []*T) error
	BatchDelete(ctx context.Context, ids []string) error
	BatchSoftDelete(ctx context.Context, ids []string) error
	BatchHardDelete(ctx context.Context, ids []string) error
	BatchRestore(ctx context.Context, ids []string) error

	// 查询操作
	FindAll(ctx context.Context, opts *QueryOptions) ([]*T, *models.PageResult, error)
	Count(ctx context.Context, opts *QueryOptions) (int64, error)
	Exists(ctx context.Context, id string) (bool, error)

	// 事务操作
	WithTx(tx *gorm.DB) BaseRepository[T]
}

// QueryOptions 查询选项配置
type QueryOptions struct {
	// 分页参数
	Page      int32  `json:"page"`       // 页码，从1开始
	PageSize  int32  `json:"page_size"`  // 每页大小
	OrderBy   string `json:"order_by"`   // 排序字段
	OrderDesc bool   `json:"order_desc"` // 是否降序
	FetchAll  bool   `json:"fetch_all"`  // 是否获取所有数据（不分页）

	// 过滤条件
	Filters map[string]interface{} `json:"filters"` // 字段过滤条件
	Search  string                 `json:"search"`  // 全文搜索关键词

	// 软删除控制
	IncludeDeleted bool `json:"include_deleted"` // 是否包含已删除记录
	DeletedOnly    bool `json:"deleted_only"`    // 仅查询已删除记录

	// 预加载关联
	Preloads []string `json:"preloads"` // 需要预加载的关联字段
}

// NewQueryOptions 创建默认查询选项
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Page:      1,
		PageSize:  20,
		OrderBy:   "created_at",
		OrderDesc: true,
		Filters:   make(map[string]interface{}),
		Preloads:  make([]string, 0),
	}
}

// WithPage 设置分页参数
func (opts *QueryOptions) WithPage(page, pageSize int32) *QueryOptions {
	// 参数验证和边界检查
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20 // 默认每页20条
	}

	if pageSize > 100 {
		pageSize = 100 // 限制最大每页100条，防止性能问题
	}

	opts.Page = page
	opts.PageSize = pageSize

	return opts
}

// WithOrder 设置排序参数
func (opts *QueryOptions) WithOrder(orderBy string, desc bool) *QueryOptions {
	// 基础安全验证，防止SQL注入风险
	if orderBy != "" {
		// 简单的字段名验证：只允许字母、数字、下划线和点号
		validField := true

		for _, char := range orderBy {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '_' || char == '.') {
				validField = false
				break
			}
		}

		if validField {
			opts.OrderBy = orderBy
			opts.OrderDesc = desc
		}
		// 如果字段名无效，保持原有排序设置
	}

	return opts
}

// WithFilter 添加过滤条件
func (opts *QueryOptions) WithFilter(field string, value interface{}) *QueryOptions {
	opts.Filters[field] = value
	return opts
}

// WithSearch 设置搜索关键词
func (opts *QueryOptions) WithSearch(searchTerm string) *QueryOptions {
	// 清理和验证搜索关键词
	searchTerm = strings.TrimSpace(searchTerm)

	// 限制搜索关键词长度，防止过长查询影响性能
	if len(searchTerm) > 100 {
		searchTerm = searchTerm[:100]
	}

	opts.Search = searchTerm

	return opts
}

// WithPreload 添加预加载关联
func (opts *QueryOptions) WithPreload(associations ...string) *QueryOptions {
	opts.Preloads = append(opts.Preloads, associations...)
	return opts
}

// WithFetchAll 设置是否获取所有数据（不分页）
func (opts *QueryOptions) WithFetchAll(fetchAll bool) *QueryOptions {
	opts.FetchAll = fetchAll
	return opts
}

// Validate 验证查询选项的有效性
// 检查参数冲突并返回错误
func (opts *QueryOptions) Validate() error {
	// 检查软删除选项冲突
	if opts.DeletedOnly && opts.IncludeDeleted {
		return ErrConflictingQueryOptions
	}

	return nil
}

// BaseRepositoryImpl 基础仓储实现，提供通用的GORM操作
type BaseRepositoryImpl[T any] struct {
	db          *gorm.DB
	modelSchema *modelSchema // 缓存的模型 Schema 信息
}

// NewBaseRepository 创建基础仓储实例
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{
		db:          db,
		modelSchema: parseModelSchema[T](db), // 初始化时解析并缓存 Schema
	}
}

// Create 创建实体
func (r *BaseRepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return err
	}

	return nil
}

// GetByID 根据ID获取实体
func (r *BaseRepositoryImpl[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, err
	}

	return &entity, nil
}

// Update 更新实体
func (r *BaseRepositoryImpl[T]) Update(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return err
	}

	return nil
}

// Delete 软删除实体
func (r *BaseRepositoryImpl[T]) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(new(T)).Error; err != nil {
		return err
	}

	return nil
}

// BatchCreate 批量创建实体
func (r *BaseRepositoryImpl[T]) BatchCreate(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		return err
	}

	return nil
}

// BatchGetByIDs 根据ID列表批量获取实体
func (r *BaseRepositoryImpl[T]) BatchGetByIDs(ctx context.Context, ids []string) ([]*T, error) {
	if len(ids) == 0 {
		return make([]*T, 0), nil
	}

	var entities []*T
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

// BatchUpdate 批量更新实体
func (r *BaseRepositoryImpl[T]) BatchUpdate(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}

	// GORM 不直接支持批量更新，使用事务逐个更新
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := tx.Save(entity).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// BatchDelete 批量软删除实体
func (r *BaseRepositoryImpl[T]) BatchDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(new(T)).Error; err != nil {
		return err
	}

	return nil
}

// FindAll 查询所有实体（支持分页和过滤）
func (r *BaseRepositoryImpl[T]) FindAll(
	ctx context.Context,
	opts *QueryOptions,
) ([]*T, *models.PageResult, error) {
	if opts == nil {
		opts = NewQueryOptions()
	}

	// 验证查询选项
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	// 构建查询
	query := r.buildQuery(ctx, opts)

	// 计算总数
	var total int64

	countQuery := r.buildCountQuery(ctx, opts)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// 分页查询
	var entities []*T

	// 如果 FetchAll 为 true，则不分页，直接返回所有数据
	if opts.FetchAll {
		if err := query.Find(&entities).Error; err != nil {
			return nil, nil, err
		}
	} else {
		offset := (opts.Page - 1) * opts.PageSize
		if err := query.Offset(int(offset)).Limit(int(opts.PageSize)).Find(&entities).Error; err != nil {
			return nil, nil, err
		}
	}

	// 构建分页结果
	pageResult := models.NewPageResult(int32(total), opts.Page, opts.PageSize)

	return entities, pageResult, nil
}

// Count 统计实体数量
func (r *BaseRepositoryImpl[T]) Count(ctx context.Context, opts *QueryOptions) (int64, error) {
	if opts == nil {
		opts = NewQueryOptions()
	}

	// 验证查询选项
	if err := opts.Validate(); err != nil {
		return 0, err
	}

	var count int64

	query := r.buildCountQuery(ctx, opts)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Exists 检查实体是否存在
func (r *BaseRepositoryImpl[T]) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// WithTx 使用指定事务创建新的仓储实例
func (r *BaseRepositoryImpl[T]) WithTx(tx *gorm.DB) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{
		db: tx,
	}
}

// buildQuery 构建查询语句
func (r *BaseRepositoryImpl[T]) buildQuery(ctx context.Context, opts *QueryOptions) *gorm.DB {
	query := r.db.WithContext(ctx).Model(new(T))

	// 处理软删除逻辑
	query = r.applySoftDeleteFilter(query, opts)

	// 处理过滤条件
	for field, value := range opts.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// 处理搜索（需要子类实现具体搜索逻辑）
	if opts.Search != "" {
		// 基础实现：按ID搜索，具体搜索逻辑由子类覆盖
		query = query.Where("id LIKE ?", "%"+opts.Search+"%")
	}

	// 处理预加载
	for _, preload := range opts.Preloads {
		query = query.Preload(preload)
	}

	// 处理排序
	if opts.OrderBy != "" {
		orderClause := opts.OrderBy
		if opts.OrderDesc {
			orderClause += " DESC"
		} else {
			orderClause += " ASC"
		}

		query = query.Order(orderClause)
	}

	return query
}

// buildCountQuery 构建计数查询语句
func (r *BaseRepositoryImpl[T]) buildCountQuery(ctx context.Context, opts *QueryOptions) *gorm.DB {
	query := r.db.WithContext(ctx).Model(new(T))

	// 处理软删除逻辑
	query = r.applySoftDeleteFilter(query, opts)

	// 处理过滤条件
	for field, value := range opts.Filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// 处理搜索
	if opts.Search != "" {
		query = query.Where("id LIKE ?", "%"+opts.Search+"%")
	}

	return query
}

// applySoftDeleteFilter 应用软删除过滤逻辑
func (r *BaseRepositoryImpl[T]) applySoftDeleteFilter(query *gorm.DB, opts *QueryOptions) *gorm.DB {
	if opts.DeletedOnly {
		// 仅查询已删除记录
		deletedAtCol := r.getDeletedAtColumn()
		return query.Unscoped().Where(fmt.Sprintf("%s IS NOT NULL", deletedAtCol))
	} else if opts.IncludeDeleted {
		// 查询所有记录（包含已删除）
		return query.Unscoped()
	} else {
		// 默认：仅查询未删除记录（GORM 自动处理 DeletedAt 字段）
		return query
	}
}

// hasSoftDelete 检查模型是否支持软删除
// 使用缓存的 Schema 信息，避免重复反射
func (r *BaseRepositoryImpl[T]) hasSoftDelete() bool {
	return r.modelSchema.hasSoftDelete
}

// getDeletedAtColumn 获取 DeletedAt 字段的列名
func (r *BaseRepositoryImpl[T]) getDeletedAtColumn() string {
	return r.modelSchema.deletedAtColumn
}

// SoftDelete 软删除实体（物理存在，逻辑删除）
func (r *BaseRepositoryImpl[T]) SoftDelete(ctx context.Context, id string) error {
	if !r.hasSoftDelete() {
		// 如果不支持软删除，则执行物理删除
		return r.Delete(ctx, id)
	}

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(new(T))
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// HardDelete 物理删除实体（从数据库中彻底删除）
func (r *BaseRepositoryImpl[T]) HardDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(new(T))
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Restore 恢复已软删除的实体
func (r *BaseRepositoryImpl[T]) Restore(ctx context.Context, id string) error {
	if !r.hasSoftDelete() {
		return gorm.ErrUnsupportedDriver
	}

	// 使用 Unscoped 查找已删除的记录并恢复
	deletedAtCol := r.getDeletedAtColumn()
	result := r.db.WithContext(ctx).Unscoped().Model(new(T)).
		Where(fmt.Sprintf("id = ? AND %s IS NOT NULL", deletedAtCol), id).
		Update(deletedAtCol, nil)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// BatchSoftDelete 批量软删除实体
func (r *BaseRepositoryImpl[T]) BatchSoftDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	if !r.hasSoftDelete() {
		// 如果不支持软删除，则执行物理删除
		return r.BatchDelete(ctx, ids)
	}

	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(new(T))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// BatchHardDelete 批量物理删除实体
func (r *BaseRepositoryImpl[T]) BatchHardDelete(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Delete(new(T))
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// BatchRestore 批量恢复已软删除的实体
func (r *BaseRepositoryImpl[T]) BatchRestore(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	if !r.hasSoftDelete() {
		return gorm.ErrUnsupportedDriver
	}

	deletedAtCol := r.getDeletedAtColumn()
	result := r.db.WithContext(ctx).Unscoped().Model(new(T)).
		Where(fmt.Sprintf("id IN ? AND %s IS NOT NULL", deletedAtCol), ids).
		Update(deletedAtCol, nil)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
