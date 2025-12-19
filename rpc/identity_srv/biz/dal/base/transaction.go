package base

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TransactionManager 事务管理器接口
type TransactionManager interface {
	// WithTransaction 在事务中执行操作
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error

	// BeginTx 开始事务
	BeginTx(ctx context.Context) (*gorm.DB, error)

	// CommitTx 提交事务
	CommitTx(tx *gorm.DB) error

	// RollbackTx 回滚事务
	RollbackTx(tx *gorm.DB) error
}

// TransactionManagerImpl 事务管理器实现
type TransactionManagerImpl struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &TransactionManagerImpl{
		db: db,
	}
}

// WithTransaction 在事务中执行操作（推荐使用）
func (tm *TransactionManagerImpl) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context, tx *gorm.DB) error,
) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

// BeginTx 开始事务
func (tm *TransactionManagerImpl) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("开启事务失败: %w", tx.Error)
	}

	return tx, nil
}

// CommitTx 提交事务
func (tm *TransactionManagerImpl) CommitTx(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// RollbackTx 回滚事务
func (tm *TransactionManagerImpl) RollbackTx(tx *gorm.DB) error {
	if err := tx.Rollback().Error; err != nil {
		return fmt.Errorf("回滚事务失败: %w", err)
	}

	return nil
}
